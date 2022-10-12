package generate

import (
	"github.com/vegaprotocol/devopstools/secrets"
)

func GenerateVegaNodeSecrets() (*secrets.VegaNodePrivate, error) {

	ethereumWallet, err := GenerateNewEthereumWallet()
	if err != nil {
		return nil, err
	}
	vegaWallet, err := GenerateVegaWallet()
	if err != nil {
		return nil, err
	}
	tendermintNodeKeys := GenerateTendermintKeys()
	tendermintValidatorKeys := GenerateTendermintKeys()

	walletBinaryPassphrase, err := GeneratePassword()
	if err != nil {
		return nil, err
	}
	binaryWallets, err := CreateBinaryWallets(
		tendermintValidatorKeys.PublicKey,
		vegaWallet.RecoveryPhrase,
		ethereumWallet.PrivateKey,
		walletBinaryPassphrase,
	)
	if err != nil {
		return nil, err
	}

	newNodeSecrets := &secrets.VegaNodePrivate{
		EthereumAddress:               ethereumWallet.Address,
		EthereumPrivateKey:            ethereumWallet.PrivateKey,
		EthereumMnemonic:              ethereumWallet.Mnemonic,
		VegaId:                        vegaWallet.Id,
		VegaPubKey:                    vegaWallet.PublicKey,
		VegaPrivateKey:                vegaWallet.PrivateKey,
		VegaRecoveryPhrase:            vegaWallet.RecoveryPhrase,
		TendermintNodeId:              tendermintNodeKeys.Address,
		TendermintNodePubKey:          tendermintNodeKeys.PublicKey,
		TendermintNodePrivateKey:      tendermintNodeKeys.PrivateKey,
		TendermintValidatorAddress:    tendermintValidatorKeys.Address,
		TendermintValidatorPubKey:     tendermintValidatorKeys.PublicKey,
		TendermintValidatorPrivateKey: tendermintValidatorKeys.PrivateKey,
		WalletBinaryPassphrase:        walletBinaryPassphrase,
		BinaryWallets:                 binaryWallets,
	}

	return newNodeSecrets, nil
}
