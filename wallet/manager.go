package wallet

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/types"
)

type WalletManager struct {
	ethClientManager  *ethutils.EthereumClientManager
	walletSecretStore secrets.WalletSecretStore
}

func NewWalletManager(
	ethClientManager *ethutils.EthereumClientManager,
	walletSecretStore secrets.WalletSecretStore,
) *WalletManager {
	return &WalletManager{
		ethClientManager:  ethClientManager,
		walletSecretStore: walletSecretStore,
	}
}

func (wm *WalletManager) GetNetworkMainEthWallet(
	ethNetwork types.ETHNetwork,
	vegaNetwork string,
) (*EthWallet, error) {
	var (
		secretPath string
		errMsg     = "failed to get Main Ethereum Wallet for %s network, %w"
	)
	switch vegaNetwork {
	case "fairground", "devnet", "stagnet3":
		secretPath = "OldMain"
	default:
		secretPath = fmt.Sprintf("%s/main", vegaNetwork)
	}
	ethWallet, err := wm.getEthereumWallet(ethNetwork, secretPath)
	if err != nil {
		return nil, fmt.Errorf(errMsg, vegaNetwork, err)
	}
	return ethWallet, nil
}

func (wm *WalletManager) GetAssetMainEthWallet(ethNetwork types.ETHNetwork) (*EthWallet, error) {
	wallet, err := wm.getEthereumWallet(ethNetwork, "AssetMain")
	if err != nil {
		return nil, fmt.Errorf("failed to get Asset Main Ethereum Wallet, %w", err)
	}
	return wallet, nil
}

func (wm *WalletManager) GetEthWhaleWallet(ethNetwork types.ETHNetwork) (*EthWallet, error) {
	wallet, err := wm.getEthereumWallet(ethNetwork, "EthWhale")
	if err != nil {
		return nil, fmt.Errorf("failed to get Ethereum Whale Wallet, %w", err)
	}
	return wallet, nil
}

func (wm *WalletManager) getEthereumWallet(
	ethNetwork types.ETHNetwork,
	secretPath string,
) (*EthWallet, error) {
	walletPrivate, err := wm.walletSecretStore.GetEthereumWallet(secretPath)
	if err != nil {
		return nil, err
	}
	ethClient, err := wm.ethClientManager.GetEthClient(ethNetwork)
	if err != nil {
		return nil, err
	}
	ethWallet, err := NewEthWallet(ethClient, walletPrivate)
	if err != nil {
		return nil, err
	}
	return ethWallet, nil
}
