package validator

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/generate"
	"github.com/vegaprotocol/devopstools/networktools"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
)

type JoinArgs struct {
	*ValidatorArgs
	VegaNetworkName string
	NodeId          string
	GenerateSecrets bool
	UnstakeFromOld  bool
	Stake           bool
	SelfDelegate    bool
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
	joinArgs.ValidatorArgs = &validatorArgs

	ValidatorCmd.AddCommand(joinCmd)
	joinCmd.PersistentFlags().StringVar(&joinArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := joinCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	joinCmd.PersistentFlags().StringVar(&joinArgs.NodeId, "node", "", "Node for which execute actions, e.g. n01")
	if err := joinCmd.MarkPersistentFlagRequired("node"); err != nil {
		log.Fatalf("%v\n", err)
	}
	joinCmd.PersistentFlags().BoolVar(&joinArgs.GenerateSecrets, "generate-new-secrets", false, "Generate new secrets and push them to the Vault. Note: stake from the old vegaPubKey is not removed")
	joinCmd.PersistentFlags().BoolVar(&joinArgs.UnstakeFromOld, "unstake-from-old-secrets", false, "Unstake from old vegaPubKey. Used together with --generate-new-secrets")
	joinCmd.PersistentFlags().BoolVar(&joinArgs.Stake, "stake", false, "Stake Vega token to validator's VegaPub key. Skip if there is enough stake already.")
	joinCmd.PersistentFlags().BoolVar(&joinArgs.SelfDelegate, "self-delegate", false, "Delegate from node's vegaPubKey to node's id. You need to stake to node's vegaPubKey first.")
}

func RunJoin(args JoinArgs) error {
	var (
		oldNodeSecrets     *secrets.VegaNodePrivate
		currentNodeSecrets *secrets.VegaNodePrivate
		minValidatorStake  *big.Int

		networkParams *types.NetworkParams
		network       *networktools.NetworkTools
		secretStore   secrets.NodeSecretStore
		err           error
	)

	args.Logger.Info("executing Join",
		zap.String("network", args.VegaNetworkName), zap.String("node", args.NodeId), zap.Bool("generate secrets", args.GenerateSecrets),
		zap.Bool("unstake from old vegaPubKey", args.UnstakeFromOld), zap.Bool("stake", args.Stake), zap.Bool("self delegate", args.SelfDelegate),
	)

	//
	// Prepare
	//
	network, err = networktools.NewNetworkTools(args.VegaNetworkName, args.Logger)
	if err != nil {
		return err
	}
	// Get Minimum Validator Stake
	networkParams, err = network.GetNetworkParams()
	if err != nil {
		return err
	}
	minValidatorStake, err = networkParams.GetMinimumValidatorStake()
	if err != nil {
		return err
	}
	// Get Node Secrets
	secretStore, err = args.GetNodeSecretStore()
	if err != nil {
		return err
	}
	currentNodeSecrets, err = secretStore.GetVegaNode(args.VegaNetworkName, args.NodeId)

	//
	// Node Secrets
	//
	if !args.GenerateSecrets {
		if err != nil || currentNodeSecrets == nil {
			return fmt.Errorf("failed to get secrets for node %s in %s network, please use --generate-new-secrets to generate secrets for node, %w",
				args.NodeId, args.VegaNetworkName, err)
		}
	} else {
		//
		// Generate new secrets for node
		//
		oldNodeSecrets = currentNodeSecrets
		currentNodeSecrets, err = generate.GenerateVegaNodeSecrets()
		if err != nil {
			return err
		}
		if err = secretStore.StoreVegaNode(args.VegaNetworkName, args.NodeId, currentNodeSecrets); err != nil {
			return err
		}
		args.Logger.Info("generated and stored new secrets for node",
			zap.String("new vegaPubKey", currentNodeSecrets.VegaPubKey),
			zap.String("new eth wallet", currentNodeSecrets.EthereumAddress),
		)

		//
		// Get Smart Contracts for Network
		//
		ethClientManager, err := args.GetEthereumClientManager()
		if err != nil {
			return err
		}
		smartContracts, err := network.GetSmartContracts(ethClientManager)
		if err != nil {
			return err
		}
		//
		// Get Ethereum Wallet
		//
		walletManager, err := args.GetWalletManager()
		if err != nil {
			return err
		}
		ethNetwork, err := network.GetEthNetwork()
		if err != nil {
			return err
		}
		mainWallet, err := walletManager.GetNetworkMainEthWallet(ethNetwork, args.VegaNetworkName)
		if err != nil {
			return err
		}

		if args.UnstakeFromOld {
			if oldNodeSecrets == nil {
				args.Logger.Info("Skip unstake from old: there was no previous vegaPubKey")
			} else {
				if err = smartContracts.RemoveStake(mainWallet, oldNodeSecrets.VegaPubKey); err != nil {
					return fmt.Errorf("failed to remove stake from old vega pub key %s, %w", oldNodeSecrets.VegaPubKey, err)
				}
				args.Logger.Info("Removed stake from old vega pub key", zap.String("old vegaPubKey", oldNodeSecrets.VegaPubKey))
			}
		}

		if err = smartContracts.TopUpStakeForOne(mainWallet, currentNodeSecrets.VegaPubKey, minValidatorStake); err != nil {
			return fmt.Errorf("failed to top up stake, %w", err)
		}
		args.Logger.Info("Staked to new vega pub key", zap.String("vegaPubKey", currentNodeSecrets.VegaPubKey),
			zap.String("amount", minValidatorStake.String()))
	}

	if args.SelfDelegate {
		var (
			epochValidator   *vega.Node
			stakedByOperator = big.NewInt(0)
			pendingStake     = big.NewInt(0)
			ok               bool
		)
		dataNodeClient, err := network.GetDataNodeClient()
		if err != nil {
			return fmt.Errorf("failed to self-delegate, %w", err)
		}
		//
		// Get current delegation
		//
		epoch, err := dataNodeClient.GetCurrentEpoch()
		if err != nil {
			return fmt.Errorf("failed to self-delegate, %w", err)
		}
		for _, v := range epoch.Validators {
			if v.Id == currentNodeSecrets.VegaId {
				epochValidator = v
				break
			}
		}
		if epochValidator != nil {
			stakedByOperator, ok = stakedByOperator.SetString(epochValidator.StakedByOperator, 0)
			if !ok {
				args.Logger.Error("failed to parse StakedByOperator", zap.String("node", currentNodeSecrets.Name),
					zap.String("StakedByOperator", epochValidator.StakedByOperator), zap.Error(err))
			}
			pendingStake, ok = pendingStake.SetString(epochValidator.PendingStake, 0)
			if !ok {
				args.Logger.Error("failed to parse PendingStake", zap.String("node", currentNodeSecrets.Name),
					zap.String("PendingStake", epochValidator.PendingStake), zap.Error(err))
			}
			args.Logger.Info("found validator", zap.String("vega Id", currentNodeSecrets.VegaId),
				zap.String("vega pub key", currentNodeSecrets.VegaPubKey), zap.String("stakedByOperator", stakedByOperator.String()),
				zap.String("pendingStake", pendingStake.String()), zap.String("minValidatorStake", minValidatorStake.String()))
			stakedByOperator = stakedByOperator.Add(stakedByOperator, pendingStake)
		}
		if stakedByOperator.Cmp(minValidatorStake) >= 0 {
			args.Logger.Info("no need to stake", zap.String("validator", currentNodeSecrets.Name))
		} else {
			//
			// Get Current Total Stake
			//
			partyTotalStake, err := dataNodeClient.GetPartyTotalStake(currentNodeSecrets.VegaPubKey)
			if err != nil {
				return fmt.Errorf("failed to self-delegate, %w", err)
			}
			if partyTotalStake.Cmp(minValidatorStake) < 0 {
				return fmt.Errorf("failed to self-delegate, party %s has %s stake which is less than required min %s",
					currentNodeSecrets.VegaPubKey, partyTotalStake.String(), minValidatorStake.String())
			}
			//
			// Submit Delegate Transaction
			//
			args.Logger.Info("data", zap.String("VegaId", currentNodeSecrets.VegaId), zap.String("VegaPubKey", currentNodeSecrets.VegaPubKey))
			vegawallet, err := wallet.NewVegaWallet(&secrets.VegaWalletPrivate{
				Id:             currentNodeSecrets.VegaId,
				PublicKey:      currentNodeSecrets.VegaPubKey,
				PrivateKey:     currentNodeSecrets.VegaPrivateKey,
				RecoveryPhrase: currentNodeSecrets.VegaRecoveryPhrase,
			})
			if err != nil {
				return fmt.Errorf("failed to self-delegate, %w", err)
			}
			lastBlockData, err := dataNodeClient.LastBlockData()
			if err != nil {
				return fmt.Errorf("failed to self-delegate, %w", err)
			}
			signedTx, err := vegawallet.SignTxWithPoW(&walletpb.SubmitTransactionRequest{
				PubKey: currentNodeSecrets.VegaPubKey,
				Command: &walletpb.SubmitTransactionRequest_DelegateSubmission{
					DelegateSubmission: &commandspb.DelegateSubmission{
						NodeId: currentNodeSecrets.VegaId,
						Amount: minValidatorStake.String(),
					},
				},
			}, lastBlockData)
			if err != nil {
				return fmt.Errorf("failed to self-delegate, %w", err)
			}

			submitResponse, err := dataNodeClient.SubmitTransaction(&vegaapipb.SubmitTransactionRequest{
				Tx:   signedTx,
				Type: vegaapipb.SubmitTransactionRequest_TYPE_SYNC,
			})
			if err != nil {
				args.Logger.Error("failed to submit a trasnaction", zap.String("node", currentNodeSecrets.Name), zap.Error(err))
				return fmt.Errorf("failed to self-delegate, %w", err)
			}
			args.Logger.Info("tx", zap.Any("submitResponse", submitResponse))
			if !submitResponse.Success {
				args.Logger.Error("transaction submission failure", zap.String("node", currentNodeSecrets.Name), zap.Error(err))
				return fmt.Errorf("failed to self-delegate, %w", err)
			}
		}
	}

	return nil
}
