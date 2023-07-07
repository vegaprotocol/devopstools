package snapshotcompatibility

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegacapsule"
	"github.com/vegaprotocol/devopstools/vegacmd"
)

type ProductNewSnapshotArgs struct {
	*SnapshotCompatibilityArgs

	VegacapsuleHome   string
	VegacapsuleBinary string
	VegaBinary        string
}

var produceNewSnapshotArgs ProductNewSnapshotArgs

var produceNewSnapshotCmd = &cobra.Command{
	Use:   "produce-new-snapshot",
	Short: "Produces new snapshot for the snapshot-compatibility network, convert it to JSON and save in given file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runProduceNewSnapshot(produceNewSnapshotArgs.Logger,
			produceNewSnapshotArgs.VegaBinary,
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
		StringVar(&produceNewSnapshotArgs.VegacapsuleBinary, "vegacapsule-binry", "vegacapsule", "The vegacapsule binary path")
	produceNewSnapshotCmd.PersistentFlags().
		StringVar(&produceNewSnapshotArgs.VegaBinary, "vega-binry", "vega", "The vega binary path")
}

func runProduceNewSnapshot(
	logger *zap.Logger,
	vegaBinary, vegacapsuleBinary, vegacapsuleHome string,
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
	producedSnapshot, err := moveNullChainNetworkForward(logger, vegaBinary, *validatorNode)
	if err != nil {
		return fmt.Errorf("failed to move network forward: %w", err)
	}
	logger.Info("New snapshot found")

	_ = producedSnapshot
	return nil
}

func moveNullChainNetworkForward(
	logger *zap.Logger,
	vegaBinary string,
	nodeDetails vegacapsule.NodeDetails,
) (*vegacmd.CoreToolsSnapshot, error) {
	stopChannel := make(<-chan struct{})

	nullchainPort, err := tools.ReadStructuredFileValue(
		"toml",
		filepath.Join(nodeDetails.Vega.HomeDir, VegaCoreConfigPath),
		"Blockchain.Null.Port",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get nullchain port: %w", err)
	}

	if nullchainPort == "" {
		return nil, fmt.Errorf("nullchain port is invalid: empty port")
	}

	initialLatestSnapshot, err := vegacmd.LatestCoreSnapshot(
		vegaBinary,
		vegacmd.CoreSnashotInput{VegaHome: nodeDetails.Vega.HomeDir},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial latest snapshot: %w", err)
	}

	go func(logger *zap.Logger, stopChannel <-chan struct{}, port string) {
		ticker := time.NewTicker(500 * time.Millisecond)
		forwardBody := []byte(`{"forward": "5s"}`)
		for {
			select {
			case <-stopChannel:
				return
			case <-ticker.C:
				logger.Info("Moving the chain 5 seconds into the future")
				if _, err := http.Post(fmt.Sprintf("http://localhost:%s", port), "application/json", bytes.NewBuffer(forwardBody)); err != nil {
					logger.Info(
						"failed to send post request to move the chain 5 sec into the future",
						zap.Error(err),
					)
				}
			}
		}
	}(logger, stopChannel, nullchainPort)

	var latestSnapshot *vegacmd.CoreToolsSnapshot
	for i := 0; i < 120; i++ {
		latestSnapshot, err = vegacmd.LatestCoreSnapshot(
			vegaBinary,
			vegacmd.CoreSnashotInput{VegaHome: nodeDetails.Vega.HomeDir},
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get latest snapshot when waiting for a new one: %w",
				err,
			)
		}

		if initialLatestSnapshot.Height < latestSnapshot.Height {
			logger.Info(fmt.Sprintf("Found new snapshot at block %d", latestSnapshot.Height))
			break
		}

		logger.Info("... still waiting for a new snapshot")
		time.Sleep(1 * time.Second)
	}

	return latestSnapshot, nil
}
