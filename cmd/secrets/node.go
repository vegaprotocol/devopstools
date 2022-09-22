package secrets

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type NodeArgs struct {
	*SecretsArgs
	VegaNetworkName string
	NodeId          string
}

var nodeArgs NodeArgs

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Get Vega Network /statistics",
	Long:  `Get Vega Network /statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunNode(nodeArgs); err != nil {
			nodeArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	nodeArgs.SecretsArgs = &secretsArgs

	SecretsCmd.AddCommand(nodeCmd)
	nodeCmd.PersistentFlags().StringVar(&nodeArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := nodeCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	nodeCmd.PersistentFlags().StringVar(&nodeArgs.NodeId, "node", "", "Node id, e.g. n01")
	if err := nodeCmd.MarkPersistentFlagRequired("node"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunNode(args NodeArgs) error {
	secretStore, err := args.GetNodeSecretStore()
	if err != nil {
		return err
	}
	secret, err := secretStore.GetVegaNode(args.VegaNetworkName, args.NodeId)
	if err != nil {
		return err
	}
	byteSecret, err := json.MarshalIndent(secret, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to parse node '%s' secret data for network '%s', %w", args.NodeId, args.VegaNetworkName, err)
	}
	fmt.Println(string(byteSecret))
	return nil
}
