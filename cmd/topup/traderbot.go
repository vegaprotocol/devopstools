package topup

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"

	"code.vegaprotocol.io/vega/protos/vega"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/networktools"
	"go.uber.org/zap"
)

type TraderbotArgs struct {
	*TopUpArgs
	VegaNetworkName string
}

var traderbotArgs TraderbotArgs

// traderbotCmd represents the traderbot command
var topUpTraderbotCmd = &cobra.Command{
	Use:   "traderbot",
	Short: "TopUp traderbot for network",
	Long:  `TopUp traderbot for network`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunTopUpTraderbot(traderbotArgs); err != nil {
			traderbotArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	traderbotArgs.TopUpArgs = &topUpArgs

	TopUpCmd.AddCommand(topUpTraderbotCmd)
	topUpTraderbotCmd.PersistentFlags().StringVar(&traderbotArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := topUpTraderbotCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunTopUpTraderbot(args TraderbotArgs) error {
	networktools, err := networktools.NewNetworkTools(args.VegaNetworkName, args.Logger)
	if err != nil {
		return err
	}
	traders, err := networktools.GetTraderbotTraders()
	if err != nil {
		return err
	}
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()
	networkAssets, err := network.DataNodeClient.GetAssets()
	if err != nil {
		return err
	}

	ethBalanceBefore, err := network.EthClient.BalanceAt(context.Background(), network.NetworkMainWallet.Address, nil)
	if err != nil {
		return err
	}
	humanEthBalanceBefore := ethutils.WeiToEther(ethBalanceBefore)
	args.Logger.Info("eth balance of main wallet - before", zap.String("wallet", network.NetworkMainWallet.Address.Hex()),
		zap.String("balance", humanEthBalanceBefore.String()))
	defer func() {
		ethBalanceAfter, err := network.EthClient.BalanceAt(context.Background(), network.NetworkMainWallet.Address, nil)
		if err != nil {
			return
		}
		humanEthBalanceAfter := ethutils.WeiToEther(ethBalanceAfter)
		humanDiff := new(big.Float).Sub(humanEthBalanceBefore, humanEthBalanceAfter)
		args.Logger.Info("eth balance of main wallet - after", zap.String("wallet", network.NetworkMainWallet.Address.Hex()),
			zap.String("balance", humanEthBalanceAfter.String()), zap.String("cost", humanDiff.String()))
	}()

	resultsChannel := make(chan error, len(traders.ByERC20TokenHexAddress)+len(traders.ByFakeAssetId))
	var wg sync.WaitGroup

	// Trigger ERC20 TopUps
	for tokenHexAddress, vegaPubKeys := range traders.ByERC20TokenHexAddress {
		wg.Add(1)
		go func(tokenHexAddress string, vegaPubKeys []string) {
			defer wg.Done()
			err := depositERC20TokenToParties(network, tokenHexAddress, vegaPubKeys, big.NewFloat(1000000000), args.Logger)
			if err != nil {
				resultsChannel <- err
			}
		}(tokenHexAddress, vegaPubKeys)
	}
	// Trigger Fake Assets TopUps
	for assetId, vegaPubKeys := range traders.ByFakeAssetId {
		asset := networkAssets[assetId]
		wg.Add(1)
		go func(assetId string, asset *vega.AssetDetails, vegaPubKeys []string) {
			defer wg.Done()
			err := depositFakeAssetToParties(networktools, assetId, asset, vegaPubKeys, big.NewFloat(1000000), args.Logger)
			resultsChannel <- err
		}(assetId, asset, vegaPubKeys)
	}
	wg.Wait()
	close(resultsChannel)

	failed := false
	for err := range resultsChannel {
		//for _, err := range network.AssetMainWallet.WaitForQueue() {
		if err != nil {
			failed = true
			args.Logger.Error("transaciton failed", zap.Error(err))
		}
	}
	if failed {
		return fmt.Errorf("failed to top up all the parties")
	}
	fmt.Printf("DONE\n")
	return nil
}
