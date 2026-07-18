package compoundinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/money"
)

func TestPresentWithPeriodicRate(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(112683, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.InDelta(t, 1000.0, present.ToDecimal().InexactFloat64(), 0.1)
}

func TestPresentWithNominalRate(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(112683, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.InDelta(t, 1000.0, present.ToDecimal().InexactFloat64(), 0.1)
}

func TestPresentWithAnnuallyRate(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Annually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(1), Annually)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(112683, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	assert.InDelta(t, 1000.0, present.ToDecimal().InexactFloat64(), 0.1)
}

func TestPresentWithDailyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.10), Daily, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(365), Daily)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(110500, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	// 10% effective annual compounded daily over 365 days is exactly one year:
	// PV = 1105 / 1.10 = 1004.5455
	assert.InDelta(t, 1004.5455, present.ToDecimal().InexactFloat64(), 0.01)
}

func TestPresentWithQuarterlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), QuarterlyOne, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(4), QuarterlyOne)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	// 12% nominal quarterly => i = 0.03, PV = 1000 / 1.03^4 = 888.4870
	assert.InDelta(t, 888.4870, present.ToDecimal().InexactFloat64(), 0.01)
}

func TestPresentWithSemiAnnuallyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.10), SemiAnnually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(2), SemiAnnually)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(551250, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	// 10% effective annual compounded semi-annually over 2 periods is exactly
	// one year: PV = 5512.50 / 1.10 = 5011.3636 (5000 would only be correct
	// for a 10% nominal rate, not the effective annual rate used here).
	assert.InDelta(t, 5011.3636, present.ToDecimal().InexactFloat64(), 0.01)
}

func TestPresentWithZeroFuture(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.05), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.Error(t, err)

	assert.Equal(t, 0.0, present.ToDecimal().InexactFloat64())
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
			// PV = 2342.70 / 1.02^8 = 1999.4719
			expected: 1999.4719,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(decimal.MustFromFloat64(tc.rate), tc.freq, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(decimal.MustFromFloat64(tc.periods), tc.freq)
			require.NoError(t, err)

			presentMoney, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			futureMoney, err := money.New(int64(tc.future*100), 2, money.USD)
			require.NoError(t, err)
			ci, err := New(presentMoney, futureMoney, rateInterest, period)
			require.NoError(t, err)

			present, err := ci.Present()
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, present.ToDecimal().InexactFloat64(), 0.01)
		})
	}
}

func TestPresentWithTwelveMonthlyPeriods(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)
	// PV = 1000 / 1.01^12 = 887.4492
	assert.InDelta(t, 887.4492, present.ToDecimal().InexactFloat64(), 0.01)
}

func TestPresentWithBimonthlyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.06), Bimonthly, RateEffectyNominal)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(6), Bimonthly)
	require.NoError(t, err)

	presentMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	present, err := ci.Present()
	require.NoError(t, err)

	// 6% nominal bimonthly => i = 0.01, PV = 1000 / 1.01^6 = 942.0452
	assert.InDelta(t, 942.0452, present.ToDecimal().InexactFloat64(), 0.01)
}

func TestPresentWithDifferentRateTypes(t *testing.T) {
	testCases := []struct {
		name     string
		rate     float64
		typeRate TypeRate
		expected float64
	}{
		// periodic 1% monthly: PV = 1000 / 1.01^12 = 887.4492
		{"periodic", 0.01, RateEffectyPeriodic, 887.4492},
		// 12% nominal monthly => i = 0.01: same as above
		{"nominal", 0.12, RateEffectyNominal, 887.4492},
		// 12.68% effective annual over 12 months = exactly one year:
		// PV = 1000 / 1.1268 = 887.4689
		{"annual", 0.1268, RateEffectyAnnually, 887.4689},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(decimal.MustFromFloat64(tc.rate), Monthly, tc.typeRate)
			require.NoError(t, err)

			period, err := NewPeriod(decimal.MustFromFloat64(12), Monthly)
			require.NoError(t, err)

			presentMoney, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			futureMoney, err := money.New(100000, 2, money.USD)
			require.NoError(t, err)
			ci, err := New(presentMoney, futureMoney, rateInterest, period)
			require.NoError(t, err)

			present, err := ci.Present()
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, present.ToDecimal().InexactFloat64(), 0.01)
		})
	}
}

func TestPresentConsistencyWithFuture(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(decimal.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	originalPresent := 1000.0
	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	future, err := ci.Future()
	require.NoError(t, err)

	presentMoney2, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	ci2, err := New(presentMoney2, future, rateInterest, period)
	require.NoError(t, err)

	calculatedPresent, err := ci2.Present()
	require.NoError(t, err)

	assert.InDelta(t, originalPresent, calculatedPresent.ToDecimal().InexactFloat64(), 0.1)
}
