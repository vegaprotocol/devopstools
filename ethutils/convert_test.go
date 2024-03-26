package ethutils_test

import (
	"math/big"
	"testing"

	"github.com/vegaprotocol/devopstools/ethutils"
)

func TestEtherToWei(t *testing.T) {
	tests := map[string]struct {
		amount         string
		decimals       uint8
		expectedResult string
	}{
		"regular": {amount: "3.14159", expectedResult: "3141590000000000000"},
		"small":   {amount: "0.000000000000014159", expectedResult: "14159"},
		"big":     {amount: "789223400", expectedResult: "789223400000000000000000000"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			amount := new(big.Float)
			_, ok := amount.SetString(tc.amount)
			if !ok {
				t.Fatalf("failed to convert string '%s' to big.Float", tc.amount)
			}
			result := ethutils.EtherToWei(amount)
			if result.String() != tc.expectedResult {
				t.Fatalf("result '%s' != expected result '%s'", result.String(), tc.expectedResult)
			}
		})
	}
}
