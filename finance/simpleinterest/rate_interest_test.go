package simpleinterest

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestRateInterest(t *testing.T) {
	numPeriods, _ := decimal.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	present, _ := money.New(5_000, 0, money.COP)
	interest, _ := money.New(500, 0, money.COP)

	simpleInterest := New(money.Money{}, present, interest, decimal.Decimal{}, periods)
	expectedRate, _ := decimal.NewFromFloat64(500.0 / (5000.0 * 2.0))

	rate, err := simpleInterest.RateInterest()
	require.NoError(t, err)

	assert.Equal(t, expectedRate.String(), rate.String())

	// Test error case: present=0
	zero, _ := money.New(0, 0, money.COP)
	simpleInterest.present = zero
	_, err = simpleInterest.RateInterest()
	assert.Error(t, err)

	// Test error case: periods=0
	simpleInterest.present = present
	zeroPeriod, _ := decimal.NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.RateInterest()
	assert.Error(t, err)
}

func TestRateInterestPropagatesPeriodError(t *testing.T) {
	// An empty Period has no valid time unit set, so getPeriod fails and
	// RateInterest must surface that error.
	present, _ := money.New(5_000, 0, money.COP)
	interest, _ := money.New(500, 0, money.COP)

	simpleInterest := New(money.Money{}, present, interest, decimal.Decimal{}, Period{})

	_, err := simpleInterest.RateInterest()
	assert.Error(t, err)
}

func TestRateInterestPropagatesOverflowFromDivision(t *testing.T) {
	// interest / (present × periods) overflows decimal128's 128-bit
	// coefficient when interest is astronomically larger than present.
	numPeriods, _ := decimal.NewFromInt64(1, 0)
	periods := NewPeriod(numPeriods, Days)

	present, _ := money.New(1, 19, money.COP)
	interest, _ := money.New(math.MaxInt64, 0, money.COP)

	simpleInterest := New(money.Money{}, present, interest, decimal.Decimal{}, periods)

	_, err := simpleInterest.RateInterest()
	assert.Error(t, err)
}

func TestRateInterestWithPresentAndFuture(t *testing.T) {
	numPeriods, _ := decimal.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	future, _ := money.New(5_500, 0, money.COP)
	present, _ := money.New(5_000, 0, money.COP)

	simpleInterest := New(future, present, money.Money{}, decimal.Decimal{}, periods)

	rate, err := simpleInterest.RateInterestWithPresentAndFuture()
	require.NoError(t, err)

	// r = (FV/PV - 1) / n = (5500/5000 - 1) / 2 = 0.05
	assert.Equal(t, 0.05, rate.InexactFloat64())

	// Test error case: present=0
	zeroValue, _ := money.New(0, 0, money.COP)
	simpleInterest.present = zeroValue
	_, err = simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)

	// Test error case: periods=0
	simpleInterest.present = present
	zeroPeriod, _ := decimal.NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)
}

func TestRateInterestWithPresentAndFuturePropagatesPeriodError(t *testing.T) {
	// An empty Period has no valid time unit set, so getPeriod fails and
	// RateInterestWithPresentAndFuture must surface that error.
	future, _ := money.New(5_500, 0, money.COP)
	present, _ := money.New(5_000, 0, money.COP)

	simpleInterest := New(future, present, money.Money{}, decimal.Decimal{}, Period{})

	_, err := simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)
}

func TestRateInterestWithPresentAndFuturePropagatesOverflowFromRatio(t *testing.T) {
	// future/present computed at an extreme magnitude mismatch overflows
	// decimal128's 128-bit coefficient.
	numPeriods, _ := decimal.NewFromInt64(1, 0)
	periods := NewPeriod(numPeriods, Days)

	present, _ := money.New(1, 19, money.COP)
	future, _ := money.New(math.MaxInt64, 0, money.COP)

	simpleInterest := New(future, present, money.Money{}, decimal.Decimal{}, periods)

	_, err := simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)
}

func TestRateInterestWithPresentAndFuturePropagatesOverflowFromDivision(t *testing.T) {
	// (future/present - 1) divided by an astronomically tiny number of
	// periods overflows decimal128's 128-bit coefficient.
	tinyPeriods, err := decimal.NewFromString("0.0000000000000000001")
	require.NoError(t, err)
	periods := NewPeriod(tinyPeriods, Days)

	future, _ := money.New(1_000_000, 2, money.COP)
	present, _ := money.New(1_000, 2, money.COP)

	simpleInterest := New(future, present, money.Money{}, decimal.Decimal{}, periods)

	_, err = simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)
}
