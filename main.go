package main

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/cmd/ops"
)

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd.RootCmd.AddCommand(ops.OpsCmd)
}
