package datanode

import (
	"context"
	"fmt"

	e "github.com/vegaprotocol/devopstools/errors"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"

	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetCurrentReferralProgram(ctx context.Context) (*v2.ReferralProgram, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.CallTimeout)
	defer cancelRequest()

	response, err := c.GetCurrentReferralProgram(reqCtx, &dataapipb.GetCurrentReferralProgramRequest{})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	return response.CurrentReferralProgram, nil
}

func (n *DataNode) ListReferralSets(ctx context.Context) (map[string]*v2.ReferralSet, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.CallTimeout)
	defer cancelRequest()

	response, err := c.ListReferralSets(reqCtx, &dataapipb.ListReferralSetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	referralSets := map[string]*v2.ReferralSet{}
	for _, edge := range response.ReferralSets.Edges {
		referralSets[edge.Node.Referrer] = edge.Node
	}
	return referralSets, nil
}

func (n *DataNode) ListReferralSetReferees(ctx context.Context) (map[string]v2.ReferralSetReferee, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.CallTimeout)
	defer cancelRequest()

	referralSetReferees := map[string]v2.ReferralSetReferee{}
	var pagination *v2.Pagination = nil
	for {
		res, err := c.ListReferralSetReferees(reqCtx, &dataapipb.ListReferralSetRefereesRequest{
			Pagination: pagination,
		})
		if err != nil {
			return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
		}

		for _, edge := range res.ReferralSetReferees.Edges {
			referralSetReferees[edge.Node.Referee] = *edge.Node
		}
		if res.ReferralSetReferees.PageInfo.HasNextPage {
			pagination = &v2.Pagination{
				After: &res.ReferralSetReferees.PageInfo.EndCursor,
			}
		} else {
			break
		}
	}
	return referralSetReferees, nil
}
