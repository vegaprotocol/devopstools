package vegacapsule

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
	"go.uber.org/zap"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"
	vctools "github.com/vegaprotocol/devopstools/vegacapsule"
)

type StartDataNodeFromNetworkHistoryArgs struct {
	*VegacapsuleArgs

	vegacapsuleBinary string
	networkHomePath   string
	baseOnGroup       string
	waitForReplay     bool
}

var startDataNodeFromNetworkHistoryArgs StartDataNodeFromNetworkHistoryArgs

// traderbotCmd represents the traderbot command
var startDataNodeFromNetworkHistoryCmd = &cobra.Command{
	Use:   "start-datanode-from-network-history",
	Short: "The command starts new data-node in the vegacapsule network based on the existing data-node.",
	Long: `The command :
- obtains information about started vegacapsule network,
- take information about all of the snapshots from all of the data-nodes
- find the latest snapshot height and hash,
- update vega & data-node config for new node
- start new data-node with the latest snapshot`,
	Example: `
	# Return last block for vegacapsule network
	devopstools vegacapsule startDataNodeFromNetworkHistoryArgs --network-home-path /home/daniel/www/networkdata/testnet --no-secrets
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := startDataNodeFromNetworkHistory(
			startDataNodeFromNetworkHistoryArgs.Logger,
			startDataNodeFromNetworkHistoryArgs.vegacapsuleBinary,
			startDataNodeFromNetworkHistoryArgs.baseOnGroup,
			startDataNodeFromNetworkHistoryArgs.networkHomePath)

		if err != nil {
			startDataNodeFromNetworkHistoryArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	startDataNodeFromNetworkHistoryArgs.VegacapsuleArgs = &vegacapsuleArgs

	VegacapsuleCmd.AddCommand(startDataNodeFromNetworkHistoryCmd)

	startDataNodeFromNetworkHistoryCmd.PersistentFlags().StringVar(
		&startDataNodeFromNetworkHistoryArgs.vegacapsuleBinary,
		"vegacapsule-bin",
		"vegacapsule",
		"TBD")

	startDataNodeFromNetworkHistoryCmd.PersistentFlags().StringVar(
		&startDataNodeFromNetworkHistoryArgs.baseOnGroup,
		"base-on-group",
		"",
		"TBD")

	startDataNodeFromNetworkHistoryCmd.PersistentFlags().BoolVar(
		&startDataNodeFromNetworkHistoryArgs.waitForReplay,
		"wait-for-replay",
		false,
		"TBD")

	startDataNodeFromNetworkHistoryCmd.PersistentFlags().StringVar(
		&startDataNodeFromNetworkHistoryArgs.networkHomePath,
		"network-home-path",
		"",
		"TBD")
}

func startDataNodeFromNetworkHistory(logger *zap.Logger, vegacapsuleBinary, baseOnGroup, networkHomePath string) error {
	logger.Info("Listening nodes from vegacapsule")
	vcNodes, err := vctools.ListNodes(vegacapsuleBinary, networkHomePath)
	if err != nil {
		return fmt.Errorf("failed to list vegacapsule nodes: %w", err)
	}
	logger.Info("Filtering data node from the network")

	var dataNode *vctools.NodeDetails
	for _, node := range vcNodes {
		if node.DataNode != nil {
			dataNode = &node
			break
		}
	}

	if dataNode == nil {
		return fmt.Errorf("no data-node exist on network")
	}

	grpcPort, err := tools.ReadStructuredFileValue("toml", dataNode.DataNode.ConfigFilePath, "API.Port")
	if err != nil {
		return fmt.Errorf("failed to read API.Port from the data node config file(%s): %w", dataNode.DataNode.ConfigFilePath, err)
	}

	logger.Info("Creating data-node gRPC connection")
	dataNodeClient := datanode.NewDataNode([]string{fmt.Sprintf("127.0.0.1:%s", grpcPort)}, 3*time.Second, logger)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	dataNodeClient.MustDialConnection(ctx)

	logger.Info("Collecting available snapshots")
	snapshots, err := dataNodeClient.ListCoreSnapshots()
	if err != nil {
		return fmt.Errorf("failed to list snapshots from the data-node: %w", err)
	}

	if len(snapshots) < 1 {
		return fmt.Errorf("no snapshot available for core")
	}

	logger.Info("Adding new data-node to the network")
	newNodeDetails, err := vctools.AddNodes(logger, vegacapsuleBinary, vctools.AddNodesBaseOn{
		Group: baseOnGroup,
	}, false, networkHomePath)

	if err != nil {
		return fmt.Errorf("failed to add new node to the vegacapsule: %w", err)
	}

	if newNodeDetails.DataNode == nil {
		return fmt.Errorf("new node does not include data-node")
	}

	logger.Info("Updating core config", zap.String("config-file", newNodeDetails.Vega.ConfigFilePath))
	if err := updateCoreConfig(logger, newNodeDetails.Vega.ConfigFilePath); err != nil {
		return fmt.Errorf("failed to update core config: %w", err)
	}

	selectedSnapshot := snapshots[0]
	logger.Info("Updating tendermint config", zap.String("config-file", newNodeDetails.Tendermint.ConfigFilePath))
	if err := updateTendermintConfig(logger, newNodeDetails.Tendermint.ConfigFilePath, selectedSnapshot.BlockHash, selectedSnapshot.BlockHeight); err != nil {
		return fmt.Errorf("failed to update tendermint config: %w", err)
	}

	logger.Info("Updating data-node config", zap.String("config-file", newNodeDetails.DataNode.ConfigFilePath))
	if err := updateDataNodeConfig(logger, newNodeDetails.DataNode.ConfigFilePath); err != nil {
		return fmt.Errorf("failed to update data-node config: %w", err)
	}

	if err := vctools.StartNode(logger, newNodeDetails.Name, vegacapsuleBinary, networkHomePath); err != nil {
		return fmt.Errorf("failed to start the %s node: %w", newNodeDetails.Name, err)
	}

	return nil
}

func updateCoreConfig(logger *zap.Logger, configPath string) error {
	coreConfigNode, err := dasel.NewFromFile(configPath, "toml")
	if err != nil {
		return fmt.Errorf("failed to read core config: %w", err)
	}
	logger.Info("Setting Snapsgot.StartHeight to -1", zap.String("config-file", configPath))
	if err := coreConfigNode.Put("Snapshot.StartHeight", -1); err != nil {
		return fmt.Errorf("failed to set Snapshot.StartHeight in the vega node config: %w", err)
	}
	if err := coreConfigNode.WriteToFile(configPath, "toml", []storage.ReadWriteOption{}); err != nil {
		return fmt.Errorf("failed to write the %s file: %w", configPath, err)
	}

	return nil
}

func updateTendermintConfig(logger *zap.Logger, configPath, snapshotHash string, snapshotHeight uint64) error {
	tmConfigNode, err := dasel.NewFromFile(configPath, "toml")
	if err != nil {
		return fmt.Errorf("failed to read core config: %w", err)
	}

	logger.Info("Setting statesync.enable to true", zap.String("config-file", configPath))
	if err := tmConfigNode.Put("statesync.enable", true); err != nil {
		return fmt.Errorf("failed to set statesync.enable field in the thendermint config: %w", err)
	}

	logger.Info(fmt.Sprintf("Setting statesync.trust_hash to %s", snapshotHash), zap.String("config-file", configPath))
	if err := tmConfigNode.Put("statesync.trust_hash", snapshotHash); err != nil {
		return fmt.Errorf("failed to set statesync.trust_hash field in the thendermint config: %w", err)
	}

	logger.Info(fmt.Sprintf("Setting statesync.trust_height to %s", snapshotHeight), zap.String("config-file", configPath))
	if err := tmConfigNode.Put("statesync.trust_height", snapshotHeight); err != nil {
		return fmt.Errorf("failed to set statesync.trust_height field in the thendermint config: %w", err)
	}

	if err := tmConfigNode.WriteToFile(configPath, "toml", []storage.ReadWriteOption{}); err != nil {
		return fmt.Errorf("failed to write the %s file: %w", configPath, err)
	}

	return nil
}

func updateDataNodeConfig(logger *zap.Logger, configPath string) error {
	dataNodeConfigNode, err := dasel.NewFromFile(configPath, "toml")
	if err != nil {
		return fmt.Errorf("failed to read core config: %w", err)
	}

	logger.Info("Setting AutoInitialiseFromNetworkHistory to true", zap.String("config-file", configPath))
	if err := dataNodeConfigNode.Put("AutoInitialiseFromNetworkHistory", true); err != nil {
		return fmt.Errorf("failed to set AutoInitialiseFromNetworkHistory in the data-node config: %w", err)
	}
	if err := dataNodeConfigNode.WriteToFile(configPath, "toml", []storage.ReadWriteOption{}); err != nil {
		return fmt.Errorf("failed to write the %s file: %w", configPath, err)
	}

	return nil
}
