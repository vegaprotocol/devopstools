package benchmark

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"

	"github.com/spf13/cobra"
)

type Args struct {
	*rootCmd.RootArgs
}

var benchmarkArgs Args

var Cmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark various things",
	Long:  `Benchmark various things`,
}

func init() {
	benchmarkArgs.RootArgs = &rootCmd.Args
}
