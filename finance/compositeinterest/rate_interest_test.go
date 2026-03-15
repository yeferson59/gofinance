package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestInterestWithPeriodicRate(t *testing.T) {
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

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.01, interestRate.InexactFloat64(), 0.0001)
}

func TestInterestWithNominalRate(t *testing.T) {
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

	interestRate, err := ci.Interest()
	t.Log(interestRate)
	require.NoError(t, err)

	assert.True(t, interestRate.InexactFloat64() > 0)
	assert.InDelta(t, 0.12, interestRate.InexactFloat64(), 0.0001)
}

func TestInterestWithAnnuallyRate(t *testing.T) {
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

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate.InexactFloat64() > 0)
}

func TestInterestWithQuarterlyCompounding(t *testing.T) {
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

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate.InexactFloat64() > 0)
	assert.True(t, interestRate.InexactFloat64() == 0.12)
}

func TestInterestWithDailyCompounding(t *testing.T) {
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

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate.InexactFloat64() > 0)
	assert.True(t, interestRate.InexactFloat64() == 0.10)
}

func TestInterestWithSemiAnnuallyCompounding(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.10), SemiAnnually, RateEffectyAnnually)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(2), SemiAnnually)
	require.NoError(t, err)

	presentMoney, err := money.New(500000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(605025, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate.InexactFloat64() > 0)
}

func TestInterestWithEqualPresentAndFuture(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.05), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(1), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.05, interestRate.InexactFloat64(), 0.0001)
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

			interestRate, err := ci.Interest()
			require.NoError(t, err)

			assert.True(t, interestRate.InexactFloat64() > 0)
			assert.InDelta(t, tc.rate, interestRate.InexactFloat64(), 0.001)
		})
	}
}

func TestInterestErrorWithInvalidPeriod(t *testing.T) {
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

	interestRate, err := ci.Interest()
	require.NoError(t, err)
	assert.True(t, interestRate.InexactFloat64() > 0)
}

func TestInterestWithBimonthlyCompounding(t *testing.T) {
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

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate.InexactFloat64() > 0)
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

			interestRate, err := ci.Interest()
			require.NoError(t, err)

			assert.True(t, interestRate.InexactFloat64() > 0)
		})
	}
}

func TestInterestFutureGreaterThanPresent(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.05), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
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
			presentMoney, err := money.New(int64(tc.present*100), 2, money.USD)
			require.NoError(t, err)
			futureMoney, err := money.New(int64(tc.future*100), 2, money.USD)
			require.NoError(t, err)
			ci, err := New(presentMoney, futureMoney, rateInterest, period)
			require.NoError(t, err)

			interestRate, err := ci.Interest()
			require.NoError(t, err)

			assert.True(t, interestRate.InexactFloat64() > 0, "interest rate should be positive when future > present")
		})
	}
}

func TestInterestNegativeWhenFutureLessThanPresent(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.05), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	presentMoney, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(50000, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.True(t, interestRate.InexactFloat64() == 0.05, "interest rate should be negative when future < present")
}
