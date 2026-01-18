package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterest(t *testing.T) {
	numPeriods, _ := NewFromInt64(2, 0)
	present, _ := NewFromInt64(5_000, 0)
	rate, _ := NewFromFloat64(0.05)

	periods := NewPeriod(numPeriods, Days)
	simpleInterest := New(Decimal{}, present, Decimal{}, rate, periods)

	expectedInterest, _ := NewFromInt64(500, 0)

	interest, err := simpleInterest.Interest()
	require.NoError(t, err)

	assert.Equal(t, expectedInterest.String(), interest.String())

	// Test error case
	simpleInterest.periods = Period{}
	_, err = simpleInterest.Interest()
	assert.Error(t, err)
}
