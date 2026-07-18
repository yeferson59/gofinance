package annuities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/finance/compoundinterest"
	"github.com/yeferson59/gofinance/money"
)

// newMonthlyPeriodicAnnuity builds an Annuity with a monthly periodic rate
// (so the rate is used as-is, with no frequency/type conversion), for
// exercising the error branches of the Periods* methods.
func newMonthlyPeriodicAnnuity(t *testing.T, value, present, future, rate float64) Annuity {
	t.Helper()

	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(rate), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	valueMoney := money.MustMoneyFromFloat64(value, money.USD)
	presentMoney := money.MustMoneyFromFloat64(present, money.USD)
	futureMoney := money.MustMoneyFromFloat64(future, money.USD)

	annuity, err := New(valueMoney, presentMoney, futureMoney, period, rateInterest)
	require.NoError(t, err)

	return annuity
}

func TestAnnuityPeriodsWithPresent(t *testing.T) {
	// PMT = 1000, PV = 10000, i = 0.01
	// n = ln(PMT / (PMT - PV×i)) / ln(1+i) = ln(1000/900) / ln(1.01) = 10.5886
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	periods, err := annuity.PeriodsWithPresent()
	require.NoError(t, err)
	assert.InDelta(t, 10.5886, periods.InexactFloat64(), 0.0001)
}

func TestAnnuityPeriodsWithFuture(t *testing.T) {
	// PMT = 1000, FV = 15000, i = 0.01
	// n = ln((FV×i + PMT) / PMT) / ln(1+i) = ln(1.15) / ln(1.01) = 14.0460
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	periods, err := annuity.PeriodsWithFuture()
	require.NoError(t, err)
	assert.InDelta(t, 14.0460, periods.InexactFloat64(), 0.0001)
}

func TestAnnuityAnticipatePeriodsWithPresent(t *testing.T) {
	// PMT = 1000, PV = 10000, i = 0.01
	// n = ln(PMT / (PMT - [PV/(1+i)]×i)) / ln(1+i) = 10.4781
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	periods, err := annuity.AnticipatePeriodsWithPresent()
	require.NoError(t, err)
	assert.InDelta(t, 10.4781, periods.InexactFloat64(), 0.0001)
}

func TestAnnuityAnticipatePeriodsWithFuture(t *testing.T) {
	// PMT = 1000, FV = 15000, i = 0.01
	// n = ln(([FV/(1+i)]×i + PMT) / PMT) / ln(1+i) = 13.9161
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	periods, err := annuity.AnticipatePeriodsWithFuture()
	require.NoError(t, err)
	assert.InDelta(t, 13.9161, periods.InexactFloat64(), 0.0001)
}

func TestAnnuityPeriodsWithPresentZeroRate(t *testing.T) {
	// With rate = 0, ln(1+i) = ln(1) = 0, so PeriodsWithPresent must return
	// an error instead of dividing by that zero.
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	_, err = annuity.PeriodsWithPresent()
	assert.Error(t, err)
}

func TestAnnuityPeriodsWithFutureZeroRate(t *testing.T) {
	// With rate = 0, ln(1+i) = ln(1) = 0, so PeriodsWithFuture must return
	// an error instead of dividing by that zero.
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	_, err = annuity.PeriodsWithFuture()
	assert.Error(t, err)
}

func TestAnnuityPeriodsWithPresentAcrossFrequencies(t *testing.T) {
	// PMT = 1000, 12% nominal at each frequency, n = 12 periods of that
	// frequency. Present values are the same ones verified in
	// TestAnnuityWithDifferentCompoundingFrequencies, so recovering the
	// periods from them must round-trip back to 12.
	testCases := []struct {
		name      string
		frequency compoundinterest.CompoundingFrequency
		present   float64
	}{
		{"daily", compoundinterest.Daily, 11974.3955},
		{"monthly", compoundinterest.Monthly, 11255.0775},
		{"quarterly", compoundinterest.QuarterlyOne, 9954.0040},
		{"annually", compoundinterest.Annually, 6194.3742},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), tc.frequency)
			require.NoError(t, err)

			rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), tc.frequency, compoundinterest.RateEffectyNominal)
			require.NoError(t, err)

			value, err := money.New(100000, 2, money.USD)
			require.NoError(t, err)
			present := money.MustMoneyFromFloat64(tc.present, money.USD)
			future, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			annuity, err := New(value, present, future, period, rateInterest)
			require.NoError(t, err)

			periods, err := annuity.PeriodsWithPresent()
			require.NoError(t, err)
			assert.InDelta(t, 12.0, periods.InexactFloat64(), 0.01)
		})
	}
}

func TestAnnuityPeriodsWithFutureAcrossFrequencies(t *testing.T) {
	// PMT = 1000, 12% nominal at each frequency, n = 12 periods of that
	// frequency. Future values are FV = PMT × [(1+i)^12 - 1]/i, so recovering
	// the periods from them must round-trip back to 12.
	testCases := []struct {
		name      string
		frequency compoundinterest.CompoundingFrequency
		future    float64
	}{
		{"daily", compoundinterest.Daily, 12021.7224},
		{"monthly", compoundinterest.Monthly, 12682.5030},
		{"quarterly", compoundinterest.QuarterlyOne, 14192.0296},
		{"annually", compoundinterest.Annually, 24133.1333},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), tc.frequency)
			require.NoError(t, err)

			rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), tc.frequency, compoundinterest.RateEffectyNominal)
			require.NoError(t, err)

			value, err := money.New(100000, 2, money.USD)
			require.NoError(t, err)
			present, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			future := money.MustMoneyFromFloat64(tc.future, money.USD)
			annuity, err := New(value, present, future, period, rateInterest)
			require.NoError(t, err)

			periods, err := annuity.PeriodsWithFuture()
			require.NoError(t, err)
			assert.InDelta(t, 12.0, periods.InexactFloat64(), 0.01)
		})
	}
}

func TestAnnuityPeriodsPropagatesRateInterestPeriodsError(t *testing.T) {
	// A zero-value Annuity has an invalid (empty) period frequency, so
	// GetEqualsRateInterestPeriods fails and every Periods* method must
	// surface that error.
	var annuity Annuity

	_, err := annuity.PeriodsWithPresent()
	assert.Error(t, err)

	_, err = annuity.PeriodsWithFuture()
	assert.Error(t, err)

	_, err = annuity.AnticipatePeriodsWithPresent()
	assert.Error(t, err)

	_, err = annuity.AnticipatePeriodsWithFuture()
	assert.Error(t, err)
}

func TestAnnuityPeriodsPropagatesPresentAndFutureErrors(t *testing.T) {
	// With present = future = 0, Present()/Future() can't derive either
	// value (both branches of their own ErrInvalidOperation check trigger),
	// so PeriodsWithPresent/PeriodsWithFuture and their anticipated variants
	// must surface that error instead of dividing by an empty value.
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 0, 0, 0.05)

	_, err := annuity.PeriodsWithPresent()
	assert.Error(t, err)

	_, err = annuity.PeriodsWithFuture()
	assert.Error(t, err)

	_, err = annuity.AnticipatePeriodsWithPresent()
	assert.Error(t, err)

	_, err = annuity.AnticipatePeriodsWithFuture()
	assert.Error(t, err)
}

func TestAnnuityPeriodsWithPresentPropagatesDivideByZero(t *testing.T) {
	// PMT - PV×i = 0 when PV×i exactly equals PMT (1000 = 20000 × 0.05),
	// so the ratio's denominator is zero.
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 20000, 0, 0.05)

	_, err := annuity.PeriodsWithPresent()
	assert.Error(t, err)
}

func TestAnnuityPeriodsWithPresentPropagatesLnError(t *testing.T) {
	// PV×i (30000 × 0.05 = 1500) exceeds PMT (1000), so the ratio's
	// denominator (and therefore the ratio itself) is negative, which has
	// no logarithm.
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 30000, 0, 0.05)

	_, err := annuity.PeriodsWithPresent()
	assert.Error(t, err)
}

func TestAnnuityPeriodsWithFuturePropagatesDivideByZero(t *testing.T) {
	// The ratio's denominator is PMT itself, so a zero payment divides by
	// zero.
	annuity := newMonthlyPeriodicAnnuity(t, 0, 0, 1500000, 0.5)

	_, err := annuity.PeriodsWithFuture()
	assert.Error(t, err)
}

func TestAnnuityPeriodsWithFuturePropagatesLnError(t *testing.T) {
	// A negative PMT with FV×i large enough to flip the numerator positive
	// makes the ratio negative, which has no logarithm.
	annuity := newMonthlyPeriodicAnnuity(t, -1000, 0, 1000000, 0.5)

	_, err := annuity.PeriodsWithFuture()
	assert.Error(t, err)
}

func TestAnnuityAnticipatePeriodsWithPresentZeroRate(t *testing.T) {
	// With rate = 0, ln(1+i) = ln(1) = 0, so AnticipatePeriodsWithPresent
	// must return an error instead of dividing by that zero.
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 20000, 0, 0)

	_, err := annuity.AnticipatePeriodsWithPresent()
	assert.Error(t, err)
}

func TestAnnuityAnticipatePeriodsWithPresentPropagatesDivideByZero(t *testing.T) {
	// Reduced to the ordinary present value (present/(1+i) = 21000/1.05 =
	// 20000) times the rate exactly equals PMT (1000), so the ratio's
	// denominator is zero.
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 21000, 0, 0.05)

	_, err := annuity.AnticipatePeriodsWithPresent()
	assert.Error(t, err)
}

func TestAnnuityAnticipatePeriodsWithPresentPropagatesLnError(t *testing.T) {
	// Same as the divide-by-zero case but scaled up (31500/1.05 = 30000)
	// so the reduced present value's contribution exceeds PMT, making the
	// ratio negative.
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 31500, 0, 0.05)

	_, err := annuity.AnticipatePeriodsWithPresent()
	assert.Error(t, err)
}

func TestAnnuityAnticipatePeriodsWithFutureZeroRate(t *testing.T) {
	// With rate = 0, ln(1+i) = ln(1) = 0, so AnticipatePeriodsWithFuture
	// must return an error instead of dividing by that zero.
	annuity := newMonthlyPeriodicAnnuity(t, 1000, 0, 1500000, 0)

	_, err := annuity.AnticipatePeriodsWithFuture()
	assert.Error(t, err)
}

func TestAnnuityAnticipatePeriodsWithFuturePropagatesDivideByZero(t *testing.T) {
	// The ratio's denominator is PMT itself, so a zero payment divides by
	// zero regardless of the future/(1+i) reduction.
	annuity := newMonthlyPeriodicAnnuity(t, 0, 0, 1500000, 0.5)

	_, err := annuity.AnticipatePeriodsWithFuture()
	assert.Error(t, err)
}

func TestAnnuityAnticipatePeriodsWithFuturePropagatesLnError(t *testing.T) {
	// A negative PMT with the reduced future value's contribution large
	// enough to flip the numerator positive makes the ratio negative,
	// which has no logarithm.
	annuity := newMonthlyPeriodicAnnuity(t, -1000, 0, 1500000, 0.5)

	_, err := annuity.AnticipatePeriodsWithFuture()
	assert.Error(t, err)
}
