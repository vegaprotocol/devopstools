package remote

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/ssh"
	"go.uber.org/zap"
)

type DownloadSnapshotArgs struct {
	*RemoteArgs

	VegaHome    string
	Destination string
}

var downloadSnapshotArgs DownloadSnapshotArgs

var downloadSnapshotCmd = &cobra.Command{
	Use:   "download-snapshot",
	Short: "Download snapshot from the server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := downloadSnapshot(
			downloadSnapshotArgs.Logger,
			downloadSnapshotArgs.ServerHost,
			downloadSnapshotArgs.ServerUser,
			downloadSnapshotArgs.ServerKey,
			downloadSnapshotArgs.VegaHome,
			downloadSnapshotArgs.Destination,
		); err != nil {
			downloadSnapshotArgs.Logger.Fatal("failed to download snapshot", zap.Error(err))
		}
	},
}

func init() {
	downloadSnapshotArgs.RemoteArgs = &remoteArgs

	downloadSnapshotCmd.PersistentFlags().StringVar(
		&downloadSnapshotArgs.VegaHome,
		"remote-vega-home-path",
		"/home/vega/vega_home",
		"Path to the vega-home on the remote server")
	downloadSnapshotCmd.PersistentFlags().StringVar(
		&downloadSnapshotArgs.Destination,
		"destination-path",
		"",
		"Local path where snapshot is downloaded")

	if err := downloadSnapshotCmd.MarkPersistentFlagRequired("destination-path"); err != nil {
		panic(err)
	}

	RemoteCmd.AddCommand(downloadSnapshotCmd)
}

func downloadSnapshot(logger *zap.Logger, sshHost, sshUser, keyPath, vegaHome, destinationPath string) error {
	// not using filepath.Join because path must be Linux style despite the local system
	dirId := time.Now().Unix()
	remoteSnapshotPath := fmt.Sprintf("%s/%s", vegaHome, "state/node/snapshots")
	remoteTmpLocation := fmt.Sprintf("%s/%s-%d", "/tmp", "/snapshots", dirId)

	logger.Info("Creating ssh client")
	client, err := ssh.GetSSHConnection(sshHost, sshUser, keyPath)
	if err != nil {
		return fmt.Errorf("failed to create ssh client: %w", err)
	}

	logger.Info(
		"Coying snapshot on remote machine to tmp location",
		zap.String("source", fmt.Sprintf("%s:%s", sshHost, remoteSnapshotPath)),
		zap.String("destination", fmt.Sprintf("%s:%s", sshHost, remoteTmpLocation)),
	)
	out, err := ssh.RunCommandWithClient(client, fmt.Sprintf("sudo cp -r '%s' '%s'", remoteSnapshotPath, remoteTmpLocation))
	if err != nil {
		return fmt.Errorf(
			"failed to copy remote snapshots from '%s' to '%s', stdout: '%s': %w",
			remoteSnapshotPath,
			remoteTmpLocation,
			out,
			err,
		)
	}
	logger.Info("Copy done")

	logger.Info(
		"Downloading remote snapshots",
		zap.String("source", fmt.Sprintf("%s:%s", sshHost, remoteTmpLocation)),
		zap.String("destination", destinationPath),
	)
	if err := ssh.Download(sshHost, sshUser, keyPath, remoteTmpLocation, destinationPath, logger); err != nil {
		return fmt.Errorf("failed to download snapshot from remote: %w", err)
	}
	logger.Info("Downloading finished")

	downloadedPath := filepath.Join(destinationPath, "/", fmt.Sprintf("snapshots-%d", dirId))
	newDestinationPath := filepath.Join(destinationPath, "snapshots")
	logger.Info("Removing old directory", zap.String("path", newDestinationPath))
	if err := os.RemoveAll(newDestinationPath); err != nil {
		return fmt.Errorf("failed to remove old downloaded snapshot from '%s': %w", newDestinationPath, err)
	}

	logger.Info("Renaming downloaded file", zap.String("source", downloadedPath), zap.String("destination", newDestinationPath))
	if err := os.Rename(downloadedPath, newDestinationPath); err != nil {
		return fmt.Errorf("failed to rename dir from '%s' to '%s': %w", downloadedPath, newDestinationPath, err)
	}

	return nil
}
