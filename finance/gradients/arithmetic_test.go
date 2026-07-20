package gradients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestArithmeticPresent(t *testing.T) {
	// A = 1000, G = 100, i = 10%, n = 5
	// PV = 1000×[1-1.1^-5]/0.1 + 100/0.1×{[1-1.1^-5]/0.1 - 5×1.1^-5} = 4476.9669
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.10), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	series, err := NewArithmetic(
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(100, money.USD),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 4476.9669, present.ToDecimal().InexactFloat64(), 0.01)
}

func TestArithmeticFuture(t *testing.T) {
	// FV = 1000×[1.1^5-1]/0.1 + 100/0.1×{[1.1^5-1]/0.1 - 5} = 7210.20
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.10), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	series, err := NewArithmetic(
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(100, money.USD),
		period, rate)
	require.NoError(t, err)

	future, err := series.Future()
	require.NoError(t, err)
	assert.InDelta(t, 7210.20, future.ToDecimal().InexactFloat64(), 0.01)
}

func TestArithmeticZeroGradientMatchesOrdinaryAnnuity(t *testing.T) {
	// With G = 0 the series collapses to a plain ordinary annuity of A per
	// period, so its PV must equal A × [1-(1+i)^-n]/i.
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.10), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	series, err := NewArithmetic(
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(0, money.USD),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 3790.7868, present.ToDecimal().InexactFloat64(), 0.01)
}

func TestArithmeticCurrencyMismatch(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.10), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	_, err = NewArithmetic(
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(100, money.EUR),
		period, rate)
	require.ErrorIs(t, err, money.ErrCurrencyMismatch)
}

func TestArithmeticConfigBuilder(t *testing.T) {
	present, err := NewArithmeticSeries().
		FirstPayment(1000, money.USD).
		Gradient(100, money.USD).
		Rate(0.10).
		Periods(5).
		Annually().
		Present()
	require.NoError(t, err)
	assert.InDelta(t, 4476.9669, present.ToDecimal().InexactFloat64(), 0.01)

	future := NewArithmeticSeries().
		FirstPayment(1000, money.USD).
		Gradient(100, money.USD).
		Rate(0.10).
		Periods(5).
		Annually().
		MustFuture()
	assert.InDelta(t, 7210.20, future.ToDecimal().InexactFloat64(), 0.01)
}
