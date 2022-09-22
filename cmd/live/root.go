package live

import (
	"log"

	"github.com/spf13/cobra"
)

type LiveArgs struct {
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
	LiveCmd.PersistentFlags().StringVar(&liveArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := LiveCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}
