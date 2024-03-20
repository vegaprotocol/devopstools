package snapshotcompatibility

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegacapsule"
	"github.com/vegaprotocol/devopstools/vegacmd"

	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
	"go.uber.org/zap"
)

type LoadSnapshotArgs struct {
	*SnapshotCompatibilityArgs

	VegacapsuleHome   string
	VegacapsuleBinary string
	VegaBinary        string
	SnapshotLocation  string
}

func (args LoadSnapshotArgs) Check() error {
	if _, err := vegacapsule.ListNodes(args.VegacapsuleBinary, args.VegacapsuleHome); err != nil {
		return fmt.Errorf("no vegacapsule network generated: %w", err)
	}

	version, err := tools.ExecuteBinary(args.VegaBinary, []string{"version"}, nil)
	if err != nil {
		return fmt.Errorf(
			"failed to execute vega binary, try to provide a different binry with the --vega-binary flag: %s, %w",
			version,
			err,
		)
	}

	if !tools.FileExists(args.SnapshotLocation) {
		return fmt.Errorf("cannot find downloaded mainnet snapshot in %s", args.SnapshotLocation)
	}

	return nil
}

var loadSnapshotArgs LoadSnapshotArgs

var loadSnapshotCmd = &cobra.Command{
	Use:   "load-snapshot",
	Short: "Download the snapshot from the mainnet and load it to local network",
	Run: func(cmd *cobra.Command, args []string) {
		if err := loadSnapshotArgs.Check(); err != nil {
			loadSnapshotArgs.Logger.Fatal("invalid input", zap.Error(err))
			return
		}

		if err := runLoadSnapshot(loadSnapshotArgs.Logger,
			loadSnapshotArgs.SnapshotLocation,
			loadSnapshotArgs.VegaBinary,
			loadSnapshotArgs.VegacapsuleBinary,
			loadSnapshotArgs.VegacapsuleHome,
		); err != nil {
			loadSnapshotArgs.Logger.Fatal(
				"failed to prepare for snapshot compatibility pipeline",
				zap.Error(err),
			)
			return
		}
	},
}

func init() {
	loadSnapshotArgs.SnapshotCompatibilityArgs = &snapshotCompatibilityArgs

	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.VegacapsuleHome, "vegacapsule-home", "", "The custom vegacapsule home")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.VegaBinary, "vega-binary", "vega", "Path to the vega executable on local machine")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.SnapshotLocation, "snapshot-location", "/tmp/mainnet-snapshot", "Path to the downloaded snapshot from mainnet server")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.VegacapsuleBinary, "vegacapsule-binary", "vegacapsule", "Path to the vegacapsule executable")
}

// TODO: Check vegacapule nodes and return list of all of the validators
func vegacapsuleValidatorsCoreHomePaths(
	vegacapsuleBinary, vegacapsuleHome string,
) ([]string, error) {
	result := []string{}

	nodes, err := vegacapsule.ListNodes(vegacapsuleBinary, vegacapsuleHome)
	if err != nil {
		return nil, fmt.Errorf("failed to list vegacapsule nodes: %w", err)
	}

	for _, node := range nodes {
		if node.Mode != vegacapsule.VegaModeValidator {
			continue
		}
		result = append(result, node.Vega.HomeDir)
	}

	return result, nil
}

func runLoadSnapshot(
	logger *zap.Logger,
	snapshotLocation string,
	vegaBinary, vegacapsuleBinary, vegacapsuleHome string,
) error {
	validatorHomePaths, err := vegacapsuleValidatorsCoreHomePaths(
		vegacapsuleBinary,
		vegacapsuleHome,
	)
	if err != nil {
		return fmt.Errorf("failed to get validator home paths: %w", err)
	}

	for _, validatorHomePath := range validatorHomePaths {
		snapshotPath := filepath.Join(validatorHomePath, VegaSnapshotPath)
		logger.Info(fmt.Sprintf("Creating folder for snapshot DB locally: %s", snapshotPath))
		if err := os.MkdirAll(snapshotPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create snapshot db(%s): %w", snapshotPath, err)
		}
	}

	logger.Info("Selecting block height for the restart")
	restartHeight, err := tools.RetryReturn(3, 5*time.Second, func() (int, error) {
		return selectSnapshotHeight(vegaBinary, filepath.Join(snapshotLocation, "snapshots"))
	})
	if err != nil {
		return fmt.Errorf("failed to select height for network start: %w", err)
	}
	logger.Info("Selected height for restart", zap.Int("height", restartHeight))

	for _, validatorHomePath := range validatorHomePaths {
		logger.Info(
			"Updating core config",
			zap.String("node_home", validatorHomePath),
			zap.Int("height", restartHeight),
		)
		if err := updateCoreConfig(validatorHomePath, restartHeight); err != nil {
			return fmt.Errorf(
				"failed to update core config for node(%s): %w",
				validatorHomePath,
				err,
			)
		}
		logger.Info("Config updated")

		snapshotDestination := filepath.Join(validatorHomePath, VegaSnapshotPath)
		snapshotSource := filepath.Join(snapshotLocation, "snapshots")

		logger.Info("Calling unsafe reset all for core", zap.String("node", validatorHomePath))
		if _, err := tools.ExecuteBinary(vegaBinary, []string{"unsafe_reset_all", "--home", validatorHomePath}, nil); err != nil {
			return fmt.Errorf("failed to unsafe reset all for home %s: %w", validatorHomePath, err)
		}
		logger.Info("Unsafe reset all successful")

		logger.Info(
			"Loading snapshot database into the core node",
			zap.String("source", snapshotSource),
			zap.String("destination", snapshotDestination),
		)

		if err := tools.CopyDir(snapshotSource, snapshotDestination); err != nil {
			return fmt.Errorf(
				"failed to copy snapshot from temporary location to node %s: %w",
				snapshotDestination,
				err,
			)
		}

		logger.Info("Snapshot database loaded")

		// force protocol-upgrade flag to be true in the snapshots
		logger.Info("Forcing protocol-upgrade flag to true in snapshot")
		snapshotToJSONArgs := []string{
			"tools",
			"snapshot",
			"--set-pup",
			"--home", validatorHomePath,
		}
		if _, err := tools.ExecuteBinary(vegaBinary, snapshotToJSONArgs, nil); err != nil {
			return fmt.Errorf("failed to set protocol-upgrade flag in latest snapshot: %w", err)
		}
	}
	return nil
}

func selectSnapshotHeight(vegaBinary, snapshotDbLocation string) (int, error) {
	snapshotsResponse, err := vegacmd.ListCoreSnapshots(
		vegaBinary,
		vegacmd.CoreSnashotInput{SnapshotDbPath: snapshotDbLocation},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to list core snapshots: %w", err)
	}

	snapshotList := snapshotsResponse.Snapshots
	sort.Slice(snapshotList, func(i, j int) bool {
		return snapshotList[i].Height > snapshotList[j].Height
	})
	if len(snapshotList) < 1 {
		return 0, fmt.Errorf(
			"not enough snapshots: expected at least 1 snapshots, %d got",
			len(snapshotList),
		)
	}
	return snapshotList[0].Height, nil
}

func updateCoreConfig(coreHome string, startHeight int) error {
	configPath := filepath.Join(coreHome, VegaCoreConfigPath)
	coreConfigNode, err := dasel.NewFromFile(configPath, "toml")
	if err != nil {
		return fmt.Errorf("failed to read core config: %w", err)
	}
	if err := coreConfigNode.Put("Snapshot.StartHeight", startHeight); err != nil {
		return fmt.Errorf("failed to set Snapshot.StartHeight in the vega node config: %w", err)
	}
	if err := coreConfigNode.WriteToFile(configPath, "toml", []storage.ReadWriteOption{}); err != nil {
		return fmt.Errorf("failed to write the %s file: %w", configPath, err)
	}

	return nil
}
