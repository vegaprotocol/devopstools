package datanode

import (
	"context"
	"fmt"

	e "github.com/vegaprotocol/devopstools/errors"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"

	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) ListGovernanceData(ctx context.Context, req *dataapipb.ListGovernanceDataRequest) (*dataapipb.ListGovernanceDataResponse, error) {
	if n.Client.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Client.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.Client.CallTimeout)
	defer cancel()

	response, err := c.ListGovernanceData(reqCtx, req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return response, nil
}

func (n *DataNode) ListVotes(ctx context.Context, req *dataapipb.ListVotesRequest) (*dataapipb.ListVotesResponse, error) {
	if n.Client.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Client.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.Client.CallTimeout)
	defer cancel()

	response, err := c.ListVotes(reqCtx, req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return response, nil
}

func (n *DataNode) GetGovernanceData(ctx context.Context, req *dataapipb.GetGovernanceDataRequest) (*dataapipb.GetGovernanceDataResponse, error) {
	if n.Client.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Client.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.Client.CallTimeout)
	defer cancel()

	response, err := c.GetGovernanceData(reqCtx, req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return response, nil
}
