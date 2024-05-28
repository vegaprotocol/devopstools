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

var ActiveMarkets = []vega.Market_State{
	vega.Market_STATE_ACTIVE,
	vega.Market_STATE_SUSPENDED,
}

func (n *DataNode) GetAllMarkets(ctx context.Context) ([]*vega.Market, error) {
	res, err := n.ListMarkets(ctx, &dataapipb.ListMarketsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all markets, %w", err)
	}
	result := make([]*vega.Market, len(res.Markets.Edges))
	for i, edge := range res.Markets.Edges {
		result[i] = edge.Node
	}
	return result, nil
}

func (n *DataNode) GetMarket(req *dataapipb.GetMarketRequest) (response *dataapipb.GetMarketResponse, err error) {
	msg := "gRPC call failed (data-node): GetMarket: %w"
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

	response, err = c.GetMarket(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

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

func (n *DataNode) GetMarketById(marketId string) (*vega.Market, error) {
	if marketId == "" {
		return nil, fmt.Errorf("market id cannot be empty")
	}

	marketResponse, err := n.GetMarket(&dataapipb.GetMarketRequest{
		MarketId: marketId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get market: %w", err)
	}

	return marketResponse.GetMarket(), nil
}

func (n *DataNode) ListMarkets(ctx context.Context, req *dataapipb.ListMarketsRequest) (*dataapipb.ListMarketsResponse, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.CallTimeout)
	defer cancelRequest()

	response, err := c.ListMarkets(reqCtx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	return response, nil
}
