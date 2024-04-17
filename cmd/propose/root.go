package propose

import (
	"log"

	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var proposeArgs Args

var Cmd = &cobra.Command{
	Use:   "propose",
	Short: "Submit and vote on Vega Network Proposals",
	Long:  "Submit and vote on Vega Network Proposals",
}

func init() {
	proposeArgs.RootArgs = &rootCmd.Args

	Cmd.PersistentFlags().StringVar(&proposeArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := Cmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
