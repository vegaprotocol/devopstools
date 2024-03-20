package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"

	"google.golang.org/grpc/connectivity"
)

func (dn *DataNode) LastNetworkHistorySegment() (*dataapipb.HistorySegment, error) {
	if dn == nil || dn.Conn == nil {
		return nil, fmt.Errorf("data-node object cannot be nil")
	}

	if dn.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf("data-node connection is not ready")
	}

	c := dataapipb.NewTradingDataServiceClient(dn.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), dn.CallTimeout)
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
