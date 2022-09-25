package erc20assetpool

import (
	"fmt"
)

type ERC20AssetPoolVersion string

const (
	ERC20AssetPoolV1 ERC20AssetPoolVersion = "v1"
)

func (n ERC20AssetPoolVersion) IsValid() error {
	switch n {
	case ERC20AssetPoolV1:
		return nil
	}
	return fmt.Errorf("Invalid ERC20 Token Version %s", n)
}
