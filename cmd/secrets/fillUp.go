package secrets

import (
	"fmt"
	"log"
	"os"

	"code.vegaprotocol.io/vega/wallet/wallet"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type FillUpArgs struct {
	*SecretsArgs
	VegaNetworkName string
	NodeId          string
}

var fillUpArgs FillUpArgs

// fillUpCmd represents the fillUp command
var fillUpCmd = &cobra.Command{
	Use:   "fill-up",
	Short: "Fill Up missing data in Node Secrets",
	Long:  `Fill Up missing data in Node Secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunFillUp(fillUpArgs); err != nil {
			fillUpArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	fillUpArgs.SecretsArgs = &secretsArgs

	SecretsCmd.AddCommand(fillUpCmd)
	fillUpCmd.PersistentFlags().StringVar(&fillUpArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := fillUpCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	fillUpCmd.PersistentFlags().StringVar(&fillUpArgs.NodeId, "node", "", "Node id, e.g. n01")
	if err := fillUpCmd.MarkPersistentFlagRequired("node"); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func RunFillUp(args FillUpArgs) error {
	secretStore, err := args.GetNodeSecretStore()
	if err != nil {
		return err
	}
	nodeSecrets, err := secretStore.GetVegaNode(args.VegaNetworkName, args.NodeId)
	if err != nil {
		return err
	}
	if len(nodeSecrets.VegaId) == 0 {
		vegaWallet, err := wallet.ImportHDWallet("my wallet", nodeSecrets.VegaRecoveryPhrase, wallet.LatestVersion)
		if err != nil {
			return fmt.Errorf("failed to get vegawallet with recovery phrase %w", err)
		}
		vegaId := vegaWallet.ID()

		nodeSecrets.VegaId = vegaId

		if err := secretStore.StoreVegaNode(args.VegaNetworkName, args.NodeId, nodeSecrets); err != nil {
			return err
		}

		args.Logger.Info("filled in", zap.String("vegaId", vegaId))
	}

	return nil
}
