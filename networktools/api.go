package networktools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/vegaapi/core"
	"github.com/vegaprotocol/devopstools/vegaapi/datanode"

	"go.uber.org/zap"
)

func (network *NetworkTools) GetVegaCoreClient() (vegaapi.VegaCoreClient, error) {
	addresses := network.GetNetworkGRPCVegaCore()
	if len(addresses) == 0 {
		return nil, fmt.Errorf("there is no single healthy Vega Core GRPC endpoint for '%s'", network.Name)
	}
	node := core.NewCoreClient(addresses, 3*time.Second, network.logger)

	network.logger.Debug("Attempting to connect to Vega gRPC node...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	node.MustDialConnection(ctx) // blocking
	network.logger.Debug("Connected to Vega gRPC node", zap.String("node", node.Target()))

	return node, nil
}

func (network *NetworkTools) GetVegaCoreClientForNode(nodeId string) (vegaapi.VegaCoreClient, error) {
	addresses := network.GetNetworkGRPCVegaCore()
	if len(addresses) == 0 {
		return nil, fmt.Errorf("there is no single healthy Vega Core GRPC endpoint for '%s'", network.Name)
	}
	host := ""
	for _, address := range addresses {
		if strings.HasPrefix(address, nodeId) {
			host = address
			break
		}
	}
	if len(host) == 0 {
		return nil, fmt.Errorf("host '%s' is not healthy in '%s' network", nodeId, network.Name)
	}
	node := core.NewCoreClient([]string{host}, 3*time.Second, network.logger)

	network.logger.Debug("Attempting to connect to Vega gRPC node...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	node.MustDialConnection(ctx) // blocking
	network.logger.Debug("Connected to Vega gRPC node", zap.String("node", node.Target()))

	return node, nil
}

func (network *NetworkTools) GetDataNodeClient() (vegaapi.DataNodeClient, error) {
	addresses := network.GetNetworkGRPCDataNodes()
	if len(addresses) == 0 {
		return nil, fmt.Errorf("there is no single healthy Data-Node GRPC endpoint for '%s'", network.Name)
	}
	node := datanode.NewDataNode(addresses, 3*time.Second, network.logger)

	network.logger.Debug("Attempting to connect to Vega gRPC node...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	node.MustDialConnection(ctx) // blocking
	network.logger.Debug("Connected to Vega gRPC node", zap.String("node", node.Target()))

	return node, nil
}
