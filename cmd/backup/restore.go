package backup


import (
	// "log"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type RestoreArgs struct {
	*BackupRootArgs
}

var restoreArgs RestoreArgs

// provideLPCmd represents the provideLP command
var performRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore vega node from remote S3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if err := DoRestore(restoreArgs); err != nil {
			backupArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	restoreArgs.BackupRootArgs = &backupRootArgs

	BackupRootCmd.AddCommand(performRestoreCmd)
}

func DoRestore(args RestoreArgs) error {

	return nil
}