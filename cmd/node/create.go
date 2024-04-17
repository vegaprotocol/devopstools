package node

import (
	"fmt"
	"log"
	"os"

	"github.com/vegaprotocol/devopstools/generation"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type CreateArgs struct {
	*Args
	VegaNetworkName string
	NodeId          string
	Force           bool
	Stake           bool
}

var createArgs CreateArgs

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new node",
	Long:  "Create new node",
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunCreateNode(createArgs); err != nil {
			createArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	createArgs.Args = &args

	Cmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVar(&createArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := createCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	createCmd.PersistentFlags().StringVar(&createArgs.NodeId, "node", "", "Node for which create secrets, e.g. n01")
	if err := createCmd.MarkPersistentFlagRequired("node"); err != nil {
		log.Fatalf("%v\n", err)
	}
	createCmd.PersistentFlags().BoolVar(&createArgs.Force, "force", false, "Force to push new secrets, even if the secrets already exist")
}

func RunCreateNode(args CreateArgs) error {
	secretStore, err := args.GetNodeSecretStore()
	if err != nil {
		return err
	}
	oldNodeData, err := secretStore.GetVegaNode(args.VegaNetworkName, args.NodeId)
	if err == nil && oldNodeData != nil && !args.Force {
		return fmt.Errorf("secrets for node %s already exists, use --force to override and put new secrets for it", args.NodeId)
	}

	newSecrets, err := generation.GenerateVegaNodeSecrets()
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
