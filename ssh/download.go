package ssh

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"go.uber.org/zap"
)

func Download(
	serverHost string,
	sshUsername string,
	sshPrivateKeyfile string,
	srcFilepath string,
	dstFilepath string,
	logger *zap.Logger,
) error {
	for i := 0; i < 3; i++ {
		if i > 0 {
			logger.Info("Retry download in a second")
			time.Sleep(time.Second)
		}

		sshFlags := []string{
			"-i", sshPrivateKeyfile,
			"-o", "StrictHostKeyChecking=no",
		}

		rsyncCmd := exec.Command(
			"rsync",
			"-avz",
			"--quiet",
			"-e", fmt.Sprintf("ssh %s", strings.Join(sshFlags, " ")),
			"--rsync-path", "sudo rsync",
			fmt.Sprintf("%s@%s:%s", sshUsername, serverHost, srcFilepath),
			dstFilepath,
		)
		if err := rsyncCmd.Run(); err == nil {
			return nil
		} else {
			logger.Error("Failed to download file", zap.Error(err))
		}
	}
	return fmt.Errorf("Failed to download file")
}
