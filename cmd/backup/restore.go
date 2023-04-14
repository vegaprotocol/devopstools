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

	postgresqlDBUser     string
	postgresqlDBPassword string

	pgBackrestFullRestore bool
}

var restoreArgs RestoreArgs

var performRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore vega node from remote S3 bucket",
	Long: `
	TBD
	TBD:
	CREATE USER vega_backup_manager WITH ENCRYPTED PASSWORD 'examplePassword';
	ALTER USER vega_backup_manager  WITH SUPERUSER;
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

	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.postgresqlDBUser, "db-user", "vega_backup_manager", "The super user for postgresql db")
	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.postgresqlDBPassword, "db-pass", "examplePassword", "Password for the db-user")

	if err := performRestoreCmd.MarkPersistentFlagRequired("id"); err != nil {
		log.Fatalf("%v\n", err)
	}
	BackupRootCmd.AddCommand(performRestoreCmd)
}

func DoRestore(args RestoreArgs) error {
	args.Logger.Info("Ensuring postgresql is running")
	if !systemctl.IsRunning(args.Logger, "postgresql") {
		return fmt.Errorf("the postgresql service is not running")
	}

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

	args.Logger.Info("Creating stanza")
	if err := pgbackrest.CreateStanza(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary); err != nil {
		return fmt.Errorf("failed to create pgbackrest stanza: %w", err)
	}

	args.Logger.Info("Collecting the system facts")
	sysInfo, err := collectSystemInfo(args.postgresqlDBUser, args.postgresqlDBPassword)
	if err != nil {
		return fmt.Errorf("failed to collect system info: %w", err)
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

		postmasterPidFile := filepath.Join(sysInfo.PostgreSqlDataDir, "postmaster.pid")

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
			// Remove all custom tablespaces. We have to do it in case any of the tablespace is in custom location
			for _, tablespaceLocation := range sysInfo.CustomTablespaces {
				args.Logger.Info("Removing custom tablespace", zap.String("location", tablespaceLocation))
				if err := os.RemoveAll(tablespaceLocation); err != nil {
					postgresqlFailed = true
					args.Logger.Error("failed remove custom tablespace", zap.Error(err), zap.String("location", tablespaceLocation))
					return
				}
			}

			// Sometimes We link pg_wal to another location. We have to remove all wal files if they are in custom location
			if sysInfo.PostgreSqlPgWalDir.IsLink {
				args.Logger.Info("Removing content of linked pg_wal dir", zap.String("location", sysInfo.PostgreSqlPgWalDir.LinkTarget))
				if err := os.RemoveAll(sysInfo.PostgreSqlPgWalDir.LinkTarget); err != nil {
					postgresqlFailed = true
					args.Logger.Error("failed to remove linked pg_wall directory content", zap.Error(err), zap.String("location", sysInfo.PostgreSqlPgWalDir.LinkTarget))
					return
				}
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

	patchSystem(*sysInfo, args.pgBackrestFullRestore)

	args.Logger.Info(
		"Backup finished",
		zap.Bool("chain_data_successfull", !chainDataFailed),
		zap.Bool("postgresql_successfull", !postgresqlFailed),
	)

	return nil
}

type systemInfo struct {
	PostgreSqlDataDir  string
	PostgreSqlPgWalDir struct {
		IsLink     bool
		LinkTarget string
	}
	CustomTablespaces map[string]string
}

// We are collecting system information to detect all customization, and then after restore procedure we have to
// revert them, because they are done for some reason.
func collectSystemInfo(postgresqlDbUser, postgresqlDbPass string) (*systemInfo, error) {
	psqlClient, err := postgresql.Client(postgresqlDbUser, postgresqlDbPass, "localhost", 5432, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgresql client: %w", err)
	}

	postgresqlDataDir, err := postgresql.GetDataDirectory(psqlClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgresql data_directory: %w", err)
	}

	customTablespaces, err := postgresql.GetCustomTablespaces(psqlClient)
	if err != nil {
		return nil, fmt.Errorf("failed to list custom postgresql tablespaces: %w", err)
	}

	result := &systemInfo{
		PostgreSqlDataDir: postgresqlDataDir,
		CustomTablespaces: customTablespaces,
	}

	pgWalPath := filepath.Join(postgresqlDataDir, "pg_wal")
	walInfo, err := os.Lstat(pgWalPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to stat pg_wall dir: %w", err)
		}

		result.PostgreSqlPgWalDir.IsLink = false
		return result, nil
	}

	if walInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		result.PostgreSqlPgWalDir.IsLink = true

		link, err := os.Readlink(pgWalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read pg_wal link: %w", err)
		}
		result.PostgreSqlPgWalDir.LinkTarget = link

		// Make sure destination exists and it is a directory
		walDestination, err := os.Stat(link)
		if err != nil {
			return nil, fmt.Errorf("failed to access the pg_wal link destination: %w", err)
		}
		if !walDestination.IsDir() {
			return nil, fmt.Errorf("the pg_wal link destination(%s) is not a dir", link)
		}
	} else {
		// PG_WAL is not a link, make sure can be accessed and it is a directory
		walDestination, err := os.Stat(pgWalPath)
		if err != nil {
			return nil, fmt.Errorf("failed to stat the pg_wal to ensure it's a dir: %w", err)
		}
		if !walDestination.IsDir() {
			return nil, fmt.Errorf("the pg_wal(%s) is not a dir", pgWalPath)
		}
	}

	return result, nil
}

// After restoring of the chain data
func patchSystem(sysInfo systemInfo, fullRestore bool) error {
	// For delta restore we do not need to touch filesystem
	if !fullRestore {
		return nil
	}

	// TODO: Implement it
	return nil
}
