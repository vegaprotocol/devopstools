package datanode

import (
	"time"

	"github.com/vegaprotocol/devopstools/vegaapi/core"

	"go.uber.org/zap"
)

// DataNode stores state for a Vega Data node.
type DataNode struct {
	*core.Client
}

// New returns a new node.
func New(hosts []string, callTimeout time.Duration, logger *zap.Logger) *DataNode {
	return &DataNode{
		Client: core.NewClient(hosts, callTimeout, logger),
	}
}
