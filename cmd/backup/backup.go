package backup


import (
	// "log"
	"os"

	"github.com/spf13/cobra"
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

	return nil
}