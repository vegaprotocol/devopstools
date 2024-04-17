package erc20token

import (
	"fmt"
	"math/big"

	"github.com/vegaprotocol/devopstools/ethutils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Info struct {
	Address         string
	TotalSupply     *big.Float
	Name            string
	Symbol          string
	Decimals        uint8
	BurnEnabled     bool
	FaucetAmount    *big.Float
	FaucetCallLimit *big.Int
}

func (t *ERC20Token) GetInfo() (result Info, err error) {
	var (
		totalSupply  *big.Int
		faucetAmount *big.Int
	)
	result.Address = t.Address.Hex()

	result.BurnEnabled, err = t.BurnEnabled(&bind.CallOpts{})
	if err != nil {
		err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
		return
	}
	faucetAmount, err = t.FaucetAmount(&bind.CallOpts{})
	if err != nil {
		err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
		return
	}
	result.FaucetAmount = ethutils.TokenToFullTokens(faucetAmount, result.Decimals)
	result.FaucetCallLimit, err = t.FaucetCallLimit(&bind.CallOpts{})
	if err != nil {
		err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
		return
	}

	totalSupply, err = t.TotalSupply(&bind.CallOpts{})
	if err != nil {
		err = fmt.Errorf("failed to get info about %s, %w", result.Name, err)
		return
	}
	result.TotalSupply = ethutils.TokenToFullTokens(totalSupply, result.Decimals)

	return
}
