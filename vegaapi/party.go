package vegaapi

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	e "github.com/vegaprotocol/devopstools/errors"
	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetPartyStake(partyId string) (string, error) {
	res, err := n.GetStake(&dataapipb.GetStakeRequest{
		PartyId: partyId,
	})
	if err != nil {
		return "", err
	}
	return res.CurrentStakeAvailable, nil
}

// GetStake returns stakes for the given party.
func (n *DataNode) GetStake(req *dataapipb.GetStakeRequest) (response *dataapipb.GetStakeResponse, err error) {
	msg := "gRPC call failed (data-node): GetStake: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.callTimeout)
	defer cancel()

	response, err = c.GetStake(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}
