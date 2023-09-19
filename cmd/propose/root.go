package propose

import (
	"log"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type ProposeArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var proposeArgs ProposeArgs

// Root Command for Propose
var ProposeCmd = &cobra.Command{
	Use:   "propose",
	Short: "Submit and vote on Vega Network Proposals",
	Long:  ``,
}

func init() {
	proposeArgs.RootArgs = &rootCmd.Args

	ProposeCmd.PersistentFlags().StringVar(&proposeArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := ProposeCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
