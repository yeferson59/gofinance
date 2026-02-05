package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestInterest(t *testing.T) {
	numPeriods, _ := money.NewFromInt64(2, 0)
	present, _ := money.New(5_000, 0, money.COP)
	rate, _ := money.NewFromFloat64(0.05)

	periods := NewPeriod(numPeriods, Days)
	simpleInterest := New(money.Money{}, present, money.Money{}, rate, periods)

	expectedInterest, _ := money.New(500, 0, money.COP)

	interest, err := simpleInterest.Interest()
	require.NoError(t, err)

	assert.Equal(t, expectedInterest.String(), interest.String())

	// Test error case
	simpleInterest.periods = Period{}
	_, err = simpleInterest.Interest()
	assert.Error(t, err)
}
