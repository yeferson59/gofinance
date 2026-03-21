package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestRateInterest(t *testing.T) {
	numPeriods, _ := money.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	present, _ := money.New(5_000, 0, money.COP)
	interest, _ := money.New(500, 0, money.COP)

	simpleInterest := New(money.Money{}, present, interest, money.Decimal{}, periods)
	expectedRate, _ := money.NewFromFloat64(500.0 / (5000.0 * 2.0))

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
	zeroPeriod, _ := money.NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.RateInterest()
	assert.Error(t, err)
}

func TestRateInterestWithPresentAndFuture(t *testing.T) {
	numPeriods, _ := money.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	future, _ := money.New(5_500, 0, money.COP)
	present, _ := money.New(5_000, 0, money.COP)

	simpleInterest := New(future, present, money.Money{}, money.Decimal{}, periods)

	rate, err := simpleInterest.RateInterestWithPresentAndFuture()
	require.NoError(t, err)

	// Check that the rate is approximately 0.05
	zero, _ := money.NewFromInt64(0, 0)
	maxValue, _ := money.NewFromFloat64(0.06)
	assert.True(t, rate.Cmp(zero) >= 0, "rate should be >= 0")
	assert.True(t, rate.Cmp(maxValue) <= 0, "rate should be <= 0.06")

	// Test error case: present=0
	zeroValue, _ := money.New(0, 0, money.COP)
	simpleInterest.present = zeroValue
	_, err = simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)

	// Test error case: periods=0
	simpleInterest.present = present
	zeroPeriod, _ := money.NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)
}
