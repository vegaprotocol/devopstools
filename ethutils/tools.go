package ethutils

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"golang.org/x/crypto/sha3"
)

//
// Helper functions
//

func GetByteCodeHash(byteCode []byte) (*string, error) {
	if len(byteCode) < 1 {
		return nil, fmt.Errorf("failed to get bytecode for smart contracts: empty bytecode")
	}

	// the bytecode of the contract is appended which is deployment specific. We only care about
	// the contract code itself and so we need to strip this meta-data before hashing it. For the version
	// of Solidity we use, the format is [contract-bytecode]a264[CBOR-encoded meta-data]
	asHex := strings.Split(hex.EncodeToString(byteCode), "a264")
	if len(asHex) != 2 {
		return nil, fmt.Errorf("unexpected format of solidity bytecode")
	}

	// Back to bytes for hashing
	smartContractByteCode, err := hex.DecodeString(asHex[0])
	if err != nil {
		return nil, err
	}

	// convert to hash
	hasher := sha3.New256()
	if _, err = hasher.Write(smartContractByteCode); err != nil {
		return nil, err
	}
	hashByte := hasher.Sum(nil)
	hash := hex.EncodeToString(hashByte)

	return &hash, nil
}

func GetGoBindings(
	abi string,
	pkgName string,
	creationHexByteCode string,
) (*string, error) {
	var (
		abis    = []string{abi}
		bins    = []string{creationHexByteCode}
		types   = []string{pkgName}
		sigs    []map[string]string
		libs    = make(map[string]string)
		aliases = make(map[string]string)
		lang    = bind.LangGo
		err     error
	)

	// Generate the contract binding
	code, err := bind.Bind(types, abis, bins, sigs, pkgName, lang, libs, aliases)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Go Binding %s: %w", pkgName, err)
	}
	return &code, nil
}

func CleanupCreationByteCode(bytecode string) string {
	// Remove constructor argument values
	// Note: go-ethereum generted go-binding Deploy function will append constructor arguments to this bytecode
	reg := regexp.MustCompile(`^(.*a2646970667358221220[0-9A-Fa-f]{64,64}64736f6c634300080[0-9A-Fa-f]0033).*$`)
	return reg.ReplaceAllString(bytecode, "${1}")
}
