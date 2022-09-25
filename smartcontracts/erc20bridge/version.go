package erc20bridge

import (
	"fmt"
)

type ERC20BridgeVersion string

const (
	ERC20BridgeV1 ERC20BridgeVersion = "v1"
	ERC20BridgeV2 ERC20BridgeVersion = "v2"
)

func (n ERC20BridgeVersion) IsValid() error {
	switch n {
	case ERC20BridgeV1, ERC20BridgeV2:
		return nil
	}
	return fmt.Errorf("Invalid ERC20 Bridge Version %s", n)
}
