// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ClaimCodes

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

// Claim is an auto generated low-level Go binding around an user-defined struct.
type Claim struct {
	Amount  *big.Int
	Tranche uint8
	Expiry  uint32
}

// Signature is an auto generated low-level Go binding around an user-defined struct.
type Signature struct {
	V uint8
	R [32]byte
	S [32]byte
}

// ClaimCodesMetaData contains all meta data concerning the ClaimCodes contract.
var ClaimCodesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vesting_address\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes2[]\",\"name\":\"countries\",\"type\":\"bytes2[]\"}],\"name\":\"allow_countries\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes2\",\"name\":\"\",\"type\":\"bytes2\"}],\"name\":\"allowed_countries\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes2[]\",\"name\":\"countries\",\"type\":\"bytes2[]\"}],\"name\":\"block_countries\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structSignature\",\"name\":\"sig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"tranche\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"expiry\",\"type\":\"uint32\"}],\"internalType\":\"structClaim\",\"name\":\"clm\",\"type\":\"tuple\"},{\"internalType\":\"bytes2\",\"name\":\"country\",\"type\":\"bytes2\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"claim_targeted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structSignature\",\"name\":\"sig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"tranche\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"expiry\",\"type\":\"uint32\"}],\"internalType\":\"structClaim\",\"name\":\"clm\",\"type\":\"tuple\"},{\"internalType\":\"bytes2\",\"name\":\"country\",\"type\":\"bytes2\"}],\"name\":\"claim_untargeted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"commit_untargeted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"commitments\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"controller\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"destroy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"issuers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"permit_issuer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"}],\"name\":\"revoke_issuer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_controller\",\"type\":\"address\"}],\"name\":\"swap_controller\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b506040516111c03803806111c083398101604081905261002f91610058565b60601b6001600160601b03191660805260008054336001600160a01b0319909116179055610088565b60006020828403121561006a57600080fd5b81516001600160a01b038116811461008157600080fd5b9392505050565b60805160601c61111a6100a66000396000610937015261111a6000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c806383197ef01161008c5780639cef4355116100665780639cef43551461020e578063b7d00d0c14610221578063c342cd9b14610234578063f77c47911461024757600080fd5b806383197ef014610198578063839df945146101a057806394a4ff3f146101fb57600080fd5b806334b4947b116100bd57806334b4947b1461011f57806338a7543e146101575780637eb1b5d71461018557600080fd5b806304df9479146100e457806326bc922a146100f9578063277817c71461010c575b600080fd5b6100f76100f2366004610f2d565b610267565b005b6100f7610107366004610e3f565b610279565b6100f761011a366004610e15565b61034b565b61014261012d366004610eb4565b60036020526000908152604090205460ff1681565b60405190151581526020015b60405180910390f35b610177610165366004610df3565b60026020526000908152604090205481565b60405190815260200161014e565b6100f7610193366004610df3565b610398565b6100f7610423565b6101d66101ae366004610ecf565b60016020526000908152604090205473ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161014e565b6100f7610209366004610df3565b61044a565b6100f761021c366004610e3f565b610495565b6100f761022f366004610ee8565b610562565b6100f7610242366004610ecf565b61056f565b6000546101d69073ffffffffffffffffffffffffffffffffffffffff1681565b610273848483856105d6565b50505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461029d57600080fd5b60005b81811015610346576000600360008585858181106102c0576102c06110b5565b90506020020160208101906102d59190610eb4565b7fffff000000000000000000000000000000000000000000000000000000000000168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558061033e8161104d565b9150506102a0565b505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461036f57600080fd5b73ffffffffffffffffffffffffffffffffffffffff909116600090815260026020526040902055565b60005473ffffffffffffffffffffffffffffffffffffffff1633146103bc57600080fd5b73ffffffffffffffffffffffffffffffffffffffff81166103dc57600080fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff16331461044757600080fd5b33ff5b60005473ffffffffffffffffffffffffffffffffffffffff16331461046e57600080fd5b73ffffffffffffffffffffffffffffffffffffffff16600090815260026020526040812055565b60005473ffffffffffffffffffffffffffffffffffffffff1633146104b957600080fd5b60005b81811015610346576001600360008585858181106104dc576104dc6110b5565b90506020020160208101906104f19190610eb4565b7fffff000000000000000000000000000000000000000000000000000000000000168152602081019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169115159190911790558061055a8161104d565b9150506104bc565b61034683836000846105d6565b60008181526001602052604090205473ffffffffffffffffffffffffffffffffffffffff161561059e57600080fd5b600090815260016020526040902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001633179055565b426105e76060850160408601610f84565b63ffffffff1611610659576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f436c61696d20636f64652068617320657870697265640000000000000000000060448201526064015b60405180910390fd5b7fffff000000000000000000000000000000000000000000000000000000000000811660009081526003602052604090205460ff1661071a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f436c61696d20636f6465206973206e6f7420617661696c61626c6520696e207960448201527f6f757220636f756e7472790000000000000000000000000000000000000000006064820152608401610650565b6000610727468585610a1b565b905060006107358287610b35565b905073ffffffffffffffffffffffffffffffffffffffff81166107b4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f496e76616c696420636c61696d20636f646500000000000000000000000000006044820152606401610650565b6107be8685610be2565b93508373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610856576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f43616e6e6f7420636c61696d20746f20796f757273656c6600000000000000006044820152606401610650565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602052604090205485358110156108e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f4f7574206f662066756e647300000000000000000000000000000000000000006044820152606401610650565b6108f2863582611036565b600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663e2de6e6d868860200160208101906109859190610faa565b60405160e084901b7fffffffff0000000000000000000000000000000000000000000000000000000016815273ffffffffffffffffffffffffffffffffffffffff909216600483015260ff16602482015288356044820152606401600060405180830381600087803b1580156109fa57600080fd5b505af1158015610a0e573d6000803e3d6000fd5b5050505050505050505050565b600080848435610a316040870160208801610faa565b610a416060880160408901610f84565b604051602001610ab49493929190938452602084019290925260f81b7fff0000000000000000000000000000000000000000000000000000000000000016604083015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016604182015260450190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052905073ffffffffffffffffffffffffffffffffffffffff831615610b25578083604051602001610b13929190610fcd565b60405160208183030381529060405290505b8051602090910120949350505050565b60007f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a060408301351115610b6b57506000610bdc565b600183610b7b6020850185610faa565b604080516000815260208181018084529490945260ff9092168282015291850135606082015290840135608082015260a0016020604051602081039080840390855afa158015610bcf573d6000803e3d6000fd5b5050506020604051035190505b92915050565b604080830135600090815260016020529081205473ffffffffffffffffffffffffffffffffffffffff90811690831615610c9c5773ffffffffffffffffffffffffffffffffffffffff811660011415610c97576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f436c61696d20636f646520616c7265616479207370656e7400000000000000006044820152606401610650565b610d3e565b73ffffffffffffffffffffffffffffffffffffffff8116331480610cd4575073ffffffffffffffffffffffffffffffffffffffff8116155b610d3a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f436c61696d20636f646520616c7265616479207370656e7400000000000000006044820152606401610650565b3392505b505060409182013560009081526001602081905292902080547fffffffffffffffffffffffff00000000000000000000000000000000000000001690921790915590565b803573ffffffffffffffffffffffffffffffffffffffff81168114610da657600080fd5b919050565b80357fffff00000000000000000000000000000000000000000000000000000000000081168114610da657600080fd5b600060608284031215610ded57600080fd5b50919050565b600060208284031215610e0557600080fd5b610e0e82610d82565b9392505050565b60008060408385031215610e2857600080fd5b610e3183610d82565b946020939093013593505050565b60008060208385031215610e5257600080fd5b823567ffffffffffffffff80821115610e6a57600080fd5b818501915085601f830112610e7e57600080fd5b813581811115610e8d57600080fd5b8660208260051b8501011115610ea257600080fd5b60209290920196919550909350505050565b600060208284031215610ec657600080fd5b610e0e82610dab565b600060208284031215610ee157600080fd5b5035919050565b600080600060e08486031215610efd57600080fd5b610f078585610ddb565b9250610f168560608601610ddb565b9150610f2460c08501610dab565b90509250925092565b6000806000806101008587031215610f4457600080fd5b610f4e8686610ddb565b9350610f5d8660608701610ddb565b9250610f6b60c08601610dab565b9150610f7960e08601610d82565b905092959194509250565b600060208284031215610f9657600080fd5b813563ffffffff81168114610e0e57600080fd5b600060208284031215610fbc57600080fd5b813560ff81168114610e0e57600080fd5b6000835160005b81811015610fee5760208187018101518583015201610fd4565b81811115610ffd576000828501525b5060609390931b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000169190920190815260140192915050565b60008282101561104857611048611086565b500390565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82141561107f5761107f611086565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fdfea26469706673582212202cc95c42f0a51f091c5ebc99a835d7215c8d0dff3046a6ca534be0104e81178664736f6c63430008070033",
}

// ClaimCodesABI is the input ABI used to generate the binding from.
// Deprecated: Use ClaimCodesMetaData.ABI instead.
var ClaimCodesABI = ClaimCodesMetaData.ABI

// ClaimCodesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ClaimCodesMetaData.Bin instead.
var ClaimCodesBin = ClaimCodesMetaData.Bin

// DeployClaimCodes deploys a new Ethereum contract, binding an instance of ClaimCodes to it.
func DeployClaimCodes(auth *bind.TransactOpts, backend bind.ContractBackend, vesting_address common.Address) (common.Address, *types.Transaction, *ClaimCodes, error) {
	parsed, err := ClaimCodesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ClaimCodesBin), backend, vesting_address)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ClaimCodes{ClaimCodesCaller: ClaimCodesCaller{contract: contract}, ClaimCodesTransactor: ClaimCodesTransactor{contract: contract}, ClaimCodesFilterer: ClaimCodesFilterer{contract: contract}}, nil
}

// ClaimCodes is an auto generated Go binding around an Ethereum contract.
type ClaimCodes struct {
	ClaimCodesCaller     // Read-only binding to the contract
	ClaimCodesTransactor // Write-only binding to the contract
	ClaimCodesFilterer   // Log filterer for contract events
}

// ClaimCodesCaller is an auto generated read-only Go binding around an Ethereum contract.
type ClaimCodesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClaimCodesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ClaimCodesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClaimCodesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ClaimCodesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClaimCodesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ClaimCodesSession struct {
	Contract     *ClaimCodes       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ClaimCodesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ClaimCodesCallerSession struct {
	Contract *ClaimCodesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ClaimCodesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ClaimCodesTransactorSession struct {
	Contract     *ClaimCodesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ClaimCodesRaw is an auto generated low-level Go binding around an Ethereum contract.
type ClaimCodesRaw struct {
	Contract *ClaimCodes // Generic contract binding to access the raw methods on
}

// ClaimCodesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ClaimCodesCallerRaw struct {
	Contract *ClaimCodesCaller // Generic read-only contract binding to access the raw methods on
}

// ClaimCodesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ClaimCodesTransactorRaw struct {
	Contract *ClaimCodesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewClaimCodes creates a new instance of ClaimCodes, bound to a specific deployed contract.
func NewClaimCodes(address common.Address, backend bind.ContractBackend) (*ClaimCodes, error) {
	contract, err := bindClaimCodes(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ClaimCodes{ClaimCodesCaller: ClaimCodesCaller{contract: contract}, ClaimCodesTransactor: ClaimCodesTransactor{contract: contract}, ClaimCodesFilterer: ClaimCodesFilterer{contract: contract}}, nil
}

// NewClaimCodesCaller creates a new read-only instance of ClaimCodes, bound to a specific deployed contract.
func NewClaimCodesCaller(address common.Address, caller bind.ContractCaller) (*ClaimCodesCaller, error) {
	contract, err := bindClaimCodes(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ClaimCodesCaller{contract: contract}, nil
}

// NewClaimCodesTransactor creates a new write-only instance of ClaimCodes, bound to a specific deployed contract.
func NewClaimCodesTransactor(address common.Address, transactor bind.ContractTransactor) (*ClaimCodesTransactor, error) {
	contract, err := bindClaimCodes(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ClaimCodesTransactor{contract: contract}, nil
}

// NewClaimCodesFilterer creates a new log filterer instance of ClaimCodes, bound to a specific deployed contract.
func NewClaimCodesFilterer(address common.Address, filterer bind.ContractFilterer) (*ClaimCodesFilterer, error) {
	contract, err := bindClaimCodes(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ClaimCodesFilterer{contract: contract}, nil
}

// bindClaimCodes binds a generic wrapper to an already deployed contract.
func bindClaimCodes(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ClaimCodesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClaimCodes *ClaimCodesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClaimCodes.Contract.ClaimCodesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClaimCodes *ClaimCodesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClaimCodes.Contract.ClaimCodesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClaimCodes *ClaimCodesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClaimCodes.Contract.ClaimCodesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClaimCodes *ClaimCodesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClaimCodes.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClaimCodes *ClaimCodesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClaimCodes.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClaimCodes *ClaimCodesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClaimCodes.Contract.contract.Transact(opts, method, params...)
}

// AllowedCountries is a free data retrieval call binding the contract method 0x34b4947b.
//
// Solidity: function allowed_countries(bytes2 ) view returns(bool)
func (_ClaimCodes *ClaimCodesCaller) AllowedCountries(opts *bind.CallOpts, arg0 [2]byte) (bool, error) {
	var out []interface{}
	err := _ClaimCodes.contract.Call(opts, &out, "allowed_countries", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedCountries is a free data retrieval call binding the contract method 0x34b4947b.
//
// Solidity: function allowed_countries(bytes2 ) view returns(bool)
func (_ClaimCodes *ClaimCodesSession) AllowedCountries(arg0 [2]byte) (bool, error) {
	return _ClaimCodes.Contract.AllowedCountries(&_ClaimCodes.CallOpts, arg0)
}

// AllowedCountries is a free data retrieval call binding the contract method 0x34b4947b.
//
// Solidity: function allowed_countries(bytes2 ) view returns(bool)
func (_ClaimCodes *ClaimCodesCallerSession) AllowedCountries(arg0 [2]byte) (bool, error) {
	return _ClaimCodes.Contract.AllowedCountries(&_ClaimCodes.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(address)
func (_ClaimCodes *ClaimCodesCaller) Commitments(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ClaimCodes.contract.Call(opts, &out, "commitments", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(address)
func (_ClaimCodes *ClaimCodesSession) Commitments(arg0 [32]byte) (common.Address, error) {
	return _ClaimCodes.Contract.Commitments(&_ClaimCodes.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(address)
func (_ClaimCodes *ClaimCodesCallerSession) Commitments(arg0 [32]byte) (common.Address, error) {
	return _ClaimCodes.Contract.Commitments(&_ClaimCodes.CallOpts, arg0)
}

// Controller is a free data retrieval call binding the contract method 0xf77c4791.
//
// Solidity: function controller() view returns(address)
func (_ClaimCodes *ClaimCodesCaller) Controller(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ClaimCodes.contract.Call(opts, &out, "controller")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Controller is a free data retrieval call binding the contract method 0xf77c4791.
//
// Solidity: function controller() view returns(address)
func (_ClaimCodes *ClaimCodesSession) Controller() (common.Address, error) {
	return _ClaimCodes.Contract.Controller(&_ClaimCodes.CallOpts)
}

// Controller is a free data retrieval call binding the contract method 0xf77c4791.
//
// Solidity: function controller() view returns(address)
func (_ClaimCodes *ClaimCodesCallerSession) Controller() (common.Address, error) {
	return _ClaimCodes.Contract.Controller(&_ClaimCodes.CallOpts)
}

// Issuers is a free data retrieval call binding the contract method 0x38a7543e.
//
// Solidity: function issuers(address ) view returns(uint256)
func (_ClaimCodes *ClaimCodesCaller) Issuers(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ClaimCodes.contract.Call(opts, &out, "issuers", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Issuers is a free data retrieval call binding the contract method 0x38a7543e.
//
// Solidity: function issuers(address ) view returns(uint256)
func (_ClaimCodes *ClaimCodesSession) Issuers(arg0 common.Address) (*big.Int, error) {
	return _ClaimCodes.Contract.Issuers(&_ClaimCodes.CallOpts, arg0)
}

// Issuers is a free data retrieval call binding the contract method 0x38a7543e.
//
// Solidity: function issuers(address ) view returns(uint256)
func (_ClaimCodes *ClaimCodesCallerSession) Issuers(arg0 common.Address) (*big.Int, error) {
	return _ClaimCodes.Contract.Issuers(&_ClaimCodes.CallOpts, arg0)
}

// AllowCountries is a paid mutator transaction binding the contract method 0x9cef4355.
//
// Solidity: function allow_countries(bytes2[] countries) returns()
func (_ClaimCodes *ClaimCodesTransactor) AllowCountries(opts *bind.TransactOpts, countries [][2]byte) (*types.Transaction, error) {
	return _ClaimCodes.contract.Transact(opts, "allow_countries", countries)
}

// AllowCountries is a paid mutator transaction binding the contract method 0x9cef4355.
//
// Solidity: function allow_countries(bytes2[] countries) returns()
func (_ClaimCodes *ClaimCodesSession) AllowCountries(countries [][2]byte) (*types.Transaction, error) {
	return _ClaimCodes.Contract.AllowCountries(&_ClaimCodes.TransactOpts, countries)
}

// AllowCountries is a paid mutator transaction binding the contract method 0x9cef4355.
//
// Solidity: function allow_countries(bytes2[] countries) returns()
func (_ClaimCodes *ClaimCodesTransactorSession) AllowCountries(countries [][2]byte) (*types.Transaction, error) {
	return _ClaimCodes.Contract.AllowCountries(&_ClaimCodes.TransactOpts, countries)
}

// BlockCountries is a paid mutator transaction binding the contract method 0x26bc922a.
//
// Solidity: function block_countries(bytes2[] countries) returns()
func (_ClaimCodes *ClaimCodesTransactor) BlockCountries(opts *bind.TransactOpts, countries [][2]byte) (*types.Transaction, error) {
	return _ClaimCodes.contract.Transact(opts, "block_countries", countries)
}

// BlockCountries is a paid mutator transaction binding the contract method 0x26bc922a.
//
// Solidity: function block_countries(bytes2[] countries) returns()
func (_ClaimCodes *ClaimCodesSession) BlockCountries(countries [][2]byte) (*types.Transaction, error) {
	return _ClaimCodes.Contract.BlockCountries(&_ClaimCodes.TransactOpts, countries)
}

// BlockCountries is a paid mutator transaction binding the contract method 0x26bc922a.
//
// Solidity: function block_countries(bytes2[] countries) returns()
func (_ClaimCodes *ClaimCodesTransactorSession) BlockCountries(countries [][2]byte) (*types.Transaction, error) {
	return _ClaimCodes.Contract.BlockCountries(&_ClaimCodes.TransactOpts, countries)
}

// ClaimTargeted is a paid mutator transaction binding the contract method 0x04df9479.
//
// Solidity: function claim_targeted((uint8,bytes32,bytes32) sig, (uint256,uint8,uint32) clm, bytes2 country, address target) returns()
func (_ClaimCodes *ClaimCodesTransactor) ClaimTargeted(opts *bind.TransactOpts, sig Signature, clm Claim, country [2]byte, target common.Address) (*types.Transaction, error) {
	return _ClaimCodes.contract.Transact(opts, "claim_targeted", sig, clm, country, target)
}

// ClaimTargeted is a paid mutator transaction binding the contract method 0x04df9479.
//
// Solidity: function claim_targeted((uint8,bytes32,bytes32) sig, (uint256,uint8,uint32) clm, bytes2 country, address target) returns()
func (_ClaimCodes *ClaimCodesSession) ClaimTargeted(sig Signature, clm Claim, country [2]byte, target common.Address) (*types.Transaction, error) {
	return _ClaimCodes.Contract.ClaimTargeted(&_ClaimCodes.TransactOpts, sig, clm, country, target)
}

// ClaimTargeted is a paid mutator transaction binding the contract method 0x04df9479.
//
// Solidity: function claim_targeted((uint8,bytes32,bytes32) sig, (uint256,uint8,uint32) clm, bytes2 country, address target) returns()
func (_ClaimCodes *ClaimCodesTransactorSession) ClaimTargeted(sig Signature, clm Claim, country [2]byte, target common.Address) (*types.Transaction, error) {
	return _ClaimCodes.Contract.ClaimTargeted(&_ClaimCodes.TransactOpts, sig, clm, country, target)
}

// ClaimUntargeted is a paid mutator transaction binding the contract method 0xb7d00d0c.
//
// Solidity: function claim_untargeted((uint8,bytes32,bytes32) sig, (uint256,uint8,uint32) clm, bytes2 country) returns()
func (_ClaimCodes *ClaimCodesTransactor) ClaimUntargeted(opts *bind.TransactOpts, sig Signature, clm Claim, country [2]byte) (*types.Transaction, error) {
	return _ClaimCodes.contract.Transact(opts, "claim_untargeted", sig, clm, country)
}

// ClaimUntargeted is a paid mutator transaction binding the contract method 0xb7d00d0c.
//
// Solidity: function claim_untargeted((uint8,bytes32,bytes32) sig, (uint256,uint8,uint32) clm, bytes2 country) returns()
func (_ClaimCodes *ClaimCodesSession) ClaimUntargeted(sig Signature, clm Claim, country [2]byte) (*types.Transaction, error) {
	return _ClaimCodes.Contract.ClaimUntargeted(&_ClaimCodes.TransactOpts, sig, clm, country)
}

// ClaimUntargeted is a paid mutator transaction binding the contract method 0xb7d00d0c.
//
// Solidity: function claim_untargeted((uint8,bytes32,bytes32) sig, (uint256,uint8,uint32) clm, bytes2 country) returns()
func (_ClaimCodes *ClaimCodesTransactorSession) ClaimUntargeted(sig Signature, clm Claim, country [2]byte) (*types.Transaction, error) {
	return _ClaimCodes.Contract.ClaimUntargeted(&_ClaimCodes.TransactOpts, sig, clm, country)
}

// CommitUntargeted is a paid mutator transaction binding the contract method 0xc342cd9b.
//
// Solidity: function commit_untargeted(bytes32 s) returns()
func (_ClaimCodes *ClaimCodesTransactor) CommitUntargeted(opts *bind.TransactOpts, s [32]byte) (*types.Transaction, error) {
	return _ClaimCodes.contract.Transact(opts, "commit_untargeted", s)
}

// CommitUntargeted is a paid mutator transaction binding the contract method 0xc342cd9b.
//
// Solidity: function commit_untargeted(bytes32 s) returns()
func (_ClaimCodes *ClaimCodesSession) CommitUntargeted(s [32]byte) (*types.Transaction, error) {
	return _ClaimCodes.Contract.CommitUntargeted(&_ClaimCodes.TransactOpts, s)
}

// CommitUntargeted is a paid mutator transaction binding the contract method 0xc342cd9b.
//
// Solidity: function commit_untargeted(bytes32 s) returns()
func (_ClaimCodes *ClaimCodesTransactorSession) CommitUntargeted(s [32]byte) (*types.Transaction, error) {
	return _ClaimCodes.Contract.CommitUntargeted(&_ClaimCodes.TransactOpts, s)
}

// Destroy is a paid mutator transaction binding the contract method 0x83197ef0.
//
// Solidity: function destroy() returns()
func (_ClaimCodes *ClaimCodesTransactor) Destroy(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClaimCodes.contract.Transact(opts, "destroy")
}

// Destroy is a paid mutator transaction binding the contract method 0x83197ef0.
//
// Solidity: function destroy() returns()
func (_ClaimCodes *ClaimCodesSession) Destroy() (*types.Transaction, error) {
	return _ClaimCodes.Contract.Destroy(&_ClaimCodes.TransactOpts)
}

// Destroy is a paid mutator transaction binding the contract method 0x83197ef0.
//
// Solidity: function destroy() returns()
func (_ClaimCodes *ClaimCodesTransactorSession) Destroy() (*types.Transaction, error) {
	return _ClaimCodes.Contract.Destroy(&_ClaimCodes.TransactOpts)
}

// PermitIssuer is a paid mutator transaction binding the contract method 0x277817c7.
//
// Solidity: function permit_issuer(address issuer, uint256 amount) returns()
func (_ClaimCodes *ClaimCodesTransactor) PermitIssuer(opts *bind.TransactOpts, issuer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ClaimCodes.contract.Transact(opts, "permit_issuer", issuer, amount)
}

// PermitIssuer is a paid mutator transaction binding the contract method 0x277817c7.
//
// Solidity: function permit_issuer(address issuer, uint256 amount) returns()
func (_ClaimCodes *ClaimCodesSession) PermitIssuer(issuer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ClaimCodes.Contract.PermitIssuer(&_ClaimCodes.TransactOpts, issuer, amount)
}

// PermitIssuer is a paid mutator transaction binding the contract method 0x277817c7.
//
// Solidity: function permit_issuer(address issuer, uint256 amount) returns()
func (_ClaimCodes *ClaimCodesTransactorSession) PermitIssuer(issuer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ClaimCodes.Contract.PermitIssuer(&_ClaimCodes.TransactOpts, issuer, amount)
}

// RevokeIssuer is a paid mutator transaction binding the contract method 0x94a4ff3f.
//
// Solidity: function revoke_issuer(address issuer) returns()
func (_ClaimCodes *ClaimCodesTransactor) RevokeIssuer(opts *bind.TransactOpts, issuer common.Address) (*types.Transaction, error) {
	return _ClaimCodes.contract.Transact(opts, "revoke_issuer", issuer)
}

// RevokeIssuer is a paid mutator transaction binding the contract method 0x94a4ff3f.
//
// Solidity: function revoke_issuer(address issuer) returns()
func (_ClaimCodes *ClaimCodesSession) RevokeIssuer(issuer common.Address) (*types.Transaction, error) {
	return _ClaimCodes.Contract.RevokeIssuer(&_ClaimCodes.TransactOpts, issuer)
}

// RevokeIssuer is a paid mutator transaction binding the contract method 0x94a4ff3f.
//
// Solidity: function revoke_issuer(address issuer) returns()
func (_ClaimCodes *ClaimCodesTransactorSession) RevokeIssuer(issuer common.Address) (*types.Transaction, error) {
	return _ClaimCodes.Contract.RevokeIssuer(&_ClaimCodes.TransactOpts, issuer)
}

// SwapController is a paid mutator transaction binding the contract method 0x7eb1b5d7.
//
// Solidity: function swap_controller(address _controller) returns()
func (_ClaimCodes *ClaimCodesTransactor) SwapController(opts *bind.TransactOpts, _controller common.Address) (*types.Transaction, error) {
	return _ClaimCodes.contract.Transact(opts, "swap_controller", _controller)
}

// SwapController is a paid mutator transaction binding the contract method 0x7eb1b5d7.
//
// Solidity: function swap_controller(address _controller) returns()
func (_ClaimCodes *ClaimCodesSession) SwapController(_controller common.Address) (*types.Transaction, error) {
	return _ClaimCodes.Contract.SwapController(&_ClaimCodes.TransactOpts, _controller)
}

// SwapController is a paid mutator transaction binding the contract method 0x7eb1b5d7.
//
// Solidity: function swap_controller(address _controller) returns()
func (_ClaimCodes *ClaimCodesTransactorSession) SwapController(_controller common.Address) (*types.Transaction, error) {
	return _ClaimCodes.Contract.SwapController(&_ClaimCodes.TransactOpts, _controller)
}
