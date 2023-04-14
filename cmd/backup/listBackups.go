package backup

import (
	// "log"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ListBackupsArgs struct {
	*BackupRootArgs

	outputJSON     bool
	localStateFile string
}

var listBackupsArgs ListBackupsArgs

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
	listBackupsCmd.PersistentFlags().BoolVar(&listBackupsArgs.outputJSON, "json", false, "Print output in JSON")
	listBackupsCmd.PersistentFlags().StringVar(&listBackupsArgs.localStateFile, "local-state-file", "/tmp/vega-backup-state.json", "Local state file for the vega backup")

	BackupRootCmd.AddCommand(listBackupsCmd)
}

func DoListBackups(args ListBackupsArgs) error {
	state, err := LoadFromLocal(args.localStateFile)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	stateString, err := state.AsJSON()
	if err != nil {
		return fmt.Errorf("failed to convert state to JSON: %w", err)
	}
	if args.outputJSON {
		fmt.Println(stateString)
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Started", "Finished", "PSQL Status", "Chain Status"})

	// TODO: Add sorting by finished date
	for _, backup := range state.Backups {
		table.Append([]string{
			backup.ID.String(),
			backup.Started.Format(time.RFC3339),
			backup.Finished.Format(time.RFC3339),
			string(backup.Postgresql.Status),
			string(backup.Status),
		})
	}
	table.Render() // Send output
	return nil
}
