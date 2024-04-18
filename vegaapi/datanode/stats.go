package datanode

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	e "github.com/vegaprotocol/devopstools/errors"
	"github.com/vegaprotocol/devopstools/types"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"

	"google.golang.org/grpc/connectivity"
)

var ErrAssetNotFound = errors.New("asset not found")

func (n *DataNode) GeneralAccountBalance(ctx context.Context, partyID, assetID string) (*big.Int, error) {
	whaleAccounts, err := n.ListAccounts(ctx, partyID, vega.AccountType_ACCOUNT_TYPE_GENERAL, &assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve general accounts: %w", err)
	}

	whaleFundsAsSubUnits := big.NewInt(0)
	if len(whaleAccounts) > 0 {
		whaleFundsAsSubUnits = whaleAccounts[0].Balance
	}

	return whaleFundsAsSubUnits, nil
}

func (n *DataNode) ListNetworkParameters(ctx context.Context) (*types.NetworkParams, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	response, err := c.ListNetworkParameters(reqCtx, &dataapipb.ListNetworkParametersRequest{})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	networkParams := map[string]string{}
	for _, edge := range response.NetworkParameters.Edges {
		networkParams[edge.Node.Key] = edge.Node.Value
	}
	return types.NewNetworkParams(networkParams), nil
}

func (n *DataNode) GetCurrentEpoch(ctx context.Context) (*vega.Epoch, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	response, err := c.GetEpoch(reqCtx, &dataapipb.GetEpochRequest{})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	return response.Epoch, nil
}

func (n *DataNode) ListAssets(ctx context.Context) (map[string]*vega.AssetDetails, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	response, err := c.ListAssets(ctx, &dataapipb.ListAssetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	assets := map[string]*vega.AssetDetails{}
	for _, edge := range response.Assets.Edges {
		assets[edge.Node.Id] = edge.Node.Details
	}
	return assets, nil
}

func (n *DataNode) ERC20AssetWithAddress(ctx context.Context, address string) (*vega.AssetDetails, *vega.ERC20, error) {
	assets, err := n.ListAssets(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, asset := range assets {
		erc20Token := asset.GetErc20()
		if erc20Token == nil {
			continue
		}
		if erc20Token.ContractAddress == address {
			return asset, erc20Token, nil
		}
	}

	return nil, nil, ErrAssetNotFound
}
