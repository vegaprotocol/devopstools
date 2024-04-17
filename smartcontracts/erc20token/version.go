package erc20token

import (
	"fmt"
)

type Version string

const (
	Base    Version = "TokenBase"
	Old     Version = "TokenOld" // deprecated - don't ever use, remove if you can
	Other   Version = "TokenOther"
	Minimal Version = "TokenMinimal"
)

func (n Version) IsValid() error {
	switch n {
	case Base, Old, Other, Minimal:
		return nil
	}
	return fmt.Errorf("invalid ERC20 Token Version %s", n)
}
