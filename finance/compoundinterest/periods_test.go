package compoundinterest

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// newPeriodsCase builds a CompoundInterest whose period value is zero so that
// Periods() is forced to derive the answer from present, future, and rate via
// n = ln(FV/PV) / ln(1+i) instead of returning the configured period.
func newPeriodsCase(t *testing.T, present, future int64, rate float64, freq CompoundingFrequency, typeRate TypeRate) CompoundInterest {
	t.Helper()

	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(rate), freq, typeRate)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(0), freq)
	require.NoError(t, err)

	presentMoney, err := money.New(present, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(future, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	return ci
}

func TestPeriodsReturnsConfiguredPeriod(t *testing.T) {
	// When a non-zero period is configured, Periods() returns it directly
	// without computing anything from present/future/rate.
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(18), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(112683, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.Equal(t, 18.0, periods.InexactFloat64())
}

func TestPeriodsWithPeriodicRate(t *testing.T) {
	// n = ln(1126.83/1000) / ln(1.01) = 12.0004
	ci := newPeriodsCase(t, 100000, 112683, 0.01, Monthly, RateEffectyPeriodic)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 12.0004, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithNominalRate(t *testing.T) {
	// 12% nominal monthly => i = 0.01, n = ln(1.12683) / ln(1.01) = 12.0004
	ci := newPeriodsCase(t, 100000, 112683, 0.12, Monthly, RateEffectyNominal)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 12.0004, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithAnnuallyRate(t *testing.T) {
	// 12.68% effective annual => n = ln(1.12683) / ln(1.1268) = 1.0002
	ci := newPeriodsCase(t, 100000, 112683, 0.1268, Annually, RateEffectyAnnually)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 1.0002, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithQuarterlyCompounding(t *testing.T) {
	// 12% nominal quarterly => i = 0.03, n = ln(1.1268) / ln(1.03) = 4.0388
	ci := newPeriodsCase(t, 100000, 112680, 0.12, Quarterly, RateEffectyNominal)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 4.0388, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithDailyCompounding(t *testing.T) {
	// 10% effective annual => daily i = 1.10^(1/365) - 1
	// n = ln(1.105) / ln(1+i) = 382.3678 days (more than a 365-day year,
	// because 10.5% growth exceeds the 10% effective annual rate)
	ci := newPeriodsCase(t, 100000, 110500, 0.10, Daily, RateEffectyAnnually)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 382.3678, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithSemiAnnuallyCompounding(t *testing.T) {
	// 10% effective annual => semi-annual i = 1.10^(1/2) - 1
	// n = ln(1.1025) / ln(1+i) = 2.0476
	ci := newPeriodsCase(t, 500000, 551250, 0.10, SemiAnnually, RateEffectyAnnually)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 2.0476, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithFivePercentPeriodicRate(t *testing.T) {
	// n = ln(1.06167) / ln(1.05) = 1.2265
	ci := newPeriodsCase(t, 100000, 106167, 0.05, Monthly, RateEffectyPeriodic)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 1.2265, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithMultipleDataSets(t *testing.T) {
	testCases := []struct {
		name     string
		present  float64
		future   float64
		rate     float64
		freq     CompoundingFrequency
		typeRate TypeRate
		periods  float64
	}{
		{
			name:     "monthly compounding",
			present:  1000,
			future:   1126.825,
			rate:     0.01,
			freq:     Monthly,
			typeRate: RateEffectyPeriodic,
			// 1126.825 truncates to 1126.82 when stored as cents:
			// n = ln(1.12682) / ln(1.01) = 11.9996
			periods: 11.9996,
		},
		{
			name:     "annual compounding",
			present:  1000,
			future:   1050.0,
			rate:     0.05,
			freq:     Annually,
			typeRate: RateEffectyPeriodic,
			// n = ln(1.05) / ln(1.05) = 1
			periods: 1.0,
		},
		{
			name:     "quarterly compounding",
			present:  2000,
			future:   2342.7,
			rate:     0.02,
			freq:     Quarterly,
			typeRate: RateEffectyPeriodic,
			// 2342.70 loses a cent to float truncation when stored (2342.69):
			// n = ln(2342.69/2000) / ln(1.02) = 7.9864
			periods: 7.9864,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ci := newPeriodsCase(t, int64(tc.present*100), int64(tc.future*100), tc.rate, tc.freq, tc.typeRate)

			periods, err := ci.Periods()
			require.NoError(t, err)

			assert.InDelta(t, tc.periods, periods.InexactFloat64(), 0.0001)
		})
	}
}

func TestPeriodsWithTenPercentGrowth(t *testing.T) {
	// n = ln(1.10) / ln(1.01) = 9.5786
	ci := newPeriodsCase(t, 100000, 110000, 0.01, Monthly, RateEffectyPeriodic)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 9.5786, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithBimonthlyCompounding(t *testing.T) {
	// 6% nominal bimonthly => i = 0.01, n = ln(1.06136) / ln(1.01) = 5.9848
	ci := newPeriodsCase(t, 100000, 106136, 0.06, Bimonthly, RateEffectyNominal)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 5.9848, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithDifferentRateTypes(t *testing.T) {
	testCases := []struct {
		name     string
		rate     float64
		typeRate TypeRate
		expected float64
	}{
		// n = ln(1.12683) / ln(1.01) = 12.0004
		{"periodic", 0.01, RateEffectyPeriodic, 12.0004},
		// 12% nominal monthly => i = 0.01: same as above
		{"nominal", 0.12, RateEffectyNominal, 12.0004},
		// 12.68% effective annual => monthly i = 1.1268^(1/12) - 1
		// n = ln(1.12683) / ln(1+i) = 12.0027
		{"annual", 0.1268, RateEffectyAnnually, 12.0027},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ci := newPeriodsCase(t, 100000, 112683, tc.rate, Monthly, tc.typeRate)

			periods, err := ci.Periods()
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, periods.InexactFloat64(), 0.0001)
		})
	}
}

func TestPeriodsToTripleInvestment(t *testing.T) {
	// n = ln(3) / ln(1.05) = 22.5171
	ci := newPeriodsCase(t, 100000, 300000, 0.05, Monthly, RateEffectyPeriodic)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 22.5171, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithFractionalPeriods(t *testing.T) {
	// 10% nominal monthly => i = 0.1/12, n = ln(1.05114) / ln(1+i) = 6.0099
	ci := newPeriodsCase(t, 100000, 105114, 0.10, Monthly, RateEffectyNominal)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 6.0099, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsConsistencyWithFuture(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	originalPeriods := 12.0
	period, err := NewPeriod(decimal.MustFromFloat64(originalPeriods), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	// Recover the periods from the computed future value using a zero period
	// so that Periods() derives the answer instead of echoing the input.
	zeroPeriod, err := NewPeriod(decimal.MustFromFloat64(0), Monthly)
	require.NoError(t, err)
	ci2, err := New(presentMoney, future, rateInterest, zeroPeriod)
	require.NoError(t, err)

	calculatedPeriods, err := ci2.Periods()
	require.NoError(t, err)

	assert.InDelta(t, originalPeriods, calculatedPeriods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithSmallInterestRate(t *testing.T) {
	// n = ln(1.06167) / ln(1.001) = 59.8731
	ci := newPeriodsCase(t, 100000, 106167, 0.001, Monthly, RateEffectyPeriodic)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 59.8731, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsWithHighInterestRate(t *testing.T) {
	// n = ln(10.9605) / ln(1.10) = 25.1211
	ci := newPeriodsCase(t, 100000, 1096050, 0.10, Monthly, RateEffectyPeriodic)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 25.1211, periods.InexactFloat64(), 0.0001)
}

func TestPeriodsPropagatesPeriodError(t *testing.T) {
	// A zero-value CompoundInterest has an invalid (empty) period
	// frequency, so both the shortcut check and GetEqualsRateInterestPeriods
	// fail, and Periods must surface that error.
	var ci CompoundInterest

	_, err := ci.Periods()
	require.Error(t, err)
}

func TestPeriodsPropagatesOverflowFromRatio(t *testing.T) {
	// future/present computed at an extreme magnitude mismatch overflows
	// decimal128's 128-bit coefficient.
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(0), Monthly)
	require.NoError(t, err)

	present, err := money.New(1, 19, money.USD)
	require.NoError(t, err)
	future, err := money.New(math.MaxInt64, 0, money.USD)
	require.NoError(t, err)

	ci, err := New(present, future, rateInterest, period)
	require.NoError(t, err)

	_, err = ci.Periods()
	require.Error(t, err)
}

func TestPeriodsPropagatesLnErrorWhenRatioIsNegative(t *testing.T) {
	// A negative future/present ratio has no logarithm.
	ci := newPeriodsCase(t, 100000, -50000, 0.01, Monthly, RateEffectyPeriodic)

	_, err := ci.Periods()
	require.Error(t, err)
}

func TestPeriodsWithZeroPresent(t *testing.T) {
	// With no configured period and no present value there is nothing to
	// derive the periods from, so Periods() must return an error.
	ci := newPeriodsCase(t, 0, 1096050, 0.10, Monthly, RateEffectyPeriodic)

	periods, err := ci.Periods()
	require.Error(t, err)
	assert.Equal(t, ErrInvalidOperation, err)
	assert.Equal(t, 0.0, periods.InexactFloat64())
}
