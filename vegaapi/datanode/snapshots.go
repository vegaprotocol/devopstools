package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegaeventspb "code.vegaprotocol.io/vega/protos/vega/events/v1"

	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) ListCoreSnapshots(ctx context.Context) ([]vegaeventspb.CoreSnapshotData, error) {
	if n == nil || n.Conn == nil {
		return nil, fmt.Errorf("data-node object cannot be nil")
	}

	if n.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf("data-node connection is not ready")
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	response, err := c.ListCoreSnapshots(reqCtx, &dataapipb.ListCoreSnapshotsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshot from the data-node: %w", err)
	}

	if response.CoreSnapshots == nil || len(response.CoreSnapshots.Edges) < 1 {
		return nil, fmt.Errorf("returned empty response from the List snapshot endpoint")
	}

	var result []vegaeventspb.CoreSnapshotData

	for idx, edge := range response.CoreSnapshots.Edges {
		if edge.Node == nil {
			continue
		}

		result = append(result, *(response.CoreSnapshots.Edges[idx].Node))
	}

	return result, nil
}
