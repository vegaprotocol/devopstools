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
	vega.Market_STATE_PENDING,
}

func (n *DataNode) GetAllMarketsWithState(ctx context.Context, states []vega.Market_State) ([]*vega.Market, error) {
	res, err := n.ListMarkets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all markets, %w", err)
	}
	result := []*vega.Market{}

	for _, edge := range res {
		if !slices.Contains(states, edge.State) {
			continue
		}

		result = append(result, edge)
	}
	return result, nil
}

func (n *DataNode) ListMarkets(ctx context.Context) ([]*vega.Market, error) {
	if n.Client.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Client.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.Client.CallTimeout)
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

func (n *DataNode) GetLatestMarketData(ctx context.Context, marketId string) (*vega.MarketData, error) {
	if n.Client.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Client.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.Client.CallTimeout)
	defer cancelRequest()

	response, err := c.GetLatestMarketData(reqCtx, &dataapipb.GetLatestMarketDataRequest{
		MarketId: marketId,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get latest market data: %w", err)
	}

	return response.MarketData, nil
}

func AssetsIdsByMarket(market *vega.Market) []string {
	instrument := market.TradableInstrument.Instrument
	if instrument.GetFuture() != nil {
		return []string{instrument.GetFuture().SettlementAsset}
	}

	if instrument.GetPerpetual() != nil {
		return []string{instrument.GetPerpetual().SettlementAsset}
	}

	if instrument.GetSpot() != nil {
		return []string{instrument.GetSpot().BaseAsset, instrument.GetSpot().QuoteAsset}
	}

	return []string{}
}
