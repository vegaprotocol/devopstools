package ssh

import (
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/tools"
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

		args := []string{
			"-avz",
			"--quiet",
			"-e", fmt.Sprintf("ssh -i %s", sshPrivateKeyfile),
			"--rsync-path", "sudo rsync",
			fmt.Sprintf("%s@%s:%s", sshUsername, serverHost, srcFilepath),
			dstFilepath,
		}

		out, err := tools.ExecuteBinary("rsync", args, nil)

		if err != nil {
			logger.Error("Failed to download file", zap.Error(err), zap.String("stdout", string(out)))
			time.Sleep(1 * time.Second)
		} else {
			return nil
		}
	}

	return fmt.Errorf("failed to download file")
}
