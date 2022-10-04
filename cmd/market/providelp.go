package market

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	vegatypes "code.vegaprotocol.io/vega/core/types"
	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/cmd/topup"
	"github.com/vegaprotocol/devopstools/networktools"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
)

type ProvideLPArgs struct {
	*MarketArgs
}

var provideLPArgs ProvideLPArgs

// provideLPCmd represents the provideLP command
var provideLPCmd = &cobra.Command{
	Use:   "provide-lp",
	Short: "Provide Liquidity Provision",
	Long:  `Provide Liquidity Provision`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunProvideLP(provideLPArgs); err != nil {
			provideLPArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	provideLPArgs.MarketArgs = &marketArgs

	MarketCmd.AddCommand(provideLPCmd)
}

func RunProvideLP(args ProvideLPArgs) error {
	var (
		logger = args.Logger
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
	markets, err := network.DataNodeClient.GetAllMarkets()
	if err != nil {
		return err
	}

	failed := false
	for _, marketName := range []string{"AAPL", "AAVEDAI", "BTCUSD", "ETHBTC", "TSLA", "UNIDAI"} {
		market := getMarket(markets, proposerVegawallet.PublicKey, fmt.Sprintf("auto:%s", strings.ToLower(marketName)))
		if market == nil {
			logger.Info("market does not exists", zap.String("market", marketName))
		} else {
			assetId := market.TradableInstrument.Instrument.GetFuture().SettlementAsset
			if err := topup.DepositAssetToParites(
				network, tools, assetId, big.NewFloat(100000), []string{proposerVegawallet.PublicKey}, args.Logger,
			); err != nil {
				return fmt.Errorf("failed to deposit assets to provider account, %w", err)
			}
			time.Sleep(5)
			if err := provideLiquidity(marketName, network.DataNodeClient, proposerVegawallet, args.Logger, market.Id); err != nil {
				failed = true
			}
		}

	}
	if failed {
		return fmt.Errorf("at least one Provision LP failed")
	}

	return nil
}

func provideLiquidity(
	name string,
	dataNodeClient *vegaapi.DataNode,
	proposerVegawallet *wallet.VegaWallet,
	logger *zap.Logger,
	marketId string,
) error {
	// Prepare vegawallet Transaction Request
	walletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: proposerVegawallet.PublicKey,
		Command: &walletpb.SubmitTransactionRequest_LiquidityProvisionSubmission{
			LiquidityProvisionSubmission: &commandspb.LiquidityProvisionSubmission{
				Fee:              "0.01",
				MarketId:         marketId,
				CommitmentAmount: "10000000",
				Buys: []*vega.LiquidityOrder{
					{
						Reference:  vegatypes.PeggedReferenceBestBid,
						Proportion: 10,
						Offset:     "1000",
					},
					{
						Reference:  vegatypes.PeggedReferenceBestBid,
						Proportion: 10,
						Offset:     "2000",
					},
				},
				Sells: []*vega.LiquidityOrder{
					{
						Reference:  vegatypes.PeggedReferenceBestAsk,
						Proportion: 10,
						Offset:     "2000",
					},
					{
						Reference:  vegatypes.PeggedReferenceBestAsk,
						Proportion: 10,
						Offset:     "1000",
					},
				},
			},
		},
	}

	// Sign + Proof of Work vegawallet Transaction request
	lastBlockData, err := dataNodeClient.LastBlockData()
	if err != nil {
		return err
	}
	signedTx, err := proposerVegawallet.SignTxWithPoW(&walletTxReq, lastBlockData)
	if err != nil {
		logger.Error("Failed to sign a trasnaction", zap.String("market", name),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", &walletTxReq), zap.Error(err))
		return err
	}

	// wrap in vega Transaction Request
	submitReq := &vegaapipb.SubmitTransactionRequest{
		Tx:   signedTx,
		Type: vegaapipb.SubmitTransactionRequest_TYPE_SYNC,
	}

	// Submit Transaction
	logger.Info("Submit transaction", zap.String("market", name),
		zap.String("proposer", proposerVegawallet.PublicKey))
	submitResponse, err := dataNodeClient.SubmitTransaction(submitReq)
	if err != nil {
		logger.Error("Failed to submit a trasnaction", zap.String("market", name),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", submitReq), zap.Error(err))
		return err
	}
	if !submitResponse.Success {
		logger.Error("Transaction submission response is not successful", zap.String("market", name),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", submitReq), zap.Any("response", submitResponse))
		return err
	}
	logger.Info("Successful Submision of Provide LP", zap.String("market", name),
		zap.String("proposer", proposerVegawallet.PublicKey), zap.String("txHash", submitResponse.TxHash))
	return nil
}
