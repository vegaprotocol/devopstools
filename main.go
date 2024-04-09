package main

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/cmd/asset"
	"github.com/vegaprotocol/devopstools/cmd/backup"
	"github.com/vegaprotocol/devopstools/cmd/benchmark"
	"github.com/vegaprotocol/devopstools/cmd/bots"
	"github.com/vegaprotocol/devopstools/cmd/incentive"
	"github.com/vegaprotocol/devopstools/cmd/market"
	"github.com/vegaprotocol/devopstools/cmd/network"
	"github.com/vegaprotocol/devopstools/cmd/node"
	"github.com/vegaprotocol/devopstools/cmd/parties"
	"github.com/vegaprotocol/devopstools/cmd/propose"
	"github.com/vegaprotocol/devopstools/cmd/snapshotcompatibility"
	"github.com/vegaprotocol/devopstools/cmd/validator"
	"github.com/vegaprotocol/devopstools/cmd/vegacapsule"
	"github.com/vegaprotocol/devopstools/cmd/version"
)

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd.RootCmd.AddCommand(network.Cmd)
	rootCmd.RootCmd.AddCommand(node.Cmd)
	rootCmd.RootCmd.AddCommand(parties.Cmd)
	rootCmd.RootCmd.AddCommand(bots.Cmd)
	rootCmd.RootCmd.AddCommand(asset.Cmd)
	rootCmd.RootCmd.AddCommand(market.Cmd)
	rootCmd.RootCmd.AddCommand(validator.Cmd)
	rootCmd.RootCmd.AddCommand(vegacapsule.Cmd)
	rootCmd.RootCmd.AddCommand(backup.Cmd)
	rootCmd.RootCmd.AddCommand(benchmark.Cmd)
	rootCmd.RootCmd.AddCommand(snapshotcompatibility.Cmd)
	rootCmd.RootCmd.AddCommand(version.Cmd)
	rootCmd.RootCmd.AddCommand(propose.Cmd)
	rootCmd.RootCmd.AddCommand(incentive.Cmd)
}
