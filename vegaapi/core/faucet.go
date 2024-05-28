package core

import (
	"context"
	"fmt"

	e "github.com/vegaprotocol/devopstools/errors"

	"code.vegaprotocol.io/vega/libs/proto"
	vgrand "code.vegaprotocol.io/vega/libs/rand"
	types "code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"

	"google.golang.org/grpc/connectivity"
)

func (n *Client) DepositBuiltinAsset(ctx context.Context, vegaAssetId string, partyId string, amount string, signAny func([]byte) ([]byte, string, error)) (bool, error) {
	if n.Conn.GetState() != connectivity.Ready {
		return false, e.ErrConnectionNotReady
	}

	c := vegaapipb.NewCoreServiceClient(n.Conn)
	reqCtx, cancel := context.WithTimeout(ctx, n.CallTimeout)
	defer cancel()

	chainEvent := &commandspb.ChainEvent{
		Nonce: vgrand.NewNonce(),
		Event: &commandspb.ChainEvent_Builtin{
			Builtin: &types.BuiltinAssetEvent{
				Action: &types.BuiltinAssetEvent_Deposit{
					Deposit: &types.BuiltinAssetDeposit{
						VegaAssetId: vegaAssetId,
						PartyId:     partyId,
						Amount:      amount,
					},
				},
			},
		},
	}
	msg, err := proto.Marshal(chainEvent)
	if err != nil {
		return false, fmt.Errorf("failed to deposit built-in asset, %w", err)
	}

	sig, sigPubKey, err := signAny(msg)
	if err != nil {
		return false, fmt.Errorf("failed to deposit built-in asset, %w", err)
	}

	resp, err := c.PropagateChainEvent(reqCtx, &vegaapipb.PropagateChainEventRequest{
		Event:     msg,
		PubKey:    sigPubKey,
		Signature: sig,
	})
	if err != nil {
		return false, fmt.Errorf("failed to deposit built-in asset, %w", err)
	}
	return resp.Success, nil
}
