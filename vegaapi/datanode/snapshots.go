package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegaeventspb "code.vegaprotocol.io/vega/protos/vega/events/v1"

	"google.golang.org/grpc/connectivity"
)

func (dn *DataNode) ListCoreSnapshots() ([]vegaeventspb.CoreSnapshotData, error) {
	if dn == nil || dn.Conn == nil {
		return nil, fmt.Errorf("data-node object cannot be nil")
	}

	if dn.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf("data-node connection is not ready")
	}

	c := dataapipb.NewTradingDataServiceClient(dn.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), dn.CallTimeout)
	defer cancel()

	response, err := c.ListCoreSnapshots(ctx, &dataapipb.ListCoreSnapshotsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshot from the data-node: %w", err)
	}

	if response.CoreSnapshots == nil || len(response.CoreSnapshots.Edges) < 1 {
		return nil, fmt.Errorf("returned empty response from the List snapshot endpoint")
	}

	result := []vegaeventspb.CoreSnapshotData{}

	for idx, edge := range response.CoreSnapshots.Edges {
		if edge.Node == nil {
			continue
		}

		result = append(result, *(response.CoreSnapshots.Edges[idx].Node))
	}

	return result, nil
}
