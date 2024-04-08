package vegacapsule

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs

	networkHomePath   string
	vegacapsuleBinary string
}

var vegacapsuleArgs Args

var Cmd = &cobra.Command{
	Use:   "vegacapsule",
	Short: "Set of commands that extends the vegacapsule for a various usecases",
}

func init() {
	vegacapsuleArgs.RootArgs = &rootCmd.Args

	Cmd.PersistentFlags().StringVar(
		&vegacapsuleArgs.networkHomePath,
		"network-home-path",
		"",
		"Custom path for the network")

	Cmd.PersistentFlags().StringVar(
		&vegacapsuleArgs.vegacapsuleBinary,
		"vegacapsule-bin",
		"vegacapsule",
		"Path to the vegacapsule binary")
}
