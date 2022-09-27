package veganetwork

import (
	"fmt"

	"code.vegaprotocol.io/vega/protos/vega"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/smartcontracts"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"go.uber.org/zap"
)

type VegaNetwork struct {
	Network        string
	SmartContracts *smartcontracts.VegaNetworkSmartContracts

	ValidatorsById map[string]*vega.Node

	// wallets
	NodeSecrets       map[string]secrets.VegaNodePrivate
	NetworkMainWallet *secrets.EthereumWalletPrivate
	AssetMainWallet   *secrets.EthereumWalletPrivate

	MarketsCreator *secrets.VegaWalletPrivate
	VegaTokenWhale *secrets.VegaWalletPrivate

	// network params/config
	NetworkParams  *NetworkParams
	EthereumConfig *vega.EthereumConfig
	EthNetwork     types.ETHNetwork

	// clients
	DataNodeClient        *vegaapi.DataNode
	SmartContractsManager *smartcontracts.SmartContractsManager
	EthClient             *ethclient.Client
	NodeSecretStore       secrets.NodeSecretStore
}

func NewVegaNetwork(
	network string,
	dataNodeClient *vegaapi.DataNode,
	nodeSecretStore secrets.NodeSecretStore,
	smartContractsManager *smartcontracts.SmartContractsManager,
	logger *zap.Logger,
) (*VegaNetwork, error) {
	var (
		n = &VegaNetwork{
			Network:               network,
			DataNodeClient:        dataNodeClient,
			SmartContractsManager: smartContractsManager,
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

	return n, nil
}

func (n *VegaNetwork) Disconnect() {
}
