package tools

import (
	"math/big"
)

func AsIntStringFromFloat(amount *big.Float) string {
	return amount.Text('f', 0)
}
