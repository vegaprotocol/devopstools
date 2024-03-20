package datanode

import (
	"context"
	"fmt"

	e "github.com/vegaprotocol/devopstools/errors"

	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"

	"google.golang.org/grpc/connectivity"
)

//
// CURRENT PROGRAM
//

func (n *DataNode) GetCurrentVolumeDiscountProgram() (*dataapipb.VolumeDiscountProgram, error) {
	res, err := n.GetCurrentVolumeDiscountProgramRaw(&dataapipb.GetCurrentVolumeDiscountProgramRequest{})
	if err != nil {
		return nil, err
	}
	return res.CurrentVolumeDiscountProgram, nil
}

func (n *DataNode) GetCurrentVolumeDiscountProgramRaw(req *dataapipb.GetCurrentVolumeDiscountProgramRequest) (response *dataapipb.GetCurrentVolumeDiscountProgramResponse, err error) {
	msg := "gRPC call failed (data-node): GetCurrentVolumeDiscountProgram: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err = c.GetCurrentVolumeDiscountProgram(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

//
// STATS
//

func (n *DataNode) GetVolumeDiscountStats() ([]*dataapipb.VolumeDiscountStatsEdge, error) {
	res, err := n.GetVolumeDiscountStatsRaw(&dataapipb.GetVolumeDiscountStatsRequest{})
	if err != nil {
		return nil, err
	}
	return res.Stats.Edges, nil
}

func (n *DataNode) GetVolumeDiscountStatsRaw(req *dataapipb.GetVolumeDiscountStatsRequest) (response *dataapipb.GetVolumeDiscountStatsResponse, err error) {
	msg := "gRPC call failed (data-node): GetVolumeDiscountStats: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := dataapipb.NewTradingDataServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err = c.GetVolumeDiscountStats(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}
