package validator

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/generation"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vega"
	"github.com/vegaprotocol/devopstools/vegaapi/core"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	vegapb "code.vegaprotocol.io/vega/protos/vega"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	walletpkg "code.vegaprotocol.io/vega/wallet/pkg"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type JoinArgs struct {
	*Args
	NodeId                      string
	GenerateSecrets             bool
	UnstakeFromOld              bool
	Stake                       bool
	SelfDelegate                bool
	GetEthAddressToSubmitBundle bool
	SendEthereumEvents          bool
}

var joinArgs JoinArgs

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Validator actions required during Validator join procedure",
	Long:  `Validator actions required during Validator join procedure`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunJoin(joinArgs); err != nil {
			joinArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	joinArgs.Args = &args

	Cmd.AddCommand(joinCmd)
	joinCmd.PersistentFlags().StringVar(&joinArgs.NodeId, "node", "", "Node for which execute actions, e.g. n01")
	if err := joinCmd.MarkPersistentFlagRequired("node"); err != nil {
		log.Fatalf("%v\n", err)
	}
	joinCmd.PersistentFlags().BoolVar(&joinArgs.GenerateSecrets, "generate-new-secrets", false, "Generate new secrets and push them to the Vault. Note: stake from the old vegaPubKey is not removed")
	joinCmd.PersistentFlags().BoolVar(&joinArgs.UnstakeFromOld, "unstake-from-old-secrets", false, "Unstake from old vegaPubKey. Used together with --generate-new-secrets")
	joinCmd.PersistentFlags().BoolVar(&joinArgs.Stake, "stake", false, "Stake Vega token to validator's VegaPub key. Skip if there is enough stake already.")
	joinCmd.PersistentFlags().BoolVar(&joinArgs.SelfDelegate, "self-delegate", false, "Delegate from node's vegaPubKey to node's id. You need to stake to node's vegaPubKey first.")
	joinCmd.PersistentFlags().BoolVar(&joinArgs.GetEthAddressToSubmitBundle, "get-eth-to-submit-bundle", false, "Prints ethereum address of a wallet that will be used to submit Multisig Control bundle.")
	joinCmd.PersistentFlags().BoolVar(&joinArgs.SendEthereumEvents, "send-ethereum-events", false, "Send ethereum events to the node. Required for new Validator by network parameter")
}

func RunJoin(args JoinArgs) error {
	ctx, cancelCommand := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelCommand()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Debug("Network file loaded", zap.String("name", cfg.Name.String()))

	args.Logger.Info("executing Join",
		zap.String("node", args.NodeId),
		zap.Bool("generate secrets", args.GenerateSecrets),
		zap.Bool("unstake from old vegaPubKey", args.UnstakeFromOld),
		zap.Bool("stake", args.Stake),
		zap.Bool("self delegate", args.SelfDelegate),
		zap.Bool("send ethereum events", args.SendEthereumEvents),
	)

	endpoints := config.ListDatanodeGRPCEndpoints(cfg)
	if len(endpoints) == 0 {
		return fmt.Errorf("no gRPC endpoint found on configured datanodes")
	}
	logger.Debug("gRPC endpoints found in network file", zap.Strings("endpoints", endpoints))

	logger.Debug("Looking for healthy gRPC endpoints...")
	healthyEndpoints := tools.FilterHealthyGRPCEndpoints(endpoints)
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

	minValidatorStake, err := networkParams.GetMinimumValidatorStake()
	if err != nil {
		return err
	}

	currentNodeSecrets, nodeAlreadyExists := config.FindNodeByID(cfg, args.NodeId)

	if !args.GenerateSecrets {
		if !nodeAlreadyExists {
			return fmt.Errorf("failed to get secrets for node %q in %q network, please use --generate-new-secrets to generate secrets for node, %w", args.NodeId, cfg.Name, err)
		}
	} else {
		previousNodeSecrets := currentNodeSecrets
		currentNodeSecrets, err := generation.GenerateVegaNodeSecrets()
		if err != nil {
			return err
		}

		currentNodeSecrets.ID = previousNodeSecrets.ID

		cfg = config.UpsertNode(cfg, currentNodeSecrets)

		if err := config.SaveConfig(args.NetworkFile, cfg); err != nil {
			return fmt.Errorf("could not save network file at %q: %w", args.NetworkFile, err)
		}

		args.Logger.Info("generated and stored new secrets for node",
			zap.String("new vegaPubKey", currentNodeSecrets.Secrets.VegaPubKey),
			zap.String("new eth wallet", currentNodeSecrets.Secrets.EthereumAddress),
		)

		primaryEthConfig, err := networkParams.PrimaryEthereumConfig()
		if err != nil {
			return fmt.Errorf("could not get primary ethereum configuration from network paramters: %w", err)
		}

		primaryChainClient, err := ethereum.NewPrimaryChainClient(ctx, cfg.Bridges.Primary, primaryEthConfig, logger.Named("primary-chain-client"))
		if err != nil {
			return fmt.Errorf("could not initialize primary ethereum chain client: %w", err)
		}

		if args.UnstakeFromOld {
			if !nodeAlreadyExists {
				args.Logger.Info("Skip unstake from old: there was no previous vegaPubKey")
			} else {
				if err = primaryChainClient.RemoveMinterStake(ctx, previousNodeSecrets.Secrets.VegaPubKey); err != nil {
					return fmt.Errorf("failed to remove stake from old vega public key %s: %w", previousNodeSecrets.Secrets.VegaPubKey, err)
				}
				args.Logger.Info("Removed stake from old vega pub key",
					zap.String("old vegaPubKey", previousNodeSecrets.Secrets.VegaPubKey),
				)
			}
		}

		missingStakes := map[string]*types.Amount{currentNodeSecrets.Secrets.VegaPubKey: minValidatorStake}
		if err = primaryChainClient.StakeFromMinter(ctx, missingStakes); err != nil {
			return fmt.Errorf("failed to top up stake, %w", err)
		}

		args.Logger.Info("Staked to new vega pub key",
			zap.String("public-key", currentNodeSecrets.Secrets.VegaPubKey),
			zap.String("amount", minValidatorStake.String()),
		)
	}

	if args.SelfDelegate {
		var epochValidator *vegapb.Node

		stakedByOperator := types.NewAmount(18)

		epoch, err := datanodeClient.GetCurrentEpoch()
		if err != nil {
			return fmt.Errorf("failed to self-delegate, %w", err)
		}

		for _, v := range epoch.Validators {
			if v.Id == currentNodeSecrets.Secrets.VegaId {
				epochValidator = v
				break
			}
		}

		if epochValidator != nil {
			stakedByOperator = types.ParseAmountFromSubUnit(epochValidator.StakedByOperator, 18)
			pendingStake := types.ParseAmountFromSubUnit(epochValidator.PendingStake, 18)

			args.Logger.Info("found validator",
				zap.String("vega Id", currentNodeSecrets.Secrets.VegaId),
				zap.String("vega pub key", currentNodeSecrets.Secrets.VegaPubKey),
				zap.String("stakedByOperator", stakedByOperator.String()),
				zap.String("pendingStake", pendingStake.String()),
				zap.String("minValidatorStake", minValidatorStake.String()),
			)
			stakedByOperator.Add(pendingStake.AsMainUnit())
		}

		if stakedByOperator.Cmp(minValidatorStake) >= 0 {
			args.Logger.Info("no need to stake", zap.String("validator", currentNodeSecrets.Metadata.Name))
		} else {
			partyTotalStakeAsSubUnit, err := datanodeClient.GetPartyTotalStake(currentNodeSecrets.Secrets.VegaPubKey)
			if err != nil {
				return fmt.Errorf("failed to self-delegate, %w", err)
			}

			partyTotalStake := types.NewAmountFromSubUnit(partyTotalStakeAsSubUnit, 18)

			if partyTotalStake.Cmp(minValidatorStake) < 0 {
				return fmt.Errorf("failed to self-delegate, party %s has %s stake which is less than required min %s",
					currentNodeSecrets.Secrets.VegaPubKey, partyTotalStake.String(), minValidatorStake.String())
			}

			vegawallet, err := vega.LoadWallet(currentNodeSecrets.Secrets.VegaId, currentNodeSecrets.Secrets.VegaRecoveryPhrase)
			if err != nil {
				return fmt.Errorf("failed to self-delegate, %w", err)
			}

			if err := vega.GenerateKeysUpToKey(vegawallet, currentNodeSecrets.Secrets.VegaPubKey); err != nil {
				return fmt.Errorf("could not generate keys: %w", err)
			}

			request := walletpb.SubmitTransactionRequest{
				Command: &walletpb.SubmitTransactionRequest_DelegateSubmission{
					DelegateSubmission: &commandspb.DelegateSubmission{
						NodeId: currentNodeSecrets.Secrets.VegaId,
						Amount: minValidatorStake.String(),
					},
				},
			}

			response, err := walletpkg.SendTransaction(ctx, vegawallet, currentNodeSecrets.Secrets.VegaPubKey, &request, datanodeClient)
			if err != nil {
				return err
			}

			args.Logger.Info("tx", zap.Any("submitResponse", response))
			if !response.Success {
				args.Logger.Error("transaction submission failure", zap.String("node", currentNodeSecrets.Metadata.Name), zap.Error(err))
				return fmt.Errorf("failed to self-delegate, %w", err)
			}
		}
	}

	if args.SendEthereumEvents {
		eventNum, err := networkParams.GetMinimumEthereumEventsForNewValidator()
		if err != nil {
			return err
		}

		coreClient := core.NewClient([]string{currentNodeSecrets.API.VegaGRPCURL}, 3*time.Second, args.Logger.Named("core"))

		logger.Debug("Connecting to a datanode's gRPC endpoint...")
		dialCtx, cancelDialing := context.WithTimeout(ctx, 2*time.Second)
		defer cancelDialing()
		datanodeClient.MustDialConnection(dialCtx) // blocking
		logger.Debug("Connected to a datanode's gRPC node", zap.String("node", datanodeClient.Target()))

		faucetSecrets := cfg.Network.Wallets.Faucet

		faucetVegaWallet, err := vega.LoadWallet(faucetSecrets.Name, faucetSecrets.RecoveryPhrase)
		if err != nil {
			return err
		}

		if err := vega.GenerateKeysUpToKey(faucetVegaWallet, faucetSecrets.PublicKey); err != nil {
			return err
		}

		vegaAssetId := "XYZepsilon"
		partyId := currentNodeSecrets.Secrets.VegaPubKey
		amount := "3"

		for i := 0; i < eventNum+10; i += 1 {
			result, err := coreClient.DepositBuiltinAsset(vegaAssetId, partyId, amount, func(data []byte) ([]byte, string, error) {
				sig, err := faucetVegaWallet.SignAny(faucetSecrets.PublicKey, data)
				return sig, faucetSecrets.PublicKey, err
			})
			if err != nil {
				return err
			}
			status := "OK"
			if !result {
				status = "FAILED"
			}
			args.Logger.Info("Deposit", zap.String("status", status), zap.String("vegaAssetId", vegaAssetId),
				zap.String("partyId", partyId), zap.String("amount", amount),
			)
		}
	}

	return nil
}
