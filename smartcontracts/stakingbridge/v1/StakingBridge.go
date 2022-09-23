// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package StakingBridge

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StakingBridgeMetaData contains all meta data concerning the StakingBridge contract.
var StakingBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"Stake_Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"Stake_Removed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"Stake_Transferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"remove_stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"stake_balance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staking_token\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"total_staked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"new_address\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"transfer_stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610b58380380610b588339818101604052810190610032919061008d565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550506100ff565b600081519050610087816100e8565b92915050565b60006020828403121561009f57600080fd5b60006100ad84828501610078565b91505092915050565b60006100c1826100c8565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6100f1816100b6565b81146100fc57600080fd5b50565b610a4a8061010e6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c8063274abf34146100675780632dc7d74c1461009757806348c66e13146100b557806383c592cf146100d1578063af7568dd146100ed578063dd01ba0b1461010b575b600080fd5b610081600480360381019061007c91906106e0565b610127565b60405161008e9190610892565b60405180910390f35b61009f610182565b6040516100ac9190610817565b60405180910390f35b6100cf60048036038101906100ca919061076e565b6101ab565b005b6100eb60048036038101906100e691906107bd565b6102e4565b005b6100f5610456565b6040516101029190610892565b60405180910390f35b610125600480360381019061012091906107bd565b610507565b005b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600083815260200190815260200160002054905092915050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b82600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000838152602001908152602001600020600082825461020b9190610903565b9250508190555082600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000838152602001908152602001600020600082825461027291906108ad565b92505081905550808273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f296aca09e6f616abedcd9cd45ac378207310452b7a713289374fd1b35e2c2fbe866040516102d79190610892565b60405180910390a4505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b815260040161034193929190610832565b602060405180830381600087803b15801561035b57600080fd5b505af115801561036f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610393919061071c565b61039c57600080fd5b81600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600083815260200190815260200160002060008282546103fc91906108ad565b92505081905550803373ffffffffffffffffffffffffffffffffffffffff167f9e3e33edf5dcded4adabc51b1266225d00fa41516bfcad69513fa4eca69519da8460405161044a9190610892565b60405180910390a35050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016104b29190610817565b60206040518083038186803b1580156104ca57600080fd5b505afa1580156104de573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105029190610745565b905090565b81600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600083815260200190815260200160002060008282546105679190610903565b9250508190555060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb33846040518363ffffffff1660e01b81526004016105c9929190610869565b602060405180830381600087803b1580156105e357600080fd5b505af11580156105f7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061061b919061071c565b61062457600080fd5b803373ffffffffffffffffffffffffffffffffffffffff167fa131d16963736e4c641f27a7f82f2e350b5971e555ae06ae906892bbba0a09398460405161066b9190610892565b60405180910390a35050565b600081359050610686816109b8565b92915050565b60008151905061069b816109cf565b92915050565b6000813590506106b0816109e6565b92915050565b6000813590506106c5816109fd565b92915050565b6000815190506106da816109fd565b92915050565b600080604083850312156106f357600080fd5b600061070185828601610677565b9250506020610712858286016106a1565b9150509250929050565b60006020828403121561072e57600080fd5b600061073c8482850161068c565b91505092915050565b60006020828403121561075757600080fd5b6000610765848285016106cb565b91505092915050565b60008060006060848603121561078357600080fd5b6000610791868287016106b6565b93505060206107a286828701610677565b92505060406107b3868287016106a1565b9150509250925092565b600080604083850312156107d057600080fd5b60006107de858286016106b6565b92505060206107ef858286016106a1565b9150509250929050565b61080281610937565b82525050565b6108118161097f565b82525050565b600060208201905061082c60008301846107f9565b92915050565b600060608201905061084760008301866107f9565b61085460208301856107f9565b6108616040830184610808565b949350505050565b600060408201905061087e60008301856107f9565b61088b6020830184610808565b9392505050565b60006020820190506108a76000830184610808565b92915050565b60006108b88261097f565b91506108c38361097f565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156108f8576108f7610989565b5b828201905092915050565b600061090e8261097f565b91506109198361097f565b92508282101561092c5761092b610989565b5b828203905092915050565b60006109428261095f565b9050919050565b60008115159050919050565b6000819050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6109c181610937565b81146109cc57600080fd5b50565b6109d881610949565b81146109e357600080fd5b50565b6109ef81610955565b81146109fa57600080fd5b50565b610a068161097f565b8114610a1157600080fd5b5056fea26469706673582212206b33e3e443456e1fdcb04dca474d33e6008802f4d31ef8464a30604b516b197564736f6c63430008010033",
}

// StakingBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingBridgeMetaData.ABI instead.
var StakingBridgeABI = StakingBridgeMetaData.ABI

// StakingBridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingBridgeMetaData.Bin instead.
var StakingBridgeBin = StakingBridgeMetaData.Bin

// DeployStakingBridge deploys a new Ethereum contract, binding an instance of StakingBridge to it.
func DeployStakingBridge(auth *bind.TransactOpts, backend bind.ContractBackend, token common.Address) (common.Address, *types.Transaction, *StakingBridge, error) {
	parsed, err := StakingBridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingBridgeBin), backend, token)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakingBridge{StakingBridgeCaller: StakingBridgeCaller{contract: contract}, StakingBridgeTransactor: StakingBridgeTransactor{contract: contract}, StakingBridgeFilterer: StakingBridgeFilterer{contract: contract}}, nil
}

// StakingBridge is an auto generated Go binding around an Ethereum contract.
type StakingBridge struct {
	StakingBridgeCaller     // Read-only binding to the contract
	StakingBridgeTransactor // Write-only binding to the contract
	StakingBridgeFilterer   // Log filterer for contract events
}

// StakingBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingBridgeSession struct {
	Contract     *StakingBridge    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingBridgeCallerSession struct {
	Contract *StakingBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StakingBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingBridgeTransactorSession struct {
	Contract     *StakingBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StakingBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingBridgeRaw struct {
	Contract *StakingBridge // Generic contract binding to access the raw methods on
}

// StakingBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingBridgeCallerRaw struct {
	Contract *StakingBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// StakingBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingBridgeTransactorRaw struct {
	Contract *StakingBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingBridge creates a new instance of StakingBridge, bound to a specific deployed contract.
func NewStakingBridge(address common.Address, backend bind.ContractBackend) (*StakingBridge, error) {
	contract, err := bindStakingBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingBridge{StakingBridgeCaller: StakingBridgeCaller{contract: contract}, StakingBridgeTransactor: StakingBridgeTransactor{contract: contract}, StakingBridgeFilterer: StakingBridgeFilterer{contract: contract}}, nil
}

// NewStakingBridgeCaller creates a new read-only instance of StakingBridge, bound to a specific deployed contract.
func NewStakingBridgeCaller(address common.Address, caller bind.ContractCaller) (*StakingBridgeCaller, error) {
	contract, err := bindStakingBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeCaller{contract: contract}, nil
}

// NewStakingBridgeTransactor creates a new write-only instance of StakingBridge, bound to a specific deployed contract.
func NewStakingBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingBridgeTransactor, error) {
	contract, err := bindStakingBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeTransactor{contract: contract}, nil
}

// NewStakingBridgeFilterer creates a new log filterer instance of StakingBridge, bound to a specific deployed contract.
func NewStakingBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingBridgeFilterer, error) {
	contract, err := bindStakingBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeFilterer{contract: contract}, nil
}

// bindStakingBridge binds a generic wrapper to an already deployed contract.
func bindStakingBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingBridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingBridge *StakingBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingBridge.Contract.StakingBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingBridge *StakingBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingBridge.Contract.StakingBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingBridge *StakingBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingBridge.Contract.StakingBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingBridge *StakingBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingBridge *StakingBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingBridge *StakingBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingBridge.Contract.contract.Transact(opts, method, params...)
}

// StakeBalance is a free data retrieval call binding the contract method 0x274abf34.
//
// Solidity: function stake_balance(address target, bytes32 vega_public_key) view returns(uint256)
func (_StakingBridge *StakingBridgeCaller) StakeBalance(opts *bind.CallOpts, target common.Address, vega_public_key [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _StakingBridge.contract.Call(opts, &out, "stake_balance", target, vega_public_key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeBalance is a free data retrieval call binding the contract method 0x274abf34.
//
// Solidity: function stake_balance(address target, bytes32 vega_public_key) view returns(uint256)
func (_StakingBridge *StakingBridgeSession) StakeBalance(target common.Address, vega_public_key [32]byte) (*big.Int, error) {
	return _StakingBridge.Contract.StakeBalance(&_StakingBridge.CallOpts, target, vega_public_key)
}

// StakeBalance is a free data retrieval call binding the contract method 0x274abf34.
//
// Solidity: function stake_balance(address target, bytes32 vega_public_key) view returns(uint256)
func (_StakingBridge *StakingBridgeCallerSession) StakeBalance(target common.Address, vega_public_key [32]byte) (*big.Int, error) {
	return _StakingBridge.Contract.StakeBalance(&_StakingBridge.CallOpts, target, vega_public_key)
}

// StakingToken is a free data retrieval call binding the contract method 0x2dc7d74c.
//
// Solidity: function staking_token() view returns(address)
func (_StakingBridge *StakingBridgeCaller) StakingToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingBridge.contract.Call(opts, &out, "staking_token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingToken is a free data retrieval call binding the contract method 0x2dc7d74c.
//
// Solidity: function staking_token() view returns(address)
func (_StakingBridge *StakingBridgeSession) StakingToken() (common.Address, error) {
	return _StakingBridge.Contract.StakingToken(&_StakingBridge.CallOpts)
}

// StakingToken is a free data retrieval call binding the contract method 0x2dc7d74c.
//
// Solidity: function staking_token() view returns(address)
func (_StakingBridge *StakingBridgeCallerSession) StakingToken() (common.Address, error) {
	return _StakingBridge.Contract.StakingToken(&_StakingBridge.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0xaf7568dd.
//
// Solidity: function total_staked() view returns(uint256)
func (_StakingBridge *StakingBridgeCaller) TotalStaked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingBridge.contract.Call(opts, &out, "total_staked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalStaked is a free data retrieval call binding the contract method 0xaf7568dd.
//
// Solidity: function total_staked() view returns(uint256)
func (_StakingBridge *StakingBridgeSession) TotalStaked() (*big.Int, error) {
	return _StakingBridge.Contract.TotalStaked(&_StakingBridge.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0xaf7568dd.
//
// Solidity: function total_staked() view returns(uint256)
func (_StakingBridge *StakingBridgeCallerSession) TotalStaked() (*big.Int, error) {
	return _StakingBridge.Contract.TotalStaked(&_StakingBridge.CallOpts)
}

// RemoveStake is a paid mutator transaction binding the contract method 0xdd01ba0b.
//
// Solidity: function remove_stake(uint256 amount, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeTransactor) RemoveStake(opts *bind.TransactOpts, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.contract.Transact(opts, "remove_stake", amount, vega_public_key)
}

// RemoveStake is a paid mutator transaction binding the contract method 0xdd01ba0b.
//
// Solidity: function remove_stake(uint256 amount, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeSession) RemoveStake(amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.RemoveStake(&_StakingBridge.TransactOpts, amount, vega_public_key)
}

// RemoveStake is a paid mutator transaction binding the contract method 0xdd01ba0b.
//
// Solidity: function remove_stake(uint256 amount, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeTransactorSession) RemoveStake(amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.RemoveStake(&_StakingBridge.TransactOpts, amount, vega_public_key)
}

// Stake is a paid mutator transaction binding the contract method 0x83c592cf.
//
// Solidity: function stake(uint256 amount, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeTransactor) Stake(opts *bind.TransactOpts, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.contract.Transact(opts, "stake", amount, vega_public_key)
}

// Stake is a paid mutator transaction binding the contract method 0x83c592cf.
//
// Solidity: function stake(uint256 amount, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeSession) Stake(amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.Stake(&_StakingBridge.TransactOpts, amount, vega_public_key)
}

// Stake is a paid mutator transaction binding the contract method 0x83c592cf.
//
// Solidity: function stake(uint256 amount, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeTransactorSession) Stake(amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.Stake(&_StakingBridge.TransactOpts, amount, vega_public_key)
}

// TransferStake is a paid mutator transaction binding the contract method 0x48c66e13.
//
// Solidity: function transfer_stake(uint256 amount, address new_address, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeTransactor) TransferStake(opts *bind.TransactOpts, amount *big.Int, new_address common.Address, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.contract.Transact(opts, "transfer_stake", amount, new_address, vega_public_key)
}

// TransferStake is a paid mutator transaction binding the contract method 0x48c66e13.
//
// Solidity: function transfer_stake(uint256 amount, address new_address, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeSession) TransferStake(amount *big.Int, new_address common.Address, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.TransferStake(&_StakingBridge.TransactOpts, amount, new_address, vega_public_key)
}

// TransferStake is a paid mutator transaction binding the contract method 0x48c66e13.
//
// Solidity: function transfer_stake(uint256 amount, address new_address, bytes32 vega_public_key) returns()
func (_StakingBridge *StakingBridgeTransactorSession) TransferStake(amount *big.Int, new_address common.Address, vega_public_key [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.TransferStake(&_StakingBridge.TransactOpts, amount, new_address, vega_public_key)
}

// StakingBridgeStakeDepositedIterator is returned from FilterStakeDeposited and is used to iterate over the raw logs and unpacked data for StakeDeposited events raised by the StakingBridge contract.
type StakingBridgeStakeDepositedIterator struct {
	Event *StakingBridgeStakeDeposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingBridgeStakeDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingBridgeStakeDeposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingBridgeStakeDeposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingBridgeStakeDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingBridgeStakeDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingBridgeStakeDeposited represents a StakeDeposited event raised by the StakingBridge contract.
type StakingBridgeStakeDeposited struct {
	User          common.Address
	Amount        *big.Int
	VegaPublicKey [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeDeposited is a free log retrieval operation binding the contract event 0x9e3e33edf5dcded4adabc51b1266225d00fa41516bfcad69513fa4eca69519da.
//
// Solidity: event Stake_Deposited(address indexed user, uint256 amount, bytes32 indexed vega_public_key)
func (_StakingBridge *StakingBridgeFilterer) FilterStakeDeposited(opts *bind.FilterOpts, user []common.Address, vega_public_key [][32]byte) (*StakingBridgeStakeDepositedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vega_public_keyRule []interface{}
	for _, vega_public_keyItem := range vega_public_key {
		vega_public_keyRule = append(vega_public_keyRule, vega_public_keyItem)
	}

	logs, sub, err := _StakingBridge.contract.FilterLogs(opts, "Stake_Deposited", userRule, vega_public_keyRule)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeStakeDepositedIterator{contract: _StakingBridge.contract, event: "Stake_Deposited", logs: logs, sub: sub}, nil
}

// WatchStakeDeposited is a free log subscription operation binding the contract event 0x9e3e33edf5dcded4adabc51b1266225d00fa41516bfcad69513fa4eca69519da.
//
// Solidity: event Stake_Deposited(address indexed user, uint256 amount, bytes32 indexed vega_public_key)
func (_StakingBridge *StakingBridgeFilterer) WatchStakeDeposited(opts *bind.WatchOpts, sink chan<- *StakingBridgeStakeDeposited, user []common.Address, vega_public_key [][32]byte) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vega_public_keyRule []interface{}
	for _, vega_public_keyItem := range vega_public_key {
		vega_public_keyRule = append(vega_public_keyRule, vega_public_keyItem)
	}

	logs, sub, err := _StakingBridge.contract.WatchLogs(opts, "Stake_Deposited", userRule, vega_public_keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingBridgeStakeDeposited)
				if err := _StakingBridge.contract.UnpackLog(event, "Stake_Deposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStakeDeposited is a log parse operation binding the contract event 0x9e3e33edf5dcded4adabc51b1266225d00fa41516bfcad69513fa4eca69519da.
//
// Solidity: event Stake_Deposited(address indexed user, uint256 amount, bytes32 indexed vega_public_key)
func (_StakingBridge *StakingBridgeFilterer) ParseStakeDeposited(log types.Log) (*StakingBridgeStakeDeposited, error) {
	event := new(StakingBridgeStakeDeposited)
	if err := _StakingBridge.contract.UnpackLog(event, "Stake_Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingBridgeStakeRemovedIterator is returned from FilterStakeRemoved and is used to iterate over the raw logs and unpacked data for StakeRemoved events raised by the StakingBridge contract.
type StakingBridgeStakeRemovedIterator struct {
	Event *StakingBridgeStakeRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingBridgeStakeRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingBridgeStakeRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingBridgeStakeRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingBridgeStakeRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingBridgeStakeRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingBridgeStakeRemoved represents a StakeRemoved event raised by the StakingBridge contract.
type StakingBridgeStakeRemoved struct {
	User          common.Address
	Amount        *big.Int
	VegaPublicKey [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeRemoved is a free log retrieval operation binding the contract event 0xa131d16963736e4c641f27a7f82f2e350b5971e555ae06ae906892bbba0a0939.
//
// Solidity: event Stake_Removed(address indexed user, uint256 amount, bytes32 indexed vega_public_key)
func (_StakingBridge *StakingBridgeFilterer) FilterStakeRemoved(opts *bind.FilterOpts, user []common.Address, vega_public_key [][32]byte) (*StakingBridgeStakeRemovedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vega_public_keyRule []interface{}
	for _, vega_public_keyItem := range vega_public_key {
		vega_public_keyRule = append(vega_public_keyRule, vega_public_keyItem)
	}

	logs, sub, err := _StakingBridge.contract.FilterLogs(opts, "Stake_Removed", userRule, vega_public_keyRule)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeStakeRemovedIterator{contract: _StakingBridge.contract, event: "Stake_Removed", logs: logs, sub: sub}, nil
}

// WatchStakeRemoved is a free log subscription operation binding the contract event 0xa131d16963736e4c641f27a7f82f2e350b5971e555ae06ae906892bbba0a0939.
//
// Solidity: event Stake_Removed(address indexed user, uint256 amount, bytes32 indexed vega_public_key)
func (_StakingBridge *StakingBridgeFilterer) WatchStakeRemoved(opts *bind.WatchOpts, sink chan<- *StakingBridgeStakeRemoved, user []common.Address, vega_public_key [][32]byte) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vega_public_keyRule []interface{}
	for _, vega_public_keyItem := range vega_public_key {
		vega_public_keyRule = append(vega_public_keyRule, vega_public_keyItem)
	}

	logs, sub, err := _StakingBridge.contract.WatchLogs(opts, "Stake_Removed", userRule, vega_public_keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingBridgeStakeRemoved)
				if err := _StakingBridge.contract.UnpackLog(event, "Stake_Removed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStakeRemoved is a log parse operation binding the contract event 0xa131d16963736e4c641f27a7f82f2e350b5971e555ae06ae906892bbba0a0939.
//
// Solidity: event Stake_Removed(address indexed user, uint256 amount, bytes32 indexed vega_public_key)
func (_StakingBridge *StakingBridgeFilterer) ParseStakeRemoved(log types.Log) (*StakingBridgeStakeRemoved, error) {
	event := new(StakingBridgeStakeRemoved)
	if err := _StakingBridge.contract.UnpackLog(event, "Stake_Removed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingBridgeStakeTransferredIterator is returned from FilterStakeTransferred and is used to iterate over the raw logs and unpacked data for StakeTransferred events raised by the StakingBridge contract.
type StakingBridgeStakeTransferredIterator struct {
	Event *StakingBridgeStakeTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakingBridgeStakeTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingBridgeStakeTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakingBridgeStakeTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakingBridgeStakeTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingBridgeStakeTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingBridgeStakeTransferred represents a StakeTransferred event raised by the StakingBridge contract.
type StakingBridgeStakeTransferred struct {
	From          common.Address
	Amount        *big.Int
	To            common.Address
	VegaPublicKey [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeTransferred is a free log retrieval operation binding the contract event 0x296aca09e6f616abedcd9cd45ac378207310452b7a713289374fd1b35e2c2fbe.
//
// Solidity: event Stake_Transferred(address indexed from, uint256 amount, address indexed to, bytes32 indexed vega_public_key)
func (_StakingBridge *StakingBridgeFilterer) FilterStakeTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address, vega_public_key [][32]byte) (*StakingBridgeStakeTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var vega_public_keyRule []interface{}
	for _, vega_public_keyItem := range vega_public_key {
		vega_public_keyRule = append(vega_public_keyRule, vega_public_keyItem)
	}

	logs, sub, err := _StakingBridge.contract.FilterLogs(opts, "Stake_Transferred", fromRule, toRule, vega_public_keyRule)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeStakeTransferredIterator{contract: _StakingBridge.contract, event: "Stake_Transferred", logs: logs, sub: sub}, nil
}

// WatchStakeTransferred is a free log subscription operation binding the contract event 0x296aca09e6f616abedcd9cd45ac378207310452b7a713289374fd1b35e2c2fbe.
//
// Solidity: event Stake_Transferred(address indexed from, uint256 amount, address indexed to, bytes32 indexed vega_public_key)
func (_StakingBridge *StakingBridgeFilterer) WatchStakeTransferred(opts *bind.WatchOpts, sink chan<- *StakingBridgeStakeTransferred, from []common.Address, to []common.Address, vega_public_key [][32]byte) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var vega_public_keyRule []interface{}
	for _, vega_public_keyItem := range vega_public_key {
		vega_public_keyRule = append(vega_public_keyRule, vega_public_keyItem)
	}

	logs, sub, err := _StakingBridge.contract.WatchLogs(opts, "Stake_Transferred", fromRule, toRule, vega_public_keyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingBridgeStakeTransferred)
				if err := _StakingBridge.contract.UnpackLog(event, "Stake_Transferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStakeTransferred is a log parse operation binding the contract event 0x296aca09e6f616abedcd9cd45ac378207310452b7a713289374fd1b35e2c2fbe.
//
// Solidity: event Stake_Transferred(address indexed from, uint256 amount, address indexed to, bytes32 indexed vega_public_key)
func (_StakingBridge *StakingBridgeFilterer) ParseStakeTransferred(log types.Log) (*StakingBridgeStakeTransferred, error) {
	event := new(StakingBridgeStakeTransferred)
	if err := _StakingBridge.contract.UnpackLog(event, "Stake_Transferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
