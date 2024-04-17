package generation

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/vegaprotocol/devopstools/config"

	"code.vegaprotocol.io/vega/core/nodewallets"
	"code.vegaprotocol.io/vega/core/nodewallets/registry"
	"code.vegaprotocol.io/vega/paths"
	"code.vegaprotocol.io/vega/wallet/wallet"
	"code.vegaprotocol.io/vega/wallet/wallets"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

func CreateBinaryWallets(
	tendermintValidatorPubKey string,
	vegaWalletRecoveryPhrase string,
	ethereumPrivateKey string,
	walletBinaryPassphrase string,
) (*config.BinaryWallets, error) {
	homeDir, err := os.MkdirTemp("", "wallets")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(homeDir)

	vegaHome := path.Join(homeDir, "vega_home")
	vegaPaths := paths.New(vegaHome)

	nodewalletPath, err := createNodewallet(vegaPaths, walletBinaryPassphrase)
	if err != nil {
		return nil, err
	}

	_, err = nodewallets.ImportTendermintPubkey(
		vegaPaths, walletBinaryPassphrase, tendermintValidatorPubKey, true,
	)
	if err != nil {
		return nil, err
	}

	vegaWalletPath, err := createVegawallet(vegaWalletRecoveryPhrase, vegaHome, walletBinaryPassphrase)
	if err != nil {
		return nil, err
	}

	data, err := nodewallets.ImportVegaWallet(
		vegaPaths, walletBinaryPassphrase, walletBinaryPassphrase, vegaWalletPath, true,
	)
	if err != nil {
		return nil, err
	}

	vegaWalletPath = data["walletFilePath"]

	ethWalletPath, err := createEthereumWallet(ethereumPrivateKey, vegaHome, walletBinaryPassphrase)
	if err != nil {
		return nil, err
	}

	data, err = nodewallets.ImportEthereumWallet(
		vegaPaths, walletBinaryPassphrase, walletBinaryPassphrase,
		"", "", ethWalletPath, true,
	)
	if err != nil {
		return nil, err
	}
	ethWalletPath = data["walletFilePath"]

	nodewalletBase64, err := encodeToBase64(nodewalletPath)
	if err != nil {
		return nil, err
	}
	vegaWalletBase64, err := encodeToBase64(vegaWalletPath)
	if err != nil {
		return nil, err
	}
	ethereumWalletBase64, err := encodeToBase64(ethWalletPath)
	if err != nil {
		return nil, err
	}

	return &config.BinaryWallets{
		NodewalletPath:       strings.TrimPrefix(nodewalletPath, vegaHome)[1:],
		VegaWalletPath:       strings.TrimPrefix(vegaWalletPath, vegaHome)[1:],
		EthereumWalletPath:   strings.TrimPrefix(ethWalletPath, vegaHome)[1:],
		NodewalletBase64:     nodewalletBase64,
		VegaWalletBase64:     vegaWalletBase64,
		EthereumWalletBase64: ethereumWalletBase64,
	}, nil
}

func createNodewallet(vegaPaths paths.Paths, walletPassphrase string) (string, error) {
	nodewalletRegistry, err := registry.NewLoader(vegaPaths, walletPassphrase)
	if err != nil {
		return "", err
	}
	nodewalletPath := nodewalletRegistry.RegistryFilePath()

	return nodewalletPath, nil
}

func createVegawallet(vegaWalletRecoveryPhrase string, vegaHome string, walletPassphrase string) (string, error) {
	// create a wallet with name
	vegaWalletName := "isolatedValidatorWallet"
	vegaWallet, err := wallet.ImportHDWallet(vegaWalletName, vegaWalletRecoveryPhrase, wallet.LatestVersion)
	if err != nil {
		return "", err
	}
	_, err = vegaWallet.GenerateKeyPair(nil)
	if err != nil {
		return "", err
	}
	// create vegaWallet store and store new hdwallet
	vegaWalletStore, err := wallets.InitialiseStore(vegaHome, false)
	if err != nil {
		return "", err
	}
	if err := vegaWalletStore.CreateWallet(context.Background(), vegaWallet, walletPassphrase); err != nil {
		return "", err
	}
	if err := vegaWalletStore.UnlockWallet(context.Background(), vegaWalletName, walletPassphrase); err != nil {
		return "", err
	}
	_, err = vegaWalletStore.GetWallet(context.Background(), vegaWalletName)
	if err != nil {
		return "", err
	}
	// get path of the new wallet in store
	vegaWalletPath := vegaWalletStore.GetWalletPath(vegaWalletName)

	return vegaWalletPath, nil
}

func createEthereumWallet(
	ethereumPrivateKey string,
	vegaHome string,
	walletPassphrase string,
) (string, error) {
	ethPrivateKey, err := crypto.HexToECDSA(ethereumPrivateKey)
	if err != nil {
		return "", err
	}
	ethKeystoreHome := path.Join(vegaHome, "eth-keystore")
	ethKeystore := keystore.NewKeyStore(ethKeystoreHome, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err = ethKeystore.ImportECDSA(ethPrivateKey, walletPassphrase)
	if err != nil {
		return "", err
	}
	files, err := os.ReadDir(ethKeystoreHome)
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", fmt.Errorf("Expected to have one file in directory %s\n", ethKeystoreHome)
	}
	ethWalletOrigPath := path.Join(ethKeystoreHome, files[0].Name())
	ethWalletPath := path.Join(path.Dir(ethWalletOrigPath), "validatorEthWallet")
	if err := os.Rename(ethWalletOrigPath, ethWalletPath); err != nil {
		return "", err
	}

	return ethWalletPath, nil
}

func encodeToBase64(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	content, _ := io.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded, nil
}
