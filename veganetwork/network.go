package veganetwork

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/secrets"
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

	VegaTokenWhale *wallet.VegaWallet

	// network params/config
	NetworkParams *types.NetworkParams

	// clients
	DataNodeClient vegaapi.DataNodeClient

	PrimaryEthereumConfig   *vega.EthereumConfig
	PrimaryEthNetwork       types.ETHNetwork
	PrimaryEthClientManager *ethutils.EthereumClientManager
	PrimaryEthClient        *ethclient.Client
	PrimarySmartContracts   *veganetworksmartcontracts.VegaNetworkSmartContracts

	WalletManager *wallet.Manager

	logger *zap.Logger
}

func NewVegaNetwork(network string, dataNodeClient vegaapi.DataNodeClient, nodeSecretStore secrets.NodeSecretStore, primaryEthClientManager *ethutils.EthereumClientManager, walletManager *wallet.Manager, logger *zap.Logger) (*VegaNetwork, error) {
	var (
		n = &VegaNetwork{
			Network:                 network,
			DataNodeClient:          dataNodeClient,
			PrimaryEthClientManager: primaryEthClientManager,
			WalletManager:           walletManager,
			logger:                  logger,
		}
		errMsg = "failed to create VegaNetwork for: %w"
		err    error
	)

	if err = n.RefreshNetworkParams(); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	// Node Secrets
	n.NodeSecrets, err = nodeSecretStore.GetAllVegaNode(n.Network)
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
	n.NetworkMainWallet, err = n.WalletManager.GetNetworkMainEthWallet(n.PrimaryEthNetwork, n.Network)
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

	return nil
}
