package veganetworksmartcontracts

import (
	context2 "context"
	"fmt"
	"math/big"
	"time"

	"github.com/vegaprotocol/devopstools/ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func (m *VegaNetworkSmartContracts) TopUpStakeForOne(minterWallet *ethereum.EthWallet, partyPubKey string, amount *big.Int) error {
	return m.TopUpStake(minterWallet, map[string]*big.Int{partyPubKey: amount})
}

func (m *VegaNetworkSmartContracts) TopUpStake(minterWallet *ethereum.EthWallet, requiredStakeByParty map[string]*big.Int) error {
	vegaTokenInfo, err := m.VegaToken.GetInfo()
	if err != nil {
		return fmt.Errorf("failed to get info about Vega Token, %w", err)
	}
	var (
		missingStakeByParty = map[string]*big.Int{}
		totalMissingStake   = new(big.Int)
	)

	m.logger.Info("Top Up Stake",
		zap.String("ethWallet", minterWallet.Address.Hex()),
		zap.String("token", vegaTokenInfo.Name),
		zap.String("tokenAddress", vegaTokenInfo.Address),
		zap.String("stakingBridgeAddress", m.StakingBridge.Address.Hex()),
		zap.Int("partiesCount", len(requiredStakeByParty)),
	)

	//
	// Get missing Stake
	//
	for partyPubKey, partyRequiredStake := range requiredStakeByParty {
		partyCurrentStake, err := m.StakingBridge.GetStakeBalance(partyPubKey)
		if err != nil {
			return err
		}
		if partyCurrentStake.Cmp(partyRequiredStake) < 0 {
			missingStakeByParty[partyPubKey] = new(big.Int).Sub(partyRequiredStake, partyCurrentStake)
			totalMissingStake = totalMissingStake.Add(totalMissingStake, missingStakeByParty[partyPubKey])
		}
	}
	m.logger.Debug("Missing Stake by Parties", zap.String("totalMissingStake", totalMissingStake.String()), zap.Any("missingStakeByParty", missingStakeByParty))

	//
	// Mint + IncreaseAllowance for minterWallet
	//
	var (
		mintTx      *ethTypes.Transaction
		allowanceTx *ethTypes.Transaction
	)
	// Mint
	minterTokenBalance, err := m.VegaToken.BalanceOf(&bind.CallOpts{}, minterWallet.Address)
	if err != nil {
		return fmt.Errorf("failed to get balance of minterWallet %s, %w", minterWallet.Address, err)
	}
	if minterTokenBalance.Cmp(totalMissingStake) < 0 {
		diff := new(big.Int).Sub(totalMissingStake, minterTokenBalance)
		opts := minterWallet.GetTransactOpts(context2.Background())
		m.logger.Info("Mint", zap.String("token", vegaTokenInfo.Name), zap.String("ethWallet", minterWallet.Address.Hex()), zap.String("amount", diff.String()))
		m.logger.Debug("Mint", zap.String("balanceBefore", minterTokenBalance.String()), zap.String("requiredBalance", totalMissingStake.String()))

		mintTx, err = m.VegaToken.Mint(opts, minterWallet.Address, diff)
		if err != nil {
			return fmt.Errorf("failed to mint %s of %s (%s) for %s: %w", diff.String(),
				vegaTokenInfo.Name, vegaTokenInfo.Address, minterWallet.Address.Hex(), err)
		}
		m.logger.Info("Mint", zap.String("tx", mintTx.Hash().Hex()))
	} else {
		m.logger.Info("No need to mint", zap.String("token", vegaTokenInfo.Name), zap.String("ethWallet", minterWallet.Address.Hex()), zap.String("minterTokenBalance", minterTokenBalance.String()))
		m.logger.Debug("No need to mint", zap.String("balanceBefore", minterTokenBalance.String()), zap.String("requiredBalance", totalMissingStake.String()))
	}
	// Increase Allowance
	minterTokenAllowance, err := m.VegaToken.Allowance(&bind.CallOpts{}, minterWallet.Address, m.StakingBridge.Address)
	if err != nil {
		return fmt.Errorf("failed to get allowance for staking bridge %s, %w", m.StakingBridge.Address.Hex(), err)
	}
	if minterTokenAllowance.Cmp(totalMissingStake) < 0 {
		diff := new(big.Int).Sub(totalMissingStake, minterTokenAllowance)
		opts := minterWallet.GetTransactOpts(context2.Background())
		m.logger.Info("Increasing allowance", zap.String("token", vegaTokenInfo.Name), zap.String("ethWallet", minterWallet.Address.Hex()), zap.String("amount", diff.String()))
		m.logger.Debug("Increasing allowance", zap.String("allowanceBefore", minterTokenAllowance.String()), zap.String("requiredAllowance", totalMissingStake.String()))
		allowanceTx, err = m.VegaToken.IncreaseAllowance(opts, m.StakingBridge.Address, diff)
		if err != nil {
			return fmt.Errorf("failed to increase allowance: %w", err)
		}
		m.logger.Info("Increase allowance", zap.String("tx", allowanceTx.Hash().Hex()))
	} else {
		m.logger.Info("No need to increase allowance", zap.String("token", vegaTokenInfo.Name), zap.String("ethWallet", minterWallet.Address.Hex()), zap.String("minterTokenAllowance", minterTokenAllowance.String()))
		m.logger.Debug("No need to increase allowance", zap.String("balanceBefore", minterTokenAllowance.String()), zap.String("requiredBalance", totalMissingStake.String()))
	}
	// wait for Mint and IncreaseAllowance
	if mintTx != nil {
		m.logger.Sugar().Infof("Wait for mint transaction  (%s) ...", mintTx.Hash().Hex())
		if err = ethereum.WaitForTransaction(context2.Background(), m.EthClient, mintTx, time.Minute*2); err != nil {
			m.logger.Error("failed to mint", zap.String("token", vegaTokenInfo.Name), zap.Error(err))
			return fmt.Errorf("transaction failed to mints: %w", err)
		}
		m.logger.Sugar().Infoln("success")
	}
	if allowanceTx != nil {
		m.logger.Sugar().Infof("Wait for increase allowance transaction  (%s) ... ", allowanceTx.Hash().Hex())
		if err = ethereum.WaitForTransaction(context2.Background(), m.EthClient, allowanceTx, time.Minute); err != nil {
			m.logger.Sugar().Infoln("failed")
			m.logger.Error("failed to increase allowance", zap.String("ethWallet", minterWallet.Address.Hex()),
				zap.String("token", vegaTokenInfo.Name), zap.String("tokenAddress", vegaTokenInfo.Address), zap.Error(err))
			return fmt.Errorf("transaction failed to increase allowance: %w", err)
		}
		m.logger.Sugar().Infoln("success")
	}

	//
	// Stake to parties
	//
	var (
		stakeTxs          = map[string]*ethTypes.Transaction{}
		stakeFailureCount = 0
		stakeSuccessCount = 0
		stakeOKCount      = 0
	)
	for partyPubKey := range requiredStakeByParty {
		if partyMissingStake, ok := missingStakeByParty[partyPubKey]; !ok {
			stakeOKCount += 1
			m.logger.Info("No need to top up", zap.String("partyPubKey", partyPubKey))
		} else {
			opts := minterWallet.GetTransactOpts(context2.Background())
			tx, err := m.StakingBridge.Stake(opts, partyMissingStake, partyPubKey)
			if err != nil {
				stakeFailureCount += 1
				m.logger.Error("failed to stake", zap.String("partyPubKey", partyPubKey), zap.String("amount", partyMissingStake.String()),
					zap.Any("tx", tx), zap.String("ethWallet", minterWallet.Address.Hex()),
					zap.String("token", vegaTokenInfo.Name), zap.String("tokenAddress", vegaTokenInfo.Address), zap.Error(err))
			} else {
				stakeTxs[partyPubKey] = tx
			}
		}
	}
	// wait for transactions
	for partyPubKey, tx := range stakeTxs {
		if tx == nil {
			continue
		}
		m.logger.Sugar().Infof("Wait for stake for %s transaction  (%s) ... ", partyPubKey, tx.Hash().Hex())
		if err := ethereum.WaitForTransaction(context2.Background(), m.EthClient, tx, time.Minute*2); err != nil {
			stakeFailureCount += 1
			m.logger.Sugar().Infoln("failed")
			m.logger.Error("failed to stake", zap.String("partyPubKey", partyPubKey),
				zap.Any("tx", tx), zap.String("ethWallet", minterWallet.Address.Hex()),
				zap.String("token", vegaTokenInfo.Name), zap.String("tokenAddress", vegaTokenInfo.Address), zap.Error(err))
		} else {
			stakeSuccessCount += 1
			m.logger.Sugar().Infoln("success")
		}
	}
	m.logger.Info("Stake summary", zap.Int("success", stakeSuccessCount), zap.Int("failure", stakeFailureCount),
		zap.Int("already ok", stakeOKCount))

	return nil
}

func (m *VegaNetworkSmartContracts) RemoveStake(ethWallet *ethereum.EthWallet, partyPubKey string) error {
	fmt.Printf("---- m: %v, ethWallet %v", m, ethWallet)
	currentStake, err := m.StakingBridge.StakeBalance(&bind.CallOpts{}, ethWallet.Address, partyPubKey)
	if err != nil {
		return err
	}
	if currentStake.Cmp(big.NewInt(0)) > 0 {
		opts := ethWallet.GetTransactOpts(context2.Background())
		tx, err := m.StakingBridge.RemoveStake(opts, currentStake, partyPubKey)
		if err != nil {
			m.logger.Error("failed to remove stake", zap.String("ethWallet", ethWallet.Address.Hex()), zap.String("partyPubKey", partyPubKey),
				zap.String("amount", currentStake.String()), zap.Error(err))
			return fmt.Errorf("failed to remove stake %s by %s from %s , %w", currentStake.String(), ethWallet.Address.Hex(), partyPubKey, err)
		}
		m.logger.Sugar().Infof("Wait for remove stake %d by %s from %s transaction  (%s) ... ", currentStake.String(), ethWallet.Address.Hex(), partyPubKey, tx.Hash().Hex())
		if err := ethereum.WaitForTransaction(context2.Background(), m.EthClient, tx, time.Minute*2); err != nil {
			m.logger.Sugar().Infoln("failed")
			m.logger.Error("failed to remove stake", zap.String("ethWallet", ethWallet.Address.Hex()), zap.String("partyPubKey", partyPubKey),
				zap.String("amount", currentStake.String()), zap.Any("tx", tx), zap.Error(err))
		} else {
			m.logger.Sugar().Infoln("success")
		}
	} else {
		m.logger.Info("no stake to remove", zap.String("ethWallet", ethWallet.Address.Hex()), zap.String("partyPubKey", partyPubKey))
	}
	return nil
}
