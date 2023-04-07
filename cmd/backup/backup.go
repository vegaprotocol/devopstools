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

	BackupRootCmd.AddCommand(performBackupCmd)
}

func DoBackup(args BackupArgs) error {
	pgBackrestConfig, err := pgbackrest.ReadConfig("/etc/pgbackrest.conf")
	if err != nil {
		return fmt.Errorf("failed to read pgbackrest config: %w", err)
	}

	if err := pgbackrest.CheckPgBackRestSetup("pgbackrest", pgBackrestConfig); err != nil {
		return fmt.Errorf("failed to check pgbackrest setup: %w", err)
	}

	fmt.Printf("%#v", pgBackrestConfig)

	return nil
}
