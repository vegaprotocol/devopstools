package stakingbridge

import (
	"math/big"

	StakingBridge_V1 "github.com/vegaprotocol/devopstools/smartcontracts/stakingbridge/v1"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type StakingBridgeCommon interface {
	StakeBalance(opts *bind.CallOpts, target common.Address, vega_public_key [32]byte) (*big.Int, error)
	StakingToken(opts *bind.CallOpts) (common.Address, error)
	TotalStaked(opts *bind.CallOpts) (*big.Int, error)
	RemoveStake(opts *bind.TransactOpts, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error)
	Stake(opts *bind.TransactOpts, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error)
	TransferStake(opts *bind.TransactOpts, amount *big.Int, new_address common.Address, vega_public_key [32]byte) (*types.Transaction, error)
}

type StakingBridge struct {
	StakingBridgeCommon
	Address common.Address
	Version StakingBridgeVersion
	client  *ethclient.Client

	// Minimal implementation
	v1 *StakingBridge_V1.StakingBridge
}

func NewStakingBridge(
	ethClient *ethclient.Client,
	hexAddress string,
	version StakingBridgeVersion,
) (*StakingBridge, error) {
	var err error
	result := &StakingBridge{
		Address: common.HexToAddress(hexAddress),
		Version: version,
		client:  ethClient,
	}
	switch version {
	case StakingBridgeV1:
		result.v1, err = StakingBridge_V1.NewStakingBridge(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.StakingBridgeCommon = result.v1
	}

	return result, nil
}
