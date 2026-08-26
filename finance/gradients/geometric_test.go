package gradients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestGeometricPresent(t *testing.T) {
	// A = 1000, g = 8%, i = 10%, n = 5
	// PV = 1000×[1-((1.08)/(1.10))^5]/(0.10-0.08) = 4383.1433
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.10), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.MustFromFloat64(0.08),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 4383.1433, present.GetDecimal().InexactFloat64(), 0.01)
}

func TestGeometricFuture(t *testing.T) {
	// FV = 1000×[1.10^5 - 1.08^5]/(0.10-0.08) = 7059.0962
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.10), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.MustFromFloat64(0.08),
		period, rate)
	require.NoError(t, err)

	future, err := series.Future()
	require.NoError(t, err)
	assert.InDelta(t, 7059.0962, future.GetDecimal().InexactFloat64(), 0.01)
}

func TestGeometricPresentGrowthEqualsRate(t *testing.T) {
	// When g == i, PV = A×n/(1+i) = 1000×5/1.10 = 4545.4545
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.10), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.MustFromFloat64(0.10),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 4545.4545, present.GetDecimal().InexactFloat64(), 0.01)
}

func TestGeometricFutureGrowthEqualsRate(t *testing.T) {
	// When g == i, FV = A×n×(1+i)^(n-1) = 1000×5×1.10^4 = 7320.50
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.10), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.MustFromFloat64(0.10),
		period, rate)
	require.NoError(t, err)

	future, err := series.Future()
	require.NoError(t, err)
	assert.InDelta(t, 7320.50, future.GetDecimal().InexactFloat64(), 0.01)
}

func TestGeometricZeroGrowthMatchesOrdinaryAnnuity(t *testing.T) {
	// With g = 0 the series collapses to a plain ordinary annuity of A per
	// period, so its PV must equal A × [1-(1+i)^-n]/i.
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.10), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.Zero,
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 3790.7868, present.GetDecimal().InexactFloat64(), 0.01)
}

func TestGeometricConfigBuilder(t *testing.T) {
	present, err := NewGeometricSeries().
		FirstPayment(1000, money.USD).
		GrowthRate(0.08).
		Rate(0.10).
		Periods(5).
		Annually().
		Present()
	require.NoError(t, err)
	assert.InDelta(t, 4383.1433, present.GetDecimal().InexactFloat64(), 0.01)

	future := NewGeometricSeries().
		FirstPayment(1000, money.USD).
		GrowthRate(0.08).
		Rate(0.10).
		Periods(5).
		Annually().
		MustFuture()
	assert.InDelta(t, 7059.0962, future.GetDecimal().InexactFloat64(), 0.01)
}
