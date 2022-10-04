package vegaapi

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	e "github.com/vegaprotocol/devopstools/errors"
	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetAllNetworkParameters() (map[string]string, error) {
	res, err := n.ListNetworkParameters(&dataapipb.ListNetworkParametersRequest{})
	if err != nil {
		return nil, err
	}
	networkParams := map[string]string{}
	for _, edge := range res.NetworkParameters.Edges {
		networkParams[edge.Node.Key] = edge.Node.Value
	}
	return networkParams, nil
}

func (n *DataNode) GetCurrentEpoch() (*vega.Epoch, error) {
	res, err := n.GetEpoch(&dataapipb.GetEpochRequest{})
	if err != nil {
		return nil, err
	}
	return res.Epoch, nil
}

func (n *DataNode) GetAssets() (map[string]*vega.AssetDetails, error) {
	res, err := n.ListAssets(&dataapipb.ListAssetsRequest{})
	if err != nil {
		return nil, err
	}
	assetList := map[string]*vega.AssetDetails{}
	for _, edge := range res.Assets.Edges {
		assetList[edge.Node.Id] = edge.Node.Details
	}
	return assetList, nil
}

func (n *DataNode) ListNetworkParameters(req *dataapipb.ListNetworkParametersRequest) (response *dataapipb.ListNetworkParametersResponse, err error) {
	msg := "gRPC call failed (data-node): ListNetworkParameters: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
	defer cancel()

	response, err = c.ListNetworkParameters(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *DataNode) GetNetworkParameter(req *dataapipb.GetNetworkParameterRequest) (response *dataapipb.GetNetworkParameterResponse, err error) {
	msg := "gRPC call failed (data-node): GetNetworkParameter: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
	defer cancel()

	response, err = c.GetNetworkParameter(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *DataNode) GetEpoch(req *dataapipb.GetEpochRequest) (response *dataapipb.GetEpochResponse, err error) {
	msg := "gRPC call failed (data-node): GetEpoch: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
	defer cancel()

	response, err = c.GetEpoch(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *DataNode) ListAssets(req *dataapipb.ListAssetsRequest) (response *dataapipb.ListAssetsResponse, err error) {
	msg := "gRPC call failed (data-node): ListAssets: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
	defer cancel()

	response, err = c.ListAssets(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *DataNode) ListGovernanceData(req *dataapipb.ListGovernanceDataRequest) (response *dataapipb.ListGovernanceDataResponse, err error) {
	msg := "gRPC call failed (data-node): ListGovernanceData: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
	defer cancel()

	response, err = c.ListGovernanceData(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}
