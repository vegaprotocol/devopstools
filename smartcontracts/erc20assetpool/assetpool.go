package erc20assetpool

import (
	"math/big"

	ERC20AssetPool_V1 "github.com/vegaprotocol/devopstools/smartcontracts/erc20assetpool/v1"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20AssetPoolCommon interface {
	Erc20BridgeAddress(opts *bind.CallOpts) (common.Address, error)
	MultisigControlAddress(opts *bind.CallOpts) (common.Address, error)
	SetBridgeAddress(opts *bind.TransactOpts, new_address common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error)
	SetMultisigControl(opts *bind.TransactOpts, new_address common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error)
	Withdraw(opts *bind.TransactOpts, token_address common.Address, target common.Address, amount *big.Int) (*types.Transaction, error)
}

type ERC20AssetPool struct {
	ERC20AssetPoolCommon
	Address common.Address
	Version ERC20AssetPoolVersion
	client  *ethclient.Client

	// Minimal implementation
	v1 *ERC20AssetPool_V1.ERC20AssetPool
}

func NewERC20AssetPool(
	ethClient *ethclient.Client,
	hexAddress string,
	version ERC20AssetPoolVersion,
) (*ERC20AssetPool, error) {
	var err error
	result := &ERC20AssetPool{
		Address: common.HexToAddress(hexAddress),
		Version: version,
		client:  ethClient,
	}
	switch version {
	case ERC20AssetPoolV1:
		result.v1, err = ERC20AssetPool_V1.NewERC20AssetPool(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.ERC20AssetPoolCommon = result.v1
	}

	return result, nil
}
