package backup

import (
	// "log"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/cmd/backup/pgbackrest"
	"go.uber.org/zap"
)

type BackupArgs struct {
	*BackupRootArgs
	localStateFile   string
	postgresqlUser   string
	pgBackrestBinary string
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

	state := LoadOrCreateNew(args.localStateFile)
	if state.Locked {
		return fmt.Errorf("backup operation is locked in the state file")
	}

	args.Logger.Info("Ensuring pgbackrest stanza exists")
	if err := pgbackrest.CreateStanza(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary); err != nil {
		return fmt.Errorf("failed to create pgbackrest stanza: %w", err)
	}

	args.Logger.Info("Checking pgbackrest stanza configuration")
	if err := pgbackrest.Check(*args.Logger, args.postgresqlUser, backupArgs.pgBackrestBinary); err != nil {
		return fmt.Errorf("failed to check pgbackrest stanza: %w", err)
	}

	// fmt.Printf("%#v", pgBackrestConfig)

	return nil
}
