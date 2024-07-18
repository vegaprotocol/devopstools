package bots

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"code.vegaprotocol.io/vega/core/netparams"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	v1 "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"
	"code.vegaprotocol.io/vega/wallet/wallet"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type AmmArgs struct {
	*Args

	timeout            time.Duration
	marketId           string
	poolSize           int
	vegaRecoveryPhrase string
}

const (
	ammSubmissionAmount     = 5000
	ammWhaleTopUpMultiplier = 5
	ammAgentTopUpMultiplier = 3
)

var ammArgs AmmArgs

var ammCmd = &cobra.Command{
	Use:   "amm",
	Short: "Setup the amm pool on given markets",
	Long:  "Setup the amm pool on given markets",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setupAmm(ammArgs); err != nil {
			ammArgs.Logger.Error("Could not top up bots", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	ammArgs.Args = &args

	ammCmd.PersistentFlags().DurationVarP(&ammArgs.timeout, "timeout", "t", 15*time.Minute, "Timeout for the amm command")
	ammCmd.PersistentFlags().StringVarP(&ammArgs.marketId, "market-id", "m", "", "The market id")
	ammCmd.PersistentFlags().IntVarP(&ammArgs.poolSize, "size", "s", 20, "Timeout for the top up command")
	ammCmd.PersistentFlags().StringVarP(&ammArgs.vegaRecoveryPhrase, "vega-recovery-phrase", "p", "", "The vega recovery phrase for the amm pool submitter wallet")

	Cmd.AddCommand(ammCmd)
}

func setupAmm(args AmmArgs) error {

	ctx, cancelCommand := context.WithTimeout(context.Background(), topUpArgs.timeout)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(ctx, args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Info("Network file loaded", zap.String("name", cfg.Name.String()))

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Info("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Debug("Looking for healthy gRPC endpoints...")
	healthyEndpoints := tools.FilterHealthyGRPCEndpoints(endpoints)
	if len(healthyEndpoints) == 0 {
		return fmt.Errorf("no healthy gRPC endpoint found on configured datanodes")
	}
	logger.Info("Healthy gRPC endpoints found", zap.Strings("endpoints", healthyEndpoints))

	datanodeClient := datanode.New(healthyEndpoints, 3*time.Second, args.Logger.Named("datanode"))

	logger.Info("Connecting to a datanode's gRPC endpoint...")
	dialCtx, cancelDialing := context.WithTimeout(ctx, 2*time.Second)
	defer cancelDialing()
	datanodeClient.MustDialConnection(dialCtx) // blocking
	logger.Info("Connected to a datanode's gRPC node", zap.String("node", datanodeClient.Target()))

	logger.Debug("Retrieving network parameters...")
	networkParams, err := datanodeClient.ListNetworkParameters(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}
	logger.Info("Network parameters retrieved")

	chainClients, err := ethereum.NewChainClients(ctx, cfg, networkParams, logger)
	if err != nil {
		return err
	}

	assets, err := datanodeClient.ListAssets(ctx)
	if err != nil {
		return fmt.Errorf("could not list assets from datanode: %w", err)
	}

	if len(assets) == 0 {
		return errors.New("no asset found on datanode")
	}

	activeMarkets, err := datanodeClient.GetAllMarketsWithState(ctx, datanode.ActiveMarkets)
	if err != nil {
		return fmt.Errorf("failed to get active markets")
	}

	var marketDetails *vegapb.Market
	for idx, market := range activeMarkets {
		if market.Id == args.marketId {
			marketDetails = activeMarkets[idx]
			break
		}
	}

	if marketDetails == nil {
		return fmt.Errorf("failed to find the %s market", args.marketId)
	}

	marketAssetsIds := datanode.AssetsIdsByMarket(marketDetails)
	if len(marketAssetsIds) == 0 {
		return fmt.Errorf("cannot obtain asset ids for given market")
	}

	submitterWallet, err := vega.LoadWallet("submitter", args.vegaRecoveryPhrase)
	if err != nil {
		return fmt.Errorf("failed to load submitter wallet: %w", err)
	}

	if err := vega.GenerateKeysUpToIndex(submitterWallet, uint32(args.poolSize+1)); err != nil {
		return fmt.Errorf("failed to generate %d keys from the submitter wallet: %w", args.poolSize, err)
	}

	topUpRegistry, err := determineAmmAgentsTopUpAmount(ctx, datanodeClient, assets, submitterWallet, args.poolSize, marketAssetsIds)
	if err != nil {
		return fmt.Errorf("failed to determine amm agents top up amount: %w", err)
	}

	whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	if err := vega.GenerateKeysUpToKey(whaleWallet, cfg.Network.Wallets.VegaTokenWhale.PublicKey); err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	whalePublicKey := cfg.Network.Wallets.VegaTokenWhale.PublicKey

	logger.Debug("Determining amounts to top up for whale...")
	whaleTopUpsByAsset, err := determineAmountsToTopUpForWhale(ctx, datanodeClient, whalePublicKey, assets, topUpRegistry, logger)
	if err != nil {
		return fmt.Errorf("failed to determine top up amounts for the whale: %w", err)
	}
	logger.Info("Amounts to top up for whale computed")

	if len(whaleTopUpsByAsset) == 0 {
		logger.Info("No top-up required for the whale")
	} else {
		logger.Debug("Depositing assets to whale...")
		if err := depositAssetsToWhale(ctx, whaleTopUpsByAsset, assets, datanodeClient, whalePublicKey, chainClients, logger); err != nil {
			return err
		}
		logger.Info("Whale has now enough funds to transfer to trading bots")
	}

	transfersToDo := countTransfersToDo(topUpRegistry)
	if networkParams.GetMaxTransfersPerEpoch() < int64(transfersToDo) {
		logger.Debug("Preparing network for transfers from whale to trading bots...")
		updateParams := map[string]string{
			netparams.TransferMaxCommandsPerEpoch: fmt.Sprintf("%d", transfersToDo*30),
		}
		if _, err := vega.UpdateNetworkParameters(ctx, whaleWallet, whalePublicKey, datanodeClient, updateParams, logger); err != nil {
			return fmt.Errorf("failed to prepare network for transfers: %w", err)
		}
		logger.Info("Network ready for transfers from whale to trading bots")
	}

	logger.Info("Transferring assets from whale to trading bots...")
	if err := transferAssetsFromWhaleToBots(ctx, datanodeClient, whaleWallet, whalePublicKey, topUpRegistry, logger); err != nil {
		return fmt.Errorf("failed to transfer assets from whale to one or more bots: %w", err)
	}
	logger.Info("Transfers done")

	marketData, err := datanodeClient.GetLatestMarketDate(ctx, args.marketId)
	if err != nil {
		return fmt.Errorf("failed to get market data for market %s: %w", args.marketId, err)
	}

	markPrice, _ := big.NewFloat(0).SetString(marketData.MarkPrice),
	upperBound := big.NewFloat(0).Mul(
		markPrice,
		big.NewFloat(1.01),
	).String()
	lowerBound := big.NewFloat(0).Mul(
		markPrice,
		big.NewFloat(0.99),
	).String()
	leverage := "10"

	// // Market ID for which to create an AMM.
	// MarketId string `protobuf:"bytes,1,opt,name=market_id,json=marketId,proto3" json:"market_id,omitempty"`
	// // Amount to be committed to the AMM.
	// CommitmentAmount string `protobuf:"bytes,2,opt,name=commitment_amount,json=commitmentAmount,proto3" json:"commitment_amount,omitempty"`
	// // Slippage tolerance used for rebasing the AMM if its base price crosses with existing order
	// SlippageTolerance string `protobuf:"bytes,3,opt,name=slippage_tolerance,json=slippageTolerance,proto3" json:"slippage_tolerance,omitempty"`
	// // Concentrated liquidity parameters defining the shape of the AMM's volume curves.
	// ConcentratedLiquidityParameters *SubmitAMM_ConcentratedLiquidityParameters `protobuf:"bytes,4,opt,name=concentrated_liquidity_parameters,json=concentratedLiquidityParameters,proto3" json:"concentrated_liquidity_parameters,omitempty"`
	// // Nominated liquidity fee factor, which is an input to the calculation of taker fees on the market.
	// ProposedFee string `protobuf:"bytes,5,opt,name=proposed_fee,json=proposedFee,proto3" json:"proposed_fee,omitempty"`
	for assetId := range topUpRegistry {
		for partyId, amount := range topUpRegistry[assetId].AmountsByParty {

			walletTxReq := walletpb.SubmitTransactionRequest{
				Command: &walletpb.SubmitTransactionRequest_SubmitAmm{
					SubmitAmm: &v1.SubmitAMM{
						MarketId:          args.marketId,
						CommitmentAmount:  amount.String(),
						SlippageTolerance: "0.5",
						ProposedFee:       "0.001",
						ConcentratedLiquidityParameters: &v1.SubmitAMM_ConcentratedLiquidityParameters{
							UpperBound:           &upperBound,
							LowerBound:           &lowerBound,
							Base:                 marketData.MarkPrice,
							LeverageAtUpperBound: &leverage,
							LeverageAtLowerBound: &leverage,
						},
					},
				},
			}

			logger.Info("Sending transaction to the network")
			resp, err := walletpkg.SendTransaction(ctx, whaleWallet, whaleKey, &walletTxReq, datanodeClient)
			if err != nil {
				return fmt.Errorf("failed to submit batch proposal with signature: %w", err)
			}
			logger.Sugar().Infof("Batch proposal transaction ID: %s", resp.TxHash)

		}
	}

	return nil
}

func determineAmmAgentsTopUpAmount(
	ctx context.Context,
	datanodeClient vegaapi.DataNodeClient,
	assets map[string]*vegapb.AssetDetails,
	submitterWallet wallet.Wallet,
	poolSize int,
	assetsIds []string,
) (map[string]AssetToTopUp, error) {
	result := map[string]AssetToTopUp{}

	ammAgentsPubKeys := submitterWallet.ListPublicKeys()
	if len(ammAgentsPubKeys) < poolSize {
		return nil, fmt.Errorf("not enough wallets generated for amm pool of %d", poolSize)
	}

	for i := 0; i <= poolSize; i++ {
		agentPubKey := ammAgentsPubKeys[i].Key()

		for _, assetId := range assetsIds {
			assetDetails, assetFound := assets[assetId]
			if !assetFound {
				return nil, fmt.Errorf("cannot find asset(%s) required for amm agent: %s", assetId, agentPubKey)
			}

			if assetDetails.GetErc20() == nil {
				return nil, fmt.Errorf("asset %s is not erc20 asset", assetId)
			}

			if _, assetResultFound := result[assetId]; !assetResultFound {
				result[assetId] = AssetToTopUp{
					Symbol:          assetDetails.Symbol,
					ContractAddress: assetDetails.GetErc20().ContractAddress,
					VegaAssetId:     assetId,
					TotalAmount:     types.NewAmount(assetDetails.Decimals),
					AmountsByParty:  map[string]*types.Amount{},
				}
			}

			currentAmmAgentGeneralAccount, err := datanodeClient.ListAccounts(ctx, agentPubKey, vegapb.AccountType_ACCOUNT_TYPE_GENERAL, &assetId)
			if err != nil {
				return nil, fmt.Errorf("failed to get general account balance for %s part", agentPubKey)
			}

			currentAmountWithoutZeros := big.NewInt(0)
			if len(currentAmmAgentGeneralAccount) > 0 {
				currentAmountWithoutZeros = big.NewInt(0).Div(
					currentAmmAgentGeneralAccount[0].Balance,
					big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(assetDetails.Decimals)), big.NewInt(0)),
				)
			}

			if currentAmountWithoutZeros.Cmp(big.NewInt(ammSubmissionAmount)) < 0 {
				result[assetId].TotalAmount.Add(big.NewFloat(0).Mul(big.NewFloat(ammAgentTopUpMultiplier), big.NewFloat(float64(ammSubmissionAmount)))) // Maybe we should add only missing amount, but for now meh...

				result[assetId].AmountsByParty[agentPubKey] = types.NewAmountFromMainUnit(
					big.NewFloat(0).Mul(big.NewFloat(ammAgentTopUpMultiplier), big.NewFloat(float64(ammSubmissionAmount))),
					assetDetails.Decimals,
				)
			}
		}
	}

	// multiplier for whale top up
	for _, assetId := range assetsIds {
		result[assetId].TotalAmount.Mul(big.NewFloat(float64(ammWhaleTopUpMultiplier)))
	}

	return result, nil
}

// logger.Sugar().Info("Getting traders from bots API: %s", cfg.Bots.Research.RESTURL)
// tradingBots, err := bots.RetrieveTradingBots(ctx, cfg.Bots.Research.RESTURL, cfg.Bots.Research.APIKey, logger.Named("trading-bots"))
// if err != nil {
// 	return fmt.Errorf("failed to retrieve trading bots: %w", err)
// }
// logger.Debug("Trading bots found", zap.Strings("traders", maps.Keys(tradingBots)))

// logger.Info("Determining amounts to top up for trading bots...")
// topUpsByAsset, err := determineAmountsToTopUpByAsset(assets, tradingBots, logger)
// if err != nil {
// 	return fmt.Errorf("failed to determine top up amounts for the each traders: %w", err)
// }
// logger.Info("Amounts to top up for trading bots computed")

// if len(topUpsByAsset) == 0 {
// 	logger.Info("No top-up required for the trading bots")
// 	return nil
// }

// whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
// if err != nil {
// 	return fmt.Errorf("could not initialized whale wallet: %w", err)
// }

// if err := vega.GenerateKeysUpToKey(whaleWallet, cfg.Network.Wallets.VegaTokenWhale.PublicKey); err != nil {
// 	return fmt.Errorf("could not initialized whale wallet: %w", err)
// }

// whalePublicKey := cfg.Network.Wallets.VegaTokenWhale.PublicKey

// logger.Debug("Determining amounts to top up for whale...")
// whaleTopUpsByAsset, err := determineAmountsToTopUpForWhale(ctx, datanodeClient, whalePublicKey, assets, topUpsByAsset, logger)
// if err != nil {
// 	return fmt.Errorf("failed to determine top up amounts for the whale: %w", err)
// }
// logger.Info("Amounts to top up for whale computed")

// if len(whaleTopUpsByAsset) == 0 {
// 	logger.Info("No top-up required for the whale")
// } else {
// 	logger.Debug("Depositing assets to whale...")
// 	if err := depositAssetsToWhale(ctx, whaleTopUpsByAsset, assets, datanodeClient, whalePublicKey, chainClients, logger); err != nil {
// 		return err
// 	}
// 	logger.Info("Whale has now enough funds to transfer to trading bots")
// }

// transfersToDo := countTransfersToDo(topUpsByAsset)
// if networkParams.GetMaxTransfersPerEpoch() < int64(transfersToDo) {
// 	logger.Debug("Preparing network for transfers from whale to trading bots...")
// 	updateParams := map[string]string{
// 		netparams.TransferMaxCommandsPerEpoch: fmt.Sprintf("%d", transfersToDo*30),
// 	}
// 	if _, err := vega.UpdateNetworkParameters(ctx, whaleWallet, whalePublicKey, datanodeClient, updateParams, logger); err != nil {
// 		return fmt.Errorf("failed to prepare network for transfers: %w", err)
// 	}
// 	logger.Info("Network ready for transfers from whale to trading bots")
// }

// logger.Info("Transferring assets from whale to trading bots...")
// if err := transferAssetsFromWhaleToBots(ctx, datanodeClient, whaleWallet, whalePublicKey, topUpsByAsset, logger); err != nil {
// 	return fmt.Errorf("failed to transfer assets from whale to one or more bots: %w", err)
// }
// logger.Info("Transfers done")

// logger.Info("Trading bots have been topped up successfully")

// func depositAssetsToWhale(ctx context.Context, whaleTopUpsByAsset map[string]*types.Amount, assets map[string]*vegapb.AssetDetails, datanodeClient *datanode.DataNode, publicKey string, chainClients *ethereum.ChainsClient, logger *zap.Logger) error {
// 	for assetID, requiredAmount := range whaleTopUpsByAsset {
// 		asset := assets[assetID]

// 		erc20Details := asset.GetErc20()
// 		if erc20Details == nil {
// 			return fmt.Errorf("asset %q is not an ERC20 token", asset.Name)
// 		}

// 		logger.Info("Depositing asset on whale wallet...",
// 			zap.String("asset-name", asset.Name),
// 			zap.String("asset-contract-address", erc20Details.ContractAddress),
// 			zap.String("amount-to-deposit", requiredAmount.String()),
// 			zap.String("chain-id", erc20Details.ChainId),
// 		)

// 		var chainClient *ethereum.ChainClient
// 		switch erc20Details.ChainId {
// 		case chainClients.PrimaryChain.ID():
// 			chainClient = chainClients.PrimaryChain
// 		case chainClients.EVMChain.ID():
// 			chainClient = chainClients.EVMChain
// 		default:
// 			return fmt.Errorf("asset with chain ID %q does not match any configured Ethereum chain", erc20Details.ChainId)
// 		}

// 		deposits := map[string]*types.Amount{
// 			publicKey: requiredAmount,
// 		}

// 		if err := chainClient.DepositERC20AssetFromMinter(ctx, erc20Details.ContractAddress, deposits); err != nil {
// 			return fmt.Errorf("failed to deposit asset %q on whale %s: %w", asset.Name, publicKey, err)
// 		}

// 		logger.Info("Asset deposited on whale wallet",
// 			zap.String("asset-name", asset.Name),
// 			zap.String("asset-contract-address", erc20Details.ContractAddress),
// 			zap.String("amount-deposited", requiredAmount.String()),
// 			zap.String("chain-id", erc20Details.ChainId),
// 		)
// 	}

// 	logger.Info("Ensuring whale received funds")
// 	if err := ensureWhaleReceivedFunds(ctx, datanodeClient, publicKey, whaleTopUpsByAsset, logger); err != nil {
// 		return err
// 	}
// 	logger.Info("Foound ")

// 	return nil
// }

// func determineAmountsToTopUpByAsset(assets map[string]*vegapb.AssetDetails, botsMap map[string]bots.TradingBot, logger *zap.Logger) (map[string]AssetToTopUp, error) {
// 	topUpRegistry := map[string]AssetToTopUp{}

// 	wantedTokenEntries := []bots.BotTraderWantedToken{}

// 	for _, bot := range botsMap {
// 		wantedTokenEntries = append(wantedTokenEntries, bot.Parameters.WantedTokens...)
// 	}

// 	for _, wantedToken := range wantedTokenEntries {
// 		assetDetails, assetExists := assets[wantedToken.VegaAssetId]
// 		if !assetExists {
// 			return nil, fmt.Errorf(
// 				"trading bot is using the asset %q but it does not exist on the network",
// 				wantedToken.VegaAssetId,
// 			)
// 		}

// 		if _, assetRegistryExists := topUpRegistry[wantedToken.VegaAssetId]; !assetRegistryExists {
// 			topUpRegistry[wantedToken.VegaAssetId] = AssetToTopUp{
// 				Symbol:          assetDetails.Symbol,
// 				ContractAddress: wantedToken.AssetERC20Asset,
// 				VegaAssetId:     wantedToken.VegaAssetId,
// 				TotalAmount:     types.NewAmount(assetDetails.Decimals),
// 				AmountsByParty:  map[string]*types.Amount{},
// 			}
// 		}

// 		requiredTopUp := computeRequiredAmount(
// 			wantedToken.Balance,
// 			wantedToken.WantedTokens,
// 		)

// 		if requiredTopUp.Cmp(big.NewFloat(0)) == 0 {
// 			logger.Info(
// 				"Party does not need a top up for the asset",
// 				zap.String("party-id", wantedToken.PartyId),
// 				zap.Float64("current-funds", wantedToken.Balance),
// 				zap.Float64("required-funds", wantedToken.WantedTokens),
// 				zap.String("asset", assetDetails.Name),
// 			)

// 			continue
// 		}

// 		currentEntry := topUpRegistry[wantedToken.VegaAssetId]

// 		if _, partyForAssetExist := topUpRegistry[wantedToken.VegaAssetId].AmountsByParty[wantedToken.PartyId]; !partyForAssetExist {
// 			currentEntry.AmountsByParty[wantedToken.PartyId] = types.NewAmountFromMainUnit(big.NewFloat(0), assetDetails.Decimals)
// 		}

// 		currentEntry.AmountsByParty[wantedToken.PartyId].Add(requiredTopUp)
// 		currentEntry.TotalAmount.Add(requiredTopUp)

// 		logger.Info(
// 			"Party requires a top up for the asset",
// 			zap.String("party-id", wantedToken.PartyId),
// 			zap.Float64("current-funds", wantedToken.Balance),
// 			zap.Float64("required-funds", wantedToken.WantedTokens),
// 			zap.String("asset", assetDetails.Name),
// 			zap.String("top-up", requiredTopUp.String()),
// 		)

// 		topUpRegistry[wantedToken.VegaAssetId] = currentEntry
// 	}

// 	return topUpRegistry, nil
// }

// func computeRequiredAmount(currentBalance, wantedBalance float64) *big.Float {
// 	if wantedBalance < 0.01 || wantedBalance > currentBalance {
// 		return big.NewFloat(wantedBalance * TopUpFactorForTradingBots)
// 	}

// 	// Top up not required.
// 	// Do not write to that object.
// 	return big.NewFloat(0)
// }

// func determineAmountsToTopUpForWhale(ctx context.Context, datanodeClient *datanode.DataNode, publicKey string, assets map[string]*vegapb.AssetDetails, topUpsByAsset map[string]AssetToTopUp, logger *zap.Logger) (map[string]*types.Amount, error) {
// 	result := map[string]*types.Amount{}

// 	for _, assetToTopUp := range topUpsByAsset {
// 		assetDetails, assetExists := assets[assetToTopUp.VegaAssetId]
// 		if !assetExists {
// 			return nil, fmt.Errorf("whale needs to topup the asset %q but it does not exist on the network", assetToTopUp.VegaAssetId)
// 		}

// 		whaleFundsAsSubUnit, err := datanodeClient.GeneralAccountBalance(ctx, publicKey, assetToTopUp.VegaAssetId)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to retrieve whale's general account balance for asset %q: %w", assetToTopUp.VegaAssetId, err)
// 		}
// 		whaleFunds := types.NewAmountFromSubUnit(whaleFundsAsSubUnit, assetDetails.Decimals)

// 		if whaleFunds.Cmp(assetToTopUp.TotalAmount) > -1 {
// 			logger.Info("Whale does not need a top up the asset",
// 				zap.String("asset", assetToTopUp.Symbol),
// 				zap.String("current-funds", whaleFunds.String()),
// 				zap.String("required-funds", assetToTopUp.TotalAmount.String()),
// 			)
// 			continue
// 		}

// 		topUpAmount := assetToTopUp.TotalAmount.Copy()
// 		topUpAmount.Mul(big.NewFloat(TopUpFactorForWhale))

// 		logger.Info("Whale requires a top up for the asset",
// 			zap.String("asset", assetToTopUp.Symbol),
// 			zap.String("current-funds", whaleFunds.String()),
// 			zap.String("required-funds", assetToTopUp.TotalAmount.String()),
// 			zap.String("top-up", topUpAmount.String()),
// 		)

// 		result[assetToTopUp.VegaAssetId] = topUpAmount
// 	}

// 	return result, nil
// }

// func countTransfersToDo(topUpRegistry map[string]AssetToTopUp) int {
// 	transferNumbers := 0

// 	for _, entry := range topUpRegistry {
// 		transferNumbers = transferNumbers + len(entry.AmountsByParty)
// 	}

// 	return transferNumbers + (10 * transferNumbers / 100)
// }

// func ensureWhaleReceivedFunds(
// 	ctx context.Context,
// 	datanodeClient *datanode.DataNode,
// 	publicKey string,
// 	whaleTopUpsByAsset map[string]*types.Amount,
// 	logger *zap.Logger,
// ) error {
// 	ticker := time.NewTicker(30 * time.Second)

// 	for {
// 		allDepositsFinalized := true

// 		for assetID, requiredAmount := range whaleTopUpsByAsset {
// 			requiredAmountAsSubUnit := requiredAmount.AsSubUnit()
// 			logger.Sugar().Infof("Checking if deposit of asset %s has been finalized", assetID)

// 			balanceAsSubUnit, err := datanodeClient.GeneralAccountBalance(ctx, publicKey, assetID)
// 			if err != nil {
// 				logger.Warn(fmt.Sprintf("failed to retrieve whale's general account balance for asset %q", assetID), zap.Error(err))
// 				allDepositsFinalized = false
// 			}

// 			if balanceAsSubUnit.Cmp(requiredAmountAsSubUnit) < 0 {
// 				logger.Sugar().Infof("Deposit for asset %s not finalized yet", assetID)
// 				allDepositsFinalized = false
// 			}
// 		}

// 		if allDepositsFinalized {
// 			logger.Info("All deposits finalized")
// 			return nil
// 		}

// 		select {
// 		case <-ticker.C:
// 			continue
// 		case <-ctx.Done():
// 			return fmt.Errorf("deposit not finalized in given time")
// 		}
// 	}
// }

// func transferAssetsFromWhaleToBots(ctx context.Context, datanodeClient *datanode.DataNode, whaleWallet wallet.Wallet, whalePublicKey string, registry map[string]AssetToTopUp, logger *zap.Logger) error {
// 	result := &multierror.Error{}

// 	for assetID, entry := range registry {
// 		for botPartyId, amount := range entry.AmountsByParty {
// 			err := tools.RetryRun(15, 6*time.Second, func() error {
// 				request := walletpb.SubmitTransactionRequest{
// 					Command: &walletpb.SubmitTransactionRequest_Transfer{
// 						Transfer: &vegacmd.Transfer{
// 							Reference:       fmt.Sprintf("Transfer from whale to %s", botPartyId),
// 							FromAccountType: vegapb.AccountType_ACCOUNT_TYPE_GENERAL,
// 							ToAccountType:   vegapb.AccountType_ACCOUNT_TYPE_GENERAL,
// 							To:              botPartyId,
// 							Asset:           assetID,
// 							Amount:          amount.StringWithDecimals(),
// 							Kind: &vegacmd.Transfer_OneOff{
// 								OneOff: &vegacmd.OneOffTransfer{
// 									DeliverOn: 0,
// 								},
// 							},
// 						},
// 					},
// 				}

// 				resp, err := walletpkg.SendTransaction(ctx, whaleWallet, whalePublicKey, &request, datanodeClient)
// 				if err != nil {
// 					return fmt.Errorf("failed to send the submit transfer transaction for bot %q: %w", botPartyId, err)
// 				}

// 				if !resp.Success {
// 					return fmt.Errorf("the transfer transaction(%s) failed for bot %q with reason %s: %w", resp.TxHash, botPartyId, resp.Data, err)
// 				}

// 				logger.Info("Assets have been sent from whale to the trading bot",
// 					zap.String("asset", entry.Symbol),
// 					zap.String("amount", amount.String()),
// 					zap.String("bot", botPartyId),
// 					zap.String("transaction", resp.TxHash),
// 				)

// 				return nil
// 			})
// 			if err != nil {
// 				result.Errors = append(result.Errors, fmt.Errorf("topping up bot %q failed: %w", botPartyId, err))

// 				logger.Error("Topping up bot failed",
// 					zap.String("bot", botPartyId),
// 					zap.Error(err),
// 				)
// 			}

// 			time.Sleep(100 * time.Millisecond)
// 		}
// 	}

// 	return result.ErrorOrNil()
// }
