package erc20assetpool

import (
	"fmt"
)

type ERC20AssetPoolVersion string

const (
	V1 ERC20AssetPoolVersion = "v1"
)

func (n ERC20AssetPoolVersion) IsValid() error {
	switch n {
	case V1:
		return nil
	}
	return fmt.Errorf("Invalid ERC20 Token Version %s", n)
}
