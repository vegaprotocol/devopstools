package network

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
	Use:   "network",
	Short: "Vega Network commands",
	Long:  `Vega Network commands`,
}

func init() {
	args.RootArgs = &rootCmd.Args

	Cmd.PersistentFlags().StringVar(&args.NetworkFile, "network-file", "./network.toml", "Path the the network file")

	if err := Cmd.MarkPersistentFlagRequired("network-file"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
