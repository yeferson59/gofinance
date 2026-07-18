package returns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/money"
)

func usd(amount float64) money.Money {
	return money.MustMoneyFromFloat64(amount, money.USD)
}

func TestCAGR(t *testing.T) {
	// Doubling over 5 years: 2^(1/5) − 1 ≈ 0.148698.
	cagr, err := CAGR(usd(1000), usd(2000), decimal.MustFromFloat64(5))
	require.NoError(t, err)
	assert.InDelta(t, 0.148698, cagr.InexactFloat64(), 1e-6)
}

func TestCAGRDecline(t *testing.T) {
	// Halving over 2 years is a negative CAGR.
	cagr, err := CAGR(usd(1000), usd(500), decimal.MustFromFloat64(2))
	require.NoError(t, err)
	assert.InDelta(t, -0.292893, cagr.InexactFloat64(), 1e-6)
}

func TestCAGRErrors(t *testing.T) {
	_, err := CAGR(usd(1000), money.MustMoneyFromFloat64(2000, money.EUR), decimal.MustFromFloat64(5))
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	_, err = CAGR(usd(0), usd(2000), decimal.MustFromFloat64(5))
	assert.ErrorIs(t, err, ErrNonPositiveValue)

	_, err = CAGR(usd(1000), usd(2000), decimal.MustFromFloat64(0))
	assert.ErrorIs(t, err, ErrNonPositivePeriods)
}

func TestROI(t *testing.T) {
	roi, err := ROI(usd(1000), usd(1250))
	require.NoError(t, err)
	assert.InDelta(t, 0.25, roi.InexactFloat64(), 1e-9)
}

func TestROILoss(t *testing.T) {
	roi, err := ROI(usd(1000), usd(900))
	require.NoError(t, err)
	assert.InDelta(t, -0.10, roi.InexactFloat64(), 1e-9)
}

func TestROIErrors(t *testing.T) {
	_, err := ROI(usd(1000), money.MustMoneyFromFloat64(1250, money.EUR))
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	_, err = ROI(usd(0), usd(1250))
	assert.ErrorIs(t, err, ErrNonPositiveValue)
}

func TestHoldingPeriodReturn(t *testing.T) {
	// Price 1000 → 1100 plus 50 income = 15% total.
	hpr, err := HoldingPeriodReturn(usd(1000), usd(1100), usd(50))
	require.NoError(t, err)
	assert.InDelta(t, 0.15, hpr.InexactFloat64(), 1e-9)
}

func TestHoldingPeriodReturnCurrencyMismatch(t *testing.T) {
	_, err := HoldingPeriodReturn(usd(1000), usd(1100), money.MustMoneyFromFloat64(50, money.EUR))
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)
}

func TestAnnualized(t *testing.T) {
	// 20% cumulative over 3 periods: 1.2^(1/3) − 1 ≈ 0.062659.
	a, err := Annualized(decimal.MustFromFloat64(0.20), decimal.MustFromFloat64(3))
	require.NoError(t, err)
	assert.InDelta(t, 0.062659, a.InexactFloat64(), 1e-6)
}

func TestAnnualizedErrors(t *testing.T) {
	_, err := Annualized(decimal.MustFromFloat64(0.20), decimal.MustFromFloat64(0))
	assert.ErrorIs(t, err, ErrNonPositivePeriods)

	// Losing more than everything is undefined.
	_, err = Annualized(decimal.MustFromFloat64(-1.5), decimal.MustFromFloat64(3))
	assert.ErrorIs(t, err, ErrNonPositiveValue)
}

func TestMustHelpers(t *testing.T) {
	assert.InDelta(t, 0.25, MustROI(usd(1000), usd(1250)).InexactFloat64(), 1e-9)
	assert.Panics(t, func() { MustCAGR(usd(0), usd(1), decimal.MustFromFloat64(1)) })
}
