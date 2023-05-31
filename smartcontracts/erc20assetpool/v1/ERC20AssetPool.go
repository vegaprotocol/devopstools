// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ERC20AssetPool

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

// ERC20AssetPoolMetaData contains all meta data concerning the ERC20AssetPool contract.
var ERC20AssetPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"multisig_control\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"new_address\",\"type\":\"address\"}],\"name\":\"Bridge_Address_Set\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"new_address\",\"type\":\"address\"}],\"name\":\"Multisig_Control_Set\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"erc20_bridge_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multisig_control_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"new_address\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_bridge_address\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"new_address\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_multisig_control\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token_address\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610d6a380380610d6a83398101604081905261002f916100d4565b6001600160a01b0381166100895760405162461bcd60e51b815260206004820152601f60248201527f696e76616c6964204d756c7469736967436f6e74726f6c206164647265737300604482015260640160405180910390fd5b600080546001600160a01b0319166001600160a01b038316908117825560405190917f1143e675ad794f5bd81a05b165b166be3a4e91f17d065f08809d88cefbd6540691a250610104565b6000602082840312156100e657600080fd5b81516001600160a01b03811681146100fd57600080fd5b9392505050565b610c57806101136000396000f3fe60806040526004361061005e5760003560e01c8063b82d5abd11610043578063b82d5abd14610137578063d9caed121461018d578063e98dfffd146101ad57600080fd5b806363bb28e0146100f5578063aeed8f951461011757600080fd5b366100f0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602160248201527f7468697320636f6e747261637420646f6573206e6f742061636365707420455460448201527f480000000000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b600080fd5b34801561010157600080fd5b50610115610110366004610a07565b6101da565b005b34801561012357600080fd5b50610115610132366004610a07565b6103ea565b34801561014357600080fd5b506000546101649073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b34801561019957600080fd5b506101156101a8366004610af0565b6106dd565b3480156101b957600080fd5b506001546101649073ffffffffffffffffffffffffffffffffffffffff1681565b6040805173ffffffffffffffffffffffffffffffffffffffff85166020820152908101839052606080820152601260808201527f7365745f6272696467655f61646472657373000000000000000000000000000060a082015260009060c001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526000547fba73659a00000000000000000000000000000000000000000000000000000000835290925073ffffffffffffffffffffffffffffffffffffffff169063ba73659a906102c090859085908890600401610ba6565b602060405180830381600087803b1580156102da57600080fd5b505af11580156102ee573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103129190610bdc565b610378576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e61747572657300000000000000000000000000000000000060448201526064016100e7565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff86169081179091556040517ff57d83269802e794395bfd99d93c82bf997d03ec73f8d4c607e8ff18b8cc62f290600090a250505050565b73ffffffffffffffffffffffffffffffffffffffff8316610467576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f696e76616c6964204d756c7469736967436f6e74726f6c20616464726573730060448201526064016100e7565b823b6104cf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f6e65772061646472657373206d75737420626520636f6e74726163740000000060448201526064016100e7565b6040805173ffffffffffffffffffffffffffffffffffffffff85166020820152908101839052606080820152601460808201527f7365745f6d756c74697369675f636f6e74726f6c00000000000000000000000060a082015260009060c001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152908290526000547fba73659a00000000000000000000000000000000000000000000000000000000835290925073ffffffffffffffffffffffffffffffffffffffff169063ba73659a906105b590859085908890600401610ba6565b602060405180830381600087803b1580156105cf57600080fd5b505af11580156105e3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106079190610bdc565b61066d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e61747572657300000000000000000000000000000000000060448201526064016100e7565b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff8616908117825560405190917f1143e675ad794f5bd81a05b165b166be3a4e91f17d065f08809d88cefbd6540691a250505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461075e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6d73672e73656e646572206e6f7420617574686f72697a65642062726964676560448201526064016100e7565b823b6107c6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f746f6b656e5f61646472657373206d75737420626520636f6e7472616374000060448201526064016100e7565b60405173ffffffffffffffffffffffffffffffffffffffff8381166024830152604482018390526000918291861690606401604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052516108799190610c05565b6000604051808303816000865af19150503d80600081146108b6576040519150601f19603f3d011682016040523d82523d6000602084013e6108bb565b606091505b509150915081610927576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f746f6b656e207472616e73666572206661696c6564000000000000000000000060448201526064016100e7565b8051156109a857808060200190518101906109429190610bdc565b6109a8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f746f6b656e207472616e73666572206661696c6564000000000000000000000060448201526064016100e7565b5050505050565b803573ffffffffffffffffffffffffffffffffffffffff811681146109d357600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080600060608486031215610a1c57600080fd5b610a25846109af565b925060208401359150604084013567ffffffffffffffff80821115610a4957600080fd5b818601915086601f830112610a5d57600080fd5b813581811115610a6f57610a6f6109d8565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610ab557610ab56109d8565b81604052828152896020848701011115610ace57600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b600080600060608486031215610b0557600080fd5b610b0e846109af565b9250610b1c602085016109af565b9150604084013590509250925092565b60005b83811015610b47578181015183820152602001610b2f565b83811115610b56576000848401525b50505050565b60008151808452610b74816020860160208601610b2c565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b606081526000610bb96060830186610b5c565b8281036020840152610bcb8186610b5c565b915050826040830152949350505050565b600060208284031215610bee57600080fd5b81518015158114610bfe57600080fd5b9392505050565b60008251610c17818460208701610b2c565b919091019291505056fea2646970667358221220ee4f009f7640390507946a8d3cff89ca7e3bc1a9708d6044a8fe8fa687f6f66d64736f6c63430008080033",
}

// ERC20AssetPoolABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20AssetPoolMetaData.ABI instead.
var ERC20AssetPoolABI = ERC20AssetPoolMetaData.ABI

// ERC20AssetPoolBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC20AssetPoolMetaData.Bin instead.
var ERC20AssetPoolBin = ERC20AssetPoolMetaData.Bin

// DeployERC20AssetPool deploys a new Ethereum contract, binding an instance of ERC20AssetPool to it.
func DeployERC20AssetPool(auth *bind.TransactOpts, backend bind.ContractBackend, multisig_control common.Address) (common.Address, *types.Transaction, *ERC20AssetPool, error) {
	parsed, err := ERC20AssetPoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC20AssetPoolBin), backend, multisig_control)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20AssetPool{ERC20AssetPoolCaller: ERC20AssetPoolCaller{contract: contract}, ERC20AssetPoolTransactor: ERC20AssetPoolTransactor{contract: contract}, ERC20AssetPoolFilterer: ERC20AssetPoolFilterer{contract: contract}}, nil
}

// ERC20AssetPool is an auto generated Go binding around an Ethereum contract.
type ERC20AssetPool struct {
	ERC20AssetPoolCaller     // Read-only binding to the contract
	ERC20AssetPoolTransactor // Write-only binding to the contract
	ERC20AssetPoolFilterer   // Log filterer for contract events
}

// ERC20AssetPoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20AssetPoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20AssetPoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20AssetPoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20AssetPoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20AssetPoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20AssetPoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20AssetPoolSession struct {
	Contract     *ERC20AssetPool   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20AssetPoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20AssetPoolCallerSession struct {
	Contract *ERC20AssetPoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ERC20AssetPoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20AssetPoolTransactorSession struct {
	Contract     *ERC20AssetPoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ERC20AssetPoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20AssetPoolRaw struct {
	Contract *ERC20AssetPool // Generic contract binding to access the raw methods on
}

// ERC20AssetPoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20AssetPoolCallerRaw struct {
	Contract *ERC20AssetPoolCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20AssetPoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20AssetPoolTransactorRaw struct {
	Contract *ERC20AssetPoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20AssetPool creates a new instance of ERC20AssetPool, bound to a specific deployed contract.
func NewERC20AssetPool(address common.Address, backend bind.ContractBackend) (*ERC20AssetPool, error) {
	contract, err := bindERC20AssetPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20AssetPool{ERC20AssetPoolCaller: ERC20AssetPoolCaller{contract: contract}, ERC20AssetPoolTransactor: ERC20AssetPoolTransactor{contract: contract}, ERC20AssetPoolFilterer: ERC20AssetPoolFilterer{contract: contract}}, nil
}

// NewERC20AssetPoolCaller creates a new read-only instance of ERC20AssetPool, bound to a specific deployed contract.
func NewERC20AssetPoolCaller(address common.Address, caller bind.ContractCaller) (*ERC20AssetPoolCaller, error) {
	contract, err := bindERC20AssetPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20AssetPoolCaller{contract: contract}, nil
}

// NewERC20AssetPoolTransactor creates a new write-only instance of ERC20AssetPool, bound to a specific deployed contract.
func NewERC20AssetPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20AssetPoolTransactor, error) {
	contract, err := bindERC20AssetPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20AssetPoolTransactor{contract: contract}, nil
}

// NewERC20AssetPoolFilterer creates a new log filterer instance of ERC20AssetPool, bound to a specific deployed contract.
func NewERC20AssetPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20AssetPoolFilterer, error) {
	contract, err := bindERC20AssetPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20AssetPoolFilterer{contract: contract}, nil
}

// bindERC20AssetPool binds a generic wrapper to an already deployed contract.
func bindERC20AssetPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC20AssetPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20AssetPool *ERC20AssetPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20AssetPool.Contract.ERC20AssetPoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20AssetPool *ERC20AssetPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.ERC20AssetPoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20AssetPool *ERC20AssetPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.ERC20AssetPoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20AssetPool *ERC20AssetPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20AssetPool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20AssetPool *ERC20AssetPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20AssetPool *ERC20AssetPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.contract.Transact(opts, method, params...)
}

// Erc20BridgeAddress is a free data retrieval call binding the contract method 0xe98dfffd.
//
// Solidity: function erc20_bridge_address() view returns(address)
func (_ERC20AssetPool *ERC20AssetPoolCaller) Erc20BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20AssetPool.contract.Call(opts, &out, "erc20_bridge_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Erc20BridgeAddress is a free data retrieval call binding the contract method 0xe98dfffd.
//
// Solidity: function erc20_bridge_address() view returns(address)
func (_ERC20AssetPool *ERC20AssetPoolSession) Erc20BridgeAddress() (common.Address, error) {
	return _ERC20AssetPool.Contract.Erc20BridgeAddress(&_ERC20AssetPool.CallOpts)
}

// Erc20BridgeAddress is a free data retrieval call binding the contract method 0xe98dfffd.
//
// Solidity: function erc20_bridge_address() view returns(address)
func (_ERC20AssetPool *ERC20AssetPoolCallerSession) Erc20BridgeAddress() (common.Address, error) {
	return _ERC20AssetPool.Contract.Erc20BridgeAddress(&_ERC20AssetPool.CallOpts)
}

// MultisigControlAddress is a free data retrieval call binding the contract method 0xb82d5abd.
//
// Solidity: function multisig_control_address() view returns(address)
func (_ERC20AssetPool *ERC20AssetPoolCaller) MultisigControlAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20AssetPool.contract.Call(opts, &out, "multisig_control_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MultisigControlAddress is a free data retrieval call binding the contract method 0xb82d5abd.
//
// Solidity: function multisig_control_address() view returns(address)
func (_ERC20AssetPool *ERC20AssetPoolSession) MultisigControlAddress() (common.Address, error) {
	return _ERC20AssetPool.Contract.MultisigControlAddress(&_ERC20AssetPool.CallOpts)
}

// MultisigControlAddress is a free data retrieval call binding the contract method 0xb82d5abd.
//
// Solidity: function multisig_control_address() view returns(address)
func (_ERC20AssetPool *ERC20AssetPoolCallerSession) MultisigControlAddress() (common.Address, error) {
	return _ERC20AssetPool.Contract.MultisigControlAddress(&_ERC20AssetPool.CallOpts)
}

// SetBridgeAddress is a paid mutator transaction binding the contract method 0x63bb28e0.
//
// Solidity: function set_bridge_address(address new_address, uint256 nonce, bytes signatures) returns()
func (_ERC20AssetPool *ERC20AssetPoolTransactor) SetBridgeAddress(opts *bind.TransactOpts, new_address common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20AssetPool.contract.Transact(opts, "set_bridge_address", new_address, nonce, signatures)
}

// SetBridgeAddress is a paid mutator transaction binding the contract method 0x63bb28e0.
//
// Solidity: function set_bridge_address(address new_address, uint256 nonce, bytes signatures) returns()
func (_ERC20AssetPool *ERC20AssetPoolSession) SetBridgeAddress(new_address common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.SetBridgeAddress(&_ERC20AssetPool.TransactOpts, new_address, nonce, signatures)
}

// SetBridgeAddress is a paid mutator transaction binding the contract method 0x63bb28e0.
//
// Solidity: function set_bridge_address(address new_address, uint256 nonce, bytes signatures) returns()
func (_ERC20AssetPool *ERC20AssetPoolTransactorSession) SetBridgeAddress(new_address common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.SetBridgeAddress(&_ERC20AssetPool.TransactOpts, new_address, nonce, signatures)
}

// SetMultisigControl is a paid mutator transaction binding the contract method 0xaeed8f95.
//
// Solidity: function set_multisig_control(address new_address, uint256 nonce, bytes signatures) returns()
func (_ERC20AssetPool *ERC20AssetPoolTransactor) SetMultisigControl(opts *bind.TransactOpts, new_address common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20AssetPool.contract.Transact(opts, "set_multisig_control", new_address, nonce, signatures)
}

// SetMultisigControl is a paid mutator transaction binding the contract method 0xaeed8f95.
//
// Solidity: function set_multisig_control(address new_address, uint256 nonce, bytes signatures) returns()
func (_ERC20AssetPool *ERC20AssetPoolSession) SetMultisigControl(new_address common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.SetMultisigControl(&_ERC20AssetPool.TransactOpts, new_address, nonce, signatures)
}

// SetMultisigControl is a paid mutator transaction binding the contract method 0xaeed8f95.
//
// Solidity: function set_multisig_control(address new_address, uint256 nonce, bytes signatures) returns()
func (_ERC20AssetPool *ERC20AssetPoolTransactorSession) SetMultisigControl(new_address common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.SetMultisigControl(&_ERC20AssetPool.TransactOpts, new_address, nonce, signatures)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token_address, address target, uint256 amount) returns()
func (_ERC20AssetPool *ERC20AssetPoolTransactor) Withdraw(opts *bind.TransactOpts, token_address common.Address, target common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20AssetPool.contract.Transact(opts, "withdraw", token_address, target, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token_address, address target, uint256 amount) returns()
func (_ERC20AssetPool *ERC20AssetPoolSession) Withdraw(token_address common.Address, target common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.Withdraw(&_ERC20AssetPool.TransactOpts, token_address, target, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token_address, address target, uint256 amount) returns()
func (_ERC20AssetPool *ERC20AssetPoolTransactorSession) Withdraw(token_address common.Address, target common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.Withdraw(&_ERC20AssetPool.TransactOpts, token_address, target, amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ERC20AssetPool *ERC20AssetPoolTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20AssetPool.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ERC20AssetPool *ERC20AssetPoolSession) Receive() (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.Receive(&_ERC20AssetPool.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ERC20AssetPool *ERC20AssetPoolTransactorSession) Receive() (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.Receive(&_ERC20AssetPool.TransactOpts)
}

// ERC20AssetPoolBridgeAddressSetIterator is returned from FilterBridgeAddressSet and is used to iterate over the raw logs and unpacked data for BridgeAddressSet events raised by the ERC20AssetPool contract.
type ERC20AssetPoolBridgeAddressSetIterator struct {
	Event *ERC20AssetPoolBridgeAddressSet // Event containing the contract specifics and raw log

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
func (it *ERC20AssetPoolBridgeAddressSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20AssetPoolBridgeAddressSet)
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
		it.Event = new(ERC20AssetPoolBridgeAddressSet)
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
func (it *ERC20AssetPoolBridgeAddressSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20AssetPoolBridgeAddressSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20AssetPoolBridgeAddressSet represents a BridgeAddressSet event raised by the ERC20AssetPool contract.
type ERC20AssetPoolBridgeAddressSet struct {
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBridgeAddressSet is a free log retrieval operation binding the contract event 0xf57d83269802e794395bfd99d93c82bf997d03ec73f8d4c607e8ff18b8cc62f2.
//
// Solidity: event Bridge_Address_Set(address indexed new_address)
func (_ERC20AssetPool *ERC20AssetPoolFilterer) FilterBridgeAddressSet(opts *bind.FilterOpts, new_address []common.Address) (*ERC20AssetPoolBridgeAddressSetIterator, error) {

	var new_addressRule []interface{}
	for _, new_addressItem := range new_address {
		new_addressRule = append(new_addressRule, new_addressItem)
	}

	logs, sub, err := _ERC20AssetPool.contract.FilterLogs(opts, "Bridge_Address_Set", new_addressRule)
	if err != nil {
		return nil, err
	}
	return &ERC20AssetPoolBridgeAddressSetIterator{contract: _ERC20AssetPool.contract, event: "Bridge_Address_Set", logs: logs, sub: sub}, nil
}

// WatchBridgeAddressSet is a free log subscription operation binding the contract event 0xf57d83269802e794395bfd99d93c82bf997d03ec73f8d4c607e8ff18b8cc62f2.
//
// Solidity: event Bridge_Address_Set(address indexed new_address)
func (_ERC20AssetPool *ERC20AssetPoolFilterer) WatchBridgeAddressSet(opts *bind.WatchOpts, sink chan<- *ERC20AssetPoolBridgeAddressSet, new_address []common.Address) (event.Subscription, error) {

	var new_addressRule []interface{}
	for _, new_addressItem := range new_address {
		new_addressRule = append(new_addressRule, new_addressItem)
	}

	logs, sub, err := _ERC20AssetPool.contract.WatchLogs(opts, "Bridge_Address_Set", new_addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20AssetPoolBridgeAddressSet)
				if err := _ERC20AssetPool.contract.UnpackLog(event, "Bridge_Address_Set", log); err != nil {
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

// ParseBridgeAddressSet is a log parse operation binding the contract event 0xf57d83269802e794395bfd99d93c82bf997d03ec73f8d4c607e8ff18b8cc62f2.
//
// Solidity: event Bridge_Address_Set(address indexed new_address)
func (_ERC20AssetPool *ERC20AssetPoolFilterer) ParseBridgeAddressSet(log types.Log) (*ERC20AssetPoolBridgeAddressSet, error) {
	event := new(ERC20AssetPoolBridgeAddressSet)
	if err := _ERC20AssetPool.contract.UnpackLog(event, "Bridge_Address_Set", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20AssetPoolMultisigControlSetIterator is returned from FilterMultisigControlSet and is used to iterate over the raw logs and unpacked data for MultisigControlSet events raised by the ERC20AssetPool contract.
type ERC20AssetPoolMultisigControlSetIterator struct {
	Event *ERC20AssetPoolMultisigControlSet // Event containing the contract specifics and raw log

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
func (it *ERC20AssetPoolMultisigControlSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20AssetPoolMultisigControlSet)
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
		it.Event = new(ERC20AssetPoolMultisigControlSet)
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
func (it *ERC20AssetPoolMultisigControlSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20AssetPoolMultisigControlSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20AssetPoolMultisigControlSet represents a MultisigControlSet event raised by the ERC20AssetPool contract.
type ERC20AssetPoolMultisigControlSet struct {
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMultisigControlSet is a free log retrieval operation binding the contract event 0x1143e675ad794f5bd81a05b165b166be3a4e91f17d065f08809d88cefbd65406.
//
// Solidity: event Multisig_Control_Set(address indexed new_address)
func (_ERC20AssetPool *ERC20AssetPoolFilterer) FilterMultisigControlSet(opts *bind.FilterOpts, new_address []common.Address) (*ERC20AssetPoolMultisigControlSetIterator, error) {

	var new_addressRule []interface{}
	for _, new_addressItem := range new_address {
		new_addressRule = append(new_addressRule, new_addressItem)
	}

	logs, sub, err := _ERC20AssetPool.contract.FilterLogs(opts, "Multisig_Control_Set", new_addressRule)
	if err != nil {
		return nil, err
	}
	return &ERC20AssetPoolMultisigControlSetIterator{contract: _ERC20AssetPool.contract, event: "Multisig_Control_Set", logs: logs, sub: sub}, nil
}

// WatchMultisigControlSet is a free log subscription operation binding the contract event 0x1143e675ad794f5bd81a05b165b166be3a4e91f17d065f08809d88cefbd65406.
//
// Solidity: event Multisig_Control_Set(address indexed new_address)
func (_ERC20AssetPool *ERC20AssetPoolFilterer) WatchMultisigControlSet(opts *bind.WatchOpts, sink chan<- *ERC20AssetPoolMultisigControlSet, new_address []common.Address) (event.Subscription, error) {

	var new_addressRule []interface{}
	for _, new_addressItem := range new_address {
		new_addressRule = append(new_addressRule, new_addressItem)
	}

	logs, sub, err := _ERC20AssetPool.contract.WatchLogs(opts, "Multisig_Control_Set", new_addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20AssetPoolMultisigControlSet)
				if err := _ERC20AssetPool.contract.UnpackLog(event, "Multisig_Control_Set", log); err != nil {
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

// ParseMultisigControlSet is a log parse operation binding the contract event 0x1143e675ad794f5bd81a05b165b166be3a4e91f17d065f08809d88cefbd65406.
//
// Solidity: event Multisig_Control_Set(address indexed new_address)
func (_ERC20AssetPool *ERC20AssetPoolFilterer) ParseMultisigControlSet(log types.Log) (*ERC20AssetPoolMultisigControlSet, error) {
	event := new(ERC20AssetPoolMultisigControlSet)
	if err := _ERC20AssetPool.contract.UnpackLog(event, "Multisig_Control_Set", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
