// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ERC20BridgeRestricted

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

// ERC20BridgeRestrictedMetaData contains all meta data concerning the ERC20BridgeRestricted contract.
var ERC20BridgeRestrictedMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"erc20_asset_pool\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user_address\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"Asset_Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lifetime_limit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdraw_threshold\",\"type\":\"uint256\"}],\"name\":\"Asset_Limits_Updated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"vega_asset_id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Listed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Removed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user_address\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Withdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Bridge_Resumed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Bridge_Stopped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdraw_delay\",\"type\":\"uint256\"}],\"name\":\"Bridge_Withdraw_Delay_Set\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"Depositor_Exempted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"Depositor_Exemption_Revoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"default_withdraw_delay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"deposit_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"erc20_asset_pool_address\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"exempt_depositor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"get_asset_deposit_lifetime_limit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"vega_asset_id\",\"type\":\"bytes32\"}],\"name\":\"get_asset_source\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_multisig_control_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"get_vega_asset_id\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"get_withdraw_threshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"global_resume\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"global_stop\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"is_asset_listed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"}],\"name\":\"is_exempt_depositor\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"is_stopped\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"vega_asset_id\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"lifetime_limit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdraw_threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"list_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"remove_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"revoke_exempt_depositor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"lifetime_limit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_asset_limits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delay\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_withdraw_delay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"creation\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"withdraw_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052620697806006553480156200001857600080fd5b50604051620023ba380380620023ba8339810160408190526200003b91620000bc565b6001600160a01b038116620000965760405162461bcd60e51b815260206004820152601a60248201527f696e76616c696420617373657420706f6f6c2061646472657373000000000000604482015260640160405180910390fd5b600080546001600160a01b0319166001600160a01b0392909216919091179055620000ee565b600060208284031215620000cf57600080fd5b81516001600160a01b0381168114620000e757600080fd5b9392505050565b6122bc80620000fe6000396000f3fe608060405234801561001057600080fd5b506004361061016c5760003560e01c80639356aab8116100cd578063c76de35811610081578063e272e9d011610066578063e272e9d014610395578063e8a7bce0146103a2578063f7683932146103d857600080fd5b8063c76de3581461036f578063d72ed5291461038257600080fd5b8063a06b5d39116100b2578063a06b5d3914610329578063b76fbb751461035f578063c58dc3b91461036757600080fd5b80639356aab8146102f65780639dfd3c881461031657600080fd5b806341fb776d116101245780636a1c6fa4116101095780636a1c6fa41461025a578063786b0bc0146102625780637fd27b7f146102bd57600080fd5b806341fb776d146102345780635a2467281461024757600080fd5b8063354a897a11610155578063354a897a146101d45780633ad90635146102185780633f4f199d1461022b57600080fd5b80630ff3562c1461017157806315c0df9d14610186575b600080fd5b61018461017f366004611dce565b6103eb565b005b6101bf610194366004611e43565b73ffffffffffffffffffffffffffffffffffffffff1660009081526009602052604090205460ff1690565b60405190151581526020015b60405180910390f35b61020a6101e2366004611e43565b73ffffffffffffffffffffffffffffffffffffffff1660009081526005602052604090205490565b6040519081526020016101cb565b610184610226366004611e67565b610762565b61020a60065481565b610184610242366004611f15565b610ae0565b610184610255366004611f88565b610d61565b610184610ef9565b610298610270366004611fdb565b60009081526002602052604090205473ffffffffffffffffffffffffffffffffffffffff1690565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101cb565b6101bf6102cb366004611e43565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604090205460ff1690565b6000546102989073ffffffffffffffffffffffffffffffffffffffff1681565b610184610324366004611ff4565b610fcf565b61020a610337366004611e43565b73ffffffffffffffffffffffffffffffffffffffff1660009081526003602052604090205490565b6101846111ef565b6102986112c9565b61018461037d366004612040565b6112d8565b610184610390366004611ff4565b61155e565b6008546101bf9060ff1681565b61020a6103b0366004611e43565b73ffffffffffffffffffffffffffffffffffffffff1660009081526007602052604090205490565b6101846103e6366004612099565b61177a565b73ffffffffffffffffffffffffffffffffffffffff861661046d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f696e76616c696420617373657420736f7572636500000000000000000000000060448201526064015b60405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff861660009081526001602052604090205460ff16156104fd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f617373657420616c7265616479206c69737465640000000000000000000000006044820152606401610464565b6040805173ffffffffffffffffffffffffffffffffffffffff88166020820152908101869052606081018590526080810184905260a0810183905260c080820152600a60e08201527f6c6973745f61737365740000000000000000000000000000000000000000000061010082015260009061012001604051602081830303815290604052905061058c611c2e565b73ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b81526004016105c893929190612148565b602060405180830381600087803b1580156105e257600080fd5b505af11580156105f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061061a919061217e565b610680576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610464565b73ffffffffffffffffffffffffffffffffffffffff8716600081815260016020818152604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169093179092558983526002815281832080547fffffffffffffffffffffffff00000000000000000000000000000000000000001685179055838352600381528183208a9055600581528183208990556007815291819020879055518581528892917f4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e91015b60405180910390a350505050505050565b60085460ff16156107cf576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f6272696467652073746f707065640000000000000000000000000000000000006044820152606401610464565b73ffffffffffffffffffffffffffffffffffffffff861660009081526007602052604090205485108061080f5750426006548461080c91906121a0565b11155b610875576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f6c61726765207769746864726177206973206e6f74206f6c6420656e6f7567686044820152606401610464565b6040805173ffffffffffffffffffffffffffffffffffffffff808916602083015291810187905290851660608201526080810184905260a0810183905260c080820152600e60e08201527f77697468647261775f6173736574000000000000000000000000000000000000610100820152600090610120016040516020818303038152906040529050610906611c2e565b73ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b815260040161094293929190612148565b602060405180830381600087803b15801561095c57600080fd5b505af1158015610970573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610994919061217e565b6109fa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610464565b6000546040517fd9caed1200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff89811660048301528781166024830152604482018990529091169063d9caed1290606401600060405180830381600087803b158015610a7657600080fd5b505af1158015610a8a573d6000803e3d6000fd5b5050604080518981526020810187905273ffffffffffffffffffffffffffffffffffffffff808c169450891692507fa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df9101610751565b73ffffffffffffffffffffffffffffffffffffffff861660009081526001602052604090205460ff16610b6f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6173736574206e6f74206c6973746564000000000000000000000000000000006044820152606401610464565b6040805173ffffffffffffffffffffffffffffffffffffffff88166020820152908101869052606081018590526080810184905260a080820152601060c08201527f7365745f61737365745f6c696d6974730000000000000000000000000000000060e0820152600090610100016040516020818303038152906040529050610bf6611c2e565b73ffffffffffffffffffffffffffffffffffffffff1663ba73659a848484886040518563ffffffff1660e01b8152600401610c3494939291906121df565b602060405180830381600087803b158015610c4e57600080fd5b505af1158015610c62573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c86919061217e565b610cec576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610464565b73ffffffffffffffffffffffffffffffffffffffff871660008181526005602090815260408083208a9055600782529182902088905581518981529081018890527ffc7eab762b8751ad85c101fd1025c763b4e8d48f2093f506629b606618e884fe910160405180910390a250505050505050565b6040805160208101869052908101849052606080820152601260808201527f7365745f77697468647261775f64656c6179000000000000000000000000000060a082015260009060c0016040516020818303038152906040529050610dc4611c2e565b73ffffffffffffffffffffffffffffffffffffffff1663ba73659a848484886040518563ffffffff1660e01b8152600401610e0294939291906121df565b602060405180830381600087803b158015610e1c57600080fd5b505af1158015610e30573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e54919061217e565b610eba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610464565b60068590556040518581527f1c7e8f73a01b8af4e18dd34455a42a45ad742bdb79cfda77bbdf50db2391fc889060200160405180910390a15050505050565b3360009081526009602052604090205460ff16610f72576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f73656e646572206e6f74206578656d70740000000000000000000000000000006044820152606401610464565b3360008181526009602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517fe74b113dca87276d976f476a9b4b9da3c780a3262eaabad051ee4e98912936a49190a2565b60085460ff161561103c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f62726964676520616c72656164792073746f70706564000000000000000000006044820152606401610464565b600083604051602001611086918152604060208201819052600b908201527f676c6f62616c5f73746f70000000000000000000000000000000000000000000606082015260800190565b604051602081830303815290604052905061109f611c2e565b73ffffffffffffffffffffffffffffffffffffffff1663ba73659a848484886040518563ffffffff1660e01b81526004016110dd94939291906121df565b602060405180830381600087803b1580156110f757600080fd5b505af115801561110b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061112f919061217e565b611195576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610464565b600880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556040517f129d99581c8e70519df1f0733d3212f33d0ed3ea6144adacc336c647f1d3638290600090a150505050565b3360009081526009602052604090205460ff1615611269576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f73656e64657220616c7265616479206578656d707400000000000000000000006044820152606401610464565b3360008181526009602052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517ff56e0868b913034a60dbca9c89ee79f8b0fa18dadbc5f6665f2f9a2cf3f51cdb9190a2565b60006112d3611c2e565b905090565b73ffffffffffffffffffffffffffffffffffffffff831660009081526001602052604090205460ff16611367576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6173736574206e6f74206c6973746564000000000000000000000000000000006044820152606401610464565b6040805173ffffffffffffffffffffffffffffffffffffffff85166020820152908101839052606080820152600c60808201527f72656d6f76655f6173736574000000000000000000000000000000000000000060a082015260009060c00160405160208183030381529060405290506113df611c2e565b73ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b815260040161141b93929190612148565b602060405180830381600087803b15801561143557600080fd5b505af1158015611449573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061146d919061217e565b6114d3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610464565b73ffffffffffffffffffffffffffffffffffffffff84166000818152600160205260409081902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517f58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c9906115509086815260200190565b60405180910390a250505050565b60085460ff166115ca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f627269646765206e6f742073746f7070656400000000000000000000000000006044820152606401610464565b600083604051602001611614918152604060208201819052600d908201527f676c6f62616c5f726573756d6500000000000000000000000000000000000000606082015260800190565b604051602081830303815290604052905061162d611c2e565b73ffffffffffffffffffffffffffffffffffffffff1663ba73659a848484886040518563ffffffff1660e01b815260040161166b94939291906121df565b602060405180830381600087803b15801561168557600080fd5b505af1158015611699573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116bd919061217e565b611723576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f626164207369676e6174757265730000000000000000000000000000000000006044820152606401610464565b600880547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556040517f79c02b0e60e0f00fe0370791204f2f175fe3f06f4816f3506ad4fa1b8e8cde0f90600090a150505050565b60085460ff16156117e7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f6272696467652073746f707065640000000000000000000000000000000000006044820152606401610464565b3360009081526009602052604090205460ff166118f45773ffffffffffffffffffffffffffffffffffffffff831660008181526005602090815260408083205433845260048352818420948452939091529020546118469084906121a0565b11156118ae576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601b60248201527f6465706f736974206f766572206c69666574696d65206c696d697400000000006044820152606401610464565b33600090815260046020908152604080832073ffffffffffffffffffffffffffffffffffffffff87168452909152812080548492906118ee9084906121a0565b90915550505b73ffffffffffffffffffffffffffffffffffffffff831660009081526001602052604090205460ff16611983576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6173736574206e6f74206c6973746564000000000000000000000000000000006044820152606401610464565b823b6119eb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f61737365745f736f75726365206d75737420626520636f6e74726163740000006044820152606401610464565b6000805460405133602482015273ffffffffffffffffffffffffffffffffffffffff9182166044820152606481018590528291861690608401604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529181526020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f23b872dd0000000000000000000000000000000000000000000000000000000017905251611aa5919061224d565b6000604051808303816000865af19150503d8060008114611ae2576040519150601f19603f3d011682016040523d82523d6000602084013e611ae7565b606091505b509150915081611b53576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f746f6b656e207472616e73666572206661696c656400000000000000000000006044820152606401610464565b805115611bd45780806020019051810190611b6e919061217e565b611bd4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f746f6b656e207472616e73666572206661696c656400000000000000000000006044820152606401610464565b604080518581526020810185905273ffffffffffffffffffffffffffffffffffffffff87169133917f3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf910160405180910390a35050505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663b82d5abd6040518163ffffffff1660e01b815260040160206040518083038186803b158015611c9757600080fd5b505afa158015611cab573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112d39190612269565b73ffffffffffffffffffffffffffffffffffffffff81168114611cf157600080fd5b50565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f830112611d3457600080fd5b813567ffffffffffffffff80821115611d4f57611d4f611cf4565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908282118183101715611d9557611d95611cf4565b81604052838152866020858801011115611dae57600080fd5b836020870160208301376000602085830101528094505050505092915050565b60008060008060008060c08789031215611de757600080fd5b8635611df281611ccf565b95506020870135945060408701359350606087013592506080870135915060a087013567ffffffffffffffff811115611e2a57600080fd5b611e3689828a01611d23565b9150509295509295509295565b600060208284031215611e5557600080fd5b8135611e6081611ccf565b9392505050565b60008060008060008060c08789031215611e8057600080fd5b8635611e8b81611ccf565b9550602087013594506040870135611ea281611ccf565b9350606087013592506080870135915060a087013567ffffffffffffffff811115611e2a57600080fd5b60008083601f840112611ede57600080fd5b50813567ffffffffffffffff811115611ef657600080fd5b602083019150836020828501011115611f0e57600080fd5b9250929050565b60008060008060008060a08789031215611f2e57600080fd5b8635611f3981611ccf565b9550602087013594506040870135935060608701359250608087013567ffffffffffffffff811115611f6a57600080fd5b611f7689828a01611ecc565b979a9699509497509295939492505050565b60008060008060608587031215611f9e57600080fd5b8435935060208501359250604085013567ffffffffffffffff811115611fc357600080fd5b611fcf87828801611ecc565b95989497509550505050565b600060208284031215611fed57600080fd5b5035919050565b60008060006040848603121561200957600080fd5b83359250602084013567ffffffffffffffff81111561202757600080fd5b61203386828701611ecc565b9497909650939450505050565b60008060006060848603121561205557600080fd5b833561206081611ccf565b925060208401359150604084013567ffffffffffffffff81111561208357600080fd5b61208f86828701611d23565b9150509250925092565b6000806000606084860312156120ae57600080fd5b83356120b981611ccf565b95602085013595506040909401359392505050565b60005b838110156120e95781810151838201526020016120d1565b838111156120f8576000848401525b50505050565b600081518084526121168160208601602086016120ce565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60608152600061215b60608301866120fe565b828103602084015261216d81866120fe565b915050826040830152949350505050565b60006020828403121561219057600080fd5b81518015158114611e6057600080fd5b600082198211156121da577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b500190565b606081528360608201528385608083013760006080858301015260007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8601168201608083820301602084015261223b60808201866120fe565b91505082604083015295945050505050565b6000825161225f8184602087016120ce565b9190910192915050565b60006020828403121561227b57600080fd5b8151611e6081611ccf56fea2646970667358221220564040cd5f92f777f0b532760585171690d4896a828abf9437b37135b0636cdb64736f6c63430008080033",
}

// ERC20BridgeRestrictedABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20BridgeRestrictedMetaData.ABI instead.
var ERC20BridgeRestrictedABI = ERC20BridgeRestrictedMetaData.ABI

// ERC20BridgeRestrictedBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC20BridgeRestrictedMetaData.Bin instead.
var ERC20BridgeRestrictedBin = ERC20BridgeRestrictedMetaData.Bin

// DeployERC20BridgeRestricted deploys a new Ethereum contract, binding an instance of ERC20BridgeRestricted to it.
func DeployERC20BridgeRestricted(auth *bind.TransactOpts, backend bind.ContractBackend, erc20_asset_pool common.Address) (common.Address, *types.Transaction, *ERC20BridgeRestricted, error) {
	parsed, err := ERC20BridgeRestrictedMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC20BridgeRestrictedBin), backend, erc20_asset_pool)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20BridgeRestricted{ERC20BridgeRestrictedCaller: ERC20BridgeRestrictedCaller{contract: contract}, ERC20BridgeRestrictedTransactor: ERC20BridgeRestrictedTransactor{contract: contract}, ERC20BridgeRestrictedFilterer: ERC20BridgeRestrictedFilterer{contract: contract}}, nil
}

// ERC20BridgeRestricted is an auto generated Go binding around an Ethereum contract.
type ERC20BridgeRestricted struct {
	ERC20BridgeRestrictedCaller     // Read-only binding to the contract
	ERC20BridgeRestrictedTransactor // Write-only binding to the contract
	ERC20BridgeRestrictedFilterer   // Log filterer for contract events
}

// ERC20BridgeRestrictedCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20BridgeRestrictedCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BridgeRestrictedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20BridgeRestrictedTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BridgeRestrictedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20BridgeRestrictedFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BridgeRestrictedSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20BridgeRestrictedSession struct {
	Contract     *ERC20BridgeRestricted // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ERC20BridgeRestrictedCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20BridgeRestrictedCallerSession struct {
	Contract *ERC20BridgeRestrictedCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ERC20BridgeRestrictedTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20BridgeRestrictedTransactorSession struct {
	Contract     *ERC20BridgeRestrictedTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ERC20BridgeRestrictedRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20BridgeRestrictedRaw struct {
	Contract *ERC20BridgeRestricted // Generic contract binding to access the raw methods on
}

// ERC20BridgeRestrictedCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20BridgeRestrictedCallerRaw struct {
	Contract *ERC20BridgeRestrictedCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20BridgeRestrictedTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20BridgeRestrictedTransactorRaw struct {
	Contract *ERC20BridgeRestrictedTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20BridgeRestricted creates a new instance of ERC20BridgeRestricted, bound to a specific deployed contract.
func NewERC20BridgeRestricted(address common.Address, backend bind.ContractBackend) (*ERC20BridgeRestricted, error) {
	contract, err := bindERC20BridgeRestricted(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestricted{ERC20BridgeRestrictedCaller: ERC20BridgeRestrictedCaller{contract: contract}, ERC20BridgeRestrictedTransactor: ERC20BridgeRestrictedTransactor{contract: contract}, ERC20BridgeRestrictedFilterer: ERC20BridgeRestrictedFilterer{contract: contract}}, nil
}

// NewERC20BridgeRestrictedCaller creates a new read-only instance of ERC20BridgeRestricted, bound to a specific deployed contract.
func NewERC20BridgeRestrictedCaller(address common.Address, caller bind.ContractCaller) (*ERC20BridgeRestrictedCaller, error) {
	contract, err := bindERC20BridgeRestricted(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedCaller{contract: contract}, nil
}

// NewERC20BridgeRestrictedTransactor creates a new write-only instance of ERC20BridgeRestricted, bound to a specific deployed contract.
func NewERC20BridgeRestrictedTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20BridgeRestrictedTransactor, error) {
	contract, err := bindERC20BridgeRestricted(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedTransactor{contract: contract}, nil
}

// NewERC20BridgeRestrictedFilterer creates a new log filterer instance of ERC20BridgeRestricted, bound to a specific deployed contract.
func NewERC20BridgeRestrictedFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20BridgeRestrictedFilterer, error) {
	contract, err := bindERC20BridgeRestricted(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedFilterer{contract: contract}, nil
}

// bindERC20BridgeRestricted binds a generic wrapper to an already deployed contract.
func bindERC20BridgeRestricted(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BridgeRestrictedABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20BridgeRestricted.Contract.ERC20BridgeRestrictedCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.ERC20BridgeRestrictedTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.ERC20BridgeRestrictedTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20BridgeRestricted.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.contract.Transact(opts, method, params...)
}

// DefaultWithdrawDelay is a free data retrieval call binding the contract method 0x3f4f199d.
//
// Solidity: function default_withdraw_delay() view returns(uint256)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) DefaultWithdrawDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "default_withdraw_delay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultWithdrawDelay is a free data retrieval call binding the contract method 0x3f4f199d.
//
// Solidity: function default_withdraw_delay() view returns(uint256)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) DefaultWithdrawDelay() (*big.Int, error) {
	return _ERC20BridgeRestricted.Contract.DefaultWithdrawDelay(&_ERC20BridgeRestricted.CallOpts)
}

// DefaultWithdrawDelay is a free data retrieval call binding the contract method 0x3f4f199d.
//
// Solidity: function default_withdraw_delay() view returns(uint256)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) DefaultWithdrawDelay() (*big.Int, error) {
	return _ERC20BridgeRestricted.Contract.DefaultWithdrawDelay(&_ERC20BridgeRestricted.CallOpts)
}

// Erc20AssetPoolAddress is a free data retrieval call binding the contract method 0x9356aab8.
//
// Solidity: function erc20_asset_pool_address() view returns(address)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) Erc20AssetPoolAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "erc20_asset_pool_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Erc20AssetPoolAddress is a free data retrieval call binding the contract method 0x9356aab8.
//
// Solidity: function erc20_asset_pool_address() view returns(address)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) Erc20AssetPoolAddress() (common.Address, error) {
	return _ERC20BridgeRestricted.Contract.Erc20AssetPoolAddress(&_ERC20BridgeRestricted.CallOpts)
}

// Erc20AssetPoolAddress is a free data retrieval call binding the contract method 0x9356aab8.
//
// Solidity: function erc20_asset_pool_address() view returns(address)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) Erc20AssetPoolAddress() (common.Address, error) {
	return _ERC20BridgeRestricted.Contract.Erc20AssetPoolAddress(&_ERC20BridgeRestricted.CallOpts)
}

// GetAssetDepositLifetimeLimit is a free data retrieval call binding the contract method 0x354a897a.
//
// Solidity: function get_asset_deposit_lifetime_limit(address asset_source) view returns(uint256)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) GetAssetDepositLifetimeLimit(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "get_asset_deposit_lifetime_limit", asset_source)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAssetDepositLifetimeLimit is a free data retrieval call binding the contract method 0x354a897a.
//
// Solidity: function get_asset_deposit_lifetime_limit(address asset_source) view returns(uint256)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) GetAssetDepositLifetimeLimit(asset_source common.Address) (*big.Int, error) {
	return _ERC20BridgeRestricted.Contract.GetAssetDepositLifetimeLimit(&_ERC20BridgeRestricted.CallOpts, asset_source)
}

// GetAssetDepositLifetimeLimit is a free data retrieval call binding the contract method 0x354a897a.
//
// Solidity: function get_asset_deposit_lifetime_limit(address asset_source) view returns(uint256)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) GetAssetDepositLifetimeLimit(asset_source common.Address) (*big.Int, error) {
	return _ERC20BridgeRestricted.Contract.GetAssetDepositLifetimeLimit(&_ERC20BridgeRestricted.CallOpts, asset_source)
}

// GetAssetSource is a free data retrieval call binding the contract method 0x786b0bc0.
//
// Solidity: function get_asset_source(bytes32 vega_asset_id) view returns(address)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) GetAssetSource(opts *bind.CallOpts, vega_asset_id [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "get_asset_source", vega_asset_id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAssetSource is a free data retrieval call binding the contract method 0x786b0bc0.
//
// Solidity: function get_asset_source(bytes32 vega_asset_id) view returns(address)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) GetAssetSource(vega_asset_id [32]byte) (common.Address, error) {
	return _ERC20BridgeRestricted.Contract.GetAssetSource(&_ERC20BridgeRestricted.CallOpts, vega_asset_id)
}

// GetAssetSource is a free data retrieval call binding the contract method 0x786b0bc0.
//
// Solidity: function get_asset_source(bytes32 vega_asset_id) view returns(address)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) GetAssetSource(vega_asset_id [32]byte) (common.Address, error) {
	return _ERC20BridgeRestricted.Contract.GetAssetSource(&_ERC20BridgeRestricted.CallOpts, vega_asset_id)
}

// GetMultisigControlAddress is a free data retrieval call binding the contract method 0xc58dc3b9.
//
// Solidity: function get_multisig_control_address() view returns(address)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) GetMultisigControlAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "get_multisig_control_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMultisigControlAddress is a free data retrieval call binding the contract method 0xc58dc3b9.
//
// Solidity: function get_multisig_control_address() view returns(address)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) GetMultisigControlAddress() (common.Address, error) {
	return _ERC20BridgeRestricted.Contract.GetMultisigControlAddress(&_ERC20BridgeRestricted.CallOpts)
}

// GetMultisigControlAddress is a free data retrieval call binding the contract method 0xc58dc3b9.
//
// Solidity: function get_multisig_control_address() view returns(address)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) GetMultisigControlAddress() (common.Address, error) {
	return _ERC20BridgeRestricted.Contract.GetMultisigControlAddress(&_ERC20BridgeRestricted.CallOpts)
}

// GetVegaAssetId is a free data retrieval call binding the contract method 0xa06b5d39.
//
// Solidity: function get_vega_asset_id(address asset_source) view returns(bytes32)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) GetVegaAssetId(opts *bind.CallOpts, asset_source common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "get_vega_asset_id", asset_source)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetVegaAssetId is a free data retrieval call binding the contract method 0xa06b5d39.
//
// Solidity: function get_vega_asset_id(address asset_source) view returns(bytes32)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) GetVegaAssetId(asset_source common.Address) ([32]byte, error) {
	return _ERC20BridgeRestricted.Contract.GetVegaAssetId(&_ERC20BridgeRestricted.CallOpts, asset_source)
}

// GetVegaAssetId is a free data retrieval call binding the contract method 0xa06b5d39.
//
// Solidity: function get_vega_asset_id(address asset_source) view returns(bytes32)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) GetVegaAssetId(asset_source common.Address) ([32]byte, error) {
	return _ERC20BridgeRestricted.Contract.GetVegaAssetId(&_ERC20BridgeRestricted.CallOpts, asset_source)
}

// GetWithdrawThreshold is a free data retrieval call binding the contract method 0xe8a7bce0.
//
// Solidity: function get_withdraw_threshold(address asset_source) view returns(uint256)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) GetWithdrawThreshold(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "get_withdraw_threshold", asset_source)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWithdrawThreshold is a free data retrieval call binding the contract method 0xe8a7bce0.
//
// Solidity: function get_withdraw_threshold(address asset_source) view returns(uint256)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) GetWithdrawThreshold(asset_source common.Address) (*big.Int, error) {
	return _ERC20BridgeRestricted.Contract.GetWithdrawThreshold(&_ERC20BridgeRestricted.CallOpts, asset_source)
}

// GetWithdrawThreshold is a free data retrieval call binding the contract method 0xe8a7bce0.
//
// Solidity: function get_withdraw_threshold(address asset_source) view returns(uint256)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) GetWithdrawThreshold(asset_source common.Address) (*big.Int, error) {
	return _ERC20BridgeRestricted.Contract.GetWithdrawThreshold(&_ERC20BridgeRestricted.CallOpts, asset_source)
}

// IsAssetListed is a free data retrieval call binding the contract method 0x7fd27b7f.
//
// Solidity: function is_asset_listed(address asset_source) view returns(bool)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) IsAssetListed(opts *bind.CallOpts, asset_source common.Address) (bool, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "is_asset_listed", asset_source)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAssetListed is a free data retrieval call binding the contract method 0x7fd27b7f.
//
// Solidity: function is_asset_listed(address asset_source) view returns(bool)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) IsAssetListed(asset_source common.Address) (bool, error) {
	return _ERC20BridgeRestricted.Contract.IsAssetListed(&_ERC20BridgeRestricted.CallOpts, asset_source)
}

// IsAssetListed is a free data retrieval call binding the contract method 0x7fd27b7f.
//
// Solidity: function is_asset_listed(address asset_source) view returns(bool)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) IsAssetListed(asset_source common.Address) (bool, error) {
	return _ERC20BridgeRestricted.Contract.IsAssetListed(&_ERC20BridgeRestricted.CallOpts, asset_source)
}

// IsExemptDepositor is a free data retrieval call binding the contract method 0x15c0df9d.
//
// Solidity: function is_exempt_depositor(address depositor) view returns(bool)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) IsExemptDepositor(opts *bind.CallOpts, depositor common.Address) (bool, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "is_exempt_depositor", depositor)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsExemptDepositor is a free data retrieval call binding the contract method 0x15c0df9d.
//
// Solidity: function is_exempt_depositor(address depositor) view returns(bool)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) IsExemptDepositor(depositor common.Address) (bool, error) {
	return _ERC20BridgeRestricted.Contract.IsExemptDepositor(&_ERC20BridgeRestricted.CallOpts, depositor)
}

// IsExemptDepositor is a free data retrieval call binding the contract method 0x15c0df9d.
//
// Solidity: function is_exempt_depositor(address depositor) view returns(bool)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) IsExemptDepositor(depositor common.Address) (bool, error) {
	return _ERC20BridgeRestricted.Contract.IsExemptDepositor(&_ERC20BridgeRestricted.CallOpts, depositor)
}

// IsStopped is a free data retrieval call binding the contract method 0xe272e9d0.
//
// Solidity: function is_stopped() view returns(bool)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCaller) IsStopped(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ERC20BridgeRestricted.contract.Call(opts, &out, "is_stopped")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsStopped is a free data retrieval call binding the contract method 0xe272e9d0.
//
// Solidity: function is_stopped() view returns(bool)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) IsStopped() (bool, error) {
	return _ERC20BridgeRestricted.Contract.IsStopped(&_ERC20BridgeRestricted.CallOpts)
}

// IsStopped is a free data retrieval call binding the contract method 0xe272e9d0.
//
// Solidity: function is_stopped() view returns(bool)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedCallerSession) IsStopped() (bool, error) {
	return _ERC20BridgeRestricted.Contract.IsStopped(&_ERC20BridgeRestricted.CallOpts)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xf7683932.
//
// Solidity: function deposit_asset(address asset_source, uint256 amount, bytes32 vega_public_key) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) DepositAsset(opts *bind.TransactOpts, asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "deposit_asset", asset_source, amount, vega_public_key)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xf7683932.
//
// Solidity: function deposit_asset(address asset_source, uint256 amount, bytes32 vega_public_key) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) DepositAsset(asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.DepositAsset(&_ERC20BridgeRestricted.TransactOpts, asset_source, amount, vega_public_key)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xf7683932.
//
// Solidity: function deposit_asset(address asset_source, uint256 amount, bytes32 vega_public_key) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) DepositAsset(asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.DepositAsset(&_ERC20BridgeRestricted.TransactOpts, asset_source, amount, vega_public_key)
}

// ExemptDepositor is a paid mutator transaction binding the contract method 0xb76fbb75.
//
// Solidity: function exempt_depositor() returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) ExemptDepositor(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "exempt_depositor")
}

// ExemptDepositor is a paid mutator transaction binding the contract method 0xb76fbb75.
//
// Solidity: function exempt_depositor() returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) ExemptDepositor() (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.ExemptDepositor(&_ERC20BridgeRestricted.TransactOpts)
}

// ExemptDepositor is a paid mutator transaction binding the contract method 0xb76fbb75.
//
// Solidity: function exempt_depositor() returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) ExemptDepositor() (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.ExemptDepositor(&_ERC20BridgeRestricted.TransactOpts)
}

// GlobalResume is a paid mutator transaction binding the contract method 0xd72ed529.
//
// Solidity: function global_resume(uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) GlobalResume(opts *bind.TransactOpts, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "global_resume", nonce, signatures)
}

// GlobalResume is a paid mutator transaction binding the contract method 0xd72ed529.
//
// Solidity: function global_resume(uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) GlobalResume(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.GlobalResume(&_ERC20BridgeRestricted.TransactOpts, nonce, signatures)
}

// GlobalResume is a paid mutator transaction binding the contract method 0xd72ed529.
//
// Solidity: function global_resume(uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) GlobalResume(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.GlobalResume(&_ERC20BridgeRestricted.TransactOpts, nonce, signatures)
}

// GlobalStop is a paid mutator transaction binding the contract method 0x9dfd3c88.
//
// Solidity: function global_stop(uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) GlobalStop(opts *bind.TransactOpts, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "global_stop", nonce, signatures)
}

// GlobalStop is a paid mutator transaction binding the contract method 0x9dfd3c88.
//
// Solidity: function global_stop(uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) GlobalStop(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.GlobalStop(&_ERC20BridgeRestricted.TransactOpts, nonce, signatures)
}

// GlobalStop is a paid mutator transaction binding the contract method 0x9dfd3c88.
//
// Solidity: function global_stop(uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) GlobalStop(nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.GlobalStop(&_ERC20BridgeRestricted.TransactOpts, nonce, signatures)
}

// ListAsset is a paid mutator transaction binding the contract method 0x0ff3562c.
//
// Solidity: function list_asset(address asset_source, bytes32 vega_asset_id, uint256 lifetime_limit, uint256 withdraw_threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) ListAsset(opts *bind.TransactOpts, asset_source common.Address, vega_asset_id [32]byte, lifetime_limit *big.Int, withdraw_threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "list_asset", asset_source, vega_asset_id, lifetime_limit, withdraw_threshold, nonce, signatures)
}

// ListAsset is a paid mutator transaction binding the contract method 0x0ff3562c.
//
// Solidity: function list_asset(address asset_source, bytes32 vega_asset_id, uint256 lifetime_limit, uint256 withdraw_threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) ListAsset(asset_source common.Address, vega_asset_id [32]byte, lifetime_limit *big.Int, withdraw_threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.ListAsset(&_ERC20BridgeRestricted.TransactOpts, asset_source, vega_asset_id, lifetime_limit, withdraw_threshold, nonce, signatures)
}

// ListAsset is a paid mutator transaction binding the contract method 0x0ff3562c.
//
// Solidity: function list_asset(address asset_source, bytes32 vega_asset_id, uint256 lifetime_limit, uint256 withdraw_threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) ListAsset(asset_source common.Address, vega_asset_id [32]byte, lifetime_limit *big.Int, withdraw_threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.ListAsset(&_ERC20BridgeRestricted.TransactOpts, asset_source, vega_asset_id, lifetime_limit, withdraw_threshold, nonce, signatures)
}

// RemoveAsset is a paid mutator transaction binding the contract method 0xc76de358.
//
// Solidity: function remove_asset(address asset_source, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) RemoveAsset(opts *bind.TransactOpts, asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "remove_asset", asset_source, nonce, signatures)
}

// RemoveAsset is a paid mutator transaction binding the contract method 0xc76de358.
//
// Solidity: function remove_asset(address asset_source, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) RemoveAsset(asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.RemoveAsset(&_ERC20BridgeRestricted.TransactOpts, asset_source, nonce, signatures)
}

// RemoveAsset is a paid mutator transaction binding the contract method 0xc76de358.
//
// Solidity: function remove_asset(address asset_source, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) RemoveAsset(asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.RemoveAsset(&_ERC20BridgeRestricted.TransactOpts, asset_source, nonce, signatures)
}

// RevokeExemptDepositor is a paid mutator transaction binding the contract method 0x6a1c6fa4.
//
// Solidity: function revoke_exempt_depositor() returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) RevokeExemptDepositor(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "revoke_exempt_depositor")
}

// RevokeExemptDepositor is a paid mutator transaction binding the contract method 0x6a1c6fa4.
//
// Solidity: function revoke_exempt_depositor() returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) RevokeExemptDepositor() (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.RevokeExemptDepositor(&_ERC20BridgeRestricted.TransactOpts)
}

// RevokeExemptDepositor is a paid mutator transaction binding the contract method 0x6a1c6fa4.
//
// Solidity: function revoke_exempt_depositor() returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) RevokeExemptDepositor() (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.RevokeExemptDepositor(&_ERC20BridgeRestricted.TransactOpts)
}

// SetAssetLimits is a paid mutator transaction binding the contract method 0x41fb776d.
//
// Solidity: function set_asset_limits(address asset_source, uint256 lifetime_limit, uint256 threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) SetAssetLimits(opts *bind.TransactOpts, asset_source common.Address, lifetime_limit *big.Int, threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "set_asset_limits", asset_source, lifetime_limit, threshold, nonce, signatures)
}

// SetAssetLimits is a paid mutator transaction binding the contract method 0x41fb776d.
//
// Solidity: function set_asset_limits(address asset_source, uint256 lifetime_limit, uint256 threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) SetAssetLimits(asset_source common.Address, lifetime_limit *big.Int, threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.SetAssetLimits(&_ERC20BridgeRestricted.TransactOpts, asset_source, lifetime_limit, threshold, nonce, signatures)
}

// SetAssetLimits is a paid mutator transaction binding the contract method 0x41fb776d.
//
// Solidity: function set_asset_limits(address asset_source, uint256 lifetime_limit, uint256 threshold, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) SetAssetLimits(asset_source common.Address, lifetime_limit *big.Int, threshold *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.SetAssetLimits(&_ERC20BridgeRestricted.TransactOpts, asset_source, lifetime_limit, threshold, nonce, signatures)
}

// SetWithdrawDelay is a paid mutator transaction binding the contract method 0x5a246728.
//
// Solidity: function set_withdraw_delay(uint256 delay, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) SetWithdrawDelay(opts *bind.TransactOpts, delay *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "set_withdraw_delay", delay, nonce, signatures)
}

// SetWithdrawDelay is a paid mutator transaction binding the contract method 0x5a246728.
//
// Solidity: function set_withdraw_delay(uint256 delay, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) SetWithdrawDelay(delay *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.SetWithdrawDelay(&_ERC20BridgeRestricted.TransactOpts, delay, nonce, signatures)
}

// SetWithdrawDelay is a paid mutator transaction binding the contract method 0x5a246728.
//
// Solidity: function set_withdraw_delay(uint256 delay, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) SetWithdrawDelay(delay *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.SetWithdrawDelay(&_ERC20BridgeRestricted.TransactOpts, delay, nonce, signatures)
}

// WithdrawAsset is a paid mutator transaction binding the contract method 0x3ad90635.
//
// Solidity: function withdraw_asset(address asset_source, uint256 amount, address target, uint256 creation, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactor) WithdrawAsset(opts *bind.TransactOpts, asset_source common.Address, amount *big.Int, target common.Address, creation *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.contract.Transact(opts, "withdraw_asset", asset_source, amount, target, creation, nonce, signatures)
}

// WithdrawAsset is a paid mutator transaction binding the contract method 0x3ad90635.
//
// Solidity: function withdraw_asset(address asset_source, uint256 amount, address target, uint256 creation, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedSession) WithdrawAsset(asset_source common.Address, amount *big.Int, target common.Address, creation *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.WithdrawAsset(&_ERC20BridgeRestricted.TransactOpts, asset_source, amount, target, creation, nonce, signatures)
}

// WithdrawAsset is a paid mutator transaction binding the contract method 0x3ad90635.
//
// Solidity: function withdraw_asset(address asset_source, uint256 amount, address target, uint256 creation, uint256 nonce, bytes signatures) returns()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedTransactorSession) WithdrawAsset(asset_source common.Address, amount *big.Int, target common.Address, creation *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20BridgeRestricted.Contract.WithdrawAsset(&_ERC20BridgeRestricted.TransactOpts, asset_source, amount, target, creation, nonce, signatures)
}

// ERC20BridgeRestrictedAssetDepositedIterator is returned from FilterAssetDeposited and is used to iterate over the raw logs and unpacked data for AssetDeposited events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetDepositedIterator struct {
	Event *ERC20BridgeRestrictedAssetDeposited // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedAssetDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedAssetDeposited)
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
		it.Event = new(ERC20BridgeRestrictedAssetDeposited)
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
func (it *ERC20BridgeRestrictedAssetDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedAssetDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedAssetDeposited represents a AssetDeposited event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetDeposited struct {
	UserAddress   common.Address
	AssetSource   common.Address
	Amount        *big.Int
	VegaPublicKey [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAssetDeposited is a free log retrieval operation binding the contract event 0x3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf.
//
// Solidity: event Asset_Deposited(address indexed user_address, address indexed asset_source, uint256 amount, bytes32 vega_public_key)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterAssetDeposited(opts *bind.FilterOpts, user_address []common.Address, asset_source []common.Address) (*ERC20BridgeRestrictedAssetDepositedIterator, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Asset_Deposited", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedAssetDepositedIterator{contract: _ERC20BridgeRestricted.contract, event: "Asset_Deposited", logs: logs, sub: sub}, nil
}

// WatchAssetDeposited is a free log subscription operation binding the contract event 0x3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf.
//
// Solidity: event Asset_Deposited(address indexed user_address, address indexed asset_source, uint256 amount, bytes32 vega_public_key)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchAssetDeposited(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedAssetDeposited, user_address []common.Address, asset_source []common.Address) (event.Subscription, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Asset_Deposited", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedAssetDeposited)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Deposited", log); err != nil {
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

// ParseAssetDeposited is a log parse operation binding the contract event 0x3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf.
//
// Solidity: event Asset_Deposited(address indexed user_address, address indexed asset_source, uint256 amount, bytes32 vega_public_key)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseAssetDeposited(log types.Log) (*ERC20BridgeRestrictedAssetDeposited, error) {
	event := new(ERC20BridgeRestrictedAssetDeposited)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeRestrictedAssetLimitsUpdatedIterator is returned from FilterAssetLimitsUpdated and is used to iterate over the raw logs and unpacked data for AssetLimitsUpdated events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetLimitsUpdatedIterator struct {
	Event *ERC20BridgeRestrictedAssetLimitsUpdated // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedAssetLimitsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedAssetLimitsUpdated)
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
		it.Event = new(ERC20BridgeRestrictedAssetLimitsUpdated)
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
func (it *ERC20BridgeRestrictedAssetLimitsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedAssetLimitsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedAssetLimitsUpdated represents a AssetLimitsUpdated event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetLimitsUpdated struct {
	AssetSource       common.Address
	LifetimeLimit     *big.Int
	WithdrawThreshold *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterAssetLimitsUpdated is a free log retrieval operation binding the contract event 0xfc7eab762b8751ad85c101fd1025c763b4e8d48f2093f506629b606618e884fe.
//
// Solidity: event Asset_Limits_Updated(address indexed asset_source, uint256 lifetime_limit, uint256 withdraw_threshold)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterAssetLimitsUpdated(opts *bind.FilterOpts, asset_source []common.Address) (*ERC20BridgeRestrictedAssetLimitsUpdatedIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Asset_Limits_Updated", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedAssetLimitsUpdatedIterator{contract: _ERC20BridgeRestricted.contract, event: "Asset_Limits_Updated", logs: logs, sub: sub}, nil
}

// WatchAssetLimitsUpdated is a free log subscription operation binding the contract event 0xfc7eab762b8751ad85c101fd1025c763b4e8d48f2093f506629b606618e884fe.
//
// Solidity: event Asset_Limits_Updated(address indexed asset_source, uint256 lifetime_limit, uint256 withdraw_threshold)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchAssetLimitsUpdated(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedAssetLimitsUpdated, asset_source []common.Address) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Asset_Limits_Updated", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedAssetLimitsUpdated)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Limits_Updated", log); err != nil {
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

// ParseAssetLimitsUpdated is a log parse operation binding the contract event 0xfc7eab762b8751ad85c101fd1025c763b4e8d48f2093f506629b606618e884fe.
//
// Solidity: event Asset_Limits_Updated(address indexed asset_source, uint256 lifetime_limit, uint256 withdraw_threshold)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseAssetLimitsUpdated(log types.Log) (*ERC20BridgeRestrictedAssetLimitsUpdated, error) {
	event := new(ERC20BridgeRestrictedAssetLimitsUpdated)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Limits_Updated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeRestrictedAssetListedIterator is returned from FilterAssetListed and is used to iterate over the raw logs and unpacked data for AssetListed events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetListedIterator struct {
	Event *ERC20BridgeRestrictedAssetListed // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedAssetListedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedAssetListed)
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
		it.Event = new(ERC20BridgeRestrictedAssetListed)
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
func (it *ERC20BridgeRestrictedAssetListedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedAssetListedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedAssetListed represents a AssetListed event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetListed struct {
	AssetSource common.Address
	VegaAssetId [32]byte
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetListed is a free log retrieval operation binding the contract event 0x4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e.
//
// Solidity: event Asset_Listed(address indexed asset_source, bytes32 indexed vega_asset_id, uint256 nonce)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterAssetListed(opts *bind.FilterOpts, asset_source []common.Address, vega_asset_id [][32]byte) (*ERC20BridgeRestrictedAssetListedIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}
	var vega_asset_idRule []interface{}
	for _, vega_asset_idItem := range vega_asset_id {
		vega_asset_idRule = append(vega_asset_idRule, vega_asset_idItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Asset_Listed", asset_sourceRule, vega_asset_idRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedAssetListedIterator{contract: _ERC20BridgeRestricted.contract, event: "Asset_Listed", logs: logs, sub: sub}, nil
}

// WatchAssetListed is a free log subscription operation binding the contract event 0x4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e.
//
// Solidity: event Asset_Listed(address indexed asset_source, bytes32 indexed vega_asset_id, uint256 nonce)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchAssetListed(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedAssetListed, asset_source []common.Address, vega_asset_id [][32]byte) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}
	var vega_asset_idRule []interface{}
	for _, vega_asset_idItem := range vega_asset_id {
		vega_asset_idRule = append(vega_asset_idRule, vega_asset_idItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Asset_Listed", asset_sourceRule, vega_asset_idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedAssetListed)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Listed", log); err != nil {
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

// ParseAssetListed is a log parse operation binding the contract event 0x4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e.
//
// Solidity: event Asset_Listed(address indexed asset_source, bytes32 indexed vega_asset_id, uint256 nonce)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseAssetListed(log types.Log) (*ERC20BridgeRestrictedAssetListed, error) {
	event := new(ERC20BridgeRestrictedAssetListed)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Listed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeRestrictedAssetRemovedIterator is returned from FilterAssetRemoved and is used to iterate over the raw logs and unpacked data for AssetRemoved events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetRemovedIterator struct {
	Event *ERC20BridgeRestrictedAssetRemoved // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedAssetRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedAssetRemoved)
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
		it.Event = new(ERC20BridgeRestrictedAssetRemoved)
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
func (it *ERC20BridgeRestrictedAssetRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedAssetRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedAssetRemoved represents a AssetRemoved event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetRemoved struct {
	AssetSource common.Address
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetRemoved is a free log retrieval operation binding the contract event 0x58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c9.
//
// Solidity: event Asset_Removed(address indexed asset_source, uint256 nonce)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterAssetRemoved(opts *bind.FilterOpts, asset_source []common.Address) (*ERC20BridgeRestrictedAssetRemovedIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Asset_Removed", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedAssetRemovedIterator{contract: _ERC20BridgeRestricted.contract, event: "Asset_Removed", logs: logs, sub: sub}, nil
}

// WatchAssetRemoved is a free log subscription operation binding the contract event 0x58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c9.
//
// Solidity: event Asset_Removed(address indexed asset_source, uint256 nonce)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchAssetRemoved(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedAssetRemoved, asset_source []common.Address) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Asset_Removed", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedAssetRemoved)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Removed", log); err != nil {
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

// ParseAssetRemoved is a log parse operation binding the contract event 0x58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c9.
//
// Solidity: event Asset_Removed(address indexed asset_source, uint256 nonce)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseAssetRemoved(log types.Log) (*ERC20BridgeRestrictedAssetRemoved, error) {
	event := new(ERC20BridgeRestrictedAssetRemoved)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Removed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeRestrictedAssetWithdrawnIterator is returned from FilterAssetWithdrawn and is used to iterate over the raw logs and unpacked data for AssetWithdrawn events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetWithdrawnIterator struct {
	Event *ERC20BridgeRestrictedAssetWithdrawn // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedAssetWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedAssetWithdrawn)
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
		it.Event = new(ERC20BridgeRestrictedAssetWithdrawn)
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
func (it *ERC20BridgeRestrictedAssetWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedAssetWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedAssetWithdrawn represents a AssetWithdrawn event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedAssetWithdrawn struct {
	UserAddress common.Address
	AssetSource common.Address
	Amount      *big.Int
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetWithdrawn is a free log retrieval operation binding the contract event 0xa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df.
//
// Solidity: event Asset_Withdrawn(address indexed user_address, address indexed asset_source, uint256 amount, uint256 nonce)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterAssetWithdrawn(opts *bind.FilterOpts, user_address []common.Address, asset_source []common.Address) (*ERC20BridgeRestrictedAssetWithdrawnIterator, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Asset_Withdrawn", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedAssetWithdrawnIterator{contract: _ERC20BridgeRestricted.contract, event: "Asset_Withdrawn", logs: logs, sub: sub}, nil
}

// WatchAssetWithdrawn is a free log subscription operation binding the contract event 0xa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df.
//
// Solidity: event Asset_Withdrawn(address indexed user_address, address indexed asset_source, uint256 amount, uint256 nonce)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchAssetWithdrawn(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedAssetWithdrawn, user_address []common.Address, asset_source []common.Address) (event.Subscription, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Asset_Withdrawn", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedAssetWithdrawn)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Withdrawn", log); err != nil {
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

// ParseAssetWithdrawn is a log parse operation binding the contract event 0xa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df.
//
// Solidity: event Asset_Withdrawn(address indexed user_address, address indexed asset_source, uint256 amount, uint256 nonce)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseAssetWithdrawn(log types.Log) (*ERC20BridgeRestrictedAssetWithdrawn, error) {
	event := new(ERC20BridgeRestrictedAssetWithdrawn)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Asset_Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeRestrictedBridgeResumedIterator is returned from FilterBridgeResumed and is used to iterate over the raw logs and unpacked data for BridgeResumed events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedBridgeResumedIterator struct {
	Event *ERC20BridgeRestrictedBridgeResumed // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedBridgeResumedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedBridgeResumed)
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
		it.Event = new(ERC20BridgeRestrictedBridgeResumed)
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
func (it *ERC20BridgeRestrictedBridgeResumedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedBridgeResumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedBridgeResumed represents a BridgeResumed event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedBridgeResumed struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterBridgeResumed is a free log retrieval operation binding the contract event 0x79c02b0e60e0f00fe0370791204f2f175fe3f06f4816f3506ad4fa1b8e8cde0f.
//
// Solidity: event Bridge_Resumed()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterBridgeResumed(opts *bind.FilterOpts) (*ERC20BridgeRestrictedBridgeResumedIterator, error) {

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Bridge_Resumed")
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedBridgeResumedIterator{contract: _ERC20BridgeRestricted.contract, event: "Bridge_Resumed", logs: logs, sub: sub}, nil
}

// WatchBridgeResumed is a free log subscription operation binding the contract event 0x79c02b0e60e0f00fe0370791204f2f175fe3f06f4816f3506ad4fa1b8e8cde0f.
//
// Solidity: event Bridge_Resumed()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchBridgeResumed(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedBridgeResumed) (event.Subscription, error) {

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Bridge_Resumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedBridgeResumed)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Bridge_Resumed", log); err != nil {
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

// ParseBridgeResumed is a log parse operation binding the contract event 0x79c02b0e60e0f00fe0370791204f2f175fe3f06f4816f3506ad4fa1b8e8cde0f.
//
// Solidity: event Bridge_Resumed()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseBridgeResumed(log types.Log) (*ERC20BridgeRestrictedBridgeResumed, error) {
	event := new(ERC20BridgeRestrictedBridgeResumed)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Bridge_Resumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeRestrictedBridgeStoppedIterator is returned from FilterBridgeStopped and is used to iterate over the raw logs and unpacked data for BridgeStopped events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedBridgeStoppedIterator struct {
	Event *ERC20BridgeRestrictedBridgeStopped // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedBridgeStoppedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedBridgeStopped)
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
		it.Event = new(ERC20BridgeRestrictedBridgeStopped)
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
func (it *ERC20BridgeRestrictedBridgeStoppedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedBridgeStoppedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedBridgeStopped represents a BridgeStopped event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedBridgeStopped struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterBridgeStopped is a free log retrieval operation binding the contract event 0x129d99581c8e70519df1f0733d3212f33d0ed3ea6144adacc336c647f1d36382.
//
// Solidity: event Bridge_Stopped()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterBridgeStopped(opts *bind.FilterOpts) (*ERC20BridgeRestrictedBridgeStoppedIterator, error) {

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Bridge_Stopped")
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedBridgeStoppedIterator{contract: _ERC20BridgeRestricted.contract, event: "Bridge_Stopped", logs: logs, sub: sub}, nil
}

// WatchBridgeStopped is a free log subscription operation binding the contract event 0x129d99581c8e70519df1f0733d3212f33d0ed3ea6144adacc336c647f1d36382.
//
// Solidity: event Bridge_Stopped()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchBridgeStopped(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedBridgeStopped) (event.Subscription, error) {

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Bridge_Stopped")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedBridgeStopped)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Bridge_Stopped", log); err != nil {
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

// ParseBridgeStopped is a log parse operation binding the contract event 0x129d99581c8e70519df1f0733d3212f33d0ed3ea6144adacc336c647f1d36382.
//
// Solidity: event Bridge_Stopped()
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseBridgeStopped(log types.Log) (*ERC20BridgeRestrictedBridgeStopped, error) {
	event := new(ERC20BridgeRestrictedBridgeStopped)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Bridge_Stopped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeRestrictedBridgeWithdrawDelaySetIterator is returned from FilterBridgeWithdrawDelaySet and is used to iterate over the raw logs and unpacked data for BridgeWithdrawDelaySet events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedBridgeWithdrawDelaySetIterator struct {
	Event *ERC20BridgeRestrictedBridgeWithdrawDelaySet // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedBridgeWithdrawDelaySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedBridgeWithdrawDelaySet)
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
		it.Event = new(ERC20BridgeRestrictedBridgeWithdrawDelaySet)
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
func (it *ERC20BridgeRestrictedBridgeWithdrawDelaySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedBridgeWithdrawDelaySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedBridgeWithdrawDelaySet represents a BridgeWithdrawDelaySet event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedBridgeWithdrawDelaySet struct {
	WithdrawDelay *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterBridgeWithdrawDelaySet is a free log retrieval operation binding the contract event 0x1c7e8f73a01b8af4e18dd34455a42a45ad742bdb79cfda77bbdf50db2391fc88.
//
// Solidity: event Bridge_Withdraw_Delay_Set(uint256 withdraw_delay)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterBridgeWithdrawDelaySet(opts *bind.FilterOpts) (*ERC20BridgeRestrictedBridgeWithdrawDelaySetIterator, error) {

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Bridge_Withdraw_Delay_Set")
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedBridgeWithdrawDelaySetIterator{contract: _ERC20BridgeRestricted.contract, event: "Bridge_Withdraw_Delay_Set", logs: logs, sub: sub}, nil
}

// WatchBridgeWithdrawDelaySet is a free log subscription operation binding the contract event 0x1c7e8f73a01b8af4e18dd34455a42a45ad742bdb79cfda77bbdf50db2391fc88.
//
// Solidity: event Bridge_Withdraw_Delay_Set(uint256 withdraw_delay)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchBridgeWithdrawDelaySet(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedBridgeWithdrawDelaySet) (event.Subscription, error) {

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Bridge_Withdraw_Delay_Set")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedBridgeWithdrawDelaySet)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Bridge_Withdraw_Delay_Set", log); err != nil {
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

// ParseBridgeWithdrawDelaySet is a log parse operation binding the contract event 0x1c7e8f73a01b8af4e18dd34455a42a45ad742bdb79cfda77bbdf50db2391fc88.
//
// Solidity: event Bridge_Withdraw_Delay_Set(uint256 withdraw_delay)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseBridgeWithdrawDelaySet(log types.Log) (*ERC20BridgeRestrictedBridgeWithdrawDelaySet, error) {
	event := new(ERC20BridgeRestrictedBridgeWithdrawDelaySet)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Bridge_Withdraw_Delay_Set", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeRestrictedDepositorExemptedIterator is returned from FilterDepositorExempted and is used to iterate over the raw logs and unpacked data for DepositorExempted events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedDepositorExemptedIterator struct {
	Event *ERC20BridgeRestrictedDepositorExempted // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedDepositorExemptedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedDepositorExempted)
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
		it.Event = new(ERC20BridgeRestrictedDepositorExempted)
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
func (it *ERC20BridgeRestrictedDepositorExemptedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedDepositorExemptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedDepositorExempted represents a DepositorExempted event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedDepositorExempted struct {
	Depositor common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDepositorExempted is a free log retrieval operation binding the contract event 0xf56e0868b913034a60dbca9c89ee79f8b0fa18dadbc5f6665f2f9a2cf3f51cdb.
//
// Solidity: event Depositor_Exempted(address indexed depositor)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterDepositorExempted(opts *bind.FilterOpts, depositor []common.Address) (*ERC20BridgeRestrictedDepositorExemptedIterator, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Depositor_Exempted", depositorRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedDepositorExemptedIterator{contract: _ERC20BridgeRestricted.contract, event: "Depositor_Exempted", logs: logs, sub: sub}, nil
}

// WatchDepositorExempted is a free log subscription operation binding the contract event 0xf56e0868b913034a60dbca9c89ee79f8b0fa18dadbc5f6665f2f9a2cf3f51cdb.
//
// Solidity: event Depositor_Exempted(address indexed depositor)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchDepositorExempted(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedDepositorExempted, depositor []common.Address) (event.Subscription, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Depositor_Exempted", depositorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedDepositorExempted)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Depositor_Exempted", log); err != nil {
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

// ParseDepositorExempted is a log parse operation binding the contract event 0xf56e0868b913034a60dbca9c89ee79f8b0fa18dadbc5f6665f2f9a2cf3f51cdb.
//
// Solidity: event Depositor_Exempted(address indexed depositor)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseDepositorExempted(log types.Log) (*ERC20BridgeRestrictedDepositorExempted, error) {
	event := new(ERC20BridgeRestrictedDepositorExempted)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Depositor_Exempted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeRestrictedDepositorExemptionRevokedIterator is returned from FilterDepositorExemptionRevoked and is used to iterate over the raw logs and unpacked data for DepositorExemptionRevoked events raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedDepositorExemptionRevokedIterator struct {
	Event *ERC20BridgeRestrictedDepositorExemptionRevoked // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeRestrictedDepositorExemptionRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeRestrictedDepositorExemptionRevoked)
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
		it.Event = new(ERC20BridgeRestrictedDepositorExemptionRevoked)
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
func (it *ERC20BridgeRestrictedDepositorExemptionRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeRestrictedDepositorExemptionRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeRestrictedDepositorExemptionRevoked represents a DepositorExemptionRevoked event raised by the ERC20BridgeRestricted contract.
type ERC20BridgeRestrictedDepositorExemptionRevoked struct {
	Depositor common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDepositorExemptionRevoked is a free log retrieval operation binding the contract event 0xe74b113dca87276d976f476a9b4b9da3c780a3262eaabad051ee4e98912936a4.
//
// Solidity: event Depositor_Exemption_Revoked(address indexed depositor)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) FilterDepositorExemptionRevoked(opts *bind.FilterOpts, depositor []common.Address) (*ERC20BridgeRestrictedDepositorExemptionRevokedIterator, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.FilterLogs(opts, "Depositor_Exemption_Revoked", depositorRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeRestrictedDepositorExemptionRevokedIterator{contract: _ERC20BridgeRestricted.contract, event: "Depositor_Exemption_Revoked", logs: logs, sub: sub}, nil
}

// WatchDepositorExemptionRevoked is a free log subscription operation binding the contract event 0xe74b113dca87276d976f476a9b4b9da3c780a3262eaabad051ee4e98912936a4.
//
// Solidity: event Depositor_Exemption_Revoked(address indexed depositor)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) WatchDepositorExemptionRevoked(opts *bind.WatchOpts, sink chan<- *ERC20BridgeRestrictedDepositorExemptionRevoked, depositor []common.Address) (event.Subscription, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ERC20BridgeRestricted.contract.WatchLogs(opts, "Depositor_Exemption_Revoked", depositorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeRestrictedDepositorExemptionRevoked)
				if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Depositor_Exemption_Revoked", log); err != nil {
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

// ParseDepositorExemptionRevoked is a log parse operation binding the contract event 0xe74b113dca87276d976f476a9b4b9da3c780a3262eaabad051ee4e98912936a4.
//
// Solidity: event Depositor_Exemption_Revoked(address indexed depositor)
func (_ERC20BridgeRestricted *ERC20BridgeRestrictedFilterer) ParseDepositorExemptionRevoked(log types.Log) (*ERC20BridgeRestrictedDepositorExemptionRevoked, error) {
	event := new(ERC20BridgeRestrictedDepositorExemptionRevoked)
	if err := _ERC20BridgeRestricted.contract.UnpackLog(event, "Depositor_Exemption_Revoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
