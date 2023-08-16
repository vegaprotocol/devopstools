package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	e "github.com/vegaprotocol/devopstools/errors"
	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetAllMarkets() ([]*vega.Market, error) {
	res, err := n.ListMarkets(&dataapipb.ListMarketsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all markets, %w", err)
	}
	result := make([]*vega.Market, len(res.Markets.Edges))
	for i, edge := range res.Markets.Edges {
		result[i] = edge.Node
	}
	return result, nil
}

// === TradingDataService ===

// // PartyAccounts returns accounts for the given party.
// func (n *DataNode) PartyAccounts(req *dataapipb.PartyAccountsRequest) (response *dataapipb.PartyAccountsResponse, err error) {
// 	msg := "gRPC call failed (data-node): PartyAccounts: %w"
// 	if n == nil {
// 		err = fmt.Errorf(msg, e.ErrNil)
// 		return
// 	}

// 	if n.Conn.GetState() != connectivity.Ready {
// 		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
// 		return
// 	}

// 	c := dataapipb.NewTradingDataServiceClient(n.Conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
// 	defer cancel()

// 	response, err = c.PartyAccounts(ctx, req)
// 	if err != nil {
// 		err = fmt.Errorf(msg, e.ErrorDetail(err))
// 	}
// 	return
// }

// // GetMarket returns market data for the specified market.
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

// ListMarkets returns all markets.
func (n *DataNode) ListMarkets(req *dataapipb.ListMarketsRequest) (response *dataapipb.ListMarketsResponse, err error) {
	msg := "gRPC call failed (data-node): ListMarkets: %w"
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

	response, err = c.ListMarkets(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

// // PositionsByParty returns positions for the given party.
// func (n *DataNode) PositionsByParty(req *dataapipb.PositionsByPartyRequest) (response *dataapipb.PositionsByPartyResponse, err error) {
// 	msg := "gRPC call failed (data-node): PositionsByParty: %w"
// 	if n == nil {
// 		err = fmt.Errorf(msg, e.ErrNil)
// 		return
// 	}

// 	if n.Conn.GetState() != connectivity.Ready {
// 		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
// 		return
// 	}

// 	c := dataapipb.NewTradingDataServiceClient(n.Conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
// 	defer cancel()

// 	response, err = c.PositionsByParty(ctx, req)
// 	if err != nil {
// 		err = fmt.Errorf(msg, e.ErrorDetail(err))
// 	}
// 	return
// }

// // AssetByID returns the specified asset.
// func (n *DataNode) AssetByID(req *dataapipb.AssetByIDRequest) (response *dataapipb.AssetByIDResponse, err error) {
// 	msg := "gRPC call failed (data-node): AssetByID: %w"
// 	if n == nil {
// 		err = fmt.Errorf(msg, e.ErrNil)
// 		return
// 	}

// 	if n.Conn.GetState() != connectivity.Ready {
// 		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
// 		return
// 	}

// 	c := dataapipb.NewTradingDataServiceClient(n.Conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
// 	defer cancel()

// 	response, err = c.AssetByID(ctx, req)
// 	if err != nil {
// 		err = fmt.Errorf(msg, e.ErrorDetail(err))
// 	}
// 	return
// }
