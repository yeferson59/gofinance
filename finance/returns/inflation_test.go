package returns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestRealValue(t *testing.T) {
	// $1,000 in 10 years at 3% inflation ≈ $744.09 today.
	real, err := RealValue(usd(1000), money.MustFromFloat64(0.03), money.MustFromFloat64(10))
	require.NoError(t, err)
	assert.InDelta(t, 744.09, real.InexactFloat64(), 0.01)
	assert.Equal(t, money.USD, real.Currency())
}

func TestNominalValueRoundTrip(t *testing.T) {
	// NominalValue is the inverse of RealValue.
	nominal, err := NominalValue(usd(744.09), money.MustFromFloat64(0.03), money.MustFromFloat64(10))
	require.NoError(t, err)
	assert.InDelta(t, 1000.0, nominal.InexactFloat64(), 0.01)
}

func TestRealRate(t *testing.T) {
	// Fisher: (1.08/1.03) − 1 ≈ 0.048544.
	r, err := RealRate(money.MustFromFloat64(0.08), money.MustFromFloat64(0.03))
	require.NoError(t, err)
	assert.InDelta(t, 0.048544, r.InexactFloat64(), 1e-6)
}

func TestInflationErrors(t *testing.T) {
	_, err := RealValue(usd(1000), money.MustFromFloat64(-1.5), money.MustFromFloat64(10))
	assert.ErrorIs(t, err, ErrInvalidInflationRate)

	_, err = RealRate(money.MustFromFloat64(0.08), money.MustFromFloat64(-2))
	assert.ErrorIs(t, err, ErrInvalidInflationRate)
}

func TestMustInflationHelpers(t *testing.T) {
	assert.InDelta(t, 744.09, MustRealValue(usd(1000), money.MustFromFloat64(0.03), money.MustFromFloat64(10)).InexactFloat64(), 0.01)
	assert.Panics(t, func() { MustRealRate(money.MustFromFloat64(0.08), money.MustFromFloat64(-2)) })
}
