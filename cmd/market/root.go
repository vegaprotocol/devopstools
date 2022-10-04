package market

import (
	"log"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type MarketArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var marketArgs MarketArgs

// Root Command for Market
var MarketCmd = &cobra.Command{
	Use:   "market",
	Short: "Use template to quickly get what you need",
	Long: `This section contains multiple built-up templates that you can use to quickly achive what you want.
	If you see that the command you created might be useful, the move it elsewhere, and leave the template in an original state.`,
}

func init() {
	marketArgs.RootArgs = &rootCmd.Args

	MarketCmd.PersistentFlags().StringVar(&marketArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := MarketCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
