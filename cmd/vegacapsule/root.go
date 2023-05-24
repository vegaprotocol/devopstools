package vegacapsule

import (
	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type VegacapsuleArgs struct {
	*rootCmd.RootArgs

	networkHomePath   string
	vegacapsuleBinary string
}

var vegacapsuleArgs VegacapsuleArgs

// Root Command for OPS
var VegacapsuleCmd = &cobra.Command{
	Use:   "vegacapsule",
	Short: "Set of commands that extends the vegacapsule for a various usecases",
}

func init() {
	vegacapsuleArgs.RootArgs = &rootCmd.Args

	VegacapsuleCmd.PersistentFlags().StringVar(
		&vegacapsuleArgs.networkHomePath,
		"network-home-path",
		"",
		"Custom path for the network")

	VegacapsuleCmd.PersistentFlags().StringVar(
		&vegacapsuleArgs.vegacapsuleBinary,
		"vegacapsule-bin",
		"vegacapsule",
		"Path to the vegacapsule binary")
}
