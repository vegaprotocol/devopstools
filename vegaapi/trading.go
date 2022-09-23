package vegaapi

// === TradingDataService ===

// // PartyAccounts returns accounts for the given party.
// func (n *DataNode) PartyAccounts(req *dataapipb.PartyAccountsRequest) (response *dataapipb.PartyAccountsResponse, err error) {
// 	msg := "gRPC call failed (data-node): PartyAccounts: %w"
// 	if n == nil {
// 		err = fmt.Errorf(msg, e.ErrNil)
// 		return
// 	}

// 	if n.conn.GetState() != connectivity.Ready {
// 		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
// 		return
// 	}

// 	c := dataapipb.NewTradingDataServiceClient(n.conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
// 	defer cancel()

// 	response, err = c.PartyAccounts(ctx, req)
// 	if err != nil {
// 		err = fmt.Errorf(msg, e.ErrorDetail(err))
// 	}
// 	return
// }

// // MarketDataByID returns market data for the specified market.
// func (n *DataNode) MarketDataByID(req *dataapipb.MarketDataByIDRequest) (response *dataapipb.MarketDataByIDResponse, err error) {
// 	msg := "gRPC call failed (data-node): MarketDataByID: %w"
// 	if n == nil {
// 		err = fmt.Errorf(msg, e.ErrNil)
// 		return
// 	}

// 	if n.conn.GetState() != connectivity.Ready {
// 		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
// 		return
// 	}

// 	c := dataapipb.NewTradingDataServiceClient(n.conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
// 	defer cancel()

// 	response, err = c.MarketDataByID(ctx, req)
// 	if err != nil {
// 		err = fmt.Errorf(msg, e.ErrorDetail(err))
// 	}
// 	return
// }

// // Markets returns all markets.
// func (n *DataNode) Markets(req *dataapipb.MarketsRequest) (response *dataapipb.MarketsResponse, err error) {
// 	msg := "gRPC call failed (data-node): Markets: %w"
// 	if n == nil {
// 		err = fmt.Errorf(msg, e.ErrNil)
// 		return
// 	}

// 	if n.conn.GetState() != connectivity.Ready {
// 		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
// 		return
// 	}

// 	c := dataapipb.NewTradingDataServiceClient(n.conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
// 	defer cancel()

// 	response, err = c.Markets(ctx, req)
// 	if err != nil {
// 		err = fmt.Errorf(msg, e.ErrorDetail(err))
// 	}
// 	return
// }

// // PositionsByParty returns positions for the given party.
// func (n *DataNode) PositionsByParty(req *dataapipb.PositionsByPartyRequest) (response *dataapipb.PositionsByPartyResponse, err error) {
// 	msg := "gRPC call failed (data-node): PositionsByParty: %w"
// 	if n == nil {
// 		err = fmt.Errorf(msg, e.ErrNil)
// 		return
// 	}

// 	if n.conn.GetState() != connectivity.Ready {
// 		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
// 		return
// 	}

// 	c := dataapipb.NewTradingDataServiceClient(n.conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
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

// 	if n.conn.GetState() != connectivity.Ready {
// 		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
// 		return
// 	}

// 	c := dataapipb.NewTradingDataServiceClient(n.conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
// 	defer cancel()

// 	response, err = c.AssetByID(ctx, req)
// 	if err != nil {
// 		err = fmt.Errorf(msg, e.ErrorDetail(err))
// 	}
// 	return
// }
