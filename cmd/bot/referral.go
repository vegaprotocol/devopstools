package bot

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/bots"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/generate"
	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
)

const waitTimeout = 5 * time.Minute

type ReferralArgs struct {
	*BotArgs
	SetupBotsInReferralProgram bool
	Assets                     []string
	NumberOfTeams              uint32
	NumberOfMembersPerTeam     uint32
}

var referralArgs ReferralArgs

// referralCmd represents the referral command
var referralCmd = &cobra.Command{
	Use:   "referral",
	Short: "manage bots in referral program",
	Long:  `manage bots in referral program`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runReferral(referralArgs); err != nil {
			referralArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	referralArgs.BotArgs = &botArgs

	BotCmd.AddCommand(referralCmd)
	referralCmd.PersistentFlags().BoolVar(
		&referralArgs.SetupBotsInReferralProgram,
		"setup",
		false,
		"Setup bots in referral program. By default it is dry run",
	)
	referralCmd.PersistentFlags().StringSliceVar(
		&referralArgs.Assets,
		"assets-symbols",
		[]string{},
		"Assets, bots operates. Script will ignore all other bots and won't add them to the teams. Empty == all markets included",
	)
	referralCmd.PersistentFlags().Uint32Var(
		&referralArgs.NumberOfTeams,
		"max-number-of-teams",
		10,
		"Maximum number of teams. However we create one team per one market maker",
	)
	referralCmd.PersistentFlags().Uint32Var(
		&referralArgs.NumberOfMembersPerTeam,
		"max-number-of-team-members",
		15,
		"Maximum number of team members. It is limited by number of traders on the research bots",
	)
}

type TeamMember struct {
	Name   string
	Wallet *wallet.VegaWallet
}

type Team struct {
	Leader TeamMember

	Members []TeamMember
}

func NewTeam(leader TeamMember) Team {
	return Team{
		Leader:  leader,
		Members: []TeamMember{},
	}
}

func runReferral(args ReferralArgs) error {
	logger := args.Logger
	start := time.Now()
	logger.Info("Start referral", zap.Time("start", start))
	if !args.SetupBotsInReferralProgram {
		logger.Info("DRY RUN - use --setup flag to run for real")
	}
	logger.Info("Connecting to nework")
	network, err := args.ConnectToVegaNetwork(args.VegaNetworkName)
	if err != nil {
		return err
	}
	defer network.Disconnect()
	logger.Info("Connected to network", zap.String("network", args.VegaNetworkName), zap.Duration("since start", time.Since(start)))

	logger.Info("Getting bots")
	botsAPIToken := args.BotsAPIToken
	if len(botsAPIToken) == 0 {
		botsAPIToken = network.BotsApiToken
	}

	traders, err := bots.GetResearchBots(args.VegaNetworkName, botsAPIToken)
	if err != nil {
		return err
	}
	logger.Info("Got bots", zap.Int("count", len(traders)), zap.Duration("since start", time.Since(start)))

	wantedMarketsIds, err := findMarketMarketsForAssets(network.DataNodeClient, args.Assets)
	if err != nil {
		return fmt.Errorf("failed to find markets for wanted assets")
	}

	teams, err := prepareTeams(traders, int(args.NumberOfTeams), int(args.NumberOfMembersPerTeam), wantedMarketsIds)
	if err != nil {
		return fmt.Errorf("failed to prepare teams: %w", err)
	}

	logger.Info("Teams design created.")
	for teamNo, team := range teams {
		teamMembersPublicKeys := []string{}

		for _, member := range team.Members {
			teamMembersPublicKeys = append(teamMembersPublicKeys, member.Wallet.PublicKey)
		}

		logger.Sugar().Infof(
			"Team #%d has leader %s and %d members: \n\t - %s",
			teamNo,
			team.Leader.Wallet.PublicKey,
			len(team.Members),
			strings.Join(teamMembersPublicKeys, "\n\t - "),
		)
	}

	logger.Info("Staking to teams leaders")
	if err := stakeToTeamLeaders(logger, teams, network, !args.SetupBotsInReferralProgram); err != nil {
		return fmt.Errorf("failed to stake to the team leaders: %w", err)
	}
	logger.Info("Staked tokens to the team leaders")

	logger.Info("Creating referral sets")
	if err := createReferralSets(logger, teams, network.DataNodeClient, !args.SetupBotsInReferralProgram); err != nil {
		return fmt.Errorf("failed to create referral sets: %w", err)
	}
	logger.Info("Teams created")

	logger.Info("Waiting for referral set to be created")
	if err := waitForReferralSets(logger, teams, network.DataNodeClient, !args.SetupBotsInReferralProgram, waitTimeout); err != nil {
		return fmt.Errorf("failed to wait for referral sets: %w")
	}
	logger.Info("All referral sets are ready")

	logger.Info("Joining members to the cluster")
	if err := joinMemberTeams(logger, teams, network.DataNodeClient, !args.SetupBotsInReferralProgram); err != nil {
		return fmt.Errorf("failed to join members to the teams: %w", err)
	}
	logger.Info("FINISHED")

	return nil
}

func waitForReferralSets(
	logger *zap.Logger,
	teams []Team,
	dataNodeClient vegaapi.DataNodeClient,
	dryRun bool,
	timeout time.Duration,
) error {
	wantedTeamsLeaderPubKeys := make([]string, len(teams))
	for _, team := range teams {
		wantedTeamsLeaderPubKeys = append(wantedTeamsLeaderPubKeys, team.Leader.Wallet.PublicKey)
	}

	ticker := time.NewTicker(10 * time.Second)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ticker.C:
			referralSets, err := dataNodeClient.GetReferralSets()
			if err != nil {
				logger.Sugar().Errorf("Cannot get referral sets: %s", err.Error())
				continue
			}
			if len(referralSets) < 1 {
				logger.Info("No referral sets found yet on the network")
				continue
			}

			for referrer, referralSet := range referralSets {
				if !slices.Contains(wantedTeamsLeaderPubKeys, referrer) {
					continue // team already confirmed or external team
				}

				logger.Sugar().Infof("Found team with id %s for leader %s", referralSet.Id, referrer)
				wantedTeamsLeaderPubKeys = slices.DeleteFunc(wantedTeamsLeaderPubKeys, func(item string) bool {
					return item == referrer
				})
			}

			stillWaitingFor := []string{}
			for _, partyId := range wantedTeamsLeaderPubKeys {
				if len(partyId) > 0 {
					stillWaitingFor = append(stillWaitingFor, partyId)
				}
			}

			if len(stillWaitingFor) < 1 {
				logger.Info("All required teams exist")
				return nil
			}
			logger.Sugar().Infof("Still waiting for referral sets for the %#v leaders", stillWaitingFor)
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout exceeded")
		}
	}

	return nil
}

// prepareTeams should generate deterministic teams based on the response from the /traders endpoint
func prepareTeams(traders bots.ResearchBots, numberOfTeams int, numberOfMembers int, includedMarkets []string) ([]Team, error) {
	if numberOfTeams < 1 {
		return nil, fmt.Errorf("you must create at least one team")
	}

	if numberOfMembers < 1 {
		return nil, fmt.Errorf("you must add at least one member to each of the team")
	}

	// traderId contains the market id in the name it is trading on
	filteredTraders := bots.ResearchBots{}
	for traderId, trader := range traders {
		for _, marketId := range includedMarkets {
			if strings.Contains(traderId, marketId) {
				filteredTraders[traderId] = trader
				break
			}
		}
	}

	teams := []Team{}
	// Create reams with leaders
	for _, trader := range filteredTraders {
		if trader.WalletData.Index == bots.MarketMakerWalletIndex {
			wallet, err := trader.GetWallet()
			if err != nil {
				return nil, fmt.Errorf("failed to create wallet for the %s research-bot when creating teams: %w", trader.Name, err)
			}
			teams = append(teams, NewTeam(TeamMember{
				Name:   trader.Name,
				Wallet: wallet,
			}))
		}

		if len(teams) >= numberOfTeams {
			break // enough teams for now
		}
	}

	if len(teams) < numberOfTeams {
		return nil, fmt.Errorf("not enough traders to create %d teams: there should be at least 1 market maker per team", numberOfTeams)
	}

	// add members to the teams
	potentialMembers := []string{}
	for traderId, trader := range filteredTraders {
		if trader.WalletData.Index != bots.MarketMakerWalletIndex {
			potentialMembers = append(potentialMembers, traderId)
		}
	}

	// not enough candidates, add add as many as we can to first teams
	if len(potentialMembers) < numberOfMembers*numberOfTeams {
		teamIndex := 0
		for _, candidateId := range potentialMembers {
			if len(teams[teamIndex].Members) >= numberOfMembers {
				teamIndex += 1
			}

			if len(teams) <= teamIndex {
				break // no more teams
			}

			trader := filteredTraders[candidateId]
			wallet, err := trader.GetWallet()
			if err != nil {
				return nil, fmt.Errorf("failed to get wallet for %s trader when assigning members to the teams: %w", trader.PubKey, err)
			}
			teams[teamIndex].Members = append(teams[teamIndex].Members, TeamMember{
				Name:   trader.Name,
				Wallet: wallet,
			})
		}

		return teams, nil
	}

	// mix bots in teams
	for index, candidateId := range potentialMembers {
		teamIndex := index % len(teams)

		if len(teams[teamIndex].Members) >= numberOfMembers {
			continue
		}

		trader := filteredTraders[candidateId]
		wallet, err := trader.GetWallet()
		if err != nil {
			return nil, fmt.Errorf("failed to get wallet for %s trader when assigning members to the teams: %w", trader.PubKey, err)
		}
		teams[teamIndex].Members = append(teams[teamIndex].Members, TeamMember{
			Name:   trader.Name,
			Wallet: wallet,
		})
	}

	return teams, nil
}

func findMarketMarketsForAssets(dataNodeClient vegaapi.DataNodeClient, assetsSymbols []string) ([]string, error) {
	allMarkets, err := dataNodeClient.GetAllMarkets()
	if err != nil {
		return nil, fmt.Errorf("failed to get all markets: %w", err)
	}
	allAssets, err := dataNodeClient.GetAssets()
	if err != nil {
		return nil, fmt.Errorf("failed to get all assets: %w", err)
	}

	result := []string{}
	for _, market := range allMarkets {
		settlementAsset := ""

		if market.GetTradableInstrument() != nil && market.GetTradableInstrument().GetInstrument() != nil {
			instrument := market.GetTradableInstrument().GetInstrument()

			if instrument.GetFuture() != nil {
				settlementAsset = instrument.GetFuture().SettlementAsset
			} else if instrument.GetPerpetual() != nil {
				settlementAsset = instrument.GetPerpetual().SettlementAsset
			}
		}

		if settlementAsset == "" {
			continue
		}

		assetDetails, assetFound := allAssets[settlementAsset]
		if !assetFound {
			return nil, fmt.Errorf("failed to find the %s asset on the network which is used for market %s", settlementAsset, market.Id)
		}

		if len(assetsSymbols) > 0 && !slices.Contains(assetsSymbols, assetDetails.Symbol) {
			continue // we are not interested in this market
		}

		result = append(result, market.Id)
	}

	return result, nil
}

func joinMemberTeams(
	logger *zap.Logger,
	teams []Team,
	dataNodeClient vegaapi.DataNodeClient,
	dryRun bool,
) error {

	logger.Info("Join Referral Sets (teams)")
	referralSets, err := dataNodeClient.GetReferralSets()
	if err != nil {
		return fmt.Errorf("failed to get referral sets, %w", err)
	}

	referralSetReferees, err := dataNodeClient.GetReferralSetReferees()
	if err != nil {
		return fmt.Errorf("failed to get referral set referees: %w", err)
	}

	for _, team := range teams {
		start := time.Now()
		referralSet, isLeaderInTheTeam := referralSets[team.Leader.Wallet.PublicKey]
		if !isLeaderInTheTeam {
			return fmt.Errorf("the team leader %s is not in the team", team.Leader.Wallet.PublicKey)
		}

		for _, member := range team.Members {
			if referralSet, ok := referralSetReferees[member.Wallet.PublicKey]; ok {
				logger.Debug("Party already belong to a team", zap.String("pub key", member.Wallet.PublicKey),
					zap.String("team", referralSet.ReferralSetId), zap.String("team lead", referralSet.Referee))
				continue
			}

			if dryRun {
				logger.Info("DRY RUN - skip joining a team by", zap.String("pub key", member.Wallet.PublicKey))
				continue
			}

			walletTxReq := walletpb.SubmitTransactionRequest{
				PubKey: member.Wallet.PublicKey,
				Command: &walletpb.SubmitTransactionRequest_ApplyReferralCode{
					ApplyReferralCode: &commandspb.ApplyReferralCode{
						Id: referralSet.Id,
					},
				},
			}
			if err := governance.SubmitTx(fmt.Sprintf("join referral team %s", referralSet.Id),
				dataNodeClient, member.Wallet, logger, &walletTxReq); err != nil {
				return fmt.Errorf("failedy to apply referral code, %w", err)
			}

			logger.Info("Joined Referral Sets (teams)", zap.Duration("since start", time.Since(start)))
		}
	}

	return nil
}

func stakeToTeamLeaders(
	logger *zap.Logger,
	teams []Team,
	network *veganetwork.VegaNetwork,
	dryRun bool,
) error {
	minStake := big.NewInt(0)
	dataNodeClient := network.DataNodeClient

	// get Referrals Tiers
	program, err := dataNodeClient.GetCurrentReferralProgram()
	if err != nil {
		return fmt.Errorf("failed to create referral sets, failed to get referral program: %w", err)
	}

	tiersMinStakeAmounts := make([]*big.Int, len(program.StakingTiers))
	for i, stakingTier := range program.StakingTiers {
		if stakeAmount, ok := new(big.Int).SetString(stakingTier.MinimumStakedTokens, 0); ok {
			tiersMinStakeAmounts[i] = ethutils.VegaTokenFromFullTokens(new(big.Float).SetInt(stakeAmount))
		} else {
			return fmt.Errorf("failed to convert %s to big.Int", stakingTier.MinimumStakedTokens)
		}
		if minStake.Cmp(tiersMinStakeAmounts[i]) > 0 {
			minStake = tiersMinStakeAmounts[i]
		}
	}

	logger.Debug("Got Referral Program Staking Tiers", zap.Int("tiers count", len(tiersMinStakeAmounts)),
		zap.String("minStake", minStake.String()), zap.Any("tiers", tiersMinStakeAmounts))

	stakeByPubKey := map[string]*big.Int{}

	for _, team := range teams {
		currentStake, err := dataNodeClient.GetPartyTotalStake(team.Leader.Wallet.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to create referral sets, failed to get stake for %: %w", team.Leader.Wallet.PublicKey, err)
		}
		if currentStake.Cmp(minStake) < 0 {
			// rand stake
			rndIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(tiersMinStakeAmounts))))
			if err != nil {
				return err
			}
			expectedStake := tiersMinStakeAmounts[rndIdx.Int64()]
			stakeByPubKey[team.Leader.Wallet.PublicKey] = expectedStake
			logger.Debug("Need to top up", zap.String("wallet", team.Leader.Wallet.PublicKey),
				zap.String("current stake", currentStake.String()), zap.String("min stake", minStake.String()),
				zap.String("expected stake", expectedStake.String()))
		} else {
			stakeByPubKey[team.Leader.Wallet.PublicKey] = currentStake
			logger.Debug("No need to top up", zap.String("wallet", team.Leader.Wallet.PublicKey),
				zap.String("current stake", currentStake.String()), zap.String("min stake", minStake.String()))
		}
	}

	if dryRun {
		logger.Info("DRY RUN - not running stake\n")
	} else {
		if err := doStake(stakeByPubKey, network, logger); err != nil {
			return err
		}
	}

	return nil
}

func createReferralSets(
	logger *zap.Logger,
	teams []Team,
	dataNodeClient vegaapi.DataNodeClient,
	dryRun bool,
) error {
	referralSets, err := dataNodeClient.GetReferralSets()
	if err != nil {
		return fmt.Errorf("failed to get referral sets, %w", err)
	}
	for _, team := range teams {
		if referralSet, ok := referralSets[team.Leader.Wallet.PublicKey]; ok {
			logger.Debug("party is already team lead", zap.String("pub key", team.Leader.Wallet.PublicKey),
				zap.String("team id", referralSet.Id))
			continue
		}

		// create referral set
		if dryRun {
			logger.Info("DRY RUN - skip creation of referral set for wallet", zap.String("pub key", team.Leader.Wallet.PublicKey))
		} else {
			if err := createReferralSet(team.Leader.Wallet, dataNodeClient, logger); err != nil {
				return fmt.Errorf("failed to create referral set %s, %w", team.Leader.Wallet.PublicKey, err)
			}
		}
	}

	return nil
}

func createReferralSet(
	creatorVegawallet *wallet.VegaWallet,
	dataNodeClient vegaapi.DataNodeClient,
	logger *zap.Logger,
) error {
	errorMsg := fmt.Errorf("failed to create referral set by %s", creatorVegawallet.PrivateKey)
	teamName, err := generate.GenerateName()
	if err != nil {
		return fmt.Errorf("%s, %w", errorMsg, err)
	}
	teamName = fmt.Sprintf("Bots Team: %s", teamName)
	teamURL, err := generate.GenerateRandomWikiURL()
	if err != nil {
		return fmt.Errorf("%s, %w", errorMsg, err)
	}
	teamAvatar, err := generate.GenerateAvatarURL()
	if err != nil {
		return fmt.Errorf("%s, %w", errorMsg, err)
	}

	walletTxReq := walletpb.SubmitTransactionRequest{
		PubKey: creatorVegawallet.PublicKey,
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
	if err := governance.SubmitTx(fmt.Sprintf("create referral team %s", teamName),
		dataNodeClient, creatorVegawallet, logger, &walletTxReq); err != nil {
		return fmt.Errorf("%s, %w", errorMsg, err)
	}
	return nil
}

func doStake(
	expectedStakeForParty map[string]*big.Int,
	network *veganetwork.VegaNetwork,
	logger *zap.Logger,
) error {
	if len(expectedStakeForParty) == 0 {
		logger.Info("No parties need staking")
		return nil
	}
	//
	// Get missing stake
	//
	dataNodeClient := network.DataNodeClient
	totalMissingStake := big.NewInt(0)
	missingStakeByPubKey := map[string]*big.Int{}

	for pubKey, expectedStake := range expectedStakeForParty {
		currentStake, err := dataNodeClient.GetPartyTotalStake(pubKey)
		if err != nil {
			return fmt.Errorf("failed to create referral sets, failed to get stake for %s, %w", pubKey, err)
		}
		if currentStake.Cmp(expectedStake) < 0 {
			missingStake := new(big.Int).Sub(expectedStake, currentStake)
			missingStakeByPubKey[pubKey] = missingStake
			totalMissingStake = totalMissingStake.Add(totalMissingStake, missingStake)

			logger.Debug("Need to top up", zap.String("wallet", pubKey),
				zap.String("current stake", currentStake.String()),
				zap.String("expected stake", expectedStake.String()))
		} else {
			logger.Debug("No need to top up", zap.String("wallet", pubKey),
				zap.String("current stake", currentStake.String()),
				zap.String("expected stake", expectedStake.String()))
		}
	}
	if len(missingStakeByPubKey) == 0 {
		logger.Info("No part needs staking")
		return nil
	}
	//
	// Prepare source Ethereum wallet
	//
	vegaToken := network.SmartContracts.VegaToken
	wallet := network.NetworkMainWallet
	stakingBridge := network.SmartContracts.StakingBridge

	// TODO - check if there any pending transactions

	// check BALANCE
	currentBalance, err := vegaToken.BalanceOf(&bind.CallOpts{}, wallet.Address)
	if err != nil {
		return fmt.Errorf("failed to get balance, %w", err)
	}
	if currentBalance.Cmp(totalMissingStake) < 0 {
		return fmt.Errorf("wallet %s doesn't have enought tokens; current: %s, expected: %s",
			wallet.Address.Hex(), currentBalance.String(), totalMissingStake.String())
	}
	// check ALLOWANCE
	currentAllowance, err := vegaToken.Allowance(&bind.CallOpts{}, wallet.Address, stakingBridge.Address)
	if err != nil {
		return fmt.Errorf("failed to get allowance for staking bridge %s for party %s, %w",
			stakingBridge.Address.Hex(), wallet.Address, err)
	}
	if currentAllowance.Cmp(totalMissingStake) < 0 {
		// increase ALLOWANCE
		increaseAllowanceAmount := ethutils.VegaTokenFromFullTokens(big.NewFloat(1000000))
		if increaseAllowanceAmount.Cmp(totalMissingStake) < 0 {
			return fmt.Errorf("Want to stake more than 1kk tokens total - it is too much")
		}
		logger.Info("increasing allowance", zap.String("amount", increaseAllowanceAmount.String()),
			zap.String("ethWallet", wallet.Address.Hex()),
			zap.String("allowanceBefore", currentAllowance.String()),
			zap.String("tokenAddress", vegaToken.Address.Hex()))
		opts := wallet.GetTransactOpts()
		allowanceTx, err := vegaToken.IncreaseAllowance(opts, stakingBridge.Address, increaseAllowanceAmount)
		if err != nil {
			return fmt.Errorf("failed to increase allowance: %w", err)
		}
		// WAIT
		if err = ethutils.WaitForTransaction(network.EthClient, allowanceTx, time.Minute); err != nil {
			logger.Error("failed to increase allowance", zap.String("ethWallet", wallet.Address.Hex()),
				zap.String("tokenAddress", vegaToken.Address.Hex()), zap.Error(err))
			return fmt.Errorf("transaction failed to increase allowance: %w", err)
		}
		logger.Info("successfully increased allowance", zap.String("ethWallet", wallet.Address.Hex()),
			zap.String("tokenAddress", vegaToken.Address.Hex()))

	} else {
		logger.Info("no need to increase allowance", zap.String("tokenAddress", vegaToken.Address.Hex()),
			zap.String("ethWallet", wallet.Address.Hex()),
			zap.String("currentAllowance", currentAllowance.String()),
			zap.String("requiredAllowance", totalMissingStake.String()))
	}

	//
	// Stake
	//
	var (
		stakeTxs     = make(map[string]*ethtypes.Transaction, len(missingStakeByPubKey))
		failedCount  = 0
		successCount = 0
	)
	logger.Info("Sending stake transactions", zap.Int("count", len(missingStakeByPubKey)))
	for pubKey, missingStake := range missingStakeByPubKey {
		opts := wallet.GetTransactOpts()
		tx, err := stakingBridge.Stake(opts, missingStake, pubKey)
		if err != nil {
			failedCount += 1
			logger.Error("Failed to stake", zap.String("pub key", pubKey),
				zap.String("staking bridge", stakingBridge.Address.Hex()),
				zap.String("nonce", opts.Nonce.String()),
				zap.String("wallet", wallet.Address.Hex()), zap.Error(err))
		} else {
			stakeTxs[pubKey] = tx
			logger.Debug("Sent stake transaction", zap.String("pub key", pubKey),
				zap.String("staking bridge", stakingBridge.Address.Hex()),
				zap.String("nonce", opts.Nonce.String()),
				zap.String("wallet", wallet.Address.Hex()), zap.Any("tx", tx))
		}
	}
	logger.Info("Waiting for stake transactions", zap.Int("count", len(missingStakeByPubKey)),
		zap.Int("failed count", failedCount))
	//
	// wait for transactions to mint
	//
	for pubKey, tx := range stakeTxs {
		if tx == nil {
			continue
		}
		logger.Debug("waiting", zap.Any("tx", tx))
		if err = ethutils.WaitForTransaction(network.EthClient, tx, time.Minute); err != nil {
			failedCount += 1
			logger.Error("Stake transaction failed", zap.String("pub key", pubKey),
				zap.Uint64("nonce", tx.Nonce()), zap.Any("tx", tx), zap.Error(err))
		} else {
			successCount += 1
			logger.Debug("Stake transaction successful", zap.String("pub key", pubKey),
				zap.Uint64("nonce", tx.Nonce()), zap.String("tx", tx.Hash().Hex()))
		}
	}
	logger.Info("Finished staking", zap.Int("success count", successCount), zap.Int("failed count", failedCount),
		zap.Int("not needed count", len(expectedStakeForParty)-successCount-failedCount))
	if failedCount > 0 {
		return fmt.Errorf("stake for %d parties failed", failedCount)
	}
	return nil
}
