package topup

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/ethutils"
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
	traders, err := getTraders(args.VegaNetworkName)
	if err != nil {
		return err
	}
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()

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
		wg.Add(1)
		go func(assetId string, vegaPubKeys []string) {
			defer wg.Done()
			err := depositFakeAssetToParties(network, assetId, vegaPubKeys, args.Logger)
			resultsChannel <- err
		}(assetId, vegaPubKeys)
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

type traderbotResponse struct {
	Traders map[string]struct {
		PubKey     string `json:"pubKey"`
		Parameters struct {
			// MarketBase                              string `json:"marketBase"`
			// MarketQuote                             string `json:"marketQuote"`
			MarketSettlementEthereumContractAddress string `json:"marketSettlementEthereumContractAddress"`
			MarketSettlementVegaAssetID             string `json:"marketSettlementVegaAssetID"`
		} `json:"parameters"`
	} `json:"traders"`
}

type Traders struct {
	ByERC20TokenHexAddress map[string][]string
	ByFakeAssetId          map[string][]string
}

func getTraders(network string) (*Traders, error) {
	// TODO curl the traderbot endpoint - easy
	byteAssets, err := ioutil.ReadFile("traderbot.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read file with traders, %w", err)
	}

	payload := traderbotResponse{}

	if err = json.Unmarshal(byteAssets, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse file with traders, %w", err)
	}

	result := Traders{
		ByERC20TokenHexAddress: map[string][]string{},
		ByFakeAssetId:          map[string][]string{},
	}

	for _, trader := range payload.Traders {
		tokenHexAddress := trader.Parameters.MarketSettlementEthereumContractAddress
		if len(tokenHexAddress) > 0 {
			_, ok := result.ByERC20TokenHexAddress[tokenHexAddress]
			if ok {
				result.ByERC20TokenHexAddress[tokenHexAddress] = append(result.ByERC20TokenHexAddress[tokenHexAddress], trader.PubKey)
			} else {
				result.ByERC20TokenHexAddress[tokenHexAddress] = []string{trader.PubKey}
			}
		} else {
			assetId := trader.Parameters.MarketSettlementVegaAssetID
			_, ok := result.ByFakeAssetId[assetId]
			if ok {
				result.ByFakeAssetId[assetId] = append(result.ByFakeAssetId[assetId], trader.PubKey)
			} else {
				result.ByFakeAssetId[assetId] = []string{trader.PubKey}
			}
		}
	}

	return &result, nil
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

func depositFakeAssetToParties(network *veganetwork.VegaNetwork, assetId string, vegaPubKeys []string, logger *zap.Logger) error {
	logger.Debug("topping up fake", zap.String("assetId", assetId), zap.Int("parties-count", len(vegaPubKeys)))
	// TODO implement - easy
	return nil
}
