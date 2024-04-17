package ethutils

import (
	"fmt"
	"os"
	"path/filepath"
)

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

func _(
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

	if err = os.WriteFile(downloadURLFile, []byte(data.DownloadURL), 0o644); err != nil {
		return fmt.Errorf("failed to write Download URL to file %s for %s. %w", downloadURLFile, name, err)
	}

	// Save Source Code
	for sourceCodeName, sourceCode := range data.SourceCode {
		sourceCodeFile := filepath.Join(dir, sourceCodeName)
		if err = os.WriteFile(sourceCodeFile, []byte(sourceCode), 0o644); err != nil {
			return fmt.Errorf("failed to write Source Code to file %s for %s. %w", sourceCodeFile, name, err)
		}
		fmt.Printf("- %s: updated Source Code: %s\n", name, sourceCodeFile)
	}

	// Save Byte Code
	if err = os.WriteFile(binaryFile, []byte(data.CreationHexByteCode[2:]), 0o644); err != nil {
		return fmt.Errorf("failed to write Byte Code to file %s for %s. %w", binaryFile, name, err)
	}
	fmt.Printf("- %s: updated Binary Code: %s\n", name, binaryFile)

	// Save Byte Code Hash
	if err = os.WriteFile(hashFile, []byte(data.ByteCodeHash), 0o644); err != nil {
		return fmt.Errorf("failed to write Byte Code Hash to file %s for %s. %w", hashFile, name, err)
	}
	fmt.Printf("- %s: updated Hash: %s\n", name, hashFile)

	// Save ABI
	if err = os.WriteFile(abiFile, []byte(data.ABI), 0o644); err != nil {
		return fmt.Errorf("failed to write ABI to file %s for %s. %w", abiFile, name, err)
	}
	fmt.Printf("- %s: updated ABI: %s\n", name, abiFile)

	// Save Go Bindings
	if err = os.WriteFile(goBindingsFile, []byte(data.GoBindings), 0o644); err != nil {
		return fmt.Errorf("failed to write Go Bindings to file %s for %s. %w", goBindingsFile, name, err)
	}
	fmt.Printf("- %s: updated Go Bindings: %s\n", name, goBindingsFile)

	return nil
}
