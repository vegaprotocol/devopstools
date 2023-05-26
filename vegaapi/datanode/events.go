package datanode

import (
	"context"
	"fmt"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegaeventspb "code.vegaprotocol.io/vega/protos/vega/events/v1"
	"google.golang.org/grpc/connectivity"
)

func (dn *DataNode) ListProtocolUpgradeProposals() ([]vegaeventspb.ProtocolUpgradeEvent, error) {
	if dn == nil {
		return nil, fmt.Errorf("data-node object cannot be nil")
	}

	if dn.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf("data-node connection is not ready")
	}

	c := dataapipb.NewTradingDataServiceClient(dn.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), dn.CallTimeout)
	defer cancel()

	result := []vegaeventspb.ProtocolUpgradeEvent{}

	response, err := c.ListProtocolUpgradeProposals(ctx, &dataapipb.ListProtocolUpgradeProposalsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to call list protocol upgrade proposalst: %w", err)
	}

	if response == nil || response.ProtocolUpgradeProposals == nil || len(response.ProtocolUpgradeProposals.Edges) < 1 {
		return result, nil
	}

	for _, edge := range response.ProtocolUpgradeProposals.Edges {
		if edge == nil || edge.Node == nil {
			continue
		}

		result = append(result, *edge.Node)
	}

	return result, nil
}
