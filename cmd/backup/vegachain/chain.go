package vegachain

import (
	"fmt"

	"github.com/seqsense/s3sync"
	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
)

const (
	VegaHome       = "/home/vega/vega_home"
	VisorHome      = "/home/vega/vegavisor_home"
	TendermintHome = "/home/vega/tendermint_home"
)

type S3ManagerLogger struct {
	logger *zap.Logger
}

func (l *S3ManagerLogger) Log(v ...interface{}) {
	// We want to log files only when debug mode is disabled
	if l.logger.Level() == zap.DebugLevel {
		l.logger.Info(fmt.Sprint(v...))
	}
}

func (l *S3ManagerLogger) Logf(format string, v ...interface{}) {
	if l.logger.Level() == zap.DebugLevel {
		args := []interface{}{format}
		args = append(args, v)

		l.logger.Info(fmt.Sprint(args...))
	}
}

func BackupChainData(logger *zap.Logger, s3CmdBinary string, destinationPath, destinationBucket, snapshotDestinationPath string) error {
	debug := logger.Level() == zap.DebugLevel

	vegaHomeS3DestinationPath := fmt.Sprintf("s3://%s/%s/vega_home/", destinationBucket, destinationPath)
	visorHomeS3DestinationPath := fmt.Sprintf("s3://%s/%s/vegavisor_home/", destinationBucket, destinationPath)
	tendermintHomeS3DestinationPath := fmt.Sprintf("s3://%s/%s/tendermint_home/", destinationBucket, destinationPath)

	s3sync.SetLogger(&S3ManagerLogger{logger: logger})

	// TODO: Check if possible parallelism for each dir
	logger.Info(
		fmt.Sprintf("Starting backup for %s", VegaHome),
		zap.String("source", VegaHome),
		zap.String("destination", vegaHomeS3DestinationPath),
	)
	if err := S3Sync(s3CmdBinary, VegaHome, vegaHomeS3DestinationPath, debug); err != nil {
		return fmt.Errorf("failed to backup vega data: %w", err)
	}
	logger.Info(
		fmt.Sprintf("Backup for %s finished", VegaHome),
		zap.String("source", VegaHome),
		zap.String("destination", vegaHomeS3DestinationPath),
	)

	if tools.FileExists(VisorHome) {
		logger.Info(
			fmt.Sprintf("Starting backup for %s", VisorHome),
			zap.String("source", VisorHome),
			zap.String("destination", visorHomeS3DestinationPath),
		)
		if err := S3Sync(s3CmdBinary, VisorHome, visorHomeS3DestinationPath, debug); err != nil {
			return fmt.Errorf("failed to backup visor data: %w", err)
		}
		logger.Info(
			fmt.Sprintf("Backup for %s finished", VisorHome),
			zap.String("source", VisorHome),
			zap.String("destination", visorHomeS3DestinationPath),
		)
	} else {
		logger.Info("Backup for vegavisor not required")
	}

	logger.Info(
		fmt.Sprintf("Starting backup for %s", TendermintHome),
		zap.String("source", TendermintHome),
		zap.String("destination", tendermintHomeS3DestinationPath),
	)
	if err := S3Sync(s3CmdBinary, TendermintHome, tendermintHomeS3DestinationPath, debug); err != nil {
		return fmt.Errorf("failed to backup tendermint data: %w", err)
	}
	logger.Info(
		fmt.Sprintf("Backup for %s finished", TendermintHome),
		zap.String("source", TendermintHome),
		zap.String("destination", tendermintHomeS3DestinationPath),
	)

	snapshotSource := fmt.Sprintf("s3://%s/%s/", destinationBucket, destinationPath)
	snapshotDestination := fmt.Sprintf("s3://%s/%s/", destinationBucket, snapshotDestinationPath)
	logger.Info(
		"Creating vega chain backup snapshot",
		zap.String("source", snapshotSource),
		zap.String("destination", snapshotDestination),
	)
	if err := S3Sync(s3CmdBinary, snapshotSource, snapshotDestination, debug); err != nil {
		return fmt.Errorf("failed to create backup snapshot: %w", err)
	}
	logger.Info(
		"Vega chain backup snapshot finished",
		zap.String("source", snapshotSource),
		zap.String("destination", snapshotDestination),
	)

	return nil
}