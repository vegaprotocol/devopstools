// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TokenBase

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

// TokenBaseMetaData contains all meta data concerning the TokenBase contract.
var TokenBaseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burnEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"faucet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"faucetAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"faucetCallLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"burnEnabled_\",\"type\":\"bool\"}],\"name\":\"setBurnEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"faucetAmount_\",\"type\":\"uint256\"}],\"name\":\"setFaucetAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"faucetCallLimit_\",\"type\":\"uint256\"}],\"name\":\"setFaucetCallLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001b7938038062001b79833981016040819052620000349162000345565b838360036200004483826200045e565b5060046200005382826200045e565b50506006805460ff191660ff8516908117909155620000759150600a6200063d565b600755620151806008556009805460ff19166001179055620000983382620000db565b620000a5600033620001c4565b620000d17f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a633620001c4565b505050506200066b565b6001600160a01b038216620001365760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640160405180910390fd5b80600260008282546200014a919062000655565b90915550506001600160a01b038216600090815260208190526040812080548392906200017990849062000655565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35b5050565b620001d0828262000253565b620001c05760008281526005602090815260408083206001600160a01b03851684529091529020805460ff191660011790556200020a3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b505050565b60008281526005602090815260408083206001600160a01b038516845290915290205460ff165b92915050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620002a857600080fd5b81516001600160401b0380821115620002c557620002c562000280565b604051601f8301601f19908116603f01168101908282118183101715620002f057620002f062000280565b816040528381526020925086838588010111156200030d57600080fd5b600091505b8382101562000331578582018301518183018401529082019062000312565b600093810190920192909252949350505050565b600080600080608085870312156200035c57600080fd5b84516001600160401b03808211156200037457600080fd5b620003828883890162000296565b955060208701519150808211156200039957600080fd5b50620003a88782880162000296565b935050604085015160ff81168114620003c057600080fd5b6060959095015193969295505050565b600181811c90821680620003e557607f821691505b6020821081036200040657634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200024e57600081815260208120601f850160051c81016020861015620004355750805b601f850160051c820191505b81811015620004565782815560010162000441565b505050505050565b81516001600160401b038111156200047a576200047a62000280565b62000492816200048b8454620003d0565b846200040c565b602080601f831160018114620004ca5760008415620004b15750858301515b600019600386901b1c1916600185901b17855562000456565b600085815260208120601f198616915b82811015620004fb57888601518255948401946001909101908401620004da565b50858210156200051a5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b634e487b7160e01b600052601160045260246000fd5b600181815b80851115620005815781600019048211156200056557620005656200052a565b808516156200057357918102915b93841c939080029062000545565b509250929050565b6000826200059a575060016200027a565b81620005a9575060006200027a565b8160018114620005c25760028114620005cd57620005ed565b60019150506200027a565b60ff841115620005e157620005e16200052a565b50506001821b6200027a565b5060208310610133831016604e8410600b841016171562000612575081810a6200027a565b6200061e838362000540565b80600019048211156200063557620006356200052a565b029392505050565b60006200064e60ff84168362000589565b9392505050565b808201808211156200027a576200027a6200052a565b6114fe806200067b6000396000f3fe608060405234801561001057600080fd5b50600436106101cf5760003560e01c806370a0823111610104578063a217fddf116100a2578063d547741f11610071578063d547741f146103cb578063dd62ed3e146103de578063de5f72fd14610417578063ee496cf01461041f57600080fd5b8063a217fddf14610388578063a457c2d714610390578063a9059cbb146103a3578063d5391393146103b657600080fd5b806381d2fd9c116100de57806381d2fd9c1461035157806391d148541461036457806395d89b41146103775780639c2814301461037f57600080fd5b806370a082311461030257806379cc67901461032b5780637b2c835f1461033e57600080fd5b8063313ce5671161017157806340c10f191161014b57806340c10f19146102bc57806342966c68146102cf5780635dc96d16146102e257806361106b6b146102ef57600080fd5b8063313ce5671461028157806336568abe1461029657806339509351146102a957600080fd5b806318160ddd116101ad57806318160ddd1461022457806323b872dd14610236578063248a9ca3146102495780632f2ff15d1461026c57600080fd5b806301ffc9a7146101d457806306fdde03146101fc578063095ea7b314610211575b600080fd5b6101e76101e23660046111ac565b610428565b60405190151581526020015b60405180910390f35b61020461045f565b6040516101f391906111fa565b6101e761021f366004611249565b6104f1565b6002545b6040519081526020016101f3565b6101e7610244366004611273565b610509565b6102286102573660046112af565b60009081526005602052604090206001015490565b61027f61027a3660046112c8565b61052d565b005b60065460405160ff90911681526020016101f3565b61027f6102a43660046112c8565b610558565b6101e76102b7366004611249565b6105db565b61027f6102ca366004611249565b61061a565b61027f6102dd3660046112af565b61063d565b6009546101e79060ff1681565b61027f6102fd3660046112af565b61068f565b6102286103103660046112f4565b6001600160a01b031660009081526020819052604090205490565b61027f610339366004611249565b6106ae565b61027f61034c36600461130f565b6106d1565b61027f61035f3660046112af565b6106fe565b6101e76103723660046112c8565b61071d565b610204610748565b61022860075481565b610228600081565b6101e761039e366004611249565b610757565b6101e76103b1366004611249565b6107e9565b6102286000805160206114a983398151915281565b61027f6103d93660046112c8565b6107f7565b6102286103ec366004611331565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b61027f61081d565b61022860085481565b60006001600160e01b03198216637965db0b60e01b148061045957506301ffc9a760e01b6001600160e01b03198316145b92915050565b60606003805461046e9061135b565b80601f016020809104026020016040519081016040528092919081815260200182805461049a9061135b565b80156104e75780601f106104bc576101008083540402835291602001916104e7565b820191906000526020600020905b8154815290600101906020018083116104ca57829003601f168201915b5050505050905090565b6000336104ff818585610907565b5060019392505050565b600033610517858285610a2b565b610522858585610abd565b506001949350505050565b6000828152600560205260409020600101546105498133610c8b565b6105538383610cef565b505050565b6001600160a01b03811633146105cd5760405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201526e103937b632b9903337b91039b2b63360891b60648201526084015b60405180910390fd5b6105d78282610d75565b5050565b3360008181526001602090815260408083206001600160a01b03871684529091528120549091906104ff90829086906106159087906113ab565b610907565b6000805160206114a98339815191526106338133610c8b565b6105538383610ddc565b60095460ff166106825760405162461bcd60e51b815260206004820152601060248201526f189d5c9b881a5cc8191a5cd8589b195960821b60448201526064016105c4565b61068c3382610ebb565b50565b6000805160206114a98339815191526106a88133610c8b565b50600855565b6000805160206114a98339815191526106c78133610c8b565b6105538383610ebb565b6000805160206114a98339815191526106ea8133610c8b565b506009805460ff1916911515919091179055565b6000805160206114a98339815191526107178133610c8b565b50600755565b60009182526005602090815260408084206001600160a01b0393909316845291905290205460ff1690565b60606004805461046e9061135b565b3360008181526001602090815260408083206001600160a01b0387168452909152812054909190838110156107dc5760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f77604482015264207a65726f60d81b60648201526084016105c4565b6105228286868403610907565b6000336104ff818585610abd565b6000828152600560205260409020600101546108138133610c8b565b6105538383610d75565b6000600754116108645760405162461bcd60e51b815260206004820152601260248201527119985d58d95d081a5cc8191a5cd8589b195960721b60448201526064016105c4565b600854336000908152600a60205260409020544291610882916113ab565b11156108e75760405162461bcd60e51b815260206004820152602e60248201527f6d75737420776169742066617563657443616c6c4c696d69742062657477656560448201526d6e206661756365742063616c6c7360901b60648201526084016105c4565b336000818152600a6020526040902042905561090590600754610ddc565b565b6001600160a01b0383166109695760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b60648201526084016105c4565b6001600160a01b0382166109ca5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b60648201526084016105c4565b6001600160a01b0383811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b6001600160a01b038381166000908152600160209081526040808320938616835292905220546000198114610ab75781811015610aaa5760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e636500000060448201526064016105c4565b610ab78484848403610907565b50505050565b6001600160a01b038316610b215760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604482015264647265737360d81b60648201526084016105c4565b6001600160a01b038216610b835760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201526265737360e81b60648201526084016105c4565b6001600160a01b03831660009081526020819052604090205481811015610bfb5760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604482015265616c616e636560d01b60648201526084016105c4565b6001600160a01b03808516600090815260208190526040808220858503905591851681529081208054849290610c329084906113ab565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610c7e91815260200190565b60405180910390a3610ab7565b610c95828261071d565b6105d757610cad816001600160a01b03166014611009565b610cb8836020611009565b604051602001610cc99291906113be565b60408051601f198184030181529082905262461bcd60e51b82526105c4916004016111fa565b610cf9828261071d565b6105d75760008281526005602090815260408083206001600160a01b03851684529091529020805460ff19166001179055610d313390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b610d7f828261071d565b156105d75760008281526005602090815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b6001600160a01b038216610e325760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f20616464726573730060448201526064016105c4565b8060026000828254610e4491906113ab565b90915550506001600160a01b03821660009081526020819052604081208054839290610e719084906113ab565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b6001600160a01b038216610f1b5760405162461bcd60e51b815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f206164647265736044820152607360f81b60648201526084016105c4565b6001600160a01b03821660009081526020819052604090205481811015610f8f5760405162461bcd60e51b815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e604482015261636560f01b60648201526084016105c4565b6001600160a01b0383166000908152602081905260408120838303905560028054849290610fbe908490611433565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b60606000611018836002611446565b6110239060026113ab565b67ffffffffffffffff81111561103b5761103b611465565b6040519080825280601f01601f191660200182016040528015611065576020820181803683370190505b509050600360fc1b816000815181106110805761108061147b565b60200101906001600160f81b031916908160001a905350600f60fb1b816001815181106110af576110af61147b565b60200101906001600160f81b031916908160001a90535060006110d3846002611446565b6110de9060016113ab565b90505b6001811115611156576f181899199a1a9b1b9c1cb0b131b232b360811b85600f16601081106111125761111261147b565b1a60f81b8282815181106111285761112861147b565b60200101906001600160f81b031916908160001a90535060049490941c9361114f81611491565b90506110e1565b5083156111a55760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e7460448201526064016105c4565b9392505050565b6000602082840312156111be57600080fd5b81356001600160e01b0319811681146111a557600080fd5b60005b838110156111f15781810151838201526020016111d9565b50506000910152565b60208152600082518060208401526112198160408501602087016111d6565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461124457600080fd5b919050565b6000806040838503121561125c57600080fd5b6112658361122d565b946020939093013593505050565b60008060006060848603121561128857600080fd5b6112918461122d565b925061129f6020850161122d565b9150604084013590509250925092565b6000602082840312156112c157600080fd5b5035919050565b600080604083850312156112db57600080fd5b823591506112eb6020840161122d565b90509250929050565b60006020828403121561130657600080fd5b6111a58261122d565b60006020828403121561132157600080fd5b813580151581146111a557600080fd5b6000806040838503121561134457600080fd5b61134d8361122d565b91506112eb6020840161122d565b600181811c9082168061136f57607f821691505b60208210810361138f57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b8082018082111561045957610459611395565b7f416363657373436f6e74726f6c3a206163636f756e74200000000000000000008152600083516113f68160178501602088016111d6565b7001034b99036b4b9b9b4b733903937b6329607d1b60179184019182015283516114278160288401602088016111d6565b01602801949350505050565b8181038181111561045957610459611395565b600081600019048311821515161561146057611460611395565b500290565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000816114a0576114a0611395565b50600019019056fe9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6a2646970667358221220005e810e174b4809442fffd26c79a3ce0576dee93286c5c6be3d0be978bc6a7464736f6c63430008100033",
}

// TokenBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenBaseMetaData.ABI instead.
var TokenBaseABI = TokenBaseMetaData.ABI

// TokenBaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TokenBaseMetaData.Bin instead.
var TokenBaseBin = TokenBaseMetaData.Bin

// DeployTokenBase deploys a new Ethereum contract, binding an instance of TokenBase to it.
func DeployTokenBase(auth *bind.TransactOpts, backend bind.ContractBackend, name_ string, symbol_ string, decimals_ uint8, totalSupply_ *big.Int) (common.Address, *types.Transaction, *TokenBase, error) {
	parsed, err := TokenBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TokenBaseBin), backend, name_, symbol_, decimals_, totalSupply_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenBase{TokenBaseCaller: TokenBaseCaller{contract: contract}, TokenBaseTransactor: TokenBaseTransactor{contract: contract}, TokenBaseFilterer: TokenBaseFilterer{contract: contract}}, nil
}

// TokenBase is an auto generated Go binding around an Ethereum contract.
type TokenBase struct {
	TokenBaseCaller     // Read-only binding to the contract
	TokenBaseTransactor // Write-only binding to the contract
	TokenBaseFilterer   // Log filterer for contract events
}

// TokenBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenBaseSession struct {
	Contract     *TokenBase        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenBaseCallerSession struct {
	Contract *TokenBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TokenBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenBaseTransactorSession struct {
	Contract     *TokenBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TokenBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenBaseRaw struct {
	Contract *TokenBase // Generic contract binding to access the raw methods on
}

// TokenBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenBaseCallerRaw struct {
	Contract *TokenBaseCaller // Generic read-only contract binding to access the raw methods on
}

// TokenBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenBaseTransactorRaw struct {
	Contract *TokenBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenBase creates a new instance of TokenBase, bound to a specific deployed contract.
func NewTokenBase(address common.Address, backend bind.ContractBackend) (*TokenBase, error) {
	contract, err := bindTokenBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenBase{TokenBaseCaller: TokenBaseCaller{contract: contract}, TokenBaseTransactor: TokenBaseTransactor{contract: contract}, TokenBaseFilterer: TokenBaseFilterer{contract: contract}}, nil
}

// NewTokenBaseCaller creates a new read-only instance of TokenBase, bound to a specific deployed contract.
func NewTokenBaseCaller(address common.Address, caller bind.ContractCaller) (*TokenBaseCaller, error) {
	contract, err := bindTokenBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenBaseCaller{contract: contract}, nil
}

// NewTokenBaseTransactor creates a new write-only instance of TokenBase, bound to a specific deployed contract.
func NewTokenBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenBaseTransactor, error) {
	contract, err := bindTokenBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenBaseTransactor{contract: contract}, nil
}

// NewTokenBaseFilterer creates a new log filterer instance of TokenBase, bound to a specific deployed contract.
func NewTokenBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenBaseFilterer, error) {
	contract, err := bindTokenBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenBaseFilterer{contract: contract}, nil
}

// bindTokenBase binds a generic wrapper to an already deployed contract.
func bindTokenBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenBaseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenBase *TokenBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenBase.Contract.TokenBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenBase *TokenBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenBase.Contract.TokenBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenBase *TokenBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenBase.Contract.TokenBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenBase *TokenBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenBase *TokenBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenBase *TokenBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenBase.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TokenBase *TokenBaseCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TokenBase *TokenBaseSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TokenBase.Contract.DEFAULTADMINROLE(&_TokenBase.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_TokenBase *TokenBaseCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _TokenBase.Contract.DEFAULTADMINROLE(&_TokenBase.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_TokenBase *TokenBaseCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_TokenBase *TokenBaseSession) MINTERROLE() ([32]byte, error) {
	return _TokenBase.Contract.MINTERROLE(&_TokenBase.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_TokenBase *TokenBaseCallerSession) MINTERROLE() ([32]byte, error) {
	return _TokenBase.Contract.MINTERROLE(&_TokenBase.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TokenBase *TokenBaseCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TokenBase *TokenBaseSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TokenBase.Contract.Allowance(&_TokenBase.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TokenBase *TokenBaseCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TokenBase.Contract.Allowance(&_TokenBase.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TokenBase *TokenBaseCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TokenBase *TokenBaseSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TokenBase.Contract.BalanceOf(&_TokenBase.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TokenBase *TokenBaseCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TokenBase.Contract.BalanceOf(&_TokenBase.CallOpts, account)
}

// BurnEnabled is a free data retrieval call binding the contract method 0x5dc96d16.
//
// Solidity: function burnEnabled() view returns(bool)
func (_TokenBase *TokenBaseCaller) BurnEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "burnEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// BurnEnabled is a free data retrieval call binding the contract method 0x5dc96d16.
//
// Solidity: function burnEnabled() view returns(bool)
func (_TokenBase *TokenBaseSession) BurnEnabled() (bool, error) {
	return _TokenBase.Contract.BurnEnabled(&_TokenBase.CallOpts)
}

// BurnEnabled is a free data retrieval call binding the contract method 0x5dc96d16.
//
// Solidity: function burnEnabled() view returns(bool)
func (_TokenBase *TokenBaseCallerSession) BurnEnabled() (bool, error) {
	return _TokenBase.Contract.BurnEnabled(&_TokenBase.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TokenBase *TokenBaseCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TokenBase *TokenBaseSession) Decimals() (uint8, error) {
	return _TokenBase.Contract.Decimals(&_TokenBase.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TokenBase *TokenBaseCallerSession) Decimals() (uint8, error) {
	return _TokenBase.Contract.Decimals(&_TokenBase.CallOpts)
}

// FaucetAmount is a free data retrieval call binding the contract method 0x9c281430.
//
// Solidity: function faucetAmount() view returns(uint256)
func (_TokenBase *TokenBaseCaller) FaucetAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "faucetAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FaucetAmount is a free data retrieval call binding the contract method 0x9c281430.
//
// Solidity: function faucetAmount() view returns(uint256)
func (_TokenBase *TokenBaseSession) FaucetAmount() (*big.Int, error) {
	return _TokenBase.Contract.FaucetAmount(&_TokenBase.CallOpts)
}

// FaucetAmount is a free data retrieval call binding the contract method 0x9c281430.
//
// Solidity: function faucetAmount() view returns(uint256)
func (_TokenBase *TokenBaseCallerSession) FaucetAmount() (*big.Int, error) {
	return _TokenBase.Contract.FaucetAmount(&_TokenBase.CallOpts)
}

// FaucetCallLimit is a free data retrieval call binding the contract method 0xee496cf0.
//
// Solidity: function faucetCallLimit() view returns(uint256)
func (_TokenBase *TokenBaseCaller) FaucetCallLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "faucetCallLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FaucetCallLimit is a free data retrieval call binding the contract method 0xee496cf0.
//
// Solidity: function faucetCallLimit() view returns(uint256)
func (_TokenBase *TokenBaseSession) FaucetCallLimit() (*big.Int, error) {
	return _TokenBase.Contract.FaucetCallLimit(&_TokenBase.CallOpts)
}

// FaucetCallLimit is a free data retrieval call binding the contract method 0xee496cf0.
//
// Solidity: function faucetCallLimit() view returns(uint256)
func (_TokenBase *TokenBaseCallerSession) FaucetCallLimit() (*big.Int, error) {
	return _TokenBase.Contract.FaucetCallLimit(&_TokenBase.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TokenBase *TokenBaseCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TokenBase *TokenBaseSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TokenBase.Contract.GetRoleAdmin(&_TokenBase.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_TokenBase *TokenBaseCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _TokenBase.Contract.GetRoleAdmin(&_TokenBase.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TokenBase *TokenBaseCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TokenBase *TokenBaseSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TokenBase.Contract.HasRole(&_TokenBase.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_TokenBase *TokenBaseCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _TokenBase.Contract.HasRole(&_TokenBase.CallOpts, role, account)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TokenBase *TokenBaseCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TokenBase *TokenBaseSession) Name() (string, error) {
	return _TokenBase.Contract.Name(&_TokenBase.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TokenBase *TokenBaseCallerSession) Name() (string, error) {
	return _TokenBase.Contract.Name(&_TokenBase.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TokenBase *TokenBaseCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TokenBase *TokenBaseSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TokenBase.Contract.SupportsInterface(&_TokenBase.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TokenBase *TokenBaseCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TokenBase.Contract.SupportsInterface(&_TokenBase.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TokenBase *TokenBaseCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TokenBase *TokenBaseSession) Symbol() (string, error) {
	return _TokenBase.Contract.Symbol(&_TokenBase.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TokenBase *TokenBaseCallerSession) Symbol() (string, error) {
	return _TokenBase.Contract.Symbol(&_TokenBase.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TokenBase *TokenBaseCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenBase.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TokenBase *TokenBaseSession) TotalSupply() (*big.Int, error) {
	return _TokenBase.Contract.TotalSupply(&_TokenBase.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TokenBase *TokenBaseCallerSession) TotalSupply() (*big.Int, error) {
	return _TokenBase.Contract.TotalSupply(&_TokenBase.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TokenBase *TokenBaseTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TokenBase *TokenBaseSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.Approve(&_TokenBase.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TokenBase *TokenBaseTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.Approve(&_TokenBase.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_TokenBase *TokenBaseTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_TokenBase *TokenBaseSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.Burn(&_TokenBase.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_TokenBase *TokenBaseTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.Burn(&_TokenBase.TransactOpts, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_TokenBase *TokenBaseTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "burnFrom", account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_TokenBase *TokenBaseSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.BurnFrom(&_TokenBase.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_TokenBase *TokenBaseTransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.BurnFrom(&_TokenBase.TransactOpts, account, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TokenBase *TokenBaseTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TokenBase *TokenBaseSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.DecreaseAllowance(&_TokenBase.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TokenBase *TokenBaseTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.DecreaseAllowance(&_TokenBase.TransactOpts, spender, subtractedValue)
}

// Faucet is a paid mutator transaction binding the contract method 0xde5f72fd.
//
// Solidity: function faucet() returns()
func (_TokenBase *TokenBaseTransactor) Faucet(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "faucet")
}

// Faucet is a paid mutator transaction binding the contract method 0xde5f72fd.
//
// Solidity: function faucet() returns()
func (_TokenBase *TokenBaseSession) Faucet() (*types.Transaction, error) {
	return _TokenBase.Contract.Faucet(&_TokenBase.TransactOpts)
}

// Faucet is a paid mutator transaction binding the contract method 0xde5f72fd.
//
// Solidity: function faucet() returns()
func (_TokenBase *TokenBaseTransactorSession) Faucet() (*types.Transaction, error) {
	return _TokenBase.Contract.Faucet(&_TokenBase.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TokenBase *TokenBaseTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TokenBase *TokenBaseSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TokenBase.Contract.GrantRole(&_TokenBase.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_TokenBase *TokenBaseTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TokenBase.Contract.GrantRole(&_TokenBase.TransactOpts, role, account)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TokenBase *TokenBaseTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TokenBase *TokenBaseSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.IncreaseAllowance(&_TokenBase.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TokenBase *TokenBaseTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.IncreaseAllowance(&_TokenBase.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TokenBase *TokenBaseTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TokenBase *TokenBaseSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.Mint(&_TokenBase.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TokenBase *TokenBaseTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.Mint(&_TokenBase.TransactOpts, to, amount)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_TokenBase *TokenBaseTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_TokenBase *TokenBaseSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TokenBase.Contract.RenounceRole(&_TokenBase.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_TokenBase *TokenBaseTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TokenBase.Contract.RenounceRole(&_TokenBase.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TokenBase *TokenBaseTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TokenBase *TokenBaseSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TokenBase.Contract.RevokeRole(&_TokenBase.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_TokenBase *TokenBaseTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _TokenBase.Contract.RevokeRole(&_TokenBase.TransactOpts, role, account)
}

// SetBurnEnabled is a paid mutator transaction binding the contract method 0x7b2c835f.
//
// Solidity: function setBurnEnabled(bool burnEnabled_) returns()
func (_TokenBase *TokenBaseTransactor) SetBurnEnabled(opts *bind.TransactOpts, burnEnabled_ bool) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "setBurnEnabled", burnEnabled_)
}

// SetBurnEnabled is a paid mutator transaction binding the contract method 0x7b2c835f.
//
// Solidity: function setBurnEnabled(bool burnEnabled_) returns()
func (_TokenBase *TokenBaseSession) SetBurnEnabled(burnEnabled_ bool) (*types.Transaction, error) {
	return _TokenBase.Contract.SetBurnEnabled(&_TokenBase.TransactOpts, burnEnabled_)
}

// SetBurnEnabled is a paid mutator transaction binding the contract method 0x7b2c835f.
//
// Solidity: function setBurnEnabled(bool burnEnabled_) returns()
func (_TokenBase *TokenBaseTransactorSession) SetBurnEnabled(burnEnabled_ bool) (*types.Transaction, error) {
	return _TokenBase.Contract.SetBurnEnabled(&_TokenBase.TransactOpts, burnEnabled_)
}

// SetFaucetAmount is a paid mutator transaction binding the contract method 0x81d2fd9c.
//
// Solidity: function setFaucetAmount(uint256 faucetAmount_) returns()
func (_TokenBase *TokenBaseTransactor) SetFaucetAmount(opts *bind.TransactOpts, faucetAmount_ *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "setFaucetAmount", faucetAmount_)
}

// SetFaucetAmount is a paid mutator transaction binding the contract method 0x81d2fd9c.
//
// Solidity: function setFaucetAmount(uint256 faucetAmount_) returns()
func (_TokenBase *TokenBaseSession) SetFaucetAmount(faucetAmount_ *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.SetFaucetAmount(&_TokenBase.TransactOpts, faucetAmount_)
}

// SetFaucetAmount is a paid mutator transaction binding the contract method 0x81d2fd9c.
//
// Solidity: function setFaucetAmount(uint256 faucetAmount_) returns()
func (_TokenBase *TokenBaseTransactorSession) SetFaucetAmount(faucetAmount_ *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.SetFaucetAmount(&_TokenBase.TransactOpts, faucetAmount_)
}

// SetFaucetCallLimit is a paid mutator transaction binding the contract method 0x61106b6b.
//
// Solidity: function setFaucetCallLimit(uint256 faucetCallLimit_) returns()
func (_TokenBase *TokenBaseTransactor) SetFaucetCallLimit(opts *bind.TransactOpts, faucetCallLimit_ *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "setFaucetCallLimit", faucetCallLimit_)
}

// SetFaucetCallLimit is a paid mutator transaction binding the contract method 0x61106b6b.
//
// Solidity: function setFaucetCallLimit(uint256 faucetCallLimit_) returns()
func (_TokenBase *TokenBaseSession) SetFaucetCallLimit(faucetCallLimit_ *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.SetFaucetCallLimit(&_TokenBase.TransactOpts, faucetCallLimit_)
}

// SetFaucetCallLimit is a paid mutator transaction binding the contract method 0x61106b6b.
//
// Solidity: function setFaucetCallLimit(uint256 faucetCallLimit_) returns()
func (_TokenBase *TokenBaseTransactorSession) SetFaucetCallLimit(faucetCallLimit_ *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.SetFaucetCallLimit(&_TokenBase.TransactOpts, faucetCallLimit_)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_TokenBase *TokenBaseTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_TokenBase *TokenBaseSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.Transfer(&_TokenBase.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_TokenBase *TokenBaseTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.Transfer(&_TokenBase.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_TokenBase *TokenBaseTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_TokenBase *TokenBaseSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.TransferFrom(&_TokenBase.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_TokenBase *TokenBaseTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenBase.Contract.TransferFrom(&_TokenBase.TransactOpts, from, to, amount)
}

// TokenBaseApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TokenBase contract.
type TokenBaseApprovalIterator struct {
	Event *TokenBaseApproval // Event containing the contract specifics and raw log

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
func (it *TokenBaseApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenBaseApproval)
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
		it.Event = new(TokenBaseApproval)
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
func (it *TokenBaseApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenBaseApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenBaseApproval represents a Approval event raised by the TokenBase contract.
type TokenBaseApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TokenBase *TokenBaseFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TokenBaseApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TokenBase.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TokenBaseApprovalIterator{contract: _TokenBase.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TokenBase *TokenBaseFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TokenBaseApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TokenBase.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenBaseApproval)
				if err := _TokenBase.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_TokenBase *TokenBaseFilterer) ParseApproval(log types.Log) (*TokenBaseApproval, error) {
	event := new(TokenBaseApproval)
	if err := _TokenBase.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenBaseRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the TokenBase contract.
type TokenBaseRoleAdminChangedIterator struct {
	Event *TokenBaseRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *TokenBaseRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenBaseRoleAdminChanged)
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
		it.Event = new(TokenBaseRoleAdminChanged)
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
func (it *TokenBaseRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenBaseRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenBaseRoleAdminChanged represents a RoleAdminChanged event raised by the TokenBase contract.
type TokenBaseRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TokenBase *TokenBaseFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*TokenBaseRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _TokenBase.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &TokenBaseRoleAdminChangedIterator{contract: _TokenBase.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TokenBase *TokenBaseFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *TokenBaseRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _TokenBase.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenBaseRoleAdminChanged)
				if err := _TokenBase.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_TokenBase *TokenBaseFilterer) ParseRoleAdminChanged(log types.Log) (*TokenBaseRoleAdminChanged, error) {
	event := new(TokenBaseRoleAdminChanged)
	if err := _TokenBase.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenBaseRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the TokenBase contract.
type TokenBaseRoleGrantedIterator struct {
	Event *TokenBaseRoleGranted // Event containing the contract specifics and raw log

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
func (it *TokenBaseRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenBaseRoleGranted)
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
		it.Event = new(TokenBaseRoleGranted)
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
func (it *TokenBaseRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenBaseRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenBaseRoleGranted represents a RoleGranted event raised by the TokenBase contract.
type TokenBaseRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TokenBase *TokenBaseFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TokenBaseRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TokenBase.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TokenBaseRoleGrantedIterator{contract: _TokenBase.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TokenBase *TokenBaseFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *TokenBaseRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TokenBase.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenBaseRoleGranted)
				if err := _TokenBase.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_TokenBase *TokenBaseFilterer) ParseRoleGranted(log types.Log) (*TokenBaseRoleGranted, error) {
	event := new(TokenBaseRoleGranted)
	if err := _TokenBase.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenBaseRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the TokenBase contract.
type TokenBaseRoleRevokedIterator struct {
	Event *TokenBaseRoleRevoked // Event containing the contract specifics and raw log

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
func (it *TokenBaseRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenBaseRoleRevoked)
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
		it.Event = new(TokenBaseRoleRevoked)
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
func (it *TokenBaseRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenBaseRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenBaseRoleRevoked represents a RoleRevoked event raised by the TokenBase contract.
type TokenBaseRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TokenBase *TokenBaseFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*TokenBaseRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TokenBase.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &TokenBaseRoleRevokedIterator{contract: _TokenBase.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TokenBase *TokenBaseFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *TokenBaseRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _TokenBase.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenBaseRoleRevoked)
				if err := _TokenBase.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_TokenBase *TokenBaseFilterer) ParseRoleRevoked(log types.Log) (*TokenBaseRoleRevoked, error) {
	event := new(TokenBaseRoleRevoked)
	if err := _TokenBase.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenBaseTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TokenBase contract.
type TokenBaseTransferIterator struct {
	Event *TokenBaseTransfer // Event containing the contract specifics and raw log

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
func (it *TokenBaseTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenBaseTransfer)
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
		it.Event = new(TokenBaseTransfer)
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
func (it *TokenBaseTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenBaseTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenBaseTransfer represents a Transfer event raised by the TokenBase contract.
type TokenBaseTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TokenBase *TokenBaseFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TokenBaseTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenBase.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TokenBaseTransferIterator{contract: _TokenBase.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TokenBase *TokenBaseFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TokenBaseTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TokenBase.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenBaseTransfer)
				if err := _TokenBase.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_TokenBase *TokenBaseFilterer) ParseTransfer(log types.Log) (*TokenBaseTransfer, error) {
	event := new(TokenBaseTransfer)
	if err := _TokenBase.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
