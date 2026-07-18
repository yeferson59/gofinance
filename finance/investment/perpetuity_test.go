package investment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/money"
)

func TestPerpetuity(t *testing.T) {
	// $100/period forever at 5% → 100/0.05 = 2000.
	pv, err := Perpetuity(usd(100), decimal.MustFromFloat64(0.05))
	require.NoError(t, err)
	assert.InDelta(t, 2000.0, pv.InexactFloat64(), 1e-9)
	assert.Equal(t, money.USD, pv.Currency())
}

func TestPerpetuityNonPositiveRate(t *testing.T) {
	_, err := Perpetuity(usd(100), decimal.Zero)
	assert.ErrorIs(t, err, ErrNonPositiveRate)
}

func TestGrowingPerpetuity(t *testing.T) {
	// D1=100, r=8%, g=3% → 100/(0.05) = 2000.
	pv, err := GrowingPerpetuity(usd(100), decimal.MustFromFloat64(0.08), decimal.MustFromFloat64(0.03))
	require.NoError(t, err)
	assert.InDelta(t, 2000.0, pv.InexactFloat64(), 1e-9)
}

func TestGrowingPerpetuityRateBelowGrowth(t *testing.T) {
	_, err := GrowingPerpetuity(usd(100), decimal.MustFromFloat64(0.03), decimal.MustFromFloat64(0.03))
	assert.ErrorIs(t, err, ErrRateBelowGrowth)
}

func TestMustPerpetuity(t *testing.T) {
	assert.InDelta(t, 2000.0, MustPerpetuity(usd(100), decimal.MustFromFloat64(0.05)).InexactFloat64(), 1e-9)
	assert.Panics(t, func() { MustGrowingPerpetuity(usd(100), decimal.MustFromFloat64(0.03), decimal.MustFromFloat64(0.05)) })
}
