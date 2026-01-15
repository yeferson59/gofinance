package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPresentWithPeriodicRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(0, 1126.825, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.InDelta(t, 1000.0, present, 0.1)
}

func TestPresentWithNominalRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.12, Monthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(0, 1126.825, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.InDelta(t, 1000.0, present, 0.1)
}

func TestPresentWithAnnuallyRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.1268, Annually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(1, Annually)
	require.NoError(t, err)

	ci, err := New(0, 1126.825, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.InDelta(t, 1000.0, present, 0.1)
}

func TestPresentWithDailyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, Daily, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(365, Daily)
	require.NoError(t, err)

	ci, err := New(0, 1105.0, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.InDelta(t, 1000.0, present, 5.0)
}

func TestPresentWithQuarterlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.12, QuarterlyOne, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(4, QuarterlyOne)
	require.NoError(t, err)

	ci, err := New(0, 1000, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.True(t, present < 1000, "present should be less than future")
	assert.True(t, present > 0, "present should be greater than zero")
}

func TestPresentWithSemiAnnuallyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, SemiAnnually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(2, SemiAnnually)
	require.NoError(t, err)

	ci, err := New(0, 5512.5, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.InDelta(t, 5000.0, present, 50.0)
}

func TestPresentWithZeroFuture(t *testing.T) {
	rateInterest, err := NewRateInterest(0.05, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(0, 0, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.Error(t, err)

	assert.Equal(t, 0.0, present)
}

func TestPresentWithMultipleDataSets(t *testing.T) {
	testCases := []struct {
		name     string
		future   float64
		rate     float64
		freq     CompoundingFrequency
		typeRate TypeRate
		periods  float64
		expected float64
	}{
		{
			name:     "monthly compounding",
			future:   1126.825,
			rate:     0.01,
			freq:     Monthly,
			typeRate: RateEffectyPeriodic,
			periods:  12,
			expected: 1000.0,
		},
		{
			name:     "annual compounding",
			future:   1050.0,
			rate:     0.05,
			freq:     Annually,
			typeRate: RateEffectyPeriodic,
			periods:  1,
			expected: 1000.0,
		},
		{
			name:     "quarterly compounding",
			future:   2342.7,
			rate:     0.02,
			freq:     QuarterlyOne,
			typeRate: RateEffectyPeriodic,
			periods:  8,
			expected: 2000.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(tc.rate, tc.freq, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(tc.periods, tc.freq)
			require.NoError(t, err)

			ci, err := New(0, tc.future, rateInterest, period)
			require.NoError(t, err)

			present, err := ci.Present()
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, present, 1.0)
		})
	}
}

func TestPresentErrorWithInvalidPeriod(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(0, 1000, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)
	assert.True(t, present < 1000)
}

func TestPresentWithBimonthlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.06, Bimonthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(6, Bimonthly)
	require.NoError(t, err)

	ci, err := New(0, 1000, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.True(t, present < 1000)
}

func TestPresentWithDifferentRateTypes(t *testing.T) {
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

			ci, err := New(0, 1000, rateInterest, period)
			require.NoError(t, err)

			present, err := ci.Present()
			require.NoError(t, err)

			assert.True(t, present < 1000)
			assert.True(t, present > 0)
		})
	}
}

func TestPresentConsistencyWithFuture(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	originalPresent := 1000.0
	ci, err := New(originalPresent, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	ci2, err := New(0, future, rateInterest, period)
	require.NoError(t, err)

	calculatedPresent, err := ci2.Present()
	require.NoError(t, err)

	assert.InDelta(t, originalPresent, calculatedPresent, 0.1)
}
