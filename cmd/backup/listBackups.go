package backup


import (
	// "log"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ListBackupsArgs struct {
	*BackupRootArgs
}

var listBackupsArgs ListBackupsArgs

// provideLPCmd represents the provideLP command
var listBackupsCmd = &cobra.Command{
	Use:   "list-backups",
	Short: "List backups in remote S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if err := DoListBackups(listBackupsArgs); err != nil {
			backupArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	listBackupsArgs.BackupRootArgs = &backupRootArgs

	BackupRootCmd.AddCommand(listBackupsCmd)
}

func DoListBackups(args ListBackupsArgs) error {

	return nil
}