package stakingbridge

import (
	"fmt"
)

type StakingBridgeVersion string

const (
	V1 StakingBridgeVersion = "v1"
)

func (n StakingBridgeVersion) IsValid() error {
	switch n {
	case V1:
		return nil
	}
	return fmt.Errorf("Invalid Staking Bridge Version %s", n)
}
