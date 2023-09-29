package snapshotcompatibility

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/storage"
	"go.uber.org/zap"
	sshlib "golang.org/x/crypto/ssh"

	"github.com/vegaprotocol/devopstools/ssh"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/vegacapsule"
	"github.com/vegaprotocol/devopstools/vegacmd"
)

var (
	VegaSnapshotPath   = filepath.Join("state", "node", "snapshots")
	VegaCoreConfigPath = filepath.Join("config", "node", "config.toml")
)

type LoadSnapshotArgs struct {
	*SnapshotCompatibilityArgs

	VegacapsuleHome   string
	VegacapsuleBinary string
	VegaBinary        string

	SnapshotServerHost     string
	SnapshotRemoteLocation string
	SnapshotServerUser     string
	SnapshotServerKeyFile  string
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

	currentUser, _ := tools.WhoAmI()
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.VegacapsuleHome, "vegacapsule-home", "", "The custom vegacapsule home")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.VegaBinary, "vega-binary", "vega", "Path to the vega executable")
	loadSnapshotCmd.PersistentFlags().
		StringVar(&loadSnapshotArgs.VegacapsuleBinary, "vegacapsule-binary", "vegacapsule", "Path to the vegacapsule executable")
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
	snapshotServerHost, snapshotServerUser, snapshotServerKeyFile, snapshotRemoteLocation string,
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

	snapshotRemoteTempLocation := fmt.Sprintf("/tmp/snapshots-db-%d/snapshots", time.Now().UnixMicro())

	logger.Info("Creating SSH client",
		zap.String("host", snapshotServerHost),
		zap.String("user", snapshotServerUser),
		zap.String("keyfile", snapshotServerKeyFile))

	sshClient, err := ssh.GetSSHConnection(snapshotServerHost, snapshotServerUser, snapshotServerKeyFile)
	if err != nil {
		return fmt.Errorf("failed to create ssh connection client: %w", err)
	}

	err = tools.RetryRun(3, 5*time.Second, func() error {
		logger.Info("Trying to copy snapshot.db on remote server into temporary location")
		return copySnapshotOnRemote(logger,
			sshClient,
			snapshotRemoteLocation,
			snapshotRemoteTempLocation)
	})
	if err != nil {
		return fmt.Errorf("failed to copy snapshot on remote temp location: %w", err)
	}
	defer remoteCleanup(logger, sshClient, filepath.Dir(snapshotRemoteTempLocation))

	tempDir, err := os.MkdirTemp("", "devopstools-snapshot-compatibility")
	if err != nil {
		return fmt.Errorf("failed to create temporary dir to download snapshot db: %w", err)
	}
	defer os.RemoveAll(tempDir)

	logger.Info(
		"Downloading the snapshot db",
		zap.String("server", fmt.Sprintf("%s@%s", snapshotServerUser, snapshotServerHost)),
		zap.String("source path", snapshotRemoteTempLocation),
		zap.String("destination", tempDir),
	)

	err = tools.RetryRun(3, 5*time.Second, func() error {
		logger.Info("Trying to download snapshot-db")
		return ssh.Download(
			snapshotServerHost,
			snapshotServerUser,
			snapshotServerKeyFile,
			snapshotRemoteTempLocation,
			tempDir,
			logger,
		)
	})
	if err != nil {
		return fmt.Errorf("failed to download snapshot db: %w", err)
	}
	logger.Info("Snapshot database downloaded")

	logger.Info("Selecting block height for the restart")
	restartHeight, err := selectSnapshotHeight(vegaBinary, filepath.Join(tempDir, "snapshots"))
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
		snapshotSource := filepath.Join(tempDir, "snapshots")

		logger.Info("Calling unsafe reset all for core", zap.String("node", validatorHomePath))
		if _, err := tools.ExecuteBinary(vegaBinary, []string{"unsafe_reset_all", "--home", validatorHomePath}, nil); err != nil {
			return fmt.Errorf("failed to unsafe reset all for home %s: %w", validatorHomePath, err)
		}
		logger.Info("Unsafe reset all succesfull")

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
	}
	return nil
}

func remoteCleanup(logger *zap.Logger, sshClient *sshlib.Client, filePath string) error {
	cleanupCommand := fmt.Sprintf("rm -fr %s || echo 'OK'", filePath)
	logger.Info("Cleaning the remote server up",
		zap.String("filepath", filePath))

	if stdout, err := ssh.RunCommandWithClient(sshClient, cleanupCommand); err != nil {
		return fmt.Errorf("failed to cleanup on remote: output: %s: %s", stdout, err)
	}

	return nil
}

func copySnapshotOnRemote(logger *zap.Logger,
	sshClient *sshlib.Client,
	source, destination string,
) error {
	cleanupCommand := fmt.Sprintf("rm -fr %s || echo 'OK'", destination)
	logger.Info("Cleaning the destination up",
		zap.String("filepath", destination))

	if stdout, err := ssh.RunCommandWithClient(sshClient, cleanupCommand); err != nil {
		return fmt.Errorf("failed to cleanup on remote: output: %s: %s", stdout, err)
	}

	mkdirCommand := fmt.Sprintf("mkdir -p %s", filepath.Dir(destination))
	logger.Info("Ensuring destination directory exists",
		zap.String("command", mkdirCommand))

	if stdout, err := ssh.RunCommandWithClient(sshClient, mkdirCommand); err != nil {
		return fmt.Errorf("failed to ensure destination directory exists: output: %s: %s", stdout, err)
	}

	copyCommand := fmt.Sprintf("cp -r %s %s", source, destination)

	logger.Info("Copying snapshot db to temp location",
		zap.String("command", copyCommand))
	if stdout, err := ssh.RunCommandWithClient(sshClient, copyCommand); err != nil {
		return fmt.Errorf("failed to copy snapshot from %s to %s on remote: output: %s: %s", source, destination, stdout, err)
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
