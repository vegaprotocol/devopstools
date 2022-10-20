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

type PartiesArgs struct {
	*TopUpArgs
	VegaNetworkName      string
	ERC20TokenHexAddress string
	FakeTokenAssetId     string
	PartiesVegaPubKeys   []string
	Amount               string
}

var partiesArgs PartiesArgs

// partiesCmd represents the parties command
var topUpPartiesCmd = &cobra.Command{
	Use:   "parties",
	Short: "TopUp parties for network",
	Long:  `TopUp parties for network`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunTopUpParties(partiesArgs); err != nil {
			partiesArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	partiesArgs.TopUpArgs = &topUpArgs

	TopUpCmd.AddCommand(topUpPartiesCmd)
	topUpPartiesCmd.PersistentFlags().StringVar(&partiesArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := topUpPartiesCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	topUpPartiesCmd.PersistentFlags().StringVar(&partiesArgs.ERC20TokenHexAddress, "erc20-token-address", "", "Vega Network name")
	topUpPartiesCmd.PersistentFlags().StringVar(&partiesArgs.FakeTokenAssetId, "fake-asset-id", "", "Vega Network name")
	topUpPartiesCmd.PersistentFlags().StringSliceVar(&partiesArgs.PartiesVegaPubKeys, "parties", nil, "Comma separated list of parties pub keys")
	if err := topUpPartiesCmd.MarkPersistentFlagRequired("parties"); err != nil {
		log.Fatalf("%v\n", err)
	}
	topUpPartiesCmd.PersistentFlags().StringVar(&partiesArgs.Amount, "amount", "", "Amount to top up. IMPORTANT: without decimal places. You can use float")
	if err := topUpPartiesCmd.MarkPersistentFlagRequired("amount"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunTopUpParties(args PartiesArgs) error {
	var (
		traders = &networktools.Traders{
			ByERC20TokenHexAddress: map[string][]string{},
			ByFakeAssetId:          map[string][]string{},
		}
		amount = new(big.Float)
	)
	if len(args.ERC20TokenHexAddress) > 0 {
		traders.ByERC20TokenHexAddress[args.ERC20TokenHexAddress] = args.PartiesVegaPubKeys
	} else if len(args.FakeTokenAssetId) > 0 {
		traders.ByFakeAssetId[args.FakeTokenAssetId] = args.PartiesVegaPubKeys
	} else {
		return fmt.Errorf("missing argument: either --erc20-token-address or --fake-asset-id needs to be set")
	}

	if _, ok := amount.SetString(args.Amount); !ok {
		return fmt.Errorf("failed to parse amount '%s' to float", args.Amount)
	}

	networktools, err := networktools.NewNetworkTools(args.VegaNetworkName, args.Logger)
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
			err := depositERC20TokenToParties(network, tokenHexAddress, vegaPubKeys, amount, args.Logger)
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
			err := depositFakeAssetToParties(networktools, assetId, asset, vegaPubKeys, amount, args.Logger)
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
