package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	v1 "code.vegaprotocol.io/vega/protos/vega/events/v1"
	e "github.com/vegaprotocol/devopstools/errors"
	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) ListAMMs(ctx context.Context, partyId *string, marketId *string, activeOnly bool) ([]*v1.AMM, error) {
	if n.Client.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Client.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.Client.CallTimeout)
	defer cancelRequest()

	var status *v1.AMM_Status
	if activeOnly {
		status = v1.AMM_STATUS_ACTIVE.Enum()
	}

	response, err := c.ListAMMs(reqCtx, &dataapipb.ListAMMsRequest{
		PartyId:  partyId,
		MarketId: marketId,
		Status:   status,
		Pagination: &dataapipb.Pagination{
			NewestFirst: &([]bool{true})[0],
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get latest market data: %w", err)
	}

	result := []*v1.AMM{}

	for i := range response.Amms.Edges {
		result = append(result, response.Amms.Edges[i].Node)
	}

	return result, nil
}
