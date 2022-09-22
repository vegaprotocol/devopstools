package main

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/cmd/live"
	"github.com/vegaprotocol/devopstools/cmd/network"
	"github.com/vegaprotocol/devopstools/cmd/ops"
	"github.com/vegaprotocol/devopstools/cmd/secrets"
)

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd.RootCmd.AddCommand(ops.OpsCmd)
	rootCmd.RootCmd.AddCommand(network.NetworkCmd)
	rootCmd.RootCmd.AddCommand(live.LiveCmd)
	rootCmd.RootCmd.AddCommand(secrets.SecretsCmd)
}
