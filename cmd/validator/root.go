package validator

import (
	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type ValidatorArgs struct {
	*rootCmd.RootArgs
}

var validatorArgs ValidatorArgs

// Root Command for OPS
var ValidatorCmd = &cobra.Command{
	Use:   "validator",
	Short: "Manage validators",
	Long:  `Manage validators`,
}

func init() {
	validatorArgs.RootArgs = &rootCmd.Args
}
