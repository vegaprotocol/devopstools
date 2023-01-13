package secrets

import (
	"log"
	"os"

	"code.vegaprotocol.io/vega/wallet/wallet"
	"github.com/ipfs/kubo/config"
	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/generate"
	"github.com/vegaprotocol/devopstools/secrets"
	"go.uber.org/zap"
)

type FillOutArgs struct {
	*SecretsArgs
	VegaNetworkName string
	NodeId          string
}

var fillOutArgs FillOutArgs

// fillOutCmd represents the fillOut command
var fillOutCmd = &cobra.Command{
	Use:   "fill-out",
	Short: "Fill Out missing data in Node Secrets",
	Long:  `Fill Out missing data in Node Secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunFillOut(fillOutArgs); err != nil {
			fillOutArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	fillOutArgs.SecretsArgs = &secretsArgs

	SecretsCmd.AddCommand(fillOutCmd)
	fillOutCmd.PersistentFlags().StringVar(&fillOutArgs.VegaNetworkName, "network", "", "Vega Network name")
	if err := fillOutCmd.MarkPersistentFlagRequired("network"); err != nil {
		log.Fatalf("%v\n", err)
	}
	fillOutCmd.PersistentFlags().StringVar(&fillOutArgs.NodeId, "node", "", "Node id, e.g. n01")
}

func RunFillOut(args FillOutArgs) error {
	var (
		nodes map[string]*secrets.VegaNodePrivate
	)
	secretStore, err := args.GetNodeSecretStore()
	if err != nil {
		return err
	}
	if len(args.NodeId) > 0 {
		nodeSecrets, err := secretStore.GetVegaNode(args.VegaNetworkName, args.NodeId)
		if err != nil {
			return err
		}
		nodes = map[string]*secrets.VegaNodePrivate{
			args.NodeId: nodeSecrets,
		}
	} else {
		nodes, err = secretStore.GetAllVegaNode(args.VegaNetworkName)
		if err != nil {
			return err
		}
	}
	for nodeId, node := range nodes {
		updates, err := fillOutNodeData(node)
		if err != nil {
			return err
		}
		if len(updates) == 0 {
			args.Logger.Info("nothing to update", zap.String("nodes", nodeId))
			continue
		}
		if err := secretStore.StoreVegaNode(args.VegaNetworkName, nodeId, node); err != nil {
			return err
		}
		args.Logger.Info("updated", zap.String("node", nodeId), zap.Any("fields", updates))
	}

	return nil
}

func fillOutNodeData(nodeSecrets *secrets.VegaNodePrivate) (updatedFields map[string]string, err error) {
	var (
		vegaWallet *wallet.HDWallet
	)
	updatedFields = map[string]string{}

	vegaWallet, err = wallet.ImportHDWallet("my wallet", nodeSecrets.VegaRecoveryPhrase, wallet.LatestVersion)
	if err != nil {
		return
	}

	if len(nodeSecrets.VegaId) == 0 {
		nodeSecrets.VegaId = vegaWallet.ID()
		updatedFields["VegaId"] = nodeSecrets.VegaId
	}
	if nodeSecrets.VegaPubKeyIndex == nil {
		vegaPubKeyIndex := uint64(1)
		nodeSecrets.VegaPubKeyIndex = &vegaPubKeyIndex
		updatedFields["VegaPubKeyIndex"] = "1"
	}
	if len(nodeSecrets.Name) == 0 {
		nodeSecrets.Name, err = generate.GenerateName()
		if err != nil {
			return
		}
		updatedFields["Name"] = nodeSecrets.Name
	}
	if len(nodeSecrets.Country) == 0 {
		nodeSecrets.Country, err = generate.GenerateCountryCode()
		if err != nil {
			return
		}
		updatedFields["Country"] = nodeSecrets.Country
	}
	if len(nodeSecrets.InfoURL) == 0 {
		nodeSecrets.InfoURL, err = generate.GenerateRandomWikiURL()
		if err != nil {
			return
		}
		updatedFields["InfoURL"] = nodeSecrets.InfoURL
	}
	if len(nodeSecrets.AvatarURL) == 0 {
		nodeSecrets.AvatarURL, err = generate.GenerateAvatarURL()
		if err != nil {
			return
		}
		updatedFields["AvatarURL"] = nodeSecrets.AvatarURL
	}
	if len(nodeSecrets.NetworkHistoryPeerId) == 0 || len(nodeSecrets.NetworkHistoryPrivateKey) == 0 {
		var deHistory config.Identity
		deHistory, err = generate.GenerateNetworkHistoryIdentity(nodeSecrets.VegaRecoveryPhrase)
		if err != nil {
			return
		}
		nodeSecrets.NetworkHistoryPeerId = deHistory.PeerID
		nodeSecrets.NetworkHistoryPrivateKey = deHistory.PrivKey
		updatedFields["NetworkHistoryPeerId"] = nodeSecrets.NetworkHistoryPeerId
		updatedFields["NetworkHistoryPrivateKey"] = nodeSecrets.NetworkHistoryPrivateKey
	}

	return
}
