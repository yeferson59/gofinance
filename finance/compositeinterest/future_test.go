package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFutureWithPeriodicRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.True(t, future > 1000, "future should be greater than present")
	assert.InDelta(t, 1126.825, future, 0.01)
}

func TestFutureWithNominalRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.12, Monthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.True(t, future > 1000, "future should be greater than present")
	assert.InDelta(t, 1126.825, future, 0.01)
}

func TestFutureWithAnnuallyRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.1268, Annually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(1, Annually)
	require.NoError(t, err)

	ci, err := New(1000, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.True(t, future > 1000, "future should be greater than present")
	assert.InDelta(t, 1126.8, future, 0.1)
}

func TestFutureWithDailyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, Daily, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(365, Daily)
	require.NoError(t, err)

	ci, err := New(1000, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.True(t, future > 1000, "future should be greater than present")
	assert.True(t, future < 1105, "future should be less than 1105 for 10% annual rate")
}

func TestFutureWithQuarterlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.12, QuarterlyOne, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(4, QuarterlyOne)
	require.NoError(t, err)

	ci, err := New(1000, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.True(t, future > 1000, "future should be greater than present")
}

func TestFutureWithSemiAnnuallyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, SemiAnnually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(2, SemiAnnually)
	require.NoError(t, err)

	ci, err := New(5000, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.True(t, future > 5000, "future should be greater than present")
	assert.InDelta(t, 5512.5, future, 100.0)
}

func TestFutureWithZeroPresent(t *testing.T) {
	rateInterest, err := NewRateInterest(0.05, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(0, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.Equal(t, 0.0, future)
}

func TestFutureWithMultipleDataSets(t *testing.T) {
	testCases := []struct {
		name     string
		present  float64
		rate     float64
		freq     CompoundingFrequency
		typeRate TypeRate
		periods  float64
		expected float64
	}{
		{
			name:     "monthly compounding",
			present:  1000,
			rate:     0.01,
			freq:     Monthly,
			typeRate: RateEffectyPeriodic,
			periods:  12,
			expected: 1126.825,
		},
		{
			name:     "annual compounding",
			present:  1000,
			rate:     0.05,
			freq:     Annually,
			typeRate: RateEffectyPeriodic,
			periods:  1,
			expected: 1050.0,
		},
		{
			name:     "quarterly compounding",
			present:  2000,
			rate:     0.02,
			freq:     QuarterlyOne,
			typeRate: RateEffectyPeriodic,
			periods:  8,
			expected: 2342.7,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(tc.rate, tc.freq, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(tc.periods, tc.freq)
			require.NoError(t, err)

			ci, err := New(tc.present, 0, rateInterest, period)
			require.NoError(t, err)

			future, err := ci.Future()
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, future, 1.0)
		})
	}
}

func TestFutureWithBimonthlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.06, Bimonthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(6, Bimonthly)
	require.NoError(t, err)

	ci, err := New(1000, 0, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.True(t, future > 1000)
}

func TestFutureWithDifferentRateTypes(t *testing.T) {
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

			ci, err := New(1000, 0, rateInterest, period)
			require.NoError(t, err)

			future, err := ci.Future()
			require.NoError(t, err)

			assert.True(t, future > 1000)
		})
	}
}
