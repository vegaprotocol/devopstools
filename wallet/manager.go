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
	"github.com/vegaprotocol/devopstools/vega"

	"code.vegaprotocol.io/vega/wallet/wallet"
)

type Manager struct {
	ethClientManager  *ethutils.EthereumClientManager
	walletSecretStore secrets.WalletSecretStore
}

func NewWalletManager(
	ethClientManager *ethutils.EthereumClientManager,
	walletSecretStore secrets.WalletSecretStore,
) *Manager {
	return &Manager{
		ethClientManager:  ethClientManager,
		walletSecretStore: walletSecretStore,
	}
}

func (wm *Manager) GetNetworkMainEthWallet(
	ethNetwork types.ETHNetwork,
	vegaNetwork string,
) (*ethereum.Wallet, error) {
	var (
		secretPath = fmt.Sprintf("%s/main", vegaNetwork)
		errMsg     = "failed to get Main Ethereum Wallet for %s network, %w"
	)
	ethWallet, err := wm.getEthereumWallet(ethNetwork, secretPath)
	if err != nil {
		return nil, fmt.Errorf(errMsg, vegaNetwork, err)
	}
	return ethWallet, nil
}

func (wm *Manager) GetAssetMainEthWallet(ethNetwork types.ETHNetwork) (*ethereum.Wallet, error) {
	wallet, err := wm.getEthereumWallet(ethNetwork, "AssetMain")
	if err != nil {
		return nil, fmt.Errorf("failed to get Asset Main Ethereum Wallet, %w", err)
	}
	return wallet, nil
}

func (wm *Manager) GetEthWhaleWallet(ethNetwork types.ETHNetwork) (*ethereum.Wallet, error) {
	wallet, err := wm.getEthereumWallet(ethNetwork, "EthWhale")
	if err != nil {
		return nil, fmt.Errorf("failed to get Ethereum Whale Wallet, %w", err)
	}
	return wallet, nil
}

func (wm *Manager) getEthereumWallet(
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

func (wm *Manager) GetVegaTokenWhaleVegaWallet() (wallet.Wallet, error) {
	return wm.getVegaWallet("vegaTokenWhale")
}

func (wm *Manager) getVegaWallet(secretPath string) (wallet.Wallet, error) {
	walletPrivate, err := wm.walletSecretStore.GetVegaWallet(secretPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get VegaWallet %s, get secret failed, %w", secretPath, err)
	}

	w, err := vega.LoadWallet(walletPrivate.Id, walletPrivate.RecoveryPhrase)
	if err != nil {
		return nil, fmt.Errorf("could not load vega wallet %q: %w", secretPath, err)
	}

	if err := vega.GenerateKeysUpToKey(w, walletPrivate.PublicKey); err != nil {
		return nil, fmt.Errorf("could not generate key %q: %w", secretPath, err)
	}

	return w, nil
}
