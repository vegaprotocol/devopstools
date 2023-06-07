package benchmark

import (
	"github.com/spf13/cobra"
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
)

type BenchmarkArgs struct {
	*rootCmd.RootArgs
}

var benchmarkArgs BenchmarkArgs

// Root Command for Benchmark
var BenchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark various things",
	Long:  `Benchmark various things`,
}

func init() {
	benchmarkArgs.RootArgs = &rootCmd.Args
}
