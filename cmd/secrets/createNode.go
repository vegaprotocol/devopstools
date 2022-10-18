package secrets

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/generate"
	"github.com/vegaprotocol/devopstools/networktools"
	"go.uber.org/zap"
)

type CreateNodeArgs struct {
	*SecretsArgs
	VegaNetworkName string
	NodeId          string
	Force           bool
	Stake           bool
}

var createNodeArgs CreateNodeArgs

// createNodeCmd represents the createNode command
var createNodeCmd = &cobra.Command{
	Use:   "create-node",
	Short: "Create New Secrets for Node",
	Long:  `Create New Secrets for Node`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunCreateNode(createNodeArgs); err != nil {
			createNodeArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	createNodeArgs.SecretsArgs = &secretsArgs

	SecretsCmd.AddCommand(createNodeCmd)
	createNodeCmd.PersistentFlags().StringVar(&createNodeArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := createNodeCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	createNodeCmd.PersistentFlags().StringVar(&createNodeArgs.NodeId, "node", "", "Node for which create secrets, e.g. n01")
	if err := createNodeCmd.MarkPersistentFlagRequired("node"); err != nil {
		log.Fatalf("%v\n", err)
	}
	createNodeCmd.PersistentFlags().BoolVar(&createNodeArgs.Force, "force", false, "Force to push new secrets, even if the secrets already exist")
	createNodeCmd.PersistentFlags().BoolVar(&createNodeArgs.Stake, "stake", false, "Stake Vega token to newly created VegaPub key. If replacing keys, then unstake will be called first.")
}

func RunCreateNode(args CreateNodeArgs) error {
	secretStore, err := args.GetNodeSecretStore()
	if err != nil {
		return err
	}
	oldNodeData, err := secretStore.GetVegaNode(args.VegaNetworkName, args.NodeId)
	if err == nil && oldNodeData != nil && !args.Force {
		return fmt.Errorf("secrets for node %s already exists, use --force to override and put new secrets for it", args.NodeId)
	}

	newSecrets, err := generate.GenerateVegaNodeSecrets()
	if err != nil {
		return err
	}
	byteResult, err := json.MarshalIndent(newSecrets, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to parse createNode '%s' secret data for network '%s', %w", args.NodeId, args.VegaNetworkName, err)
	}
	fmt.Println(string(byteResult))

	if err = secretStore.StoreVegaNode(args.VegaNetworkName, args.NodeId, newSecrets); err != nil {
		return err
	}

	if args.Stake {
		//
		// Get Smart Contracts for Network
		//
		network, err := networktools.NewNetworkTools(args.VegaNetworkName, args.Logger)
		if err != nil {
			return err
		}
		ethClientManager, err := args.GetEthereumClientManager()
		if err != nil {
			return err
		}
		smartContracts, err := network.GetSmartContracts(ethClientManager)
		if err != nil {
			return err
		}
		//
		// Get Minimum Validator Stake
		//
		networkParams, err := network.GetNetworkParams()
		if err != nil {
			return err
		}
		minStake, err := networkParams.GetMinimumValidatorStake()
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

		//
		// Remove stake if needed, then top up the new vega pub key
		//
		if oldNodeData != nil {
			_ = smartContracts.RemoveStake(mainWallet, oldNodeData.VegaPubKey)
		}
		if err := smartContracts.TopUpStakeForONe(mainWallet, newSecrets.VegaPubKey, minStake); err != nil {
			return err
		}
	}
	return nil
}
