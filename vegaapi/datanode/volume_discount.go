package datanode

import (
	"context"
	"fmt"

	e "github.com/vegaprotocol/devopstools/errors"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"

	"google.golang.org/grpc/connectivity"
)

func (n *DataNode) GetCurrentVolumeDiscountProgram(ctx context.Context) (*dataapipb.VolumeDiscountProgram, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return nil, e.ErrConnectionNotReady
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	response, err := c.GetCurrentVolumeDiscountProgram(reqCtx, &dataapipb.GetCurrentVolumeDiscountProgramRequest{})
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", e.ErrorDetail(err))
	}

	return response.CurrentVolumeDiscountProgram, nil
}
