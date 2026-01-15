package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterestWithPeriodicRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1126.825, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.01, interestRate, 0.0001)
}

func TestInterestWithNominalRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.12, Monthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1126.825, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	t.Log(interestRate)
	require.NoError(t, err)

	assert.True(t, interestRate > 0)
	assert.InDelta(t, 0.12, interestRate, 0.0001)
}

func TestInterestWithAnnuallyRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.1268, Annually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(1, Annually)
	require.NoError(t, err)

	ci, err := New(1000, 1126.825, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate > 0)
}

func TestInterestWithQuarterlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.12, QuarterlyOne, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(4, QuarterlyOne)
	require.NoError(t, err)

	ci, err := New(1000, 1126.8, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate > 0)
	assert.True(t, interestRate == 0.12)
}

func TestInterestWithDailyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, Daily, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(365, Daily)
	require.NoError(t, err)

	ci, err := New(1000, 1105, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate > 0)
	assert.True(t, interestRate == 0.10)
}

func TestInterestWithSemiAnnuallyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.10, SemiAnnually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(2, SemiAnnually)
	require.NoError(t, err)

	ci, err := New(5000, 6050.25, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate > 0)
}

func TestInterestWithEqualPresentAndFuture(t *testing.T) {
	rateInterest, err := NewRateInterest(0.05, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(1, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1000, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.05, interestRate, 0.0001)
}

func TestInterestWithMultipleDataSets(t *testing.T) {
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

			interestRate, err := ci.Interest()
			require.NoError(t, err)

			assert.True(t, interestRate > 0)
			assert.InDelta(t, tc.rate, interestRate, 0.001)
		})
	}
}

func TestInterestErrorWithInvalidPeriod(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 1100, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)
	assert.True(t, interestRate > 0)
}

func TestInterestWithBimonthlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(0.06, Bimonthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(6, Bimonthly)
	require.NoError(t, err)

	ci, err := New(1000, 1061.36, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate > 0)
}

func TestInterestWithDifferentRateTypes(t *testing.T) {
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

			interestRate, err := ci.Interest()
			require.NoError(t, err)

			assert.True(t, interestRate > 0)
		})
	}
}

func TestInterestFutureGreaterThanPresent(t *testing.T) {
	rateInterest, err := NewRateInterest(0.05, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		present float64
		future  float64
	}{
		{"small amounts", 100, 110},
		{"medium amounts", 1000, 1200},
		{"large amounts", 10000, 15000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ci, err := New(tc.present, tc.future, rateInterest, period)
			require.NoError(t, err)

			interestRate, err := ci.Interest()
			require.NoError(t, err)

			assert.True(t, interestRate > 0, "interest rate should be positive when future > present")
		})
	}
}

func TestInterestNegativeWhenFutureLessThanPresent(t *testing.T) {
	rateInterest, err := NewRateInterest(0.05, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	ci, err := New(1000, 500, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate == 0.05, "interest rate should be negative when future < present")
}
