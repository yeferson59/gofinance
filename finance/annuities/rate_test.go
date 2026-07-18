package annuities

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/finance/compoundinterest"
	"github.com/yeferson59/gofinance/money"
)

func TestAnnuityRateWithPresent(t *testing.T) {
	// PMT = 1000, PV = 11255.0775, n = 12 => the periodic rate that
	// reproduces this present value is i = 0.01 (12% nominal monthly).
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present := money.MustMoneyFromFloat64(11255.0775, money.USD)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	rate, err := annuity.RateWithPresent()
	require.NoError(t, err)
	assert.InDelta(t, 0.01, rate.InexactFloat64(), 0.0001)
}

func TestAnnuityRateWithFuture(t *testing.T) {
	// PMT = 1000, FV = 12682.5030, n = 12 => i = 0.01
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future := money.MustMoneyFromFloat64(12682.5030, money.USD)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	rate, err := annuity.RateWithFuture()
	require.NoError(t, err)
	assert.InDelta(t, 0.01, rate.InexactFloat64(), 0.0001)
}

func TestAnnuityAnticipateRateWithPresent(t *testing.T) {
	// PMT = 1000, PV_due = PV_ordinary × (1+i) = 11255.0775 × 1.01 =
	// 11367.6283, n = 12 => i = 0.01
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present := money.MustMoneyFromFloat64(11367.6283, money.USD)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	rate, err := annuity.AnticipateRateWithPresent()
	require.NoError(t, err)
	assert.InDelta(t, 0.01, rate.InexactFloat64(), 0.0001)
}

func TestAnnuityAnticipateRateWithFuture(t *testing.T) {
	// PMT = 1000, FV_due = FV_ordinary × (1+i) = 12682.5030 × 1.01 =
	// 12809.3280, n = 12 => i = 0.01
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future := money.MustMoneyFromFloat64(12809.3280, money.USD)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	rate, err := annuity.AnticipateRateWithFuture()
	require.NoError(t, err)
	assert.InDelta(t, 0.01, rate.InexactFloat64(), 0.0001)
}

func TestAnnuityRateWithPresentAcrossFrequencies(t *testing.T) {
	// PMT = 1000, n = 12 periods; the present values are the same ones
	// verified in TestAnnuityWithDifferentCompoundingFrequencies for a 12%
	// nominal rate at each frequency, so the recovered periodic rate must
	// match that frequency's periodic rate.
	testCases := []struct {
		name         string
		frequency    compoundinterest.CompoundingFrequency
		present      float64
		periodicRate float64
	}{
		{"daily", compoundinterest.Daily, 11974.3955, 0.12 / 365},
		{"monthly", compoundinterest.Monthly, 11255.0775, 0.01},
		{"quarterly", compoundinterest.QuarterlyOne, 9954.0040, 0.03},
		{"annually", compoundinterest.Annually, 6194.3742, 0.12},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), tc.frequency)
			require.NoError(t, err)

			rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), tc.frequency, compoundinterest.RateEffectyPeriodic)
			require.NoError(t, err)

			value, err := money.New(100000, 2, money.USD)
			require.NoError(t, err)
			present := money.MustMoneyFromFloat64(tc.present, money.USD)
			future, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			annuity, err := New(value, present, future, period, rateInterest)
			require.NoError(t, err)

			rate, err := annuity.RateWithPresent()
			require.NoError(t, err)
			assert.InDelta(t, tc.periodicRate, rate.InexactFloat64(), 0.0001)
		})
	}
}

func TestAnnuityRateWithPresentLongTerm(t *testing.T) {
	// A 30-year, 360-month mortgage-like scenario: PMT = 1000, i = 0.005,
	// n = 360. Verifies the solver stays accurate (and doesn't overflow)
	// for a large number of periods.
	i := 0.005
	n := 360.0
	pv := 1000.0 * (1 - math.Pow(1+i, -n)) / i

	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(n), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present := money.MustMoneyFromFloat64(pv, money.USD)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	rate, err := annuity.RateWithPresent()
	require.NoError(t, err)
	assert.InDelta(t, i, rate.InexactFloat64(), 0.0001)
}

// For a positive payment, presentValueOrdinary(payment, rate, periods) and
// futureValueOrdinary(payment, rate, periods) are strictly positive for
// every rate > -1 (the entire domain rateBounds searches), so a negative
// target can never be reached: both bracket endpoints land on the same
// side of it, and RateWithPresent/RateWithFuture (and their anticipated
// counterparts) must report ErrRateNotFound instead of guessing.

func TestAnnuityRateWithPresentNotFound(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present := money.MustMoneyFromFloat64(-100, money.USD)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	_, err = annuity.RateWithPresent()
	assert.True(t, errors.Is(err, ErrRateNotFound))
}

func TestAnnuityRateWithFutureNotFound(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future := money.MustMoneyFromFloat64(-100, money.USD)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	_, err = annuity.RateWithFuture()
	assert.True(t, errors.Is(err, ErrRateNotFound))
}

func TestAnnuityAnticipateRateWithPresentNotFound(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present := money.MustMoneyFromFloat64(-100, money.USD)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	_, err = annuity.AnticipateRateWithPresent()
	assert.True(t, errors.Is(err, ErrRateNotFound))
}

func TestAnnuityAnticipateRateWithFutureNotFound(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future := money.MustMoneyFromFloat64(-100, money.USD)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	_, err = annuity.AnticipateRateWithFuture()
	assert.True(t, errors.Is(err, ErrRateNotFound))
}
