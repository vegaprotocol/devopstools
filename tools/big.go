package tools

import (
	"fmt"
	"math/big"
	"strings"
)

func AsIntStringFromFloat(amount *big.Float) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", amount), "0"), ".")
}
