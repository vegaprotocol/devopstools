package network

import (
	"log"

	"github.com/spf13/cobra"
)

type NetworkArgs struct {
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
	NetworkCmd.PersistentFlags().StringVar(&networkArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := NetworkCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
