package script

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type ScriptArgs struct {
	*rootCmd.RootArgs
}

var scriptArgs ScriptArgs

// Root Command for Script
var ScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Use template to quickly get what you need",
	Long: `This section contains multiple built-up templates that you can use to quickly achive what you want.
	If you see that the command you created might be useful, the move it elsewhere, and leave the template in an original state.`,
}

func init() {
	scriptArgs.RootArgs = &rootCmd.Args
}
