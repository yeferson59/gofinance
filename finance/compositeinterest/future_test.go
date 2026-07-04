package compositeinterest

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestFutureWithPeriodicRate(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)

	tx := require.New(t)

	tx.NoError(err)

	period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
	tx.NoError(err)

	present, err := money.New(100000, 2, money.USD)
	tx.NoError(err)
	futureValue, err := money.New(0, 2, money.USD)
	tx.NoError(err)
	ci, err := New(present, futureValue, rateInterest, period)
	tx.NoError(err)

	future, err := ci.Future()
	tx.NoError(err)

	tx.True(future.ToDecimal().InexactFloat64() > 1000, "future should be greater than present")
	tx.InDelta(1126.825, future.ToDecimal().InexactFloat64(), 0.01)
}

func TestFutureWithNominalRate(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	present, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureValue, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(present, futureValue, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.True(t, future.ToDecimal().InexactFloat64() > 1000, "future should be greater than present")
	assert.InDelta(t, 1126.825, future.ToDecimal().InexactFloat64(), 0.01)
}

func TestFutureWithAnnuallyRate(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Annually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(1), Annually)
	require.NoError(t, err)

	present, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureValue, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(present, futureValue, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	assert.True(t, future.ToDecimal().InexactFloat64() > 1000, "future should be greater than present")
	assert.InDelta(t, 1126.8, future.ToDecimal().InexactFloat64(), 0.1)
}

func TestFutureWithDailyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.10), Daily, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(365), Daily)
	require.NoError(t, err)

	present, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureValue, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(present, futureValue, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	// 10% effective annual compounded daily over 365 days is exactly one year:
	// FV = 1000 × 1.10 = 1100
	assert.InDelta(t, 1100.0, future.ToDecimal().InexactFloat64(), 0.01)
}

func TestFutureWithQuarterlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), QuarterlyOne, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(4), QuarterlyOne)
	require.NoError(t, err)

	present, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureValue, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(present, futureValue, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	// 12% nominal quarterly => i = 0.03, FV = 1000 × 1.03^4 = 1125.5088
	assert.InDelta(t, 1125.5088, future.ToDecimal().InexactFloat64(), 0.01)
}

func TestFutureWithSemiAnnuallyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.10), SemiAnnually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(2), SemiAnnually)
	require.NoError(t, err)

	present, err := money.New(500000, 2, money.USD)
	require.NoError(t, err)
	futureValue, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(present, futureValue, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	// 10% effective annual compounded semi-annually over 2 periods is exactly
	// one year: FV = 5000 × 1.10 = 5500 (not 5512.50, which would be a 10%
	// nominal rate instead of the effective annual rate used here).
	assert.InDelta(t, 5500.0, future.ToDecimal().InexactFloat64(), 0.01)
}

func TestFutureWithZeroPresent(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.05), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureValue, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(present, futureValue, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.Error(t, err)

	assert.Equal(t, 0.0, future.ToDecimal().InexactFloat64())
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
			// FV = 2000 × 1.02^8 = 2343.3188
			expected: 2343.3188,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(money.MustFromFloat64(tc.rate), tc.freq, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(money.MustFromFloat64(tc.periods), tc.freq)
			require.NoError(t, err)

			presentMoney, err := money.New(int64(tc.present*100), 2, money.USD)
			require.NoError(t, err)
			futureValue, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			ci, err := New(presentMoney, futureValue, rateInterest, period)
			require.NoError(t, err)

			future, err := ci.Future()
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, future.ToDecimal().InexactFloat64(), 0.01)
		})
	}
}

func TestFutureWithBimonthlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.06), Bimonthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(6), Bimonthly)
	require.NoError(t, err)

	present, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureValue, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(present, futureValue, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	// 6% nominal bimonthly => i = 0.01, FV = 1000 × 1.01^6 = 1061.5202
	assert.InDelta(t, 1061.5202, future.ToDecimal().InexactFloat64(), 0.01)
}

func TestFutureWithDifferentRateTypes(t *testing.T) {
	testCases := []struct {
		name     string
		rate     float64
		typeRate TypeRate
		expected float64
	}{
		// periodic 1% monthly: FV = 1000 × 1.01^12 = 1126.8250
		{"periodic", 0.01, RateEffectyPeriodic, 1126.8250},
		// 12% nominal monthly => i = 0.01: same as above
		{"nominal", 0.12, RateEffectyNominal, 1126.8250},
		// 12.68% effective annual over 12 months = exactly one year:
		// FV = 1000 × 1.1268 = 1126.80
		{"annual", 0.1268, RateEffectyAnnually, 1126.80},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(money.MustFromFloat64(tc.rate), Monthly, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
			require.NoError(t, err)

			present, err := money.New(100000, 2, money.USD)
			require.NoError(t, err)
			futureValue, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			ci, err := New(present, futureValue, rateInterest, period)
			require.NoError(t, err)

			future, err := ci.Future()
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, future.ToDecimal().InexactFloat64(), 0.01)
		})
	}
}

func createCompositeInterest(numberPeriod float64, periodF CompoundingFrequency, rateInterest float64, compounding CompoundingFrequency, typeRate TypeRate, present, future float64, tx *require.Assertions) CompositeInterest {
	period, err := NewPeriod(money.MustFromFloat64(numberPeriod), periodF)
	tx.NoError(err)

	rate, err := NewRateInterest(money.MustFromFloat64(rateInterest), compounding, typeRate)
	tx.NoError(err)

	presentMoney, err := money.New(int64(present*100), 2, money.USD)
	tx.NoError(err)
	futureMoney, err := money.New(int64(future*100), 2, money.USD)
	tx.NoError(err)
	cr, err := New(presentMoney, futureMoney, rate, period)
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

		tx.InDelta(expected, math.Round(future.ToDecimal().InexactFloat64()), 0.00001)
	})

	t.Run("run simple operation with monthly compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with bimonth compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, Bimonthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_673_902.24980029

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with quarterly one compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, QuarterlyOne, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_658_908.265198686

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with quarterly two compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, QuarterlyTwo, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_644_246.304611462

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with semi annually compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, SemiAnnually, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_615_871.864700002

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with annually compounding frequency for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(30, Monthly, 0.18, Annually, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_537_600.666784153

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with day period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(912.5, Daily, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_718_420.985860214

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with bimonthly period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(15, Bimonthly, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with quarterly one period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(10, QuarterlyOne, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with quarterly two period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(7.5, QuarterlyTwo, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with semi-annually period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(5, SemiAnnually, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})

	t.Run("run simple operation with annually period for future", func(t *testing.T) {
		tx := require.New(t)

		cr := createCompositeInterest(2.5, Annually, 0.18, Monthly, RateEffectyNominal, 3_000_000, 0, tx)

		future, err := cr.Future()
		tx.NoError(err)

		expected := 4_689_240.66147257

		tx.InDelta(expected, future.ToDecimal().InexactFloat64(), 0.00001)
	})
}
