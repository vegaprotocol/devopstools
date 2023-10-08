package market

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/networktools"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"go.uber.org/zap"
)

type ProvideMinLiquidityArgs struct {
	*MarketArgs

	MarketName string
}

var provideMinLiquidityArgs ProvideMinLiquidityArgs

// provideMinLiquidityCmd represents the provideMinLiquidity command
var provideMinLiquidityCmd = &cobra.Command{
	Use:   "provide-min-liquidity",
	Short: "Provide Minimum required, missing Liquidity to the Markets",
	Long:  `Provide Minimum required, missing Liquidity to the Markets`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunProvideMinLiquidity(provideMinLiquidityArgs); err != nil {
			provideMinLiquidityArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	provideMinLiquidityArgs.MarketArgs = &marketArgs

	MarketCmd.AddCommand(provideMinLiquidityCmd)
	provideMinLiquidityCmd.PersistentFlags().StringVar(&provideMinLiquidityArgs.MarketName, "name", "", "Optional Market Name. If set then, provide LP to that market only")
}

func RunProvideMinLiquidity(args ProvideMinLiquidityArgs) error {
	var (
		logger = args.Logger
		market vegaapi.MarketInfo
	)
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()
	tools, err := networktools.NewNetworkTools(args.VegaNetworkName, args.Logger)
	if err != nil {
		return err
	}

	proposerVegawallet := network.VegaTokenWhale

	tradeableMarkets, err := network.DataNodeClient.GetTradeableMakertInfo()
	if err != nil {
		return err
	}

	if len(args.MarketName) > 0 {
		foundMarket := false
		for _, marketInfo := range tradeableMarkets {
			if marketInfo.Market.TradableInstrument.Instrument.Name == args.MarketName {
				market = marketInfo
				foundMarket = true
			}
		}
		if !foundMarket {
			return fmt.Errorf("there is no market with name %s", args.MarketName)
		}
	} else {
		return fmt.Errorf("You need to provide filters to get single market")
	}

	logger.Info("Found market", zap.String("id", market.Market.Id),
		zap.String("supplied stake", market.MarketData.SuppliedStake),
		zap.String("target stake", market.MarketData.TargetStake),
		zap.String("asset", market.SettlementAsset.Details.Name),
		zap.Uint64("asset decimal", market.SettlementAsset.Details.Decimals),
		zap.Uint64("market decimal places", market.Market.DecimalPlaces),
	)

	_ = tools
	_ = proposerVegawallet

	return nil
}

func ProvideMinLiquidity(
	marketId string,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	// errMsg := fmt.Sprintf("failed to provide Minimum Liquidity to market %s", marketId)
	// marketData, err := dataNodeClient.GetLatestMarketDataById(marketId)
	// if err != nil {
	// 	return fmt.Errorf("%s, %w", errMsg, err)
	// }

	// _ = marketData

	return nil
}
