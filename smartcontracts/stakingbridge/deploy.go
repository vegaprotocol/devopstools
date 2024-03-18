package stakingbridge

import (
	"fmt"

	StakingBridge_V1 "github.com/vegaprotocol/devopstools/smartcontracts/stakingbridge/v1"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

func DeployStakingBridge(
	version StakingBridgeVersion,
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	vegaTokenAddress common.Address,
) (address common.Address, tx *ethTypes.Transaction, err error) {
	switch version {
	case StakingBridgeV1:
		address, tx, _, err = StakingBridge_V1.DeployStakingBridge(auth, backend, vegaTokenAddress)
	default:
		err = fmt.Errorf("Invalid Staking Bridge Version %s", version)
	}
	return
}
