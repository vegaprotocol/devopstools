package datanode

import (
	"context"
	"fmt"

	e "github.com/vegaprotocol/devopstools/errors"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"

	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) LastNetworkHistorySegment(ctx context.Context) (*dataapipb.HistorySegment, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	response, err := c.GetMostRecentNetworkHistorySegment(reqCtx, &dataapipb.GetMostRecentNetworkHistorySegmentRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get most recent network history segment: %w", err)
	}
	if response.Segment == nil {
		return nil, fmt.Errorf("empty response from get most recent network history segment")
	}

	return response.Segment, nil
}
