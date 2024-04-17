package vega

import (
	"context"
	"fmt"

	"github.com/vegaprotocol/devopstools/governance"
	"github.com/vegaprotocol/devopstools/types"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"code.vegaprotocol.io/vega/wallet/wallet"

	"go.uber.org/zap"
)

func UpdateNetworkParameters(ctx context.Context, whaleWallet wallet.Wallet, whalePublicKey string, datanodeClient *datanode.DataNode, updateParams map[string]string, logger *zap.Logger) (*types.NetworkParams, error) {
	networkParameters, err := datanodeClient.GetAllNetworkParameters()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve network parameters from datanode: %w", err)
	}

	updateCount, err := governance.ProposeAndVoteOnNetworkParameters(ctx, updateParams, whaleWallet, whalePublicKey, networkParameters, datanodeClient, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to propose and vote for network parameter update proposals: %w", err)
	}

	if updateCount == 0 {
		logger.Debug("No network parameter update is required before issuing transfers")
		return nil, nil
	}

	updatedNetworkParameters, err := datanodeClient.GetAllNetworkParameters()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve updated network parameters from datanode: %w", err)
	}

	for name, expectedValue := range updateParams {
		updatedValue := updatedNetworkParameters.Params[name]
		if updatedValue != expectedValue {
			return nil, fmt.Errorf("failed to update network parameter %q, current value: %q, expected value: %q", name, updatedValue, expectedValue)
		}
	}

	return updatedNetworkParameters, nil
}
