package bots

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/vegaprotocol/devopstools/bots"
	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/generation"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/networktools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"code.vegaprotocol.io/vega/core/netparams"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"
	"code.vegaprotocol.io/vega/wallet/wallet"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type ReferralArgs struct {
	*Args
	Setup                 bool
	Assets                []string
	NumberOfSets          uint32
	NumberOfMembersPerSet uint32
}

var referralArgs ReferralArgs

// referralCmd represents the referral command
var referralCmd = &cobra.Command{
	Use:   "referral",
	Short: "Make the selected bots participate to the referral program",
	Long:  "Make the selected bots participate to the referral program",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runReferral(referralArgs); err != nil {
			referralArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	referralArgs.Args = &args

	Cmd.AddCommand(referralCmd)
	referralCmd.PersistentFlags().BoolVar(
		&referralArgs.Setup,
		"setup",
		false,
		"Setup bots in referral program. By default, it is dry run",
	)
	referralCmd.PersistentFlags().StringSliceVar(
		&referralArgs.Assets,
		"assets-symbols",
		[]string{},
		"Used to only select the bots that operate on the specified assets. Other bots won't be added to the referral sets. If left empty, all bots are included.",
	)
	referralCmd.PersistentFlags().Uint32Var(
		&referralArgs.NumberOfSets,
		"max-number-of-teams",
		10,
		"Maximum number of referral sets. However, we create one referral set per market maker",
	)
	referralCmd.PersistentFlags().Uint32Var(
		&referralArgs.NumberOfMembersPerSet,
		"max-number-of-team-members",
		15,
		"Maximum number of referees per referral set. It is limited by number of traders on the research bots",
	)
}

type Party struct {
	Name      string
	Wallet    wallet.Wallet
	PublicKey string
}

type ReferralSet struct {
	Leader   Party
	Referees []Party
}

func NewReferralSet(referrer Party) ReferralSet {
	return ReferralSet{
		Leader:   referrer,
		Referees: []Party{},
	}
}

func runReferral(args ReferralArgs) error {
	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	if !args.Setup {
		logger.Warn("DRY RUN - use --setup flag to run for real")
	}

	cfg, err := config.Load(args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Info("Network file loaded", zap.String("name", cfg.Name))

	whaleWallet, err := vega.LoadWallet(cfg.Network.Wallets.VegaTokenWhale.Name, cfg.Network.Wallets.VegaTokenWhale.RecoveryPhrase)
	if err != nil {
		return fmt.Errorf("could not initialized whale wallet: %w", err)
	}

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Info("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Info("Looking for healthy gRPC endpoints...")
	healthyEndpoints := networktools.FilterHealthyGRPCEndpoints(endpoints)
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

	logger.Sugar().Infof("Fetching traders from the %s", cfg.Bots.Research.RESTURL)
	researchBots, err := bots.RetrieveResearchBots(ctx, cfg.Bots.Research.RESTURL, cfg.Bots.Research.APIKey, logger.Named("research-bots"))
	if err != nil {
		return fmt.Errorf("failed to retrieve research bots: %w", err)
	}
	logger.Info("Research bots found", zap.Strings("traders", maps.Keys(researchBots)))

	if err := prepareNetworkParameters(ctx, whaleWallet, datanodeClient, !args.Setup, logger.Named("pepare network parameters")); err != nil {
		return fmt.Errorf("failed prepare network parameters for referral")
	}

	logger.Info("Retrieving markets for filtered assets...", zap.Strings("assets", args.Assets))
	wantedMarketsIds, err := findMarketsForAssets(ctx, datanodeClient, args.Assets)
	if err != nil {
		return fmt.Errorf("failed to find markets for wanted assets")
	}
	logger.Info("Markets retrieved", zap.Strings("market-ids", wantedMarketsIds))

	logger.Info("Getting referral sets")
	referralSets, err := datanodeClient.ListReferralSets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list referral sets: %w", err)
	}
	logger.Info("Referral sets got")

	logger.Info("Building referral sets topology...")
	newReferralSets, err := buildReferralSetsTopology(referralSets, researchBots, int(args.NumberOfSets), int(args.NumberOfMembersPerSet), wantedMarketsIds)
	if err != nil {
		return fmt.Errorf("could not build referral sets topology: %w", err)
	}
	logger.Info("Referral sets topology built")

	for setNo, referralSet := range newReferralSets {
		var refereesPublicKeys []string

		for _, referee := range referralSet.Referees {
			refereesPublicKeys = append(refereesPublicKeys, referee.PublicKey)
		}

		logger.Info(fmt.Sprintf("Referral set topology #%d", setNo),
			zap.String("referrer", referralSet.Leader.PublicKey),
			zap.Strings("referees", refereesPublicKeys),
			zap.Int("total-referees", len(referralSet.Referees)),
		)
	}

	logger.Info("Retrieving network parameters...")
	networkParams, err := datanodeClient.GetAllNetworkParameters()
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}
	logger.Info("Network parameters retrieved")

	primaryEthConfig, err := networkParams.PrimaryEthereumConfig()
	if err != nil {
		return fmt.Errorf("could not get primary ethereum configuration from network paramters: %w", err)
	}

	primaryChainClient, err := ethereum.NewPrimaryChainClient(ctx, cfg.Bridges.Primary, primaryEthConfig, logger.Named("primary-chain-client"))
	if err != nil {
		return fmt.Errorf("could not initialize primary ethereum chain client: %w", err)
	}

	logger.Info("Ensuring enough stake for referrers...")
	if err := ensureReferrersHaveEnoughStake(ctx, newReferralSets, datanodeClient, primaryChainClient, !args.Setup, logger); err != nil {
		return fmt.Errorf("failed to ensure enough stake for referrers: %w", err)
	}
	logger.Info("Referrers have enough stake")

	logger.Info("Creating new referral sets...")
	if err := createReferralSets(ctx, newReferralSets, datanodeClient, !args.Setup, logger); err != nil {
		return fmt.Errorf("failed to create referral sets: %w", err)
	}
	logger.Info("Referral sets created")

	if args.Setup {
		logger.Info("Waiting for referral set to be created")
		if err := waitForReferralSets(ctx, newReferralSets, datanodeClient, logger); err != nil {
			return fmt.Errorf("failed to wait for referral sets: %w", err)
		}
		logger.Info("All referral sets are ready")
	}

	logger.Info("Referees are joining the referral sets...")
	if err := joinReferees(ctx, newReferralSets, datanodeClient, !args.Setup, logger); err != nil {
		return fmt.Errorf("referees failed to join the referral sets: %w", err)
	}
	logger.Info("Referees joined the referral sets")

	logger.Info("Referral set created and joined by research bots successfully")

	return nil
}

func prepareNetworkParameters(ctx context.Context, whaleWallet wallet.Wallet, datanodeClient *datanode.DataNode, dryRun bool, logger *zap.Logger) error {
	_ = ctx

	networkParameters, err := datanodeClient.GetAllNetworkParameters()
	if err != nil {
		return fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}

	updateParams := map[string]string{
		netparams.SpamProtectionApplyReferralMinFunds: "0",
	}

	if dryRun {
		logger.Info("DRY RUN: Not updating network parameters")
		return nil
	}

	pubKey := vega.MustFirstKey(whaleWallet)

	updateCount, err := governance.ProposeAndVoteOnNetworkParameters(ctx, updateParams, whaleWallet, pubKey, networkParameters, datanodeClient, logger)
	if err != nil {
		return fmt.Errorf("failed to propose and vote for network parameter update proposals: %w", err)
	}

	if updateCount == 0 {
		logger.Debug("No network parameter update is required before issuing transfers")
		return nil
	}

	return nil
}

func waitForReferralSets(ctx context.Context, referralSets []ReferralSet, dataNodeClient vegaapi.DataNodeClient, logger *zap.Logger) error {
	wantedReferrerPubKeys := make([]string, len(referralSets))
	for _, referralSet := range referralSets {
		wantedReferrerPubKeys = append(wantedReferrerPubKeys, referralSet.Leader.PublicKey)
	}

	ticker := time.NewTicker(10 * time.Second)
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer func() {
		cancel()
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			referralSets, err := dataNodeClient.ListReferralSets(timeoutCtx)
			if err != nil {
				logger.Debug("Cannot retrieve current referral sets: %s", zap.Error(err))
				continue
			}

			if len(referralSets) < 1 {
				logger.Debug("No referral sets found on the network, yet")
				continue
			}

			for referrer, referralSet := range referralSets {
				if !slices.Contains(wantedReferrerPubKeys, referrer) {
					continue // referral set already confirmed or external referral set
				}

				logger.Info("Found referral set",
					zap.String("referral-set-id", referralSet.Id),
					zap.String("referrer", referrer),
				)
				wantedReferrerPubKeysNew := slices.DeleteFunc(wantedReferrerPubKeys, func(item string) bool {
					return item == referrer
				})

				wantedReferrerPubKeys = []string{}
				for _, pubKey := range wantedReferrerPubKeysNew {
					if pubKey == "" {
						continue
					}

					wantedReferrerPubKeys = append(wantedReferrerPubKeys, pubKey)
				}
			}

			if len(wantedReferrerPubKeys) == 0 {
				logger.Debug("All referral sets have been created")
				return nil
			}

			logger.Debug("Still waiting for all referral sets to be created",
				zap.Strings("waiting-for-referrers", wantedReferrerPubKeys),
			)
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout exceeded")
		}
	}
}

// buildReferralSetsTopology should generate deterministic referral set based on the response from the /traders endpoint
func buildReferralSetsTopology(existingReferralSets map[string]*v2.ReferralSet, traders bots.ResearchBots, numberOfSets int, numberOfReferees int, includedMarkets []string) ([]ReferralSet, error) {
	_ = existingReferralSets
	if numberOfSets < 1 {
		return nil, fmt.Errorf("you must create at least one referral set")
	}

	if numberOfReferees < 1 {
		return nil, fmt.Errorf("you must add at least one referee to each referral set")
	}

	// traderID contains the market id in the name it is trading on
	filteredTraders := bots.ResearchBots{}
	for traderID, trader := range traders {
		for _, marketId := range includedMarkets {
			if strings.Contains(trader.Name, marketId) {
				filteredTraders[traderID] = trader
				break
			}
		}
	}

	// TODO: Fetch teams which were already created earlier.
	// teamOwners := []string{}
	// for _, referralSet := range existingReferralSets {
	// 	teamOwners = append(teamOwners, referralSet.Referrer)
	// }

	var referralSets []ReferralSet
	// Create teams with leaders
	for _, trader := range filteredTraders {
		// Party is already owner of the team
		// if slices.Contains(teamOwners, trader.PubKey) {
		// 	continue
		// }

		if trader.IsMarketMaker() {
			w, err := trader.GetWallet()
			if err != nil {
				return nil, fmt.Errorf("failed to get wallet for %s trader when creating the referral set: %w", trader.PubKey, err)
			}
			referralSets = append(referralSets, NewReferralSet(Party{
				Name:   trader.Name,
				Wallet: w,
			}))
		}

		if len(referralSets) >= numberOfSets {
			break // enough referral sets for now
		}
	}

	if len(referralSets) < numberOfSets {
		return nil, fmt.Errorf("not enough traders to create %d referralSets: there should be at least 1 market maker per referral set", numberOfSets)
	}

	// add members to the referralSets
	var potentialMembers []string
	for traderId, trader := range filteredTraders {
		if !trader.IsMarketMaker() {
			potentialMembers = append(potentialMembers, traderId)
		}
	}

	// not enough candidates, add add as many as we can to first referral sets
	if len(potentialMembers) < numberOfReferees*numberOfSets {
		referralSetIndex := 0
		for _, candidateId := range potentialMembers {
			if len(referralSets[referralSetIndex].Referees) >= numberOfReferees {
				referralSetIndex += 1
			}

			if len(referralSets) <= referralSetIndex {
				break // no more referral sets
			}

			trader := filteredTraders[candidateId]
			w, err := trader.GetWallet()
			if err != nil {
				return nil, fmt.Errorf("failed to get wallet for %s trader when assigning members to the referral sets: %w", trader.PubKey, err)
			}
			referralSets[referralSetIndex].Referees = append(referralSets[referralSetIndex].Referees, Party{
				Name:   trader.Name,
				Wallet: w,
			})
		}

		return referralSets, nil
	}

	// mix bots in referral sets
	for index, candidateId := range potentialMembers {
		referralSetIndex := index % len(referralSets)

		if len(referralSets[referralSetIndex].Referees) >= numberOfReferees {
			continue
		}

		trader := filteredTraders[candidateId]
		w, err := trader.GetWallet()
		if err != nil {
			return nil, fmt.Errorf("failed to get wallet for %s trader when assigning members to the referral sets: %w", trader.PubKey, err)
		}
		referralSets[referralSetIndex].Referees = append(referralSets[referralSetIndex].Referees, Party{
			Name:   trader.Name,
			Wallet: w,
		})
	}

	return referralSets, nil
}

func findMarketsForAssets(ctx context.Context, dataNodeClient vegaapi.DataNodeClient, assetsSymbols []string) ([]string, error) {
	allMarkets, err := dataNodeClient.GetAllMarketsWithState(ctx, datanode.ActiveMarkets)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve markets from datanode: %w", err)
	}

	allAssets, err := dataNodeClient.ListAssets(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve assets from datanode: %w", err)
	}

	var result []string
	for _, market := range allMarkets {
		settlementAsset := ""
		quoteAsset := ""

		if market.GetTradableInstrument() != nil && market.GetTradableInstrument().GetInstrument() != nil {
			instrument := market.GetTradableInstrument().GetInstrument()

			if instrument.GetFuture() != nil {
				settlementAsset = instrument.GetFuture().SettlementAsset
			} else if instrument.GetPerpetual() != nil {
				settlementAsset = instrument.GetPerpetual().SettlementAsset
			} else if instrument.GetSpot() != nil {
				settlementAsset = instrument.GetSpot().BaseAsset
				quoteAsset = instrument.GetSpot().QuoteAsset
			}
		}

		if settlementAsset == "" {
			continue
		}

		assetDetails, assetFound := allAssets[settlementAsset]
		if !assetFound {
			return nil, fmt.Errorf("settlement asset %q for market %q not found in asset lists", settlementAsset, market.Id)
		}
		quoteAssetSymbol := "UNKNOWN"
		if quoteAsset != "" {
			quoteAssetDetails, assetFound := allAssets[quoteAsset]
			if !assetFound {
				return nil, fmt.Errorf("quote settlement asset %q for market %q not found in asset lists", settlementAsset, market.Id)
			}

			quoteAssetSymbol = quoteAssetDetails.Symbol
		}

		if len(assetsSymbols) > 0 &&
			!slices.Contains(assetsSymbols, assetDetails.Symbol) &&
			!slices.Contains(assetsSymbols, quoteAssetSymbol) {
			continue
		}

		result = append(result, market.Id)
	}

	return result, nil
}

func joinReferees(ctx context.Context, referralSets []ReferralSet, dataNodeClient vegaapi.DataNodeClient, dryRun bool, logger *zap.Logger) error {
	logger.Info("Join Referral Sets (teams)")
	currentReferralSets, err := dataNodeClient.ListReferralSets(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve current referral sets, %w", err)
	}

	referralSetReferees, err := dataNodeClient.GetReferralSetReferees(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve current referees: %w", err)
	}

	for _, set := range referralSets {
		referralSet, isReferrer := currentReferralSets[set.Leader.PublicKey]
		if !isReferrer {
			return fmt.Errorf("the referrer %s is not in the referral set", set.Leader.PublicKey)
		}

		for _, member := range set.Referees {
			refereeKey := member.PublicKey

			if _, ok := referralSetReferees[refereeKey]; ok {
				logger.Info("Party already belong to a referral set",
					zap.String("party-id", refereeKey),
					zap.String("referral-set-id", referralSetReferees[refereeKey].ReferralSetId),
				)
				continue
			}

			if dryRun {
				logger.Info("DRY RUN - skip referee from joining referral set", zap.String("party-id", refereeKey))
				continue
			}

			txReq := walletpb.SubmitTransactionRequest{
				PubKey: refereeKey,
				Command: &walletpb.SubmitTransactionRequest_ApplyReferralCode{
					ApplyReferralCode: &commandspb.ApplyReferralCode{
						Id: referralSet.Id,
					},
				},
			}

			if _, err := walletpkg.SendTransaction(ctx, member.Wallet, member.PublicKey, &txReq, dataNodeClient); err != nil {
				return fmt.Errorf("transaction to join referral referral set failed %q: %w", referralSet.Id, err)
			}

			logger.Debug("Referees joined referral set", zap.String("party-id", refereeKey))
		}
	}

	return nil
}

func ensureReferrersHaveEnoughStake(ctx context.Context, newReferralSets []ReferralSet, datanodeClient *datanode.DataNode, chainClient *ethereum.ChainClient, dryRun bool, logger *zap.Logger) error {
	minStake := types.NewAmount(18)

	logger.Debug("Retrieving current referral program...")
	program, err := datanodeClient.GetCurrentReferralProgram(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve current referral program from datanode: %w", err)
	}
	logger.Debug("Current referral program found",
		zap.String("id", program.Id),
		zap.Uint64("id", program.Version),
	)

	tiersMinStakeAmounts := make([]*types.Amount, len(program.StakingTiers))
	for i, stakingTier := range program.StakingTiers {
		stakeAmountAsMainUnit, ok := new(big.Float).SetString(stakingTier.MinimumStakedTokens)
		if !ok {
			return fmt.Errorf("failed to convert %q to big.Float", stakingTier.MinimumStakedTokens)
		}
		tiersMinStakeAmounts[i] = types.NewAmountFromMainUnit(stakeAmountAsMainUnit, 18)

		if minStake.Cmp(tiersMinStakeAmounts[i]) > 0 {
			minStake = tiersMinStakeAmounts[i]
		}
	}

	missingStakeByPubKey := map[string]*types.Amount{}

	for _, referralSet := range newReferralSets {
		logger.Debug("Retrieving current stake for party...",
			zap.String("party-id", referralSet.Leader.PublicKey),
		)
		currentStakeAsSubUnit, err := datanodeClient.GetPartyTotalStake(referralSet.Leader.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to retrieve current stake for party %s: %w", referralSet.Leader.PublicKey, err)
		}
		currentStake := types.NewAmountFromSubUnit(currentStakeAsSubUnit, 18)
		logger.Debug("Current stake for party found",
			zap.String("party-id", referralSet.Leader.PublicKey),
			zap.String("current-stake", currentStake.String()),
		)

		if currentStake.Cmp(minStake) < 0 {
			expectedStake := tiersMinStakeAmounts[rand.Intn(len(tiersMinStakeAmounts))]
			logger.Debug("Party needs more stake",
				zap.String("party-id", referralSet.Leader.PublicKey),
				zap.String("program-minimum-stake", minStake.String()),
				zap.String("current-stake", currentStake.String()),
				zap.String("expected-stake", expectedStake.String()))

			missingStake := expectedStake.Copy()
			missingStake.Sub(currentStake.AsMainUnit())
			missingStakeByPubKey[referralSet.Leader.PublicKey] = missingStake
		} else {
			logger.Debug("Party does not need more stake",
				zap.String("party-id", referralSet.Leader.PublicKey),
				zap.String("program-minimum-stake", minStake.String()),
				zap.String("current-stake", currentStake.String()),
			)
		}
	}

	if dryRun {
		logger.Warn("DRY RUN - not running stake")
		return nil
	}

	logger.Debug("Staking Vega token to parties", zap.Strings("parties", maps.Keys(missingStakeByPubKey)))
	if err := chainClient.StakeVegaTokenFromMinter(ctx, missingStakeByPubKey); err != nil {
		return fmt.Errorf("failed to stake Vega token from minter wallet: %w", err)
	}
	logger.Debug("Staking Vega token successful", zap.Strings("parties", maps.Keys(missingStakeByPubKey)))

	return nil
}

func createReferralSets(ctx context.Context, newReferralSets []ReferralSet, dataNodeClient vegaapi.DataNodeClient, dryRun bool, logger *zap.Logger) error {
	currentReferralSets, err := dataNodeClient.ListReferralSets(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieved existing referral sets: %w", err)
	}

	for _, currentReferralSet := range newReferralSets {
		referrerKey := currentReferralSet.Leader.PublicKey
		if referralSet, ok := currentReferralSets[referrerKey]; ok {
			logger.Debug("Party is already is a referrer",
				zap.String("pub key", referrerKey),
				zap.String("team id", referralSet.Id))
			continue
		}

		if dryRun {
			logger.Info("DRY RUN - skip creation of referral set for party", zap.String("party-id", referrerKey))
		} else {
			if err := createReferralSet(ctx, currentReferralSet.Leader, dataNodeClient); err != nil {
				return fmt.Errorf("failed to create referral set %s, %w", referrerKey, err)
			}
		}
	}

	return nil
}

func createReferralSet(ctx context.Context, creator Party, dataNodeClient vegaapi.DataNodeClient) error {
	teamName, err := generation.GenerateName()
	if err != nil {
		return fmt.Errorf("could not generate a name for referral set: %w", err)
	}
	teamName = fmt.Sprintf("Bots Team: %s", teamName)
	teamURL, err := generation.GenerateRandomWikiURL()
	if err != nil {
		return fmt.Errorf("could not generate a wiki URL for referral set: %w", err)
	}
	teamAvatar, err := generation.GenerateAvatarURL()
	if err != nil {
		return fmt.Errorf("could not generate an avatar URL for referral set: %w", err)
	}

	request := walletpb.SubmitTransactionRequest{
		Command: &walletpb.SubmitTransactionRequest_CreateReferralSet{
			CreateReferralSet: &commandspb.CreateReferralSet{
				IsTeam: true,
				Team: &commandspb.CreateReferralSet_Team{
					Name:      teamName,
					TeamUrl:   &teamURL,
					AvatarUrl: &teamAvatar,
				},
			},
		},
	}

	if _, err := walletpkg.SendTransaction(ctx, creator.Wallet, creator.PublicKey, &request, dataNodeClient); err != nil {
		return fmt.Errorf("transaction to create a referral set failed: %w", err)
	}
	return nil
}
