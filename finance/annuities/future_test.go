package annuities

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestAnticipateFutureFromPaymentsOnly(t *testing.T) {
	// With no principal or preset future value, AnticipateFuture must fall
	// through to contributionsAnticipateFuture: FV_due = FV_ordinary × (1+i).
	rate, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.01),
		compoundinterest.Monthly,
		compoundinterest.RateEffectyPeriodic,
	)
	require.NoError(t, err)

	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
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
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
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
	assert.InDelta(t, 23950.7533, total.GetDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityFutureWithContributionsNoPrincipal(t *testing.T) {
	// With no principal (present = 0), the result must equal the
	// contributions-only future value.
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
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
	assert.InDelta(t, 12682.5030, total.GetDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityAnticipateFutureWithContributions(t *testing.T) {
	// Same as above but contributions are made at the start of each period,
	// so they earn one extra period of interest: 11268.2503 + 12809.3280 = 24077.5783
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
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
	assert.InDelta(t, 24077.5783, total.GetDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityConfigFutureValue(t *testing.T) {
	total := NewAnnuity().
		Present(10000, money.USD).
		Value(1000, money.USD).
		AnnualRate(0.12).
		Periods(12).
		Monthly().
		MustFutureValue()

	assert.InDelta(t, 23950.7533, total.GetDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityConfigAnticipateFutureValue(t *testing.T) {
	total := NewAnnuity().
		Present(10000, money.USD).
		Value(1000, money.USD).
		AnnualRate(0.12).
		Periods(12).
		Monthly().
		MustAnticipateFutureValue()

	assert.InDelta(t, 24077.5783, total.GetDecimal().InexactFloat64(), 0.01)
}

func TestFuturePropagatesRateInterestPeriodsError(t *testing.T) {
	// A zero-value Annuity has an invalid (empty) period frequency, so
	// GetEqualsRateInterestPeriods fails inside contributionsFuture, and
	// Future must surface that error instead of a bogus zero value.
	var annuity Annuity

	_, err := annuity.Future()
	assert.Error(t, err)
}

func TestContributionsFuturePropagatesPowOverflow(t *testing.T) {
	// (1+r)^n overflows decimal128's 128-bit coefficient when both the rate
	// and the period count are astronomically large. This needs a 1000-period
	// term, which newMonthlyPeriodicAnnuity's fixed 12 periods can't express,
	// so it's built directly here instead.
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(1000), compoundinterest.Monthly)
	require.NoError(t, err)
	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(1000), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)
	annuity, err := New(
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MoneyZero,
		money.MoneyZero,
		period,
		rateInterest,
	)
	require.NoError(t, err)

	_, err = annuity.contributionsFuture()
	assert.Error(t, err)
}

func TestContributionsFutureWithZeroRate(t *testing.T) {
	// With no interest the contributions never grow, so the future value is
	// simply their sum: 1000 × 12 = 12000. The general formula divides by
	// the rate, so the limit is returned directly (TESTING_PLAN.md §2.1).
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 0, 0, 0)

	future, err := annuity.contributionsFuture()
	require.NoError(t, err)
	assert.InDelta(t, 12000.0, future.InexactFloat64(), 1e-9)
}

func TestAnticipateFutureIgnoresPrincipal(t *testing.T) {
	// Even with a principal configured, AnticipateFuture must return the
	// payments' annuity-due value — the principal's growth is only
	// reachable through PrincipalFuture / the WithContributions variants.
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 1000, 0, 0.01)

	future, err := annuity.AnticipateFuture()
	require.NoError(t, err)

	// FV_due = PMT × ((1+i)^n − 1)/i × (1+i) = 12682.5030 × 1.01 = 12809.3280
	assert.InDelta(t, 12809.3280, future.InexactFloat64(), 0.01)
}

func TestPrincipalFutureReturnsCompoundedPrincipal(t *testing.T) {
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 1000, 0, 0.01)

	principal, err := annuity.PrincipalFuture()
	require.NoError(t, err)

	// FV = PV × (1+i)^n = 1000 × 1.01^12 = 1126.8250
	assert.InDelta(t, 1126.8250, principal.InexactFloat64(), 0.01)
}

func TestContributionsAnticipateFutureWithZeroRate(t *testing.T) {
	// contributionsAnticipateFuture is contributionsFuture × (1+i), so at a
	// zero rate it matches the ordinary value: 1000 × 12 = 12000
	// (TESTING_PLAN.md §2.1).
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 0, 0, 0)

	future, err := annuity.contributionsAnticipateFuture()
	require.NoError(t, err)
	assert.InDelta(t, 12000.0, future.InexactFloat64(), 1e-9)
}

func TestPrincipalFuturePropagatesNonInvalidOperationError(t *testing.T) {
	// A zero-value Annuity fails GetEqualsRateInterestPeriods with a
	// different error than ErrInvalidOperation, so principalFuture must
	// surface it instead of silently returning zero.
	var annuity Annuity

	_, err := annuity.PrincipalFuture()
	assert.Error(t, err)
	assert.NotErrorIs(t, err, compoundinterest.ErrInvalidOperation)
}

func TestFutureWithContributionsWithZeroRate(t *testing.T) {
	// No principal and no interest: the future value is the contributions'
	// sum, 1000 × 12 = 12000 (TESTING_PLAN.md §2.1).
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 0, 0, 0)

	future, err := annuity.FutureWithContributions()
	require.NoError(t, err)
	assert.InDelta(t, 12000.0, future.InexactFloat64(), 1e-9)
}

func TestAnticipateFutureWithContributionsWithZeroRate(t *testing.T) {
	// Same as above, but through the anticipated (annuity due) variant: at a
	// zero rate the payment's timing inside the period makes no difference.
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 0, 0, 0)

	future, err := annuity.AnticipateFutureWithContributions()
	require.NoError(t, err)
	assert.InDelta(t, 12000.0, future.InexactFloat64(), 1e-9)
}
