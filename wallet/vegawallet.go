package wallet

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/secrets"

	"code.vegaprotocol.io/vega/commands"
	vgcrypto "code.vegaprotocol.io/vega/libs/crypto"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	wcommands "code.vegaprotocol.io/vega/wallet/commands"
	"code.vegaprotocol.io/vega/wallet/wallet"
)

type VegaWallet struct {
	*secrets.VegaWalletPrivate

	hdWallet *wallet.HDWallet
	keyPair  wallet.KeyPair
}

func NewVegaWallet(private *secrets.VegaWalletPrivate) (*VegaWallet, error) {
	hdWallet, err := wallet.ImportHDWallet(private.Id, private.RecoveryPhrase, wallet.LatestVersion)
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
		keyPair:           keyPair,
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

func (vw *VegaWallet) SignTxWithPoW(req *walletpb.SubmitTransactionRequest, lastBlockData *vegaapipb.LastBlockHeightResponse) (*commandspb.Transaction, error) {
	errMsg := fmt.Sprintf("failed to sing transaction %s", req)
	signedTx, err := vw.SignTx(req, lastBlockData.Height, lastBlockData.ChainId)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}

	tid := vgcrypto.RandomHash()
	powNonce, _, err := vgcrypto.PoW(lastBlockData.Hash, tid, uint(lastBlockData.SpamPowDifficulty), vgcrypto.Sha3)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", errMsg, err)
	}
	signedTx.Pow = &commandspb.ProofOfWork{
		Tid:   tid,
		Nonce: powNonce,
	}

	return signedTx, nil
}

func (vw *VegaWallet) SignAny(data []byte) ([]byte, string, error) {
	sig, err := vw.hdWallet.SignAny(vw.PublicKey, data)
	return sig, vw.PublicKey, err
}

func (vw *VegaWallet) DeriveKeyPair() (*VegaWallet, error) {
	if vw == nil || vw.hdWallet == nil {
		return nil, fmt.Errorf("hd wallet must be given for vega wallet to derive new keys")
	}

	keyPair, err := vw.hdWallet.GenerateKeyPair(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair %w", err)
	}
	return &VegaWallet{
		VegaWalletPrivate: &secrets.VegaWalletPrivate{
			Id:             keyPair.Name(),
			PublicKey:      keyPair.PublicKey(),
			PrivateKey:     keyPair.PrivateKey(),
			RecoveryPhrase: "",
		},
		hdWallet: vw.hdWallet,
		keyPair:  keyPair,
	}, nil
}
