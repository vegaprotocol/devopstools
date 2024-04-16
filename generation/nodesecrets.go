package generation

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/config"
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

	walletBinaryPassphrase, err := Password()
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

func GenerateVegaNodeSecrets2() (config.Node, error) {
	scrt, err := GenerateVegaNodeSecrets()
	if err != nil {
		return config.Node{}, err
	}

	return config.Node{
		Metadata: config.NodeMetadata{
			Name:      scrt.Name,
			Country:   scrt.Country,
			InfoURL:   scrt.InfoURL,
			AvatarURL: scrt.AvatarURL,
		},
		Secrets: config.NodeSecrets{
			EthereumAddress:               scrt.EthereumAddress,
			EthereumPrivateKey:            scrt.EthereumPrivateKey,
			EthereumMnemonic:              scrt.EthereumMnemonic,
			VegaId:                        scrt.VegaId,
			VegaPubKey:                    scrt.VegaPubKey,
			VegaPrivateKey:                scrt.VegaPrivateKey,
			VegaRecoveryPhrase:            scrt.VegaRecoveryPhrase,
			VegaPubKeyIndex:               scrt.VegaPubKeyIndex,
			DeHistoryPeerId:               scrt.DeHistoryPeerId,
			DeHistoryPrivateKey:           scrt.DeHistoryPrivateKey,
			NetworkHistoryPeerId:          scrt.NetworkHistoryPeerId,
			NetworkHistoryPrivateKey:      scrt.NetworkHistoryPrivateKey,
			TendermintNodeId:              scrt.TendermintNodeId,
			TendermintNodePubKey:          scrt.TendermintNodePubKey,
			TendermintNodePrivateKey:      scrt.TendermintNodePrivateKey,
			TendermintValidatorAddress:    scrt.TendermintValidatorAddress,
			TendermintValidatorPubKey:     scrt.TendermintValidatorPubKey,
			TendermintValidatorPrivateKey: scrt.TendermintValidatorPrivateKey,
			WalletBinaryPassphrase:        scrt.WalletBinaryPassphrase,
			BinaryWallets: &config.BinaryWallets{
				NodewalletPath:       scrt.BinaryWallets.NodewalletPath,
				NodewalletBase64:     scrt.BinaryWallets.NodewalletBase64,
				VegaWalletPath:       scrt.BinaryWallets.VegaWalletPath,
				VegaWalletBase64:     scrt.BinaryWallets.VegaWalletBase64,
				EthereumWalletPath:   scrt.BinaryWallets.EthereumWalletPath,
				EthereumWalletBase64: scrt.BinaryWallets.EthereumWalletBase64,
			},
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
