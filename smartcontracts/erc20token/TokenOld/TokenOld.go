// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TokenOld

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

// TokenOldMetaData contains all meta data concerning the TokenOld contract.
var TokenOldMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"total_supply_whole_tokens\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"faucet_amount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"bridge_address\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"vega_public_keys\",\"type\":\"bytes32[]\"}],\"name\":\"admin_deposit_bulk\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"bridge_address\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"admin_deposit_single\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"staking_bridge_address\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"vega_public_keys\",\"type\":\"bytes32[]\"}],\"name\":\"admin_stake_bulk\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"faucet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"issue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162002ecc38038062002ecc83398181016040528101906200003791906200038d565b84848482600090805190602001906200005292919062000231565b5081600190805190602001906200006b92919062000231565b5080600260006101000a81548160ff021916908360ff16021790555050505033600260016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600260019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a360008360ff16600a6200015c91906200053b565b8362000169919062000678565b9050816006819055508060058190555080600360003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200021d919062000464565b60405180910390a350505050505062000885565b8280546200023f9062000726565b90600052602060002090601f016020900481019282620002635760008555620002af565b82601f106200027e57805160ff1916838001178555620002af565b82800160010185558215620002af579182015b82811115620002ae57825182559160200191906001019062000291565b5b509050620002be9190620002c2565b5090565b5b80821115620002dd576000816000905550600101620002c3565b5090565b6000620002f8620002f284620004aa565b62000481565b90508281526020810184848401111562000317576200031662000824565b5b62000324848285620006f0565b509392505050565b600082601f8301126200034457620003436200081f565b5b815162000356848260208601620002e1565b91505092915050565b600081519050620003708162000851565b92915050565b60008151905062000387816200086b565b92915050565b600080600080600060a08688031215620003ac57620003ab6200082e565b5b600086015167ffffffffffffffff811115620003cd57620003cc62000829565b5b620003db888289016200032c565b955050602086015167ffffffffffffffff811115620003ff57620003fe62000829565b5b6200040d888289016200032c565b9450506040620004208882890162000376565b935050606062000433888289016200035f565b925050608062000446888289016200035f565b9150509295509295909350565b6200045e81620006d9565b82525050565b60006020820190506200047b600083018462000453565b92915050565b60006200048d620004a0565b90506200049b82826200075c565b919050565b6000604051905090565b600067ffffffffffffffff821115620004c857620004c7620007f0565b5b620004d38262000833565b9050602081019050919050565b6000808291508390505b600185111562000532578086048111156200050a576200050962000792565b5b60018516156200051a5780820291505b80810290506200052a8562000844565b9450620004ea565b94509492505050565b60006200054882620006d9565b91506200055583620006d9565b9250620005847fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84846200058c565b905092915050565b6000826200059e576001905062000671565b81620005ae576000905062000671565b8160018114620005c75760028114620005d25762000608565b600191505062000671565b60ff841115620005e757620005e662000792565b5b8360020a91508482111562000601576200060062000792565b5b5062000671565b5060208310610133831016604e8410600b8410161715620006425782820a9050838111156200063c576200063b62000792565b5b62000671565b620006518484846001620004e0565b925090508184048111156200066b576200066a62000792565b5b81810290505b9392505050565b60006200068582620006d9565b91506200069283620006d9565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615620006ce57620006cd62000792565b5b828202905092915050565b6000819050919050565b600060ff82169050919050565b60005b8381101562000710578082015181840152602081019050620006f3565b8381111562000720576000848401525b50505050565b600060028204905060018216806200073f57607f821691505b60208210811415620007565762000755620007c1565b5b50919050565b620007678262000833565b810181811067ffffffffffffffff82111715620007895762000788620007f0565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b60008160011c9050919050565b6200085c81620006d9565b81146200086857600080fd5b50565b6200087681620006e3565b81146200088257600080fd5b50565b61263780620008956000396000f3fe608060405234801561001057600080fd5b50600436106101375760003560e01c80638da5cb5b116100b8578063b777374c1161007c578063b777374c14610340578063bc36878e1461035c578063d779cae814610378578063dd62ed3e14610394578063de5f72fd146103c4578063f2fde38b146103ce57610137565b80638da5cb5b146102865780638f32d59b146102a457806395d89b41146102c2578063a457c2d7146102e0578063a9059cbb1461031057610137565b806339509351116100ff57806339509351146101f657806341c0e1b51461022657806370a0823114610230578063715018a614610260578063867904b41461026a57610137565b806306fdde031461013c578063095ea7b31461015a57806318160ddd1461018a57806323b872dd146101a8578063313ce567146101d8575b600080fd5b6101446103ea565b6040516101519190611ed8565b60405180910390f35b610174600480360381019061016f9190611bf0565b61047c565b6040516101819190611ebd565b60405180910390f35b610192610493565b60405161019f9190611fda565b60405180910390f35b6101c260048036038101906101bd9190611b9d565b61049d565b6040516101cf9190611ebd565b60405180910390f35b6101e061054e565b6040516101ed919061201e565b60405180910390f35b610210600480360381019061020b9190611bf0565b610565565b60405161021d9190611ebd565b60405180910390f35b61022e61060a565b005b61024a60048036038101906102459190611b30565b610676565b6040516102579190611fda565b60405180910390f35b6102686106bf565b005b610284600480360381019061027f9190611bf0565b6107c7565b005b61028e61081d565b60405161029b9190611e6b565b60405180910390f35b6102ac610847565b6040516102b99190611ebd565b60405180910390f35b6102ca61089f565b6040516102d79190611ed8565b60405180910390f35b6102fa60048036038101906102f59190611bf0565b610931565b6040516103079190611ebd565b60405180910390f35b61032a60048036038101906103259190611bf0565b6109d6565b6040516103379190611ebd565b60405180910390f35b61035a60048036038101906103559190611c9f565b6109ed565b005b61037660048036038101906103719190611c30565b610c3f565b005b610392600480360381019061038d9190611c30565b610f03565b005b6103ae60048036038101906103a99190611b5d565b6111c5565b6040516103bb9190611fda565b60405180910390f35b6103cc61124c565b005b6103e860048036038101906103e39190611b30565b61143e565b005b6060600080546103f99061221c565b80601f01602080910402602001604051908101604052809291908181526020018280546104259061221c565b80156104725780601f1061044757610100808354040283529160200191610472565b820191906000526020600020905b81548152906001019060200180831161045557829003601f168201915b5050505050905090565b6000610489338484611491565b6001905092915050565b6000600554905090565b60006104aa84848461165c565b610543843361053e85600460008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118d090919063ffffffff16565b611491565b600190509392505050565b6000600260009054906101000a900460ff16905090565b600061060033846105fb85600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118f790919063ffffffff16565b611491565b6001905092915050565b610612610847565b610651576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064890611f5a565b60405180910390fd5b600061065b61081d565b90508073ffffffffffffffffffffffffffffffffffffffff16ff5b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6106c7610847565b610706576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106fd90611f5a565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff16600260019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a36000600260016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b6107cf610847565b61080e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161080590611f5a565b60405180910390fd5b61081930838361165c565b5050565b6000600260019054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600260019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614905090565b6060600180546108ae9061221c565b80601f01602080910402602001604051908101604052809291908181526020018280546108da9061221c565b80156109275780601f106108fc57610100808354040283529160200191610927565b820191906000526020600020905b81548152906001019060200180831161090a57829003601f168201915b5050505050905090565b60006109cc33846109c785600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118d090919063ffffffff16565b611491565b6001905092915050565b60006109e333848461165c565b6001905092915050565b6109f5610847565b610a34576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a2b90611f5a565b60405180910390fd5b82600460003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610aca836005546118f790919063ffffffff16565b600581905550610b2283600360003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118f790919063ffffffff16565b600360003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef85604051610bc39190611fda565b60405180910390a38173ffffffffffffffffffffffffffffffffffffffff1663f76839323085846040518463ffffffff1660e01b8152600401610c0893929190611e86565b600060405180830381600087803b158015610c2257600080fd5b505af1158015610c36573d6000803e3d6000fd5b50505050505050565b610c47610847565b610c86576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c7d90611f5a565b60405180910390fd5b6000815184610c9591906120fc565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600460003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610d4d816005546118f790919063ffffffff16565b600581905550610da581600360003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118f790919063ffffffff16565b600360003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610e469190611fda565b60405180910390a360005b82518160ff161015610efc578373ffffffffffffffffffffffffffffffffffffffff1663f76839323087868560ff1681518110610e9157610e90612336565b5b60200260200101516040518463ffffffff1660e01b8152600401610eb793929190611e86565b600060405180830381600087803b158015610ed157600080fd5b505af1158015610ee5573d6000803e3d6000fd5b505050508080610ef49061227f565b915050610e51565b5050505050565b610f0b610847565b610f4a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f4190611f5a565b60405180910390fd5b6000815184610f5991906120fc565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600460003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550611011816005546118f790919063ffffffff16565b60058190555061106981600360003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118f790919063ffffffff16565b600360003073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161110a9190611fda565b60405180910390a360005b82518160ff1610156111be578373ffffffffffffffffffffffffffffffffffffffff166383c592cf86858460ff168151811061115457611153612336565b5b60200260200101516040518363ffffffff1660e01b8152600401611179929190611ff5565b600060405180830381600087803b15801561119357600080fd5b505af11580156111a7573d6000803e3d6000fd5b5050505080806111b69061227f565b915050611115565b5050505050565b6000600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b4262015180600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461129b91906120a6565b11156112dc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112d390611fba565b60405180910390fd5b42600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506113376006546005546118f790919063ffffffff16565b600581905550611391600654600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118f790919063ffffffff16565b600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6006546040516114349190611fda565b60405180910390a3565b611446610847565b611485576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161147c90611f5a565b60405180910390fd5b61148e81611923565b50565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415611501576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114f890611f9a565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611571576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161156890611f3a565b60405180910390fd5b80600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258360405161164f9190611fda565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614156116cc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016116c390611f7a565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561173c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161173390611efa565b60405180910390fd5b61178e81600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118d090919063ffffffff16565b600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555061182381600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118f790919063ffffffff16565b600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516118c39190611fda565b60405180910390a3505050565b6000828211156118e3576118e26122a9565b5b81836118ef9190612156565b905092915050565b600080828461190691906120a6565b905083811015611919576119186122a9565b5b8091505092915050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415611993576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161198a90611f1a565b60405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff16600260019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a380600260016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000611a66611a618461205e565b612039565b90508083825260208201905082856020860282011115611a8957611a88612399565b5b60005b85811015611ab95781611a9f8882611b06565b845260208401935060208301925050600181019050611a8c565b5050509392505050565b600081359050611ad2816125bc565b92915050565b600082601f830112611aed57611aec612394565b5b8135611afd848260208601611a53565b91505092915050565b600081359050611b15816125d3565b92915050565b600081359050611b2a816125ea565b92915050565b600060208284031215611b4657611b456123a3565b5b6000611b5484828501611ac3565b91505092915050565b60008060408385031215611b7457611b736123a3565b5b6000611b8285828601611ac3565b9250506020611b9385828601611ac3565b9150509250929050565b600080600060608486031215611bb657611bb56123a3565b5b6000611bc486828701611ac3565b9350506020611bd586828701611ac3565b9250506040611be686828701611b1b565b9150509250925092565b60008060408385031215611c0757611c066123a3565b5b6000611c1585828601611ac3565b9250506020611c2685828601611b1b565b9150509250929050565b600080600060608486031215611c4957611c486123a3565b5b6000611c5786828701611b1b565b9350506020611c6886828701611ac3565b925050604084013567ffffffffffffffff811115611c8957611c8861239e565b5b611c9586828701611ad8565b9150509250925092565b600080600060608486031215611cb857611cb76123a3565b5b6000611cc686828701611b1b565b9350506020611cd786828701611ac3565b9250506040611ce886828701611b06565b9150509250925092565b611cfb8161218a565b82525050565b611d0a8161219c565b82525050565b611d19816121a8565b82525050565b6000611d2a8261208a565b611d348185612095565b9350611d448185602086016121e9565b611d4d816123a8565b840191505092915050565b6000611d65602383612095565b9150611d70826123b9565b604082019050919050565b6000611d88602683612095565b9150611d9382612408565b604082019050919050565b6000611dab602283612095565b9150611db682612457565b604082019050919050565b6000611dce602083612095565b9150611dd9826124a6565b602082019050919050565b6000611df1602583612095565b9150611dfc826124cf565b604082019050919050565b6000611e14602483612095565b9150611e1f8261251e565b604082019050919050565b6000611e37602783612095565b9150611e428261256d565b604082019050919050565b611e56816121d2565b82525050565b611e65816121dc565b82525050565b6000602082019050611e806000830184611cf2565b92915050565b6000606082019050611e9b6000830186611cf2565b611ea86020830185611e4d565b611eb56040830184611d10565b949350505050565b6000602082019050611ed26000830184611d01565b92915050565b60006020820190508181036000830152611ef28184611d1f565b905092915050565b60006020820190508181036000830152611f1381611d58565b9050919050565b60006020820190508181036000830152611f3381611d7b565b9050919050565b60006020820190508181036000830152611f5381611d9e565b9050919050565b60006020820190508181036000830152611f7381611dc1565b9050919050565b60006020820190508181036000830152611f9381611de4565b9050919050565b60006020820190508181036000830152611fb381611e07565b9050919050565b60006020820190508181036000830152611fd381611e2a565b9050919050565b6000602082019050611fef6000830184611e4d565b92915050565b600060408201905061200a6000830185611e4d565b6120176020830184611d10565b9392505050565b60006020820190506120336000830184611e5c565b92915050565b6000612043612054565b905061204f828261224e565b919050565b6000604051905090565b600067ffffffffffffffff82111561207957612078612365565b5b602082029050602081019050919050565b600081519050919050565b600082825260208201905092915050565b60006120b1826121d2565b91506120bc836121d2565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156120f1576120f06122d8565b5b828201905092915050565b6000612107826121d2565b9150612112836121d2565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048311821515161561214b5761214a6122d8565b5b828202905092915050565b6000612161826121d2565b915061216c836121d2565b92508282101561217f5761217e6122d8565b5b828203905092915050565b6000612195826121b2565b9050919050565b60008115159050919050565b6000819050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600060ff82169050919050565b60005b838110156122075780820151818401526020810190506121ec565b83811115612216576000848401525b50505050565b6000600282049050600182168061223457607f821691505b6020821081141561224857612247612307565b5b50919050565b612257826123a8565b810181811067ffffffffffffffff8211171561227657612275612365565b5b80604052505050565b600061228a826121dc565b915060ff82141561229e5761229d6122d8565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b7f6d757374207761697420323420686f757273206265747765656e20666175636560008201527f742063616c6c7300000000000000000000000000000000000000000000000000602082015250565b6125c58161218a565b81146125d057600080fd5b50565b6125dc816121a8565b81146125e757600080fd5b50565b6125f3816121d2565b81146125fe57600080fd5b5056fea264697066735822122052eb56e95e3f90e9dc19112b1197d9e861af1086797fd594e398cf3ba3abdd4364736f6c63430008060033",
}

// TokenOldABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenOldMetaData.ABI instead.
var TokenOldABI = TokenOldMetaData.ABI

// TokenOldBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TokenOldMetaData.Bin instead.
var TokenOldBin = TokenOldMetaData.Bin

// DeployTokenOld deploys a new Ethereum contract, binding an instance of TokenOld to it.
func DeployTokenOld(auth *bind.TransactOpts, backend bind.ContractBackend, _name string, _symbol string, _decimals uint8, total_supply_whole_tokens *big.Int, faucet_amount *big.Int) (common.Address, *types.Transaction, *TokenOld, error) {
	parsed, err := TokenOldMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TokenOldBin), backend, _name, _symbol, _decimals, total_supply_whole_tokens, faucet_amount)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenOld{TokenOldCaller: TokenOldCaller{contract: contract}, TokenOldTransactor: TokenOldTransactor{contract: contract}, TokenOldFilterer: TokenOldFilterer{contract: contract}}, nil
}

// TokenOld is an auto generated Go binding around an Ethereum contract.
type TokenOld struct {
	TokenOldCaller     // Read-only binding to the contract
	TokenOldTransactor // Write-only binding to the contract
	TokenOldFilterer   // Log filterer for contract events
}

// TokenOldCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenOldCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenOldTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenOldTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenOldFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenOldFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenOldSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenOldSession struct {
	Contract     *TokenOld         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenOldCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenOldCallerSession struct {
	Contract *TokenOldCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TokenOldTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenOldTransactorSession struct {
	Contract     *TokenOldTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TokenOldRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenOldRaw struct {
	Contract *TokenOld // Generic contract binding to access the raw methods on
}

// TokenOldCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenOldCallerRaw struct {
	Contract *TokenOldCaller // Generic read-only contract binding to access the raw methods on
}

// TokenOldTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenOldTransactorRaw struct {
	Contract *TokenOldTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenOld creates a new instance of TokenOld, bound to a specific deployed contract.
func NewTokenOld(address common.Address, backend bind.ContractBackend) (*TokenOld, error) {
	contract, err := bindTokenOld(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenOld{TokenOldCaller: TokenOldCaller{contract: contract}, TokenOldTransactor: TokenOldTransactor{contract: contract}, TokenOldFilterer: TokenOldFilterer{contract: contract}}, nil
}

// NewTokenOldCaller creates a new read-only instance of TokenOld, bound to a specific deployed contract.
func NewTokenOldCaller(address common.Address, caller bind.ContractCaller) (*TokenOldCaller, error) {
	contract, err := bindTokenOld(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenOldCaller{contract: contract}, nil
}

// NewTokenOldTransactor creates a new write-only instance of TokenOld, bound to a specific deployed contract.
func NewTokenOldTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenOldTransactor, error) {
	contract, err := bindTokenOld(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenOldTransactor{contract: contract}, nil
}

// NewTokenOldFilterer creates a new log filterer instance of TokenOld, bound to a specific deployed contract.
func NewTokenOldFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenOldFilterer, error) {
	contract, err := bindTokenOld(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenOldFilterer{contract: contract}, nil
}

// bindTokenOld binds a generic wrapper to an already deployed contract.
func bindTokenOld(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenOldABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenOld *TokenOldRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenOld.Contract.TokenOldCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenOld *TokenOldRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenOld.Contract.TokenOldTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenOld *TokenOldRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenOld.Contract.TokenOldTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenOld *TokenOldCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenOld.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenOld *TokenOldTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenOld.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenOld *TokenOldTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenOld.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TokenOld *TokenOldCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenOld.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TokenOld *TokenOldSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TokenOld.Contract.Allowance(&_TokenOld.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TokenOld *TokenOldCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TokenOld.Contract.Allowance(&_TokenOld.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TokenOld *TokenOldCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenOld.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TokenOld *TokenOldSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TokenOld.Contract.BalanceOf(&_TokenOld.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TokenOld *TokenOldCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TokenOld.Contract.BalanceOf(&_TokenOld.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TokenOld *TokenOldCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenOld.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TokenOld *TokenOldSession) Decimals() (uint8, error) {
	return _TokenOld.Contract.Decimals(&_TokenOld.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TokenOld *TokenOldCallerSession) Decimals() (uint8, error) {
	return _TokenOld.Contract.Decimals(&_TokenOld.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TokenOld *TokenOldCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenOld.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TokenOld *TokenOldSession) IsOwner() (bool, error) {
	return _TokenOld.Contract.IsOwner(&_TokenOld.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TokenOld *TokenOldCallerSession) IsOwner() (bool, error) {
	return _TokenOld.Contract.IsOwner(&_TokenOld.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TokenOld *TokenOldCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TokenOld.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TokenOld *TokenOldSession) Name() (string, error) {
	return _TokenOld.Contract.Name(&_TokenOld.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TokenOld *TokenOldCallerSession) Name() (string, error) {
	return _TokenOld.Contract.Name(&_TokenOld.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenOld *TokenOldCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenOld.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenOld *TokenOldSession) Owner() (common.Address, error) {
	return _TokenOld.Contract.Owner(&_TokenOld.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TokenOld *TokenOldCallerSession) Owner() (common.Address, error) {
	return _TokenOld.Contract.Owner(&_TokenOld.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TokenOld *TokenOldCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TokenOld.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TokenOld *TokenOldSession) Symbol() (string, error) {
	return _TokenOld.Contract.Symbol(&_TokenOld.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TokenOld *TokenOldCallerSession) Symbol() (string, error) {
	return _TokenOld.Contract.Symbol(&_TokenOld.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TokenOld *TokenOldCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenOld.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TokenOld *TokenOldSession) TotalSupply() (*big.Int, error) {
	return _TokenOld.Contract.TotalSupply(&_TokenOld.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TokenOld *TokenOldCallerSession) TotalSupply() (*big.Int, error) {
	return _TokenOld.Contract.TotalSupply(&_TokenOld.CallOpts)
}

// AdminDepositBulk is a paid mutator transaction binding the contract method 0xbc36878e.
//
// Solidity: function admin_deposit_bulk(uint256 amount, address bridge_address, bytes32[] vega_public_keys) returns()
func (_TokenOld *TokenOldTransactor) AdminDepositBulk(opts *bind.TransactOpts, amount *big.Int, bridge_address common.Address, vega_public_keys [][32]byte) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "admin_deposit_bulk", amount, bridge_address, vega_public_keys)
}

// AdminDepositBulk is a paid mutator transaction binding the contract method 0xbc36878e.
//
// Solidity: function admin_deposit_bulk(uint256 amount, address bridge_address, bytes32[] vega_public_keys) returns()
func (_TokenOld *TokenOldSession) AdminDepositBulk(amount *big.Int, bridge_address common.Address, vega_public_keys [][32]byte) (*types.Transaction, error) {
	return _TokenOld.Contract.AdminDepositBulk(&_TokenOld.TransactOpts, amount, bridge_address, vega_public_keys)
}

// AdminDepositBulk is a paid mutator transaction binding the contract method 0xbc36878e.
//
// Solidity: function admin_deposit_bulk(uint256 amount, address bridge_address, bytes32[] vega_public_keys) returns()
func (_TokenOld *TokenOldTransactorSession) AdminDepositBulk(amount *big.Int, bridge_address common.Address, vega_public_keys [][32]byte) (*types.Transaction, error) {
	return _TokenOld.Contract.AdminDepositBulk(&_TokenOld.TransactOpts, amount, bridge_address, vega_public_keys)
}

// AdminDepositSingle is a paid mutator transaction binding the contract method 0xb777374c.
//
// Solidity: function admin_deposit_single(uint256 amount, address bridge_address, bytes32 vega_public_key) returns()
func (_TokenOld *TokenOldTransactor) AdminDepositSingle(opts *bind.TransactOpts, amount *big.Int, bridge_address common.Address, vega_public_key [32]byte) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "admin_deposit_single", amount, bridge_address, vega_public_key)
}

// AdminDepositSingle is a paid mutator transaction binding the contract method 0xb777374c.
//
// Solidity: function admin_deposit_single(uint256 amount, address bridge_address, bytes32 vega_public_key) returns()
func (_TokenOld *TokenOldSession) AdminDepositSingle(amount *big.Int, bridge_address common.Address, vega_public_key [32]byte) (*types.Transaction, error) {
	return _TokenOld.Contract.AdminDepositSingle(&_TokenOld.TransactOpts, amount, bridge_address, vega_public_key)
}

// AdminDepositSingle is a paid mutator transaction binding the contract method 0xb777374c.
//
// Solidity: function admin_deposit_single(uint256 amount, address bridge_address, bytes32 vega_public_key) returns()
func (_TokenOld *TokenOldTransactorSession) AdminDepositSingle(amount *big.Int, bridge_address common.Address, vega_public_key [32]byte) (*types.Transaction, error) {
	return _TokenOld.Contract.AdminDepositSingle(&_TokenOld.TransactOpts, amount, bridge_address, vega_public_key)
}

// AdminStakeBulk is a paid mutator transaction binding the contract method 0xd779cae8.
//
// Solidity: function admin_stake_bulk(uint256 amount, address staking_bridge_address, bytes32[] vega_public_keys) returns()
func (_TokenOld *TokenOldTransactor) AdminStakeBulk(opts *bind.TransactOpts, amount *big.Int, staking_bridge_address common.Address, vega_public_keys [][32]byte) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "admin_stake_bulk", amount, staking_bridge_address, vega_public_keys)
}

// AdminStakeBulk is a paid mutator transaction binding the contract method 0xd779cae8.
//
// Solidity: function admin_stake_bulk(uint256 amount, address staking_bridge_address, bytes32[] vega_public_keys) returns()
func (_TokenOld *TokenOldSession) AdminStakeBulk(amount *big.Int, staking_bridge_address common.Address, vega_public_keys [][32]byte) (*types.Transaction, error) {
	return _TokenOld.Contract.AdminStakeBulk(&_TokenOld.TransactOpts, amount, staking_bridge_address, vega_public_keys)
}

// AdminStakeBulk is a paid mutator transaction binding the contract method 0xd779cae8.
//
// Solidity: function admin_stake_bulk(uint256 amount, address staking_bridge_address, bytes32[] vega_public_keys) returns()
func (_TokenOld *TokenOldTransactorSession) AdminStakeBulk(amount *big.Int, staking_bridge_address common.Address, vega_public_keys [][32]byte) (*types.Transaction, error) {
	return _TokenOld.Contract.AdminStakeBulk(&_TokenOld.TransactOpts, amount, staking_bridge_address, vega_public_keys)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TokenOld *TokenOldTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TokenOld *TokenOldSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.Approve(&_TokenOld.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TokenOld *TokenOldTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.Approve(&_TokenOld.TransactOpts, spender, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TokenOld *TokenOldTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TokenOld *TokenOldSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.DecreaseAllowance(&_TokenOld.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TokenOld *TokenOldTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.DecreaseAllowance(&_TokenOld.TransactOpts, spender, subtractedValue)
}

// Faucet is a paid mutator transaction binding the contract method 0xde5f72fd.
//
// Solidity: function faucet() returns()
func (_TokenOld *TokenOldTransactor) Faucet(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "faucet")
}

// Faucet is a paid mutator transaction binding the contract method 0xde5f72fd.
//
// Solidity: function faucet() returns()
func (_TokenOld *TokenOldSession) Faucet() (*types.Transaction, error) {
	return _TokenOld.Contract.Faucet(&_TokenOld.TransactOpts)
}

// Faucet is a paid mutator transaction binding the contract method 0xde5f72fd.
//
// Solidity: function faucet() returns()
func (_TokenOld *TokenOldTransactorSession) Faucet() (*types.Transaction, error) {
	return _TokenOld.Contract.Faucet(&_TokenOld.TransactOpts)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TokenOld *TokenOldTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TokenOld *TokenOldSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.IncreaseAllowance(&_TokenOld.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TokenOld *TokenOldTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.IncreaseAllowance(&_TokenOld.TransactOpts, spender, addedValue)
}

// Issue is a paid mutator transaction binding the contract method 0x867904b4.
//
// Solidity: function issue(address account, uint256 value) returns()
func (_TokenOld *TokenOldTransactor) Issue(opts *bind.TransactOpts, account common.Address, value *big.Int) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "issue", account, value)
}

// Issue is a paid mutator transaction binding the contract method 0x867904b4.
//
// Solidity: function issue(address account, uint256 value) returns()
func (_TokenOld *TokenOldSession) Issue(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.Issue(&_TokenOld.TransactOpts, account, value)
}

// Issue is a paid mutator transaction binding the contract method 0x867904b4.
//
// Solidity: function issue(address account, uint256 value) returns()
func (_TokenOld *TokenOldTransactorSession) Issue(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.Issue(&_TokenOld.TransactOpts, account, value)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_TokenOld *TokenOldTransactor) Kill(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "kill")
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_TokenOld *TokenOldSession) Kill() (*types.Transaction, error) {
	return _TokenOld.Contract.Kill(&_TokenOld.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_TokenOld *TokenOldTransactorSession) Kill() (*types.Transaction, error) {
	return _TokenOld.Contract.Kill(&_TokenOld.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenOld *TokenOldTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenOld *TokenOldSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenOld.Contract.RenounceOwnership(&_TokenOld.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TokenOld *TokenOldTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TokenOld.Contract.RenounceOwnership(&_TokenOld.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TokenOld *TokenOldTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TokenOld *TokenOldSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.Transfer(&_TokenOld.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TokenOld *TokenOldTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.Transfer(&_TokenOld.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TokenOld *TokenOldTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TokenOld *TokenOldSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.TransferFrom(&_TokenOld.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TokenOld *TokenOldTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenOld.Contract.TransferFrom(&_TokenOld.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenOld *TokenOldTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenOld.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenOld *TokenOldSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenOld.Contract.TransferOwnership(&_TokenOld.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TokenOld *TokenOldTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TokenOld.Contract.TransferOwnership(&_TokenOld.TransactOpts, newOwner)
}

// TokenOldApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TokenOld contract.
type TokenOldApprovalIterator struct {
	Event *TokenOldApproval // Event containing the contract specifics and raw log

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
func (it *TokenOldApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenOldApproval)
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
		it.Event = new(TokenOldApproval)
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
func (it *TokenOldApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenOldApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenOldApproval represents a Approval event raised by the TokenOld contract.
type TokenOldApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TokenOld *TokenOldFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TokenOldApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TokenOld.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TokenOldApprovalIterator{contract: _TokenOld.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TokenOld *TokenOldFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TokenOldApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TokenOld.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenOldApproval)
				if err := _TokenOld.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_TokenOld *TokenOldFilterer) ParseApproval(log types.Log) (*TokenOldApproval, error) {
	event := new(TokenOldApproval)
	if err := _TokenOld.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenOldOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TokenOld contract.
type TokenOldOwnershipTransferredIterator struct {
	Event *TokenOldOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TokenOldOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenOldOwnershipTransferred)
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
		it.Event = new(TokenOldOwnershipTransferred)
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
func (it *TokenOldOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenOldOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenOldOwnershipTransferred represents a OwnershipTransferred event raised by the TokenOld contract.
type TokenOldOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenOld *TokenOldFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TokenOldOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenOld.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TokenOldOwnershipTransferredIterator{contract: _TokenOld.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenOld *TokenOldFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TokenOldOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TokenOld.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenOldOwnershipTransferred)
				if err := _TokenOld.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TokenOld *TokenOldFilterer) ParseOwnershipTransferred(log types.Log) (*TokenOldOwnershipTransferred, error) {
	event := new(TokenOldOwnershipTransferred)
	if err := _TokenOld.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenOldTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TokenOld contract.
type TokenOldTransferIterator struct {
	Event *TokenOldTransfer // Event containing the contract specifics and raw log

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
func (it *TokenOldTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenOldTransfer)
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
		it.Event = new(TokenOldTransfer)
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
func (it *TokenOldTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenOldTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenOldTransfer represents a Transfer event raised by the TokenOld contract.
type TokenOldTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TokenOld *TokenOldFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenOldTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenOld.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenOldTransferIterator{contract: _TokenOld.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TokenOld *TokenOldFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TokenOldTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenOld.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenOldTransfer)
				if err := _TokenOld.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_TokenOld *TokenOldFilterer) ParseTransfer(log types.Log) (*TokenOldTransfer, error) {
	event := new(TokenOldTransfer)
	if err := _TokenOld.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
