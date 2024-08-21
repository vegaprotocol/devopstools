package stakingbridge

import (
	"math/big"

	StakingBridge_V1 "github.com/vegaprotocol/devopstools/smartcontracts/stakingbridge/v1"
	StakingBridge_V2 "github.com/vegaprotocol/devopstools/smartcontracts/stakingbridge/v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Common interface {
	StakeBalance(opts *bind.CallOpts, target common.Address, vega_public_key [32]byte) (*big.Int, error)
	StakingToken(opts *bind.CallOpts) (common.Address, error)
	TotalStaked(opts *bind.CallOpts) (*big.Int, error)
	RemoveStake(opts *bind.TransactOpts, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error)
	Stake(opts *bind.TransactOpts, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error)
	TransferStake(opts *bind.TransactOpts, amount *big.Int, new_address common.Address, vega_public_key [32]byte) (*types.Transaction, error)
}

type StakingBridge struct {
	Common
	Address common.Address
	Version Version
	client  *ethclient.Client

	// Minimal implementation
	v1 *StakingBridge_V1.StakingBridge
	v2 *StakingBridge_V2.StakingBridge
}

func NewStakingBridge(
	ethClient *ethclient.Client,
	hexAddress string,
	version Version,
) (*StakingBridge, error) {
	var err error
	result := &StakingBridge{
		Address: common.HexToAddress(hexAddress),
		Version: version,
		client:  ethClient,
	}
	switch version {
	case V1:
		result.v1, err = StakingBridge_V1.NewStakingBridge(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.Common = result.v1
	case V2:
		result.v2, err = StakingBridge_V2.NewStakingBridge(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.Common = result.v2
	}

	return result, nil
}
