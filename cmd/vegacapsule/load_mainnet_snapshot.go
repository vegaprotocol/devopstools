package vegacapsule

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/tools"
	vctools "github.com/vegaprotocol/devopstools/vegacapsule"
	"github.com/vegaprotocol/devopstools/vegasnapshot"
	"go.uber.org/zap"
)

type LoadMainnetSnapshotArgs struct {
	*VegacapsuleArgs

	snapshotSourcePath string
	workDirPath        string
}

var loadMainnetSnapshotArgs LoadMainnetSnapshotArgs

// traderbotCmd represents the traderbot command
var loadMainnetSnapshotCmd = &cobra.Command{
	Use:   "load-mainnet-snapshot",
	Short: "Load mainnet snapshot into the generated vegacapsule-network",
	Long:  `Snapshot must be downloaded to local file system. To download snapshot see the 'devopstools remote download-snapshot' command.`,

	Run: func(cmd *cobra.Command, args []string) {
		if err := loadSnapshot(
			loadMainnetSnapshotArgs.Logger,
			loadMainnetSnapshotArgs.vegacapsuleBinary,
			loadMainnetSnapshotArgs.networkHomePath,
			loadMainnetSnapshotArgs.snapshotSourcePath,
			loadMainnetSnapshotArgs.workDirPath,
		); err != nil {
			loadMainnetSnapshotArgs.Logger.Fatal("failed to load snapshot", zap.Error(err))
		}
	},
}

func init() {
	loadMainnetSnapshotArgs.VegacapsuleArgs = &vegacapsuleArgs

	loadMainnetSnapshotCmd.PersistentFlags().StringVar(
		&loadMainnetSnapshotArgs.snapshotSourcePath,
		"snapshot-source-path",
		"",
		"Path to the snapshot source downloaded from the mainnet")

	loadMainnetSnapshotCmd.PersistentFlags().StringVar(
		&loadMainnetSnapshotArgs.workDirPath,
		"work-dir-path",
		"./",
		"Path to the work dir")

	if err := loadMainnetSnapshotCmd.MarkPersistentFlagRequired("snapshot-source-path"); err != nil {
		panic(err)
	}

	VegacapsuleCmd.AddCommand(loadMainnetSnapshotCmd)
}

func loadSnapshot(logger *zap.Logger, vegacapsuleBinary, vegacapsuleHomePath, snapshotSourcePath, workDirPath string) error {
	snapshot, err := vegasnapshot.OpenSnapshotDB(snapshotSourcePath)
	if err != nil {
		return fmt.Errorf("failed to open db snapshot: %w", err)
	}

	snapshotJSONOutput := filepath.Join(workDirPath, "./snapshot.json")

	if err := snapshot.WriteSnapshotAsJSON(snapshotJSONOutput); err != nil {
		return fmt.Errorf("failed to write snapshot to the JSON file: %w", err)
	}

	nodesDetails, err := vctools.ListNodes(vegacapsuleBinary, vegacapsuleHomePath)
	if err != nil {
		return fmt.Errorf("failed to list nodes for vegacapsule network: %w", err)
	}

	for _, node := range nodesDetails {
		vegaHomePath := node.Vega.HomeDir
		snapshotDirPath := filepath.Join(vegaHomePath, "/state/node/snapshots")

		logger.Info(fmt.Sprintf("Ensuring snapshot directory exists for node %s", node.Name), zap.String("path", snapshotDirPath))
		if err := os.MkdirAll(snapshotDirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to ensure snapshot dir exists for node %s: %w", node.Name, err)
		}

		if err := tools.CopyDirectory(snapshotSourcePath, snapshotDirPath); err != nil {
			return fmt.Errorf("failed to copy mainnet snapshot for node %s: %w", node.Name, err)
		}
	}

	return nil
}
