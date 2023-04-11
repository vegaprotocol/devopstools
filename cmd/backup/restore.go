package backup

import (
	// "log"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/cmd/backup/pgbackrest"
	"github.com/vegaprotocol/devopstools/cmd/backup/postgresql"
	"go.uber.org/zap"
)

type RestoreArgs struct {
	*BackupRootArgs
	localStateFile string

	backupID string

	postgresqlUser       string
	pgBackrestBinary     string
	pgBackrestConfigFile string

	postgresqlDBUser     string
	postgresqlDBPassword string

	pgBackrestDeltaRestore bool
}

var restoreArgs RestoreArgs

// provideLPCmd represents the provideLP command
var performRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore vega node from remote S3 bucket",
	Long: `
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
	performRestoreCmd.PersistentFlags().BoolVar(&restoreArgs.pgBackrestDeltaRestore, "delta", false, "Perform the delta restore for postgresql")
	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.pgBackrestConfigFile, "pgbackrest-config-file", "/etc/pgbackrest.conf", "Location of pgbackrest config file")

	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.postgresqlDBUser, "db-user", "vega_backup_manager", "The super user for postgresql db")
	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.postgresqlDBPassword, "db-pass", "examplePassword", "Password for the db-user")

	if err := performRestoreCmd.MarkPersistentFlagRequired("id"); err != nil {
		log.Fatalf("%v\n", err)
	}
	BackupRootCmd.AddCommand(performRestoreCmd)
}

func DoRestore(args RestoreArgs) error {
	state, err := LoadFromLocal(args.localStateFile)
	if err != nil {
		return fmt.Errorf("failed to load backups state: %w", err)
	}

	if len(state.PgBackrestConfig) < 1 {
		return fmt.Errorf("missing pgbackrest config in the state")
	}

	if err := os.WriteFile(args.pgBackrestConfigFile, []byte(state.PgBackrestConfig), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write pgbackrest config from state: %w", err)
	}

	pgBackrestConfig, err := pgbackrest.ReadConfig(args.pgBackrestConfigFile)
	if err != nil {
		return fmt.Errorf("failed to read pgbackrest config: %w", err)
	}

	args.Logger.Info("Verifying stanza setup")
	if err := pgbackrest.CheckPgBackRestSetup(backupArgs.pgBackrestBinary, pgBackrestConfig); err != nil {
		return fmt.Errorf("failed to check pgbackrest setup: %w", err)
	}

	if err := pgbackrest.CreateStanza(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary); err != nil {
		return fmt.Errorf("failed to create pgbackrest stanza: %w", err)
	}

	sysInfo, err := collectSystemInfo(args.postgresqlDBUser, args.postgresqlDBPassword)
	if err != nil {
		return fmt.Errorf("failed to collect system info: %w", err)
	}

	fmt.Printf("%#v", sysInfo)

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

	pgWalPath := filepath.Join(postgresqlDataDir, "pg_wal2")
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
			return nil, fmt.Errorf("the pg_wal file(%s) is not a dir", pgWalPath)
		}
	}

	return result, nil
}

// After restoring of the chain data
func patchSystem(sysInfo systemInfo) error {

	// TODO: Implement it
	return nil
}
