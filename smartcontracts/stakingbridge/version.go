package stakingbridge

import (
	"fmt"
)

type StakingBridgeVersion string

const (
	StakingBridgeV1 StakingBridgeVersion = "v1"
)

func (n StakingBridgeVersion) IsValid() error {
	switch n {
	case StakingBridgeV1:
		return nil
	}
	return fmt.Errorf("Invalid Staking Bridge Version %s", n)
}
