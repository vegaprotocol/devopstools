package erc20token

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type ERC20TokenInfo struct {
	Address         string
	TotalSupply     *big.Int
	Name            string
	Symbol          string
	Decimals        uint8
	BurnEnabled     bool
	FaucetAmount    *big.Int
	FaucetCallLimit *big.Int
}

func (t *ERC20Token) GetInfo() (result ERC20TokenInfo, err error) {
	result.Address = t.Address.Hex()

	if t.Version != ERC20TokenMinimal {
		result.Name, err = t.Name(&bind.CallOpts{})
		if err != nil {
			err = fmt.Errorf("failed to get info about %s, %w", t.Address, err)
			return
		}
		result.Symbol, err = t.Symbol(&bind.CallOpts{})
		if err != nil {
			err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
			return
		}
		result.Decimals, err = t.Decimals(&bind.CallOpts{})
		if err != nil {
			err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
			return
		}
	}
	if t.Version == ERC20TokenBase {
		result.BurnEnabled, err = t.BurnEnabled(&bind.CallOpts{})
		if err != nil {
			err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
			return
		}
		result.FaucetAmount, err = t.FaucetAmount(&bind.CallOpts{})
		if err != nil {
			err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
			return
		}
		result.FaucetCallLimit, err = t.FaucetCallLimit(&bind.CallOpts{})
		if err != nil {
			err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
			return
		}
	}

	result.TotalSupply, err = t.TotalSupply(&bind.CallOpts{})
	if err != nil {
		err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
		return
	}

	return
}
