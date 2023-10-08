package market

import (
	"fmt"
	"math/big"
	"os"

	vegatypes "code.vegaprotocol.io/vega/core/types"
	"code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/governance"
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

	//
	// Find Market
	//

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
		if market.MarketData == nil {
			return fmt.Errorf("missing Market Data for market %v", market)
		}
	} else {
		return fmt.Errorf("You need to provide filters to get single market")
	}

	suppliedStake, ok := new(big.Int).SetString(market.MarketData.SuppliedStake, 10)
	if !ok {
		return fmt.Errorf("failed to parse supplied stake %s for market %v", market.MarketData.SuppliedStake, market)
	}
	targetStake, ok := new(big.Int).SetString(market.MarketData.TargetStake, 10)
	if !ok {
		return fmt.Errorf("failed to parse target stake %s for market %v", market.MarketData.TargetStake, market)
	}
	logger.Info("Found market", zap.String("id", market.Market.Id),
		zap.String("supplied stake", suppliedStake.String()),
		zap.String("target stake", targetStake.String()),
		zap.String("asset", market.SettlementAsset.Details.Name),
		zap.Uint64("asset decimal", market.SettlementAsset.Details.Decimals),
		zap.Uint64("market decimal places", market.Market.DecimalPlaces),
	)

	//
	// Check if current liquidity is enough
	//

	if targetStake.Cmp(suppliedStake) <= 0 {
		logger.Info("Supplied Stake is above Target Stake, no need to provide more liquidity",
			zap.String("supplied stake", suppliedStake.String()),
			zap.String("target stake", targetStake.String()))
		return nil
	}

	//
	// Calculate amount to stake
	//

	missingAmount := new(big.Int).Sub(targetStake, suppliedStake)
	extraAmount := new(big.Int).Div(targetStake, big.NewInt(30))
	amountToStake := new(big.Int).Add(missingAmount, extraAmount)

	logger.Info("Supplied Stake is below Target Stake, need to provide more liquidity",
		zap.String("supplied stake", suppliedStake.String()),
		zap.String("target stake", targetStake.String()),
		zap.String("amount to increase commitment by (with extra)", amountToStake.String()),
	)

	//
	// Mint
	//
	lpWallet := network.VegaTokenWhale

	lpWalletBalances, err := network.DataNodeClient.GetPartyGeneralBalances(lpWallet.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to get balance for lpWallet %s, %w", lpWallet.PublicKey, err)
	}

	walletBalance, ok := lpWalletBalances[market.SettlementAsset.Id]
	if !ok {
		walletBalance = big.NewInt(0)
	}
	logger.Info("Liquidity Provider Wallet", zap.String("asset", market.SettlementAsset.Details.Name),
		zap.String("balance", walletBalance.String()), zap.Any("all balances", lpWalletBalances))

	if walletBalance.Cmp(amountToStake) < 0 {
		// TODO: need to mint etc
		return fmt.Errorf("Note implemented")
	}

	//
	// Get current commitment by lpWallet
	//
	currentCommitment := big.NewInt(0)

	currentLiquidityProvisions, err := network.DataNodeClient.LPForMarketByParty(market.Market.Id, lpWallet.PublicKey)
	if err != nil {
		return err
	}
	for _, currentLP := range currentLiquidityProvisions {
		commitmentAmount, ok := new(big.Int).SetString(currentLP.CommitmentAmount, 10)
		if !ok {
			return fmt.Errorf("failed to parse current commitment amount %s", currentLP.CommitmentAmount)
		}
		currentCommitment = currentCommitment.Add(currentCommitment, commitmentAmount)
	}

	//
	// Submit or Amend LP commitment
	//
	var (
		txReq *walletpb.SubmitTransactionRequest
	)
	if len(currentLiquidityProvisions) > 0 {
		amendAmount := new(big.Int).Add(currentCommitment, amountToStake)
		txReq = &walletpb.SubmitTransactionRequest{
			PubKey:  lpWallet.PublicKey,
			Command: NewLiquidityProvisionAmendment(market.Market.Id, amendAmount),
		}
		logger.Info("Sending LP Amendment", zap.String("amount", amendAmount.String()))
	} else {
		txReq = &walletpb.SubmitTransactionRequest{
			PubKey:  lpWallet.PublicKey,
			Command: NewLiquidityProvisionSubmission(market.Market.Id, amountToStake),
		}
		logger.Info("Sending LP Submission", zap.String("amount", amountToStake.String()))
	}
	if err = governance.SubmitTx("Provide LP by bots", network.DataNodeClient, lpWallet, logger, txReq); err != nil {
		return err
	}

	return nil
}

func NewLiquidityProvisionSubmission(
	marketId string,
	amount *big.Int,
) *walletpb.SubmitTransactionRequest_LiquidityProvisionSubmission {
	return &walletpb.SubmitTransactionRequest_LiquidityProvisionSubmission{
		LiquidityProvisionSubmission: &commandspb.LiquidityProvisionSubmission{
			Fee:              "0.01",
			MarketId:         marketId,
			CommitmentAmount: amount.String(),
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
	}
}

func NewLiquidityProvisionAmendment(
	marketId string,
	amount *big.Int,
) *walletpb.SubmitTransactionRequest_LiquidityProvisionAmendment {
	return &walletpb.SubmitTransactionRequest_LiquidityProvisionAmendment{
		LiquidityProvisionAmendment: &commandspb.LiquidityProvisionAmendment{
			MarketId:         marketId,
			CommitmentAmount: amount.String() + "0",
			Fee:              "0.01",
		},
	}
}
