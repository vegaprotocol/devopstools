package backup

import (
	// "log"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/cmd/backup/pgbackrest"
	"github.com/vegaprotocol/devopstools/cmd/backup/postgresql"
	"github.com/vegaprotocol/devopstools/cmd/backup/systemctl"
	"github.com/vegaprotocol/devopstools/cmd/backup/vegachain"
	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap"
)

type RestoreArgs struct {
	*BackupRootArgs
	localStateFile string

	backupID    string
	s3CmdBinary string

	postgresqlUser       string
	pgBackrestBinary     string
	pgBackrestConfigFile string

	postgresqlConfigFile  string
	pgBackrestFullRestore bool
}

var restoreArgs RestoreArgs

var performRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore vega node from remote S3 bucket",
	Long: `
	TBD
	TBD:
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := DoRestore(restoreArgs); err != nil {
			backupArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	restoreArgs.BackupRootArgs = &backupRootArgs

	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.backupID, "id", "postgres", "The ID of the backup to restore. See")
	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.postgresqlUser, "postgresql-user", "postgres", "The linux username who runs the postgresql")
	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.localStateFile, "local-state-file", "/tmp/vega-backup-state.json", "Local state file for the vega backup")
	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.pgBackrestBinary, "pgbackrest-bin", "pgbackrest", "The binary for pgbackrest")
	performRestoreCmd.PersistentFlags().BoolVar(&restoreArgs.pgBackrestFullRestore, "full", false, "Perform the full restore for postgresql")
	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.pgBackrestConfigFile, "pgbackrest-config-file", "/etc/pgbackrest.conf", "Location of pgbackrest config file")
	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.s3CmdBinary, "s3cmd-bin", "s3cmd", "The binary for s3cmd")

	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.postgresqlConfigFile, "postgresql-config-file", "/etc/postgresql/14/main/postgresql.conf", "The config file for postgresql")

	if err := performRestoreCmd.MarkPersistentFlagRequired("id"); err != nil {
		log.Fatalf("%v\n", err)
	}
	BackupRootCmd.AddCommand(performRestoreCmd)
}

func DoRestore(args RestoreArgs) error {
	args.Logger.Info("Checking if postgresql is running")
	postgresqlRunning := systemctl.IsRunning(args.Logger, "postgresql")

	args.Logger.Info("Loading state from local file", zap.String("file", args.localStateFile))
	state, err := LoadFromLocal(args.localStateFile)
	if err != nil {
		return fmt.Errorf("failed to load backups state: %w", err)
	}

	if len(state.PgBackrestConfig) < 1 {
		return fmt.Errorf("missing pgbackrest config in the state")
	}

	if len(state.Backups) < 1 {
		return fmt.Errorf("no backup found in the local state file: %s", args.localStateFile)
	}

	currentBackup, backupFound := state.Backups[args.backupID]
	if !backupFound {
		return fmt.Errorf("backup %s not found in the state file, run list-backups to see available backups", args.backupID)
	}

	args.Logger.Info("Writing the pgbackrest config file")
	if err := os.WriteFile(args.pgBackrestConfigFile, []byte(state.PgBackrestConfig), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write pgbackrest config from state: %w", err)
	}

	args.Logger.Info("Loading the pgbackrest config file into memory")
	pgBackrestConfig, err := pgbackrest.ReadConfig(args.pgBackrestConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read pgbackrest config: %w", err)
	}

	args.Logger.Info("Verifying s3cmd setup")
	if err := vegachain.CheckS3Setup(args.s3CmdBinary); err != nil {
		return fmt.Errorf("failed to check s3 setup: %w", err)
	}

	if err := vegachain.S3CmdInit(args.s3CmdBinary, vegachain.S3Credentials{
		Endpoint:     pgBackrestConfig.Global.R1S3Endpoint,
		Region:       pgBackrestConfig.Global.R1S3Region,
		AccessKey:    pgBackrestConfig.Global.R1S3Key,
		AccessSecret: pgBackrestConfig.Global.R1S3KeySecret,
	}); err != nil {
		return fmt.Errorf("failed to initialize s3cmd credentials: %w", err)
	}

	args.Logger.Info("Verifying stanza setup")
	if err := pgbackrest.CheckPgBackRestSetup(backupArgs.pgBackrestBinary, pgBackrestConfig); err != nil {
		return fmt.Errorf("failed to check pgbackrest setup: %w", err)
	}

	if !postgresqlRunning {
		args.Logger.Info("Postgresql is not running, stanza-create skipped")
	} else {
		args.Logger.Info("Upgrading stanza")
		if err := pgbackrest.UpgradeStanza(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary); err != nil {
			return fmt.Errorf("failed to upgrade pgbackrest stanza: %w", err)
		}
	}

	args.Logger.Info("Collecting postgresql config")
	postgresqlConfig, err := postgresql.ReadConfig(args.postgresqlConfigFile)
	if err != nil {
		return fmt.Errorf("failed to get postgresql config: %w", err)
	}

	// We will update pgbackrest config with the one from the backuped server
	args.Logger.Info("Creating backup of original pgbackrest config")
	if err := pgbackrest.BackupConfig(args.pgBackrestConfigFile, false); err != nil {
		return fmt.Errorf("failed to create backup file for pgbackrest config: %w", err)
	}

	defer func() {
		args.Logger.Info("Restoring original pgbackrest config from backup")
		if err := pgbackrest.RestoreConfigFromBackup(args.pgBackrestConfigFile); err != nil {
			args.Logger.Error("failed to restore pgbackrest config from backup", zap.Error(err))
		}
	}()

	// Update postgresql data dir, as it could be different for other server
	args.Logger.Info("Creating backup of original pgbackrest config",
		zap.String("file", args.pgBackrestConfigFile),
		zap.String("data_directory", postgresqlConfig.DataDirectory),
	)
	if err := pgbackrest.UpdatePgbackrestConfig(args.Logger, args.pgBackrestConfigFile, map[string]string{
		"pg1-path": postgresqlConfig.DataDirectory,
	}); err != nil {
		return fmt.Errorf("failed to update postgresql data directory in pgbackrest config: %w", err)
	}

	args.Logger.Info("Stopping postgresql before resoring")
	if err := systemctl.Stop(args.Logger, "postgresql"); err != nil {
		return fmt.Errorf("failed to stop postgresql: %w", err)
	}

	args.Logger.Info("Stopping vegavisor before resoring")
	if err := systemctl.Stop(args.Logger, "vegavisor"); err != nil {
		return fmt.Errorf("failed to stop vegavisor: %w", err)
	}

	defer func() {
		args.Logger.Info("Starting postgresql service")
		if err := systemctl.Start(args.Logger, "postgresql"); err != nil {
			args.Logger.Error("failed to start postgresql service", zap.Error(err))
			return
		}

		args.Logger.Info("Starting vegavisor service")
		if err := systemctl.Start(args.Logger, "vegavisor"); err != nil {
			args.Logger.Error("failed to start vegavisor service", zap.Error(err))
			return
		}

		time.Sleep(30 * time.Second)
		if !systemctl.IsRunning(args.Logger, "postgresql") {
			args.Logger.Error("the postgresql service failed within 30 seconds after start", zap.Error(err))
		}

		if !systemctl.IsRunning(args.Logger, "vegavisor") {
			args.Logger.Error("the vegavisor service failed within 30 seconds after start", zap.Error(err))
		}
	}()

	// Ensure postgresql and visor are not running
	if systemctl.IsRunning(args.Logger, "postgresql") {
		return fmt.Errorf("postgresql is still running after servise has been stopped")
	}

	if systemctl.IsRunning(args.Logger, "vegavisor") {
		return fmt.Errorf("vegavisor is still running after servise has been stopped")
	}

	var (
		wg               sync.WaitGroup
		chainDataFailed  bool
		postgresqlFailed bool

		postgresqlBranchFinished atomic.Bool
		s3BranchFinished         atomic.Bool
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer s3BranchFinished.Store(true)

		args.Logger.Info("Removing local chain data")

		args.Logger.Info("Removing local vega chain data")
		if err := vegachain.RemoveLocalChainData(args.Logger); err != nil {
			chainDataFailed = true
			args.Logger.Error("failed to remove local chain data", zap.Error(err))
			return
		}

		args.Logger.Info("Restoring vega chain data from remote")
		snapshotDestination := fmt.Sprintf("s3://%s/%s", currentBackup.VegaChain.Location.Bucket, currentBackup.VegaChain.Location.Path)
		if err := vegachain.RestoreChainData(
			args.Logger,
			args.s3CmdBinary,
			snapshotDestination,
			currentBackup.VegaChain.Components.VisorHome,
		); err != nil {
			chainDataFailed = true
			args.Logger.Error("failed to restore chain data", zap.Error(err))
			return
		}

		if !postgresqlBranchFinished.Load() {
			args.Logger.Info("The PostgreSQL backup restore process is still in progress")
		}
	}()

	go func() {
		defer wg.Done()
		defer postgresqlBranchFinished.Store(true)

		args.Logger.Info("Starting pgbackrest stanza")
		if err := tools.Retry(3, 5*time.Second, func() error {
			return pgbackrest.Start(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary)
		}); err != nil {
			postgresqlFailed = true
			args.Logger.Error("failed to start pgbackrest stanza", zap.Error(err))
			return
		}

		args.Logger.Info("Stopping postgresql before restoring it")
		if err := systemctl.Stop(args.Logger, "postgresql"); err != nil {
			postgresqlFailed = true
			args.Logger.Error("failed to stop postgresql service", zap.Error(err))
			return
		}

		args.Logger.Info("Checking postgresql has been stopped")
		if err := tools.Retry(3, 5*time.Second, func() error {
			if systemctl.IsRunning(args.Logger, "postgresql") {
				return fmt.Errorf("postgresql service is still running")
			}
			return nil
		}); err != nil {
			postgresqlFailed = true
			args.Logger.Error("failed to check postgresql has been stopped", zap.Error(err))
			return
		}

		postmasterPidFile := filepath.Join(postgresqlConfig.DataDirectory, "postmaster.pid")

		if err := tools.Retry(20, 2*time.Second, func() error {
			if tools.FileExists(postmasterPidFile) {
				return fmt.Errorf("the %s file still exists", postmasterPidFile)
			}

			return nil
		}); err != nil {
			postgresqlFailed = true
			args.Logger.Error("failed to wait until postgresql process has been stopped and the PID file is missing", zap.Error(err))
			return
		}

		if args.pgBackrestFullRestore {
			args.Logger.Info("Full backup, removing the content of postgresql data dir", zap.String("data_directory", postgresqlConfig.DataDirectory))
			if err := tools.RemoveDirectoryContents(postgresqlConfig.DataDirectory); err != nil {
				args.Logger.Error("failed to remove content of postgresql data dir", zap.Error(err), zap.String("data_directory", postgresqlConfig.DataDirectory))
			}
		}

		args.Logger.Info("Restoring the postgresql backup")
		if err := pgbackrest.Restore(
			*args.Logger,
			args.postgresqlUser,
			args.pgBackrestBinary,
			currentBackup.Postgresql.Label,
			!args.pgBackrestFullRestore,
		); err != nil {
			postgresqlFailed = true
			args.Logger.Error("failed to restore postgresql backup", zap.Error(err))
			return
		}

		if !s3BranchFinished.Load() {
			args.Logger.Info("S3 backup restore procedure is still in progress")
		}
	}()
	wg.Wait()

	args.Logger.Info(
		"Backup finished",
		zap.Bool("chain_data_successfull", !chainDataFailed),
		zap.Bool("postgresql_successfull", !postgresqlFailed),
	)

	return nil
}
