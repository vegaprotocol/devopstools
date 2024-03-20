package validator

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
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
