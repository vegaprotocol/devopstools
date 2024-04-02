package topup

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type TopUpArgs struct {
	*rootCmd.RootArgs
}

var topUpArgs TopUpArgs

var TopUpCmd = &cobra.Command{
	Use:   "topup",
	Short: "Deposit ERC20 assets to vega pub keys",
	Long:  `Deposit ERC20 assets to vega pub keys`,
}

func init() {
	topUpArgs.RootArgs = &rootCmd.Args
}
