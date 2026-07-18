package annuities

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/finance/compoundinterest"
	"github.com/yeferson59/gofinance/money"
)

func TestAnnuityPresent(t *testing.T) {
	// PMT = 1000, 12% nominal monthly => i = 0.01, n = 12
	// PV = PMT × [1 - (1+i)^-n] / i = 1000 × [1 - 1.01^-12] / 0.01 = 11255.0775
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	presentValue, err := annuity.Present()
	require.NoError(t, err)
	assert.InDelta(t, 11255.0775, presentValue.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityPresentWithZeroInterestRate(t *testing.T) {
	// With a zero interest rate the annuity formula divides by zero,
	// so Present must return an error instead of a value.
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.0), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1200000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	presentValue, err := annuity.Present()
	require.Error(t, err)
	assert.Equal(t, 0.0, presentValue.ToDecimal().InexactFloat64())
}

func TestAnnuityPresentWithSmallPeriods(t *testing.T) {
	// PMT = 100, 5% annual, n = 1 => PV = 100 / 1.05 = 95.2381
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(1), compoundinterest.Annually)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.05), compoundinterest.Annually, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(10000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(10500, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	presentValue, err := annuity.Present()
	require.NoError(t, err)
	assert.InDelta(t, 95.2381, presentValue.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityFuture(t *testing.T) {
	// Future value is derived from the underlying compound interest data:
	// FV = PV × (1+i)^n = 10000 × 1.01^12 = 11268.2503
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

	futureValue, err := annuity.Future()
	require.NoError(t, err)
	assert.InDelta(t, 11268.2503, futureValue.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityAnticipatePresent(t *testing.T) {
	// PV_due = PV_ordinary × (1+i); with the same inputs as
	// TestAnnuityPresent, PV_ordinary = 11255.0775, so
	// PV_due = 11255.0775 × 1.01 = 11367.628275.
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	presentValue, err := annuity.AnticipatePresent()
	require.NoError(t, err)
	assert.InDelta(t, 11367.628275, presentValue.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityPaymentFromPresentValue(t *testing.T) {
	// PV = 10000, i = 0.01, n = 12
	// PMT = PV × i(1+i)^n / [(1+i)^n - 1] = 888.4879
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	payment, err := annuity.PaymentFromPresentValue()
	require.NoError(t, err)
	assert.InDelta(t, 888.4879, payment.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityPaymentFromFutureValue(t *testing.T) {
	// FV = 15000, i = 0.01, n = 12
	// PMT = FV × i / [(1+i)^n - 1] = 1182.7318
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	payment, err := annuity.PaymentFromFutureValue()
	require.NoError(t, err)
	assert.InDelta(t, 1182.7318, payment.ToDecimal().InexactFloat64(), 0.01)
}

func TestNewAnnuityWithInvalidPeriod(t *testing.T) {
	// Test creation with invalid period - negative value should fail at NewPeriod
	_, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(-1), compoundinterest.Monthly)
	assert.Error(t, err)
	assert.Equal(t, "value periods must be greater or equal to zero", err.Error())

	// Since period creation failed, we can't create an annuity with it
	// The test above already validates the error handling
}

func TestAnnuityPresentMathematicalCorrectness(t *testing.T) {
	// Verify that the optimized formula produces correct results
	// Using a known case where we can verify the math
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(10), compoundinterest.Annually)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.05), compoundinterest.Annually, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	presentValue, err := annuity.Present()
	require.NoError(t, err)

	// Manual calculation: PV = A * [1 - (1+i)^-n] / i
	// Where A = 1000, i = 0.05, n = 10
	expectedPow := math.Pow(1.05, 10)
	expectedPresent := 1000 * (1 - 1/expectedPow) / 0.05

	// Should match within floating point precision
	assert.InDelta(t, expectedPresent, presentValue.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityWithDifferentCompoundingFrequencies(t *testing.T) {
	// PMT = 1000, PV = 10000, 12% nominal at each frequency, n = 12 periods
	// of that frequency => i = 0.12/m with m periods per year.
	// Annuity PV = PMT × [1 - (1+i)^-12] / i
	// Compound FV = PV × (1+i)^12
	testCases := []struct {
		name            string
		frequency       compoundinterest.CompoundingFrequency
		expectedPresent float64
		expectedFuture  float64
	}{
		{"daily", compoundinterest.Daily, 11974.3955, 10039.5235},
		{"monthly", compoundinterest.Monthly, 11255.0775, 11268.2503},
		{"quarterly", compoundinterest.QuarterlyOne, 9954.0040, 14257.6089},
		{"annually", compoundinterest.Annually, 6194.3742, 38959.7599},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), tc.frequency)
			require.NoError(t, err)

			rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), tc.frequency, compoundinterest.RateEffectyNominal)
			require.NoError(t, err)

			value, err := money.New(100000, 2, money.USD)
			require.NoError(t, err)
			present, err := money.New(1000000, 2, money.USD)
			require.NoError(t, err)
			future, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			annuity, err := New(value, present, future, period, rateInterest)
			require.NoError(t, err)

			presentValue, err := annuity.Present()
			require.NoError(t, err)
			assert.InDelta(t, tc.expectedPresent, presentValue.ToDecimal().InexactFloat64(), 0.01)

			futureValue, err := annuity.Future()
			require.NoError(t, err)
			assert.InDelta(t, tc.expectedFuture, futureValue.ToDecimal().InexactFloat64(), 0.01)
		})
	}
}

func TestPresentPropagatesRateInterestPeriodsError(t *testing.T) {
	// A zero-value Annuity has an invalid (empty) period frequency, so
	// GetEqualsRateInterestPeriods fails and Present must surface that
	// error instead of a bogus zero value.
	var annuity Annuity

	_, err := annuity.Present()
	assert.Error(t, err)
}

func TestPresentPropagatesPowOverflow(t *testing.T) {
	// (1+r)^n overflows decimal128's 128-bit coefficient when both the rate
	// and the period count are astronomically large.
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

	_, err = annuity.Present()
	assert.Error(t, err)
}

func TestAnticipatePresentPropagatesRateInterestPeriodsError(t *testing.T) {
	// A zero-value Annuity has an invalid (empty) period frequency, so
	// GetEqualsRateInterestPeriods fails and AnticipatePresent must surface
	// that error instead of a bogus zero value.
	var annuity Annuity

	_, err := annuity.AnticipatePresent()
	assert.Error(t, err)
}

func TestAnticipatePresentPropagatesPowOverflow(t *testing.T) {
	// (1+r)^(n-1) overflows decimal128's 128-bit coefficient when both the
	// rate and the period count are astronomically large.
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

	_, err = annuity.AnticipatePresent()
	assert.Error(t, err)
}

func TestAnticipatePresentWithZeroInterestRate(t *testing.T) {
	// With a zero interest rate the annuity formula divides by zero, so
	// AnticipatePresent must return an error instead of a value.
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.0), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1200000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	_, err = annuity.AnticipatePresent()
	assert.Error(t, err)
}
