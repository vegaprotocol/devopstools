package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"code.vegaprotocol.io/vega/core/netparams"
	"code.vegaprotocol.io/vega/protos/vega"

	"google.golang.org/protobuf/encoding/protojson"
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

func (p *NetworkParams) PrimaryEthereumConfig() (*vega.EthereumConfig, error) {
	param := netparams.BlockchainsPrimaryEthereumConfig
	val, ok := p.Params[param]
	if !ok {
		return nil, fmt.Errorf("missing network parameter %q", param)
	}
	result := &vega.EthereumConfig{}
	if err := json.Unmarshal([]byte(val), result); err != nil {
		return nil, fmt.Errorf("could not deserialize network parameter %q: %w", param, err)
	}

	return result, nil
}

func (p *NetworkParams) EVMChainConfig() (*vega.EVMBridgeConfigs, error) {
	param := netparams.BlockchainsEVMBridgeConfigs
	val, ok := p.Params[param]
	if !ok {
		return nil, fmt.Errorf("missing network parameter %q", param)
	}
	result := &vega.EVMBridgeConfigs{}
	if err := json.Unmarshal([]byte(val), result); err != nil {
		return nil, fmt.Errorf("could not deserialize network parameter %q: %w", param, err)
	}

	return result, nil
}

func (p *NetworkParams) GetEthereumL2Configs() (*vega.EthereumL2Configs, error) {
	val, ok := p.Params[netparams.BlockchainsEthereumL2Configs]
	if !ok {
		return nil, fmt.Errorf("the %s network parameter is missing", netparams.BlockchainsEthereumL2Configs)
	}

	result := &vega.EthereumL2Configs{}
	if err := protojson.Unmarshal([]byte(val), result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the %s parameter into go structure: %w", netparams.BlockchainsEthereumL2Configs, err)
	}

	return result, nil
}

func (p *NetworkParams) GetMinimumEthereumEventsForNewValidator() (int, error) {
	param := netparams.MinimumEthereumEventsForNewValidator
	val, ok := p.Params[param]
	if !ok {
		return -1, fmt.Errorf("failed to get MinimumEthereumEventsForNewValidator, missing '%s' network parameter", param)
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return -1, fmt.Errorf("failed to get MinimumEthereumEventsForNewValidator, value '%s' is not an integer", val)
	}
	return intVal, nil
}

func (p *NetworkParams) GetMinFundsForApplyReferral() int64 {
	val, found := p.Params[netparams.SpamProtectionApplyReferralMinFunds]
	if !found {
		return 0
	}

	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}

	return intVal
}
