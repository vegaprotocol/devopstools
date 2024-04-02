package erc20bridge

import (
	"math/big"

	ERC20Bridge_V1 "github.com/vegaprotocol/devopstools/smartcontracts/erc20bridge/v1"
	ERC20Bridge_V2 "github.com/vegaprotocol/devopstools/smartcontracts/erc20bridge/v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20BridgeCommon interface {
	GetAssetSource(opts *bind.CallOpts, vega_asset_id [32]byte) (common.Address, error)
	// GetDepositMaximum(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error) // removed in v2
	// GetDepositMinimum(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error) // removed in v2
	GetMultisigControlAddress(opts *bind.CallOpts) (common.Address, error)
	GetVegaAssetId(opts *bind.CallOpts, asset_source common.Address) ([32]byte, error)
	IsAssetListed(opts *bind.CallOpts, asset_source common.Address) (bool, error)

	DepositAsset(opts *bind.TransactOpts, asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error)
	// ListAsset(opts *bind.TransactOpts, asset_source common.Address, vega_asset_id [32]byte, nonce *big.Int, signatures []byte) (*types.Transaction, error) // changed in v2
	RemoveAsset(opts *bind.TransactOpts, asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error)
	// SetDepositMaximum(opts *bind.TransactOpts, asset_source common.Address, maximum_amount *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) // removed in v2
	// SetDepositMinimum(opts *bind.TransactOpts, asset_source common.Address, minimum_amount *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) // removed in v2
	// WithdrawAsset(opts *bind.TransactOpts, asset_source common.Address, amount *big.Int, target common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) // changed in v2
}

type ERC20BridgeNewInV2 interface {
	DefaultWithdrawDelay(opts *bind.CallOpts) (*big.Int, error)
	Erc20AssetPoolAddress(opts *bind.CallOpts) (common.Address, error)
	GetAssetDepositLifetimeLimit(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error)
	GetWithdrawThreshold(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error)
	IsExemptDepositor(opts *bind.CallOpts, depositor common.Address) (bool, error)
	IsStopped(opts *bind.CallOpts) (bool, error)

	ExemptDepositor(opts *bind.TransactOpts) (*types.Transaction, error)
	GlobalResume(opts *bind.TransactOpts, nonce *big.Int, signatures []byte) (*types.Transaction, error)
	GlobalStop(opts *bind.TransactOpts, nonce *big.Int, signatures []byte) (*types.Transaction, error)
	ListAsset(opts *bind.TransactOpts, asset_source common.Address, vega_asset_id [32]byte, lifetime_limit *big.Int, withdraw_threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error)
	RevokeExemptDepositor(opts *bind.TransactOpts) (*types.Transaction, error)
	SetAssetLimits(opts *bind.TransactOpts, asset_source common.Address, lifetime_limit *big.Int, threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error)
	SetWithdrawDelay(opts *bind.TransactOpts, delay *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error)
	WithdrawAsset(opts *bind.TransactOpts, asset_source common.Address, amount *big.Int, target common.Address, creation *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error)
}

type ERC20Bridge struct {
	ERC20BridgeCommon
	ERC20BridgeNewInV2
	Address common.Address
	Version ERC20BridgeVersion
	client  *ethclient.Client

	// Minimal implementation
	v1 *ERC20Bridge_V1.ERC20Bridge
	v2 *ERC20Bridge_V2.ERC20BridgeRestricted
}

func NewERC20Bridge(
	ethClient *ethclient.Client,
	hexAddress string,
	version ERC20BridgeVersion,
) (*ERC20Bridge, error) {
	var err error
	result := &ERC20Bridge{
		Address: common.HexToAddress(hexAddress),
		Version: version,
		client:  ethClient,
	}
	switch version {
	case ERC20BridgeV1:
		result.v1, err = ERC20Bridge_V1.NewERC20Bridge(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.ERC20BridgeCommon = result.v1
	case ERC20BridgeV2:
		result.v2, err = ERC20Bridge_V2.NewERC20BridgeRestricted(result.Address, result.client)
		if err != nil {
			return nil, err
		}
		result.ERC20BridgeCommon = result.v2
		result.ERC20BridgeNewInV2 = result.v2
	}

	return result, nil
}
