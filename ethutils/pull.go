package ethutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vegaprotocol/devopstools/etherscan"
	"github.com/vegaprotocol/devopstools/types"
)

func PullAndStoreSmartContractImmutableData(
	hexAddress string,
	ethNetwork types.ETHNetwork,
	name string,
	dir string,
	ethereumClientManager *EthereumClientManager,
) error {
	// Get clients
	ethClient, err := ethereumClientManager.GetEthClient(ethNetwork)
	if err != nil {
		return fmt.Errorf(
			"failed to pull and store Smart Contract %s (%s: %s), failed to get Ethereum Client, %w",
			name, ethNetwork, hexAddress, err,
		)
	}
	etherscanClient, err := ethereumClientManager.GetEtherscanClient(ethNetwork)
	if err != nil {
		return fmt.Errorf(
			"failed to pull and store Smart Contract %s (%s: %s), failed to get Etherscan Client, %w",
			name, ethNetwork, hexAddress, err,
		)
	}

	// Pull data
	data, err := PullSmartContractImmutableData(
		ethClient, etherscanClient, name, hexAddress,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to pull and store Smart Contract %s (%s: %s), failed to pull, %w",
			name, ethNetwork, hexAddress, err,
		)
	}

	// Store data
	if err = storeSmartContractImmutableData(name, *data, dir); err != nil {
		return fmt.Errorf(
			"failed to pull and store Smart Contract %s (%s: %s), failed to store, %w",
			name, ethNetwork, hexAddress, err,
		)
	}
	return nil
}

type SmartContractImmutableData struct {
	SourceCode          map[string]string
	Name                string
	ByteCode            []byte
	ByteCodeHash        string
	CreationHexByteCode string
	ABI                 string
	GoBindings          string
	DownloadURL         string
}

func PullSmartContractImmutableData(
	ethClient *ethclient.Client,
	etherscanClient *etherscan.EtherscanClient,
	name string,
	hexAddress string,
) (*SmartContractImmutableData, error) {
	var (
		address = common.HexToAddress(hexAddress)
	)

	// Fetch Byte Code
	byteCode, err := ethClient.CodeAt(context.Background(), address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get byte code for the \"%s\" smart contract: %w", hexAddress, err)
	}
	// Get Byte Code Hash
	byteCodeHash, err := GetByteCodeHash(byteCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get Byte Code Hash for %s. %w", hexAddress, err)
	}
	// Data from etherscan
	etherscanData, err := etherscanClient.GetSmartContractData(context.Background(), hexAddress)
	if err != nil {
		return nil, err
	}
	// Cleanup Creation Byte Code
	creationHexByteCode := CleanupCreationByteCode(etherscanData.CreationHexByteCode)
	// Get Go Bindings
	goBindings, err := GetGoBindings(etherscanData.ABI, name, creationHexByteCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get Go Bindings for %s. %w", hexAddress, err)
	}

	return &SmartContractImmutableData{
		SourceCode:          etherscanData.SourceCode,
		Name:                etherscanData.ContractName,
		ByteCode:            byteCode,
		ByteCodeHash:        *byteCodeHash,
		CreationHexByteCode: etherscanData.CreationHexByteCode,
		ABI:                 etherscanData.ABI,
		GoBindings:          *goBindings,
		DownloadURL:         etherscanData.DownloadURL,
	}, nil
}

func storeSmartContractImmutableData(
	name string,
	data SmartContractImmutableData,
	dir string,
) error {
	var (
		binaryFile      = filepath.Join(dir, name+".bin")
		hashFile        = filepath.Join(dir, "hash.txt")
		abiFile         = filepath.Join(dir, "abi.json")
		goBindingsFile  = filepath.Join(dir, name+".go")
		downloadURLFile = filepath.Join(dir, "download_url.txt")
		err             error
	)

	if err = os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to clean directory %s. %w", dir, err)
	}
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s. %w", dir, err)
	}

	if err = ioutil.WriteFile(downloadURLFile, []byte(data.DownloadURL), 0644); err != nil {
		return fmt.Errorf("failed to write Download URL to file %s for %s. %w", downloadURLFile, name, err)
	}

	// Save Source Code
	for sourceCodeName, sourceCode := range data.SourceCode {
		sourceCodeFile := filepath.Join(dir, sourceCodeName)
		if err = ioutil.WriteFile(sourceCodeFile, []byte(sourceCode), 0644); err != nil {
			return fmt.Errorf("failed to write Source Code to file %s for %s. %w", sourceCodeFile, name, err)
		}
		fmt.Printf("- %s: updated Source Code: %s\n", name, sourceCodeFile)
	}

	// Save Byte Code
	if err = ioutil.WriteFile(binaryFile, []byte(data.CreationHexByteCode[2:]), 0644); err != nil {
		return fmt.Errorf("failed to write Byte Code to file %s for %s. %w", binaryFile, name, err)
	}
	fmt.Printf("- %s: updated Binary Code: %s\n", name, binaryFile)

	// Save Byte Code Hash
	if err = ioutil.WriteFile(hashFile, []byte(data.ByteCodeHash), 0644); err != nil {
		return fmt.Errorf("failed to write Byte Code Hash to file %s for %s. %w", hashFile, name, err)
	}
	fmt.Printf("- %s: updated Hash: %s\n", name, hashFile)

	// Save ABI
	if err = ioutil.WriteFile(abiFile, []byte(data.ABI), 0644); err != nil {
		return fmt.Errorf("failed to write ABI to file %s for %s. %w", abiFile, name, err)
	}
	fmt.Printf("- %s: updated ABI: %s\n", name, abiFile)

	// Save Go Bindings
	if err = ioutil.WriteFile(goBindingsFile, []byte(data.GoBindings), 0644); err != nil {
		return fmt.Errorf("failed to write Go Bindings to file %s for %s. %w", goBindingsFile, name, err)
	}
	fmt.Printf("- %s: updated Go Bindings: %s\n", name, goBindingsFile)

	return nil
}
