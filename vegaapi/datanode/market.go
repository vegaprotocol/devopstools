package datanode

import (
	"context"
	"fmt"
	"slices"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	e "github.com/vegaprotocol/devopstools/errors"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"google.golang.org/grpc/connectivity"
)

//
// MARKET
//

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

//
// LATEST MARKET DATA
//

func (n *DataNode) ListLatestMarketData(req *dataapipb.ListLatestMarketDataRequest) (response *dataapipb.ListLatestMarketDataResponse, err error) {
	msg := "gRPC call failed (data-node): ListLatestMarketData: %w"
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

	response, err = c.ListLatestMarketData(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *DataNode) GetLatestMarketData(req *dataapipb.GetLatestMarketDataRequest) (response *dataapipb.GetLatestMarketDataResponse, err error) {
	msg := "gRPC call failed (data-node): GetLatestMarketData: %w"
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

	response, err = c.GetLatestMarketData(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

//
// MARKET INFO
//

var marketTradeableStates = []vega.Market_State{
	vega.Market_STATE_PENDING,
	vega.Market_STATE_ACTIVE,
	vega.Market_STATE_SUSPENDED,
}

func (n *DataNode) GetTradeableMakertInfo() ([]vegaapi.MarketInfo, error) {
	allMarketInfo, err := n.GetAllMakertInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get tradeable market info, %w", err)
	}
	tradeableMarketInfo := []vegaapi.MarketInfo{}
	for _, marketInfo := range allMarketInfo {
		if slices.Contains(marketTradeableStates, marketInfo.Market.State) {
			tradeableMarketInfo = append(tradeableMarketInfo, marketInfo)
		}
	}

	return tradeableMarketInfo, nil
}

func (n *DataNode) GetAllMakertInfo() ([]vegaapi.MarketInfo, error) {
	var (
		marketById     = map[string]*vega.Market{}
		marketDataById = map[string]*vega.MarketData{}
		assetById      = map[string]*vega.Asset{}
	)
	// get marketById
	marketsRes, err := n.ListMarkets(&dataapipb.ListMarketsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all markets info, %w", err)
	}
	for _, edge := range marketsRes.Markets.Edges {
		marketById[edge.Node.Id] = edge.Node
	}
	// get marketDataById
	latestMarketDataRes, err := n.ListLatestMarketData(&dataapipb.ListLatestMarketDataRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all markets info, %w", err)
	}
	for _, marketData := range latestMarketDataRes.MarketsData {
		marketDataById[marketData.Market] = marketData
	}
	// get assetById
	assetsRes, err := n.ListAssets(&dataapipb.ListAssetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all markets info, %w", err)
	}
	for _, edge := range assetsRes.Assets.Edges {
		assetById[edge.Node.Id] = edge.Node
	}
	// compose result
	allMarketInfo := []vegaapi.MarketInfo{}
	for marketId, market := range marketById {
		marketData := marketDataById[marketId] // can be null
		// get asset
		settlementAssetId := ""
		if market.TradableInstrument.Instrument.GetFuture() != nil {
			settlementAssetId = market.TradableInstrument.Instrument.GetFuture().SettlementAsset
		}
		if market.TradableInstrument.Instrument.GetPerpetual() != nil {
			settlementAssetId = market.TradableInstrument.Instrument.GetPerpetual().SettlementAsset
		}
		if len(settlementAssetId) == 0 {
			return nil, fmt.Errorf("failed to get all markets info, there is not settlement asset for market %s", market.TradableInstrument.Instrument.Name)
		}
		settlementAsset, ok := assetById[settlementAssetId]
		if !ok {
			return nil, fmt.Errorf("failed to get all markets info, there is not settlement asset with id %s for market %s",
				settlementAssetId, market.TradableInstrument.Instrument.Name)
		}
		allMarketInfo = append(allMarketInfo, vegaapi.MarketInfo{
			Market:          market,
			MarketData:      marketData,
			SettlementAsset: settlementAsset,
		})

	}

	return allMarketInfo, nil
}

// func (n *DataNode) GetTradeableMarketInfoWithName(name string) (vegaapi.MarketInfo, error) {
// 	//
// 	// GET all markets
// 	//
// 	markets, err := n.GetAllMarkets()
// 	if err != nil {
// 		return vegaapi.MarketInfo{}, fmt.Errorf("failed to get market with name %s, %w", name, err)
// 	}
// 	//
// 	// FILTER
// 	//
// 	matchedMarkets := []*vega.Market{}

// 	for _, market := range markets {
// 		if slices.Contains(marketTradeableStates, market.State) &&
// 			market.TradableInstrument.Instrument.Name == name {
// 			matchedMarkets = append(matchedMarkets, market)
// 		}
// 	}
// 	//
// 	// VALIDATE there is only ONE market matching filters
// 	//
// 	if len(matchedMarkets) == 0 {
// 		marketNames := []string{}
// 		for _, market := range markets {
// 			marketNames = append(marketNames, market.TradableInstrument.Instrument.Name)
// 		}
// 		return vegaapi.MarketInfo{}, fmt.Errorf("no market with name %s, available names: %s", name, marketNames)
// 	}
// 	if len(matchedMarkets) > 1 {
// 		return vegaapi.MarketInfo{}, fmt.Errorf("found more than one market with name %s", name)
// 	}

// 	market := matchedMarkets[0]
// 	//
// 	// GET latest market data
// 	//
// 	marketData, err := n.GetLatestMarketData(&dataapipb.GetLatestMarketDataRequest{MarketId: market.Id})
// 	if err != nil {
// 		return vegaapi.MarketInfo{}, fmt.Errorf("failed to get market with name %s, %w", name, err)
// 	}
// 	//
// 	// GET settlement asset
// 	//
// 	settlementAssetId := ""
// 	if market.TradableInstrument.Instrument.GetFuture() != nil {
// 		settlementAssetId = market.TradableInstrument.Instrument.GetFuture().SettlementAsset
// 	}
// 	if market.TradableInstrument.Instrument.GetPerpetual() != nil {
// 		settlementAssetId = market.TradableInstrument.Instrument.GetPerpetual().SettlementAsset
// 	}
// 	if len(settlementAssetId) == 0 {
// 		return vegaapi.MarketInfo{}, fmt.Errorf("there is not settlement asset for market %s", market.TradableInstrument.Instrument.Name)
// 	}
// 	settlementAsset, err := n.GetAssetById(settlementAssetId)
// 	if err != nil {
// 		return vegaapi.MarketInfo{}, fmt.Errorf("failed to get asset for id %s, %w", settlementAssetId, err)
// 	}

// 	return vegaapi.MarketInfo{
// 		Market:          market,
// 		MarketData:      marketData.MarketData,
// 		SettlementAsset: settlementAsset,
// 	}, nil
// }
