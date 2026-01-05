package compositeinterest

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFutureWithPeriodicRate(t *testing.T) {
	rateInterest, err := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)

	tx := require.New(t)

	tx.NoError(err)

	period, err := NewPeriod(12, Monthly)
	tx.NoError(err)

	ci, err := New(1000, 0, rateInterest, period)
	tx.NoError(err)

	future, err := ci.Future()
	tx.NoError(err)

	tx.True(future > 1000, "future should be greater than present")
	tx.InDelta(1126.825, future, 0.01)
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

func createCompositeInterest(numberPeriod float64, periodF CompoundingFrequency, rateInterest float64, compounding CompoundingFrequency, typeRate TypeRate, present, future float64, tx *require.Assertions) *CompositeInterest {
	period, err := NewPeriod(numberPeriod, periodF)
	tx.NoError(err)

	rate, err := NewRateInterest(rateInterest, compounding, typeRate)
	tx.NoError(err)

	cr, err := New(present, future, rate, period)
	tx.NoError(err)

	return cr
}

func TestMoreExampleFuture(t *testing.T) {
	t.Run("run simple operation for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(2.5, Annually, 0.09, Monthly, RateEffectyNominal, 4_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 5_005_087.00

		tx.InDelta(expected, math.Round(future), 0.00001)
	})

	t.Run("run simple operation with monthly compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with bimonth compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, Bimonthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_673_902.24980029

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with quarterly one compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, QuarterlyOne, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_658_908.265198686

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with quarterly two compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, QuarterlyTwo, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_644_246.304611462

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with semi annually compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, SemiAnnually, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_615_871.864700002

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with annually compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, Annually, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_537_600.666784153

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with bimonthly period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(15, Bimonthly, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with quarterly one period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(10, QuarterlyOne, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with quarterly two period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(7.5, QuarterlyTwo, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with semi-annually period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(5, SemiAnnually, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future, 0.00001)
	})

	t.Run("run simple operation with annually period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(2.5, Annually, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future, 0.00001)
	})
}
