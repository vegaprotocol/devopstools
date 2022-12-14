package secrets

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/generate"
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

	if err = secretStore.StoreVegaNode(args.VegaNetworkName, args.NodeId, newSecrets); err != nil {
		return err
	}

	args.Logger.Info("Generated new secrets for node", zap.String("network", args.VegaNetworkName), zap.String("nodeId", args.NodeId),
		zap.String("Name", newSecrets.Name),
		zap.String("Country", newSecrets.Country),
		zap.String("InfoURL", newSecrets.InfoURL),
		zap.String("AvatarURL", newSecrets.AvatarURL),
		zap.String("VegaPubKey", newSecrets.VegaPubKey),
		zap.String("VegaId", newSecrets.VegaId),
		zap.Uint64("VegaPubKeyIndex", *newSecrets.VegaPubKeyIndex),
		zap.String("EthereumAddress", newSecrets.EthereumAddress),
		zap.String("TendermintNodeId", newSecrets.TendermintNodeId),
		zap.String("TendermintNodePubKey", newSecrets.TendermintNodePubKey),
		zap.String("TendermintValidatorAddress", newSecrets.TendermintValidatorAddress),
		zap.String("TendermintValidatorPubKey", newSecrets.TendermintValidatorPubKey),
	)

	return nil
}
