package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterest(t *testing.T) {
	periods := NewPeriod(2, Days)
	simpleInterest := New(0, 5_000, 0, 0.05, periods)
	expectedInterest := 5_000 * 2 * 0.05 // 500

	interest, err := simpleInterest.Interest()
	require.NoError(t, err)

	assert.Equal(t, expectedInterest, interest)

	// Test error case
	simpleInterest.periods = &Period{}
	_, err = simpleInterest.Interest()
	assert.Error(t, err)
}
