package snapshotcompatibility

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vegaprotocol/devopstools/ssh"
	"github.com/vegaprotocol/devopstools/tools"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	sshlib "golang.org/x/crypto/ssh"
)

var (
	VegaSnapshotPath   = filepath.Join("state", "node", "snapshots")
	VegaCoreConfigPath = filepath.Join("config", "node", "config.toml")
	RemoteVegaPath     = "/home/vega/vegavisor_home/current/vega"
)

type DownloadMainnetSnapshotArgs struct {
	*Args

	RemoteVegaBinary     string
	TemporaryDestination string

	SnapshotServerHost     string
	SnapshotRemoteLocation string
	SnapshotServerUser     string
	SnapshotServerKeyFile  string
}

func (args DownloadMainnetSnapshotArgs) Check() error {
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

var downloadMainnetSnapshotArgs DownloadMainnetSnapshotArgs

var downloadMainnetSnapshotCmd = &cobra.Command{
	Use:   "download-mainnet-snapshot",
	Short: "Download the snapshot from the mainnet and put it in the local folder",
	Run: func(cmd *cobra.Command, args []string) {
		if err := downloadMainnetSnapshotArgs.Check(); err != nil {
			loadSnapshotArgs.Logger.Fatal("invalid input", zap.Error(err))
			return
		}

		if err := runDownloadMainnetSnapshot(downloadMainnetSnapshotArgs.Logger,
			downloadMainnetSnapshotArgs.SnapshotServerHost,
			downloadMainnetSnapshotArgs.SnapshotServerUser,
			downloadMainnetSnapshotArgs.SnapshotServerKeyFile,
			downloadMainnetSnapshotArgs.SnapshotRemoteLocation,
			downloadMainnetSnapshotArgs.RemoteVegaBinary,
			downloadMainnetSnapshotArgs.TemporaryDestination,
		); err != nil {
			downloadMainnetSnapshotArgs.Logger.Fatal(
				"failed to download mainnet snapshot",
				zap.Error(err),
			)

			return
		}
	},
}

func init() {
	downloadMainnetSnapshotArgs.Args = &snapshotCompatibilityArgs

	currentUser, _ := tools.WhoAmI()
	downloadMainnetSnapshotCmd.PersistentFlags().
		StringVar(
			&downloadMainnetSnapshotArgs.RemoteVegaBinary,
			"remote-vega-binary",
			RemoteVegaPath,
			"Path to the vega executable on remote machine",
		)
	downloadMainnetSnapshotCmd.PersistentFlags().
		StringVar(
			&downloadMainnetSnapshotArgs.TemporaryDestination,
			"local-temporary-destination",
			"/tmp/mainnet-snapshot",
			"Path on the local machine, where we download the mainnet snapshot",
		)
	downloadMainnetSnapshotCmd.PersistentFlags().
		StringVar(
			&downloadMainnetSnapshotArgs.SnapshotServerHost,
			"snapshot-server",
			"",
			"The source server for the snapshot",
		)
	downloadMainnetSnapshotCmd.PersistentFlags().
		StringVar(
			&downloadMainnetSnapshotArgs.SnapshotServerUser,
			"snapshot-server-user",
			currentUser,
			"The SSH user for the snapshot server",
		)
	downloadMainnetSnapshotCmd.PersistentFlags().
		StringVar(
			&downloadMainnetSnapshotArgs.SnapshotRemoteLocation,
			"snapshot-remote-location",
			"/home/vega/vega_home/state/node/snapshots",
			"The location where the snapshot is on the remote server",
		)
	downloadMainnetSnapshotCmd.PersistentFlags().
		StringVar(
			&downloadMainnetSnapshotArgs.SnapshotServerKeyFile,
			"snapshot-server-key-file",
			filepath.Join(tools.CurrentUserHomePath(), ".ssh", "id_rsa"),
			"The SSH private key used to authenticate user",
		)
}

func runDownloadMainnetSnapshot(
	logger *zap.Logger,
	snapshotServerHost, snapshotServerUser, snapshotServerKeyFile, snapshotRemoteLocation string,
	RemoteVegaBinary, snapshotTempDestination string,
) error {
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
			RemoteVegaBinary,
			sshClient,
			snapshotRemoteLocation,
			snapshotRemoteTempLocation)
	})
	if err != nil {
		return fmt.Errorf("failed to copy snapshot on remote temp location: %w", err)
	}
	defer remoteCleanup(logger, sshClient, filepath.Dir(snapshotRemoteTempLocation))

	logger.Info("Cleaning the local temp folder up", zap.String("path", snapshotTempDestination))
	if err := os.RemoveAll(snapshotTempDestination); err != nil {
		return fmt.Errorf("failed to cleanup temp destination(%s): %w", snapshotTempDestination, err)
	}

	logger.Info("Ensuring local temp path exist", zap.String("path", snapshotTempDestination))
	if err := os.MkdirAll(snapshotRemoteTempLocation, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temporary dir to download snapshot db: %w", err)
	}

	logger.Info(
		"Downloading the snapshot db",
		zap.String("server", fmt.Sprintf("%s@%s", snapshotServerUser, snapshotServerHost)),
		zap.String("source path", snapshotRemoteTempLocation),
		zap.String("destination", snapshotTempDestination),
	)

	err = tools.RetryRun(3, 5*time.Second, func() error {
		logger.Info("Trying to download snapshot-db")
		return ssh.Download(
			snapshotServerHost,
			snapshotServerUser,
			snapshotServerKeyFile,
			snapshotRemoteTempLocation,
			snapshotTempDestination,
			logger,
		)
	})
	if err != nil {
		return fmt.Errorf("failed to download snapshot db: %w", err)
	}
	logger.Info("Snapshot database downloaded")
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
	RemoteVegaBinary string,
	sshClient *sshlib.Client,
	source, destination string,
) error {
	cleanupCommand := fmt.Sprintf("rm -fr %s || echo 'OK'", destination)
	logger.Info("Cleaning the destination up",
		zap.String("filepath", destination))

	if stdout, err := ssh.RunCommandWithClient(sshClient, cleanupCommand); err != nil {
		logger.Error(
			"failed to cleanup on remote",
			zap.String("stdout", stdout),
			zap.Error(err),
		)
		return fmt.Errorf("failed to cleanup on remote: output: %s: %s", stdout, err)
	}

	mkdirCommand := fmt.Sprintf("mkdir -p %s", filepath.Dir(destination))
	logger.Info("Ensuring destination directory exists",
		zap.String("command", mkdirCommand))

	if stdout, err := ssh.RunCommandWithClient(sshClient, mkdirCommand); err != nil {
		logger.Error(
			"failed to ensure destination directory exists",
			zap.String("stdout", stdout),
			zap.Error(err),
		)
		return fmt.Errorf("failed to ensure destination directory exists: output: %s: %s", stdout, err)
	}

	copyCommand := fmt.Sprintf("cp -r %s %s", source, destination)

	logger.Info("Copying snapshot db to temp location",
		zap.String("command", copyCommand))
	if stdout, err := ssh.RunCommandWithClient(sshClient, copyCommand); err != nil {
		logger.Error(
			fmt.Sprintf("failed to copy snapshot from %s to %s on remote: output", source, destination),
			zap.String("stdout", stdout),
			zap.Error(err),
		)
		return fmt.Errorf("failed to copy snapshot from %s to %s on remote: output: %s: %s", source, destination, stdout, err)
	}

	checkCommand := fmt.Sprintf("%s tools snapshot --db-path '%s'", RemoteVegaBinary, destination)
	logger.Info("Checking if copied snapshot was copied correctly",
		zap.String("command", checkCommand))
	if stdout, err := ssh.RunCommandWithClient(sshClient, checkCommand); err != nil {
		logger.Error(
			"failed to check if snapshot was correctly copied on remote",
			zap.String("stdout", stdout),
			zap.Error(err),
		)
		return fmt.Errorf("failed to check if snapshot was correctly copied on remote: output: %s: %w", stdout, err)
	}

	return nil
}
