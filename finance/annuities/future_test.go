package annuities

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

func TestAnticipateFutureFromPaymentsOnly(t *testing.T) {
	// With no principal or preset future value, AnticipateFuture must fall
	// through to contributionsAnticipateFuture: FV_due = FV_ordinary × (1+i).
	rate, err := compositeinterest.NewRateInterest(
		money.MustFromFloat64(0.01),
		compositeinterest.Monthly,
		compositeinterest.RateEffectyPeriodic,
	)
	require.NoError(t, err)

	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	ann, err := New(money.MustMoneyFromFloat64(500, money.USD), money.MoneyZero, money.MoneyZero, period, rate)
	require.NoError(t, err)

	future, err := ann.AnticipateFuture()
	require.NoError(t, err)

	ordinary := 500 * (math.Pow(1.01, 12) - 1) / 0.01
	expected := ordinary * 1.01
	assert.InDelta(t, expected, future.InexactFloat64(), 0.01)
}

func TestAnnuityFutureWithContributions(t *testing.T) {
	// PV = 10000, PMT = 1000, i = 0.01, n = 12
	// FV = PV×(1+i)^n + PMT×[(1+i)^n - 1]/i = 11268.2503 + 12682.5030 = 23950.7533
	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.12), compositeinterest.Monthly, compositeinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	total, err := annuity.FutureWithContributions()
	require.NoError(t, err)
	assert.InDelta(t, 23950.7533, total.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityFutureWithContributionsNoPrincipal(t *testing.T) {
	// With no principal (present = 0), the result must equal the
	// contributions-only future value.
	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.12), compositeinterest.Monthly, compositeinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	total, err := annuity.FutureWithContributions()
	require.NoError(t, err)
	assert.InDelta(t, 12682.5030, total.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityAnticipateFutureWithContributions(t *testing.T) {
	// Same as above but contributions are made at the start of each period,
	// so they earn one extra period of interest: 11268.2503 + 12809.3280 = 24077.5783
	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.12), compositeinterest.Monthly, compositeinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	total, err := annuity.AnticipateFutureWithContributions()
	require.NoError(t, err)
	assert.InDelta(t, 24077.5783, total.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityConfigFutureValue(t *testing.T) {
	total := NewAnnuity().
		Present(10000, money.USD).
		Value(1000, money.USD).
		AnnualRate(0.12).
		Periods(12).
		Monthly().
		MustFutureValue()

	assert.InDelta(t, 23950.7533, total.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityConfigAnticipateFutureValue(t *testing.T) {
	total := NewAnnuity().
		Present(10000, money.USD).
		Value(1000, money.USD).
		AnnualRate(0.12).
		Periods(12).
		Monthly().
		MustAnticipateFutureValue()

	assert.InDelta(t, 24077.5783, total.ToDecimal().InexactFloat64(), 0.01)
}
