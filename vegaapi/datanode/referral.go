package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	v2 "code.vegaprotocol.io/vega/protos/data-node/api/v2"
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

func (n *DataNode) GetReferralSets() (map[string]v2.ReferralSet, error) {
	res, err := n.ListReferralSets(&dataapipb.ListReferralSetsRequest{})
	if err != nil {
		return nil, err
	}
	referralSets := map[string]v2.ReferralSet{}
	for _, edge := range res.ReferralSets.Edges {
		referralSets[edge.Node.Referrer] = *edge.Node
	}
	return referralSets, nil
}

func (n *DataNode) ListReferralSets(req *dataapipb.ListReferralSetsRequest) (response *dataapipb.ListReferralSetsResponse, err error) {
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

//
// Referees
//

func (n *DataNode) GetReferralSetReferees() (map[string]v2.ReferralSetReferee, error) {
	res, err := n.ListReferralSetReferees(&dataapipb.ListReferralSetRefereesRequest{})
	if err != nil {
		return nil, err
	}
	referralSetReferees := map[string]v2.ReferralSetReferee{}
	for _, edge := range res.ReferralSetReferees.Edges {
		referralSetReferees[edge.Node.Referee] = *edge.Node
	}
	return referralSetReferees, nil

}

func (n *DataNode) ListReferralSetReferees(req *dataapipb.ListReferralSetRefereesRequest) (response *dataapipb.ListReferralSetRefereesResponse, err error) {
	msg := "gRPC call failed (data-node): ListReferralSetReferees: %w"
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

	response, err = c.ListReferralSetReferees(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}
