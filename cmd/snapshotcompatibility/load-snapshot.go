package snapshotcompatibility

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/vegaprotocol/devopstools/tools"
)

type LoadSnapshotArgs struct {
	*SnapshotCompatibilityArgs

	VegacapsuleHome string
	VegaBinary      string

	SnapshotServerHost     string
	SnapshotRemoteLocation string
	SnapshotServerUser     string
	SnapshotServerKeyFile  string
}

func (args LoadSnapshotArgs) Check() error {
	if !tools.FileExists(args.VegacapsuleHome) {
		return fmt.Errorf("vega home (%s) does not exists", args.VegacapsuleHome)
	}

	version, err := tools.ExecuteBinary(args.VegaBinary, []string{"version"}, nil)
	if err != nil {
		return fmt.Errorf("failed to execute vega binary: %s, %w", version, err)
	}

	if args.SnapshotServerHost == "" {
		return fmt.Errorf("no snapshot server provided")
	}

	if args.SnapshotRemoteLocation == "" {
		return fmt.Errorf("snapshot remote locatino not provided")
	}

	if args.SnapshotServerUser == "" {
		return fmt.Errorf("snapshot server user not provided")
	}

	if !tools.FileExists(args.SnapshotServerKeyFile) {
		return fmt.Errorf("snapshot server key does not exist")
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
			loadSnapshotArgs.SnapshotServerHost,
			loadSnapshotArgs.SnapshotServerUser,
			loadSnapshotArgs.SnapshotServerKeyFile,
			loadSnapshotArgs.SnapshotRemoteLocation,
			loadSnapshotArgs.VegaBinary,
			loadSnapshotArgs.VegacapsuleHome); err != nil {
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

	currentUser, _ := tools.WhoAmI()
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.VegacapsuleHome, "vegacapsule-home", "", "The custom vegacapsule home")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.VegaBinary, "vega-binary", "vega", "Path to the vega executable")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.SnapshotServerHost, "snapshot-server", "", "The source server for the snapshot")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.SnapshotServerUser, "snapshot-server-user", currentUser, "The SSH user for the snapshot server")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.SnapshotRemoteLocation, "snapshot-remote-location", "/home/vega/vega_home/state/node/snapshots", "The location where the snapshot is on the remote server")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.SnapshotServerKeyFile, "snapshot-server-key-file", filepath.Join(tools.CurrentUserHomePath(), ".ssh", "id_rsa"), "The SSH private key used to authenticate user")
}

// TODO: Check vegacapule nodes and return list of all of the validators
func localSnapshotPaths(vegacapsuleHome string) []string {
	return []string{
		filepath.Join(vegacapsuleHome, "vega", "node0", "state", "node", "snapshots"),
	}
}

func runLoadSnapshot(
	logger *zap.Logger,
	snapshotServerHost, snapshotServerUser, snapshotServrKeyFile, snapshotRemoteLocation string,
	vegaBinary, vegacapsuleHome string,
) error {
	localSnapshotsDbPaths := localSnapshotPaths(vegacapsuleHome)
	for _, snapshotPath := range localSnapshotsDbPaths {
		logger.Info(fmt.Sprintf("Creating folder for snapshot DB locally: %s", snapshotPath))
		if err := os.MkdirAll(snapshotPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create snapshot db(%s): %w", snapshotPath, err)
		}
	}

	// Create TMP dir
	// Download the snapshot from the snapshotServerHost
	// run vega tools snapshot --output json ...
	// select 2-nd or 3-rd newest snapshot to load
	// update the config or vega core
	// convert selected snapshot to json
	// move the converted snapshot to new location(flag??)
	return nil
}
