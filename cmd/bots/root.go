package bots

import (
	"log"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type BotsArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var botsArgs BotsArgs

// Root Command for Bots
var BotsCmd = &cobra.Command{
	Use:   "bots",
	Short: "Work with bots",
	Long:  `Work with bots`,
}

func init() {
	botsArgs.RootArgs = &rootCmd.Args
	BotsCmd.PersistentFlags().StringVar(&botsArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := BotsCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
