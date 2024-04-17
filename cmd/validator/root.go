package validator

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
}

var validatorArgs Args

var Cmd = &cobra.Command{
	Use:   "validator",
	Short: "Manage validators",
	Long:  `Manage validators`,
}

func init() {
	validatorArgs.RootArgs = &rootCmd.Args
}
