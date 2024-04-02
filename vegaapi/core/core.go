package core

import (
	"context"
	"fmt"

	e "github.com/vegaprotocol/devopstools/errors"

	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"

	"google.golang.org/grpc/connectivity"
)

// === CoreService ===

// SubmitTransaction submits a signed v2 transaction.
func (n *CoreClient) SubmitTransaction(
	req *vegaapipb.SubmitTransactionRequest,
) (response *vegaapipb.SubmitTransactionResponse, err error) {
	msg := "gRPC call failed: SubmitTransaction: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err = c.SubmitTransaction(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

// LastBlockData gets the latest blockchain data, height, hash and pow parameters.
func (n *CoreClient) LastBlockData() (*vegaapipb.LastBlockHeightResponse, error) {
	msg := "gRPC call failed: LastBlockData: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()
	var response *vegaapipb.LastBlockHeightResponse
	response, err := c.LastBlockHeight(ctx, &vegaapipb.LastBlockHeightRequest{})
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return response, err
}

// ObserveEventBus opens a stream.
func (n *CoreClient) ObserveEventBus(
	ctx context.Context,
) (client vegaapipb.CoreService_ObserveEventBusClient, err error) {
	msg := "gRPC call failed: ObserveEventBus: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	// no timeout on streams
	client, err = c.ObserveEventBus(ctx)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
		return
	}
	return
}

func (n *CoreClient) Statistics() (*vegaapipb.StatisticsResponse, error) {
	msg := "gRPC call failed: Statistics: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()
	response, err := c.Statistics(ctx, &vegaapipb.StatisticsRequest{})
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return response, err
}

// PropagateChainEvent submits a signed v2 transaction.
func (n *CoreClient) PropagateChainEvent(
	req *vegaapipb.PropagateChainEventRequest,
) (response *vegaapipb.PropagateChainEventResponse, err error) {
	msg := "gRPC call failed: PropagateChainEvent: %w"
	if n == nil {
		err = fmt.Errorf(msg, e.ErrNil)
		return
	}

	if n.Conn.GetState() != connectivity.Ready {
		err = fmt.Errorf(msg, e.ErrConnectionNotReady)
		return
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	response, err = c.PropagateChainEvent(ctx, req)
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}
	return
}

func (n *CoreClient) CoreNetworkParameters(parameterKey string) ([]*vega.NetworkParameter, error) {
	if n == nil {
		return nil, fmt.Errorf("failed to get network parameters from the core api: %w", e.ErrNil)
	}

	if n.Conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(
			"failed to get network parameters from the core api: %w",
			e.ErrConnectionNotReady,
		)
	}

	c := vegaapipb.NewCoreStateServiceClient(n.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), n.CallTimeout)
	defer cancel()

	input := &vegaapipb.ListNetworkParametersRequest{}
	if parameterKey != "" {
		input.NetworkParameterKey = parameterKey
	}
	response, err := c.ListNetworkParameters(ctx, input)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get network parameters from the core api: %w",
			e.ErrorDetail(err),
		)
	}

	return response.NetworkParameters, nil
}
