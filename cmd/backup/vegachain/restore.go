package vegachain

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
)

func RemoveLocalChainData(logger *zap.Logger) error {
	if tools.FileExists(VegaHome) {
		logger.Info("Removing local vega home", zap.String("path", VegaHome))
		if err := tools.RemoveDirectoryContents(VegaHome); err != nil {
			return fmt.Errorf("failed to remove vega home(%s): %w", VegaHome, err)
		}
	}

	if tools.FileExists(TendermintHome) {
		logger.Info("Removing local tendermint home", zap.String("path", TendermintHome))
		if err := tools.RemoveDirectoryContents(TendermintHome); err != nil {
			return fmt.Errorf("failed to remove tendermint home(%s): %w", TendermintHome, err)
		}
	}

	if tools.FileExists(VisorHome) {
		logger.Info("Removing local vegavisor home", zap.String("path", VisorHome))
		if err := tools.RemoveDirectoryContents(VisorHome); err != nil {
			return fmt.Errorf("failed to remove vegavisor home(%s): %w", VisorHome, err)
		}
	}

	return nil
}

func RestoreChainData(logger *zap.Logger, s3CmdBinary, s3SnapshotLocation string, withVisorHome bool) error {
	vegaHomeSnapshot := S3Join(s3SnapshotLocation, "vega_home")
	visorHomeSnapshot := S3Join(s3SnapshotLocation, "vegavisor_home")
	tendermintHomeSnapshot := S3Join(s3SnapshotLocation, "tendermint_home")

	logger.Info("Restoring vega home from remote snapshot", zap.String("source", vegaHomeSnapshot), zap.String("destination", VegaHome))
	// TODO: Think to add parallel
	if err := S3Sync(logger, s3CmdBinary, vegaHomeSnapshot, HomesParentDir); err != nil {
		return fmt.Errorf("failed to download snapshot for vega home: %w", err)
	}

	if withVisorHome {
		logger.Info("Restoring vegavisor home from remote snapshot", zap.String("source", visorHomeSnapshot), zap.String("destination", VisorHome))
		if err := S3Sync(logger, s3CmdBinary, visorHomeSnapshot, HomesParentDir); err != nil {
			return fmt.Errorf("failed to download snapshot for vegavisor home: %w", err)
		}
	}

	logger.Info("Restoring tendermint home from remote snapshot", zap.String("source", tendermintHomeSnapshot), zap.String("destination", TendermintHome))
	if err := S3Sync(logger, s3CmdBinary, tendermintHomeSnapshot, HomesParentDir); err != nil {
		return fmt.Errorf("failed to download snapshot for tendermint home: %w", err)
	}

	logger.Info(fmt.Sprintf("Changing ownership for the %s directory to %s:%s", VegaHome, VegaUser, VegaGroup))
	if err := tools.ChownR(VegaHome, VegaUser, VegaGroup); err != nil {
		return fmt.Errorf("failed to change owner for %s: %w", VegaHome, err)
	}

	if withVisorHome {
		logger.Info(fmt.Sprintf("Changing ownership for the %s directory to %s:%s", VisorHome, VegaUser, VegaGroup))
		if err := tools.ChownR(VisorHome, VegaUser, VegaGroup); err != nil {
			return fmt.Errorf("failed to change owner for %s: %w", VisorHome, err)
		}
	}

	logger.Info(fmt.Sprintf("Changing ownership for the %s directory to %s:%s", TendermintHome, VegaUser, VegaGroup))
	if err := tools.ChownR(TendermintHome, VegaUser, VegaGroup); err != nil {
		return fmt.Errorf("failed to change owner for %s: %w", TendermintHome, err)
	}

	return nil
}
