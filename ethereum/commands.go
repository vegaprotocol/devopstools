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

var (
	ErrNoSignerFound         = errors.New("no signer found")
	ErrStakingBridgeDisabled = errors.New("staking bridge is disabled for this client")
)

func (c *ChainClient) StakeVegaTokenFromMinter(ctx context.Context, stakes map[string]*types.Amount) error {
	if c.stakingBridge == nil {
		return ErrStakingBridgeDisabled
	}

	stakingTokenAddress, err := c.stakingBridge.StakingToken(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("could not retrieve staking token: %w", err)
	}

	stakingTokenHexAddress := stakingTokenAddress.Hex()

	logger := c.logger.With(
		zap.Int("flow", rand.Int()),
		zap.String("asset-contract-address", stakingTokenHexAddress),
		zap.String("minter", c.minterWallet.Address.Hex()),
		zap.String("bridge", c.stakingBridge.Address.Hex()),
	)

	requiredAmountAsSubUnit := big.NewInt(0)
	for _, amount := range stakes {
		requiredAmountAsSubUnit.Add(requiredAmountAsSubUnit, amount.AsSubUnit())
	}

	token, err := erc20token.NewERC20Token(c.client, stakingTokenHexAddress)
	if err != nil {
		return fmt.Errorf("could not initialize ERC20 token contract client (%s): %w", stakingTokenHexAddress, err)
	}

	if err := c.mintWallet(ctx, c.minterWallet, token, requiredAmountAsSubUnit, logger); err != nil {
		return err
	}

	txs := map[string]*ethtypes.Transaction{}
	for partyID, amount := range stakes {
		logger.Debug("Staking Vega token...",
			zap.String("party-id", partyID),
			zap.String("amount", amount.String()),
		)

		tx, err := c.stakingBridge.Stake(c.minterWallet.GetTransactOpts(ctx), amount.AsSubUnit(), partyID)
		if err != nil {
			return fmt.Errorf("could not send transaction to stake Vega token %q to party ID (%s): %w", stakingTokenHexAddress, partyID, err)
		}
		txs[partyID] = tx
	}

	for partyID, tx := range txs {
		logger.Debug("Waiting for Ethereum transaction for party to complete...",
			zap.String("party-id", partyID),
			zap.String("tx-hash", tx.Hash().Hex()),
		)
		if err := WaitForTransaction(ctx, c.client, tx, time.Minute); err != nil {
			return fmt.Errorf("failed to stake Vega token %q to party ID (%s): %w", stakingTokenHexAddress, partyID, err)
		}
		logger.Info("Vega token staking successful for party",
			zap.String("party-id", partyID),
			zap.String("tx-hash", tx.Hash().Hex()),
		)
	}
	return nil
}

func (c *ChainClient) DepositERC20AssetFromMinter(ctx context.Context, assetContractHexAddress string, deposits map[string]*types.Amount) error {
	logger := c.logger.With(
		zap.Int("flow", rand.Int()),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", c.minterWallet.Address.Hex()),
		zap.String("bridge", c.collateralBridge.Address.Hex()),
	)

	return c.depositERC20TokenFromWallet(ctx, c.minterWallet, assetContractHexAddress, deposits, logger)
}

func (c *ChainClient) DepositERC20AssetFromAddress(ctx context.Context, minterPrivateHexKey string, assetContractHexAddress string, deposits map[string]*types.Amount) error {
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

	logger := c.logger.With(
		zap.Int("flow", rand.Int()),
		zap.String("asset-contract-address", assetContractHexAddress),
		zap.String("minter", minterWallet.Address.Hex()),
		zap.String("bridge", c.collateralBridge.Address.Hex()),
	)

	return c.depositERC20TokenFromWallet(ctx, minterWallet, assetContractHexAddress, deposits, logger)
}

func (c *ChainClient) depositERC20TokenFromWallet(ctx context.Context, minterWallet *Wallet, assetContractHexAddress string, deposits map[string]*types.Amount, logger *zap.Logger) error {
	requiredAmountAsSubUnit := big.NewInt(0)
	for _, amount := range deposits {
		requiredAmountAsSubUnit.Add(requiredAmountAsSubUnit, amount.AsSubUnit())
	}

	token, err := erc20token.NewERC20Token(c.client, assetContractHexAddress)
	if err != nil {
		return fmt.Errorf("could not initialize ERC20 token contract client (%s): %w", assetContractHexAddress, err)
	}

	if err := c.mintWallet(ctx, minterWallet, token, requiredAmountAsSubUnit, logger); err != nil {
		return err
	}

	txs := map[string]*ethtypes.Transaction{}
	for partyID, amount := range deposits {
		partyKeyB32, err := ethutils.VegaPubKeyToByte32(partyID)
		if err != nil {
			return fmt.Errorf("could not convert party ID to byte32: %w", err)
		}

		logger.Debug("Depositing asset...",
			zap.String("party-id", partyID),
			zap.String("amount", amount.String()),
		)

		tx, err := c.collateralBridge.DepositAsset(minterWallet.GetTransactOpts(ctx), token.Address, amount.AsSubUnit(), partyKeyB32)
		if err != nil {
			return fmt.Errorf("could not send transaction to deposit asset %q to party ID (%s): %w", assetContractHexAddress, partyID, err)
		}
		txs[partyID] = tx
	}

	for partyID, tx := range txs {
		logger.Debug("Waiting for Ethereum transaction for party to complete...",
			zap.String("party-id", partyID),
			zap.String("tx-hash", tx.Hash().Hex()),
		)
		if err := WaitForTransaction(ctx, c.client, tx, time.Minute); err != nil {
			return fmt.Errorf("failed to deposit asset %q to party ID (%s): %w", assetContractHexAddress, partyID, err)
		}
		logger.Info("Asset deposit successful for party",
			zap.String("party-id", partyID),
			zap.String("tx-hash", tx.Hash().Hex()),
		)
	}
	return nil
}

func (c *ChainClient) mintWallet(ctx context.Context, minterWallet *Wallet, token *erc20token.ERC20Token, requiredAmountAsSubUnit *big.Int, logger *zap.Logger) error {
	logger.Debug("Retrieving wallet's balance...")
	balanceAsSubUnit, err := token.BalanceOf(&bind.CallOpts{Context: ctx}, minterWallet.Address)
	if err != nil {
		return fmt.Errorf("failed to retrieve balance from minter's wallet: %w", err)
	}
	logger.Debug("Minter's balance retrieved", zap.String("balance-su", balanceAsSubUnit.String()))

	if balanceAsSubUnit.Cmp(requiredAmountAsSubUnit) > -1 {
		logger.Debug("No minting required",
			zap.String("current-balance-su", balanceAsSubUnit.String()),
			zap.String("required-balance-su", requiredAmountAsSubUnit.String()),
		)
	}

	amountToMintAsSubUnit := new(big.Int).Sub(requiredAmountAsSubUnit, balanceAsSubUnit)
	logger.Info("Minting required",
		zap.String("current-balance-su", balanceAsSubUnit.String()),
		zap.String("required-balance-su", requiredAmountAsSubUnit.String()),
		zap.String("amount-to-mint-su", amountToMintAsSubUnit.String()),
	)

	logger.Debug("Minting...", zap.String("amount-to-mint-su", amountToMintAsSubUnit.String()))

	mintTx, err := token.Mint(minterWallet.GetTransactOpts(ctx), minterWallet.Address, amountToMintAsSubUnit)
	if err != nil {
		return fmt.Errorf("could not send transaction to mint wallet: %w", err)
	}

	if err := WaitForTransaction(ctx, c.client, mintTx, time.Minute); err != nil {
		return fmt.Errorf("failed to mint wallet: %w", err)
	}

	logger.Info("Minting successful", zap.String("amount-minted-su", amountToMintAsSubUnit.String()))

	logger.Debug("Retrieving bridge allowance...")
	allowanceAsSubUnit, err := token.Allowance(&bind.CallOpts{Context: ctx}, minterWallet.Address, c.collateralBridge.Address)
	if err != nil {
		return fmt.Errorf("failed to retrieve allowance: %w", err)
	}
	logger.Debug("Bridge's allowance retrieved", zap.String("allowance-su", allowanceAsSubUnit.String()))

	if allowanceAsSubUnit.Cmp(requiredAmountAsSubUnit) > -1 {
		logger.Debug("No allowance increase required",
			zap.String("current-allowance-su", allowanceAsSubUnit.String()),
			zap.String("required-allowance-su", requiredAmountAsSubUnit.String()),
		)
	}

	allowanceToIncrease := new(big.Int).Sub(requiredAmountAsSubUnit, allowanceAsSubUnit)
	logger.Info("Allowance increase required",
		zap.String("current-allowance-su", allowanceAsSubUnit.String()),
		zap.String("required-allowance-su", requiredAmountAsSubUnit.String()),
		zap.String("increase-by-su", allowanceToIncrease.String()),
	)

	logger.Debug("Increasing allowance...", zap.String("increase-by-su", allowanceToIncrease.String()))

	allowanceTx, err := token.IncreaseAllowance(minterWallet.GetTransactOpts(ctx), c.collateralBridge.Address, allowanceToIncrease)
	if err != nil {
		return fmt.Errorf("could not send transaction to increase bridge's allowance: %w", err)
	}

	if err := WaitForTransaction(ctx, c.client, allowanceTx, time.Minute); err != nil {
		return fmt.Errorf("failed to increase bridge's allowance: %w", err)
	}

	logger.Info("Allowance increase successful", zap.String("increased-by-su", allowanceToIncrease.String()))

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

	var signatures []byte
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