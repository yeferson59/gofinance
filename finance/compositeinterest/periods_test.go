package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriodsWithPeriodicRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1126.825, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 12.0, periods, 0.01)
}

func TestPeriodsWithNominalRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.12, Monthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1126.825, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 12.0, periods, 0.01)
}

func TestPeriodsWithAnnuallyRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.1268, Annually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(1, Annually)
	require.NoError(t, err)

	ci, err := New(1000, 1126.825, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 1.0, periods, 0.01)
}

func TestPeriodsWithQuarterlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.12, QuarterlyOne, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(4, QuarterlyOne)
	require.NoError(t, err)

	ci, err := New(1000, 1126.8, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.True(t, periods > 0)
	assert.InDelta(t, 4.0, periods, 0.1)
}

func TestPeriodsWithDailyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, Daily, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(365, Daily)
	require.NoError(t, err)

	ci, err := New(1000, 1105, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.True(t, periods > 0)
	assert.InDelta(t, 365.0, periods, 50.0)
}

func TestPeriodsWithSemiAnnuallyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, SemiAnnually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(2, SemiAnnually)
	require.NoError(t, err)

	ci, err := New(5000, 5512.5, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 2.0, periods, 0.5)
}

func TestPeriodsWithZeroPeriodShouldError(t *testing.T) {
	rateInterest, err := NewRateInterest(0.05, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1061.67, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)
	assert.True(t, periods > 0)
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
			periods:  12,
		},
		{
			name:     "annual compounding",
			present:  1000,
			future:   1050.0,
			rate:     0.05,
			freq:     Annually,
			typeRate: RateEffectyPeriodic,
			periods:  1,
		},
		{
			name:     "quarterly compounding",
			present:  2000,
			future:   2342.7,
			rate:     0.02,
			freq:     QuarterlyOne,
			typeRate: RateEffectyPeriodic,
			periods:  8,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(tc.rate, tc.freq, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(tc.periods, tc.freq)
			require.NoError(t, err)

			ci, err := New(tc.present, tc.future, rateInterest, period)
			require.NoError(t, err)

			periods, err := ci.Periods()
			require.NoError(t, err)

			assert.InDelta(t, tc.periods, periods, 1.0)
		})
	}
}

func TestPeriodsErrorWithInvalidPeriod(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1100, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)
	assert.True(t, periods > 0)
}

func TestPeriodsWithBimonthlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.06, Bimonthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(6, Bimonthly)
	require.NoError(t, err)

	ci, err := New(1000, 1061.36, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 6.0, periods, 0.5)
}

func TestPeriodsWithDifferentRateTypes(t *testing.T) {
	testCases := []struct {
		name     string
		rate     float64
		typeRate TypeRate
	}{
		{"periodic", 0.01, RateEffectyPeriodic},
		{"nominal", 0.12, RateEffectyNominal},
		{"annual", 0.1268, RateEffectyAnnually},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(tc.rate, Monthly, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(12, Monthly)
			require.NoError(t, err)

			ci, err := New(1000, 1126.825, rateInterest, period)
			require.NoError(t, err)

			periods, err := ci.Periods()
			require.NoError(t, err)

			assert.InDelta(t, 12.0, periods, 0.1)
		})
	}
}

func TestPeriodsReturnPositiveValue(t *testing.T) {
	rateInterest, err := NewRateInterest(0.05, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(24, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 3000, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.True(t, periods > 0, "periods should always be positive")
}

func TestPeriodsWithFractionalPeriods(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, Monthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(6, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1051.14, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.True(t, periods > 0)
	assert.InDelta(t, 6.0, periods, 0.5)
}

func TestPeriodsConsistencyWithFuture(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	originalPeriods := 12.0
	period, err := NewPeriod(originalPeriods, Monthly)
	require.NoError(t, err)

	present := 1000.0
	ci, err := New(present, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	ci2, err := New(present, future, rateInterest, period)
	require.NoError(t, err)

	calculatedPeriods, err := ci2.Periods()
	require.NoError(t, err)

	assert.InDelta(t, originalPeriods, calculatedPeriods, 0.1)
}

func TestPeriodsWithSmallInterestRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.001, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(60, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1061.67, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 60.0, periods, 1.0)
}

func TestPeriodsWithHighInterestRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(24, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 10960.5, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 24.0, periods, 1.5)
}
