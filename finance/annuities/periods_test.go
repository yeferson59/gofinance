package annuities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

func TestAnnuityPeriodsWithPresent(t *testing.T) {
	// PMT = 1000, PV = 10000, i = 0.01
	// n = ln(PMT / (PMT - PV×i)) / ln(1+i) = ln(1000/900) / ln(1.01) = 10.5886
	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.12), compositeinterest.Monthly, compositeinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	periods, err := annuity.PeriodsWithPresent()
	require.NoError(t, err)
	assert.InDelta(t, 10.5886, periods.InexactFloat64(), 0.0001)
}

func TestAnnuityPeriodsWithFuture(t *testing.T) {
	// PMT = 1000, FV = 15000, i = 0.01
	// n = ln((FV×i + PMT) / PMT) / ln(1+i) = ln(1.15) / ln(1.01) = 14.0460
	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.12), compositeinterest.Monthly, compositeinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	periods, err := annuity.PeriodsWithFuture()
	require.NoError(t, err)
	assert.InDelta(t, 14.0460, periods.InexactFloat64(), 0.0001)
}

func TestAnnuityAnticipatePeriodsWithPresent(t *testing.T) {
	// PMT = 1000, PV = 10000, i = 0.01
	// n = ln(PMT / (PMT - [PV/(1+i)]×i)) / ln(1+i) = 10.4781
	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.12), compositeinterest.Monthly, compositeinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	periods, err := annuity.AnticipatePeriodsWithPresent()
	require.NoError(t, err)
	assert.InDelta(t, 10.4781, periods.InexactFloat64(), 0.0001)
}

func TestAnnuityAnticipatePeriodsWithFuture(t *testing.T) {
	// PMT = 1000, FV = 15000, i = 0.01
	// n = ln(([FV/(1+i)]×i + PMT) / PMT) / ln(1+i) = 13.9161
	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.12), compositeinterest.Monthly, compositeinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	periods, err := annuity.AnticipatePeriodsWithFuture()
	require.NoError(t, err)
	assert.InDelta(t, 13.9161, periods.InexactFloat64(), 0.0001)
}

func TestAnnuityPeriodsWithPresentZeroRate(t *testing.T) {
	// With rate = 0, ln(1+i) = ln(1) = 0, so PeriodsWithPresent must return
	// an error instead of dividing by that zero.
	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.0), compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	_, err = annuity.PeriodsWithPresent()
	assert.Error(t, err)
}

func TestAnnuityPeriodsWithFutureZeroRate(t *testing.T) {
	// With rate = 0, ln(1+i) = ln(1) = 0, so PeriodsWithFuture must return
	// an error instead of dividing by that zero.
	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.0), compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value, err := money.New(100000, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	_, err = annuity.PeriodsWithFuture()
	assert.Error(t, err)
}

func TestAnnuityPeriodsWithPresentAcrossFrequencies(t *testing.T) {
	// PMT = 1000, 12% nominal at each frequency, n = 12 periods of that
	// frequency. Present values are the same ones verified in
	// TestAnnuityWithDifferentCompoundingFrequencies, so recovering the
	// periods from them must round-trip back to 12.
	testCases := []struct {
		name      string
		frequency compositeinterest.CompoundingFrequency
		present   float64
	}{
		{"daily", compositeinterest.Daily, 11974.3955},
		{"monthly", compositeinterest.Monthly, 11255.0775},
		{"quarterly", compositeinterest.QuarterlyOne, 9954.0040},
		{"annually", compositeinterest.Annually, 6194.3742},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), tc.frequency)
			require.NoError(t, err)

			rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.12), tc.frequency, compositeinterest.RateEffectyNominal)
			require.NoError(t, err)

			value, err := money.New(100000, 2, money.USD)
			require.NoError(t, err)
			present := money.MustMoneyFromFloat64(tc.present, money.USD)
			future, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			annuity, err := New(value, present, future, period, rateInterest)
			require.NoError(t, err)

			periods, err := annuity.PeriodsWithPresent()
			require.NoError(t, err)
			assert.InDelta(t, 12.0, periods.InexactFloat64(), 0.01)
		})
	}
}

func TestAnnuityPeriodsWithFutureAcrossFrequencies(t *testing.T) {
	// PMT = 1000, 12% nominal at each frequency, n = 12 periods of that
	// frequency. Future values are FV = PMT × [(1+i)^12 - 1]/i, so recovering
	// the periods from them must round-trip back to 12.
	testCases := []struct {
		name      string
		frequency compositeinterest.CompoundingFrequency
		future    float64
	}{
		{"daily", compositeinterest.Daily, 12021.7224},
		{"monthly", compositeinterest.Monthly, 12682.5030},
		{"quarterly", compositeinterest.QuarterlyOne, 14192.0296},
		{"annually", compositeinterest.Annually, 24133.1333},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), tc.frequency)
			require.NoError(t, err)

			rateInterest, err := compositeinterest.NewRateInterest(money.MustFromFloat64(0.12), tc.frequency, compositeinterest.RateEffectyNominal)
			require.NoError(t, err)

			value, err := money.New(100000, 2, money.USD)
			require.NoError(t, err)
			present, err := money.New(0, 2, money.USD)
			require.NoError(t, err)
			future := money.MustMoneyFromFloat64(tc.future, money.USD)
			annuity, err := New(value, present, future, period, rateInterest)
			require.NoError(t, err)

			periods, err := annuity.PeriodsWithFuture()
			require.NoError(t, err)
			assert.InDelta(t, 12.0, periods.InexactFloat64(), 0.01)
		})
	}
}
