package types_test

import (
	"math/big"
	"testing"

	"github.com/vegaprotocol/devopstools/types"

	"github.com/stretchr/testify/assert"
)

func TestAmount(t *testing.T) {
	a0 := types.NewAmount(4)
	assert.Equal(t, "0", a0.String())
	assert.Equal(t, "0", a0.AsSubUnit().String())

	a1 := types.NewAmountFromSubUnit(big.NewInt(1234), 2)
	assert.Equal(t, "12.34", a1.String())
	assert.Equal(t, "1234", a1.AsSubUnit().String())

	a2 := types.NewAmountFromSubUnit(big.NewInt(1234), 6)
	assert.Equal(t, "0.001234", a2.String())
	assert.Equal(t, "1234", a2.AsSubUnit().String())

	a3 := types.NewAmountFromMainUnit(big.NewFloat(12.34), 6)
	assert.Equal(t, "12.34", a3.String())
	assert.Equal(t, "12340000", a3.AsSubUnit().String())

	a1.Add(big.NewFloat(12.34))
	assert.Equal(t, "24.68", a1.String())
	assert.Equal(t, "2468", a1.AsSubUnit().String())
}
