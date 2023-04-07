package backup


import (
	// "log"

	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type BackupRootArgs struct {
	*rootCmd.RootArgs
}

var backupRootArgs BackupRootArgs

// Root Command for Market
var BackupRootCmd = &cobra.Command{
	Use:   "backup",
	Short: "Perform vega node backup",
	Long: ``,
}

func init() {
	backupRootArgs.RootArgs = &rootCmd.Args

	// MarketCmd.PersistentFlags().StringVar(&marketArgs.VegaNetworkName, "network", "", "Vega Network name")
	// if err := MarketCmd.MarkPersistentFlagRequired("network"); err != nil {
	// 	log.Fatalf("%v\n", err)
	// }
}
