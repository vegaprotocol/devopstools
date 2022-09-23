package party

import (
	"log"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type PartyArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
	PartyId         string
}

var partyArgs PartyArgs

// Root Command for Party
var PartyCmd = &cobra.Command{
	Use:   "party",
	Short: "Vega Party commands",
	Long:  `Vega Party commands`,
}

func init() {
	partyArgs.RootArgs = &rootCmd.Args

	PartyCmd.PersistentFlags().StringVar(&partyArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := PartyCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	PartyCmd.PersistentFlags().StringVar(&partyArgs.PartyId, "party", "", "Vega Party Id")
	if err := PartyCmd.MarkPersistentFlagRequired("party"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
