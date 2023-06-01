// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ERC20Bridge

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

// ERC20BridgeMetaData contains all meta data concerning the ERC20Bridge contract.
var ERC20BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"erc20_asset_pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"multisig_control\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"new_maximum\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Deposit_Maximum_Set\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"new_minimum\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Deposit_Minimum_Set\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user_address\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"Asset_Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"vega_asset_id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Listed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Removed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user_address\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Asset_Withdrawn\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"vega_public_key\",\"type\":\"bytes32\"}],\"name\":\"deposit_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"vega_asset_id\",\"type\":\"bytes32\"}],\"name\":\"get_asset_source\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"get_deposit_maximum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"get_deposit_minimum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_multisig_control_address\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"get_vega_asset_id\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"}],\"name\":\"is_asset_listed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"vega_asset_id\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"list_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"remove_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maximum_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_deposit_maximum\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimum_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"set_deposit_minimum\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset_source\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signatures\",\"type\":\"bytes\"}],\"name\":\"withdraw_asset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200234b3803806200234b83398181016040528101906200003791906200016f565b81600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050620001b6565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620000f282620000c5565b9050919050565b6200010481620000e5565b81146200011057600080fd5b50565b6000815190506200012481620000f9565b92915050565b60006200013782620000c5565b9050919050565b62000149816200012a565b81146200015557600080fd5b50565b60008151905062000169816200013e565b92915050565b60008060408385031215620001895762000188620000c0565b5b6000620001998582860162000113565b9250506020620001ac8582860162000158565b9150509250929050565b61218580620001c66000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c8063927ee5b911610071578063927ee5b9146101b1578063a06b5d39146101cd578063a8780cda146101fd578063c58dc3b914610219578063c76de35814610237578063f768393214610253576100b4565b80631d501b5d146100b9578063292463b1146100e95780633882b3da146101055780634322b1f214610121578063786b0bc0146101515780637fd27b7f14610181575b600080fd5b6100d360048036038101906100ce91906113d2565b61026f565b6040516100e09190611418565b60405180910390f35b61010360048036038101906100fe91906115a5565b6102b8565b005b61011f600480360381019061011a91906115a5565b6104f5565b005b61013b600480360381019061013691906113d2565b610732565b6040516101489190611418565b60405180910390f35b61016b6004803603810190610166919061165e565b61077b565b604051610178919061169a565b60405180910390f35b61019b600480360381019061019691906113d2565b6107b8565b6040516101a891906116d0565b60405180910390f35b6101cb60048036038101906101c691906116eb565b61080e565b005b6101e760048036038101906101e291906113d2565b610a85565b6040516101f49190611791565b60405180910390f35b610217600480360381019061021291906117ac565b610ace565b005b610221610db5565b60405161022e919061169a565b60405180910390f35b610251600480360381019061024c919061182f565b610dde565b005b61026d6004803603810190610268919061189e565b61102a565b005b6000600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610344576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161033b9061194e565b60405180910390fd5b600084848460405160200161035b939291906119ba565b604051602081830303815290604052905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b81526004016103c993929190611a8c565b602060405180830381600087803b1580156103e357600080fd5b505af11580156103f7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041b9190611afd565b61045a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161045190611b76565b60405180910390fd5b83600460008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508473ffffffffffffffffffffffffffffffffffffffff167f8861dbda4052c4bf7ba24aad82097ef1e6ade8e6f1b42341640db8a12abfa99385856040516104e6929190611b96565b60405180910390a25050505050565b600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610581576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105789061194e565b60405180910390fd5b600084848460405160200161059893929190611c0b565b604051602081830303815290604052905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b815260040161060693929190611a8c565b602060405180830381600087803b15801561062057600080fd5b505af1158015610634573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106589190611afd565b610697576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161068e90611b76565b60405180910390fd5b83600360008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508473ffffffffffffffffffffffffffffffffffffffff167f4ed0df0b169b573722ecdcc12333646a6efd0445c28fe277470bbfca620e8ad58585604051610723929190611b96565b60405180910390a25050505050565b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60006005600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b6000600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050919050565b6000858585856040516020016108279493929190611ca1565b604051602081830303815290604052905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b815260040161089593929190611a8c565b602060405180830381600087803b1580156108af57600080fd5b505af11580156108c3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108e79190611afd565b610926576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161091d90611b76565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d9caed128786886040518463ffffffff1660e01b815260040161098593929190611cf9565b602060405180830381600087803b15801561099f57600080fd5b505af11580156109b3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109d79190611afd565b610a16576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a0d90611da2565b60405180910390fd5b8573ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df8786604051610a75929190611b96565b60405180910390a3505050505050565b6000600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615610b5b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b5290611e0e565b60405180910390fd5b6000848484604051602001610b7293929190611e7a565b604051602081830303815290604052905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b8152600401610be093929190611a8c565b602060405180830381600087803b158015610bfa57600080fd5b505af1158015610c0e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c329190611afd565b610c71576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c6890611b76565b60405180910390fd5b6001600260008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550846005600086815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555083600660008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550838573ffffffffffffffffffffffffffffffffffffffff167f4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e85604051610da69190611418565b60405180910390a35050505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610e6a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e619061194e565b60405180910390fd5b60008383604051602001610e7f929190611f10565b604051602081830303815290604052905060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ba73659a8383866040518463ffffffff1660e01b8152600401610eed93929190611a8c565b602060405180830381600087803b158015610f0757600080fd5b505af1158015610f1b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f3f9190611afd565b610f7e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f7590611b76565b60405180910390fd5b6000600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508373ffffffffffffffffffffffffffffffffffffffff167f58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c98460405161101c9190611418565b60405180910390a250505050565b600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166110b6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110ad9061194e565b60405180910390fd5b6000600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414806111435750600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020548211155b611182576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161117990611f98565b60405180910390fd5b600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054821015611204576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111fb90612004565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff166323b872dd33600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16856040518463ffffffff1660e01b815260040161126393929190612083565b602060405180830381600087803b15801561127d57600080fd5b505af1158015611291573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112b59190611afd565b6112f4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112eb90612106565b60405180910390fd5b8273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf8484604051611353929190612126565b60405180910390a3505050565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061139f82611374565b9050919050565b6113af81611394565b81146113ba57600080fd5b50565b6000813590506113cc816113a6565b92915050565b6000602082840312156113e8576113e761136a565b5b60006113f6848285016113bd565b91505092915050565b6000819050919050565b611412816113ff565b82525050565b600060208201905061142d6000830184611409565b92915050565b61143c816113ff565b811461144757600080fd5b50565b60008135905061145981611433565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6114b282611469565b810181811067ffffffffffffffff821117156114d1576114d061147a565b5b80604052505050565b60006114e4611360565b90506114f082826114a9565b919050565b600067ffffffffffffffff8211156115105761150f61147a565b5b61151982611469565b9050602081019050919050565b82818337600083830152505050565b6000611548611543846114f5565b6114da565b90508281526020810184848401111561156457611563611464565b5b61156f848285611526565b509392505050565b600082601f83011261158c5761158b61145f565b5b813561159c848260208601611535565b91505092915050565b600080600080608085870312156115bf576115be61136a565b5b60006115cd878288016113bd565b94505060206115de8782880161144a565b93505060406115ef8782880161144a565b925050606085013567ffffffffffffffff8111156116105761160f61136f565b5b61161c87828801611577565b91505092959194509250565b6000819050919050565b61163b81611628565b811461164657600080fd5b50565b60008135905061165881611632565b92915050565b6000602082840312156116745761167361136a565b5b600061168284828501611649565b91505092915050565b61169481611394565b82525050565b60006020820190506116af600083018461168b565b92915050565b60008115159050919050565b6116ca816116b5565b82525050565b60006020820190506116e560008301846116c1565b92915050565b600080600080600060a086880312156117075761170661136a565b5b6000611715888289016113bd565b95505060206117268882890161144a565b9450506040611737888289016113bd565b93505060606117488882890161144a565b925050608086013567ffffffffffffffff8111156117695761176861136f565b5b61177588828901611577565b9150509295509295909350565b61178b81611628565b82525050565b60006020820190506117a66000830184611782565b92915050565b600080600080608085870312156117c6576117c561136a565b5b60006117d4878288016113bd565b94505060206117e587828801611649565b93505060406117f68782880161144a565b925050606085013567ffffffffffffffff8111156118175761181661136f565b5b61182387828801611577565b91505092959194509250565b6000806000606084860312156118485761184761136a565b5b6000611856868287016113bd565b93505060206118678682870161144a565b925050604084013567ffffffffffffffff8111156118885761188761136f565b5b61189486828701611577565b9150509250925092565b6000806000606084860312156118b7576118b661136a565b5b60006118c5868287016113bd565b93505060206118d68682870161144a565b92505060406118e786828701611649565b9150509250925092565b600082825260208201905092915050565b7f6173736574206e6f74206c697374656400000000000000000000000000000000600082015250565b60006119386010836118f1565b915061194382611902565b602082019050919050565b600060208201905081810360008301526119678161192b565b9050919050565b7f7365745f6465706f7369745f6d6178696d756d00000000000000000000000000600082015250565b60006119a46013836118f1565b91506119af8261196e565b602082019050919050565b60006080820190506119cf600083018661168b565b6119dc6020830185611409565b6119e96040830184611409565b81810360608301526119fa81611997565b9050949350505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015611a3e578082015181840152602081019050611a23565b83811115611a4d576000848401525b50505050565b6000611a5e82611a04565b611a688185611a0f565b9350611a78818560208601611a20565b611a8181611469565b840191505092915050565b60006060820190508181036000830152611aa68186611a53565b90508181036020830152611aba8185611a53565b9050611ac96040830184611409565b949350505050565b611ada816116b5565b8114611ae557600080fd5b50565b600081519050611af781611ad1565b92915050565b600060208284031215611b1357611b1261136a565b5b6000611b2184828501611ae8565b91505092915050565b7f626164207369676e617475726573000000000000000000000000000000000000600082015250565b6000611b60600e836118f1565b9150611b6b82611b2a565b602082019050919050565b60006020820190508181036000830152611b8f81611b53565b9050919050565b6000604082019050611bab6000830185611409565b611bb86020830184611409565b9392505050565b7f7365745f6465706f7369745f6d696e696d756d00000000000000000000000000600082015250565b6000611bf56013836118f1565b9150611c0082611bbf565b602082019050919050565b6000608082019050611c20600083018661168b565b611c2d6020830185611409565b611c3a6040830184611409565b8181036060830152611c4b81611be8565b9050949350505050565b7f77697468647261775f6173736574000000000000000000000000000000000000600082015250565b6000611c8b600e836118f1565b9150611c9682611c55565b602082019050919050565b600060a082019050611cb6600083018761168b565b611cc36020830186611409565b611cd0604083018561168b565b611cdd6060830184611409565b8181036080830152611cee81611c7e565b905095945050505050565b6000606082019050611d0e600083018661168b565b611d1b602083018561168b565b611d286040830184611409565b949350505050565b7f746f6b656e206469646e2774207472616e736665722c2072656a65637465642060008201527f627920617373657420706f6f6c2e000000000000000000000000000000000000602082015250565b6000611d8c602e836118f1565b9150611d9782611d30565b604082019050919050565b60006020820190508181036000830152611dbb81611d7f565b9050919050565b7f617373657420616c7265616479206c6973746564000000000000000000000000600082015250565b6000611df86014836118f1565b9150611e0382611dc2565b602082019050919050565b60006020820190508181036000830152611e2781611deb565b9050919050565b7f6c6973745f617373657400000000000000000000000000000000000000000000600082015250565b6000611e64600a836118f1565b9150611e6f82611e2e565b602082019050919050565b6000608082019050611e8f600083018661168b565b611e9c6020830185611782565b611ea96040830184611409565b8181036060830152611eba81611e57565b9050949350505050565b7f72656d6f76655f61737365740000000000000000000000000000000000000000600082015250565b6000611efa600c836118f1565b9150611f0582611ec4565b602082019050919050565b6000606082019050611f25600083018561168b565b611f326020830184611409565b8181036040830152611f4381611eed565b90509392505050565b7f6465706f7369742061626f7665206d6178696d756d0000000000000000000000600082015250565b6000611f826015836118f1565b9150611f8d82611f4c565b602082019050919050565b60006020820190508181036000830152611fb181611f75565b9050919050565b7f6465706f7369742062656c6f77206d696e696d756d0000000000000000000000600082015250565b6000611fee6015836118f1565b9150611ff982611fb8565b602082019050919050565b6000602082019050818103600083015261201d81611fe1565b9050919050565b6000819050919050565b600061204961204461203f84611374565b612024565b611374565b9050919050565b600061205b8261202e565b9050919050565b600061206d82612050565b9050919050565b61207d81612062565b82525050565b6000606082019050612098600083018661168b565b6120a56020830185612074565b6120b26040830184611409565b949350505050565b7f7472616e73666572206661696c656420696e206465706f736974000000000000600082015250565b60006120f0601a836118f1565b91506120fb826120ba565b602082019050919050565b6000602082019050818103600083015261211f816120e3565b9050919050565b600060408201905061213b6000830185611409565b6121486020830184611782565b939250505056fea2646970667358221220dda9939ad5c1380c7f934bf19db31051c648ae958f3a9e7ae3c88eadc59a4b8464736f6c63430008080033",
}

// ERC20BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20BridgeMetaData.ABI instead.
var ERC20BridgeABI = ERC20BridgeMetaData.ABI

// ERC20BridgeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC20BridgeMetaData.Bin instead.
var ERC20BridgeBin = ERC20BridgeMetaData.Bin

// DeployERC20Bridge deploys a new Ethereum contract, binding an instance of ERC20Bridge to it.
func DeployERC20Bridge(auth *bind.TransactOpts, backend bind.ContractBackend, erc20_asset_pool common.Address, multisig_control common.Address) (common.Address, *types.Transaction, *ERC20Bridge, error) {
	parsed, err := ERC20BridgeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC20BridgeBin), backend, erc20_asset_pool, multisig_control)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Bridge{ERC20BridgeCaller: ERC20BridgeCaller{contract: contract}, ERC20BridgeTransactor: ERC20BridgeTransactor{contract: contract}, ERC20BridgeFilterer: ERC20BridgeFilterer{contract: contract}}, nil
}

// ERC20Bridge is an auto generated Go binding around an Ethereum contract.
type ERC20Bridge struct {
	ERC20BridgeCaller     // Read-only binding to the contract
	ERC20BridgeTransactor // Write-only binding to the contract
	ERC20BridgeFilterer   // Log filterer for contract events
}

// ERC20BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20BridgeSession struct {
	Contract     *ERC20Bridge      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20BridgeCallerSession struct {
	Contract *ERC20BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ERC20BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20BridgeTransactorSession struct {
	Contract     *ERC20BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ERC20BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20BridgeRaw struct {
	Contract *ERC20Bridge // Generic contract binding to access the raw methods on
}

// ERC20BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20BridgeCallerRaw struct {
	Contract *ERC20BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20BridgeTransactorRaw struct {
	Contract *ERC20BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Bridge creates a new instance of ERC20Bridge, bound to a specific deployed contract.
func NewERC20Bridge(address common.Address, backend bind.ContractBackend) (*ERC20Bridge, error) {
	contract, err := bindERC20Bridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Bridge{ERC20BridgeCaller: ERC20BridgeCaller{contract: contract}, ERC20BridgeTransactor: ERC20BridgeTransactor{contract: contract}, ERC20BridgeFilterer: ERC20BridgeFilterer{contract: contract}}, nil
}

// NewERC20BridgeCaller creates a new read-only instance of ERC20Bridge, bound to a specific deployed contract.
func NewERC20BridgeCaller(address common.Address, caller bind.ContractCaller) (*ERC20BridgeCaller, error) {
	contract, err := bindERC20Bridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeCaller{contract: contract}, nil
}

// NewERC20BridgeTransactor creates a new write-only instance of ERC20Bridge, bound to a specific deployed contract.
func NewERC20BridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20BridgeTransactor, error) {
	contract, err := bindERC20Bridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeTransactor{contract: contract}, nil
}

// NewERC20BridgeFilterer creates a new log filterer instance of ERC20Bridge, bound to a specific deployed contract.
func NewERC20BridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20BridgeFilterer, error) {
	contract, err := bindERC20Bridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeFilterer{contract: contract}, nil
}

// bindERC20Bridge binds a generic wrapper to an already deployed contract.
func bindERC20Bridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC20BridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Bridge *ERC20BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Bridge.Contract.ERC20BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Bridge *ERC20BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ERC20BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Bridge *ERC20BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ERC20BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Bridge *ERC20BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Bridge *ERC20BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Bridge *ERC20BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.contract.Transact(opts, method, params...)
}

// GetAssetSource is a free data retrieval call binding the contract method 0x786b0bc0.
//
// Solidity: function get_asset_source(bytes32 vega_asset_id) view returns(address)
func (_ERC20Bridge *ERC20BridgeCaller) GetAssetSource(opts *bind.CallOpts, vega_asset_id [32]byte) (common.Address, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_asset_source", vega_asset_id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAssetSource is a free data retrieval call binding the contract method 0x786b0bc0.
//
// Solidity: function get_asset_source(bytes32 vega_asset_id) view returns(address)
func (_ERC20Bridge *ERC20BridgeSession) GetAssetSource(vega_asset_id [32]byte) (common.Address, error) {
	return _ERC20Bridge.Contract.GetAssetSource(&_ERC20Bridge.CallOpts, vega_asset_id)
}

// GetAssetSource is a free data retrieval call binding the contract method 0x786b0bc0.
//
// Solidity: function get_asset_source(bytes32 vega_asset_id) view returns(address)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetAssetSource(vega_asset_id [32]byte) (common.Address, error) {
	return _ERC20Bridge.Contract.GetAssetSource(&_ERC20Bridge.CallOpts, vega_asset_id)
}

// GetDepositMaximum is a free data retrieval call binding the contract method 0x1d501b5d.
//
// Solidity: function get_deposit_maximum(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCaller) GetDepositMaximum(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_deposit_maximum", asset_source)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositMaximum is a free data retrieval call binding the contract method 0x1d501b5d.
//
// Solidity: function get_deposit_maximum(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeSession) GetDepositMaximum(asset_source common.Address) (*big.Int, error) {
	return _ERC20Bridge.Contract.GetDepositMaximum(&_ERC20Bridge.CallOpts, asset_source)
}

// GetDepositMaximum is a free data retrieval call binding the contract method 0x1d501b5d.
//
// Solidity: function get_deposit_maximum(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetDepositMaximum(asset_source common.Address) (*big.Int, error) {
	return _ERC20Bridge.Contract.GetDepositMaximum(&_ERC20Bridge.CallOpts, asset_source)
}

// GetDepositMinimum is a free data retrieval call binding the contract method 0x4322b1f2.
//
// Solidity: function get_deposit_minimum(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCaller) GetDepositMinimum(opts *bind.CallOpts, asset_source common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_deposit_minimum", asset_source)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositMinimum is a free data retrieval call binding the contract method 0x4322b1f2.
//
// Solidity: function get_deposit_minimum(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeSession) GetDepositMinimum(asset_source common.Address) (*big.Int, error) {
	return _ERC20Bridge.Contract.GetDepositMinimum(&_ERC20Bridge.CallOpts, asset_source)
}

// GetDepositMinimum is a free data retrieval call binding the contract method 0x4322b1f2.
//
// Solidity: function get_deposit_minimum(address asset_source) view returns(uint256)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetDepositMinimum(asset_source common.Address) (*big.Int, error) {
	return _ERC20Bridge.Contract.GetDepositMinimum(&_ERC20Bridge.CallOpts, asset_source)
}

// GetMultisigControlAddress is a free data retrieval call binding the contract method 0xc58dc3b9.
//
// Solidity: function get_multisig_control_address() view returns(address)
func (_ERC20Bridge *ERC20BridgeCaller) GetMultisigControlAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_multisig_control_address")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMultisigControlAddress is a free data retrieval call binding the contract method 0xc58dc3b9.
//
// Solidity: function get_multisig_control_address() view returns(address)
func (_ERC20Bridge *ERC20BridgeSession) GetMultisigControlAddress() (common.Address, error) {
	return _ERC20Bridge.Contract.GetMultisigControlAddress(&_ERC20Bridge.CallOpts)
}

// GetMultisigControlAddress is a free data retrieval call binding the contract method 0xc58dc3b9.
//
// Solidity: function get_multisig_control_address() view returns(address)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetMultisigControlAddress() (common.Address, error) {
	return _ERC20Bridge.Contract.GetMultisigControlAddress(&_ERC20Bridge.CallOpts)
}

// GetVegaAssetId is a free data retrieval call binding the contract method 0xa06b5d39.
//
// Solidity: function get_vega_asset_id(address asset_source) view returns(bytes32)
func (_ERC20Bridge *ERC20BridgeCaller) GetVegaAssetId(opts *bind.CallOpts, asset_source common.Address) ([32]byte, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "get_vega_asset_id", asset_source)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetVegaAssetId is a free data retrieval call binding the contract method 0xa06b5d39.
//
// Solidity: function get_vega_asset_id(address asset_source) view returns(bytes32)
func (_ERC20Bridge *ERC20BridgeSession) GetVegaAssetId(asset_source common.Address) ([32]byte, error) {
	return _ERC20Bridge.Contract.GetVegaAssetId(&_ERC20Bridge.CallOpts, asset_source)
}

// GetVegaAssetId is a free data retrieval call binding the contract method 0xa06b5d39.
//
// Solidity: function get_vega_asset_id(address asset_source) view returns(bytes32)
func (_ERC20Bridge *ERC20BridgeCallerSession) GetVegaAssetId(asset_source common.Address) ([32]byte, error) {
	return _ERC20Bridge.Contract.GetVegaAssetId(&_ERC20Bridge.CallOpts, asset_source)
}

// IsAssetListed is a free data retrieval call binding the contract method 0x7fd27b7f.
//
// Solidity: function is_asset_listed(address asset_source) view returns(bool)
func (_ERC20Bridge *ERC20BridgeCaller) IsAssetListed(opts *bind.CallOpts, asset_source common.Address) (bool, error) {
	var out []interface{}
	err := _ERC20Bridge.contract.Call(opts, &out, "is_asset_listed", asset_source)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAssetListed is a free data retrieval call binding the contract method 0x7fd27b7f.
//
// Solidity: function is_asset_listed(address asset_source) view returns(bool)
func (_ERC20Bridge *ERC20BridgeSession) IsAssetListed(asset_source common.Address) (bool, error) {
	return _ERC20Bridge.Contract.IsAssetListed(&_ERC20Bridge.CallOpts, asset_source)
}

// IsAssetListed is a free data retrieval call binding the contract method 0x7fd27b7f.
//
// Solidity: function is_asset_listed(address asset_source) view returns(bool)
func (_ERC20Bridge *ERC20BridgeCallerSession) IsAssetListed(asset_source common.Address) (bool, error) {
	return _ERC20Bridge.Contract.IsAssetListed(&_ERC20Bridge.CallOpts, asset_source)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xf7683932.
//
// Solidity: function deposit_asset(address asset_source, uint256 amount, bytes32 vega_public_key) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) DepositAsset(opts *bind.TransactOpts, asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "deposit_asset", asset_source, amount, vega_public_key)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xf7683932.
//
// Solidity: function deposit_asset(address asset_source, uint256 amount, bytes32 vega_public_key) returns()
func (_ERC20Bridge *ERC20BridgeSession) DepositAsset(asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.DepositAsset(&_ERC20Bridge.TransactOpts, asset_source, amount, vega_public_key)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xf7683932.
//
// Solidity: function deposit_asset(address asset_source, uint256 amount, bytes32 vega_public_key) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) DepositAsset(asset_source common.Address, amount *big.Int, vega_public_key [32]byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.DepositAsset(&_ERC20Bridge.TransactOpts, asset_source, amount, vega_public_key)
}

// ListAsset is a paid mutator transaction binding the contract method 0xa8780cda.
//
// Solidity: function list_asset(address asset_source, bytes32 vega_asset_id, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) ListAsset(opts *bind.TransactOpts, asset_source common.Address, vega_asset_id [32]byte, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "list_asset", asset_source, vega_asset_id, nonce, signatures)
}

// ListAsset is a paid mutator transaction binding the contract method 0xa8780cda.
//
// Solidity: function list_asset(address asset_source, bytes32 vega_asset_id, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) ListAsset(asset_source common.Address, vega_asset_id [32]byte, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ListAsset(&_ERC20Bridge.TransactOpts, asset_source, vega_asset_id, nonce, signatures)
}

// ListAsset is a paid mutator transaction binding the contract method 0xa8780cda.
//
// Solidity: function list_asset(address asset_source, bytes32 vega_asset_id, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) ListAsset(asset_source common.Address, vega_asset_id [32]byte, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.ListAsset(&_ERC20Bridge.TransactOpts, asset_source, vega_asset_id, nonce, signatures)
}

// RemoveAsset is a paid mutator transaction binding the contract method 0xc76de358.
//
// Solidity: function remove_asset(address asset_source, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) RemoveAsset(opts *bind.TransactOpts, asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "remove_asset", asset_source, nonce, signatures)
}

// RemoveAsset is a paid mutator transaction binding the contract method 0xc76de358.
//
// Solidity: function remove_asset(address asset_source, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) RemoveAsset(asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.RemoveAsset(&_ERC20Bridge.TransactOpts, asset_source, nonce, signatures)
}

// RemoveAsset is a paid mutator transaction binding the contract method 0xc76de358.
//
// Solidity: function remove_asset(address asset_source, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) RemoveAsset(asset_source common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.RemoveAsset(&_ERC20Bridge.TransactOpts, asset_source, nonce, signatures)
}

// SetDepositMaximum is a paid mutator transaction binding the contract method 0x292463b1.
//
// Solidity: function set_deposit_maximum(address asset_source, uint256 maximum_amount, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) SetDepositMaximum(opts *bind.TransactOpts, asset_source common.Address, maximum_amount *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "set_deposit_maximum", asset_source, maximum_amount, nonce, signatures)
}

// SetDepositMaximum is a paid mutator transaction binding the contract method 0x292463b1.
//
// Solidity: function set_deposit_maximum(address asset_source, uint256 maximum_amount, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) SetDepositMaximum(asset_source common.Address, maximum_amount *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.SetDepositMaximum(&_ERC20Bridge.TransactOpts, asset_source, maximum_amount, nonce, signatures)
}

// SetDepositMaximum is a paid mutator transaction binding the contract method 0x292463b1.
//
// Solidity: function set_deposit_maximum(address asset_source, uint256 maximum_amount, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) SetDepositMaximum(asset_source common.Address, maximum_amount *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.SetDepositMaximum(&_ERC20Bridge.TransactOpts, asset_source, maximum_amount, nonce, signatures)
}

// SetDepositMinimum is a paid mutator transaction binding the contract method 0x3882b3da.
//
// Solidity: function set_deposit_minimum(address asset_source, uint256 minimum_amount, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) SetDepositMinimum(opts *bind.TransactOpts, asset_source common.Address, minimum_amount *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "set_deposit_minimum", asset_source, minimum_amount, nonce, signatures)
}

// SetDepositMinimum is a paid mutator transaction binding the contract method 0x3882b3da.
//
// Solidity: function set_deposit_minimum(address asset_source, uint256 minimum_amount, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) SetDepositMinimum(asset_source common.Address, minimum_amount *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.SetDepositMinimum(&_ERC20Bridge.TransactOpts, asset_source, minimum_amount, nonce, signatures)
}

// SetDepositMinimum is a paid mutator transaction binding the contract method 0x3882b3da.
//
// Solidity: function set_deposit_minimum(address asset_source, uint256 minimum_amount, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) SetDepositMinimum(asset_source common.Address, minimum_amount *big.Int, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.SetDepositMinimum(&_ERC20Bridge.TransactOpts, asset_source, minimum_amount, nonce, signatures)
}

// WithdrawAsset is a paid mutator transaction binding the contract method 0x927ee5b9.
//
// Solidity: function withdraw_asset(address asset_source, uint256 amount, address target, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactor) WithdrawAsset(opts *bind.TransactOpts, asset_source common.Address, amount *big.Int, target common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.contract.Transact(opts, "withdraw_asset", asset_source, amount, target, nonce, signatures)
}

// WithdrawAsset is a paid mutator transaction binding the contract method 0x927ee5b9.
//
// Solidity: function withdraw_asset(address asset_source, uint256 amount, address target, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeSession) WithdrawAsset(asset_source common.Address, amount *big.Int, target common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.WithdrawAsset(&_ERC20Bridge.TransactOpts, asset_source, amount, target, nonce, signatures)
}

// WithdrawAsset is a paid mutator transaction binding the contract method 0x927ee5b9.
//
// Solidity: function withdraw_asset(address asset_source, uint256 amount, address target, uint256 nonce, bytes signatures) returns()
func (_ERC20Bridge *ERC20BridgeTransactorSession) WithdrawAsset(asset_source common.Address, amount *big.Int, target common.Address, nonce *big.Int, signatures []byte) (*types.Transaction, error) {
	return _ERC20Bridge.Contract.WithdrawAsset(&_ERC20Bridge.TransactOpts, asset_source, amount, target, nonce, signatures)
}

// ERC20BridgeAssetDepositMaximumSetIterator is returned from FilterAssetDepositMaximumSet and is used to iterate over the raw logs and unpacked data for AssetDepositMaximumSet events raised by the ERC20Bridge contract.
type ERC20BridgeAssetDepositMaximumSetIterator struct {
	Event *ERC20BridgeAssetDepositMaximumSet // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetDepositMaximumSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetDepositMaximumSet)
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
		it.Event = new(ERC20BridgeAssetDepositMaximumSet)
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
func (it *ERC20BridgeAssetDepositMaximumSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetDepositMaximumSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetDepositMaximumSet represents a AssetDepositMaximumSet event raised by the ERC20Bridge contract.
type ERC20BridgeAssetDepositMaximumSet struct {
	AssetSource common.Address
	NewMaximum  *big.Int
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetDepositMaximumSet is a free log retrieval operation binding the contract event 0x8861dbda4052c4bf7ba24aad82097ef1e6ade8e6f1b42341640db8a12abfa993.
//
// Solidity: event Asset_Deposit_Maximum_Set(address indexed asset_source, uint256 new_maximum, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetDepositMaximumSet(opts *bind.FilterOpts, asset_source []common.Address) (*ERC20BridgeAssetDepositMaximumSetIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Deposit_Maximum_Set", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetDepositMaximumSetIterator{contract: _ERC20Bridge.contract, event: "Asset_Deposit_Maximum_Set", logs: logs, sub: sub}, nil
}

// WatchAssetDepositMaximumSet is a free log subscription operation binding the contract event 0x8861dbda4052c4bf7ba24aad82097ef1e6ade8e6f1b42341640db8a12abfa993.
//
// Solidity: event Asset_Deposit_Maximum_Set(address indexed asset_source, uint256 new_maximum, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetDepositMaximumSet(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetDepositMaximumSet, asset_source []common.Address) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Deposit_Maximum_Set", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetDepositMaximumSet)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Deposit_Maximum_Set", log); err != nil {
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

// ParseAssetDepositMaximumSet is a log parse operation binding the contract event 0x8861dbda4052c4bf7ba24aad82097ef1e6ade8e6f1b42341640db8a12abfa993.
//
// Solidity: event Asset_Deposit_Maximum_Set(address indexed asset_source, uint256 new_maximum, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetDepositMaximumSet(log types.Log) (*ERC20BridgeAssetDepositMaximumSet, error) {
	event := new(ERC20BridgeAssetDepositMaximumSet)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Deposit_Maximum_Set", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeAssetDepositMinimumSetIterator is returned from FilterAssetDepositMinimumSet and is used to iterate over the raw logs and unpacked data for AssetDepositMinimumSet events raised by the ERC20Bridge contract.
type ERC20BridgeAssetDepositMinimumSetIterator struct {
	Event *ERC20BridgeAssetDepositMinimumSet // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetDepositMinimumSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetDepositMinimumSet)
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
		it.Event = new(ERC20BridgeAssetDepositMinimumSet)
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
func (it *ERC20BridgeAssetDepositMinimumSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetDepositMinimumSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetDepositMinimumSet represents a AssetDepositMinimumSet event raised by the ERC20Bridge contract.
type ERC20BridgeAssetDepositMinimumSet struct {
	AssetSource common.Address
	NewMinimum  *big.Int
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetDepositMinimumSet is a free log retrieval operation binding the contract event 0x4ed0df0b169b573722ecdcc12333646a6efd0445c28fe277470bbfca620e8ad5.
//
// Solidity: event Asset_Deposit_Minimum_Set(address indexed asset_source, uint256 new_minimum, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetDepositMinimumSet(opts *bind.FilterOpts, asset_source []common.Address) (*ERC20BridgeAssetDepositMinimumSetIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Deposit_Minimum_Set", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetDepositMinimumSetIterator{contract: _ERC20Bridge.contract, event: "Asset_Deposit_Minimum_Set", logs: logs, sub: sub}, nil
}

// WatchAssetDepositMinimumSet is a free log subscription operation binding the contract event 0x4ed0df0b169b573722ecdcc12333646a6efd0445c28fe277470bbfca620e8ad5.
//
// Solidity: event Asset_Deposit_Minimum_Set(address indexed asset_source, uint256 new_minimum, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetDepositMinimumSet(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetDepositMinimumSet, asset_source []common.Address) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Deposit_Minimum_Set", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetDepositMinimumSet)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Deposit_Minimum_Set", log); err != nil {
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

// ParseAssetDepositMinimumSet is a log parse operation binding the contract event 0x4ed0df0b169b573722ecdcc12333646a6efd0445c28fe277470bbfca620e8ad5.
//
// Solidity: event Asset_Deposit_Minimum_Set(address indexed asset_source, uint256 new_minimum, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetDepositMinimumSet(log types.Log) (*ERC20BridgeAssetDepositMinimumSet, error) {
	event := new(ERC20BridgeAssetDepositMinimumSet)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Deposit_Minimum_Set", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeAssetDepositedIterator is returned from FilterAssetDeposited and is used to iterate over the raw logs and unpacked data for AssetDeposited events raised by the ERC20Bridge contract.
type ERC20BridgeAssetDepositedIterator struct {
	Event *ERC20BridgeAssetDeposited // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetDeposited)
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
		it.Event = new(ERC20BridgeAssetDeposited)
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
func (it *ERC20BridgeAssetDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetDeposited represents a AssetDeposited event raised by the ERC20Bridge contract.
type ERC20BridgeAssetDeposited struct {
	UserAddress   common.Address
	AssetSource   common.Address
	Amount        *big.Int
	VegaPublicKey [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAssetDeposited is a free log retrieval operation binding the contract event 0x3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf.
//
// Solidity: event Asset_Deposited(address indexed user_address, address indexed asset_source, uint256 amount, bytes32 vega_public_key)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetDeposited(opts *bind.FilterOpts, user_address []common.Address, asset_source []common.Address) (*ERC20BridgeAssetDepositedIterator, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Deposited", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetDepositedIterator{contract: _ERC20Bridge.contract, event: "Asset_Deposited", logs: logs, sub: sub}, nil
}

// WatchAssetDeposited is a free log subscription operation binding the contract event 0x3724ff5e82ddc640a08d68b0b782a5991aea0de51a8dd10a59cdbe5b3ec4e6bf.
//
// Solidity: event Asset_Deposited(address indexed user_address, address indexed asset_source, uint256 amount, bytes32 vega_public_key)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetDeposited(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetDeposited, user_address []common.Address, asset_source []common.Address) (event.Subscription, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Deposited", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetDeposited)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Deposited", log); err != nil {
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
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetDeposited(log types.Log) (*ERC20BridgeAssetDeposited, error) {
	event := new(ERC20BridgeAssetDeposited)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeAssetListedIterator is returned from FilterAssetListed and is used to iterate over the raw logs and unpacked data for AssetListed events raised by the ERC20Bridge contract.
type ERC20BridgeAssetListedIterator struct {
	Event *ERC20BridgeAssetListed // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetListedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetListed)
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
		it.Event = new(ERC20BridgeAssetListed)
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
func (it *ERC20BridgeAssetListedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetListedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetListed represents a AssetListed event raised by the ERC20Bridge contract.
type ERC20BridgeAssetListed struct {
	AssetSource common.Address
	VegaAssetId [32]byte
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetListed is a free log retrieval operation binding the contract event 0x4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e.
//
// Solidity: event Asset_Listed(address indexed asset_source, bytes32 indexed vega_asset_id, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetListed(opts *bind.FilterOpts, asset_source []common.Address, vega_asset_id [][32]byte) (*ERC20BridgeAssetListedIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}
	var vega_asset_idRule []interface{}
	for _, vega_asset_idItem := range vega_asset_id {
		vega_asset_idRule = append(vega_asset_idRule, vega_asset_idItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Listed", asset_sourceRule, vega_asset_idRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetListedIterator{contract: _ERC20Bridge.contract, event: "Asset_Listed", logs: logs, sub: sub}, nil
}

// WatchAssetListed is a free log subscription operation binding the contract event 0x4180d77d05ff0d31650c548c23f2de07a3da3ad42e3dd6edd817b438a150452e.
//
// Solidity: event Asset_Listed(address indexed asset_source, bytes32 indexed vega_asset_id, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetListed(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetListed, asset_source []common.Address, vega_asset_id [][32]byte) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}
	var vega_asset_idRule []interface{}
	for _, vega_asset_idItem := range vega_asset_id {
		vega_asset_idRule = append(vega_asset_idRule, vega_asset_idItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Listed", asset_sourceRule, vega_asset_idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetListed)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Listed", log); err != nil {
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
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetListed(log types.Log) (*ERC20BridgeAssetListed, error) {
	event := new(ERC20BridgeAssetListed)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Listed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeAssetRemovedIterator is returned from FilterAssetRemoved and is used to iterate over the raw logs and unpacked data for AssetRemoved events raised by the ERC20Bridge contract.
type ERC20BridgeAssetRemovedIterator struct {
	Event *ERC20BridgeAssetRemoved // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetRemoved)
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
		it.Event = new(ERC20BridgeAssetRemoved)
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
func (it *ERC20BridgeAssetRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetRemoved represents a AssetRemoved event raised by the ERC20Bridge contract.
type ERC20BridgeAssetRemoved struct {
	AssetSource common.Address
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetRemoved is a free log retrieval operation binding the contract event 0x58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c9.
//
// Solidity: event Asset_Removed(address indexed asset_source, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetRemoved(opts *bind.FilterOpts, asset_source []common.Address) (*ERC20BridgeAssetRemovedIterator, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Removed", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetRemovedIterator{contract: _ERC20Bridge.contract, event: "Asset_Removed", logs: logs, sub: sub}, nil
}

// WatchAssetRemoved is a free log subscription operation binding the contract event 0x58ad5e799e2df93ab408be0e5c1870d44c80b5bca99dfaf7ddf0dab5e6b155c9.
//
// Solidity: event Asset_Removed(address indexed asset_source, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetRemoved(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetRemoved, asset_source []common.Address) (event.Subscription, error) {

	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Removed", asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetRemoved)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Removed", log); err != nil {
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
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetRemoved(log types.Log) (*ERC20BridgeAssetRemoved, error) {
	event := new(ERC20BridgeAssetRemoved)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Removed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20BridgeAssetWithdrawnIterator is returned from FilterAssetWithdrawn and is used to iterate over the raw logs and unpacked data for AssetWithdrawn events raised by the ERC20Bridge contract.
type ERC20BridgeAssetWithdrawnIterator struct {
	Event *ERC20BridgeAssetWithdrawn // Event containing the contract specifics and raw log

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
func (it *ERC20BridgeAssetWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BridgeAssetWithdrawn)
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
		it.Event = new(ERC20BridgeAssetWithdrawn)
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
func (it *ERC20BridgeAssetWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BridgeAssetWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BridgeAssetWithdrawn represents a AssetWithdrawn event raised by the ERC20Bridge contract.
type ERC20BridgeAssetWithdrawn struct {
	UserAddress common.Address
	AssetSource common.Address
	Amount      *big.Int
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAssetWithdrawn is a free log retrieval operation binding the contract event 0xa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df.
//
// Solidity: event Asset_Withdrawn(address indexed user_address, address indexed asset_source, uint256 amount, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) FilterAssetWithdrawn(opts *bind.FilterOpts, user_address []common.Address, asset_source []common.Address) (*ERC20BridgeAssetWithdrawnIterator, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.FilterLogs(opts, "Asset_Withdrawn", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BridgeAssetWithdrawnIterator{contract: _ERC20Bridge.contract, event: "Asset_Withdrawn", logs: logs, sub: sub}, nil
}

// WatchAssetWithdrawn is a free log subscription operation binding the contract event 0xa79be4f3361e32d396d64c478ecef73732cb40b2a75702c3b3b3226a2c83b5df.
//
// Solidity: event Asset_Withdrawn(address indexed user_address, address indexed asset_source, uint256 amount, uint256 nonce)
func (_ERC20Bridge *ERC20BridgeFilterer) WatchAssetWithdrawn(opts *bind.WatchOpts, sink chan<- *ERC20BridgeAssetWithdrawn, user_address []common.Address, asset_source []common.Address) (event.Subscription, error) {

	var user_addressRule []interface{}
	for _, user_addressItem := range user_address {
		user_addressRule = append(user_addressRule, user_addressItem)
	}
	var asset_sourceRule []interface{}
	for _, asset_sourceItem := range asset_source {
		asset_sourceRule = append(asset_sourceRule, asset_sourceItem)
	}

	logs, sub, err := _ERC20Bridge.contract.WatchLogs(opts, "Asset_Withdrawn", user_addressRule, asset_sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BridgeAssetWithdrawn)
				if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Withdrawn", log); err != nil {
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
func (_ERC20Bridge *ERC20BridgeFilterer) ParseAssetWithdrawn(log types.Log) (*ERC20BridgeAssetWithdrawn, error) {
	event := new(ERC20BridgeAssetWithdrawn)
	if err := _ERC20Bridge.contract.UnpackLog(event, "Asset_Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
