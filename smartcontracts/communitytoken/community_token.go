// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package communitytoken

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

// CommunitytokenMetaData contains all meta data concerning the Communitytoken contract.
var CommunitytokenMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]",
	Bin: "0x608060405234801561001057600080fd5b506004361061010b5760003560e01c806340c10f19116100a257806395d89b411161007157806395d89b41146102ce578063a457c2d7146102ec578063a9059cbb1461031c578063dd62ed3e1461034c578063de5f72fd1461037c5761010b565b806340c10f19146102485780634e82cd9e14610264578063673200a71461028257806370a082311461029e5761010b565b806318d3c6f6116100de57806318d3c6f61461019a57806323b872dd146101ca578063313ce567146101fa57806339509351146102185761010b565b806306fdde0314610110578063075461721461012e578063095ea7b31461014c57806318160ddd1461017c575b600080fd5b610118610386565b604051610125919061108a565b60405180910390f35b610136610418565b60405161014391906110ed565b60405180910390f35b6101666004803603810190610161919061116f565b61043e565b60405161017391906111ca565b60405180910390f35b610184610461565b60405161019191906111f4565b60405180910390f35b6101b460048036038101906101af919061120f565b61046b565b6040516101c191906111f4565b60405180910390f35b6101e460048036038101906101df919061123c565b610483565b6040516101f191906111ca565b60405180910390f35b6102026104b2565b60405161020f91906112ab565b60405180910390f35b610232600480360381019061022d919061116f565b6104c9565b60405161023f91906111ca565b60405180910390f35b610262600480360381019061025d919061116f565b610500565b005b61026c6105f9565b60405161027991906111f4565b60405180910390f35b61029c6004803603810190610297919061120f565b6105ff565b005b6102b860048036038101906102b3919061120f565b61069d565b6040516102c591906111f4565b60405180910390f35b6102d66106e5565b6040516102e3919061108a565b60405180910390f35b6103066004803603810190610301919061116f565b610777565b60405161031391906111ca565b60405180910390f35b6103366004803603810190610331919061116f565b6107ee565b60405161034391906111ca565b60405180910390f35b610366600480360381019061036191906112c6565b610811565b60405161037391906111f4565b60405180910390f35b610384610898565b005b60606003805461039590611335565b80601f01602080910402602001604051908101604052809291908181526020018280546103c190611335565b801561040e5780601f106103e35761010080835404028352916020019161040e565b820191906000526020600020905b8154815290600101906020018083116103f157829003601f168201915b5050505050905090565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000806104496109be565b90506104568185856109c6565b600191505092915050565b6000600254905090565b60086020528060005260406000206000915090505481565b60008061048e6109be565b905061049b858285610b8f565b6104a6858585610c1b565b60019150509392505050565b6000600760009054906101000a900460ff16905090565b6000806104d46109be565b90506104f58185856104e68589610811565b6104f09190611395565b6109c6565b600191505092915050565b600073ffffffffffffffffffffffffffffffffffffffff16600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361055b57600080fd5b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146105eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105e290611437565b60405180910390fd5b6105f58282610e91565b5050565b60065481565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461065957600080fd5b80600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6060600480546106f490611335565b80601f016020809104026020016040519081016040528092919081815260200182805461072090611335565b801561076d5780601f106107425761010080835404028352916020019161076d565b820191906000526020600020905b81548152906001019060200180831161075057829003601f168201915b5050505050905090565b6000806107826109be565b905060006107908286610811565b9050838110156107d5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107cc906114c9565b60405180910390fd5b6107e282868684036109c6565b60019250505092915050565b6000806107f96109be565b9050610806818585610c1b565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000600654116108dd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108d490611535565b60405180910390fd5b4262015180600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461092c9190611395565b1061096c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610963906115c7565b60405180910390fd5b42600860003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506109bc33600654610e91565b565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610a35576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a2c90611659565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610aa4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a9b906116eb565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610b8291906111f4565b60405180910390a3505050565b6000610b9b8484610811565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610c155781811015610c07576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bfe90611757565b60405180910390fd5b610c1484848484036109c6565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610c8a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c81906117e9565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610cf9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cf09061187b565b60405180910390fd5b610d04838383610fe7565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610d8a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d819061190d565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610e7891906111f4565b60405180910390a3610e8b848484610fec565b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610f00576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ef790611979565b60405180910390fd5b610f0c60008383610fe7565b8060026000828254610f1e9190611395565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610fcf91906111f4565b60405180910390a3610fe360008383610fec565b5050565b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561102b578082015181840152602081019050611010565b8381111561103a576000848401525b50505050565b6000601f19601f8301169050919050565b600061105c82610ff1565b6110668185610ffc565b935061107681856020860161100d565b61107f81611040565b840191505092915050565b600060208201905081810360008301526110a48184611051565b905092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006110d7826110ac565b9050919050565b6110e7816110cc565b82525050565b600060208201905061110260008301846110de565b92915050565b600080fd5b611116816110cc565b811461112157600080fd5b50565b6000813590506111338161110d565b92915050565b6000819050919050565b61114c81611139565b811461115757600080fd5b50565b60008135905061116981611143565b92915050565b6000806040838503121561118657611185611108565b5b600061119485828601611124565b92505060206111a58582860161115a565b9150509250929050565b60008115159050919050565b6111c4816111af565b82525050565b60006020820190506111df60008301846111bb565b92915050565b6111ee81611139565b82525050565b600060208201905061120960008301846111e5565b92915050565b60006020828403121561122557611224611108565b5b600061123384828501611124565b91505092915050565b60008060006060848603121561125557611254611108565b5b600061126386828701611124565b935050602061127486828701611124565b92505060406112858682870161115a565b9150509250925092565b600060ff82169050919050565b6112a58161128f565b82525050565b60006020820190506112c0600083018461129c565b92915050565b600080604083850312156112dd576112dc611108565b5b60006112eb85828601611124565b92505060206112fc85828601611124565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061134d57607f821691505b6020821081036113605761135f611306565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006113a082611139565b91506113ab83611139565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156113e0576113df611366565b5b828201905092915050565b7f757365722063616e6e6f74206d696e7400000000000000000000000000000000600082015250565b6000611421601083610ffc565b915061142c826113eb565b602082019050919050565b6000602082019050818103600083015261145081611414565b9050919050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b60006114b3602583610ffc565b91506114be82611457565b604082019050919050565b600060208201905081810360008301526114e2816114a6565b9050919050565b7f666175636574206e6f7420656e61626c65640000000000000000000000000000600082015250565b600061151f601283610ffc565b915061152a826114e9565b602082019050919050565b6000602082019050818103600083015261154e81611512565b9050919050565b7f6661756365742063616e2062652075736564206f6e636520706572203234206860008201527f6f75727300000000000000000000000000000000000000000000000000000000602082015250565b60006115b1602483610ffc565b91506115bc82611555565b604082019050919050565b600060208201905081810360008301526115e0816115a4565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b6000611643602483610ffc565b915061164e826115e7565b604082019050919050565b6000602082019050818103600083015261167281611636565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b60006116d5602283610ffc565b91506116e082611679565b604082019050919050565b60006020820190508181036000830152611704816116c8565b9050919050565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b6000611741601d83610ffc565b915061174c8261170b565b602082019050919050565b6000602082019050818103600083015261177081611734565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b60006117d3602583610ffc565b91506117de82611777565b604082019050919050565b60006020820190508181036000830152611802816117c6565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b6000611865602383610ffc565b915061187082611809565b604082019050919050565b6000602082019050818103600083015261189481611858565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b60006118f7602683610ffc565b91506119028261189b565b604082019050919050565b60006020820190508181036000830152611926816118ea565b9050919050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b6000611963601f83610ffc565b915061196e8261192d565b602082019050919050565b6000602082019050818103600083015261199281611956565b905091905056fea264697066735822122001750d524b4661590f49b2ec768864d852571732b55a7fc8308254d8f7ac20db64736f6c634300080d0033",
}

// CommunitytokenABI is the input ABI used to generate the binding from.
// Deprecated: Use CommunitytokenMetaData.ABI instead.
var CommunitytokenABI = CommunitytokenMetaData.ABI

// CommunitytokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CommunitytokenMetaData.Bin instead.
var CommunitytokenBin = CommunitytokenMetaData.Bin

// DeployCommunitytoken deploys a new Ethereum contract, binding an instance of Communitytoken to it.
func DeployCommunitytoken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Communitytoken, error) {
	parsed, err := CommunitytokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CommunitytokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Communitytoken{CommunitytokenCaller: CommunitytokenCaller{contract: contract}, CommunitytokenTransactor: CommunitytokenTransactor{contract: contract}, CommunitytokenFilterer: CommunitytokenFilterer{contract: contract}}, nil
}

// Communitytoken is an auto generated Go binding around an Ethereum contract.
type Communitytoken struct {
	CommunitytokenCaller     // Read-only binding to the contract
	CommunitytokenTransactor // Write-only binding to the contract
	CommunitytokenFilterer   // Log filterer for contract events
}

// CommunitytokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type CommunitytokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommunitytokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CommunitytokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommunitytokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CommunitytokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommunitytokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CommunitytokenSession struct {
	Contract     *Communitytoken   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CommunitytokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CommunitytokenCallerSession struct {
	Contract *CommunitytokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// CommunitytokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CommunitytokenTransactorSession struct {
	Contract     *CommunitytokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// CommunitytokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type CommunitytokenRaw struct {
	Contract *Communitytoken // Generic contract binding to access the raw methods on
}

// CommunitytokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CommunitytokenCallerRaw struct {
	Contract *CommunitytokenCaller // Generic read-only contract binding to access the raw methods on
}

// CommunitytokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CommunitytokenTransactorRaw struct {
	Contract *CommunitytokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCommunitytoken creates a new instance of Communitytoken, bound to a specific deployed contract.
func NewCommunitytoken(address common.Address, backend bind.ContractBackend) (*Communitytoken, error) {
	contract, err := bindCommunitytoken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Communitytoken{CommunitytokenCaller: CommunitytokenCaller{contract: contract}, CommunitytokenTransactor: CommunitytokenTransactor{contract: contract}, CommunitytokenFilterer: CommunitytokenFilterer{contract: contract}}, nil
}

// NewCommunitytokenCaller creates a new read-only instance of Communitytoken, bound to a specific deployed contract.
func NewCommunitytokenCaller(address common.Address, caller bind.ContractCaller) (*CommunitytokenCaller, error) {
	contract, err := bindCommunitytoken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CommunitytokenCaller{contract: contract}, nil
}

// NewCommunitytokenTransactor creates a new write-only instance of Communitytoken, bound to a specific deployed contract.
func NewCommunitytokenTransactor(address common.Address, transactor bind.ContractTransactor) (*CommunitytokenTransactor, error) {
	contract, err := bindCommunitytoken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CommunitytokenTransactor{contract: contract}, nil
}

// NewCommunitytokenFilterer creates a new log filterer instance of Communitytoken, bound to a specific deployed contract.
func NewCommunitytokenFilterer(address common.Address, filterer bind.ContractFilterer) (*CommunitytokenFilterer, error) {
	contract, err := bindCommunitytoken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CommunitytokenFilterer{contract: contract}, nil
}

// bindCommunitytoken binds a generic wrapper to an already deployed contract.
func bindCommunitytoken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CommunitytokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Communitytoken *CommunitytokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Communitytoken.Contract.CommunitytokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Communitytoken *CommunitytokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Communitytoken.Contract.CommunitytokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Communitytoken *CommunitytokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Communitytoken.Contract.CommunitytokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Communitytoken *CommunitytokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Communitytoken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Communitytoken *CommunitytokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Communitytoken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Communitytoken *CommunitytokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Communitytoken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Communitytoken *CommunitytokenCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Communitytoken.contract.Call(opts, &out, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Communitytoken *CommunitytokenSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Communitytoken.Contract.Allowance(&_Communitytoken.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Communitytoken *CommunitytokenCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Communitytoken.Contract.Allowance(&_Communitytoken.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 balance)
func (_Communitytoken *CommunitytokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Communitytoken.contract.Call(opts, &out, "balanceOf", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 balance)
func (_Communitytoken *CommunitytokenSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Communitytoken.Contract.BalanceOf(&_Communitytoken.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 balance)
func (_Communitytoken *CommunitytokenCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Communitytoken.Contract.BalanceOf(&_Communitytoken.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Communitytoken *CommunitytokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Communitytoken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Communitytoken *CommunitytokenSession) Decimals() (uint8, error) {
	return _Communitytoken.Contract.Decimals(&_Communitytoken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Communitytoken *CommunitytokenCallerSession) Decimals() (uint8, error) {
	return _Communitytoken.Contract.Decimals(&_Communitytoken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Communitytoken *CommunitytokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Communitytoken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Communitytoken *CommunitytokenSession) Name() (string, error) {
	return _Communitytoken.Contract.Name(&_Communitytoken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Communitytoken *CommunitytokenCallerSession) Name() (string, error) {
	return _Communitytoken.Contract.Name(&_Communitytoken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Communitytoken *CommunitytokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Communitytoken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Communitytoken *CommunitytokenSession) Symbol() (string, error) {
	return _Communitytoken.Contract.Symbol(&_Communitytoken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Communitytoken *CommunitytokenCallerSession) Symbol() (string, error) {
	return _Communitytoken.Contract.Symbol(&_Communitytoken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Communitytoken *CommunitytokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Communitytoken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Communitytoken *CommunitytokenSession) TotalSupply() (*big.Int, error) {
	return _Communitytoken.Contract.TotalSupply(&_Communitytoken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Communitytoken *CommunitytokenCallerSession) TotalSupply() (*big.Int, error) {
	return _Communitytoken.Contract.TotalSupply(&_Communitytoken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_Communitytoken *CommunitytokenTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Communitytoken.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_Communitytoken *CommunitytokenSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Communitytoken.Contract.Approve(&_Communitytoken.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_Communitytoken *CommunitytokenTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Communitytoken.Contract.Approve(&_Communitytoken.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_Communitytoken *CommunitytokenTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Communitytoken.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_Communitytoken *CommunitytokenSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Communitytoken.Contract.Transfer(&_Communitytoken.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_Communitytoken *CommunitytokenTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Communitytoken.Contract.Transfer(&_Communitytoken.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_Communitytoken *CommunitytokenTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Communitytoken.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_Communitytoken *CommunitytokenSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Communitytoken.Contract.TransferFrom(&_Communitytoken.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool)
func (_Communitytoken *CommunitytokenTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Communitytoken.Contract.TransferFrom(&_Communitytoken.TransactOpts, _from, _to, _value)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Communitytoken *CommunitytokenTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Communitytoken.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Communitytoken *CommunitytokenSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Communitytoken.Contract.Fallback(&_Communitytoken.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Communitytoken *CommunitytokenTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Communitytoken.Contract.Fallback(&_Communitytoken.TransactOpts, calldata)
}

// CommunitytokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Communitytoken contract.
type CommunitytokenApprovalIterator struct {
	Event *CommunitytokenApproval // Event containing the contract specifics and raw log

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
func (it *CommunitytokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommunitytokenApproval)
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
		it.Event = new(CommunitytokenApproval)
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
func (it *CommunitytokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommunitytokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommunitytokenApproval represents a Approval event raised by the Communitytoken contract.
type CommunitytokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Communitytoken *CommunitytokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*CommunitytokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Communitytoken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &CommunitytokenApprovalIterator{contract: _Communitytoken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Communitytoken *CommunitytokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CommunitytokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Communitytoken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommunitytokenApproval)
				if err := _Communitytoken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Communitytoken *CommunitytokenFilterer) ParseApproval(log types.Log) (*CommunitytokenApproval, error) {
	event := new(CommunitytokenApproval)
	if err := _Communitytoken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommunitytokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Communitytoken contract.
type CommunitytokenTransferIterator struct {
	Event *CommunitytokenTransfer // Event containing the contract specifics and raw log

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
func (it *CommunitytokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommunitytokenTransfer)
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
		it.Event = new(CommunitytokenTransfer)
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
func (it *CommunitytokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommunitytokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommunitytokenTransfer represents a Transfer event raised by the Communitytoken contract.
type CommunitytokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Communitytoken *CommunitytokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CommunitytokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Communitytoken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CommunitytokenTransferIterator{contract: _Communitytoken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Communitytoken *CommunitytokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CommunitytokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Communitytoken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommunitytokenTransfer)
				if err := _Communitytoken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Communitytoken *CommunitytokenFilterer) ParseTransfer(log types.Log) (*CommunitytokenTransfer, error) {
	event := new(CommunitytokenTransfer)
	if err := _Communitytoken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
