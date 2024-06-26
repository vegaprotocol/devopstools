package snapshotcompatibility

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/vegaapi/core"
	"github.com/vegaprotocol/devopstools/vegacapsule"

	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ProduceNewSnapshotArgs struct {
	*Args

	VegacapsuleHome   string
	VegacapsuleBinary string
}

var produceNewSnapshotArgs ProduceNewSnapshotArgs

var produceNewSnapshotCmd = &cobra.Command{
	Use:   "produce-new-snapshot",
	Short: "Produces new snapshot for the snapshot-compatibility network",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runProduceNewSnapshot(produceNewSnapshotArgs.Logger,
			produceNewSnapshotArgs.VegacapsuleBinary,
			produceNewSnapshotArgs.VegacapsuleHome,
		); err != nil {
			produceNewSnapshotArgs.Logger.Fatal("failed to produce new snapshot", zap.Error(err))
		}
	},
}

func init() {
	produceNewSnapshotArgs.Args = &snapshotCompatibilityArgs
	produceNewSnapshotCmd.PersistentFlags().
		StringVar(&produceNewSnapshotArgs.VegacapsuleHome, "vegacapsule-home", "", "The custom vegacapsule home")
	produceNewSnapshotCmd.PersistentFlags().
		StringVar(&produceNewSnapshotArgs.VegacapsuleBinary, "vegacapsule-binary", "vegacapsule", "The vegacapsule binary path")
}

func runProduceNewSnapshot(logger *zap.Logger, vegacapsuleBinary, vegacapsuleHome string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fmt.Println(vegacapsuleBinary)
	nodes, err := vegacapsule.ListNodes(vegacapsuleBinary, vegacapsuleHome)
	if err != nil {
		return fmt.Errorf("failed to list vegacapsule nodes: %w", err)
	}

	var validatorNode *vegacapsule.NodeDetails

	logger.Info("Searching validator node in the vegacapsule network")
	for _, node := range nodes {
		if node.Mode == vegacapsule.VegaModeValidator {
			validatorNode = &node
			break
		}
	}

	if validatorNode == nil {
		return fmt.Errorf(
			"validator node not found in the vegacapsule network in %s",
			vegacapsuleHome,
		)
	}
	logger.Info(fmt.Sprintf("Found validator node %s", validatorNode.Name))

	logger.Info("Moving network forward until new snapshot is produced")
	err = moveNullChainNetworkForward(ctx, logger, *validatorNode)
	if err != nil {
		return fmt.Errorf("failed to move network forward: %w", err)
	}
	logger.Info("New snapshot found")

	return nil
}

func moveNullChainNetworkForward(ctx context.Context, logger *zap.Logger, nodeDetails vegacapsule.NodeDetails) error {
	coreConfigPath := filepath.Join(nodeDetails.Vega.HomeDir, VegaCoreConfigPath)

	nullchainPort, err := tools.ReadStructuredFileValue(
		"toml",
		coreConfigPath,
		"Blockchain.Null.Port",
	)
	if err != nil {
		return fmt.Errorf("failed to get nullchain port: %w", err)
	}

	coreGRPCPort, err := tools.ReadStructuredFileValue("toml", coreConfigPath, "API.Port")
	if err != nil {
		return fmt.Errorf("failed to get core grpc port: %w", err)
	}

	if nullchainPort == "" {
		return fmt.Errorf("nullchain port is invalid: empty port")
	}
	if coreGRPCPort == "" {
		return fmt.Errorf("the core gprc port is invalid: empty port")
	}

	logger.Info("Creating Core GRPC client")
	coreClient := core.NewClient(
		[]string{fmt.Sprintf("localhost:%s", coreGRPCPort)},
		5*time.Second,
		logger,
	)
	dialContext, dialCancel := context.WithTimeout(ctx, 15*time.Second)
	defer dialCancel()
	coreClient.MustDialConnectionIgnoreTime(dialContext)

	logger.Info("Getting snapshot.interval.length from network parameters")
	networkParameters, err := coreClient.CoreNetworkParameters(ctx, "snapshot.interval.length")
	if err != nil {
		return fmt.Errorf("failed to get network parameters from core api: %w", err)
	}
	if len(networkParameters) <= 0 {
		return fmt.Errorf(
			"failed to get network parameters from core api: empty list returned",
		)
	}

	snapshotLength := 0
	for _, parameter := range networkParameters {
		if parameter.Key == "snapshot.interval.length" {
			snapshotLength, err = strconv.Atoi(parameter.Value)
			if err != nil {
				return fmt.Errorf(
					"failed to convert value of snapshot.interval.length to int: %w",
					err,
				)
			}
			break
		}
	}
	logger.Info(fmt.Sprintf("Snapsgot interval is %d", snapshotLength))

	logger.Info("Fetching current network height")
	initialNetworkHeight, err := getNetworkHeight(ctx, coreClient)
	if err != nil {
		return fmt.Errorf("failed to fetch initial network height: %w", err)
	}

	expectedBlock := initialNetworkHeight + snapshotLength*2
	currentnetworkHeight := 0
	// wait about 120 secs
	for i := 0; i < 120; i++ {
		forwardBody := []byte(`{"forward": "30s"}`)
		logger.Info("Moving the chain 30 seconds into the future")
		if _, err := http.Post(fmt.Sprintf("http://localhost:%s/api/v1/forwardtime", nullchainPort), "application/json", bytes.NewBuffer(forwardBody)); err != nil {
			logger.Info(
				"failed to send post request to move the chain 5 sec into the future",
				zap.Error(err),
			)
		}

		currentnetworkHeight, err = getNetworkHeight(ctx, coreClient)
		if err != nil {
			return fmt.Errorf("failed to get network height: %w", err)
		}

		if currentnetworkHeight <= expectedBlock {
			logger.Info(
				"... still waiting",
				zap.Int("current network block", currentnetworkHeight),
				zap.Int("expected network block", expectedBlock),
			)
			time.Sleep(1 * time.Second)
			continue
		}

		logger.Info(
			"Snapshot should be produced",
			zap.Int("current network block", currentnetworkHeight),
			zap.Int("expected network block", expectedBlock),
		)
		break
	}

	if currentnetworkHeight <= expectedBlock {
		return fmt.Errorf("network did not move enough to produce snapshot")
	}

	return nil
}

func getNetworkHeight(ctx context.Context, coreClient vegaapi.VegaCoreClient) (int, error) {
	statistics, err := tools.RetryReturn(3, 3*time.Second, func() (*vegaapipb.StatisticsResponse, error) {
		return coreClient.Statistics(ctx)
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get output from the statistic core api: %w", err)
	}

	return int(statistics.Statistics.BlockHeight), nil
}
