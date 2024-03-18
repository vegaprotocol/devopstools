package bots

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
}

var args Args

var Cmd = &cobra.Command{
	Use:   "bots",
	Short: "Manage bots on the network",
	Long:  "Manage bots on the network",
}

func init() {
	args.RootArgs = &rootCmd.Args
}
