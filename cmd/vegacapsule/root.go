package vegacapsule

import (
	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type VegacapsuleArgs struct {
	*rootCmd.RootArgs
}

var vegacapsuleArgs VegacapsuleArgs

// Root Command for OPS
var VegacapsuleCmd = &cobra.Command{
	Use:   "vegacapsule",
	Short: "Set of commands that extends the vegacapsule for a various usecases",
}

func init() {
	vegacapsuleArgs.RootArgs = &rootCmd.Args
}
