package bots

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/api"
	"github.com/vegaprotocol/devopstools/bots"
	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/networktools"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"code.vegaprotocol.io/vega/core/netparams"
	vegapb "code.vegaprotocol.io/vega/protos/vega"
	vegacmd "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"
	"code.vegaprotocol.io/vega/wallet/wallet"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

const (
	TopUpFactorForTradingBots = 3.0
	TopUpFactorForWhale       = 10.0
)

type TopUpArgs struct {
	*Args
}

var topUpArgs TopUpArgs

var topUpCmd = &cobra.Command{
	Use:   "top-up",
	Short: "Top up bots on network with Vega transfer",
	Long:  "Top up bots on network with Vega transfer",
	Run: func(cmd *cobra.Command, args []string) {
		if err := TopUpBots(topUpArgs); err != nil {
			topUpArgs.Logger.Error("Could not top up bots", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	topUpArgs.Args = &args

	Cmd.AddCommand(topUpCmd)
}

type AssetToTopUp struct {
	Symbol          string
	ContractAddress string
	VegaAssetId     string
	TotalAmount     *types.Amount
	AmountsByParty  map[string]*types.Amount
}

func TopUpBots(args TopUpArgs) error {
	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Debug("Network file loaded", zap.String("name", cfg.Name.String()))

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Debug("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Debug("Looking for healthy gRPC endpoints...")
	healthyEndpoints := networktools.FilterHealthyGRPCEndpoints(endpoints)
	if len(healthyEndpoints) == 0 {
		return fmt.Errorf("no healthy gRPC endpoint found on configured datanodes")
	}
	logger.Debug("Healthy gRPC endpoints found", zap.Strings("endpoints", healthyEndpoints))

	datanodeClient := datanode.New(healthyEndpoints, 3*time.Second, args.Logger.Named("datanode"))

	logger.Debug("Connecting to a datanode's gRPC endpoint...")
	dialCtx, cancelDialing := context.WithTimeout(ctx, 2*time.Second)
	defer cancelDialing()
	datanodeClient.MustDialConnection(dialCtx) // blocking
	logger.Debug("Connected to a datanode's gRPC node", zap.String("node", datanodeClient.Target()))

	logger.Debug("Retrieving network parameters...")
	networkParams, err := datanodeClient.GetAllNetworkParameters()
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}
	logger.Debug("Network parameters retrieved")

	chainClients, err := ethereum.NewChainClients(ctx, cfg, networkParams, logger)
	if err != nil {
		return err
	}

	logger.Debug("Listing assets from datanode...")
	assets, err := datanodeClient.ListAssets(ctx)
	if err != nil {
		return fmt.Errorf("could not list assets from datanode: %w", err)
	}
	if len(assets) == 0 {
		return errors.New("no asset found on datanode")
	}
	logger.Debug("Assets found", zap.Strings("assets", maps.Keys(assets)))

	tradingBots, err := bots.RetrieveTradingBots(ctx, cfg.Bots.Trading.RESTURL, cfg.Bots.Trading.APIKey, logger.Named("trading-bots"))
	if err != nil {
		return fmt.Errorf("failed to retrieve trading bots: %w", err)
	}
	logger.Debug("Trading bots found", zap.Strings("traders", maps.Keys(tradingBots)))

	logger.Debug("Determining amounts to top up for trading bots...")
	topUpsByAsset, err := determineAmountsToTopUpByAsset(assets, tradingBots, logger)
	if err != nil {
		return fmt.Errorf("failed to determine top up amounts for the each traders: %w", err)
	}
	logger.Debug("Amounts to top up for trading bots computed")

	if len(topUpsByAsset) == 0 {
		logger.Info("No top-up required for the trading bots")
		return nil
	}

	whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	if err := vega.GenerateKeysUpToKey(whaleWallet, cfg.Network.Wallets.VegaTokenWhale.PublicKey); err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	whalePublicKey := cfg.Network.Wallets.VegaTokenWhale.PublicKey

	whaleClient := api.NewPartyClient(datanodeClient, whalePublicKey)

	logger.Debug("Determining amounts to top up for whale...")
	whaleTopUpsByAsset, err := determineAmountsToTopUpForWhale(ctx, whaleClient, assets, topUpsByAsset, logger)
	if err != nil {
		return fmt.Errorf("failed to determine top up amounts for the whale: %w", err)
	}
	logger.Debug("Amounts to top up for whale computed")

	if len(whaleTopUpsByAsset) == 0 {
		logger.Debug("No top-up required for the whale")
	} else {
		logger.Debug("Depositing assets to whale...")
		if err := depositAssetsToWhale(ctx, whaleTopUpsByAsset, assets, whaleClient, chainClients, logger); err != nil {
			return err
		}
		logger.Debug("Whale has now enough funds to transfer to trading bots")
	}

	logger.Debug("Preparing network for transfers from whale to trading bots...")
	numberOfTransfers := countTransfersToDo(topUpsByAsset)
	if err := prepareNetworkForTransfers(ctx, whaleWallet, whalePublicKey, datanodeClient, numberOfTransfers, logger); err != nil {
		return fmt.Errorf("failed to prepare network for transfers: %w", err)
	}
	logger.Debug("Network ready for transfers from whale to trading bots")

	logger.Debug("Transferring assets from whale to trading bots...")
	if err := transferAssetsFromWhaleToBots(ctx, datanodeClient, whaleWallet, whalePublicKey, topUpsByAsset, logger); err != nil {
		return fmt.Errorf("failed to transfer assets from whale to one or more bots: %w", err)
	}
	logger.Debug("Transfers done")

	logger.Info("Trading bots have been topped up successfully")

	return nil
}

func depositAssetsToWhale(ctx context.Context, whaleTopUpsByAsset map[string]*types.Amount, assets map[string]*vegapb.AssetDetails, whaleClient *api.PartyClient, chainClients *ethereum.ChainsClient, logger *zap.Logger) error {
	for assetID, requiredAmount := range whaleTopUpsByAsset {
		asset := assets[assetID]

		erc20Details := asset.GetErc20()
		if erc20Details == nil {
			return fmt.Errorf("asset %q is not an ERC20 token", asset.Name)
		}

		logger.Debug("Depositing asset on whale wallet...",
			zap.String("asset-name", asset.Name),
			zap.String("asset-contract-address", erc20Details.ContractAddress),
			zap.String("amount-to-deposit", requiredAmount.String()),
			zap.String("chain-id", erc20Details.ChainId),
		)

		var chainClient *ethereum.ChainClient
		switch erc20Details.ChainId {
		case chainClients.PrimaryChain.ID():
			chainClient = chainClients.PrimaryChain
		case chainClients.EVMChain.ID():
			chainClient = chainClients.EVMChain
		default:
			return fmt.Errorf("asset with chain ID %q does not match any configured Ethereum chain", erc20Details.ChainId)
		}

		deposits := map[string]*types.Amount{
			whaleClient.PartyID(): requiredAmount,
		}

		if err := chainClient.DepositERC20AssetFromMinter(ctx, erc20Details.ContractAddress, deposits); err != nil {
			return fmt.Errorf("failed to deposit asset %q on whale %s: %w", asset.Name, whaleClient.PartyID(), err)
		}

		logger.Info("Asset deposited on whale wallet",
			zap.String("asset-name", asset.Name),
			zap.String("asset-contract-address", erc20Details.ContractAddress),
			zap.String("amount-deposited", requiredAmount.String()),
			zap.String("chain-id", erc20Details.ChainId),
		)
	}

	if err := ensureWhaleReceivedFunds(ctx, whaleClient, whaleTopUpsByAsset); err != nil {
		return err
	}

	return nil
}

func determineAmountsToTopUpByAsset(assets map[string]*vegapb.AssetDetails, bots map[string]bots.TradingBot, logger *zap.Logger) (map[string]AssetToTopUp, error) {
	topUpRegistry := map[string]AssetToTopUp{}

	for _, bot := range bots {
		assetDetails, assetExists := assets[bot.Parameters.SettlementVegaAssetID]
		if !assetExists {
			return nil, fmt.Errorf(
				"trading bot is using the asset %q but it does not exist on the network",
				bot.Parameters.SettlementVegaAssetID,
			)
		}

		if _, registryExists := topUpRegistry[bot.Parameters.SettlementVegaAssetID]; !registryExists {
			topUpRegistry[bot.Parameters.SettlementVegaAssetID] = AssetToTopUp{
				Symbol:          assetDetails.Symbol,
				ContractAddress: bot.Parameters.SettlementEthereumContractAddress,
				VegaAssetId:     bot.Parameters.SettlementVegaAssetID,
				TotalAmount:     types.NewAmount(assetDetails.Decimals),
				AmountsByParty:  map[string]*types.Amount{},
			}
		}

		currentEntry := topUpRegistry[bot.Parameters.SettlementVegaAssetID]

		requiredTopUp := computeRequiredAmount(bot.Parameters.CurrentBalance, bot.Parameters.WantedTokens)

		if requiredTopUp.Cmp(big.NewFloat(0)) == 0 {
			logger.Debug(
				"Party does not need a top up for the asset",
				zap.String("party-id", bot.PubKey),
				zap.Float64("current-funds", bot.Parameters.CurrentBalance),
				zap.Float64("required-funds", bot.Parameters.WantedTokens),
				zap.String("asset", assetDetails.Name),
			)

			continue
		}

		currentEntry.AmountsByParty[bot.PubKey] = types.NewAmountFromMainUnit(requiredTopUp, assetDetails.Decimals)
		currentEntry.TotalAmount.Add(requiredTopUp)

		logger.Info(
			"Party requires a top up for the asset",
			zap.String("party-id", bot.PubKey),
			zap.Float64("current-funds", bot.Parameters.CurrentBalance),
			zap.Float64("required-funds", bot.Parameters.WantedTokens),
			zap.String("asset", assetDetails.Name),
			zap.String("top-up", requiredTopUp.String()),
		)

		topUpRegistry[bot.Parameters.SettlementVegaAssetID] = currentEntry
	}

	return topUpRegistry, nil
}

func computeRequiredAmount(currentBalance, wantedBalance float64) *big.Float {
	if wantedBalance < 0.01 || wantedBalance > currentBalance {
		return big.NewFloat(wantedBalance * TopUpFactorForTradingBots)
	}

	// Top up not required.
	// Do not write to that object.
	return big.NewFloat(0)
}

func determineAmountsToTopUpForWhale(ctx context.Context, whaleClient *api.PartyClient, assets map[string]*vegapb.AssetDetails, topUpsByAsset map[string]AssetToTopUp, logger *zap.Logger) (map[string]*types.Amount, error) {
	result := map[string]*types.Amount{}

	for _, assetToTopUp := range topUpsByAsset {
		assetDetails, assetExists := assets[assetToTopUp.VegaAssetId]
		if !assetExists {
			return nil, fmt.Errorf("whale needs to topup the asset %q but it does not exist on the network", assetToTopUp.VegaAssetId)
		}

		whaleFundsAsSubUnit, err := whaleClient.GeneralAccountBalanceForAsset(ctx, assetToTopUp.VegaAssetId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve whale's general account balance for asset %q: %w", assetToTopUp.VegaAssetId, err)
		}
		whaleFunds := types.NewAmountFromSubUnit(whaleFundsAsSubUnit, assetDetails.Decimals)

		if whaleFunds.Cmp(assetToTopUp.TotalAmount) > -1 {
			logger.Debug("Whale does not need a top up the asset",
				zap.String("asset", assetToTopUp.Symbol),
				zap.String("current-funds", whaleFunds.String()),
				zap.String("required-funds", assetToTopUp.TotalAmount.String()),
			)
			continue
		}

		topUpAmount := assetToTopUp.TotalAmount.Copy()
		topUpAmount.Mul(big.NewFloat(TopUpFactorForWhale))

		logger.Info("Whale requires a top up for the asset",
			zap.String("asset", assetToTopUp.Symbol),
			zap.String("current-funds", whaleFunds.String()),
			zap.String("required-funds", assetToTopUp.TotalAmount.String()),
			zap.String("top-up", topUpAmount.String()),
		)

		result[assetToTopUp.VegaAssetId] = topUpAmount
	}

	return result, nil
}

func countTransfersToDo(topUpRegistry map[string]AssetToTopUp) int {
	transferNumbers := 0

	for _, entry := range topUpRegistry {
		transferNumbers = transferNumbers + len(entry.AmountsByParty)
	}

	return transferNumbers + (10 * transferNumbers / 100)
}

func prepareNetworkForTransfers(ctx context.Context, whaleWallet wallet.Wallet, whalePublicKey string, datanodeClient *datanode.DataNode, numberOfTransferToBeDone int, logger *zap.Logger) error {
	networkParameters, err := datanodeClient.GetAllNetworkParameters()
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}

	updateParams := map[string]string{
		netparams.TransferMaxCommandsPerEpoch: fmt.Sprintf("%d", numberOfTransferToBeDone),
	}

	updateCount, err := governance.ProposeAndVoteOnNetworkParameters(ctx, updateParams, whaleWallet, whalePublicKey, networkParameters, datanodeClient, logger)
	if err != nil {
		return fmt.Errorf("failed to propose and vote for network parameter update proposals: %w", err)
	}

	if updateCount == 0 {
		logger.Debug("No network parameter update is required before issuing transfers")
		return nil
	}

	updatedNetworkParameters, err := datanodeClient.GetAllNetworkParameters()
	if err != nil {
		return fmt.Errorf("could not retrieve updated network parameters from datanode: %w", err)
	}

	for name, expectedValue := range updateParams {
		updatedValue := updatedNetworkParameters.Params[name]
		if updatedValue != expectedValue {
			return fmt.Errorf("failed to update network parameter %q, current value: %q, expected value: %q", name, updatedValue, expectedValue)
		}
	}

	return nil
}

func ensureWhaleReceivedFunds(ctx context.Context, whaleClient *api.PartyClient, whaleTopUpsByAsset map[string]*types.Amount) error {
	for assetID, requiredAmount := range whaleTopUpsByAsset {
		requiredAmountAsSubUnit := requiredAmount.AsSubUnit()
		err := tools.RetryRun(15, 6*time.Second, func() error {
			balanceAsSubUnit, err := whaleClient.GeneralAccountBalanceForAsset(ctx, assetID)
			if err != nil {
				return fmt.Errorf("failed to retrieve whale's general account balance for asset %q: %w", assetID, err)
			}

			if balanceAsSubUnit.Cmp(requiredAmountAsSubUnit) < 0 {
				return fmt.Errorf("whale wallet does not have enough funds for asset %q: expected %s, got %s", assetID, requiredAmountAsSubUnit.String(), balanceAsSubUnit.String())
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func transferAssetsFromWhaleToBots(ctx context.Context, datanodeClient *datanode.DataNode, whaleWallet wallet.Wallet, whalePublicKey string, registry map[string]AssetToTopUp, logger *zap.Logger) error {
	result := &multierror.Error{}

	for assetID, entry := range registry {
		for botPartyId, amount := range entry.AmountsByParty {
			err := tools.RetryRun(15, 6*time.Second, func() error {
				request := walletpb.SubmitTransactionRequest{
					Command: &walletpb.SubmitTransactionRequest_Transfer{
						Transfer: &vegacmd.Transfer{
							Reference:       fmt.Sprintf("Transfer from whale to %s", botPartyId),
							FromAccountType: vegapb.AccountType_ACCOUNT_TYPE_GENERAL,
							ToAccountType:   vegapb.AccountType_ACCOUNT_TYPE_GENERAL,
							To:              botPartyId,
							Asset:           assetID,
							Amount:          amount.String(),
							Kind: &vegacmd.Transfer_OneOff{
								OneOff: &vegacmd.OneOffTransfer{
									DeliverOn: 0,
								},
							},
						},
					},
				}

				resp, err := walletpkg.SendTransaction(ctx, whaleWallet, whalePublicKey, &request, datanodeClient)
				if err != nil {
					return fmt.Errorf("failed to send the submit transfer transaction for bot %q: %w", botPartyId, err)
				}

				if !resp.Success {
					return fmt.Errorf("the transfer transaction failed for bot %q: %w", botPartyId, err)
				}

				logger.Info("Assets have been sent from whale to the trading bot",
					zap.String("asset", entry.Symbol),
					zap.String("amount", amount.String()),
					zap.String("bot", botPartyId),
					zap.String("transaction", resp.TxHash),
				)

				return nil
			})
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("topping up bot %q failed: %w", botPartyId, err))

				logger.Error("Topping up bot failed",
					zap.String("bot", botPartyId),
					zap.Error(err),
				)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	return result.ErrorOrNil()
}
