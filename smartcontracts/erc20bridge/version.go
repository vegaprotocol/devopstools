package erc20bridge

import (
	"fmt"
)

type ERC20BridgeVersion string

const (
	V1 ERC20BridgeVersion = "v1"
	V2 ERC20BridgeVersion = "v2"
)

func (n ERC20BridgeVersion) IsValid() error {
	switch n {
	case V1, V2:
		return nil
	}
	return fmt.Errorf("Invalid ERC20 Bridge Version %s", n)
}
