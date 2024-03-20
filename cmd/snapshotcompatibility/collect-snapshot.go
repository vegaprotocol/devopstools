package snapshotcompatibility

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegacapsule"
	"github.com/vegaprotocol/devopstools/vegacmd"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
			collectSnapshotArgs.Logger.Fatal("failed to collect snapshot", zap.Error(err))
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

	logger.Info("Creating temporary directory to access the snapshot db")
	tempDir, err := os.MkdirTemp("", "snapshot-compatibility-collect-snapshot")
	if err != nil {
		return fmt.Errorf("failed to create temporary dir: %w", err)
	}
	logger.Info("Temporary dir created", zap.String("path", tempDir))
	defer os.RemoveAll(tempDir)

	snapshotDbSource := filepath.Join(validatorNode.Vega.HomeDir, VegaSnapshotPath)
	logger.Info(
		"Copying the snapshot db",
		zap.String("source", snapshotDbSource),
		zap.String("destiantion", tempDir),
	)
	if err := tools.CopyDir(snapshotDbSource, tempDir); err != nil {
		return fmt.Errorf("failed to copy db snapshot path: %w", err)
	}
	logger.Info("The snapshot db copied")

	newSnapshotDbPath := tempDir
	logger.Info("Looking for latest core snapshot")
	latestSnapshot, err := vegacmd.LatestCoreSnapshot(vegaBinary, vegacmd.CoreSnashotInput{
		SnapshotDbPath: newSnapshotDbPath,
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
		"--db-path", tempDir,
	}
	if _, err := tools.ExecuteBinary(vegaBinary, snapshotToJSONArgs, nil); err != nil {
		return fmt.Errorf("failed to convert snapshot to JSON: %w", err)
	}
	logger.Info(fmt.Sprintf("JSON snapshot saved in %s", jsonOutputFile))

	return nil
}
