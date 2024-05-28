package datanode

import (
	"context"
	"fmt"
	"slices"

	e "github.com/vegaprotocol/devopstools/errors"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"

	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetAllMarketsWithState(ctx context.Context, states []vega.Market_State) ([]*vega.Market, error) {
	res, err := n.ListMarkets(ctx, &dataapipb.ListMarketsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all markets, %w", err)
	}
	result := []*vega.Market{}

	for _, edge := range res.Markets.Edges {
		if edge.Node == nil {
			continue
		}
		if !slices.Contains(states, edge.Node.State) {
			continue
		}

		result = append(result, edge.Node)
	}
	return result, nil
}

func (n *DataNode) ListMarkets(ctx context.Context) ([]*vega.Market, error) {

	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.CallTimeout)
	defer cancelRequest()

	response, err := c.ListMarkets(reqCtx, &dataapipb.ListMarketsRequest{})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	result := make([]*vega.Market, len(response.Markets.Edges))
	for i, edge := range response.Markets.Edges {
		result[i] = edge.Node
	}
	return result, nil
}
