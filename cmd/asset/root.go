package asset

import (
	"log"

	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var args Args

var Cmd = &cobra.Command{
	Use:   "asset",
	Short: "Manage an asset on the vega network",
	Long:  "Manage an asset on the vega network",
}

func init() {
	args.RootArgs = &rootCmd.Args

	Cmd.PersistentFlags().StringVar(&args.VegaNetworkName, "network", "", "Vega Network name")

	if err := Cmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
