package vegachain

import (
	"fmt"
	"os"

	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
)

func RemoveLocalChainData() error {
	if tools.FileExists(VegaHome) {
		if err := os.RemoveAll(VegaHome); err != nil {
			return fmt.Errorf("failed to remove vega home(%s): %w", VegaHome, err)
		}
	}

	if tools.FileExists(TendermintHome) {
		if err := os.RemoveAll(TendermintHome); err != nil {
			return fmt.Errorf("failed to remove tendermint home(%s): %w", TendermintHome, err)
		}
	}

	if tools.FileExists(VisorHome) {
		if err := os.RemoveAll(VisorHome); err != nil {
			return fmt.Errorf("failed to remove vegavisor home(%s): %w", VisorHome, err)
		}
	}

	return nil
}

func RestoreChainData(logger *zap.Logger, s3CmdBinary, s3SnapshotLocation string, withVisorHome bool) error {
	// vegaHomeSnapshot = filepath.Join(s3SnapshotLocation, "vega_home")
	// visorHomeSnapshot = filepath.Join(s3SnapshotLocation, "vegavisor_home")
	// tendermintHomeSnapshot = filepath.Join(s3SnapshotLocation, "tendermint_home")

	// S3Sync(logger, s3CmdBinary, s3SnapshotLocation)
	return nil
}
