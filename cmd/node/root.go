package node

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
}

var args Args

var Cmd = &cobra.Command{
	Use:   "node",
	Short: "Manage a node on the network",
	Long:  "Manage a node on the network",
}

func init() {
	args.RootArgs = &rootCmd.Args
}
