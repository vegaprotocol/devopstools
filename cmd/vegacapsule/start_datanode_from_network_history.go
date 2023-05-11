package vegacapsule

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	v1 "code.vegaprotocol.io/vega/protos/vega/events/v1"
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
			startDataNodeFromNetworkHistoryArgs.networkHomePath,
			startDataNodeFromNetworkHistoryArgs.waitForReplay)

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
		"Path to the vegacapsule binary")

	startDataNodeFromNetworkHistoryCmd.PersistentFlags().StringVar(
		&startDataNodeFromNetworkHistoryArgs.baseOnGroup,
		"base-on-group",
		"",
		"Node set name to create")

	startDataNodeFromNetworkHistoryCmd.PersistentFlags().BoolVar(
		&startDataNodeFromNetworkHistoryArgs.waitForReplay,
		"wait-for-replay",
		false,
		"Determine if we should wait for the node after it has been started")

	startDataNodeFromNetworkHistoryCmd.PersistentFlags().StringVar(
		&startDataNodeFromNetworkHistoryArgs.networkHomePath,
		"network-home-path",
		"",
		"Custom path for the network")
}

func startDataNodeFromNetworkHistory(logger *zap.Logger, vegacapsuleBinary, baseOnGroup, networkHomePath string, waitForReplay bool) error {
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

	logger.Info("Getting most recent network history segment")
	segment, err := dataNodeClient.LastNetworkHistorySegment()
	if err != nil {
		return fmt.Errorf("failed to get most recent network history segment: %w", err)
	}
	logger.Info("Selected most recent network history segment", zap.Int64("from", segment.FromHeight), zap.Int64("to", segment.ToHeight))

	var selectedSnapshot *v1.CoreSnapshotData
	for idx, snapshot := range snapshots {
		if segment.ToHeight < int64(snapshot.BlockHeight) {
			continue
		}

		if selectedSnapshot != nil && selectedSnapshot.BlockHeight > snapshot.BlockHeight {
			continue
		}

		fmt.Printf("using snapshot for %d\n", snapshot.BlockHeight)

		selectedSnapshot = &(snapshots[idx])
	}

	if selectedSnapshot == nil {
		return fmt.Errorf("failed to select core snapshot")
	}

	logger.Info("Selected snapshot for restart", zap.Uint64("height", selectedSnapshot.BlockHeight), zap.String("hash", selectedSnapshot.BlockHash))

	logger.Info("Updating core config", zap.String("config-file", newNodeDetails.Vega.ConfigFilePath))
	if err := updateCoreConfig(logger, newNodeDetails.Vega.ConfigFilePath, selectedSnapshot.BlockHeight); err != nil {
		return fmt.Errorf("failed to update core config: %w", err)
	}

	// selectedSnapshot := snapshots[0]
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

	if !waitForReplay {
		return nil
	}

	logger.Info("Getting Gateway.Port value from the data node for old node", zap.String("config-file", dataNode.DataNode.ConfigFilePath))
	oldNodeRESTPort, err := tools.ReadStructuredFileValue("toml", dataNode.DataNode.ConfigFilePath, "Gateway.Port")
	if err != nil {
		return fmt.Errorf("failed to read Gateway.Port from the old node config file(%s): %w", dataNode.DataNode.ConfigFilePath, err)
	}

	logger.Info("Getting Gateway.Port value from the data node for new node", zap.String("config-file", newNodeDetails.DataNode.ConfigFilePath))
	newNodeRESTPort, err := tools.ReadStructuredFileValue("toml", newNodeDetails.DataNode.ConfigFilePath, "Gateway.Port")
	if err != nil {
		return fmt.Errorf("failed to read Gateway.Port from the new node config file(%s): %w", dataNode.DataNode.ConfigFilePath, err)
	}

	logger.Info("Waiting for node to replay")
	waitCtx, waitCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer waitCancel()
	if err := waitForNode(logger, waitCtx, oldNodeRESTPort, newNodeRESTPort); err != nil {
		return fmt.Errorf("failed to wait until new data node replied: %w", err)
	}

	return nil
}

func updateCoreConfig(logger *zap.Logger, configPath string, height uint64) error {
	coreConfigNode, err := dasel.NewFromFile(configPath, "toml")
	if err != nil {
		return fmt.Errorf("failed to read core config: %w", err)
	}
	logger.Info(fmt.Sprintf("Setting Snapsgot.StartHeight to %d", height), zap.String("config-file", configPath))
	if err := coreConfigNode.Put("Snapshot.StartHeight", height); err != nil {
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

func waitForNode(logger *zap.Logger, ctx context.Context, oldNodeRESTPort, newNodeRESTPort string) error {
	const threshold = 10

	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ticker.C:
			oldNodeHeight := getBlockHeight(logger, oldNodeRESTPort)
			newNodeHeight := getBlockHeight(logger, newNodeRESTPort)

			if oldNodeHeight != 0 && newNodeHeight != 0 && newNodeHeight > oldNodeHeight-threshold {
				return nil
			}

			logger.Info("Still waiting", zap.Int("expected-block", oldNodeHeight-threshold), zap.Int("current-block", newNodeHeight))
		case <-ctx.Done():
			return fmt.Errorf("context cancaled")
		}
	}
}

type statistics struct {
	Statistics struct {
		BlockHeight string `json:"blockHeight"`
	} `json:"statistics"`
}

func getBlockHeight(logger *zap.Logger, port string) int {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/statistics", port))
	if err != nil {
		logger.Debug(fmt.Sprintf("failed to call http://127.0.0.1:%s/statistics", port), zap.Error(err))
		return 0
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read response body buffer", zap.Error(err))
		return 0
	}

	stats := statistics{}
	if err := json.Unmarshal(respBytes, &stats); err != nil {
		logger.Error("failed to unmarshall data", zap.Error(err))
		return 0
	}

	result, err := strconv.Atoi(stats.Statistics.BlockHeight)
	if err != nil {
		logger.Error("failed to convert block height into int", zap.Error(err))
		return 0
	}

	return result
}
