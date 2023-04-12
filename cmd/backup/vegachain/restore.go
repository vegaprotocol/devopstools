package vegachain

import (
	"fmt"
	"os"
	"path/filepath"

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
	vegaHomeSnapshot := filepath.Join(s3SnapshotLocation, "vega_home")
	visorHomeSnapshot := filepath.Join(s3SnapshotLocation, "vegavisor_home")
	tendermintHomeSnapshot := filepath.Join(s3SnapshotLocation, "tendermint_home")

	// TODO: Think to add parallel
	if err := S3Sync(logger, s3CmdBinary, vegaHomeSnapshot, VegaHome); err != nil {
		return fmt.Errorf("failed to download snapshot for vega home: %w", err)
	}

	if withVisorHome {
		if err := S3Sync(logger, s3CmdBinary, visorHomeSnapshot, VisorHome); err != nil {
			return fmt.Errorf("failed to download snapshot for vegavisor home: %w", err)
		}
	}

	if err := S3Sync(logger, s3CmdBinary, tendermintHomeSnapshot, TendermintHome); err != nil {
		return fmt.Errorf("failed to download snapshot for tendermint home: %w", err)
	}

	return nil
}
