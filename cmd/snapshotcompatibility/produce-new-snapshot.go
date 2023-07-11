package snapshotcompatibility

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/vegaapi/core"
	"github.com/vegaprotocol/devopstools/vegacapsule"
)

type ProduceNewSnapshotArgs struct {
	*SnapshotCompatibilityArgs

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
	produceNewSnapshotArgs.SnapshotCompatibilityArgs = &snapshotCompatibilityArgs
	produceNewSnapshotCmd.PersistentFlags().
		StringVar(&produceNewSnapshotArgs.VegacapsuleHome, "vegacapsule-home", "", "The custom vegacapsule home")
	produceNewSnapshotCmd.PersistentFlags().
		StringVar(&produceNewSnapshotArgs.VegacapsuleBinary, "vegacapsule-binary", "vegacapsule", "The vegacapsule binary path")
}

func runProduceNewSnapshot(
	logger *zap.Logger,
	vegacapsuleBinary, vegacapsuleHome string,
) error {
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
	err = moveNullChainNetworkForward(logger, *validatorNode)
	if err != nil {
		return fmt.Errorf("failed to move network forward: %w", err)
	}
	logger.Info("New snapshot found")

	return nil
}

func moveNullChainNetworkForward(
	logger *zap.Logger,

	nodeDetails vegacapsule.NodeDetails,
) error {
	stopChannel := make(chan struct{})

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

	logger.Info("Crateing Core GRPC client")
	coreClient := core.NewCoreClient(
		[]string{fmt.Sprintf("localhost:%s", coreGRPCPort)},
		5*time.Second,
		logger,
	)
	dialContext, dialCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancel()
	coreClient.MustDialConnectionIgnoreTime(dialContext)

	logger.Info("Getting snapshot.interval.length from network parameters")
	networkParameters, err := coreClient.CoreNetworkParameters("snapshot.interval.length")
	if err != nil {
		return fmt.Errorf("failed to get network parameters from core api: %w", err)
	}
	if len(networkParameters) < 0 {
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
	initialNetworkHeight, err := getNetworkHeight(coreClient)
	if err != nil {
		return fmt.Errorf("failed to fetch initial network height: %w", err)
	}

	go func(logger *zap.Logger, stopChannel <-chan struct{}, port string) {
		ticker := time.NewTicker(500 * time.Millisecond)
		forwardBody := []byte(`{"forward": "30s"}`)
		for {
			select {
			case <-stopChannel:
				return
			case <-ticker.C:
				logger.Info("Moving the chain 30 seconds into the future")
				if _, err := http.Post(fmt.Sprintf("http://localhost:%s/api/v1/forwardtime", port), "application/json", bytes.NewBuffer(forwardBody)); err != nil {
					logger.Info(
						"failed to send post request to move the chain 5 sec into the future",
						zap.Error(err),
					)
				}
			}
		}
	}(logger, stopChannel, nullchainPort)

	expectedBlock := initialNetworkHeight + snapshotLength*2
	currentnetworkHeight := 0
	// wait about 120 secs
	for i := 0; i < 120; i++ {
		currentnetworkHeight, err = getNetworkHeight(coreClient)
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

	stopChannel <- struct{}{}
	if currentnetworkHeight <= expectedBlock {
		return fmt.Errorf("network did not move enough to produce snapshot")
	}

	return nil
}

func getNetworkHeight(coreClient vegaapi.VegaCoreClient) (int, error) {
	statistics, err := coreClient.Statistics()
	if err != nil {
		return 0, fmt.Errorf("failed to get output from the statistic core api: %w", err)
	}

	return int(statistics.Statistics.BlockHeight), nil
}
