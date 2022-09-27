package veganetwork

import (
	"encoding/json"
	"fmt"
	"math/big"

	"code.vegaprotocol.io/vega/protos/vega"
	"github.com/vegaprotocol/devopstools/ethutils"
)

type NetworkParams struct {
	Params map[string]string
}

func NewNetworkParams(params map[string]string) *NetworkParams {
	return &NetworkParams{
		Params: params,
	}
}

func (p *NetworkParams) GetMinimumValidatorStake() (*big.Float, error) {
	val, ok := p.Params["reward.staking.delegation.minimumValidatorStake"]
	if !ok {
		return nil, fmt.Errorf("failed to get MinimumValidatorStake, missing 'reward.staking.delegation.minimumValidatorStake' network parameter")
	}

	minimumValidatorStake := new(big.Int)
	minimumValidatorStake, ok = minimumValidatorStake.SetString(val, 0)
	if !ok {
		return nil, fmt.Errorf("failed to get MinimumValidatorStake, failed to conver 'reward.staking.delegation.minimumValidatorStake'='%s' to big.Int", val)
	}
	humanMinimumValidatorStake := ethutils.VegaTokenToFullTokens(minimumValidatorStake)
	return humanMinimumValidatorStake, nil
}

func (p *NetworkParams) GetEthereumConfig() (*vega.EthereumConfig, error) {
	val, ok := p.Params["blockchains.ethereumConfig"]
	if !ok {
		return nil, fmt.Errorf("failed to get EthereumConfig, missing 'blockchains.ethereumConfig' network parameter")
	}
	result := &vega.EthereumConfig{}
	if err := json.Unmarshal([]byte(val), result); err != nil {
		return nil, fmt.Errorf("failed to get EthereumConfig, failed to parse %v, %w", result, err)
	}

	return result, nil
}
