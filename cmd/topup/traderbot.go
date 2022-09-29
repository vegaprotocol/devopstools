package topup

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"sync"

	"code.vegaprotocol.io/vega/protos/vega"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/networktools"
	"github.com/vegaprotocol/devopstools/veganetwork"
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
			err := depositERC20TokenToParties(network, tokenHexAddress, vegaPubKeys, args.Logger)
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
			err := depositFakeAssetToParties(networktools, assetId, asset, vegaPubKeys, args.Logger)
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

func depositERC20TokenToParties(
	network *veganetwork.VegaNetwork,
	tokenHexAddress string,
	vegaPubKeys []string,
	logger *zap.Logger,
) error {
	//
	// Setup
	//
	var (
		humanDepositAmount = big.NewFloat(1000) // in full tokens, i.e. without decimals zeros
		errMsg             = fmt.Sprintf("failed to deposit %s to %d parites on %s network", tokenHexAddress, len(vegaPubKeys), network.Network)
		minterWallet       = network.NetworkMainWallet
		erc20bridge        = network.SmartContracts.ERC20Bridge
		flowId             = rand.Int()
	)
	token, err := network.SmartContractsManager.GetAssetWithAddress(tokenHexAddress)
	if err != nil {
		return fmt.Errorf("failed to get token %s, %s: %w", tokenHexAddress, errMsg, err)
	}
	tokenInfo, err := token.GetInfo()
	if err != nil {
		return fmt.Errorf("failed to get token info %s, %s: %w", tokenHexAddress, errMsg, err)
	}
	logger.Info("deposit", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
		zap.String("token-address", token.Address.Hex()), zap.String("erc20bridge", erc20bridge.Address.Hex()),
		zap.String("minter", minterWallet.Address.Hex()), zap.String("amount-per-party", humanDepositAmount.String()),
		zap.Int("party-count", len(vegaPubKeys)), zap.Any("parties", vegaPubKeys))

	//
	// Mint enough tokens and Increase Allowance
	//
	var (
		humanTotalDepositAmount = new(big.Float).Mul(humanDepositAmount, big.NewFloat(float64(len(vegaPubKeys))))
		totalDepositAmount      = ethutils.TokenFromFullTokens(humanTotalDepositAmount, tokenInfo.Decimals)
		balance                 *big.Int
		allowance               *big.Int
		mintTx                  *ethTypes.Transaction
		allowanceTx             *ethTypes.Transaction
	)
	balance, err = token.BalanceOf(&bind.CallOpts{}, minterWallet.Address)
	if err != nil {
		return fmt.Errorf("failed to get %s balance of minter %s, %s: %w", tokenInfo.Name, minterWallet.Address.Hex(), errMsg, err)
	}
	if balance.Cmp(totalDepositAmount) < 0 {
		diff := new(big.Int).Sub(totalDepositAmount, balance)
		opts := minterWallet.GetTransactOpts()
		logger.Info("minting", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
			zap.String("minter", minterWallet.Address.Hex()),
			zap.String("balance", balance.String()), zap.String("mint-amount", diff.String()),
			zap.String("required", totalDepositAmount.String()))
		mintTx, err = token.Mint(opts, minterWallet.Address, diff)
		if err != nil {
			return fmt.Errorf("failed to mint, %s: %w", errMsg, err)
		}
	} else {
		logger.Info("no need to mint", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
			zap.String("minter", minterWallet.Address.Hex()),
			zap.String("balance", balance.String()), zap.String("required", totalDepositAmount.String()))
	}
	allowance, err = token.Allowance(&bind.CallOpts{}, minterWallet.Address, erc20bridge.Address)
	if err != nil {
		return fmt.Errorf("failed to get allowance, %s: %w", errMsg, err)
	}
	if allowance.Cmp(totalDepositAmount) < 0 {
		diff := new(big.Int).Sub(totalDepositAmount, balance)
		opts := minterWallet.GetTransactOpts()
		logger.Info("increasing allowance", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
			zap.String("minter", minterWallet.Address.Hex()),
			zap.String("allowance", allowance.String()), zap.String("increasing-by", diff.String()),
			zap.String("required", totalDepositAmount.String()))
		allowanceTx, err = token.IncreaseAllowance(opts, erc20bridge.Address, diff)
		if err != nil {
			return fmt.Errorf("failed to increase allowance, %s: %w", errMsg, err)
		}
	} else {
		logger.Info("no need to increasing allowance", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
			zap.String("minter", minterWallet.Address.Hex()),
			zap.String("allowance", allowance.String()), zap.String("required", totalDepositAmount.String()))
	}
	// wait
	if mintTx != nil {
		if err = ethutils.WaitForTransact(network.EthClient, mintTx); err != nil {
			logger.Error("failed to mint", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name), zap.Error(err))
			return fmt.Errorf("transaction failed to mint, %s: %w", errMsg, err)
		}
		logger.Info("successfully minted", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name))
	}
	if allowanceTx != nil {
		if err = ethutils.WaitForTransact(network.EthClient, allowanceTx); err != nil {
			logger.Error("failed to increase allowance", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name), zap.Error(err))
			return fmt.Errorf("transaction failed to increase allowance, %s: %w", errMsg, err)
		}
		logger.Info("successfully increased allowance", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name))
	}

	//
	// DEPOSIT to ERC20 Bridge
	//
	var (
		depositTxs       = make([]*ethTypes.Transaction, len(vegaPubKeys))
		depositAmount    = ethutils.TokenFromFullTokens(humanDepositAmount, tokenInfo.Decimals)
		success, failure int
	)
	for i, pubKey := range vegaPubKeys {
		var bytePubKey [32]byte
		copy(bytePubKey[:], []byte(pubKey))
		opts := minterWallet.GetTransactOpts()
		logger.Debug("depositing", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
			zap.String("vegaPubKey", pubKey), zap.String("amount", depositAmount.String()))
		depositTxs[i], err = erc20bridge.DepositAsset(opts, token.Address, depositAmount, bytePubKey)
		if err != nil {
			failure += 1
			logger.Error("failed to deposit", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
				zap.String("vegaPubKey", pubKey), zap.String("amount", depositAmount.String()),
				zap.Error(err))
		}
	}
	// wait
	for i, tx := range depositTxs {
		if tx == nil {
			continue
		}
		logger.Debug("waiting", zap.Any("tx", tx))
		if err = ethutils.WaitForTransact(network.EthClient, tx); err != nil {
			failure += 1
			logger.Error("failed to deposit", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
				zap.Any("tx", tx),
				zap.String("vegaPubKey", vegaPubKeys[i]), zap.String("amount", depositAmount.String()), zap.Error(err))
		} else {
			success += 1
			logger.Debug("successfully deposited", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
				zap.String("vegaPubKey", vegaPubKeys[i]), zap.String("amount", depositAmount.String()))
		}
	}
	logger.Info("Summary", zap.Int("flow", flowId), zap.String("token", tokenInfo.Name),
		zap.Int("success-count", success), zap.Int("fail-count", failure))
	if failure > 0 {
		return fmt.Errorf("%s", errMsg)
	}
	return nil
}

func depositFakeAssetToParties(networktools *networktools.NetworkTools, assetId string, asset *vega.AssetDetails, vegaPubKeys []string, logger *zap.Logger) error {
	var (
		humanMintAmount = big.NewFloat(1)
		mintAmount      = ethutils.TokenFromFullTokens(humanMintAmount, uint8(asset.Decimals))
		flowId          = rand.Int()
	)
	if maxFaucetMintAmount, ok := new(big.Int).SetString(asset.GetBuiltinAsset().MaxFaucetAmountMint, 0); ok {
		if maxFaucetMintAmount.Cmp(mintAmount) < 0 {
			mintAmount = maxFaucetMintAmount
		}
	}
	logger.Info("topping up fake", zap.Int("flow", flowId), zap.String("mint amount", mintAmount.String()), zap.String("asset", asset.Name), zap.Int("parties-count", len(vegaPubKeys)))

	resultsChannel := make(chan error, len(vegaPubKeys))
	var wg sync.WaitGroup

	// Trigger ERC20 TopUps
	for _, vegaPubKeys := range vegaPubKeys {
		wg.Add(1)
		go func(vegaAssetId string, vegaPubKey string) {
			defer wg.Done()
			err := networktools.MintFakeTokens(vegaPubKey, vegaAssetId, mintAmount)
			resultsChannel <- err
			if err != nil {
				logger.Error("failed to mint", zap.Int("flow", flowId), zap.String("assetId", vegaAssetId), zap.String("vegaPubKey", vegaPubKey), zap.Error(err))
			}
		}(assetId, vegaPubKeys)
	}
	wg.Wait()
	close(resultsChannel)

	var (
		success, failure int
		failed           = false
	)
	for err := range resultsChannel {
		if err != nil {
			failed = true
			failure += 1
		} else {
			success += 1
		}
	}
	logger.Info("Summary fake asset", zap.Int("flow", flowId), zap.String("sssetId", assetId),
		zap.Int("success-count", success), zap.Int("fail-count", failure))
	if failed {
		return fmt.Errorf("failed to top up all the parties (success: %d, failure: %d)", success, failure)
	}

	return nil
}
