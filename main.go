package main

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/cmd/live"
	"github.com/vegaprotocol/devopstools/cmd/network"
	"github.com/vegaprotocol/devopstools/cmd/ops"
	"github.com/vegaprotocol/devopstools/cmd/party"
	"github.com/vegaprotocol/devopstools/cmd/secrets"
	"github.com/vegaprotocol/devopstools/cmd/smartcontracts"
)

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd.RootCmd.AddCommand(ops.OpsCmd)
	rootCmd.RootCmd.AddCommand(network.NetworkCmd)
	rootCmd.RootCmd.AddCommand(live.LiveCmd)
	rootCmd.RootCmd.AddCommand(secrets.SecretsCmd)
	rootCmd.RootCmd.AddCommand(smartcontracts.SmartContractsCmd)
	rootCmd.RootCmd.AddCommand(party.PartyCmd)
}
