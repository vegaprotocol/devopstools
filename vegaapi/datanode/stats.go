package datanode

import (
	"context"
	"fmt"

	e "github.com/vegaprotocol/devopstools/errors"
	"github.com/vegaprotocol/devopstools/types"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"

	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetAllNetworkParameters() (*types.NetworkParams, error) {
	res, err := n.ListNetworkParameters(&dataapipb.ListNetworkParametersRequest{})
	if err != nil {
		return nil, err
	}
	networkParams := map[string]string{}
	for _, edge := range res.NetworkParameters.Edges {
		networkParams[edge.Node.Key] = edge.Node.Value
	}
	return types.NewNetworkParams(networkParams), nil
}

func (n *DataNode) GetCurrentEpoch() (*vega.Epoch, error) {
	res, err := n.getEpoch(&dataapipb.GetEpochRequest{})
	if err != nil {
		return nil, err
	}
	return res.Epoch, nil
}

func (n *DataNode) ListAssets(ctx context.Context) (map[string]*vega.AssetDetails, error) {
	res, err := n.listAssets(ctx, &dataapipb.ListAssetsRequest{})
	if err != nil {
		return nil, err
	}

	assets := map[string]*vega.AssetDetails{}
	for _, edge := range res.Assets.Edges {
		assets[edge.Node.Id] = edge.Node.Details
	}
	return assets, nil
}

func (n *DataNode) ListNetworkParameters(req *dataapipb.ListNetworkParametersRequest) (response *dataapipb.ListNetworkParametersResponse, err error) {
	msg := "gRPC call failed (data-node): ListNetworkParameters: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err = c.ListNetworkParameters(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *DataNode) getEpoch(req *dataapipb.GetEpochRequest) (response *dataapipb.GetEpochResponse, err error) {
	msg := "gRPC call failed (data-node): GetEpoch: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err = c.GetEpoch(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *DataNode) listAssets(ctx context.Context, req *dataapipb.ListAssetsRequest) (*dataapipb.ListAssetsResponse, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	response, err := c.ListAssets(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	return response, nil
}
