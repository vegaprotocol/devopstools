package datanode

import (
	"context"
	"fmt"
	"math/big"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	e "github.com/vegaprotocol/devopstools/errors"
	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetPartyTotalStake(partyId string) (*big.Int, error) {
	res, err := n.GetStake(&dataapipb.GetStakeRequest{
		PartyId: partyId,
	})
	if err != nil {
		return nil, err
	}
	totalStake := new(big.Int)
	totalStake, ok := totalStake.SetString(res.CurrentStakeAvailable, 0)
	if !ok {
		return nil, fmt.Errorf("failed to convert %s to big.Int", res.CurrentStakeAvailable)
	}

	return totalStake, nil
}

func (n *DataNode) GetPartyDelegationToNode(partyId string, nodeId string) (*big.Int, error) {
	epoch, err := n.GetCurrentEpoch()
	if err != nil {
		return nil, err
	}
	amount := "0"
	for _, delegation := range epoch.Delegations {
		if delegation.Party == partyId && delegation.NodeId == nodeId {
			amount = delegation.Amount
			break
		}
	}
	result := new(big.Int)
	result, ok := result.SetString(amount, 0)
	if !ok {
		return nil, fmt.Errorf("failed to convert %s to big.Int", amount)
	}
	return result, nil
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

// ListDelegations returns delegations for the given party.
func (n *DataNode) ListDelegations(req *dataapipb.ListDelegationsRequest) (response *dataapipb.ListDelegationsResponse, err error) {
	msg := "gRPC call failed (data-node): ListDelegations: %w"
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

	response, err = c.ListDelegations(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}
