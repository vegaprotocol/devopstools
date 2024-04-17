package network

import (
	"log"

	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var networkArgs Args

var Cmd = &cobra.Command{
	Use:   "network",
	Short: "Vega Network commands",
	Long:  `Vega Network commands`,
}

func init() {
	networkArgs.RootArgs = &rootCmd.Args

	Cmd.PersistentFlags().StringVar(&networkArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := Cmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
