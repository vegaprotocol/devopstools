package veganetwork

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/smartcontracts"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/veganetworksmartcontracts"
	"github.com/vegaprotocol/devopstools/wallet"

	"code.vegaprotocol.io/vega/protos/vega"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type VegaNetwork struct {
	Network string

	ValidatorsById map[string]*vega.Node

	// wallets
	NodeSecrets       map[string]*secrets.VegaNodePrivate
	NetworkMainWallet *ethereum.Wallet
	AssetMainWallet   *ethereum.Wallet
	BotsApiToken      string

	MarketsCreator *secrets.VegaWalletPrivate
	VegaTokenWhale *wallet.VegaWallet

	// network params/config
	NetworkParams *types.NetworkParams

	// clients
	DataNodeClient vegaapi.DataNodeClient

	PrimaryEthereumConfig        *vega.EthereumConfig
	PrimaryEthNetwork            types.ETHNetwork
	PrimaryEthClientManager      *ethutils.EthereumClientManager
	PrimaryEthClient             *ethclient.Client
	PrimarySmartContractsManager *smartcontracts.Manager
	PrimarySmartContracts        *veganetworksmartcontracts.VegaNetworkSmartContracts

	EVMChainConfig                 *vega.EVMChainConfig
	SecondaryEthNetwork            types.ETHNetwork
	SecondaryEthClientManager      *ethutils.EthereumClientManager
	SecondaryEthClient             *ethclient.Client
	SecondarySmartContractsManager *smartcontracts.Manager
	SecondarySmartContracts        *veganetworksmartcontracts.VegaNetworkSmartContracts

	WalletManager *wallet.WalletManager

	NodeSecretStore    secrets.NodeSecretStore
	ServiceSecretStore secrets.ServiceSecretStore

	logger *zap.Logger
}

func NewVegaNetwork(
	network string,
	dataNodeClient vegaapi.DataNodeClient,
	nodeSecretStore secrets.NodeSecretStore,
	serviceSecretStore secrets.ServiceSecretStore,
	primaryEthClientManager, secondaryEthClientManager *ethutils.EthereumClientManager,
	primarySmartContractsManager, secondarySmartContractsManager *smartcontracts.Manager,
	walletManager *wallet.WalletManager,
	logger *zap.Logger,
) (*VegaNetwork, error) {
	var (
		n = &VegaNetwork{
			Network:                        network,
			DataNodeClient:                 dataNodeClient,
			PrimaryEthClientManager:        primaryEthClientManager,
			PrimarySmartContractsManager:   primarySmartContractsManager,
			SecondaryEthClientManager:      secondaryEthClientManager,
			SecondarySmartContractsManager: secondarySmartContractsManager,
			WalletManager:                  walletManager,
			NodeSecretStore:                nodeSecretStore,
			ServiceSecretStore:             serviceSecretStore,
			logger:                         logger,
		}
		errMsg = "failed to create VegaNetwork for: %w"
		err    error
	)

	if err = n.RefreshNetworkParams(); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

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
	n.AssetMainWallet, err = n.WalletManager.GetAssetMainEthWallet(n.PrimaryEthNetwork)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	n.NetworkMainWallet, err = n.WalletManager.GetNetworkMainEthWallet(n.PrimaryEthNetwork, n.Network)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	n.VegaTokenWhale, err = n.WalletManager.GetVegaTokenWhaleVegaWallet()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	n.BotsApiToken, err = n.ServiceSecretStore.GetBotsApiToken()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return n, nil
}

func (n *VegaNetwork) EthClientForChainID(chainID string) *ethclient.Client {
	switch chainID {
	case n.PrimaryEthereumConfig.ChainId:
		return n.PrimaryEthClient
	case n.EVMChainConfig.ChainId:
		return n.SecondaryEthClient
	default:
		panic(fmt.Sprintf("no ethereum client for chain ID %q", chainID))
	}
}

func (n *VegaNetwork) SmartContractManagerForChainID(chainID string) *smartcontracts.Manager {
	switch chainID {
	case n.PrimaryEthereumConfig.ChainId:
		return n.PrimarySmartContractsManager
	case n.EVMChainConfig.ChainId:
		return n.SecondarySmartContractsManager
	default:
		panic(fmt.Sprintf("no smart contract manager for chain ID %q", chainID))
	}
}

func (n *VegaNetwork) SmartContractForChainID(chainID string) *veganetworksmartcontracts.VegaNetworkSmartContracts {
	switch chainID {
	case n.PrimaryEthereumConfig.ChainId:
		return n.PrimarySmartContracts
	case n.EVMChainConfig.ChainId:
		return n.SecondarySmartContracts
	default:
		panic(fmt.Sprintf("no smart contract for chain ID %q", chainID))
	}
}

func (n *VegaNetwork) Disconnect() {
}

func (n *VegaNetwork) RefreshNetworkParams() error {
	// Read and parse some Network Parameters
	networkParams, err := n.DataNodeClient.GetAllNetworkParameters()
	if err != nil {
		return err
	}
	n.NetworkParams = networkParams

	n.PrimaryEthereumConfig, err = n.NetworkParams.PrimaryEthereumConfig()
	if err != nil {
		return fmt.Errorf("could not retrieve primary ethereum config from network parameters: %w", err)
	}
	n.PrimaryEthNetwork, err = types.GetEthNetworkForId(n.PrimaryEthereumConfig.ChainId)
	if err != nil {
		return fmt.Errorf("could not resolve primary ethereum network name from chain ID: %w", err)
	}
	n.PrimaryEthClient, err = n.PrimaryEthClientManager.GetEthClient(n.PrimaryEthNetwork)
	if err != nil {
		return fmt.Errorf("could not create primary ethereum client: %w", err)
	}
	n.PrimarySmartContracts, err = veganetworksmartcontracts.NewVegaNetworkSmartContracts(
		n.PrimaryEthClient,
		"", // will be taken from Staking Bridge
		"", // will be taken from ERC20 Bridge
		n.PrimaryEthereumConfig.CollateralBridgeContract.Address,
		n.PrimaryEthereumConfig.MultisigControlContract.Address,
		n.PrimaryEthereumConfig.StakingBridgeContract.Address,
		n.logger.Named("primary-smart-contracts"),
	)
	if err != nil {
		return fmt.Errorf("could not create primary smart contract connector: %w", err)
	}

	n.EVMChainConfig, err = n.NetworkParams.EVMChainConfig()
	if err != nil {
		return fmt.Errorf("could not retrieve secondary ethereum config from network parameters: %w", err)
	}
	n.SecondaryEthNetwork, err = types.GetEthNetworkForId(n.EVMChainConfig.ChainId)
	if err != nil {
		return fmt.Errorf("could not resolve secondary ethereum network name from chain ID: %w", err)
	}
	n.SecondaryEthClient, err = n.SecondaryEthClientManager.GetEthClient(n.SecondaryEthNetwork)
	if err != nil {
		return fmt.Errorf("could not create secondary ethereum client: %w", err)
	}
	n.SecondarySmartContracts, err = veganetworksmartcontracts.NewVegaNetworkSmartContracts(
		n.SecondaryEthClient,
		"", // will be taken from Staking Bridge
		"", // will be taken from ERC20 Bridge
		n.EVMChainConfig.CollateralBridgeContract.Address,
		n.EVMChainConfig.MultisigControlContract.Address,
		"",
		n.logger.Named("secondary-smart-contracts"),
	)
	if err != nil {
		return fmt.Errorf("could not create secondary smart contract connector: %w", err)
	}
	return nil
}
