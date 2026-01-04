package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateInterest(t *testing.T) {
	periods := NewPeriod(2, Days)
	simpleInterest := New(0, 5_000, 500, 0, periods) // interest=500, present=5000, periods=2
	expectedRate := 500.0 / (5000.0 * 2.0)           // 0.05

	rate, err := simpleInterest.RateInterest()
	require.NoError(t, err)

	assert.Equal(t, expectedRate, rate)

	// Test error case: present=0
	simpleInterest.present = 0
	_, err = simpleInterest.RateInterest()
	assert.Error(t, err)

	// Test error case: periods=0
	simpleInterest.present = 5_000
	simpleInterest.periods = NewPeriod(0, Days)
	_, err = simpleInterest.RateInterest()
	assert.Error(t, err)
}

func TestRateInterestWithPresentAndFuture(t *testing.T) {
	periods := NewPeriod(2, Days)
	simpleInterest := New(5_500, 5_000, 0, 0, periods)      // future=5500, present=5000, periods=2
	expectedRate := (float64(5_500)/float64(5_000) - 1) / 2 // 0.05

	rate, err := simpleInterest.RateInterestWithPresentAndFuture()
	require.NoError(t, err)

	assert.InDelta(t, expectedRate, rate, 1e-10)

	// Test error case: present=0
	simpleInterest.present = 0
	_, err = simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)

	// Test error case: periods=0
	simpleInterest.present = 5_000
	simpleInterest.periods = NewPeriod(0, Days)
	_, err = simpleInterest.RateInterestWithPresentAndFuture()
	assert.Error(t, err)
}
