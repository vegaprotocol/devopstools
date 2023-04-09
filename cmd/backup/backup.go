package backup

import (
	// "log"

	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/cmd/backup/pgbackrest"
	"github.com/vegaprotocol/devopstools/cmd/backup/vegachain"
	"go.uber.org/zap"
)

type BackupArgs struct {
	*BackupRootArgs
	localStateFile   string
	postgresqlUser   string
	pgBackrestBinary string
	pgBackrestFull   bool

	s3CmdBinary string
}

var backupArgs BackupArgs

// provideLPCmd represents the provideLP command
var performBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup vega node to remote S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if err := DoBackup(backupArgs); err != nil {
			backupArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	backupArgs.BackupRootArgs = &backupRootArgs

	performBackupCmd.PersistentFlags().StringVar(&backupArgs.localStateFile, "local-state-file", "/tmp/vega-backup-state.json", "Local state file for the vega backup")
	performBackupCmd.PersistentFlags().StringVar(&backupArgs.postgresqlUser, "postgresql-user", "postgres", "The username who runs the postgresql")
	performBackupCmd.PersistentFlags().StringVar(&backupArgs.pgBackrestBinary, "pgbackrest-bin", "pgbackrest", "The binary for pgbackrest")
	performBackupCmd.PersistentFlags().BoolVar(&backupArgs.pgBackrestFull, "full", false, "Perform the full backup")
	performBackupCmd.PersistentFlags().StringVar(&backupArgs.s3CmdBinary, "s3cmd-bin", "s3cmd", "The binary for s3cmd")

	BackupRootCmd.AddCommand(performBackupCmd)
}

func DoBackup(args BackupArgs) error {
	pgBackrestConfig, err := pgbackrest.ReadConfig("/etc/pgbackrest.conf")
	if err != nil {
		return fmt.Errorf("failed to read pgbackrest config: %w", err)
	}

	args.Logger.Info("Verifying stanza setup")
	if err := pgbackrest.CheckPgBackRestSetup(backupArgs.pgBackrestBinary, pgBackrestConfig); err != nil {
		return fmt.Errorf("failed to check pgbackrest setup: %w", err)
	}

	// args.Logger.Info("Creating session for S3 connection")
	// s3Session, err := vegachain.NewSession(vegachain.S3Credentials{
	// 	Region:       pgBackrestConfig.Global.R1S3Region,
	// 	Endpoint:     pgBackrestConfig.Global.R1S3Endpoint,
	// 	AccessKey:    pgBackrestConfig.Global.R1S3Key,
	// 	AccessSecret: pgBackrestConfig.Global.R1S3KeySecret,
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to create s3 session: %w", err)
	// }

	args.Logger.Info("Verifying s3cmd setup")
	if err := vegachain.CheckS3Setup(args.s3CmdBinary); err != nil {
		return fmt.Errorf("failed to check s3 setup: %w", err)
	}

	args.Logger.Info("Initializing s3cmd config")
	if err := vegachain.S3CmdInit(args.s3CmdBinary, vegachain.S3Credentials{
		Region:       pgBackrestConfig.Global.R1S3Region,
		Endpoint:     pgBackrestConfig.Global.R1S3Endpoint,
		AccessKey:    pgBackrestConfig.Global.R1S3Key,
		AccessSecret: pgBackrestConfig.Global.R1S3KeySecret,
	}); err != nil {
		return fmt.Errorf("failed to init s3cmd: %w", err)
	}

	currentState := LoadOrCreateNew(args.localStateFile)
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

	currentBackup, err := NewBackupEntry()
	if err != nil {
		return fmt.Errorf("failed to initialize new backup entry: %w", err)
	}
	if err := currentState.AddOrModifyEntry(currentBackup, true); err != nil {
		return fmt.Errorf("failed to add backup entry: %w", err)
	}

	defer func() {
		args.Logger.Info("Finissing current backup")
		currentBackup.Finished = time.Now()

		if err := currentState.AddOrModifyEntry(currentBackup, true); err != nil {
			args.Logger.Fatal("failed to add or modify backup entry on defer", zap.Error(err))
		}
	}()

	_, postgresqlBackupDir := filepath.Split(pgBackrestConfig.Global.R1Path)
	currentBackup.VegaChain.Location.Bucket = pgBackrestConfig.Global.R1S3Bucket
	currentBackup.VegaChain.Location.Path = fmt.Sprintf("vega_chain_snapshots/%s/%s", postgresqlBackupDir, currentBackup.ID)

	if err := vegachain.BackupChainData(
		args.Logger,
		args.s3CmdBinary,
		postgresqlBackupDir,
		currentBackup.VegaChain.Location.Bucket,
		currentBackup.VegaChain.Location.Path); err != nil {
		return fmt.Errorf("failed to backup vega chain data: %w", err)
	}

	args.Logger.Info("Ensuring pgbackrest stanza exists")
	if err := pgbackrest.CreateStanza(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary); err != nil {
		currentBackup.Status = BackupStatusFailed
		currentBackup.Postgresql.Status = BackupStatusFailed
		return fmt.Errorf("failed to create pgbackrest stanza: %w", err)
	}

	args.Logger.Info("Starting pgbackrest stanza")
	if err := pgbackrest.Start(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary); err != nil {
		currentBackup.Status = BackupStatusFailed
		currentBackup.Postgresql.Status = BackupStatusFailed
		return fmt.Errorf("failed to start pgbackrest stanza: %w", err)
	}

	args.Logger.Info("Checking pgbackrest stanza configuration")
	if err := pgbackrest.Check(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary); err != nil {
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

	var (
		wg         sync.WaitGroup
		stateMutex sync.Mutex
		failed     bool
	)

	wg.Add(2)
	go func() {
		defer wg.Done()

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

		stateMutex.Lock()
		currentBackup.Postgresql.Label = lastBackup.Label
		currentBackup.Postgresql.Status = BackupStatusSuccess
		stateMutex.Unlock()
	}()

	go func() {
		defer wg.Done()

		stateMutex.Lock()
		currentBackup.VegaChain.Started = time.Now()
		stateMutex.Unlock()

		args.Logger.Info("Starting vega chain data backup")
		time.Sleep(10 * time.Second)
		args.Logger.Info("Finished vega chain data backup")

		stateMutex.Lock()
		currentBackup.VegaChain.Finished = time.Now()
		stateMutex.Unlock()
	}()

	args.Logger.Info("Waiting for backup to finish")
	wg.Wait()

	if !failed {
		currentBackup.Status = BackupStatusSuccess
	}

	// fmt.Printf("%#v", pgBackrestConfig)

	return nil
}