package backup

import (
	// "log"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type RestoreArgs struct {
	*BackupRootArgs
	localStateFile string

	pgBackrestBinary     string
	pgBackrestConfigFile string

	pgBackrestDeltaRestore bool
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

	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.localStateFile, "local-state-file", "/tmp/vega-backup-state.json", "Local state file for the vega backup")
	performRestoreCmd.PersistentFlags().StringVar(&restoreArgs.pgBackrestBinary, "pgbackrest-bin", "pgbackrest", "The binary for pgbackrest")
	performRestoreCmd.PersistentFlags().BoolVar(&restoreArgs.pgBackrestDeltaRestore, "delta", false, "Perform the delta restore for postgresql")

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

	if err := os.WriteFile(args.localStateFile, []byte(state.PgBackrestConfig), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write pgbackrest config from state: %w", err)
	}

	return nil
}
