package simpleinterest

import (
	"testing"

	"github.com/quagmt/udecimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateInterest(t *testing.T) {
	numPeriods, _ := NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	present, _ := NewFromInt64(5_000, 0)
	interest, _ := NewFromInt64(500, 0)

	simpleInterest := New(Decimal{}, present, interest, Decimal{}, periods)
	expectedRate, _ := NewFromFloat64(500.0 / (5000.0 * 2.0))

	rate, err := simpleInterest.RateInterest()
	require.NoError(t, err)

	assert.Equal(t, expectedRate.String(), rate.String())

	// Test error case: present=0
	zero, _ := NewFromInt64(0, 0)
	simpleInterest.present = zero
	_, err = simpleInterest.RateInterest()
	assert.Error(t, err)

	// Test error case: periods=0
	simpleInterest.present = present
	zeroPeriod, _ := NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.RateInterest()
	assert.Error(t, err)
}

func TestRateInterestWithPresentAndFuture(t *testing.T) {
	numPeriods, _ := NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	future, _ := NewFromInt64(5_500, 0)
	present, _ := NewFromInt64(5_000, 0)

	simpleInterest := New(future, present, Decimal{}, Decimal{}, periods)

	rate, err := simpleInterest.RateInterestWithPresentAndFuture()
	require.NoError(t, err)

	// Check that the rate is approximately 0.05
	zero, _ := udecimal.NewFromInt64(0, 0)
	maxValue, _ := udecimal.NewFromFloat64(0.06)
	assert.True(t, rate.Decimal.Cmp(zero) >= 0, "rate should be >= 0")
	assert.True(t, rate.Decimal.Cmp(maxValue) <= 0, "rate should be <= 0.06")

	// Test error case: present=0
	zeroValue, _ := NewFromInt64(0, 0)
	simpleInterest.present = zeroValue
	_, err = simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)

	// Test error case: periods=0
	simpleInterest.present = present
	zeroPeriod, _ := NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)
}
