package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20token"
	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"

	"code.vegaprotocol.io/vega/libs/num"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
)

var ErrNoSignerFound = errors.New("no signer found")

func (c *ChainClient) DepositERC20AssetFromMinter(ctx context.Context, assetContractHexAddress string, partyID string, requiredAmount *types.Amount) error {
	flowID := rand.Int()

	requiredAmountAsSubUnit := requiredAmount.AsSubUnit()

	minterHexAddress := c.minterWallet.Address.Hex()
	bridgeHexAddress := c.collateralBridge.Address.Hex()

	token, err := erc20token.NewERC20Token(c.client, assetContractHexAddress, erc20token.ERC20TokenBase)
	if err != nil {
		return fmt.Errorf("could not initialize ERC20 token contract client (%s): %w", assetContractHexAddress, err)
	}

	if err := c.mintWallet(ctx, flowID, c.minterWallet, token, requiredAmountAsSubUnit, assetContractHexAddress, minterHexAddress, bridgeHexAddress); err != nil {
		return err
	}

	whaleKeyB32, err := ethutils.VegaPubKeyToByte32(partyID)
	if err != nil {
		return fmt.Errorf("could not convert party ID to byte32: %w", err)
	}

	c.logger.Debug("Depositing asset...",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("party-id", partyID),
		zap.String("amount", requiredAmount.String()),
	)

	depositTx, err := c.collateralBridge.DepositAsset(c.minterWallet.GetTransactOpts(ctx), token.Address, requiredAmountAsSubUnit, whaleKeyB32)
	if err != nil {
		return fmt.Errorf("could not send transaction to deposit asset %q to party ID (%s): %w", assetContractHexAddress, partyID, err)
	}

	if err := WaitForTransaction(ctx, c.client, depositTx, time.Minute); err != nil {
		return fmt.Errorf("failed to deposit asset %q to party ID (%s): %w", assetContractHexAddress, partyID, err)
	}

	c.logger.Info("Asset deposit successful",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
		zap.String("required-amount", requiredAmountAsSubUnit.String()),
	)

	return nil
}

func (c *ChainClient) DepositERC20AssetFromAddress(ctx context.Context, minterPrivateHexKey string, assetContractHexAddress string, deposits map[string]*types.Amount) error {
	flowID := rand.Int()

	minterWallet, err := tools.RetryReturn(6, 10*time.Second, func() (*Wallet, error) {
		w, err := NewWallet(ctx, c.client, minterPrivateHexKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Ethereum wallet: %w", err)
		}
		return w, nil
	})
	if err != nil {
		return err
	}

	minterHexAddress := minterWallet.Address.Hex()
	bridgeHexAddress := c.collateralBridge.Address.Hex()

	requiredAmountAsSubUnit := big.NewInt(0)
	for _, amount := range deposits {
		requiredAmountAsSubUnit.Add(requiredAmountAsSubUnit, amount.AsSubUnit())
	}

	token, err := erc20token.NewERC20Token(c.client, assetContractHexAddress, erc20token.ERC20TokenBase)
	if err != nil {
		return fmt.Errorf("could not initialize ERC20 token contract client (%s): %w", assetContractHexAddress, err)
	}

	if err := c.mintWallet(ctx, flowID, minterWallet, token, requiredAmountAsSubUnit, assetContractHexAddress, minterHexAddress, bridgeHexAddress); err != nil {
		return err
	}

	depositTxs := map[string]*ethtypes.Transaction{}
	for partyID, amount := range deposits {
		partyKeyB32, err := ethutils.VegaPubKeyToByte32(partyID)
		if err != nil {
			return fmt.Errorf("could not convert party ID to byte32: %w", err)
		}

		c.logger.Debug("Depositing asset...",
			zap.Int("flow", flowID),
			zap.String("asset-contract-address", assetContractHexAddress),
			zap.String("party-id", partyID),
			zap.String("amount", amount.String()),
		)

		depositTx, err := c.collateralBridge.DepositAsset(minterWallet.GetTransactOpts(ctx), token.Address, amount.AsSubUnit(), partyKeyB32)
		if err != nil {
			return fmt.Errorf("could not send transaction to deposit asset %q to party ID (%s): %w", assetContractHexAddress, partyID, err)
		}
		depositTxs[partyID] = depositTx
	}

	for partyID, depositTx := range depositTxs {
		if err := WaitForTransaction(ctx, c.client, depositTx, time.Minute); err != nil {
			return fmt.Errorf("failed to deposit asset %q to party ID (%s): %w", assetContractHexAddress, partyID, err)
		}
	}

	c.logger.Info("Asset deposit successful",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
	)

	return nil
}

func (c *ChainClient) mintWallet(ctx context.Context, flowID int, minterWallet *Wallet, token *erc20token.ERC20Token, requiredAmountAsSubUnit *big.Int, assetContractHexAddress string, minterHexAddress string, bridgeHexAddress string) error {
	c.logger.Debug("Retrieving wallet's balance...", zap.Int("flow", flowID), zap.String("address", minterHexAddress))
	balanceAsSubUnit, err := token.BalanceOf(&bind.CallOpts{}, c.minterWallet.Address)
	if err != nil {
		return fmt.Errorf("failed to retrieve balance from minter's wallet %s: %w", minterHexAddress, err)
	}
	c.logger.Debug("Minter's balance retrieved", zap.Int("flow", flowID), zap.String("balance-su", balanceAsSubUnit.String()))

	if balanceAsSubUnit.Cmp(requiredAmountAsSubUnit) > -1 {
		c.logger.Debug("No minting required",
			zap.Int("flow", flowID),
			zap.String("asset-contract-address", assetContractHexAddress),
			zap.String("minter", minterHexAddress),
			zap.String("current-balance-su", balanceAsSubUnit.String()),
			zap.String("required-balance-su", requiredAmountAsSubUnit.String()),
		)
	}

	amountToMintAsSubUnit := new(big.Int).Sub(requiredAmountAsSubUnit, balanceAsSubUnit)
	c.logger.Info("Minting required",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", minterHexAddress),
		zap.String("current-balance-su", balanceAsSubUnit.String()),
		zap.String("required-balance-su", requiredAmountAsSubUnit.String()),
		zap.String("amount-to-mint-su", amountToMintAsSubUnit.String()),
	)

	c.logger.Debug("Minting...",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", minterHexAddress),
		zap.String("amount-to-mint-su", amountToMintAsSubUnit.String()),
	)

	mintTx, err := token.Mint(minterWallet.GetTransactOpts(ctx), c.minterWallet.Address, amountToMintAsSubUnit)
	if err != nil {
		return fmt.Errorf("could not send transaction to mint wallet %s: %w", minterHexAddress, err)
	}

	if err := WaitForTransaction(ctx, c.client, mintTx, time.Minute); err != nil {
		return fmt.Errorf("failed to mint wallet %s: %w", minterHexAddress, err)
	}

	c.logger.Info("Minting successful",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", minterHexAddress),
		zap.String("amount-minted-su", amountToMintAsSubUnit.String()),
	)

	c.logger.Debug("Retrieving bridge allowance...",
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
	)
	allowanceAsSubUnit, err := token.Allowance(&bind.CallOpts{}, minterWallet.Address, c.collateralBridge.Address)
	if err != nil {
		return fmt.Errorf("failed to retrieve allowance: %w", err)
	}
	c.logger.Debug("Bridge's allowance retrieved", zap.String("allowance-su", allowanceAsSubUnit.String()))

	if allowanceAsSubUnit.Cmp(requiredAmountAsSubUnit) > -1 {
		c.logger.Debug("No allowance increase required",
			zap.Int("flow", flowID),
			zap.String("asset-contract-address", assetContractHexAddress),
			zap.String("minter", minterHexAddress),
			zap.String("bridge", bridgeHexAddress),
			zap.String("current-allowance-su", allowanceAsSubUnit.String()),
			zap.String("required-allowance-su", requiredAmountAsSubUnit.String()),
		)
	}

	allowanceToIncrease := new(big.Int).Sub(requiredAmountAsSubUnit, allowanceAsSubUnit)
	c.logger.Info("Allowance increase required",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
		zap.String("current-allowance-su", allowanceAsSubUnit.String()),
		zap.String("required-allowance-su", requiredAmountAsSubUnit.String()),
		zap.String("increase-by-su", allowanceToIncrease.String()),
	)

	c.logger.Debug("Increasing allowance...",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
		zap.String("increase-by-su", allowanceToIncrease.String()),
	)

	allowanceTx, err := token.IncreaseAllowance(minterWallet.GetTransactOpts(ctx), c.collateralBridge.Address, allowanceToIncrease)
	if err != nil {
		return fmt.Errorf("could not send transaction to increase bridge's allowance %s: %w", bridgeHexAddress, err)
	}

	if err := WaitForTransaction(ctx, c.client, allowanceTx, time.Minute); err != nil {
		return fmt.Errorf("failed to increase bridge's allowance %s: %w", bridgeHexAddress, err)
	}

	c.logger.Info("Allowance increase successful",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
		zap.String("increased-by-su", allowanceToIncrease.String()),
	)

	return nil
}

func (c *ChainClient) Signers(ctx context.Context) ([]common.Address, error) {
	signersAddresses, err := c.multisigControl.GetSigners(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get signers addresses from multisig control: %w", err)
	}

	if len(signersAddresses) == 0 {
		return nil, ErrNoSignerFound
	}

	return signersAddresses, nil
}

func (c *ChainClient) ListAsset(ctx context.Context, signers []*Wallet, assetID, assetHexAddress string, lifetimeLimit, withdrawThreshold *big.Int) error {
	assetIDB32, err := ethutils.VegaPubKeyToByte32(assetID)
	if err != nil {
		return fmt.Errorf("could not convert asset ID to byte32: %w", err)
	}

	assetAddress := common.HexToAddress(assetHexAddress)
	nonce := num.NewUint(ethutils.TimestampNonce()).BigInt()

	msg, err := buildListAssetMsg(c.collateralBridge.Address, assetAddress, assetIDB32, lifetimeLimit, withdrawThreshold, nonce)
	if err != nil {
		return fmt.Errorf("could not generate list an asset message for signing: %w", err)
	}

	c.logger.Debug("Generating multisig...", zap.ByteString("message", msg))
	signatures, err := generateMultisig(signers, msg)
	if err != nil {
		return fmt.Errorf("could not generate signatures to list an asset: %w", err)
	}
	c.logger.Debug("Multisig generated", zap.ByteString("signatures", signatures))

	c.logger.Debug("Listing asset...",
		zap.String("asset-id", assetID),
		zap.String("asset-contract-address", assetHexAddress),
	)

	tx, err := c.collateralBridge.ListAsset(c.minterWallet.GetTransactOpts(ctx), assetAddress, assetIDB32, lifetimeLimit, withdrawThreshold, nonce, signatures)
	if err != nil {
		return fmt.Errorf("listing asset failed: %w", err)
	}

	if err := WaitForTransaction(ctx, c.client, tx, time.Minute); err != nil {
		return fmt.Errorf("transaction to list asset failed: %w", err)
	}

	c.logger.Info("Asset listing successful",
		zap.String("asset-id", assetID),
		zap.String("asset-contract-address", assetHexAddress),
	)

	return nil
}

func generateMultisig(signers []*Wallet, msg []byte) ([]byte, error) {
	hash := crypto.Keccak256(msg)

	signatures := []byte{}
	for _, signerKey := range signers {
		signature, err := signerKey.Sign(hash)
		if err != nil {
			return nil, fmt.Errorf("failed to sign message hash: %w", err)
		}

		signatures = append(signatures, signature...)
	}
	return signatures, nil
}

func buildListAssetMsg(bridgeAddr common.Address, tokenAddress common.Address, assetIDB32 [32]byte, lifetimeLimit *big.Int, withdrawThreshold *big.Int, nonce *big.Int) ([]byte, error) {
	typAddr, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, err
	}
	typBytes32, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return nil, err
	}
	typString, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, err
	}
	typU256, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}

	args := abi.Arguments([]abi.Argument{
		{
			Name: "address",
			Type: typAddr,
		},
		{
			Name: "vega_asset_id",
			Type: typBytes32,
		},
		{
			Name: "lifetime_limit",
			Type: typU256,
		},
		{
			Name: "withdraw_treshold",
			Type: typU256,
		},
		{
			Name: "nonce",
			Type: typU256,
		},
		{
			Name: "func_name",
			Type: typString,
		},
	})

	buf, err := args.Pack([]interface{}{
		tokenAddress,
		assetIDB32,
		lifetimeLimit,
		withdrawThreshold,
		nonce,
		"list_asset",
	}...)
	if err != nil {
		return nil, fmt.Errorf("couldn't pack abi message: %w", err)
	}

	msg, err := packBufAndSubmitter(buf, bridgeAddr)
	if err != nil {
		return nil, fmt.Errorf("couldn't pack abi message: %w", err)
	}

	return msg, nil
}

func packBufAndSubmitter(buf []byte, submitter common.Address) ([]byte, error) {
	typBytes, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return nil, err
	}
	typAddr, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, err
	}

	args2 := abi.Arguments([]abi.Argument{
		{
			Name: "bytes",
			Type: typBytes,
		},
		{
			Name: "address",
			Type: typAddr,
		},
	})

	return args2.Pack(buf, submitter)
}
