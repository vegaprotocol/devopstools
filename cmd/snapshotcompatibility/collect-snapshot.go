package snapshotcompatibility

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegacapsule"
	"github.com/vegaprotocol/devopstools/vegacmd"
)

type CollectSnapshotArgs struct {
	*SnapshotCompatibilityArgs

	VegacapsuleHome    string
	VegacapsuleBinary  string
	VegaBinary         string
	SnapshotJSONOutput string
}

var collectSnapshotArgs CollectSnapshotArgs

var collectSnapshotCmd = &cobra.Command{
	Use:   "collect-snapshot",
	Short: "Collects latest snashot in the JSON format from the stopped network",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCollectSnapshot(collectSnapshotArgs.Logger,
			collectSnapshotArgs.VegaBinary,
			collectSnapshotArgs.VegacapsuleBinary,
			collectSnapshotArgs.VegacapsuleHome,
			collectSnapshotArgs.SnapshotJSONOutput,
		); err != nil {
			collectSnapshotArgs.Logger.Fatal("failed to produce new snapshot", zap.Error(err))
		}
	},
}

func init() {
	collectSnapshotArgs.SnapshotCompatibilityArgs = &snapshotCompatibilityArgs
	collectSnapshotCmd.PersistentFlags().
		StringVar(&collectSnapshotArgs.VegacapsuleHome, "vegacapsule-home", "", "The custom vegacapsule home")
	collectSnapshotCmd.PersistentFlags().
		StringVar(&collectSnapshotArgs.VegacapsuleBinary, "vegacapsule-binary", "vegacapsule", "The vegacapsule binary path")
	collectSnapshotCmd.PersistentFlags().
		StringVar(&collectSnapshotArgs.VegaBinary, "vega-binary", "vega", "The vega binary path")
	collectSnapshotCmd.PersistentFlags().
		StringVar(&collectSnapshotArgs.SnapshotJSONOutput, "snapshot-json-output", "./snapshot.json", "The file where JSON version of snapshot will be saved to")
}

func runCollectSnapshot(
	logger *zap.Logger,
	vegaBinary, vegacapsuleBinary, vegacapsuleHome string,
	jsonOutputFile string,
) error {
	var validatorNode *vegacapsule.NodeDetails

	logger.Info("Searching for validator node in the vegacapsule network")
	nodes, err := vegacapsule.ListNodes(vegacapsuleBinary, vegacapsuleHome)
	if err != nil {
		return fmt.Errorf("failed to list vegacapsule nodes: %w", err)
	}

	for _, node := range nodes {
		if node.Mode == vegacapsule.VegaModeValidator {
			validatorNode = &node
			break
		}
	}

	if validatorNode == nil {
		return fmt.Errorf("failed to find validator node in the vegacapsule network")
	}
	logger.Info("Found validator node", zap.String("name", validatorNode.Name))

	logger.Info("Looking for latest core snapshot")
	latestSnapshot, err := vegacmd.LatestCoreSnapshot(vegaBinary, vegacmd.CoreSnashotInput{
		VegaHome: validatorNode.Vega.HomeDir,
	})
	if err != nil {
		return fmt.Errorf(
			"failed to find latest snapshot for vega node %s: %w",
			validatorNode.Name,
			err,
		)
	}
	if latestSnapshot == nil {
		return fmt.Errorf("snapshot not found")
	}
	logger.Info("Snapshot selected", zap.Int("height", latestSnapshot.Height))

	logger.Info(fmt.Sprintf("Saving snapshot in JSON to %s", jsonOutputFile))
	snapshotToJSONArgs := []string{
		"tools",
		"snapshot",
		"--output", "json",
		"--block-height", fmt.Sprint(latestSnapshot.Height),
		"--snapshot-contents", jsonOutputFile,
		"--home", validatorNode.Vega.HomeDir,
	}
	if _, err := tools.ExecuteBinary(vegaBinary, snapshotToJSONArgs, nil); err != nil {
		return fmt.Errorf("failed to convert snapshot to JSON: %w", err)
	}
	logger.Info(fmt.Sprintf("JSON snapshot saved in %s", jsonOutputFile))

	return nil
}
