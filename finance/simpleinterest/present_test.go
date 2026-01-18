package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPresent(t *testing.T) {
	numPeriods, _ := NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	interest, _ := NewFromInt64(500, 0)
	rate, _ := NewFromFloat64(0.05)

	simpleInterest := New(Decimal{}, Decimal{}, interest, rate, periods)
	expectedPresent, _ := NewFromInt64(5000, 0)

	present, err := simpleInterest.Present()
	require.NoError(t, err)

	assert.Equal(t, expectedPresent.String(), present.String())

	// Test error case: period=0
	zeroPeriod, _ := NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.Present()
	assert.Error(t, err)

	// Test error case: rate=0
	numPeriods2, _ := NewFromInt64(2, 0)
	simpleInterest.periods = NewPeriod(numPeriods2, Days)
	simpleInterest.rateInterest, _ = NewFromInt64(0, 0)
	_, err = simpleInterest.Present()
	assert.Error(t, err)
}

func TestPresentWithFuture(t *testing.T) {
	numPeriods, _ := NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	future, _ := NewFromInt64(5500, 0)
	rate, _ := NewFromFloat64(0.05)

	simpleInterest := New(future, Decimal{}, Decimal{}, rate, periods)
	expectedPresent, _ := NewFromInt64(5000, 0)

	present, err := simpleInterest.PresentWithFuture()
	require.NoError(t, err)

	// Compare the string representation of Decimal values
	assert.Equal(t, expectedPresent.String(), present.String())

	// Test error case: period=0
	zeroPeriod, _ := NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.PresentWithFuture()
	assert.Error(t, err)

	// Test error case: rate=0
	numPeriods2, _ := NewFromInt64(2, 0)
	simpleInterest.periods = NewPeriod(numPeriods2, Days)
	simpleInterest.rateInterest, _ = NewFromInt64(0, 0)
	_, err = simpleInterest.PresentWithFuture()
	assert.Error(t, err)
}
