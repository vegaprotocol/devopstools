package smartcontracts

import (
	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type SmartContractsArgs struct {
	*rootCmd.RootArgs
}

var smartContractsArgs SmartContractsArgs

// Root Command for OPS
var SmartContractsCmd = &cobra.Command{
	Use:   "smart-contracts",
	Short: "General Smart Contracts tasks",
	Long:  `General Smart Contracts tasks`,
}

func init() {
	smartContractsArgs.RootArgs = &rootCmd.Args
}
