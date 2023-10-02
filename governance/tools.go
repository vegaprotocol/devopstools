package governance

import (
	"fmt"

	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"github.com/vegaprotocol/devopstools/vegaapi"
	"github.com/vegaprotocol/devopstools/wallet"
	"go.uber.org/zap"
)

func SubmitTx(
	description string,
	dataNodeClient vegaapi.DataNodeClient,
	proposerVegawallet *wallet.VegaWallet,
	logger *zap.Logger,
	walletTxReq *walletpb.SubmitTransactionRequest,
) error {
	lastBlockData, err := dataNodeClient.LastBlockData()
	if err != nil {
		return fmt.Errorf("failed to submit tx: %w", err)
	}

	// Sign + Proof of Work vegawallet Transaction request
	signedTx, err := proposerVegawallet.SignTxWithPoW(walletTxReq, lastBlockData)
	if err != nil {
		logger.Error("Failed to sign a trasnaction", zap.String("description", description),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", &walletTxReq), zap.Error(err))
		return err
	}

	// wrap in vega Transaction Request
	submitReq := &vegaapipb.SubmitTransactionRequest{
		Tx:   signedTx,
		Type: vegaapipb.SubmitTransactionRequest_TYPE_SYNC,
	}

	// Submit Transaction
	logger.Info("Submit transaction", zap.String("description", description),
		zap.String("proposer", proposerVegawallet.PublicKey))
	submitResponse, err := dataNodeClient.SubmitTransaction(submitReq)
	if err != nil {
		logger.Error("Failed to submit a trasnaction", zap.String("description", description),
			zap.String("proposer", proposerVegawallet.PublicKey),
			zap.Any("txReq", submitReq), zap.Error(err))
		return err
	}
	if !submitResponse.Success {
		logger.Error("Transaction submission response is not successful",
			zap.String("proposer", proposerVegawallet.PublicKey), zap.String("description", description),
			zap.Any("txReq", submitReq.String()), zap.String("response", fmt.Sprintf("%#v", submitResponse)))
		return err
	}
	logger.Info("Successful Submision of Transaction", zap.String("description", description),
		zap.String("proposer", proposerVegawallet.PublicKey), zap.String("txHash", submitResponse.TxHash),
		zap.Any("response", submitResponse))

	return nil
}
