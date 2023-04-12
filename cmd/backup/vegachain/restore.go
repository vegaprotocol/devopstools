package vegachain

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
)

func RemoveLocalChainData(logger *zap.Logger) error {
	if tools.FileExists(VegaHome) {
		logger.Info("Removing local vega home", zap.String("path", VegaHome))
		if err := os.RemoveAll(VegaHome); err != nil {
			return fmt.Errorf("failed to remove vega home(%s): %w", VegaHome, err)
		}
	}

	if tools.FileExists(TendermintHome) {
		logger.Info("Removing local tendermint home", zap.String("path", TendermintHome))
		if err := os.RemoveAll(TendermintHome); err != nil {
			return fmt.Errorf("failed to remove tendermint home(%s): %w", TendermintHome, err)
		}
	}

	if tools.FileExists(VisorHome) {
		logger.Info("Removing local vegavisor home", zap.String("path", VisorHome))
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

	logger.Info("Restoring vega home from remote snapshot", zap.String("source", vegaHomeSnapshot), zap.String("destination", VegaHome))
	// TODO: Think to add parallel
	if err := S3Sync(logger, s3CmdBinary, vegaHomeSnapshot, VegaHome); err != nil {
		return fmt.Errorf("failed to download snapshot for vega home: %w", err)
	}

	if withVisorHome {
		logger.Info("Restoring vegavisor home from remote snapshot", zap.String("source", visorHomeSnapshot), zap.String("destination", VisorHome))
		if err := S3Sync(logger, s3CmdBinary, visorHomeSnapshot, VisorHome); err != nil {
			return fmt.Errorf("failed to download snapshot for vegavisor home: %w", err)
		}
	}

	logger.Info("Restoring tendermint home from remote snapshot", zap.String("source", tendermintHomeSnapshot), zap.String("destination", TendermintHome))
	if err := S3Sync(logger, s3CmdBinary, tendermintHomeSnapshot, TendermintHome); err != nil {
		return fmt.Errorf("failed to download snapshot for tendermint home: %w", err)
	}

	if err := tools.ChownR(VegaHome, VegaUser, VegaGroup); err != nil {
		return fmt.Errorf("failed to change owner for %s: %w", VegaHome, err)
	}

	if withVisorHome {
		if err := tools.ChownR(VisorHome, VegaUser, VegaGroup); err != nil {
			return fmt.Errorf("failed to change owner for %s: %w", VisorHome, err)
		}
	}

	if err := tools.ChownR(TendermintHome, VegaUser, VegaGroup); err != nil {
		return fmt.Errorf("failed to change owner for %s: %w", TendermintHome, err)
	}

	return nil
}
