package wallet

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/secrets"

	"code.vegaprotocol.io/vega/wallet/wallet"
)

var (
	walletSingletons = map[string]struct {
		HDWallet *wallet.HDWallet
		KeyPair  map[uint32]*VegaWallet
	}{}

	walletByPubKey = map[string]*VegaWallet{}
)

func GetVegaWalletSingleton(
	recoveryPhrase string,
	index uint32,
) (*VegaWallet, error) {
	errMsg := fmt.Sprintf("failed to generate key pair for wallet no %d at index %d", len(walletSingletons)+1, index)
	//
	// Create HDWallet when needed
	//
	if _, ok := walletSingletons[recoveryPhrase]; !ok {
		name := fmt.Sprintf("wallet no %d", len(walletSingletons)+1)
		hdWallet, err := wallet.ImportHDWallet(name, recoveryPhrase, wallet.LatestVersion)
		if err != nil {
			return nil, fmt.Errorf("%s, %w", errMsg, err)
		}
		walletSingletons[recoveryPhrase] = struct {
			HDWallet *wallet.HDWallet
			KeyPair  map[uint32]*VegaWallet
		}{
			HDWallet: hdWallet,
			KeyPair:  map[uint32]*VegaWallet{},
		}
	}
	wallet := walletSingletons[recoveryPhrase]
	//
	// Generate Key Pair when needed
	//
	if _, ok := wallet.KeyPair[index]; !ok {
		for {
			nextKeyPair, err := wallet.HDWallet.GenerateKeyPair(nil)
			if err != nil {
				return nil, fmt.Errorf("%s, %w", errMsg, err)
			}
			newIndex := nextKeyPair.Index()
			if newIndex > index {
				return nil, fmt.Errorf("%s, can't generate wallet at index %d, because HDWallet already generated wallet at index %d", errMsg, index, newIndex)
			}
			if _, keyPairForNewIndexAlreadyExists := wallet.KeyPair[index]; keyPairForNewIndexAlreadyExists {
				return nil, fmt.Errorf("%s, double generated key pair for index %d", errMsg, newIndex)
			}
			newWallet := &VegaWallet{
				VegaWalletPrivate: &secrets.VegaWalletPrivate{
					Id:             wallet.HDWallet.ID(),
					PublicKey:      nextKeyPair.PublicKey(),
					PrivateKey:     nextKeyPair.PrivateKey(),
					RecoveryPhrase: recoveryPhrase,
				},
				hdWallet: wallet.HDWallet,
				keyPair:  nextKeyPair,
			}
			walletSingletons[recoveryPhrase].KeyPair[newIndex] = newWallet
			walletByPubKey[newWallet.PublicKey] = newWallet
			if newIndex == index {
				break
			}
		}
	}
	return walletSingletons[recoveryPhrase].KeyPair[index], nil
}
