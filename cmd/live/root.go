package live

import (
	"log"

	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type LiveArgs struct {
	*rootCmd.RootArgs
	VegaNetworkName string
}

var liveArgs LiveArgs

// Root Command for Live
var LiveCmd = &cobra.Command{
	Use:   "live",
	Short: "Get Live data from running Vega Network",
	Long:  `Get Live data from running Vega Network`,
}

func init() {
	liveArgs.RootArgs = &rootCmd.Args

	LiveCmd.PersistentFlags().StringVar(&liveArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := LiveCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
