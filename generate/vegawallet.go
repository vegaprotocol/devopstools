package generate

import (
	"fmt"

	"code.vegaprotocol.io/vega/wallet/wallet"
	"github.com/vegaprotocol/devopstools/secrets"
)

func GenerateVegaWallet() (*secrets.VegaWalletPrivate, error) {
	vegaWallet, recoveryPhrase, err := wallet.NewHDWallet("")
	if err != nil {
		return nil, fmt.Errorf("failed to generate vegawallet recovery phrase %w", err)
	}
	id := vegaWallet.ID()
	keyPair, err := vegaWallet.GenerateKeyPair(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair %w", err)
	}
	return &secrets.VegaWalletPrivate{
		Id:             id,
		PublicKey:      keyPair.PublicKey(),
		PrivateKey:     keyPair.PrivateKey(),
		RecoveryPhrase: recoveryPhrase,
	}, nil
}

func ValidateVegawalletIdAndPubKeyAndPrivKeyWithRecoveryPhrase(
	id string, publicKey string, privateKey string, recoveryPhrase string,
) error {
	vegaWallet, err := wallet.ImportHDWallet("my wallet", recoveryPhrase, wallet.LatestVersion)

	if err != nil {
		return fmt.Errorf("failed to get vegawallet with recovery phrase %w", err)
	}
	expectedId := vegaWallet.ID()
	keyPair, err := vegaWallet.GenerateKeyPair(nil)
	if err != nil {
		return fmt.Errorf("failed to get key pair for wallet from recovery phrase %w", err)
	}
	if id != expectedId {
		return fmt.Errorf("vegawallet data does not match recovery phrase, id does not match, provided='%s', expected='%s'", id, expectedId)
	}
	if publicKey != keyPair.PublicKey() {
		return fmt.Errorf("vegawallet data does not match recovery phrase, public key does not match, provided='%s', expected='%s'", publicKey, keyPair.PublicKey())
	}
	if privateKey != keyPair.PrivateKey() {
		return fmt.Errorf("vegawallet data does not match recovery phrase, private key does not match for wallet with id='%s'", id)
	}
	return nil
}
