package secrets

import (
	"fmt"
	"log"
	"os"

	"github.com/vegaprotocol/devopstools/generate"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type CreateArgs struct {
	*SecretsArgs
	VegaNetworkName string
	Count           uint16
}

var createArgs CreateArgs

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Secrets for new nodes in network",
	Long:  `Create Secrets for new nodes in network`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunCreate(createArgs); err != nil {
			createArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	createArgs.SecretsArgs = &secretsArgs

	SecretsCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVar(&createArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := createCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	createCmd.PersistentFlags().Uint16Var(&createArgs.Count, "count", 0, "Make sure this number of nodes exists in the network. Does not recreate secrets for existing nodes, also does not delete nodes if those already exists. But fills data if anything is missing.")
	if err := createCmd.MarkPersistentFlagRequired("count"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunCreate(args CreateArgs) error {
	secretStore, err := args.GetNodeSecretStore()
	if err != nil {
		return err
	}
	nodes, err := secretStore.GetAllVegaNode(args.VegaNetworkName)
	if err != nil {
		return err
	}
	for i := 0; i < int(args.Count); i += 1 {
		nodeId := fmt.Sprintf("n%02d", i)
		if args.VegaNetworkName == "sentrynode" {
			nodeId = fmt.Sprintf("sn%03d", i)
		}
		node, ok := nodes[nodeId]
		if ok {
			updates, err := fillOutNodeData(node)
			if err != nil {
				return err
			}
			if len(updates) == 0 {
				args.Logger.Info("already there", zap.String("node", nodeId))
				continue
			}
			if err := secretStore.StoreVegaNode(args.VegaNetworkName, nodeId, node); err != nil {
				return err
			}
			args.Logger.Info("updated", zap.String("node", nodeId), zap.Any("fields", updates))
		} else {
			newSecrets, err := generate.GenerateVegaNodeSecrets()
			if err != nil {
				return err
			}

			if err = secretStore.StoreVegaNode(args.VegaNetworkName, nodeId, newSecrets); err != nil {
				return err
			}
			args.Logger.Info("created new", zap.String("node", nodeId))
		}
	}

	return nil
}
