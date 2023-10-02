package incentive

import (
	"log"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type IncentiveArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var incentiveArgs IncentiveArgs

// Root Command for Incentive
var IncentiveCmd = &cobra.Command{
	Use:   "incentive",
	Short: "Setup network for incentive",
	Long:  ``,
}

func init() {
	incentiveArgs.RootArgs = &rootCmd.Args

	IncentiveCmd.PersistentFlags().StringVar(&incentiveArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := IncentiveCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
