package topup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type TraderbotArgs struct {
	*TopUpArgs
	VegaNetworkName string
}

var traderbotArgs TraderbotArgs

// traderbotCmd represents the traderbot command
var traderbotCmd = &cobra.Command{
	Use:   "traderbot",
	Short: "TopUp traderbot for network",
	Long:  `TopUp traderbot for network`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunTraderbot(traderbotArgs); err != nil {
			traderbotArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	traderbotArgs.TopUpArgs = &topUpArgs

	TopUpCmd.AddCommand(traderbotCmd)
	traderbotCmd.PersistentFlags().StringVar(&traderbotArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := traderbotCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunTraderbot(args TraderbotArgs) error {
	traders, err := getTraders(args.VegaNetworkName)
	if err != nil {
		return err
	}

	resultsChannel := make(chan error, len(traders.ByERC20TokenHexAddress)+len(traders.ByFakeAssetId))
	var wg sync.WaitGroup
	for tokenHexAddress, vegaPubKeys := range traders.ByERC20TokenHexAddress {
		wg.Add(1)
		go func(tokenHexAddress string, vegaPubKeys []string) {
			defer wg.Done()
			err := depositERC20TokenToParties(tokenHexAddress, vegaPubKeys, args.Logger)
			resultsChannel <- err
		}(tokenHexAddress, vegaPubKeys)
	}
	for assetId, vegaPubKeys := range traders.ByFakeAssetId {
		wg.Add(1)
		go func(assetId string, vegaPubKeys []string) {
			defer wg.Done()
			err := depositFakeAssetToParties(assetId, vegaPubKeys, args.Logger)
			resultsChannel <- err
		}(assetId, vegaPubKeys)
	}
	wg.Wait()
	close(resultsChannel)

	failed := false
	for err := range resultsChannel {
		if err != nil {
			failed = true
			args.Logger.Error("Error", zap.Error(err))
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

func depositERC20TokenToParties(tokenHexAddress string, vegaPubKeys []string, logger *zap.Logger) error {
	logger.Debug("topping up", zap.String("token", tokenHexAddress), zap.Any("parties", vegaPubKeys))
	// TODO implement - not that easy
	return nil
}

func depositFakeAssetToParties(assetId string, vegaPubKeys []string, logger *zap.Logger) error {
	logger.Debug("topping up fake", zap.String("assetId", assetId), zap.Any("parties", vegaPubKeys))
	// TODO implement - easy
	return nil
}
