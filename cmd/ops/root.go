package ops

import (
	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type OpsArgs struct {
	*rootCmd.RootArgs
}

var opsArgs OpsArgs

// Root Command for OPS
var OpsCmd = &cobra.Command{
	Use:   "ops",
	Short: "General ops tasks",
	Long:  `Range of OPS tasks`,
}

func init() {
	opsArgs.RootArgs = &rootCmd.Args
}
