package networktools

import (
	"context"
	"fmt"
	"time"

	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"
	"go.uber.org/zap"
)

func (network *NetworkTools) GetDataNodeClient() (vegaapi.DataNodeClient, error) {
	addresses := network.GetNetworkGRPCDataNodes()
	if len(addresses) == 0 {
		return nil, fmt.Errorf("there is no single healthy GRPC endpoint for '%s'", network.Name)
	}
	node := datanode.NewDataNode(addresses, 3*time.Second, network.logger)

	network.logger.Debug("Attempting to connect to Vega gRPC node...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	node.MustDialConnection(ctx) // blocking
	network.logger.Debug("Connected to Vega gRPC node", zap.String("node", node.Target()))

	return node, nil
}
