package tools

import (
	"math/big"
	"strconv"
)

func AsIntStringFromFloat(amount *big.Float) string {
	return amount.Text('f', 0)
}

func StrToIntOrDefault(s string, defaultValue int) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return res
}
