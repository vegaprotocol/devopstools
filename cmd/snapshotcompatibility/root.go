package snapshotcompatibility

import (
	"github.com/spf13/cobra"

	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type SnapshotCompatibilityArgs struct {
	*rootCmd.RootArgs
}

var snapshotCompatibilityArgs SnapshotCompatibilityArgs

// Root Command for Spam
var SnapshotCompatibilityCmd = &cobra.Command{
	Use:   "snapshot-compatibility",
	Short: "Set of tools for the compatibility pipeline",
}

func init() {
	snapshotCompatibilityArgs.RootArgs = &rootCmd.Args

	// SnapshotCompatibilityCmd.AddCommand(loadSnapshotCmd)
	SnapshotCompatibilityCmd.AddCommand(downloadMainnetSnapshotCmd)
	SnapshotCompatibilityCmd.AddCommand(produceNewSnapshotCmd)
	SnapshotCompatibilityCmd.AddCommand(collectSnapshotCmd)
}
