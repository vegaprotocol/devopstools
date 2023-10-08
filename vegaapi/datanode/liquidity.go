package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	e "github.com/vegaprotocol/devopstools/errors"
	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) LPForMarketByParty(marketId string, partyId string) ([]*vega.LiquidityProvision, error) {
	live := true
	res, err := n.ListLiquidityProvisions(&dataapipb.ListLiquidityProvisionsRequest{
		MarketId: &marketId,
		PartyId:  &partyId,
		Live:     &live,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get LP for market %s from party %s, %w", marketId, partyId, err)
	}
	lpList := []*vega.LiquidityProvision{}
	for _, edge := range res.LiquidityProvisions.Edges {
		lpList = append(lpList, edge.Node)
	}
	return lpList, nil
}

func (n *DataNode) ListLiquidityProvisions(req *dataapipb.ListLiquidityProvisionsRequest) (response *dataapipb.ListLiquidityProvisionsResponse, err error) {
	msg := "gRPC call failed (data-node): ListLiquidityProvisions: %w"
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

	response, err = c.ListLiquidityProvisions(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}
