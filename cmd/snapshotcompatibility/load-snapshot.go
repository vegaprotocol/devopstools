package snapshotcompatibility

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/vegaprotocol/devopstools/ssh"
	"github.com/vegaprotocol/devopstools/tools"
)

var (
	SnapshotDBSubPath      = filepath.Join("state", "node", "snapshots")
	DefaultVegacapsuleHome = filepath.Join(tools.CurrentUserHomePath(), ".vegacapsule")
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
		return fmt.Errorf(
			"failed to execute vega binary, try to provide a different binry with the --vega-binary flag: %s, %w",
			version,
			err,
		)
	}

	if args.SnapshotServerHost == "" {
		return fmt.Errorf(
			"no snapshot server provided: provide value with the --snapshot-server flag",
		)
	}

	if args.SnapshotRemoteLocation == "" {
		return fmt.Errorf(
			"snapshot remote locatino not provided: provide value with the --snapshot-remote-location flag",
		)
	}

	if args.SnapshotServerUser == "" {
		return fmt.Errorf(
			"snapshot server user not provided: provide value with the --snapshot-server-user flag",
		)
	}

	if !tools.FileExists(args.SnapshotServerKeyFile) {
		return fmt.Errorf(
			"snapshot server key does not exist: provide value with the --snapshot-server-key flag",
		)
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
		StringVar(&loadSnapshotArgs.VegacapsuleHome, "vegacapsule-home", DefaultVegacapsuleHome, "The custom vegacapsule home")
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
		filepath.Join(vegacapsuleHome, "vega", "node0", SnapshotDBSubPath),
	}
}

func runLoadSnapshot(
	logger *zap.Logger,
	snapshotServerHost, snapshotServerUser, snapshotServerKeyFile, snapshotRemoteLocation string,
	vegaBinary, vegacapsuleHome string,
) error {
	localSnapshotsDbPaths := localSnapshotPaths(vegacapsuleHome)
	for _, snapshotPath := range localSnapshotsDbPaths {
		logger.Info(fmt.Sprintf("Creating folder for snapshot DB locally: %s", snapshotPath))
		if err := os.MkdirAll(snapshotPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create snapshot db(%s): %w", snapshotPath, err)
		}
	}

	tempDir, err := os.MkdirTemp("", "devopstools-snapshot-compatibility")
	if err != nil {
		return fmt.Errorf("failed to create temporary dir to download snapshot db: %w", err)
	}
	// defer os.RemoveAll(tempDir)

	logger.Info(
		"Downloading the snapshot db",
		zap.String("server", fmt.Sprintf("%s@%s", snapshotServerUser, snapshotServerHost)),
		zap.String("source path", snapshotRemoteLocation),
		zap.String("destination", tempDir),
	)
	if err := ssh.Download(
		snapshotServerHost,
		snapshotServerUser,
		snapshotServerKeyFile,
		snapshotRemoteLocation,
		tempDir,
		logger,
	); err != nil {
		return fmt.Errorf("failed to download snapshot db: %w", err)
	}
	logger.Info("Snapshot database downloaded")

	logger.Info("Selecting block height for the restart")
	restartHeight, err := selectSnapshotHeight(vegaBinary, filepath.Join(tempDir, "snapshots"))
	if err != nil {
		return fmt.Errorf("failed to select height for network start: %w", err)
	}
	logger.Info("Selected height for restart", zap.Int("height", restartHeight))

	// run vega tools snapshot --output json ...
	// select 2-nd or 3-rd newest snapshot to load
	// update the config or vega core
	// convert selected snapshot to json
	// move the converted snapshot to new location(flag??)
	return nil
}

// {"snapshots":[{"height":5648600,"version":18830,"size":71,"hash":"80bedacff88b8069f3abfff49d42930c553632ce48ecc6f675339955edd8f24a"},
type CoreToolsSnapshot struct {
	Height  int    `json:"height"`
	Version int    `json:"version"`
	Size    int    `json:"size"`
	Hash    string `json:"hash"`
}

type CoreToolsSnapshots struct {
	Snapshots []CoreToolsSnapshot `json:"snapshots"`
}

func selectSnapshotHeight(vegaBinary, snapshotDbLocation string) (int, error) {
	result := &CoreToolsSnapshots{}

	args := []string{
		"tools",
		"snapshot",
		"--db-path", snapshotDbLocation,
		"--output", "json",
	}
	if _, err := tools.ExecuteBinary(vegaBinary, args, result); err != nil {
		return 0, fmt.Errorf("failed to execute vega %v: %w", args, err)
	}
	snapshotList := result.Snapshots
	sort.Slice(snapshotList, func(i, j int) bool {
		return snapshotList[i].Height < snapshotList[j].Height
	})
	if len(snapshotList) < 2 {
		return 0, fmt.Errorf(
			"not enough snapshots: expected at least 2 snapshots, %d got",
			len(snapshotList),
		)
	}
	return snapshotList[1].Height, nil
}
