package main

import (
	rootCmd "github.com/vegaprotocol/devopstools/cmd"
	"github.com/vegaprotocol/devopstools/cmd/asset"
	"github.com/vegaprotocol/devopstools/cmd/backup"
	"github.com/vegaprotocol/devopstools/cmd/batch"
	"github.com/vegaprotocol/devopstools/cmd/benchmark"
	"github.com/vegaprotocol/devopstools/cmd/bot"
	"github.com/vegaprotocol/devopstools/cmd/bots"
	"github.com/vegaprotocol/devopstools/cmd/helper"
	"github.com/vegaprotocol/devopstools/cmd/incentive"
	"github.com/vegaprotocol/devopstools/cmd/live"
	"github.com/vegaprotocol/devopstools/cmd/market"
	"github.com/vegaprotocol/devopstools/cmd/network"
	"github.com/vegaprotocol/devopstools/cmd/ops"
	"github.com/vegaprotocol/devopstools/cmd/parties"
	"github.com/vegaprotocol/devopstools/cmd/party"
	"github.com/vegaprotocol/devopstools/cmd/propose"
	"github.com/vegaprotocol/devopstools/cmd/script"
	"github.com/vegaprotocol/devopstools/cmd/secrets"
	"github.com/vegaprotocol/devopstools/cmd/snapshotcompatibility"
	"github.com/vegaprotocol/devopstools/cmd/spam"
	"github.com/vegaprotocol/devopstools/cmd/validator"
	"github.com/vegaprotocol/devopstools/cmd/vegacapsule"
	"github.com/vegaprotocol/devopstools/cmd/version"
)

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd.RootCmd.AddCommand(ops.OpsCmd)
	rootCmd.RootCmd.AddCommand(network.NetworkCmd)
	rootCmd.RootCmd.AddCommand(live.LiveCmd)
	rootCmd.RootCmd.AddCommand(secrets.SecretsCmd)
	rootCmd.RootCmd.AddCommand(parties.Cmd)
	rootCmd.RootCmd.AddCommand(party.PartyCmd)
	rootCmd.RootCmd.AddCommand(script.ScriptCmd)
	rootCmd.RootCmd.AddCommand(bots.Cmd)
	rootCmd.RootCmd.AddCommand(market.MarketCmd)
	rootCmd.RootCmd.AddCommand(validator.ValidatorCmd)
	rootCmd.RootCmd.AddCommand(spam.SpamCmd)
	rootCmd.RootCmd.AddCommand(helper.HelperCmd)
	rootCmd.RootCmd.AddCommand(vegacapsule.VegacapsuleCmd)
	rootCmd.RootCmd.AddCommand(backup.BackupCmd)
	rootCmd.RootCmd.AddCommand(benchmark.BenchmarkCmd)
	rootCmd.RootCmd.AddCommand(snapshotcompatibility.SnapshotCompatibilityCmd)
	rootCmd.RootCmd.AddCommand(version.VersionCmd)
	rootCmd.RootCmd.AddCommand(propose.ProposeCmd)
	rootCmd.RootCmd.AddCommand(bot.BotCmd)
	rootCmd.RootCmd.AddCommand(incentive.IncentiveCmd)
	rootCmd.RootCmd.AddCommand(batch.BatchCmd)
	rootCmd.RootCmd.AddCommand(asset.Cmd)
}
