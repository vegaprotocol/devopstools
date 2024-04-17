package core

import (
	"fmt"

	"code.vegaprotocol.io/vega/libs/proto"
	vgrand "code.vegaprotocol.io/vega/libs/rand"
	types "code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
)

func (n *Client) DepositBuiltinAsset(
	vegaAssetId string,
	partyId string,
	amount string,
	signAny func([]byte) ([]byte, string, error),
) (bool, error) {
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

	resp, err := n.PropagateChainEvent(&vegaapipb.PropagateChainEventRequest{
		Event:     msg,
		PubKey:    sigPubKey,
		Signature: sig,
	})
	if err != nil {
		return false, fmt.Errorf("failed to deposit built-in asset, %w", err)
	}
	return resp.Success, nil
}
