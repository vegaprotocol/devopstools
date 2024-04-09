package types

import (
	"math/big"
	"strings"
)

// Amount wraps an amount of a given asset to easily switch between main and sub-unit
// representation.
type Amount struct {
	amount   *big.Float
	decimals uint64
}

func (a *Amount) Cmp(other *Amount) int {
	return a.amount.Cmp(other.amount)
}

func (a *Amount) AsSubUnit() *big.Int {
	return FloatAsSubUnit(a.amount, a.decimals)
}

func (a *Amount) AsMainUnit() *big.Float {
	return new(big.Float).Copy(a.amount)
}

func (a *Amount) Add(toAdd *big.Float) {
	a.amount = big.NewFloat(0).Add(a.amount, toAdd)
}

func (a *Amount) Sub(toSub *big.Float) {
	a.amount = big.NewFloat(0).Add(a.amount, toSub)
}

func (a *Amount) Mul(toMul *big.Float) {
	a.amount = big.NewFloat(0).Mul(a.amount, toMul)
}

func (a *Amount) String() string {
	return a.amount.String()
}

func (a *Amount) Copy() *Amount {
	return &Amount{
		amount:   new(big.Float).Copy(a.amount),
		decimals: a.decimals,
	}
}

func NewAmount(decimals uint64) *Amount {
	return &Amount{
		amount:   big.NewFloat(0),
		decimals: decimals,
	}
}

func NewAmountFromMainUnit(value *big.Float, decimals uint64) *Amount {
	return &Amount{
		amount:   value,
		decimals: decimals,
	}
}

func NewAmountFromSubUnit(value *big.Int, decimals uint64) *Amount {
	return &Amount{
		amount:   AsMainUint(value, decimals),
		decimals: decimals,
	}
}

func AsMainUint(value *big.Int, decimals uint64) *big.Float {
	valueStr := value.String()

	wholeNumber := "0"
	if len(valueStr) > int(decimals) {
		wholeNumber = valueStr[:len(valueStr)-int(decimals)]
	}

	var fractionalPart string
	if len(valueStr) < int(decimals) {
		fractionalPart = strings.Repeat("0", int(decimals)-len(valueStr)) + valueStr
	} else {
		fractionalPart = valueStr[len(valueStr)-int(decimals):]
	}

	amount, _ := new(big.Float).SetString(wholeNumber + "." + fractionalPart)
	return amount
}

func FloatAsSubUnit(amount *big.Float, decimals uint64) *big.Int {
	wholeNumber, _ := amount.Int(nil)

	segments := strings.Split(amount.String(), ".")
	if len(segments) == 1 {
		// There is no decimals.
		return IntAsSubUnit(wholeNumber, decimals)
	}

	fractionalPart := segments[1]
	if len(fractionalPart) < int(decimals) {
		// We pad with zeros the missing decimal.
		fractionalPart += strings.Repeat("0", int(decimals)-len(fractionalPart))
	} else {
		// We remove extra decimals.
		fractionalPart = fractionalPart[:int(decimals)]
	}

	amountAsSubUnit, _ := new(big.Int).SetString(wholeNumber.String()+fractionalPart, 10)
	return amountAsSubUnit
}

// IntAsSubUnit return an asset's amount padded with the decimals.
// Given "12" with a decimal place of 2, it returns "1200".
func IntAsSubUnit(amount *big.Int, decimal uint64) *big.Int {
	decimalMultiplier := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	return big.NewInt(0).Mul(amount, decimalMultiplier)
}
