package network

import (
	"log"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type NetworkArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var networkArgs NetworkArgs

// Root Command for Network
var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Vega Network commands",
	Long:  `Vega Network commands`,
}

func init() {
	networkArgs.RootArgs = &rootCmd.Args

	NetworkCmd.PersistentFlags().StringVar(&networkArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := NetworkCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
