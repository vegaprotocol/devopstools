package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"code.vegaprotocol.io/vega/core/netparams"
	"code.vegaprotocol.io/vega/protos/vega"
)

type NetworkParams struct {
	Params map[string]string
}

func NewNetworkParams(params map[string]string) *NetworkParams {
	return &NetworkParams{
		Params: params,
	}
}

func (p *NetworkParams) GetMinimumValidatorStake() (*big.Int, error) {
	param := netparams.StakingAndDelegationRewardMinimumValidatorStake
	val, ok := p.Params[param]
	if !ok {
		return nil, fmt.Errorf("failed to get MinimumValidatorStake, missing '%s' network parameter", param)
	}

	minimumValidatorStake := new(big.Int)
	minimumValidatorStake, ok = minimumValidatorStake.SetString(val, 0)
	if !ok {
		return nil, fmt.Errorf("failed to get MinimumValidatorStake, failed to conver '%s'='%s' to big.Int", param, val)
	}
	return minimumValidatorStake, nil
}

func (p *NetworkParams) GetEthereumConfig() (*vega.EthereumConfig, error) {
	param := netparams.BlockchainsEthereumConfig
	val, ok := p.Params[param]
	if !ok {
		return nil, fmt.Errorf("failed to get EthereumConfig, missing '%s' network parameter", param)
	}
	result := &vega.EthereumConfig{}
	if err := json.Unmarshal([]byte(val), result); err != nil {
		return nil, fmt.Errorf("failed to get EthereumConfig, failed to parse %v, %w", result, err)
	}

	return result, nil
}
