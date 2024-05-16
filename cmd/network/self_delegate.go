package network

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/networktools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	vegapb "code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type SelfDelegateArgs struct {
	*Args
}

var selfDelegateArgs SelfDelegateArgs

// selfDelegateCmd represents the selfDelegate command
var selfDelegateCmd = &cobra.Command{
	Use:   "self-delegate",
	Short: "Execute self-delegate for validators",
	Long: `Excecute self-delegate process for every validator, and some steps for non-validators.
	It is safe to call it multiple times, cos it won't repeat steps from previous call.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunSelfDelegate(selfDelegateArgs); err != nil {
			selfDelegateArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	selfDelegateArgs.Args = &args

	Cmd.AddCommand(selfDelegateCmd)
}

func RunSelfDelegate(args SelfDelegateArgs) error {
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

	primaryEthConfig, err := networkParams.PrimaryEthereumConfig()
	if err != nil {
		return fmt.Errorf("could not get primary ethereum configuration from network paramters: %w", err)
	}

	primaryChainClient, err := ethereum.NewPrimaryChainClient(ctx, cfg.Bridges.Primary, primaryEthConfig, logger.Named("primary-chain-client"))
	if err != nil {
		return fmt.Errorf("could not initialize primary ethereum chain client: %w", err)
	}

	minValidatorStake, err := networkParams.GetMinimumValidatorStake()
	if err != nil {
		return err
	}

	missingStakeByPubKey := map[string]*types.Amount{}
	for _, node := range cfg.Nodes {
		nodeKey := node.Secrets.VegaPubKey

		currentStakeAsSubUnit, err := datanodeClient.GetPartyTotalStake(nodeKey)
		if err != nil {
			return fmt.Errorf("failed to retrieve current stake for party %s: %w", nodeKey, err)
		}

		currentStake := types.NewAmountFromSubUnit(currentStakeAsSubUnit, 18)
		logger.Debug("Current stake for party found",
			zap.String("party-id", nodeKey),
			zap.String("current-stake", currentStake.String()),
		)

		if currentStake.Cmp(minValidatorStake) < 0 {
			logger.Debug("Party needs more stake",
				zap.String("party-id", nodeKey),
				zap.String("program-minimum-stake", minValidatorStake.String()),
				zap.String("current-stake", currentStake.String()),
			)

			missingStakeByPubKey[nodeKey] = minValidatorStake.Copy()
		} else {
			logger.Debug("Party does not need more stake",
				zap.String("party-id", nodeKey),
				zap.String("program-minimum-stake", minValidatorStake.String()),
				zap.String("current-stake", currentStake.String()),
			)
		}
	}

	if err := primaryChainClient.StakeVegaTokenFromMinter(ctx, missingStakeByPubKey); err != nil {
		return fmt.Errorf("could not stake validator vega token: %w", err)
	}

	//
	// Delegate
	//

	currentEpoch, err := datanodeClient.GetCurrentEpoch()
	if err != nil {
		return fmt.Errorf("could not retrieve current epoch: %w", err)
	}

	validatorsById := map[string]*vegapb.Node{}
	for _, validator := range currentEpoch.Validators {
		validatorsById[validator.Id] = validator
	}

	var wg sync.WaitGroup
	for _, node := range cfg.Nodes {
		name := node.Metadata.Name

		validator, isValidator := validatorsById[node.Secrets.VegaId]
		if !isValidator {
			continue
		}

		stakedByOperatorAsSubUnit, ok := new(big.Int).SetString(validator.StakedByOperator, 0)
		if !ok {
			logger.Error("failed to parse StakedByOperator",
				zap.String("node", name),
				zap.String("value", validator.StakedByOperator),
				zap.Error(err),
			)
		}
		stakedByOperator := types.NewAmountFromSubUnit(stakedByOperatorAsSubUnit, 18)

		if minValidatorStake.Cmp(stakedByOperator) <= 0 {
			logger.Info("node already self-delegated enough",
				zap.String("node", name),
				zap.String("staked-by-operator", stakedByOperator.String()),
			)
			continue
		}

		partyTotalStakeAsSubUnit, err := datanodeClient.GetPartyTotalStake(node.Secrets.VegaId)
		if err != nil {
			return err
		}
		partyTotalStake := types.NewAmountFromSubUnit(partyTotalStakeAsSubUnit, 18)

		if partyTotalStake.Cmp(partyTotalStake) < 0 {
			logger.Warn("party doesn't have visible stake yet - you might need to wait till next epoch",
				zap.String("node", name),
				zap.String("partyTotalStake", partyTotalStakeAsSubUnit.String()),
			)

			// TODO: write wait functionality
			// continue
		}

		wg.Add(1)
		go func(name string, nodeSecrets config.NodeSecrets, hasVisibleStake bool, datanodeClient vegaapi.DataNodeClient) {
			defer wg.Done()
			delegateStake(ctx, nodeSecrets, hasVisibleStake, minValidatorStake, datanodeClient, logger.With(zap.String("node", name)))
		}(name, node.Secrets, partyTotalStake.Cmp(minValidatorStake) >= 0, datanodeClient)
	}
	wg.Wait()

	logger.Info("Self-delegation successful")

	return nil
}

func delegateStake(ctx context.Context, nodeSecrets config.NodeSecrets, hasVisibleStake bool, minValidatorStake *types.Amount, datanodeClient vegaapi.DataNodeClient, logger *zap.Logger) {
	vegawallet, err := vega.LoadWallet(nodeSecrets.VegaId, nodeSecrets.VegaRecoveryPhrase)
	if err != nil {
		logger.Error("failed to create wallet", zap.Error(err))
		return
	}

	if err := vega.GenerateKeysUpToKey(vegawallet, nodeSecrets.VegaPubKey); err != nil {
		logger.Error("failed to generate wallet keys", zap.String("key", nodeSecrets.VegaPubKey), zap.Error(err))
		return
	}

	request := walletpb.SubmitTransactionRequest{
		PubKey: nodeSecrets.VegaPubKey,
		Command: &walletpb.SubmitTransactionRequest_DelegateSubmission{
			DelegateSubmission: &commandspb.DelegateSubmission{
				NodeId: nodeSecrets.VegaId,
				Amount: minValidatorStake.String(),
			},
		},
	}
	logger.Info("Wallet request", zap.Any("request", request))

	attemptLeft := 5
	for 0 < attemptLeft {
		submitResponse, err := walletpkg.SendTransaction(ctx, vegawallet, nodeSecrets.VegaPubKey, &request, datanodeClient)
		if err != nil {
			logger.Error("Failed to submit a transaction", zap.Int("attempt-left", attemptLeft), zap.Error(err))
			continue
		}

		logger.Info("Transaction response", zap.Any("response", submitResponse))
		if submitResponse.Success {
			logger.Info("successful delegation", zap.Int("attempt-left", attemptLeft), zap.Any("response", submitResponse))
			return
		}

		if hasVisibleStake {
			return
		}
		attemptLeft -= 1
		if attemptLeft > 0 {
			logger.Warn("Transaction submission failure, retrying in 1 minute...",
				zap.Int("attempt-left", attemptLeft),
				zap.Error(err),
			)
			time.Sleep(time.Minute)
		}
	}

	logger.Error("Transaction failed")
}
