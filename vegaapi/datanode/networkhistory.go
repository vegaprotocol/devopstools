package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"

	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) LastNetworkHistorySegment() (*dataapipb.HistorySegment, error) {
	if n == nil || n.Conn == nil {
		return nil, fmt.Errorf("data-node object cannot be nil")
	}

	if n.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf("data-node connection is not ready")
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err := c.GetMostRecentNetworkHistorySegment(ctx, &dataapipb.GetMostRecentNetworkHistorySegmentRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get most recent network history segment: %w", err)
	}
	if response.Segment == nil {
		return nil, fmt.Errorf("empty response from get most recent network history segment")
	}

	return response.Segment, nil
}
