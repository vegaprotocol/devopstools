package backup

import (
	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type BackupArgs struct {
	*rootCmd.RootArgs

	configPath string
}

var backupArgs BackupArgs

var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create backup or restore",
}

func init() {
	backupArgs.RootArgs = &rootCmd.Args

	BackupCmd.PersistentFlags().StringVar(
		&backupArgs.configPath,
		"config-path",
		"",
		"Path to the config file")
}
