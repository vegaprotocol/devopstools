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
	res, err := n.getCurrentReferralProgramRaw(ctx, &dataapipb.GetCurrentReferralProgramRequest{})
	if err != nil {
		return nil, err
	}
	return res.CurrentReferralProgram, nil
}

func (n *DataNode) getCurrentReferralProgramRaw(ctx context.Context, req *v2.GetCurrentReferralProgramRequest) (*v2.GetCurrentReferralProgramResponse, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.CallTimeout)
	defer cancelRequest()

	response, err := c.GetCurrentReferralProgram(reqCtx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	return response, nil
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

func (n *DataNode) GetReferralSetReferees(ctx context.Context) (map[string]v2.ReferralSetReferee, error) {
	referralSetReferees := map[string]v2.ReferralSetReferee{}
	var pagination *v2.Pagination = nil
	for {
		res, err := n.ListReferralSetReferees(ctx, &dataapipb.ListReferralSetRefereesRequest{
			Pagination: pagination,
		})
		if err != nil {
			return nil, err
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

func (n *DataNode) ListReferralSetReferees(ctx context.Context, req *v2.ListReferralSetRefereesRequest) (*v2.ListReferralSetRefereesResponse, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancelRequest := context.WithTimeout(ctx, n.CallTimeout)
	defer cancelRequest()

	response, err := c.ListReferralSetReferees(reqCtx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}
	return response, nil
}
