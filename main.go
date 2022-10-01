package main

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/cmd/erc20token"
	"github.com/vegaprotocol/devopstools/cmd/live"
	"github.com/vegaprotocol/devopstools/cmd/market"
	"github.com/vegaprotocol/devopstools/cmd/network"
	"github.com/vegaprotocol/devopstools/cmd/ops"
	"github.com/vegaprotocol/devopstools/cmd/party"
	"github.com/vegaprotocol/devopstools/cmd/script"
	"github.com/vegaprotocol/devopstools/cmd/secrets"
	"github.com/vegaprotocol/devopstools/cmd/smartcontracts"
	"github.com/vegaprotocol/devopstools/cmd/topup"
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
	rootCmd.RootCmd.AddCommand(erc20token.ERC20tokenCmd)
	rootCmd.RootCmd.AddCommand(script.ScriptCmd)
	rootCmd.RootCmd.AddCommand(topup.TopUpCmd)
	rootCmd.RootCmd.AddCommand(market.MarketCmd)
}
