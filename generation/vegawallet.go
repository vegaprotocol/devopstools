package generation

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/secrets"

	"code.vegaprotocol.io/vega/wallet/wallet"
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
