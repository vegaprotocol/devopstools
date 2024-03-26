package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/vegaprotocol/devopstools/ethutils"
	"github.com/vegaprotocol/devopstools/smartcontracts/erc20token"
	"github.com/vegaprotocol/devopstools/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"go.uber.org/zap"
)

func (c *ChainClient) DepositERC20AssetToWhale(ctx context.Context, partyID string, assetContractAddress string, requiredAmount *types.Amount) error {
	flowID := rand.Int()

	token, err := erc20token.NewERC20Token(c.client, assetContractAddress, erc20token.ERC20TokenBase)
	if err != nil {
		return fmt.Errorf("could not initialize ERC20 token contract (%s): %w", assetContractAddress, err)
	}

	requiredAmountAsSubUnit := requiredAmount.AsSubUnit()

	minterHexAddress := c.minterWallet.Address.Hex()
	bridgeHexAddress := c.collateralBridge.Address.Hex()

	c.logger.Debug("Retrieving minter's balance...", zap.String("address", minterHexAddress))
	balanceAsSubUnit, err := token.BalanceOf(&bind.CallOpts{}, c.minterWallet.Address)
	if err != nil {
		return fmt.Errorf("failed to retrieve balance from minter's wallet %s: %w", minterHexAddress, err)
	}
	c.logger.Debug("Minter's balance retrieved", zap.String("balance-su", balanceAsSubUnit.String()))

	if balanceAsSubUnit.Cmp(requiredAmountAsSubUnit) > -1 {
		c.logger.Debug("No minting required",
			zap.Int("flow", flowID),
			zap.String("asset-contract-address", assetContractAddress),
			zap.String("minter", minterHexAddress),
			zap.String("current-balance-su", balanceAsSubUnit.String()),
			zap.String("required-balance-su", requiredAmountAsSubUnit.String()),
		)
	}

	amountToMintAsSubUnit := new(big.Int).Sub(requiredAmountAsSubUnit, balanceAsSubUnit)
	c.logger.Info("Minting required",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractAddress),
		zap.String("minter", minterHexAddress),
		zap.String("current-balance-su", balanceAsSubUnit.String()),
		zap.String("required-balance-su", requiredAmountAsSubUnit.String()),
		zap.String("amount-to-mint-su", amountToMintAsSubUnit.String()),
	)

	c.logger.Debug("Minting...",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractAddress),
		zap.String("minter", minterHexAddress),
		zap.String("amount-to-mint-su", amountToMintAsSubUnit.String()),
	)

	mintTx, err := token.Mint(c.minterWallet.GetTransactOpts(ctx), c.minterWallet.Address, amountToMintAsSubUnit)
	if err != nil {
		return fmt.Errorf("could not send transaction to mint wallet %s: %w", minterHexAddress, err)
	}

	if err := WaitForTransaction(ctx, c.client, mintTx, time.Minute); err != nil {
		return fmt.Errorf("failed to mint wallet %s: %w", minterHexAddress, err)
	}

	c.logger.Info("Minting successful",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractAddress),
		zap.String("minter", minterHexAddress),
		zap.String("amount-minted-su", amountToMintAsSubUnit.String()),
	)

	c.logger.Debug("Retrieving bridge allowance...",
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
	)
	allowanceAsSubUnit, err := token.Allowance(&bind.CallOpts{}, c.minterWallet.Address, c.collateralBridge.Address)
	if err != nil {
		return fmt.Errorf("failed to retrieve allowance: %w", err)
	}
	c.logger.Debug("Bridge's allowance retrieved", zap.String("allowance-su", allowanceAsSubUnit.String()))

	if allowanceAsSubUnit.Cmp(requiredAmountAsSubUnit) > -1 {
		c.logger.Debug("No allowance increase required",
			zap.Int("flow", flowID),
			zap.String("asset-contract-address", assetContractAddress),
			zap.String("minter", minterHexAddress),
			zap.String("bridge", bridgeHexAddress),
			zap.String("current-allowance-su", allowanceAsSubUnit.String()),
			zap.String("required-allowance-su", requiredAmountAsSubUnit.String()),
		)
	}

	allowanceToIncrease := new(big.Int).Sub(requiredAmountAsSubUnit, allowanceAsSubUnit)
	c.logger.Info("Allowance increase required",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractAddress),
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
		zap.String("current-allowance-su", allowanceAsSubUnit.String()),
		zap.String("required-allowance-su", requiredAmountAsSubUnit.String()),
		zap.String("increase-by-su", allowanceToIncrease.String()),
	)

	c.logger.Debug("Increasing allowance...",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractAddress),
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
		zap.String("increase-by-su", allowanceToIncrease.String()),
	)

	allowanceTx, err := token.IncreaseAllowance(c.minterWallet.GetTransactOpts(ctx), c.collateralBridge.Address, allowanceToIncrease)
	if err != nil {
		return fmt.Errorf("could not send transaction to increase bridge's allowance %s: %w", bridgeHexAddress, err)
	}

	if err := WaitForTransaction(ctx, c.client, allowanceTx, time.Minute); err != nil {
		return fmt.Errorf("failed to increase bridge's allowance %s: %w", bridgeHexAddress, err)
	}

	c.logger.Info("Allowance increase successful",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractAddress),
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
		zap.String("increased-by-su", allowanceToIncrease.String()),
	)

	whaleKeyB32, err := ethutils.VegaPubKeyToByte32(partyID)
	if err != nil {
		return fmt.Errorf("could not convert party ID to byte32: %w", err)
	}

	c.logger.Debug("Depositing asset...",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractAddress),
		zap.String("party-id", partyID),
		zap.String("amount", requiredAmount.String()),
	)

	depositTx, err := c.collateralBridge.DepositAsset(c.minterWallet.GetTransactOpts(ctx), token.Address, requiredAmountAsSubUnit, whaleKeyB32)
	if err != nil {
		return fmt.Errorf("could not send transaction to deposit asset %q to party ID (%s): %w", assetContractAddress, partyID, err)
	}

	if err := WaitForTransaction(ctx, c.client, depositTx, time.Minute); err != nil {
		return fmt.Errorf("failed to deposit asset %q to party ID (%s): %w", assetContractAddress, partyID, err)
	}

	c.logger.Info("Asset deposit successful",
		zap.Int("flow", flowID),
		zap.String("asset-contract-address", assetContractAddress),
		zap.String("minter", minterHexAddress),
		zap.String("bridge", bridgeHexAddress),
		zap.String("increased-by-su", allowanceToIncrease.String()),
	)

	return nil
}
