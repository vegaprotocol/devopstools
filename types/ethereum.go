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
	ETHLocal   ETHNetwork = "local"
)

func (n ETHNetwork) IsValid() error {
	switch n {
	case ETHMainnet, ETHSepolia, ETHGoerli, ETHRopsten, ETHLocal:
		return nil
	}
	return fmt.Errorf("Invalid Ethereum network %s", n)
}
