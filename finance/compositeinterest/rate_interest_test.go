package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

// newInterestCase builds a CompositeInterest whose rate value is zero so that
// Interest() is forced to derive the periodic rate from present, future, and
// periods via r = (FV/PV)^(1/n) - 1 instead of returning the configured rate.
func newInterestCase(t *testing.T, present, future int64, freq CompoundingFrequency, typeRate TypeRate, periods float64) CompositeInterest {
	t.Helper()

	rateInterest, err := NewRateInterest(money.MustFromFloat64(0), freq, typeRate)
	require.NoError(t, err)

	period, err := NewPeriod(money.MustFromFloat64(periods), freq)
	require.NoError(t, err)

	presentMoney, err := money.New(present, 2, money.USD)
	require.NoError(t, err)
	futureMoney, err := money.New(future, 2, money.USD)
	require.NoError(t, err)
	ci, err := New(presentMoney, futureMoney, rateInterest, period)
	require.NoError(t, err)

	return ci
}

func TestInterestReturnsConfiguredRate(t *testing.T) {
	// When a non-zero rate is configured, Interest() returns its raw value
	// directly without deriving anything from present/future/periods.
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

	assert.Equal(t, 0.01, interestRate.InexactFloat64())
}

func TestInterestWithPeriodicRate(t *testing.T) {
	// r = (1126.83/1000)^(1/12) - 1 = 0.0100004
	ci := newInterestCase(t, 100000, 112683, Monthly, RateEffectyPeriodic, 12)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.0100004, interestRate.InexactFloat64(), 0.0000001)
}

func TestInterestWithNominalRate(t *testing.T) {
	// The rate type does not change the derivation:
	// r = (1126.83/1000)^(1/12) - 1 = 0.0100004
	ci := newInterestCase(t, 100000, 112683, Monthly, RateEffectyNominal, 12)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.0100004, interestRate.InexactFloat64(), 0.0000001)
}

func TestInterestWithAnnuallyRate(t *testing.T) {
	// r = (1126.83/1000)^(1/1) - 1 = 0.12683
	ci := newInterestCase(t, 100000, 112683, Annually, RateEffectyAnnually, 1)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.12683, interestRate.InexactFloat64(), 0.0000001)
}

func TestInterestWithQuarterlyCompounding(t *testing.T) {
	// r = (1126.80/1000)^(1/4) - 1 = 0.0302953
	ci := newInterestCase(t, 100000, 112680, QuarterlyOne, RateEffectyNominal, 4)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.0302953, interestRate.InexactFloat64(), 0.0000001)
}

func TestInterestWithDailyCompounding(t *testing.T) {
	// r = (1105/1000)^(1/365) - 1 = 0.0002736 per day
	ci := newInterestCase(t, 100000, 110500, Daily, RateEffectyAnnually, 365)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.0002736, interestRate.InexactFloat64(), 0.0000001)
}

func TestInterestWithSemiAnnuallyCompounding(t *testing.T) {
	// r = (6050.25/5000)^(1/2) - 1 = 0.1000227 per semester
	ci := newInterestCase(t, 500000, 605025, SemiAnnually, RateEffectyAnnually, 2)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.1000227, interestRate.InexactFloat64(), 0.0000001)
}

func TestInterestWithEqualPresentAndFuture(t *testing.T) {
	// If the money does not grow, the derived rate is exactly zero:
	// r = (1000/1000)^(1/1) - 1 = 0
	ci := newInterestCase(t, 100000, 100000, Monthly, RateEffectyPeriodic, 1)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.Equal(t, 0.0, interestRate.InexactFloat64())
}

func TestInterestWithMultipleDataSets(t *testing.T) {
	testCases := []struct {
		name     string
		present  float64
		future   float64
		freq     CompoundingFrequency
		typeRate TypeRate
		periods  float64
		expected float64
	}{
		{
			name:     "monthly compounding",
			present:  1000,
			future:   1126.825,
			freq:     Monthly,
			typeRate: RateEffectyPeriodic,
			periods:  12,
			// 1126.825 truncates to 1126.82 when stored as cents:
			// r = 1.12682^(1/12) - 1 = 0.0099996
			expected: 0.0099996243,
		},
		{
			name:     "annual compounding",
			present:  1000,
			future:   1050.0,
			freq:     Annually,
			typeRate: RateEffectyPeriodic,
			periods:  1,
			// r = 1.05 - 1 = 0.05
			expected: 0.05,
		},
		{
			name:     "quarterly compounding",
			present:  2000,
			future:   2342.7,
			freq:     QuarterlyOne,
			typeRate: RateEffectyPeriodic,
			periods:  8,
			// 2342.70 loses a cent to float truncation when stored (2342.69):
			// r = (2342.69/2000)^(1/8) - 1 = 0.0199658
			expected: 0.0199657850,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ci := newInterestCase(t, int64(tc.present*100), int64(tc.future*100), tc.freq, tc.typeRate, tc.periods)

			interestRate, err := ci.Interest()
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, interestRate.InexactFloat64(), 0.0000001)
		})
	}
}

func TestInterestWithTenPercentGrowth(t *testing.T) {
	// r = (1100/1000)^(1/12) - 1 = 0.0079741
	ci := newInterestCase(t, 100000, 110000, Monthly, RateEffectyPeriodic, 12)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.0079741, interestRate.InexactFloat64(), 0.0000001)
}

func TestInterestWithBimonthlyCompounding(t *testing.T) {
	// r = (1061.36/1000)^(1/6) - 1 = 0.0099746
	ci := newInterestCase(t, 100000, 106136, Bimonthly, RateEffectyNominal, 6)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, 0.0099746, interestRate.InexactFloat64(), 0.0000001)
}

func TestInterestWithDifferentRateTypes(t *testing.T) {
	// The configured rate type does not affect the derivation formula when
	// the rate value is zero: all cases share r = (FV/PV)^(1/n) - 1.
	testCases := []struct {
		name     string
		typeRate TypeRate
	}{
		{"periodic", RateEffectyPeriodic},
		{"nominal", RateEffectyNominal},
		{"annual", RateEffectyAnnually},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ci := newInterestCase(t, 100000, 112683, Monthly, tc.typeRate, 12)

			interestRate, err := ci.Interest()
			require.NoError(t, err)

			// r = (1126.83/1000)^(1/12) - 1 = 0.0100004
			assert.InDelta(t, 0.0100004, interestRate.InexactFloat64(), 0.0000001)
		})
	}
}

func TestInterestFutureGreaterThanPresent(t *testing.T) {
	testCases := []struct {
		name     string
		present  float64
		future   float64
		expected float64
	}{
		// r = (110/100)^(1/12) - 1
		{"small amounts", 100, 110, 0.0079741},
		// r = (1200/1000)^(1/12) - 1
		{"medium amounts", 1000, 1200, 0.0153095},
		// r = (15000/10000)^(1/12) - 1
		{"large amounts", 10000, 15000, 0.0343661},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ci := newInterestCase(t, int64(tc.present*100), int64(tc.future*100), Monthly, RateEffectyPeriodic, 12)

			interestRate, err := ci.Interest()
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, interestRate.InexactFloat64(), 0.0000001)
		})
	}
}

func TestInterestNegativeWhenFutureLessThanPresent(t *testing.T) {
	// r = (500/1000)^(1/12) - 1 = -0.0561257: the derived periodic rate is
	// negative when the money shrinks.
	ci := newInterestCase(t, 100000, 50000, Monthly, RateEffectyPeriodic, 12)

	interestRate, err := ci.Interest()
	require.NoError(t, err)

	assert.InDelta(t, -0.0561257, interestRate.InexactFloat64(), 0.0000001)
}
