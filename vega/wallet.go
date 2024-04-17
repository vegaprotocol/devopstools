package vega

import (
	"fmt"

	"code.vegaprotocol.io/vega/wallet/wallet"
)

var walletSingletons = map[string]wallet.Wallet{}

func LoadWallet(name, recoveryPhrase string) (wallet.Wallet, error) {
	w, err := wallet.ImportHDWallet(name, recoveryPhrase, wallet.Version2)
	if err != nil {
		return nil, fmt.Errorf("could not import HD wallet: %w", err)
	}
	if _, err := w.GenerateKeyPair(nil); err != nil {
		return nil, fmt.Errorf("could not generate key pair: %w", err)
	}

	return w, nil
}

func LoadWalletAsSingleton(recoveryPhrase string) (wallet.Wallet, error) {
	if _, ok := walletSingletons[recoveryPhrase]; !ok {
		name := fmt.Sprintf("wallet no %d", len(walletSingletons)+1)
		hdWallet, err := wallet.ImportHDWallet(name, recoveryPhrase, wallet.LatestVersion)
		if err != nil {
			return nil, fmt.Errorf("could not import HD wallet %q: %w", name, err)
		}

		walletSingletons[recoveryPhrase] = hdWallet
	}

	return walletSingletons[recoveryPhrase], nil
}

func GenerateKeysUpToIndex(w wallet.Wallet, index uint32) error {
	keyPairs := w.ListKeyPairs()
	maxIteration := int(index) - len(keyPairs)

	for i := 0; i < maxIteration; i++ {
		_, err := w.GenerateKeyPair(nil)
		if err != nil {
			return fmt.Errorf("could not generate key pair for index %d: %w", index, err)
		}
	}

	return nil
}

func GenerateKeysUpToKey(w wallet.Wallet, pubKey string) error {
	for _, keyPair := range w.ListKeyPairs() {
		if keyPair.PublicKey() == pubKey {
			return nil
		}
	}

	maxIteration := 20
	for i := 0; i < maxIteration; i++ {
		keyPair, err := w.GenerateKeyPair(nil)
		if err != nil {
			return fmt.Errorf("could not generate key pair for index: %w", err)
		}
		if keyPair.PublicKey() == pubKey {
			return nil
		}
	}

	return fmt.Errorf("could not generate key pair for %d keys", maxIteration)
}
