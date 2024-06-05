package node

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/vegaprotocol/devopstools/config"
	"github.com/vegaprotocol/devopstools/generation"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type CreateArgs struct {
	*Args
	NodeID string
	Force  bool
	Stake  bool
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

	createCmd.PersistentFlags().StringVar(&createArgs.NodeID, "node", "", "Node ID for which create secrets, e.g. n01")
	if err := createCmd.MarkPersistentFlagRequired("node"); err != nil {
		log.Fatalf("%v\n", err)
	}
	createCmd.PersistentFlags().BoolVar(&createArgs.Force, "force", false, "Force to push new secrets, even if the secrets already exist")
}

func RunCreateNode(args CreateArgs) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := args.Logger.Named("command")

	cfg, err := config.Load(ctx, args.NetworkFile)
	if err != nil {
		return fmt.Errorf("could not load network file at %q: %w", args.NetworkFile, err)
	}
	logger.Debug("Network file loaded", zap.String("name", cfg.Name.String()))

	for _, node := range cfg.Nodes {
		if node.ID == args.NodeID && !args.Force {
			return fmt.Errorf("node %s already exists, use --force to override and put new secrets for it", args.NodeID)
		}
	}

	newNode, err := generation.GenerateVegaNodeSecrets()
	if err != nil {
		return fmt.Errorf("failed to generate vega node: %w", err)
	}

	newNode.ID = args.NodeID

	cfg = config.UpsertNode(cfg, newNode)

	if err := config.SaveConfig(args.NetworkFile, cfg); err != nil {
		return fmt.Errorf("could not save network file at %q: %w", args.NetworkFile, err)
	}

	args.Logger.Info("Generated new secrets for node", zap.String("network", cfg.Name.String()), zap.String("nodeId", args.NodeID),
		zap.String("Name", newNode.Metadata.Name),
		zap.String("Country", newNode.Metadata.Country),
		zap.String("InfoURL", newNode.Metadata.InfoURL),
		zap.String("AvatarURL", newNode.Metadata.AvatarURL),
		zap.String("VegaPubKey", newNode.Secrets.VegaPubKey),
		zap.String("VegaId", newNode.Secrets.VegaId),
		zap.Uint64("VegaPubKeyIndex", *newNode.Secrets.VegaPubKeyIndex),
		zap.String("EthereumAddress", newNode.Secrets.EthereumAddress),
		zap.String("TendermintNodeId", newNode.Secrets.TendermintNodeId),
		zap.String("TendermintNodePubKey", newNode.Secrets.TendermintNodePubKey),
		zap.String("TendermintValidatorAddress", newNode.Secrets.TendermintValidatorAddress),
		zap.String("TendermintValidatorPubKey", newNode.Secrets.TendermintValidatorPubKey),
	)

	return nil
}
