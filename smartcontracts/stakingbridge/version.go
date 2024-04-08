package stakingbridge

import (
	"fmt"
)

type Version string

const (
	V1 Version = "v1"
)

func (n Version) IsValid() error {
	switch n {
	case V1:
		return nil
	}
	return fmt.Errorf("invalid Staking Bridge Version %s", n)
}
