package generate

import (
	"encoding/hex"
	"fmt"

	"github.com/vegaprotocol/devopstools/secrets"

	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func GenerateNewEthereumWallet() (*secrets.EthereumWalletPrivate, error) {
	return GenerateEthereumWallet("", "", "")
}

func GenerateEthereumWallet(
	mnemonic string,
	seed string,
	password string,
) (*secrets.EthereumWalletPrivate, error) {
	var (
		byteSeed []byte
		err      error
	)
	if mnemonic != "" {
		if !bip39.IsMnemonicValid(mnemonic) {
			return nil, fmt.Errorf("invalid mnmenoic")
		}
	} else {
		entropy, err := bip39.NewEntropy(128)
		if err != nil {
			return nil, fmt.Errorf("failed to get entropy for Ethereum key %w", err)
		}
		mnemonic, err = bip39.NewMnemonic(entropy)
		if err != nil {
			return nil, fmt.Errorf("failed to generate mnemonic for Ethereum key %w", err)
		}
	}
	if seed == "" {
		byteSeed = bip39.NewSeed(mnemonic, password)
		seed = hex.EncodeToString(byteSeed)
	} else {
		byteSeed, err = hex.DecodeString(seed)
		if err != nil {
			return nil, fmt.Errorf("failed to decode seed to byte seed %w", err)
		}
	}
	wallet, err := hdwallet.NewFromSeed(byteSeed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate wallet from seed for Ethereum key %w", err)
	}
	// zero address
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, fmt.Errorf("failed to generate account from wallet for Ethereum key %w", err)
	}
	address := account.Address.Hex()
	privKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get priv key for Ethereum key %w", err)
	}

	return &secrets.EthereumWalletPrivate{
		Address:    address,
		Mnemonic:   mnemonic,
		Seed:       seed,
		PrivateKey: privKey,
	}, nil
}

func ValidateEthereumAddressWithPrivateKey(hexAddress string, hexPrivateKey string) error {
	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to check if '%s' address matches private key, %w", hexAddress, err)
	}

	if hexAddress != crypto.PubkeyToAddress(privateKey.PublicKey).Hex() {
		return fmt.Errorf("ethereum address '%s' does not match its private key", hexAddress)
	}
	return nil
}
