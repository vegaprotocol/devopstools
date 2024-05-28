package generation

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/config"
)

func GenerateVegaNodeSecrets() (config.Node, error) {
	metadata, err := GenerateNodeMetadata()
	if err != nil {
		return config.Node{}, err
	}
	ethereumWallet, err := GenerateNewEthereumWallet()
	if err != nil {
		return config.Node{}, err
	}
	vegaWallet, err := GenerateVegaWallet()
	if err != nil {
		return config.Node{}, err
	}

	vegaPubKeyIndex := uint64(1)
	deHistory, err := GenerateDeHistoryIdentity(vegaWallet.RecoveryPhrase)
	if err != nil {
		return config.Node{}, err
	}
	tendermintNodeKeys := GenerateTendermintKeys()
	tendermintValidatorKeys := GenerateTendermintKeys()

	walletBinaryPassphrase, err := Password()
	if err != nil {
		return config.Node{}, err
	}
	binaryWallets, err := CreateBinaryWallets(
		tendermintValidatorKeys.PublicKey,
		vegaWallet.RecoveryPhrase,
		ethereumWallet.PrivateKey,
		walletBinaryPassphrase,
	)
	if err != nil {
		return config.Node{}, err
	}

	return config.Node{
		Metadata: config.NodeMetadata{
			Name:      metadata.Name,
			Country:   metadata.Country,
			InfoURL:   metadata.InfoURL,
			AvatarURL: metadata.AvatarURL,
		},
		Secrets: config.NodeSecrets{
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
		},
		API: config.NodeAPI{},
	}, nil
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
