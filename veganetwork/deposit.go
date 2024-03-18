package veganetwork

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/vegaprotocol/devopstools/networktools"

	"code.vegaprotocol.io/vega/protos/vega"

	"go.uber.org/zap"
)

type DepositConfig struct {
	AssetId string
	Parties map[string]*big.Float
}

type DepositResult struct {
	AssetId        string
	SuccessParties []string
	FailedParties  map[string]error
	Err            error
}

func DepositAssetsToParties(
	network *VegaNetwork,
	networktools *networktools.NetworkTools,
	config []DepositConfig,
	logger *zap.Logger,
) []DepositResult {
	resultsChannel := make(chan DepositResult, len(config))
	var wg sync.WaitGroup

	for _, singleConfig := range config {
		wg.Add(1)
		go func(singleConfig DepositConfig) {
			defer wg.Done()
			resultsChannel <- DepositAssetToParties(
				network, networktools, singleConfig, logger,
			)
		}(singleConfig)
	}

	wg.Wait()
	close(resultsChannel)

	result := []DepositResult{}

	for singleResult := range resultsChannel {
		result = append(result, singleResult)
	}

	return result
}

func DepositAssetToParties(
	network *VegaNetwork,
	networktools *networktools.NetworkTools,
	config DepositConfig,
	logger *zap.Logger,
) DepositResult {
	assets, err := network.DataNodeClient.GetAssets()
	if err != nil {
		return DepositResult{
			AssetId: config.AssetId,
			Err:     fmt.Errorf("failed to get assets, %w", err),
		}
	}
	asset, ok := assets[config.AssetId]
	if !ok {
		return DepositResult{
			AssetId: config.AssetId,
			Err:     fmt.Errorf("there is no asset with id: %s", config.AssetId),
		}
	}

	if fakeAsset := asset.GetBuiltinAsset(); fakeAsset != nil {
		return DepositFakeAssetToParties(networktools, fakeAsset, config, logger)
	} else if erc20asset := asset.GetErc20(); erc20asset != nil {
		return DepositERC20TokenToParties(
			network, erc20asset.ContractAddress, config, logger,
		)
	} else {
		return DepositResult{
			AssetId: config.AssetId,
			Err:     fmt.Errorf("unsupported asset type %+v", asset.GetSource()),
		}
	}
}

func DepositERC20TokenToParties(
	_ *VegaNetwork,
	_ string,
	_ DepositConfig,
	_ *zap.Logger,
) DepositResult {
	return DepositResult{}
}

//
// FAKE ASSETS
//

func DepositFakeAssetToParties(
	_ *networktools.NetworkTools,
	_ *vega.BuiltinAsset,
	_ DepositConfig,
	_ *zap.Logger,
) DepositResult {
	return DepositResult{}
}
