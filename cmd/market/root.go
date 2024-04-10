package market

import (
	"log"

	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
	NetworkFile string
}

var args Args

var Cmd = &cobra.Command{
	Use:   "market",
	Short: "Manage a market on the network",
	Long:  "Manage a market on the network",
}

func init() {
	args.RootArgs = &rootCmd.Args

	Cmd.PersistentFlags().StringVar(&args.NetworkFile, "network-file", "./network.toml", "Path the the network file")

	if err := Cmd.MarkPersistentFlagRequired("network-file"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
