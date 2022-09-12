package ssh

import (
	"fmt"
	"io"
	"os/exec"
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
		rsyncCmd := exec.Command(
			"rsync",
			"-avz",
			"--quiet",
			"-e", fmt.Sprintf("ssh -i %s", sshPrivateKeyfile),
			"--rsync-path", "sudo rsync",
			fmt.Sprintf("%s@%s:%s", sshUsername, serverHost, srcFilepath),
			dstFilepath,
		)
		out, err := rsyncCmd.Output()
		if err == nil {
			return nil
		}
		data, _ := io.ReadAll(rsyncCmd.Stdin)
		logger.Error("Failed to download file",
			zap.String("out", string(out)), zap.String("stdin", string(data)), zap.Error(err))
	}
	return fmt.Errorf("Failed to download file")
}
