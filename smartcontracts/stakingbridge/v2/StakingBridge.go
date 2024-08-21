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
	_ = abi.ConvertType
)

// StakingBridgeMetaData contains all meta data concerning the StakingBridge contract.
var StakingBridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BalanceTooLow\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"vegaPublicKey\",\"type\":\"bytes32\"}],\"name\":\"StakeDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"vegaPublicKey\",\"type\":\"bytes32\"}],\"name\":\"StakeRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"vegaPublicKey\",\"type\":\"bytes32\"}],\"name\":\"StakeTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"vegaPublicKey\",\"type\":\"bytes32\"}],\"name\":\"removeStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"vegaPublicKey\",\"type\":\"bytes32\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"vegaPublicKey\",\"type\":\"bytes32\"}],\"name\":\"stakeBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakingToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalStaked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"vegaPublicKey\",\"type\":\"bytes32\"}],\"name\":\"transferStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakingBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingBridgeMetaData.ABI instead.
var StakingBridgeABI = StakingBridgeMetaData.ABI

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
	parsed, err := StakingBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// StakeBalance is a free data retrieval call binding the contract method 0xd4a5ea85.
//
// Solidity: function stakeBalance(address target, bytes32 vegaPublicKey) view returns(uint256)
func (_StakingBridge *StakingBridgeCaller) StakeBalance(opts *bind.CallOpts, target common.Address, vegaPublicKey [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _StakingBridge.contract.Call(opts, &out, "stakeBalance", target, vegaPublicKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeBalance is a free data retrieval call binding the contract method 0xd4a5ea85.
//
// Solidity: function stakeBalance(address target, bytes32 vegaPublicKey) view returns(uint256)
func (_StakingBridge *StakingBridgeSession) StakeBalance(target common.Address, vegaPublicKey [32]byte) (*big.Int, error) {
	return _StakingBridge.Contract.StakeBalance(&_StakingBridge.CallOpts, target, vegaPublicKey)
}

// StakeBalance is a free data retrieval call binding the contract method 0xd4a5ea85.
//
// Solidity: function stakeBalance(address target, bytes32 vegaPublicKey) view returns(uint256)
func (_StakingBridge *StakingBridgeCallerSession) StakeBalance(target common.Address, vegaPublicKey [32]byte) (*big.Int, error) {
	return _StakingBridge.Contract.StakeBalance(&_StakingBridge.CallOpts, target, vegaPublicKey)
}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_StakingBridge *StakingBridgeCaller) StakingToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingBridge.contract.Call(opts, &out, "stakingToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_StakingBridge *StakingBridgeSession) StakingToken() (common.Address, error) {
	return _StakingBridge.Contract.StakingToken(&_StakingBridge.CallOpts)
}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_StakingBridge *StakingBridgeCallerSession) StakingToken() (common.Address, error) {
	return _StakingBridge.Contract.StakingToken(&_StakingBridge.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_StakingBridge *StakingBridgeCaller) TotalStaked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingBridge.contract.Call(opts, &out, "totalStaked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_StakingBridge *StakingBridgeSession) TotalStaked() (*big.Int, error) {
	return _StakingBridge.Contract.TotalStaked(&_StakingBridge.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_StakingBridge *StakingBridgeCallerSession) TotalStaked() (*big.Int, error) {
	return _StakingBridge.Contract.TotalStaked(&_StakingBridge.CallOpts)
}

// RemoveStake is a paid mutator transaction binding the contract method 0x34f3b31d.
//
// Solidity: function removeStake(uint256 amount, bytes32 vegaPublicKey) returns()
func (_StakingBridge *StakingBridgeTransactor) RemoveStake(opts *bind.TransactOpts, amount *big.Int, vegaPublicKey [32]byte) (*types.Transaction, error) {
	return _StakingBridge.contract.Transact(opts, "removeStake", amount, vegaPublicKey)
}

// RemoveStake is a paid mutator transaction binding the contract method 0x34f3b31d.
//
// Solidity: function removeStake(uint256 amount, bytes32 vegaPublicKey) returns()
func (_StakingBridge *StakingBridgeSession) RemoveStake(amount *big.Int, vegaPublicKey [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.RemoveStake(&_StakingBridge.TransactOpts, amount, vegaPublicKey)
}

// RemoveStake is a paid mutator transaction binding the contract method 0x34f3b31d.
//
// Solidity: function removeStake(uint256 amount, bytes32 vegaPublicKey) returns()
func (_StakingBridge *StakingBridgeTransactorSession) RemoveStake(amount *big.Int, vegaPublicKey [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.RemoveStake(&_StakingBridge.TransactOpts, amount, vegaPublicKey)
}

// Stake is a paid mutator transaction binding the contract method 0x83c592cf.
//
// Solidity: function stake(uint256 amount, bytes32 vegaPublicKey) returns()
func (_StakingBridge *StakingBridgeTransactor) Stake(opts *bind.TransactOpts, amount *big.Int, vegaPublicKey [32]byte) (*types.Transaction, error) {
	return _StakingBridge.contract.Transact(opts, "stake", amount, vegaPublicKey)
}

// Stake is a paid mutator transaction binding the contract method 0x83c592cf.
//
// Solidity: function stake(uint256 amount, bytes32 vegaPublicKey) returns()
func (_StakingBridge *StakingBridgeSession) Stake(amount *big.Int, vegaPublicKey [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.Stake(&_StakingBridge.TransactOpts, amount, vegaPublicKey)
}

// Stake is a paid mutator transaction binding the contract method 0x83c592cf.
//
// Solidity: function stake(uint256 amount, bytes32 vegaPublicKey) returns()
func (_StakingBridge *StakingBridgeTransactorSession) Stake(amount *big.Int, vegaPublicKey [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.Stake(&_StakingBridge.TransactOpts, amount, vegaPublicKey)
}

// TransferStake is a paid mutator transaction binding the contract method 0x679f5e5f.
//
// Solidity: function transferStake(uint256 amount, address newAddress, bytes32 vegaPublicKey) returns()
func (_StakingBridge *StakingBridgeTransactor) TransferStake(opts *bind.TransactOpts, amount *big.Int, newAddress common.Address, vegaPublicKey [32]byte) (*types.Transaction, error) {
	return _StakingBridge.contract.Transact(opts, "transferStake", amount, newAddress, vegaPublicKey)
}

// TransferStake is a paid mutator transaction binding the contract method 0x679f5e5f.
//
// Solidity: function transferStake(uint256 amount, address newAddress, bytes32 vegaPublicKey) returns()
func (_StakingBridge *StakingBridgeSession) TransferStake(amount *big.Int, newAddress common.Address, vegaPublicKey [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.TransferStake(&_StakingBridge.TransactOpts, amount, newAddress, vegaPublicKey)
}

// TransferStake is a paid mutator transaction binding the contract method 0x679f5e5f.
//
// Solidity: function transferStake(uint256 amount, address newAddress, bytes32 vegaPublicKey) returns()
func (_StakingBridge *StakingBridgeTransactorSession) TransferStake(amount *big.Int, newAddress common.Address, vegaPublicKey [32]byte) (*types.Transaction, error) {
	return _StakingBridge.Contract.TransferStake(&_StakingBridge.TransactOpts, amount, newAddress, vegaPublicKey)
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

// FilterStakeDeposited is a free log retrieval operation binding the contract event 0xa740b666eafe67a67d1b2753cb8f8311c88f5c2bdd5077aa463a9f63d08638c4.
//
// Solidity: event StakeDeposited(address indexed user, uint256 amount, bytes32 indexed vegaPublicKey)
func (_StakingBridge *StakingBridgeFilterer) FilterStakeDeposited(opts *bind.FilterOpts, user []common.Address, vegaPublicKey [][32]byte) (*StakingBridgeStakeDepositedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vegaPublicKeyRule []interface{}
	for _, vegaPublicKeyItem := range vegaPublicKey {
		vegaPublicKeyRule = append(vegaPublicKeyRule, vegaPublicKeyItem)
	}

	logs, sub, err := _StakingBridge.contract.FilterLogs(opts, "StakeDeposited", userRule, vegaPublicKeyRule)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeStakeDepositedIterator{contract: _StakingBridge.contract, event: "StakeDeposited", logs: logs, sub: sub}, nil
}

// WatchStakeDeposited is a free log subscription operation binding the contract event 0xa740b666eafe67a67d1b2753cb8f8311c88f5c2bdd5077aa463a9f63d08638c4.
//
// Solidity: event StakeDeposited(address indexed user, uint256 amount, bytes32 indexed vegaPublicKey)
func (_StakingBridge *StakingBridgeFilterer) WatchStakeDeposited(opts *bind.WatchOpts, sink chan<- *StakingBridgeStakeDeposited, user []common.Address, vegaPublicKey [][32]byte) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vegaPublicKeyRule []interface{}
	for _, vegaPublicKeyItem := range vegaPublicKey {
		vegaPublicKeyRule = append(vegaPublicKeyRule, vegaPublicKeyItem)
	}

	logs, sub, err := _StakingBridge.contract.WatchLogs(opts, "StakeDeposited", userRule, vegaPublicKeyRule)
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
				if err := _StakingBridge.contract.UnpackLog(event, "StakeDeposited", log); err != nil {
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

// ParseStakeDeposited is a log parse operation binding the contract event 0xa740b666eafe67a67d1b2753cb8f8311c88f5c2bdd5077aa463a9f63d08638c4.
//
// Solidity: event StakeDeposited(address indexed user, uint256 amount, bytes32 indexed vegaPublicKey)
func (_StakingBridge *StakingBridgeFilterer) ParseStakeDeposited(log types.Log) (*StakingBridgeStakeDeposited, error) {
	event := new(StakingBridgeStakeDeposited)
	if err := _StakingBridge.contract.UnpackLog(event, "StakeDeposited", log); err != nil {
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

// FilterStakeRemoved is a free log retrieval operation binding the contract event 0x3df2a2c33fc4a392029fcbabd913802df02edfae3039eed78ddcc961bbf74f3e.
//
// Solidity: event StakeRemoved(address indexed user, uint256 amount, bytes32 indexed vegaPublicKey)
func (_StakingBridge *StakingBridgeFilterer) FilterStakeRemoved(opts *bind.FilterOpts, user []common.Address, vegaPublicKey [][32]byte) (*StakingBridgeStakeRemovedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vegaPublicKeyRule []interface{}
	for _, vegaPublicKeyItem := range vegaPublicKey {
		vegaPublicKeyRule = append(vegaPublicKeyRule, vegaPublicKeyItem)
	}

	logs, sub, err := _StakingBridge.contract.FilterLogs(opts, "StakeRemoved", userRule, vegaPublicKeyRule)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeStakeRemovedIterator{contract: _StakingBridge.contract, event: "StakeRemoved", logs: logs, sub: sub}, nil
}

// WatchStakeRemoved is a free log subscription operation binding the contract event 0x3df2a2c33fc4a392029fcbabd913802df02edfae3039eed78ddcc961bbf74f3e.
//
// Solidity: event StakeRemoved(address indexed user, uint256 amount, bytes32 indexed vegaPublicKey)
func (_StakingBridge *StakingBridgeFilterer) WatchStakeRemoved(opts *bind.WatchOpts, sink chan<- *StakingBridgeStakeRemoved, user []common.Address, vegaPublicKey [][32]byte) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	var vegaPublicKeyRule []interface{}
	for _, vegaPublicKeyItem := range vegaPublicKey {
		vegaPublicKeyRule = append(vegaPublicKeyRule, vegaPublicKeyItem)
	}

	logs, sub, err := _StakingBridge.contract.WatchLogs(opts, "StakeRemoved", userRule, vegaPublicKeyRule)
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
				if err := _StakingBridge.contract.UnpackLog(event, "StakeRemoved", log); err != nil {
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

// ParseStakeRemoved is a log parse operation binding the contract event 0x3df2a2c33fc4a392029fcbabd913802df02edfae3039eed78ddcc961bbf74f3e.
//
// Solidity: event StakeRemoved(address indexed user, uint256 amount, bytes32 indexed vegaPublicKey)
func (_StakingBridge *StakingBridgeFilterer) ParseStakeRemoved(log types.Log) (*StakingBridgeStakeRemoved, error) {
	event := new(StakingBridgeStakeRemoved)
	if err := _StakingBridge.contract.UnpackLog(event, "StakeRemoved", log); err != nil {
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

// FilterStakeTransferred is a free log retrieval operation binding the contract event 0xab5f7a0e10dc8589661875693d701b8779b7606ee1a15027a7114d66f1257794.
//
// Solidity: event StakeTransferred(address indexed from, uint256 amount, address indexed to, bytes32 indexed vegaPublicKey)
func (_StakingBridge *StakingBridgeFilterer) FilterStakeTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address, vegaPublicKey [][32]byte) (*StakingBridgeStakeTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var vegaPublicKeyRule []interface{}
	for _, vegaPublicKeyItem := range vegaPublicKey {
		vegaPublicKeyRule = append(vegaPublicKeyRule, vegaPublicKeyItem)
	}

	logs, sub, err := _StakingBridge.contract.FilterLogs(opts, "StakeTransferred", fromRule, toRule, vegaPublicKeyRule)
	if err != nil {
		return nil, err
	}
	return &StakingBridgeStakeTransferredIterator{contract: _StakingBridge.contract, event: "StakeTransferred", logs: logs, sub: sub}, nil
}

// WatchStakeTransferred is a free log subscription operation binding the contract event 0xab5f7a0e10dc8589661875693d701b8779b7606ee1a15027a7114d66f1257794.
//
// Solidity: event StakeTransferred(address indexed from, uint256 amount, address indexed to, bytes32 indexed vegaPublicKey)
func (_StakingBridge *StakingBridgeFilterer) WatchStakeTransferred(opts *bind.WatchOpts, sink chan<- *StakingBridgeStakeTransferred, from []common.Address, to []common.Address, vegaPublicKey [][32]byte) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var vegaPublicKeyRule []interface{}
	for _, vegaPublicKeyItem := range vegaPublicKey {
		vegaPublicKeyRule = append(vegaPublicKeyRule, vegaPublicKeyItem)
	}

	logs, sub, err := _StakingBridge.contract.WatchLogs(opts, "StakeTransferred", fromRule, toRule, vegaPublicKeyRule)
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
				if err := _StakingBridge.contract.UnpackLog(event, "StakeTransferred", log); err != nil {
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

// ParseStakeTransferred is a log parse operation binding the contract event 0xab5f7a0e10dc8589661875693d701b8779b7606ee1a15027a7114d66f1257794.
//
// Solidity: event StakeTransferred(address indexed from, uint256 amount, address indexed to, bytes32 indexed vegaPublicKey)
func (_StakingBridge *StakingBridgeFilterer) ParseStakeTransferred(log types.Log) (*StakingBridgeStakeTransferred, error) {
	event := new(StakingBridgeStakeTransferred)
	if err := _StakingBridge.contract.UnpackLog(event, "StakeTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
