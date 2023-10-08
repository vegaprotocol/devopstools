package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	e "github.com/vegaprotocol/devopstools/errors"
	"google.golang.org/grpc/connectivity"
)

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

func (n *DataNode) ListAssets(req *dataapipb.ListAssetsRequest) (response *dataapipb.ListAssetsResponse, err error) {
	msg := "gRPC call failed (data-node): ListAssets: %w"
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

	response, err = c.ListAssets(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *DataNode) GetAssetById(assetId string) (*vega.Asset, error) {
	res, err := n.GetAsset(&dataapipb.GetAssetRequest{
		AssetId: assetId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get asset for id %s, %w", assetId, err)
	}
	return res.Asset, nil
}

func (n *DataNode) GetAsset(req *dataapipb.GetAssetRequest) (response *dataapipb.GetAssetResponse, err error) {
	msg := "gRPC call failed (data-node): GetAsset: %w"
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

	response, err = c.GetAsset(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}
