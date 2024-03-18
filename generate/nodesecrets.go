package generate

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/secrets"
)

func GenerateVegaNodeSecrets() (*secrets.VegaNodePrivate, error) {
	metadata, err := GenerateNodeMetadata()
	if err != nil {
		return nil, err
	}

	ethereumWallet, err := GenerateNewEthereumWallet()
	if err != nil {
		return nil, err
	}
	vegaWallet, err := GenerateVegaWallet()
	if err != nil {
		return nil, err
	}
	vegaPubKeyIndex := uint64(1)
	deHistory, err := GenerateDeHistoryIdentity(vegaWallet.RecoveryPhrase)
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
		Name:                          metadata.Name,
		Country:                       metadata.Country,
		InfoURL:                       metadata.InfoURL,
		AvatarURL:                     metadata.AvatarURL,
		EthereumAddress:               ethereumWallet.Address,
		EthereumPrivateKey:            ethereumWallet.PrivateKey,
		EthereumMnemonic:              ethereumWallet.Mnemonic,
		VegaId:                        vegaWallet.Id,
		VegaPubKey:                    vegaWallet.PublicKey,
		VegaPrivateKey:                vegaWallet.PrivateKey,
		VegaRecoveryPhrase:            vegaWallet.RecoveryPhrase,
		VegaPubKeyIndex:               &vegaPubKeyIndex,
		DeHistoryPeerId:               deHistory.PeerID,
		DeHistoryPrivateKey:           deHistory.PrivKey,
		NetworkHistoryPeerId:          deHistory.PeerID,
		NetworkHistoryPrivateKey:      deHistory.PrivKey,
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

type VegaNodeMetadata struct {
	Name      string `json:"name"`
	Country   string `json:"country"`
	InfoURL   string `json:"info_url"`
	AvatarURL string `json:"avatar_url"`
}

func GenerateNodeMetadata() (*VegaNodeMetadata, error) {
	errMsg := "failed to generate metadata for node, %w"
	name, err := GenerateName()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	country, err := GenerateCountryCode()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	infoURL, err := GenerateRandomWikiURL()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	avatarURL, err := GenerateAvatarURL()
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	return &VegaNodeMetadata{
		Name:      name,
		Country:   country,
		InfoURL:   infoURL,
		AvatarURL: avatarURL,
	}, nil
}
