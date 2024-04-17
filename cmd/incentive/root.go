package incentive

import (
	"log"

	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var incentiveArgs Args

var Cmd = &cobra.Command{
	Use:   "incentive",
	Short: "Setup network for incentive",
	Long:  ``,
}

func init() {
	incentiveArgs.RootArgs = &rootCmd.Args

	Cmd.PersistentFlags().StringVar(&incentiveArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := Cmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
