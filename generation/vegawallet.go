package generation

import (
	"fmt"

	"code.vegaprotocol.io/vega/wallet/wallet"
)

type VegaWalletPrivate struct {
	Id             string `json:"id"`
	PublicKey      string `json:"public_key"`
	PrivateKey     string `json:"private_key"`
	RecoveryPhrase string `json:"recovery_phrase"`
}

func GenerateVegaWallet() (*VegaWalletPrivate, error) {
	vegaWallet, recoveryPhrase, err := wallet.NewHDWallet("")
	if err != nil {
		return nil, fmt.Errorf("failed to generate vegawallet recovery phrase %w", err)
	}
	id := vegaWallet.ID()
	keyPair, err := vegaWallet.GenerateKeyPair(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair %w", err)
	}
	return &VegaWalletPrivate{
		Id:             id,
		PublicKey:      keyPair.PublicKey(),
		PrivateKey:     keyPair.PrivateKey(),
		RecoveryPhrase: recoveryPhrase,
	}, nil
}
