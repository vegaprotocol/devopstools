package datanode

import (
	"time"

	"github.com/vegaprotocol/devopstools/vegaapi/core"

	"go.uber.org/zap"
)

// DataNode stores state for a Vega Data node.
type DataNode struct {
	core.CoreClient
}

// NewDataNode returns a new node.
func NewDataNode(hosts []string, callTimeout time.Duration, logger *zap.Logger) *DataNode {
	return &DataNode{
		CoreClient: *core.NewCoreClient(hosts, callTimeout, logger),
	}
}
