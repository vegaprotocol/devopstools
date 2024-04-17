package veganetworksmartcontracts

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/smartcontracts/erc20assetpool"
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20bridge"
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20token"
	"github.com/vegaprotocol/devopstools/smartcontracts/multisigcontrol"
	"github.com/vegaprotocol/devopstools/smartcontracts/stakingbridge"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type VegaNetworkSmartContracts struct {
	VegaToken       *erc20token.ERC20Token
	AssetPool       *erc20assetpool.ERC20AssetPool
	ERC20Bridge     *erc20bridge.ERC20Bridge
	MultisigControl *multisigcontrol.MultisigControl
	StakingBridge   *stakingbridge.StakingBridge

	EthClient *ethclient.Client
	logger    *zap.Logger
}

func NewVegaNetworkSmartContracts(
	ethClient *ethclient.Client,
	vegaTokenHexAddress string,
	assetPoolHexAddress string,
	erc20BridgeHexAddress string,
	multisigControlHexAddress string,
	stakingBridgeHexAddress string,
	logger *zap.Logger,
) (*VegaNetworkSmartContracts, error) {
	var (
		result = &VegaNetworkSmartContracts{
			EthClient: ethClient,
			logger:    logger,
		}
		errMsg = "failed to create new VegaNetwork SmartContracts, %w"
		err    error
	)

	// ERC20 Bridge
	if len(erc20BridgeHexAddress) == 0 {
		return nil, fmt.Errorf("missing ERC20 Bridge address")
	}
	result.ERC20Bridge, err = erc20bridge.NewERC20Bridge(result.EthClient, erc20BridgeHexAddress, erc20bridge.V2)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	// Staking Bridge
	if len(stakingBridgeHexAddress) == 0 {
		result.StakingBridge, err = stakingbridge.NewStakingBridge(result.EthClient, stakingBridgeHexAddress, stakingbridge.V1)
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
	}

	// Multisig Control
	if len(multisigControlHexAddress) == 0 {
		multisigControlAddress, err := result.ERC20Bridge.GetMultisigControlAddress(&bind.CallOpts{})
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		multisigControlHexAddress = multisigControlAddress.Hex()
	}
	result.MultisigControl, err = multisigcontrol.NewMultisigControl(result.EthClient, multisigControlHexAddress, multisigcontrol.V2, nil)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	// Asset Pool
	if len(assetPoolHexAddress) == 0 {
		assetPoolAddress, err := result.ERC20Bridge.Erc20AssetPoolAddress(&bind.CallOpts{})
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		assetPoolHexAddress = assetPoolAddress.Hex()
	}
	result.AssetPool, err = erc20assetpool.NewERC20AssetPool(result.EthClient, assetPoolHexAddress, erc20assetpool.V1)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}

	// Vega Token
	if len(vegaTokenHexAddress) == 0 {
		vegaTokenAddress, err := result.StakingBridge.StakingToken(&bind.CallOpts{})
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		vegaTokenHexAddress = vegaTokenAddress.Hex()
	}
	result.VegaToken, err = erc20token.NewERC20Token(result.EthClient, vegaTokenHexAddress, erc20token.Base)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	return result, nil
}
