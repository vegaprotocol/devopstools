package erc20token

import (
	"fmt"
)

type ERC20TokenVersion string

const (
	ERC20TokenBase    ERC20TokenVersion = "TokenBase"
	ERC20TokenOld     ERC20TokenVersion = "TokenOld" // deprecated - don't ever use, remove if you can
	ERC20TokenOther   ERC20TokenVersion = "TokenOther"
	ERC20TokenMinimal ERC20TokenVersion = "TokenMinimal"
)

func (n ERC20TokenVersion) IsValid() error {
	switch n {
	case ERC20TokenBase, ERC20TokenOld, ERC20TokenOther, ERC20TokenMinimal:
		return nil
	}
	return fmt.Errorf("invalid ERC20 Token Version %s", n)
}
