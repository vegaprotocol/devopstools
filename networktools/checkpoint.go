package networktools

import (
	"fmt"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

func (network *NetworkTools) GetCheckpointPath(vegaHome string) string {
	if len(vegaHome) > 0 {
		return fmt.Sprintf("%s/state/node/checkpoints", vegaHome)
	}
	return "/home/vega/.local/state/vega/node/checkpoints"
}

type Checkpoint struct {
	ServerHost string `json:"server"`
	Filepath   string `json:"file"`
}

func (network *NetworkTools) FindLatestCheckpoint(
	vegaHome string,
	sshUsername string,
	sshPrivateKeyfile string,
) (result *Checkpoint, err error) {
	checkpointPath := network.GetCheckpointPath(vegaHome)
	command := fmt.Sprintf("sudo find '%s' -name '20*.cp' |sort |tail -1", checkpointPath)

	sshResults := network.RunCommandOnEveryNode(sshUsername, sshPrivateKeyfile, command)

	for _, sshResult := range sshResults {
		if sshResult.Err != nil {
			network.logger.Debug("failed to get latest checkpoint from host", zap.String("host", sshResult.Host), zap.Error(sshResult.Err))
			continue
		}
		output := strings.TrimSpace(sshResult.Output)
		if len(output) == 0 {
			network.logger.Debug("no checkpoint on host", zap.String("host", sshResult.Host), zap.Error(sshResult.Err))
			continue
		}
		newResult := Checkpoint{
			ServerHost: sshResult.Host,
			Filepath:   output,
		}
		if result == nil {
			result = &newResult
		} else {
			currFilename := filepath.Base(result.Filepath)
			toCheckFilename := filepath.Base(newResult.Filepath)
			if toCheckFilename > currFilename {
				network.logger.Debug("found better results", zap.String("previous", currFilename), zap.String("new", toCheckFilename))
				result = &newResult
			}
		}
	}
	return result, nil
}
