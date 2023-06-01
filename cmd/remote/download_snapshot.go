package remote

import (
	"fmt"

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
	remoteSnapshotPath := fmt.Sprintf("%s/%s", vegaHome, "state/node/snapshots")

	logger.Info(
		"Downloading remote snapshots",
		zap.String("source", fmt.Sprintf("%s:%s", sshHost, remoteSnapshotPath)),
		zap.String("destination", destinationPath),
	)
	if err := ssh.Download(sshHost, sshUser, keyPath, remoteSnapshotPath, destinationPath, logger); err != nil {
		return fmt.Errorf("failed to download snapshot from remote: %w", err)
	}
	logger.Info("Downloading finished")

	return nil
}
