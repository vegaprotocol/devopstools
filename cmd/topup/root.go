package topup

import (
	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type TopUpArgs struct {
	*rootCmd.RootArgs
}

var topUpArgs TopUpArgs

// Root Command for OPS
var TopUpCmd = &cobra.Command{
	Use:   "topup",
	Short: "Deposit ERC20 assets to vega pub keys",
	Long:  `Deposit ERC20 assets to vega pub keys`,
}

func init() {
	topUpArgs.RootArgs = &rootCmd.Args
}
