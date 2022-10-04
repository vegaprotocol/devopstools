package veganetwork

import (
	"fmt"

	"code.vegaprotocol.io/vega/protos/vega"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/smartcontracts"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
)

type VegaNetwork struct {
	Network        string
	SmartContracts *smartcontracts.VegaNetworkSmartContracts

	ValidatorsById map[string]*vega.Node

	// wallets
	NodeSecrets       map[string]secrets.VegaNodePrivate
	NetworkMainWallet *wallet.EthWallet
	AssetMainWallet   *wallet.EthWallet

	MarketsCreator *secrets.VegaWalletPrivate
	VegaTokenWhale *wallet.VegaWallet

	// network params/config
	NetworkParams  *NetworkParams
	EthereumConfig *vega.EthereumConfig
	EthNetwork     types.ETHNetwork

	// clients
	DataNodeClient        *vegaapi.DataNode
	SmartContractsManager *smartcontracts.SmartContractsManager
	WalletManager         *wallet.WalletManager
	EthClient             *ethclient.Client
	NodeSecretStore       secrets.NodeSecretStore
}

func NewVegaNetwork(
	network string,
	dataNodeClient *vegaapi.DataNode,
	nodeSecretStore secrets.NodeSecretStore,
	smartContractsManager *smartcontracts.SmartContractsManager,
	walletManager *wallet.WalletManager,
	logger *zap.Logger,
) (*VegaNetwork, error) {
	var (
		n = &VegaNetwork{
			Network:               network,
			DataNodeClient:        dataNodeClient,
			SmartContractsManager: smartContractsManager,
			WalletManager:         walletManager,
			NodeSecretStore:       nodeSecretStore,
		}
		errMsg = "failed to create VegaNetwork for: %w"
		err    error
	)

	// Read and parse some Network Parameters
	networkParams, err := n.DataNodeClient.GetAllNetworkParameters()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	n.NetworkParams = NewNetworkParams(networkParams)
	n.EthereumConfig, err = n.NetworkParams.GetEthereumConfig()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	n.EthNetwork, err = types.GetEthNetworkForId(n.EthereumConfig.ChainId)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	// Setup Smart Contracts
	n.SmartContracts, err = n.SmartContractsManager.Connect(
		n.EthNetwork,
		"", // will be taken from Staking Bridge
		"", // will be taken from ERC20 Bridge
		n.EthereumConfig.CollateralBridgeContract.Address,
		n.EthereumConfig.MultisigControlContract.Address,
		n.EthereumConfig.StakingBridgeContract.Address,
	)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	n.EthClient = n.SmartContracts.EthClient

	// Node Secrets
	n.NodeSecrets, err = n.NodeSecretStore.GetAllVegaNode(n.Network)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	epoch, err := n.DataNodeClient.GetCurrentEpoch()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	//
	n.ValidatorsById = make(map[string]*vega.Node)
	for _, validator := range epoch.Validators {
		n.ValidatorsById[validator.Id] = validator
	}

	// Wallets
	n.AssetMainWallet, err = n.WalletManager.GetAssetMainEthWallet(n.EthNetwork)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	n.NetworkMainWallet, err = n.WalletManager.GetNetworkMainEthWallet(n.EthNetwork, n.Network)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	n.VegaTokenWhale, err = n.WalletManager.GetVegaTokenWhaleVegaWallet()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return n, nil
}

func (n *VegaNetwork) Disconnect() {
}
