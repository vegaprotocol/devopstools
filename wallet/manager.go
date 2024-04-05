package wallet

import (
	context2 "context"
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/ethereum"
	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/secrets"
	"github.com/vegaprotocol/devopstools/tools"
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

//
// ETHEREUM
//

func (wm *WalletManager) GetNetworkMainEthWallet(
	ethNetwork types.ETHNetwork,
	vegaNetwork string,
) (*ethereum.Wallet, error) {
	var (
		secretPath string = fmt.Sprintf("%s/main", vegaNetwork)
		errMsg            = "failed to get Main Ethereum Wallet for %s network, %w"
	)
	ethWallet, err := wm.getEthereumWallet(ethNetwork, secretPath)
	if err != nil {
		return nil, fmt.Errorf(errMsg, vegaNetwork, err)
	}
	return ethWallet, nil
}

func (wm *WalletManager) GetAssetMainEthWallet(ethNetwork types.ETHNetwork) (*ethereum.Wallet, error) {
	wallet, err := wm.getEthereumWallet(ethNetwork, "AssetMain")
	if err != nil {
		return nil, fmt.Errorf("failed to get Asset Main Ethereum Wallet, %w", err)
	}
	return wallet, nil
}

func (wm *WalletManager) GetEthWhaleWallet(ethNetwork types.ETHNetwork) (*ethereum.Wallet, error) {
	wallet, err := wm.getEthereumWallet(ethNetwork, "EthWhale")
	if err != nil {
		return nil, fmt.Errorf("failed to get Ethereum Whale Wallet, %w", err)
	}
	return wallet, nil
}

func (wm *WalletManager) getEthereumWallet(
	ethNetwork types.ETHNetwork,
	secretPath string,
) (*ethereum.Wallet, error) {
	walletPrivate, err := wm.walletSecretStore.GetEthereumWallet(secretPath)
	if err != nil {
		return nil, err
	}
	ethClient, err := wm.ethClientManager.GetEthClient(ethNetwork)
	if err != nil {
		return nil, err
	}

	ethWallet, err := tools.RetryReturn(6, 10*time.Second, func() (*ethereum.Wallet, error) {
		ethWallet, err := ethereum.NewWallet(context2.Background(), ethClient, walletPrivate.PrivateKey)
		if err != nil {
			return nil, err
		}
		return ethWallet, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new eth wallet: %w", err)
	}
	return ethWallet, nil
}

//
// VEGAWALLET
//

func (wm *WalletManager) GetVegaTokenWhaleVegaWallet() (*VegaWallet, error) {
	return wm.getVegaWallet("vegaTokenWhale")
}

func (wm *WalletManager) getVegaWallet(secretPath string) (*VegaWallet, error) {
	walletPrivate, err := wm.walletSecretStore.GetVegaWallet(secretPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get VegaWallet %s, get secret failed, %w", secretPath, err)
	}
	vegawallet, err := NewVegaWallet(walletPrivate)
	if err != nil {
		return nil, fmt.Errorf("failed to get VegaWallet %s, setting up secret failed, %w", secretPath, err)
	}
	return vegawallet, nil
}
