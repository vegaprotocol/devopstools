// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package MultisigControl

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

// MultisigControlMetaData contains all meta data concerning the MultisigControl contract.
var MultisigControlMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"NonceBurnt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"new_signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"SignerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old_signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"SignerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"new_threshold\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"ThresholdSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"new_signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"add_signer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"burn_nonce\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_current_threshold\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_valid_signer_count\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"is_nonce_used\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer_address\",\"type\":\"address\"}],\"name\":\"is_valid_signer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"old_signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"remove_signer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"new_threshold\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_threshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"signers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"verify_signatures\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506000805461ffff19166101f41781553381526001602081905260408220805460ff19169091179055805462010000900460ff16906002610050836100ab565b825460ff9182166101009390930a92830291909202199091161790555060408051338152600060208201527f50999ebf9b59bf3157a58816611976f2d723378ad51457d7b0413209e0cdee59910160405180910390a16100d9565b600060ff821660ff8114156100d057634e487b7160e01b600052601160045260246000fd5b60010192915050565b61137d806100e86000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c806398c5f73e11610076578063ba73659a1161005b578063ba73659a146101b0578063dbe528df146101c3578063f8e3a660146101d957600080fd5b806398c5f73e1461017f578063b04e3dd11461019257600080fd5b80635ec51639116100a75780635ec51639146101105780635f06155914610123578063736c0d5b1461015c57600080fd5b806350ac8df8146100c35780635b9fe26b146100d8575b600080fd5b6100d66100d1366004610e98565b6101ec565b005b6100fb6100e6366004610efb565b60009081526002602052604090205460ff1690565b60405190151581526020015b60405180910390f35b6100d661011e366004610f14565b6103b2565b6100fb610131366004610f89565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604090205460ff1690565b6100fb61016a366004610f89565b60016020526000908152604090205460ff1681565b6100d661018d366004610fa4565b6104b8565b60005462010000900460ff1660405160ff9091168152602001610107565b6100fb6101be366004610ff2565b610713565b60005460405161ffff9091168152602001610107565b6100d66101e7366004610fa4565b610b69565b6103e88461ffff16108015610205575060008461ffff16115b610270576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f6e6577207468726573686f6c64206f7574736964652072616e6765000000000060448201526064015b60405180910390fd5b6040805161ffff86166020820152908101849052606080820152600d60808201527f7365745f7468726573686f6c640000000000000000000000000000000000000060a082015260009060c00160405160208183030381529060405290506102da83838387610713565b610340576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610267565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00001661ffff871690811790915560408051918252602082018690527ff6d24c23627520a3b70e5dc66aa1249844b4bb407c2c153d9000a2b14a1e3c1191015b60405180910390a15050505050565b6000836040516020016103fc918152604060208201819052600a908201527f6275726e5f6e6f6e636500000000000000000000000000000000000000000000606082015260800190565b604051602081830303815290604052905061041983838387610713565b61047f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610267565b6040518481527fb33a7fc220f9e1c644c0f616b48edee1956a978a7dcb37a10f16e148969e4c0b9060200160405180910390a150505050565b6040805173ffffffffffffffffffffffffffffffffffffffff86166020820152908101849052606080820152600d60808201527f72656d6f76655f7369676e65720000000000000000000000000000000000000060a082015260009060c001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291815273ffffffffffffffffffffffffffffffffffffffff871660009081526001602052205490915060ff166105d4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f7369676e657220646f65736e27742065786973740000000000000000000000006044820152606401610267565b6105e083838387610713565b610646576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610267565b73ffffffffffffffffffffffffffffffffffffffff8516600090815260016020526040812080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055805462010000900460ff169060026106a883611122565b91906101000a81548160ff021916908360ff160217905550507f99c1d2c0ed8107e4db2e5dbfb10a2549cd2a63cbe39cf99d2adffbcd0395441885856040516103a392919073ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b600061072060418561118c565b15610787576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f62616420736967206c656e6774680000000000000000000000000000000000006044820152606401610267565b836107ee576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f6d75737420636f6e7461696e206174206c6561737420312073696700000000006044820152606401610267565b60008281526002602052604090205460ff1615610867576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6e6f6e636520616c7265616479207573656400000000000000000000000000006044820152606401610267565b60008054819062010000900460ff1667ffffffffffffffff81111561088e5761088e610fc3565b6040519080825280602002602001820160405280156108b7578160200160208202803683370190505b509050600085336040516020016108cf9291906111a0565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012090508760005b88811015610aef578181018035906020810135906040013560001a7f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08211156109b0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f4d616c6c6561626c65207369676e6174757265206572726f72000000000000006044820152606401610267565b601b8160ff1610156109ca576109c7601b82611231565b90505b6040805160008082526020820180845289905260ff841692820192909252606081018590526080810184905260019060a0016020604051602081039080840390855afa158015610a1e573d6000803e3d6000fd5b5050604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015173ffffffffffffffffffffffffffffffffffffffff811660009081526001602052919091205490925060ff1690508015610a895750610a8788828b610dca565b155b15610ad75780888a60ff1681518110610aa457610aa4611256565b73ffffffffffffffffffffffffffffffffffffffff9092166020928302919091019091015288610ad381611285565b9950505b50505050604181610ae891906112a5565b905061090b565b5060005461ffff81169060ff62010000909104811690610b139087166103e86112bd565b610b1d91906112fa565b600088815260026020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001692909111918217905560ff16945050505050949350505050565b6040805173ffffffffffffffffffffffffffffffffffffffff86166020820152908101849052606080820152600a60808201527f6164645f7369676e65720000000000000000000000000000000000000000000060a082015260009060c001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815291815273ffffffffffffffffffffffffffffffffffffffff871660009081526001602052205490915060ff1615610c86576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f7369676e657220616c72656164792065786973747300000000000000000000006044820152606401610267565b610c9283838387610713565b610cf8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610267565b73ffffffffffffffffffffffffffffffffffffffff85166000908152600160208190526040822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169091179055805462010000900460ff16906002610d5f83611285565b91906101000a81548160ff021916908360ff160217905550507f50999ebf9b59bf3157a58816611976f2d723378ad51457d7b0413209e0cdee5985856040516103a392919073ffffffffffffffffffffffffffffffffffffffff929092168252602082015260400190565b6000805b8260ff16811015610e42578373ffffffffffffffffffffffffffffffffffffffff16858281518110610e0257610e02611256565b602002602001015173ffffffffffffffffffffffffffffffffffffffff161415610e30576001915050610e48565b80610e3a8161130e565b915050610dce565b50600090505b9392505050565b60008083601f840112610e6157600080fd5b50813567ffffffffffffffff811115610e7957600080fd5b602083019150836020828501011115610e9157600080fd5b9250929050565b60008060008060608587031215610eae57600080fd5b843561ffff81168114610ec057600080fd5b935060208501359250604085013567ffffffffffffffff811115610ee357600080fd5b610eef87828801610e4f565b95989497509550505050565b600060208284031215610f0d57600080fd5b5035919050565b600080600060408486031215610f2957600080fd5b83359250602084013567ffffffffffffffff811115610f4757600080fd5b610f5386828701610e4f565b9497909650939450505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610f8457600080fd5b919050565b600060208284031215610f9b57600080fd5b610e4882610f60565b60008060008060608587031215610fba57600080fd5b610ec085610f60565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806000806060858703121561100857600080fd5b843567ffffffffffffffff8082111561102057600080fd5b61102c88838901610e4f565b9096509450602087013591508082111561104557600080fd5b818701915087601f83011261105957600080fd5b81358181111561106b5761106b610fc3565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156110b1576110b1610fc3565b816040528281528a60208487010111156110ca57600080fd5b826020860160208301376000928101602001929092525095989497509495604001359450505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600060ff821680611135576111356110f3565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261119b5761119b61115d565b500690565b604081526000835180604084015260005b818110156111ce57602081870181015160608684010152016111b1565b818111156111e0576000606083860101525b5073ffffffffffffffffffffffffffffffffffffffff93909316602083015250601f919091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601606001919050565b600060ff821660ff84168060ff0382111561124e5761124e6110f3565b019392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600060ff821660ff81141561129c5761129c6110f3565b60010192915050565b600082198211156112b8576112b86110f3565b500190565b6000817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156112f5576112f56110f3565b500290565b6000826113095761130961115d565b500490565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611340576113406110f3565b506001019056fea2646970667358221220723331279716e624d896dfec3f042f3c520b43f5df2bd24562c0040f40cfc4df64736f6c63430008080033",
}

// MultisigControlABI is the input ABI used to generate the binding from.
// Deprecated: Use MultisigControlMetaData.ABI instead.
var MultisigControlABI = MultisigControlMetaData.ABI

// MultisigControlBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MultisigControlMetaData.Bin instead.
var MultisigControlBin = MultisigControlMetaData.Bin

// DeployMultisigControl deploys a new Ethereum contract, binding an instance of MultisigControl to it.
func DeployMultisigControl(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MultisigControl, error) {
	parsed, err := MultisigControlMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MultisigControlBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultisigControl{MultisigControlCaller: MultisigControlCaller{contract: contract}, MultisigControlTransactor: MultisigControlTransactor{contract: contract}, MultisigControlFilterer: MultisigControlFilterer{contract: contract}}, nil
}

// MultisigControl is an auto generated Go binding around an Ethereum contract.
type MultisigControl struct {
	MultisigControlCaller     // Read-only binding to the contract
	MultisigControlTransactor // Write-only binding to the contract
	MultisigControlFilterer   // Log filterer for contract events
}

// MultisigControlCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultisigControlCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigControlTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultisigControlTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigControlFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultisigControlFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigControlSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultisigControlSession struct {
	Contract     *MultisigControl  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MultisigControlCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultisigControlCallerSession struct {
	Contract *MultisigControlCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// MultisigControlTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultisigControlTransactorSession struct {
	Contract     *MultisigControlTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// MultisigControlRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultisigControlRaw struct {
	Contract *MultisigControl // Generic contract binding to access the raw methods on
}

// MultisigControlCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultisigControlCallerRaw struct {
	Contract *MultisigControlCaller // Generic read-only contract binding to access the raw methods on
}

// MultisigControlTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultisigControlTransactorRaw struct {
	Contract *MultisigControlTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultisigControl creates a new instance of MultisigControl, bound to a specific deployed contract.
func NewMultisigControl(address common.Address, backend bind.ContractBackend) (*MultisigControl, error) {
	contract, err := bindMultisigControl(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultisigControl{MultisigControlCaller: MultisigControlCaller{contract: contract}, MultisigControlTransactor: MultisigControlTransactor{contract: contract}, MultisigControlFilterer: MultisigControlFilterer{contract: contract}}, nil
}

// NewMultisigControlCaller creates a new read-only instance of MultisigControl, bound to a specific deployed contract.
func NewMultisigControlCaller(address common.Address, caller bind.ContractCaller) (*MultisigControlCaller, error) {
	contract, err := bindMultisigControl(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultisigControlCaller{contract: contract}, nil
}

// NewMultisigControlTransactor creates a new write-only instance of MultisigControl, bound to a specific deployed contract.
func NewMultisigControlTransactor(address common.Address, transactor bind.ContractTransactor) (*MultisigControlTransactor, error) {
	contract, err := bindMultisigControl(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultisigControlTransactor{contract: contract}, nil
}

// NewMultisigControlFilterer creates a new log filterer instance of MultisigControl, bound to a specific deployed contract.
func NewMultisigControlFilterer(address common.Address, filterer bind.ContractFilterer) (*MultisigControlFilterer, error) {
	contract, err := bindMultisigControl(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultisigControlFilterer{contract: contract}, nil
}

// bindMultisigControl binds a generic wrapper to an already deployed contract.
func bindMultisigControl(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultisigControlABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultisigControl *MultisigControlRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultisigControl.Contract.MultisigControlCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultisigControl *MultisigControlRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultisigControl.Contract.MultisigControlTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultisigControl *MultisigControlRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultisigControl.Contract.MultisigControlTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultisigControl *MultisigControlCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultisigControl.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultisigControl *MultisigControlTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultisigControl.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultisigControl *MultisigControlTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultisigControl.Contract.contract.Transact(opts, method, params...)
}

// GetCurrentThreshold is a free data retrieval call binding the contract method 0xdbe528df.
//
// Solidity: function get_current_threshold() view returns(uint16)
func (_MultisigControl *MultisigControlCaller) GetCurrentThreshold(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _MultisigControl.contract.Call(opts, &out, "get_current_threshold")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetCurrentThreshold is a free data retrieval call binding the contract method 0xdbe528df.
//
// Solidity: function get_current_threshold() view returns(uint16)
func (_MultisigControl *MultisigControlSession) GetCurrentThreshold() (uint16, error) {
	return _MultisigControl.Contract.GetCurrentThreshold(&_MultisigControl.CallOpts)
}

// GetCurrentThreshold is a free data retrieval call binding the contract method 0xdbe528df.
//
// Solidity: function get_current_threshold() view returns(uint16)
func (_MultisigControl *MultisigControlCallerSession) GetCurrentThreshold() (uint16, error) {
	return _MultisigControl.Contract.GetCurrentThreshold(&_MultisigControl.CallOpts)
}

// GetValidSignerCount is a free data retrieval call binding the contract method 0xb04e3dd1.
//
// Solidity: function get_valid_signer_count() view returns(uint8)
func (_MultisigControl *MultisigControlCaller) GetValidSignerCount(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MultisigControl.contract.Call(opts, &out, "get_valid_signer_count")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetValidSignerCount is a free data retrieval call binding the contract method 0xb04e3dd1.
//
// Solidity: function get_valid_signer_count() view returns(uint8)
func (_MultisigControl *MultisigControlSession) GetValidSignerCount() (uint8, error) {
	return _MultisigControl.Contract.GetValidSignerCount(&_MultisigControl.CallOpts)
}

// GetValidSignerCount is a free data retrieval call binding the contract method 0xb04e3dd1.
//
// Solidity: function get_valid_signer_count() view returns(uint8)
func (_MultisigControl *MultisigControlCallerSession) GetValidSignerCount() (uint8, error) {
	return _MultisigControl.Contract.GetValidSignerCount(&_MultisigControl.CallOpts)
}

// IsNonceUsed is a free data retrieval call binding the contract method 0x5b9fe26b.
//
// Solidity: function is_nonce_used(uint256 nonce) view returns(bool)
func (_MultisigControl *MultisigControlCaller) IsNonceUsed(opts *bind.CallOpts, nonce *big.Int) (bool, error) {
	var out []interface{}
	err := _MultisigControl.contract.Call(opts, &out, "is_nonce_used", nonce)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsNonceUsed is a free data retrieval call binding the contract method 0x5b9fe26b.
//
// Solidity: function is_nonce_used(uint256 nonce) view returns(bool)
func (_MultisigControl *MultisigControlSession) IsNonceUsed(nonce *big.Int) (bool, error) {
	return _MultisigControl.Contract.IsNonceUsed(&_MultisigControl.CallOpts, nonce)
}

// IsNonceUsed is a free data retrieval call binding the contract method 0x5b9fe26b.
//
// Solidity: function is_nonce_used(uint256 nonce) view returns(bool)
func (_MultisigControl *MultisigControlCallerSession) IsNonceUsed(nonce *big.Int) (bool, error) {
	return _MultisigControl.Contract.IsNonceUsed(&_MultisigControl.CallOpts, nonce)
}

// IsValidSigner is a free data retrieval call binding the contract method 0x5f061559.
//
// Solidity: function is_valid_signer(address signer_address) view returns(bool)
func (_MultisigControl *MultisigControlCaller) IsValidSigner(opts *bind.CallOpts, signer_address common.Address) (bool, error) {
	var out []interface{}
	err := _MultisigControl.contract.Call(opts, &out, "is_valid_signer", signer_address)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidSigner is a free data retrieval call binding the contract method 0x5f061559.
//
// Solidity: function is_valid_signer(address signer_address) view returns(bool)
func (_MultisigControl *MultisigControlSession) IsValidSigner(signer_address common.Address) (bool, error) {
	return _MultisigControl.Contract.IsValidSigner(&_MultisigControl.CallOpts, signer_address)
}

// IsValidSigner is a free data retrieval call binding the contract method 0x5f061559.
//
// Solidity: function is_valid_signer(address signer_address) view returns(bool)
func (_MultisigControl *MultisigControlCallerSession) IsValidSigner(signer_address common.Address) (bool, error) {
	return _MultisigControl.Contract.IsValidSigner(&_MultisigControl.CallOpts, signer_address)
}

// Signers is a free data retrieval call binding the contract method 0x736c0d5b.
//
// Solidity: function signers(address ) view returns(bool)
func (_MultisigControl *MultisigControlCaller) Signers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _MultisigControl.contract.Call(opts, &out, "signers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Signers is a free data retrieval call binding the contract method 0x736c0d5b.
//
// Solidity: function signers(address ) view returns(bool)
func (_MultisigControl *MultisigControlSession) Signers(arg0 common.Address) (bool, error) {
	return _MultisigControl.Contract.Signers(&_MultisigControl.CallOpts, arg0)
}

// Signers is a free data retrieval call binding the contract method 0x736c0d5b.
//
// Solidity: function signers(address ) view returns(bool)
func (_MultisigControl *MultisigControlCallerSession) Signers(arg0 common.Address) (bool, error) {
	return _MultisigControl.Contract.Signers(&_MultisigControl.CallOpts, arg0)
}

// AddSigner is a paid mutator transaction binding the contract method 0xf8e3a660.
//
// Solidity: function add_signer(address new_signer, uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlTransactor) AddSigner(opts *bind.TransactOpts, new_signer common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.contract.Transact(opts, "add_signer", new_signer, nonce, signatures)
}

// AddSigner is a paid mutator transaction binding the contract method 0xf8e3a660.
//
// Solidity: function add_signer(address new_signer, uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlSession) AddSigner(new_signer common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.Contract.AddSigner(&_MultisigControl.TransactOpts, new_signer, nonce, signatures)
}

// AddSigner is a paid mutator transaction binding the contract method 0xf8e3a660.
//
// Solidity: function add_signer(address new_signer, uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlTransactorSession) AddSigner(new_signer common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.Contract.AddSigner(&_MultisigControl.TransactOpts, new_signer, nonce, signatures)
}

// BurnNonce is a paid mutator transaction binding the contract method 0x5ec51639.
//
// Solidity: function burn_nonce(uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlTransactor) BurnNonce(opts *bind.TransactOpts, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.contract.Transact(opts, "burn_nonce", nonce, signatures)
}

// BurnNonce is a paid mutator transaction binding the contract method 0x5ec51639.
//
// Solidity: function burn_nonce(uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlSession) BurnNonce(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.Contract.BurnNonce(&_MultisigControl.TransactOpts, nonce, signatures)
}

// BurnNonce is a paid mutator transaction binding the contract method 0x5ec51639.
//
// Solidity: function burn_nonce(uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlTransactorSession) BurnNonce(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.Contract.BurnNonce(&_MultisigControl.TransactOpts, nonce, signatures)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x98c5f73e.
//
// Solidity: function remove_signer(address old_signer, uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlTransactor) RemoveSigner(opts *bind.TransactOpts, old_signer common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.contract.Transact(opts, "remove_signer", old_signer, nonce, signatures)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x98c5f73e.
//
// Solidity: function remove_signer(address old_signer, uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlSession) RemoveSigner(old_signer common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.Contract.RemoveSigner(&_MultisigControl.TransactOpts, old_signer, nonce, signatures)
}

// RemoveSigner is a paid mutator transaction binding the contract method 0x98c5f73e.
//
// Solidity: function remove_signer(address old_signer, uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlTransactorSession) RemoveSigner(old_signer common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.Contract.RemoveSigner(&_MultisigControl.TransactOpts, old_signer, nonce, signatures)
}

// SetThreshold is a paid mutator transaction binding the contract method 0x50ac8df8.
//
// Solidity: function set_threshold(uint16 new_threshold, uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlTransactor) SetThreshold(opts *bind.TransactOpts, new_threshold uint16, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.contract.Transact(opts, "set_threshold", new_threshold, nonce, signatures)
}

// SetThreshold is a paid mutator transaction binding the contract method 0x50ac8df8.
//
// Solidity: function set_threshold(uint16 new_threshold, uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlSession) SetThreshold(new_threshold uint16, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.Contract.SetThreshold(&_MultisigControl.TransactOpts, new_threshold, nonce, signatures)
}

// SetThreshold is a paid mutator transaction binding the contract method 0x50ac8df8.
//
// Solidity: function set_threshold(uint16 new_threshold, uint256 nonce, bytes signatures) returns()
func (_MultisigControl *MultisigControlTransactorSession) SetThreshold(new_threshold uint16, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _MultisigControl.Contract.SetThreshold(&_MultisigControl.TransactOpts, new_threshold, nonce, signatures)
}

// VerifySignatures is a paid mutator transaction binding the contract method 0xba73659a.
//
// Solidity: function verify_signatures(bytes signatures, bytes message, uint256 nonce) returns(bool)
func (_MultisigControl *MultisigControlTransactor) VerifySignatures(opts *bind.TransactOpts, signatures []byte, message []byte, nonce *big.Int) (*types.Transaction, error) {
	return _MultisigControl.contract.Transact(opts, "verify_signatures", signatures, message, nonce)
}

// VerifySignatures is a paid mutator transaction binding the contract method 0xba73659a.
//
// Solidity: function verify_signatures(bytes signatures, bytes message, uint256 nonce) returns(bool)
func (_MultisigControl *MultisigControlSession) VerifySignatures(signatures []byte, message []byte, nonce *big.Int) (*types.Transaction, error) {
	return _MultisigControl.Contract.VerifySignatures(&_MultisigControl.TransactOpts, signatures, message, nonce)
}

// VerifySignatures is a paid mutator transaction binding the contract method 0xba73659a.
//
// Solidity: function verify_signatures(bytes signatures, bytes message, uint256 nonce) returns(bool)
func (_MultisigControl *MultisigControlTransactorSession) VerifySignatures(signatures []byte, message []byte, nonce *big.Int) (*types.Transaction, error) {
	return _MultisigControl.Contract.VerifySignatures(&_MultisigControl.TransactOpts, signatures, message, nonce)
}

// MultisigControlNonceBurntIterator is returned from FilterNonceBurnt and is used to iterate over the raw logs and unpacked data for NonceBurnt events raised by the MultisigControl contract.
type MultisigControlNonceBurntIterator struct {
	Event *MultisigControlNonceBurnt // Event containing the contract specifics and raw log

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
func (it *MultisigControlNonceBurntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigControlNonceBurnt)
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
		it.Event = new(MultisigControlNonceBurnt)
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
func (it *MultisigControlNonceBurntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigControlNonceBurntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigControlNonceBurnt represents a NonceBurnt event raised by the MultisigControl contract.
type MultisigControlNonceBurnt struct {
	Nonce *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterNonceBurnt is a free log retrieval operation binding the contract event 0xb33a7fc220f9e1c644c0f616b48edee1956a978a7dcb37a10f16e148969e4c0b.
//
// Solidity: event NonceBurnt(uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) FilterNonceBurnt(opts *bind.FilterOpts) (*MultisigControlNonceBurntIterator, error) {

	logs, sub, err := _MultisigControl.contract.FilterLogs(opts, "NonceBurnt")
	if err != nil {
		return nil, err
	}
	return &MultisigControlNonceBurntIterator{contract: _MultisigControl.contract, event: "NonceBurnt", logs: logs, sub: sub}, nil
}

// WatchNonceBurnt is a free log subscription operation binding the contract event 0xb33a7fc220f9e1c644c0f616b48edee1956a978a7dcb37a10f16e148969e4c0b.
//
// Solidity: event NonceBurnt(uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) WatchNonceBurnt(opts *bind.WatchOpts, sink chan<- *MultisigControlNonceBurnt) (event.Subscription, error) {

	logs, sub, err := _MultisigControl.contract.WatchLogs(opts, "NonceBurnt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigControlNonceBurnt)
				if err := _MultisigControl.contract.UnpackLog(event, "NonceBurnt", log); err != nil {
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

// ParseNonceBurnt is a log parse operation binding the contract event 0xb33a7fc220f9e1c644c0f616b48edee1956a978a7dcb37a10f16e148969e4c0b.
//
// Solidity: event NonceBurnt(uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) ParseNonceBurnt(log types.Log) (*MultisigControlNonceBurnt, error) {
	event := new(MultisigControlNonceBurnt)
	if err := _MultisigControl.contract.UnpackLog(event, "NonceBurnt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultisigControlSignerAddedIterator is returned from FilterSignerAdded and is used to iterate over the raw logs and unpacked data for SignerAdded events raised by the MultisigControl contract.
type MultisigControlSignerAddedIterator struct {
	Event *MultisigControlSignerAdded // Event containing the contract specifics and raw log

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
func (it *MultisigControlSignerAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigControlSignerAdded)
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
		it.Event = new(MultisigControlSignerAdded)
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
func (it *MultisigControlSignerAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigControlSignerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigControlSignerAdded represents a SignerAdded event raised by the MultisigControl contract.
type MultisigControlSignerAdded struct {
	NewSigner common.Address
	Nonce     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignerAdded is a free log retrieval operation binding the contract event 0x50999ebf9b59bf3157a58816611976f2d723378ad51457d7b0413209e0cdee59.
//
// Solidity: event SignerAdded(address new_signer, uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) FilterSignerAdded(opts *bind.FilterOpts) (*MultisigControlSignerAddedIterator, error) {

	logs, sub, err := _MultisigControl.contract.FilterLogs(opts, "SignerAdded")
	if err != nil {
		return nil, err
	}
	return &MultisigControlSignerAddedIterator{contract: _MultisigControl.contract, event: "SignerAdded", logs: logs, sub: sub}, nil
}

// WatchSignerAdded is a free log subscription operation binding the contract event 0x50999ebf9b59bf3157a58816611976f2d723378ad51457d7b0413209e0cdee59.
//
// Solidity: event SignerAdded(address new_signer, uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) WatchSignerAdded(opts *bind.WatchOpts, sink chan<- *MultisigControlSignerAdded) (event.Subscription, error) {

	logs, sub, err := _MultisigControl.contract.WatchLogs(opts, "SignerAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigControlSignerAdded)
				if err := _MultisigControl.contract.UnpackLog(event, "SignerAdded", log); err != nil {
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

// ParseSignerAdded is a log parse operation binding the contract event 0x50999ebf9b59bf3157a58816611976f2d723378ad51457d7b0413209e0cdee59.
//
// Solidity: event SignerAdded(address new_signer, uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) ParseSignerAdded(log types.Log) (*MultisigControlSignerAdded, error) {
	event := new(MultisigControlSignerAdded)
	if err := _MultisigControl.contract.UnpackLog(event, "SignerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultisigControlSignerRemovedIterator is returned from FilterSignerRemoved and is used to iterate over the raw logs and unpacked data for SignerRemoved events raised by the MultisigControl contract.
type MultisigControlSignerRemovedIterator struct {
	Event *MultisigControlSignerRemoved // Event containing the contract specifics and raw log

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
func (it *MultisigControlSignerRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigControlSignerRemoved)
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
		it.Event = new(MultisigControlSignerRemoved)
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
func (it *MultisigControlSignerRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigControlSignerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigControlSignerRemoved represents a SignerRemoved event raised by the MultisigControl contract.
type MultisigControlSignerRemoved struct {
	OldSigner common.Address
	Nonce     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignerRemoved is a free log retrieval operation binding the contract event 0x99c1d2c0ed8107e4db2e5dbfb10a2549cd2a63cbe39cf99d2adffbcd03954418.
//
// Solidity: event SignerRemoved(address old_signer, uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) FilterSignerRemoved(opts *bind.FilterOpts) (*MultisigControlSignerRemovedIterator, error) {

	logs, sub, err := _MultisigControl.contract.FilterLogs(opts, "SignerRemoved")
	if err != nil {
		return nil, err
	}
	return &MultisigControlSignerRemovedIterator{contract: _MultisigControl.contract, event: "SignerRemoved", logs: logs, sub: sub}, nil
}

// WatchSignerRemoved is a free log subscription operation binding the contract event 0x99c1d2c0ed8107e4db2e5dbfb10a2549cd2a63cbe39cf99d2adffbcd03954418.
//
// Solidity: event SignerRemoved(address old_signer, uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) WatchSignerRemoved(opts *bind.WatchOpts, sink chan<- *MultisigControlSignerRemoved) (event.Subscription, error) {

	logs, sub, err := _MultisigControl.contract.WatchLogs(opts, "SignerRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigControlSignerRemoved)
				if err := _MultisigControl.contract.UnpackLog(event, "SignerRemoved", log); err != nil {
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

// ParseSignerRemoved is a log parse operation binding the contract event 0x99c1d2c0ed8107e4db2e5dbfb10a2549cd2a63cbe39cf99d2adffbcd03954418.
//
// Solidity: event SignerRemoved(address old_signer, uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) ParseSignerRemoved(log types.Log) (*MultisigControlSignerRemoved, error) {
	event := new(MultisigControlSignerRemoved)
	if err := _MultisigControl.contract.UnpackLog(event, "SignerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultisigControlThresholdSetIterator is returned from FilterThresholdSet and is used to iterate over the raw logs and unpacked data for ThresholdSet events raised by the MultisigControl contract.
type MultisigControlThresholdSetIterator struct {
	Event *MultisigControlThresholdSet // Event containing the contract specifics and raw log

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
func (it *MultisigControlThresholdSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigControlThresholdSet)
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
		it.Event = new(MultisigControlThresholdSet)
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
func (it *MultisigControlThresholdSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigControlThresholdSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigControlThresholdSet represents a ThresholdSet event raised by the MultisigControl contract.
type MultisigControlThresholdSet struct {
	NewThreshold uint16
	Nonce        *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterThresholdSet is a free log retrieval operation binding the contract event 0xf6d24c23627520a3b70e5dc66aa1249844b4bb407c2c153d9000a2b14a1e3c11.
//
// Solidity: event ThresholdSet(uint16 new_threshold, uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) FilterThresholdSet(opts *bind.FilterOpts) (*MultisigControlThresholdSetIterator, error) {

	logs, sub, err := _MultisigControl.contract.FilterLogs(opts, "ThresholdSet")
	if err != nil {
		return nil, err
	}
	return &MultisigControlThresholdSetIterator{contract: _MultisigControl.contract, event: "ThresholdSet", logs: logs, sub: sub}, nil
}

// WatchThresholdSet is a free log subscription operation binding the contract event 0xf6d24c23627520a3b70e5dc66aa1249844b4bb407c2c153d9000a2b14a1e3c11.
//
// Solidity: event ThresholdSet(uint16 new_threshold, uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) WatchThresholdSet(opts *bind.WatchOpts, sink chan<- *MultisigControlThresholdSet) (event.Subscription, error) {

	logs, sub, err := _MultisigControl.contract.WatchLogs(opts, "ThresholdSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigControlThresholdSet)
				if err := _MultisigControl.contract.UnpackLog(event, "ThresholdSet", log); err != nil {
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

// ParseThresholdSet is a log parse operation binding the contract event 0xf6d24c23627520a3b70e5dc66aa1249844b4bb407c2c153d9000a2b14a1e3c11.
//
// Solidity: event ThresholdSet(uint16 new_threshold, uint256 nonce)
func (_MultisigControl *MultisigControlFilterer) ParseThresholdSet(log types.Log) (*MultisigControlThresholdSet, error) {
	event := new(MultisigControlThresholdSet)
	if err := _MultisigControl.contract.UnpackLog(event, "ThresholdSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
