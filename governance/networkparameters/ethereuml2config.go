package networkparameters

import (
	"fmt"

	"code.vegaprotocol.io/vega/protos/vega"
)

func CloneEthereumL2Config(currentParam *vega.EthereumL2Configs) *vega.EthereumL2Configs {
	result := &vega.EthereumL2Configs{
		Configs: []*vega.EthereumL2Config{},
	}

	for _, conf := range currentParam.Configs {
		result.Configs = append(result.Configs, &vega.EthereumL2Config{
			NetworkId:     conf.NetworkId,
			ChainId:       conf.ChainId,
			Name:          conf.Name,
			Confirmations: conf.Confirmations,
			BlockInterval: conf.BlockInterval,
		})
	}

	return result
}

func AppendEthereumL2Config(
	currentParam *vega.EthereumL2Configs,
	newConfig *vega.EthereumL2Config,
	updateExisting bool,
) (*vega.EthereumL2Configs, error) {
	if newConfig.ChainId == "" || newConfig.NetworkId == "" || newConfig.Name == "" {
		return nil, fmt.Errorf("fields ChainId, NetworkId and Name must not be empty")
	}

	result := &vega.EthereumL2Configs{
		Configs: []*vega.EthereumL2Config{},
	}

	found := false
	for _, conf := range currentParam.Configs {
		// Update existing param
		if newConfig.ChainId == conf.ChainId && newConfig.NetworkId == conf.NetworkId {
			found = true
			// update not allowed
			if !updateExisting {
				continue
			}

			result.Configs = append(result.Configs, &vega.EthereumL2Config{
				NetworkId:     newConfig.NetworkId,
				ChainId:       newConfig.ChainId,
				Name:          newConfig.Name,
				Confirmations: newConfig.Confirmations,
				BlockInterval: newConfig.BlockInterval,
			})
			continue
		}

		result.Configs = append(result.Configs, &vega.EthereumL2Config{
			NetworkId:     conf.NetworkId,
			ChainId:       conf.ChainId,
			Name:          conf.Name,
			Confirmations: conf.Confirmations,
			BlockInterval: conf.BlockInterval,
		})
	}

	if !found {
		result.Configs = append(result.Configs, &vega.EthereumL2Config{
			NetworkId:     newConfig.NetworkId,
			ChainId:       newConfig.ChainId,
			Name:          newConfig.Name,
			Confirmations: newConfig.Confirmations,
			BlockInterval: newConfig.BlockInterval,
		})
	}

	return result, nil
}
