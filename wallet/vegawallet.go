package wallet

import (
	"fmt"

	"code.vegaprotocol.io/vega/commands"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	wcommands "code.vegaprotocol.io/vega/wallet/commands"
	"code.vegaprotocol.io/vega/wallet/wallet"
	"github.com/vegaprotocol/devopstools/secrets"
)

type VegaWallet struct {
	*secrets.VegaWalletPrivate

	hdWallet *wallet.HDWallet
	keyPair  *wallet.KeyPair
}

func NewVegaWallet(
	name string,
	private *secrets.VegaWalletPrivate,
) (*VegaWallet, error) {
	hdWallet, err := wallet.ImportHDWallet(name, private.RecoveryPhrase, wallet.LatestVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to create new VegaWallet, %w", err)
	}
	keyPair, err := hdWallet.GenerateKeyPair(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair %w", err)
	}
	return &VegaWallet{
		VegaWalletPrivate: private,
		hdWallet:          hdWallet,
		keyPair:           &keyPair,
	}, nil
}

func (vw *VegaWallet) SignTx(req *walletpb.SubmitTransactionRequest, height uint64, chainID string) (*commandspb.Transaction, error) {
	marshaledInputData, err := wcommands.ToMarshaledInputData(req, height)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal input data: %w", err)
	}

	pubKey := req.GetPubKey()
	signature, err := vw.hdWallet.SignTx(pubKey, commands.BundleInputDataForSigning(marshaledInputData, chainID))
	if err != nil {
		return nil, fmt.Errorf("couldn't sign transaction: %w", err)
	}

	protoSignature := &commandspb.Signature{
		Value:   signature.Value,
		Algo:    signature.Algo,
		Version: signature.Version,
	}
	return commands.NewTransaction(pubKey, marshaledInputData, protoSignature), nil
}
