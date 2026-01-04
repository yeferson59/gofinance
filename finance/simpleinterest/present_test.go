package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPresent(t *testing.T) {
	periods := NewPeriod(2, Days)
	simpleInterest := New(0, 0, 500, 0.05, periods) // interest=500, rate=0.05, periods=2
	expectedPresent := 500 / (2 * 0.05)             // 5000

	present, err := simpleInterest.Present()
	require.NoError(t, err)

	assert.Equal(t, expectedPresent, present)

	// Test error case: period=0
	simpleInterest.periods = NewPeriod(0, Days)
	_, err = simpleInterest.Present()
	assert.Error(t, err)

	// Test error case: rate=0
	simpleInterest.periods = NewPeriod(2, Days)
	simpleInterest.rateInterest = 0
	_, err = simpleInterest.Present()
	assert.Error(t, err)
}

func TestPresentWithFuture(t *testing.T) {
	periods := NewPeriod(2, Days)
	simpleInterest := New(5500, 0, 0, 0.05, periods) // future=5500, rate=0.05, periods=2
	expectedPresent := 5500 / (1 + 2*0.05)           // 5000

	present, err := simpleInterest.PresentWithFuture()
	require.NoError(t, err)

	assert.InDelta(t, expectedPresent, present, 1e-10)

	// Test error case: period=0
	simpleInterest.periods = NewPeriod(0, Days)
	_, err = simpleInterest.PresentWithFuture()
	assert.Error(t, err)

	// Test error case: rate=0
	simpleInterest.periods = NewPeriod(2, Days)
	simpleInterest.rateInterest = 0
	_, err = simpleInterest.PresentWithFuture()
	assert.Error(t, err)
}
