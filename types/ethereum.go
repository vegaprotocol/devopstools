package types

import (
	"fmt"
)

type ETHNetwork string

const (
	ETHMainnet ETHNetwork = "mainnet"
	ETHSepolia ETHNetwork = "sepolia"
	ETHGoerli  ETHNetwork = "goerli"
	ETHRopsten ETHNetwork = "ropsten"
)

func (n ETHNetwork) IsValid() error {
	switch n {
	case ETHMainnet, ETHSepolia, ETHGoerli, ETHRopsten:
		return nil
	}
	return fmt.Errorf("invalid Ethereum network %s", n)
}

func GetEthNetworkForId(chainId string) (ETHNetwork, error) {
	switch chainId {
	case "1":
		return ETHMainnet, nil
	case "3":
		return ETHRopsten, nil
	case "5":
		return ETHGoerli, nil
	case "11155111":
		return ETHSepolia, nil
	}
	return "", fmt.Errorf("unknown Ethereum chain id: %s", chainId)
}

type ETHBridge string

const (
	PrimaryBridge   ETHBridge = "primary-bridge"
	SecondaryBridge ETHBridge = "secondary-bridge"
)

func (b ETHBridge) String() string {
	return string(b)
}
