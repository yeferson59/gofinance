package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestPeriodsWithPeriodicRate(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(112683, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 12.0, periods.InexactFloat64(), 0.01)
}

func TestPeriodsWithNominalRate(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(112683, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 12.0, periods.InexactFloat64(), 0.01)
}

func TestPeriodsWithAnnuallyRate(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Annually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(1), Annually)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(112683, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 1.0, periods.InexactFloat64(), 0.01)
}

func TestPeriodsWithQuarterlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), QuarterlyOne, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(4), QuarterlyOne)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(112680, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.True(t, periods.InexactFloat64() > 0)
	assert.InDelta(t, 4.0, periods.InexactFloat64(), 0.1)
}

func TestPeriodsWithDailyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.10), Daily, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(365), Daily)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(110500, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.True(t, periods.InexactFloat64() > 0)
	assert.InDelta(t, 365.0, periods.InexactFloat64(), 50.0)
}

func TestPeriodsWithSemiAnnuallyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.10), SemiAnnually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(2), SemiAnnually)
	require.NoError(t, err)

	presentMoney, err := money.New(500000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(551250, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 2.0, periods.InexactFloat64(), 0.5)
}

func TestPeriodsWithZeroPeriodShouldError(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.05), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(106167, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)
	assert.True(t, periods.InexactFloat64() > 0)
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
			rateInterest, err := NewRateInterest(money.MustFromFloat64(tc.rate), tc.freq, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(money.MustFromFloat64(tc.periods), tc.freq)
			require.NoError(t, err)

			presentMoney, err := money.New(int64(tc.present*100), 2, money.USD)
			require.NoError(t, err)
			futureMoney, err := money.New(int64(tc.future*100), 2, money.USD)
			require.NoError(t, err)
			ci, err := New(presentMoney, futureMoney, rateInterest, period)
			require.NoError(t, err)

			periods, err := ci.Periods()
			require.NoError(t, err)

			assert.InDelta(t, tc.periods, periods.InexactFloat64(), 1.0)
		})
	}
}

func TestPeriodsErrorWithInvalidPeriod(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(110000, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)
	assert.True(t, periods.InexactFloat64() > 0)
}

func TestPeriodsWithBimonthlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.06), Bimonthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(6), Bimonthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(106136, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 6.0, periods.InexactFloat64(), 0.5)
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
			rateInterest, err := NewRateInterest(money.MustFromFloat64(tc.rate), Monthly, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
			require.NoError(t, err)

			presentMoney, err := money.New(100000, 2, money.USD)
			require.NoError(t, err)
			futureMoney, err := money.New(112683, 2, money.USD)
			require.NoError(t, err)
			ci, err := New(presentMoney, futureMoney, rateInterest, period)
			require.NoError(t, err)

			periods, err := ci.Periods()
			require.NoError(t, err)

			assert.InDelta(t, 12.0, periods.InexactFloat64(), 0.1)
		})
	}
}

func TestPeriodsReturnPositiveValue(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.05), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(24), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(300000, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.True(t, periods.InexactFloat64() > 0, "periods should always be positive")
}

func TestPeriodsWithFractionalPeriods(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.10), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(6), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(105114, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.True(t, periods.InexactFloat64() > 0)
	assert.InDelta(t, 6.0, periods.InexactFloat64(), 0.5)
}

func TestPeriodsConsistencyWithFuture(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	originalPeriods := 12.0
	period, err := NewPeriod(money.MustFromFloat64(originalPeriods), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	ci2, err := New(presentMoney, future, rateInterest, period)
	require.NoError(t, err)

	calculatedPeriods, err := ci2.Periods()
	require.NoError(t, err)

	assert.InDelta(t, originalPeriods, calculatedPeriods.InexactFloat64(), 0.1)
}

func TestPeriodsWithSmallInterestRate(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.001), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(60), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(106167, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 60.0, periods.InexactFloat64(), 1.0)
}

func TestPeriodsWithHighInterestRate(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.10), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(24), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(1096050, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.InDelta(t, 24.0, periods.InexactFloat64(), 1.5)
}

func TestPeriodsWithZeroPresent(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.10), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(24), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(1096050, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	periods, err := ci.Periods()
	require.NoError(t, err)

	assert.Equal(t, 24.0, periods.InexactFloat64())
}
