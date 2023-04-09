package vegachain

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
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

func BackupChainData(logger *zap.Logger, sess *session.Session, destinationPath, destinationBucket, snapshotDestinationPath string) error {
	syncManager := s3sync.New(sess)

	vegaHomeS3DestinationPath := fmt.Sprintf("s3://%s/%s/vega_home", destinationBucket, destinationPath)
	visorHomeS3DestinationPath := fmt.Sprintf("s3://%s/%s/vegavisor_home", destinationBucket, destinationPath)
	tendermintHomeS3DestinationPath := fmt.Sprintf("s3://%s/%s/tendermint_home", destinationBucket, destinationPath)
	visorBackupPerformed := false

	s3sync.SetLogger(&S3ManagerLogger{logger: logger})

	// TODO: Check if possible parallelism for each dir
	logger.Info(
		fmt.Sprintf("Starting backup for %s", VegaHome),
		zap.String("source", VegaHome),
		zap.String("destination", vegaHomeS3DestinationPath),
	)
	syncManager.Sync(VegaHome, vegaHomeS3DestinationPath)
	logger.Info(
		fmt.Sprintf("Backup for %s finished", VegaHome),
		zap.String("source", VegaHome),
		zap.String("destination", vegaHomeS3DestinationPath),
	)

	if tools.FileExists(VisorHome) {
		visorBackupPerformed = true
		logger.Info(
			fmt.Sprintf("Starting backup for %s", VisorHome),
			zap.String("source", VisorHome),
			zap.String("destination", visorHomeS3DestinationPath),
		)
		syncManager.Sync(VisorHome, visorHomeS3DestinationPath)
		logger.Info(
			fmt.Sprintf("Backup for %s finished", VisorHome),
			zap.String("source", VisorHome),
			zap.String("destination", visorHomeS3DestinationPath),
		)
	}

	logger.Info(
		fmt.Sprintf("Starting backup for %s", TendermintHome),
		zap.String("source", TendermintHome),
		zap.String("destination", tendermintHomeS3DestinationPath),
	)
	syncManager.Sync(TendermintHome, tendermintHomeS3DestinationPath)
	logger.Info(
		fmt.Sprintf("Backup for %s finished", TendermintHome),
		zap.String("source", TendermintHome),
		zap.String("destination", tendermintHomeS3DestinationPath),
	)

	logger.Info("Creating backup snapshot")

	// go func() {
	// 	// TODO: Snapsgot
	// }()

	return nil
}
