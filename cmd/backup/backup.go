package backup

import (
	// "log"

	"fmt"
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

type BackupArgs struct {
	*BackupRootArgs

	pgBackrestFull    bool
	postgresqlUser    string
	postgresqlService string
	encryptionKey     string

	localStateFile       string
	pgBackrestConfigFile string
	postgresqlConfigFile string

	s3CmdBinary      string
	pgBackrestBinary string
}

var backupArgs BackupArgs

var performBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup vega node to remote S3 bucket",
	Long: `The backup command combines two methods of backup into one command:
    - we archive the vega chain data into S3 with the 's3cmd' program
    - to archive the postgresql database, we use the 'pgbackrest' software
We save both backups in the S3 bucket.

Both parts of the backup need to be synchronized, and We must collect the data 
simultaneously after We stop the Vega node.


Requirements to make a backup:

    - The 's3cmd' is installed in the system
    - The 'pgbackrest' is installed.
    - The pgbackrest stanza is configured, and the 'pgbackrest check' command does not fail.
    - The postgresql archive_command is configured
    - The PostgreSQL server is running
    - The vega_home is initialized under the /home/vega/vega_home
    - The tendermint home is initialized under the /home/vega/tendermint_home
    - The vegavisor home is initialized under /home/vega/vegavisor_home
    - The server uses the systemctl to manage postgresql and vegavisor processes


You must start the program with the user who can:

    - read pgbackrest config
    - read postgresql config
    - restart the vegavisor service
    - restart the postgresql service
    - read and write the backup state file
    - read the home directory for vega, tendermint, and visor.


Full backup

You may provide the '--full' flag is present, the --full flag is passed to the 'pgbackrest backup' subcommand
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := DoBackup(backupArgs); err != nil {
			backupArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	backupArgs.BackupRootArgs = &backupRootArgs

	performBackupCmd.PersistentFlags().BoolVar(&backupArgs.pgBackrestFull, "full", false, "Flag enforcing the full backup on the server")
	performBackupCmd.PersistentFlags().StringVar(&backupArgs.postgresqlUser, "postgresql-user", "postgres", "Linux username who runs the postgresql service")
	performBackupCmd.PersistentFlags().StringVar(&backupArgs.postgresqlService, "postgresql-service", "postgresql", "Systemctl service name for the postgresql server")
	performBackupCmd.PersistentFlags().StringVar(&backupArgs.encryptionKey, "`passphrase`", "0123456789abcdef", "Encryption key for the sensitive data in the state JSON file. It must be 16 characters in length or multiplied of 16")

	performBackupCmd.PersistentFlags().StringVar(&backupArgs.localStateFile, "local-state-file", "/tmp/vega-backup-state.json", "Local path for the backup state json file")
	performBackupCmd.PersistentFlags().StringVar(&backupArgs.pgBackrestConfigFile, "`pgbackrest-config-file`", "/etc/pgbackrest.conf", "Path to the pgbackrest config file")
	performBackupCmd.PersistentFlags().StringVar(&backupArgs.postgresqlConfigFile, "postgresql-config-file", "/etc/postgresql/14/main/postgresql.conf", "Path to the postgresql.conf file")

	performBackupCmd.PersistentFlags().StringVar(&backupArgs.s3CmdBinary, "s3cmd-bin", "s3cmd", "Path to the s3cmd binary")
	performBackupCmd.PersistentFlags().StringVar(&backupArgs.pgBackrestBinary, "pgbackrest-bin", "pgbackrest", "Path to the pgbackrest binary")

	BackupRootCmd.AddCommand(performBackupCmd)
}

func DoBackup(args BackupArgs) error {
	args.Logger.Info("Reading the pgbackrest config file", zap.String("file", args.pgBackrestConfigFile))
	pgBackrestConfig, err := pgbackrest.ReadConfig(args.pgBackrestConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read pgbackrest config: %w", err)
	}

	args.Logger.Info("Reading the postgresql config file", zap.String("file", args.postgresqlConfigFile))
	postgresqlConfig, err := postgresql.ReadConfig(args.postgresqlConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read postgresql config")
	}

	args.Logger.Info("Verifying stanza setup")
	if err := pgbackrest.CheckPgBackRestSetup(backupArgs.pgBackrestBinary, pgBackrestConfig); err != nil {
		return fmt.Errorf("failed to check pgbackrest setup: %w", err)
	}

	args.Logger.Info("Verifying s3cmd setup")
	if err := vegachain.CheckS3Setup(args.s3CmdBinary); err != nil {
		return fmt.Errorf("failed to check s3 setup: %w", err)
	}

	args.Logger.Info("Checking postgresql service")
	if !systemctl.IsRunning(args.Logger, args.postgresqlService) {
		return fmt.Errorf("postgresql service is not running")
	}

	args.Logger.Info("Initializing s3cmd config")
	if err := vegachain.S3CmdInit(args.s3CmdBinary, vegachain.S3Credentials{
		Region:       pgBackrestConfig.Global.R1S3Region,
		Endpoint:     pgBackrestConfig.Global.R1S3Endpoint,
		AccessKey:    pgBackrestConfig.Global.R1S3Key,
		AccessSecret: pgBackrestConfig.Global.R1S3KeySecret,
	}); err != nil {
		return fmt.Errorf("failed to init s3cmd config: %w", err)
	}

	currentState, err := LoadOrCreateNew(args.encryptionKey, args.localStateFile)
	if err != nil {
		return fmt.Errorf("failed to load or create new state: %w", err)
	}
	if currentState.Locked {
		return fmt.Errorf("backup operation is locked in the state file")
	}
	currentState.Locked = true
	if err := currentState.WriteLocal(args.localStateFile); err != nil {
		return fmt.Errorf("failed to write backup state to local file: %w", err)
	}

	defer func() {
		args.Logger.Info("Unlocking and writing final state file")
		currentState.Locked = false
		if err := currentState.WriteLocal(args.localStateFile); err != nil {
			args.Logger.Fatal("failed to write backup state to local file on defer", zap.Error(err))
		}
	}()

	// Sometimes postgresql in the restore command restores the default values for the archive mode and restore command
	// We do not want to do this because then we cannot create new backups
	postgresqlAutoConfigFilePath := filepath.Join(postgresqlConfig.DataDirectory, "postgresql.auto.conf")
	args.Logger.Info("Ignoring the archive_mode and restore_command parameters from the postgresql.auto.conf", zap.String("file", postgresqlAutoConfigFilePath))
	changedLines, err := postgresql.IgnoreConfigParams(postgresqlAutoConfigFilePath, []string{"archive_mode", "restore_command"}, true)
	if err != nil {
		return fmt.Errorf("failed to ignore config params in the %s file: %w", postgresqlAutoConfigFilePath, err)
	}

	if changedLines > 0 {
		args.Logger.Info("Restarting the systemctl postgresql service", zap.String("service", args.postgresqlService))
		if err := systemctl.Restart(args.Logger, args.postgresqlService); err != nil {
			return fmt.Errorf("failed to restart the postgresql servivce: %w", err)
		}

		args.Logger.Info("Checking postgresql service after restaqrt")
		if !systemctl.IsRunning(args.Logger, args.postgresqlService) {
			return fmt.Errorf("postgresql service is not running after restart")
		}
	}

	rawPgBackrestConfig, err := pgbackrest.ReadRawConfig(args.pgBackrestConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read raw pgbackrest config file: %w", err)
	}
	currentState.PgBackrestConfig = rawPgBackrestConfig

	currentBackup, err := NewBackupEntry()
	if err != nil {
		return fmt.Errorf("failed to initialize new backup entry: %w", err)
	}

	if err := currentState.AddOrModifyEntry(currentBackup, true); err != nil {
		return fmt.Errorf("failed to add backup entry: %w", err)
	}

	defer func() {
		args.Logger.Info("Finishing current backup")
		currentBackup.Finished = time.Now()

		if err := currentState.AddOrModifyEntry(currentBackup, true); err != nil {
			args.Logger.Fatal("failed to add or modify backup entry on defer", zap.Error(err))
		}
	}()

	args.Logger.Info("Ensuring pgbackrest stanza exists")
	if err := tools.Retry(3, 5*time.Second, func() error {
		return pgbackrest.CreateStanza(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary)
	}); err != nil {
		currentBackup.Status = BackupStatusFailed
		currentBackup.Postgresql.Status = BackupStatusFailed
		return fmt.Errorf("failed to create pgbackrest stanza: %w", err)
	}

	args.Logger.Info("Starting pgbackrest stanza")
	if err := tools.Retry(3, 5*time.Second, func() error {
		return pgbackrest.Start(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary)
	}); err != nil {
		currentBackup.Status = BackupStatusFailed
		currentBackup.Postgresql.Status = BackupStatusFailed
		return fmt.Errorf("failed to start pgbackrest stanza: %w", err)
	}

	args.Logger.Info("Checking pgbackrest stanza configuration")
	if err := tools.Retry(3, 5*time.Second, func() error {
		return pgbackrest.Check(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary)
	}); err != nil {
		currentBackup.Status = BackupStatusFailed
		currentBackup.Postgresql.Status = BackupStatusFailed
		return fmt.Errorf("failed to check pgbackrest stanza: %w", err)
	}

	backupType := pgbackrest.BackupIncremental
	if args.pgBackrestFull {
		backupType = pgbackrest.BackupFull
	}
	currentBackup.Postgresql.Type = backupType
	currentState.AddOrModifyEntry(currentBackup, true)

	args.Logger.Info("Stopping vegavisor service")
	if err := systemctl.Stop(args.Logger, "vegavisor"); err != nil {
		return fmt.Errorf("failed to stop vegavisor: %w", err)
	}
	defer func() {
		args.Logger.Info("Starting vegavisor service")
		if err := systemctl.Start(args.Logger, "vegavisor"); err != nil {
			args.Logger.Error("failed to start vegavisor service", zap.Error(err))
			return
		}

		time.Sleep(30 * time.Second)
		if !systemctl.IsRunning(args.Logger, "vegavisor") {
			args.Logger.Error("the vegavisor service failed within 30 seconds after start", zap.Error(err))
		}
	}()

	var (
		wg         sync.WaitGroup
		stateMutex sync.Mutex
		failed     bool

		postgresqlBranchFinished atomic.Bool
		s3BranchFinished         atomic.Bool
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		defer postgresqlBranchFinished.Store(true)

		stateMutex.Lock()
		currentBackup.Postgresql.Started = time.Now()
		stateMutex.Unlock()
		args.Logger.Info("Starting pgbackrest backup")
		if err := pgbackrest.Backup(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary, backupType); err != nil {
			stateMutex.Lock()
			currentBackup.Status = BackupStatusFailed
			currentBackup.Postgresql.Status = BackupStatusFailed
			currentBackup.Postgresql.Finished = time.Now()
			failed = true
			stateMutex.Unlock()

			args.Logger.Error("failed to backup data", zap.Error(err))
			return
		}
		args.Logger.Info("Pgbackrest backup finished")

		stateMutex.Lock()
		currentBackup.Postgresql.Finished = time.Now()
		stateMutex.Unlock()

		// We have to stop stanza to avoid issues with automatic backup of postgresql data
		args.Logger.Info("Stopping pgbackrest stanza")
		if err := pgbackrest.Stop(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary); err != nil {
			stateMutex.Lock()
			currentBackup.Status = BackupStatusFailed
			currentBackup.Postgresql.Status = BackupStatusFailed
			failed = true
			stateMutex.Unlock()
			args.Logger.Error("failed to stop pgbackrest stanza", zap.Error(err))
			return
		}

		pgBackrestBackups, err := pgbackrest.Info(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary)
		if err != nil {
			stateMutex.Lock()
			currentBackup.Status = BackupStatusFailed
			currentBackup.Postgresql.Status = BackupStatusFailed
			failed = true
			stateMutex.Unlock()
			args.Logger.Error("failed to list pgbackrest backups", zap.Error(err))
			return
		}

		lastBackup := pgbackrest.LastPgBackRestBackupInfo(pgBackrestBackups, true)
		if lastBackup == nil {
			stateMutex.Lock()
			currentBackup.Status = BackupStatusFailed
			currentBackup.Postgresql.Status = BackupStatusFailed
			failed = true
			stateMutex.Unlock()
			args.Logger.Error("failed to find last pgbackrest backup", zap.Error(err))
			return
		}
		args.Logger.Info("Found last postgresql backup label", zap.String("label", lastBackup.Label))
		if !s3BranchFinished.Load() {
			args.Logger.Info("S3 backup is still in progress")
		}

		stateMutex.Lock()
		currentBackup.Postgresql.Label = lastBackup.Label
		currentBackup.Postgresql.Status = BackupStatusSuccess
		stateMutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		defer s3BranchFinished.Store(true)

		stateMutex.Lock()
		currentBackup.VegaChain.Started = time.Now()
		stateMutex.Unlock()

		args.Logger.Info("Starting vega chain data backup")

		_, postgresqlBackupDir := filepath.Split(pgBackrestConfig.Global.R1Path)
		currentBackup.VegaChain.Location.Bucket = pgBackrestConfig.Global.R1S3Bucket
		currentBackup.VegaChain.Location.Path = fmt.Sprintf("vega_chain_snapshots/%s/%s", postgresqlBackupDir, currentBackup.ID)

		chainBackupInfo, err := vegachain.BackupChainData(
			args.Logger,
			args.s3CmdBinary,
			postgresqlBackupDir,
			currentBackup.VegaChain.Location.Bucket,
			currentBackup.VegaChain.Location.Path,
		)
		if err != nil {
			args.Logger.Info("failed to backup vega chain data", zap.Error(err))
			stateMutex.Lock()
			currentBackup.Status = BackupStatusFailed
			currentBackup.VegaChain.Status = BackupStatusFailed
			stateMutex.Unlock()

			return
		}
		args.Logger.Info("Finished vega chain data backup")
		if !postgresqlBranchFinished.Load() {
			args.Logger.Info("The PostgreSQL backup is still in progress")
		}

		stateMutex.Lock()
		currentBackup.VegaChain.Finished = time.Now()
		currentBackup.VegaChain.Status = BackupStatusSuccess
		currentBackup.VegaChain.Components.VegaHome = chainBackupInfo.WithVegaHome
		currentBackup.VegaChain.Components.TendermintHome = chainBackupInfo.WithTendermintHome
		currentBackup.VegaChain.Components.VisorHome = chainBackupInfo.WithVisorHome
		stateMutex.Unlock()
	}()

	args.Logger.Info("Waiting for backup to finish")
	wg.Wait()

	if !failed {
		currentBackup.Status = BackupStatusSuccess
	} else {
		return fmt.Errorf("one of the backup failed, postgresql-status: %s, chain-status: %s", string(currentBackup.Postgresql.Status), string(currentBackup.VegaChain.Status))
	}

	return nil
}
