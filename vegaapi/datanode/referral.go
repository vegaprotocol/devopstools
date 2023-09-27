package datanode

import (
	"context"
	"fmt"

	vegaTypes "code.vegaprotocol.io/vega/core/types"
	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	e "github.com/vegaprotocol/devopstools/errors"
	"google.golang.org/grpc/connectivity"
)

//
// CURRENT PROGRAM
//

func (n *DataNode) GetCurrentReferralProgram() (*dataapipb.ReferralProgram, error) {
	res, err := n.GetCurrentReferralProgramRaw(&dataapipb.GetCurrentReferralProgramRequest{})
	if err != nil {
		return nil, err
	}
	return res.CurrentReferralProgram, nil
}

func (n *DataNode) GetCurrentReferralProgramRaw(req *dataapipb.GetCurrentReferralProgramRequest) (response *dataapipb.GetCurrentReferralProgramResponse, err error) {
	msg := "gRPC call failed (data-node): GetCurrentReferralProgram: %w"
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

	response, err = c.GetCurrentReferralProgram(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

//
// LIST REFERRAL SETS
//

func (n *DataNode) ListReferralSets() ([]vegaTypes.ReferralSet, error) {
	res, err := n.ListReferralSetsRaw(&dataapipb.ListReferralSetsRequest{})
	if err != nil {
		return nil, err
	}
	referralSets := make([]vegaTypes.ReferralSet, len(res.ReferralSets.Edges))
	for i, edge := range res.ReferralSets.Edges {
		referralSets[i] = *edge.Node
	}
	return referralSets, nil
}

func (n *DataNode) ListReferralSetsRaw(req *dataapipb.ListReferralSetsRequest) (response *dataapipb.ListReferralSetsResponse, err error) {
	msg := "gRPC call failed (data-node): ListReferralSets: %w"
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

	response, err = c.ListReferralSets(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}
