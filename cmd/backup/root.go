package backup

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs

	configPath string
}

var backupArgs Args

var Cmd = &cobra.Command{
	Use:   "backup",
	Short: "Create backup or restore",
}

func init() {
	backupArgs.RootArgs = &rootCmd.Args

	Cmd.PersistentFlags().StringVar(
		&backupArgs.configPath,
		"config-path",
		"",
		"Path to the config file")
}
