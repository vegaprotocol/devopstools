package snapshotcompatibility

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
}

var snapshotCompatibilityArgs Args

var Cmd = &cobra.Command{
	Use:   "snapshot-compatibility",
	Short: "Set of tools for the compatibility pipeline",
}

func init() {
	snapshotCompatibilityArgs.RootArgs = &rootCmd.Args

	Cmd.AddCommand(loadSnapshotCmd)
	Cmd.AddCommand(downloadMainnetSnapshotCmd)
	Cmd.AddCommand(downloadBinaryCmd)
	Cmd.AddCommand(produceNewSnapshotCmd)
	Cmd.AddCommand(collectSnapshotCmd)
}
