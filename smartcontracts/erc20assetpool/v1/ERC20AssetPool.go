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
)

// ERC20AssetPoolMetaData contains all meta data concerning the ERC20AssetPool contract.
var ERC20AssetPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"multisig_control\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"new_address\",\"type\":\"address\"}],\"name\":\"Bridge_Address_Set\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"new_address\",\"type\":\"address\"}],\"name\":\"Multisig_Control_Set\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"erc20_bridge_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multisig_control_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"new_address\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_bridge_address\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"new_address\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_multisig_control\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token_address\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405162000fc938038062000fc983398181016040528101906100349190610120565b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167f1143e675ad794f5bd81a05b165b166be3a4e91f17d065f08809d88cefbd6540660405160405180910390a25061014d565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100ed826100c2565b9050919050565b6100fd816100e2565b811461010857600080fd5b50565b60008151905061011a816100f4565b92915050565b600060208284031215610136576101356100bd565b5b60006101448482850161010b565b91505092915050565b610e6c806200015d6000396000f3fe60806040526004361061004e5760003560e01c806363bb28e014610093578063aeed8f95146100bc578063b82d5abd146100e5578063d9caed1214610110578063e98dfffd1461014d5761008e565b3661008e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161008590610754565b60405180910390fd5b600080fd5b34801561009f57600080fd5b506100ba60048036038101906100b59190610962565b610178565b005b3480156100c857600080fd5b506100e360048036038101906100de9190610962565b610316565b005b3480156100f157600080fd5b506100fa6104ed565b60405161010791906109e0565b60405180910390f35b34801561011c57600080fd5b50610137600480360381019061013291906109fb565b610511565b6040516101449190610a69565b60405180910390f35b34801561015957600080fd5b506101626106ab565b60405161016f91906109e0565b60405180910390f35b6000838360405160200161018d929190610adf565b604051602081830303815290604052905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b81526004016101fb93929190610ba3565b602060405180830381600087803b15801561021557600080fd5b505af1158015610229573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061024d9190610c14565b61028c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161028390610c8d565b60405180910390fd5b83600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508373ffffffffffffffffffffffffffffffffffffffff167ff57d83269802e794395bfd99d93c82bf997d03ec73f8d4c607e8ff18b8cc62f260405160405180910390a250505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141561035057600080fd5b60008383604051602001610365929190610cf9565b604051602081830303815290604052905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b81526004016103d393929190610ba3565b602060405180830381600087803b1580156103ed57600080fd5b505af1158015610401573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104259190610c14565b610464576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161045b90610c8d565b60405180910390fd5b836000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508373ffffffffffffffffffffffffffffffffffffffff167f1143e675ad794f5bd81a05b165b166be3a4e91f17d065f08809d88cefbd6540660405160405180910390a250505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146105a3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161059a90610d81565b60405180910390fd5b8373ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84846040518363ffffffff1660e01b81526004016105de929190610da1565b602060405180830381600087803b1580156105f857600080fd5b505af115801561060c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106309190610c14565b5060003d6000811461064957602081146106525761065e565b6001915061065e565b60206000803e60005191505b508061069f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161069690610e16565b60405180910390fd5b60019150509392505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600082825260208201905092915050565b7f7468697320636f6e747261637420646f6573206e6f742061636365707420455460008201527f4800000000000000000000000000000000000000000000000000000000000000602082015250565b600061073e6021836106d1565b9150610749826106e2565b604082019050919050565b6000602082019050818103600083015261076d81610731565b9050919050565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006107b382610788565b9050919050565b6107c3816107a8565b81146107ce57600080fd5b50565b6000813590506107e0816107ba565b92915050565b6000819050919050565b6107f9816107e6565b811461080457600080fd5b50565b600081359050610816816107f0565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61086f82610826565b810181811067ffffffffffffffff8211171561088e5761088d610837565b5b80604052505050565b60006108a1610774565b90506108ad8282610866565b919050565b600067ffffffffffffffff8211156108cd576108cc610837565b5b6108d682610826565b9050602081019050919050565b82818337600083830152505050565b6000610905610900846108b2565b610897565b90508281526020810184848401111561092157610920610821565b5b61092c8482856108e3565b509392505050565b600082601f8301126109495761094861081c565b5b81356109598482602086016108f2565b91505092915050565b60008060006060848603121561097b5761097a61077e565b5b6000610989868287016107d1565b935050602061099a86828701610807565b925050604084013567ffffffffffffffff8111156109bb576109ba610783565b5b6109c786828701610934565b9150509250925092565b6109da816107a8565b82525050565b60006020820190506109f560008301846109d1565b92915050565b600080600060608486031215610a1457610a1361077e565b5b6000610a22868287016107d1565b9350506020610a33868287016107d1565b9250506040610a4486828701610807565b9150509250925092565b60008115159050919050565b610a6381610a4e565b82525050565b6000602082019050610a7e6000830184610a5a565b92915050565b610a8d816107e6565b82525050565b7f7365745f6272696467655f616464726573730000000000000000000000000000600082015250565b6000610ac96012836106d1565b9150610ad482610a93565b602082019050919050565b6000606082019050610af460008301856109d1565b610b016020830184610a84565b8181036040830152610b1281610abc565b90509392505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610b55578082015181840152602081019050610b3a565b83811115610b64576000848401525b50505050565b6000610b7582610b1b565b610b7f8185610b26565b9350610b8f818560208601610b37565b610b9881610826565b840191505092915050565b60006060820190508181036000830152610bbd8186610b6a565b90508181036020830152610bd18185610b6a565b9050610be06040830184610a84565b949350505050565b610bf181610a4e565b8114610bfc57600080fd5b50565b600081519050610c0e81610be8565b92915050565b600060208284031215610c2a57610c2961077e565b5b6000610c3884828501610bff565b91505092915050565b7f626164207369676e617475726573000000000000000000000000000000000000600082015250565b6000610c77600e836106d1565b9150610c8282610c41565b602082019050919050565b60006020820190508181036000830152610ca681610c6a565b9050919050565b7f7365745f6d756c74697369675f636f6e74726f6c000000000000000000000000600082015250565b6000610ce36014836106d1565b9150610cee82610cad565b602082019050919050565b6000606082019050610d0e60008301856109d1565b610d1b6020830184610a84565b8181036040830152610d2c81610cd6565b90509392505050565b7f6d73672e73656e646572206e6f7420617574686f72697a656420627269646765600082015250565b6000610d6b6020836106d1565b9150610d7682610d35565b602082019050919050565b60006020820190508181036000830152610d9a81610d5e565b9050919050565b6000604082019050610db660008301856109d1565b610dc36020830184610a84565b9392505050565b7f746f6b656e207472616e73666572206661696c65640000000000000000000000600082015250565b6000610e006015836106d1565b9150610e0b82610dca565b602082019050919050565b60006020820190508181036000830152610e2f81610df3565b905091905056fea2646970667358221220d0ba88771aea5ef02e9bbea1437104ba0ec780ec4415a7a84caed2b8559466b864736f6c63430008080033",
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
	parsed, err := abi.JSON(strings.NewReader(ERC20AssetPoolABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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
// Solidity: function withdraw(address token_address, address target, uint256 amount) returns(bool)
func (_ERC20AssetPool *ERC20AssetPoolTransactor) Withdraw(opts *bind.TransactOpts, token_address common.Address, target common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20AssetPool.contract.Transact(opts, "withdraw", token_address, target, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token_address, address target, uint256 amount) returns(bool)
func (_ERC20AssetPool *ERC20AssetPoolSession) Withdraw(token_address common.Address, target common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20AssetPool.Contract.Withdraw(&_ERC20AssetPool.TransactOpts, token_address, target, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token_address, address target, uint256 amount) returns(bool)
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
