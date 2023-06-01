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
	_ = abi.ConvertType
)

// MultisigControlMetaData contains all meta data concerning the MultisigControl contract.
var MultisigControlMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"new_signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"SignerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old_signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"SignerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"new_threshold\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"ThresholdSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"new_signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"add_signer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_current_threshold\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_valid_signer_count\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"is_nonce_used\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"signer_address\",\"type\":\"address\"}],\"name\":\"is_valid_signer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"old_signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"remove_signer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"new_threshold\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_threshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"verify_signatures\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506101f46000806101000a81548161ffff021916908361ffff16021790555060018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506000600281819054906101000a900460ff1680929190620000a8906200013f565b91906101000a81548160ff021916908360ff160217905550507f50999ebf9b59bf3157a58816611976f2d723378ad51457d7b0413209e0cdee59336000604051620000f59291906200020a565b60405180910390a162000237565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600060ff82169050919050565b60006200014c8262000132565b915060ff82141562000163576200016262000103565b5b600182019050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200019b826200016e565b9050919050565b620001ad816200018e565b82525050565b6000819050919050565b6000819050919050565b6000819050919050565b6000620001f2620001ec620001e684620001b3565b620001c7565b620001bd565b9050919050565b6200020481620001d1565b82525050565b6000604082019050620002216000830185620001a2565b620002306020830184620001f9565b9392505050565b61180880620002476000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063b04e3dd11161005b578063b04e3dd114610125578063ba73659a14610143578063dbe528df14610173578063f8e3a6601461019157610088565b806350ac8df81461008d5780635b9fe26b146100a95780635f061559146100d957806398c5f73e14610109575b600080fd5b6100a760048036038101906100a29190610bac565b6101ad565b005b6100c360048036038101906100be9190610c20565b6102d3565b6040516100d09190610c68565b60405180910390f35b6100f360048036038101906100ee9190610ce1565b6102fd565b6040516101009190610c68565b60405180910390f35b610123600480360381019061011e9190610d0e565b610353565b005b61012d610520565b60405161013a9190610d9e565b60405180910390f35b61015d60048036038101906101589190610efa565b610536565b60405161016a9190610c68565b60405180910390f35b61017b6108df565b6040516101889190610f99565b60405180910390f35b6101ab60048036038101906101a69190610d0e565b6108f6565b005b6103e88461ffff16111580156101c7575060008461ffff16115b610206576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101fd90611011565b60405180910390fd5b6000848460405160200161021b92919061108c565b604051602081830303815290604052905061023883838387610536565b610277576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161026e90611114565b60405180910390fd5b846000806101000a81548161ffff021916908361ffff1602179055507ff6d24c23627520a3b70e5dc66aa1249844b4bb407c2c153d9000a2b14a1e3c1185856040516102c4929190611134565b60405180910390a15050505050565b60006002600083815260200190815260200160002060009054906101000a900460ff169050919050565b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050919050565b600084846040516020016103689291906111b8565b6040516020818303038152906040529050600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610405576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103fc90611240565b60405180910390fd5b61041183838387610536565b610450576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161044790611114565b60405180910390fd5b6000600160008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506000600281819054906101000a900460ff16809291906104c79061128f565b91906101000a81548160ff021916908360ff160217905550507f99c1d2c0ed8107e4db2e5dbfb10a2549cd2a63cbe39cf99d2adffbcd0395441885856040516105119291906112b9565b60405180910390a15050505050565b60008060029054906101000a900460ff16905090565b6000806041868690506105499190611311565b14610589576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105809061138e565b60405180910390fd5b6002600083815260200190815260200160002060009054906101000a900460ff16156105ea576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105e1906113fa565b60405180910390fd5b60008084336040516020016106009291906114a2565b60405160208183030381529060405280519060200120905060005b87879050811015610864576000806000838b013592506020848c01013591506040848c01013560001a90507f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08260001c11156106ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106a39061151e565b60405180910390fd5b601b8160ff1610156106c857601b816106c5919061153e565b90505b6000600186838686604051600081526020016040526040516106ed949392919061158e565b6020604051602081039080840390855afa15801561070f573d6000803e3d6000fd5b505050602060405103519050600160008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1680156107cf57506003600087815260200190815260200160002060008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16155b1561084c5760016003600088815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508680610848906115d3565b9750505b5050505060418161085d91906115fd565b905061061b565b5060016002600086815260200190815260200160002060006101000a81548160ff02191690831515021790555060008054906101000a900461ffff1661ffff16600060029054906101000a900460ff1660ff166103e88460ff166108c89190611653565b6108d291906116ad565b1192505050949350505050565b60008060009054906101000a900461ffff16905090565b6000848460405160200161090b92919061172a565b6040516020818303038152906040529050600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16156109a9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109a0906117b2565b60405180910390fd5b6109b583838387610536565b6109f4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109eb90611114565b60405180910390fd5b60018060008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506000600281819054906101000a900460ff1680929190610a6a906115d3565b91906101000a81548160ff021916908360ff160217905550507f50999ebf9b59bf3157a58816611976f2d723378ad51457d7b0413209e0cdee598585604051610ab49291906112b9565b60405180910390a15050505050565b6000604051905090565b600080fd5b600080fd5b600061ffff82169050919050565b610aee81610ad7565b8114610af957600080fd5b50565b600081359050610b0b81610ae5565b92915050565b6000819050919050565b610b2481610b11565b8114610b2f57600080fd5b50565b600081359050610b4181610b1b565b92915050565b600080fd5b600080fd5b600080fd5b60008083601f840112610b6c57610b6b610b47565b5b8235905067ffffffffffffffff811115610b8957610b88610b4c565b5b602083019150836001820283011115610ba557610ba4610b51565b5b9250929050565b60008060008060608587031215610bc657610bc5610acd565b5b6000610bd487828801610afc565b9450506020610be587828801610b32565b935050604085013567ffffffffffffffff811115610c0657610c05610ad2565b5b610c1287828801610b56565b925092505092959194509250565b600060208284031215610c3657610c35610acd565b5b6000610c4484828501610b32565b91505092915050565b60008115159050919050565b610c6281610c4d565b82525050565b6000602082019050610c7d6000830184610c59565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610cae82610c83565b9050919050565b610cbe81610ca3565b8114610cc957600080fd5b50565b600081359050610cdb81610cb5565b92915050565b600060208284031215610cf757610cf6610acd565b5b6000610d0584828501610ccc565b91505092915050565b60008060008060608587031215610d2857610d27610acd565b5b6000610d3687828801610ccc565b9450506020610d4787828801610b32565b935050604085013567ffffffffffffffff811115610d6857610d67610ad2565b5b610d7487828801610b56565b925092505092959194509250565b600060ff82169050919050565b610d9881610d82565b82525050565b6000602082019050610db36000830184610d8f565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610e0782610dbe565b810181811067ffffffffffffffff82111715610e2657610e25610dcf565b5b80604052505050565b6000610e39610ac3565b9050610e458282610dfe565b919050565b600067ffffffffffffffff821115610e6557610e64610dcf565b5b610e6e82610dbe565b9050602081019050919050565b82818337600083830152505050565b6000610e9d610e9884610e4a565b610e2f565b905082815260208101848484011115610eb957610eb8610db9565b5b610ec4848285610e7b565b509392505050565b600082601f830112610ee157610ee0610b47565b5b8135610ef1848260208601610e8a565b91505092915050565b60008060008060608587031215610f1457610f13610acd565b5b600085013567ffffffffffffffff811115610f3257610f31610ad2565b5b610f3e87828801610b56565b9450945050602085013567ffffffffffffffff811115610f6157610f60610ad2565b5b610f6d87828801610ecc565b9250506040610f7e87828801610b32565b91505092959194509250565b610f9381610ad7565b82525050565b6000602082019050610fae6000830184610f8a565b92915050565b600082825260208201905092915050565b7f6e6577207468726573686f6c64206f7574736964652072616e67650000000000600082015250565b6000610ffb601b83610fb4565b915061100682610fc5565b602082019050919050565b6000602082019050818103600083015261102a81610fee565b9050919050565b61103a81610b11565b82525050565b7f7365745f7468726573686f6c6400000000000000000000000000000000000000600082015250565b6000611076600d83610fb4565b915061108182611040565b602082019050919050565b60006060820190506110a16000830185610f8a565b6110ae6020830184611031565b81810360408301526110bf81611069565b90509392505050565b7f626164207369676e617475726573000000000000000000000000000000000000600082015250565b60006110fe600e83610fb4565b9150611109826110c8565b602082019050919050565b6000602082019050818103600083015261112d816110f1565b9050919050565b60006040820190506111496000830185610f8a565b6111566020830184611031565b9392505050565b61116681610ca3565b82525050565b7f72656d6f76655f7369676e657200000000000000000000000000000000000000600082015250565b60006111a2600d83610fb4565b91506111ad8261116c565b602082019050919050565b60006060820190506111cd600083018561115d565b6111da6020830184611031565b81810360408301526111eb81611195565b90509392505050565b7f7369676e657220646f65736e2774206578697374000000000000000000000000600082015250565b600061122a601483610fb4565b9150611235826111f4565b602082019050919050565b600060208201905081810360008301526112598161121d565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061129a82610d82565b915060008214156112ae576112ad611260565b5b600182039050919050565b60006040820190506112ce600083018561115d565b6112db6020830184611031565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600061131c82610b11565b915061132783610b11565b925082611337576113366112e2565b5b828206905092915050565b7f62616420736967206c656e677468000000000000000000000000000000000000600082015250565b6000611378600e83610fb4565b915061138382611342565b602082019050919050565b600060208201905081810360008301526113a78161136b565b9050919050565b7f6e6f6e636520616c726561647920757365640000000000000000000000000000600082015250565b60006113e4601283610fb4565b91506113ef826113ae565b602082019050919050565b60006020820190508181036000830152611413816113d7565b9050919050565b600081519050919050565b600082825260208201905092915050565b60005b83811015611454578082015181840152602081019050611439565b83811115611463576000848401525b50505050565b60006114748261141a565b61147e8185611425565b935061148e818560208601611436565b61149781610dbe565b840191505092915050565b600060408201905081810360008301526114bc8185611469565b90506114cb602083018461115d565b9392505050565b7f4d616c6c61626c65207369676e6174757265206572726f720000000000000000600082015250565b6000611508601883610fb4565b9150611513826114d2565b602082019050919050565b60006020820190508181036000830152611537816114fb565b9050919050565b600061154982610d82565b915061155483610d82565b92508260ff0382111561156a57611569611260565b5b828201905092915050565b6000819050919050565b61158881611575565b82525050565b60006080820190506115a3600083018761157f565b6115b06020830186610d8f565b6115bd604083018561157f565b6115ca606083018461157f565b95945050505050565b60006115de82610d82565b915060ff8214156115f2576115f1611260565b5b600182019050919050565b600061160882610b11565b915061161383610b11565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561164857611647611260565b5b828201905092915050565b600061165e82610b11565b915061166983610b11565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156116a2576116a1611260565b5b828202905092915050565b60006116b882610b11565b91506116c383610b11565b9250826116d3576116d26112e2565b5b828204905092915050565b7f6164645f7369676e657200000000000000000000000000000000000000000000600082015250565b6000611714600a83610fb4565b915061171f826116de565b602082019050919050565b600060608201905061173f600083018561115d565b61174c6020830184611031565b818103604083015261175d81611707565b90509392505050565b7f7369676e657220616c7265616479206578697374730000000000000000000000600082015250565b600061179c601583610fb4565b91506117a782611766565b602082019050919050565b600060208201905081810360008301526117cb8161178f565b905091905056fea2646970667358221220fb9265433efee6e553a6ea1fa62168c751c6afcb346f6c8347f5063fc369ece064736f6c63430008080033",
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
	parsed, err := MultisigControlMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
