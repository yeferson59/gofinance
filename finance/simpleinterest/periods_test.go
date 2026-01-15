package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriods(t *testing.T) {
	simpleInterest := New(0, 5_000, 500, 0.05, Period{}) // interest=500, present=5000, rate=0.05
	expectedPeriods := 500 / (5_000 * 0.05)              // 2

	period, err := simpleInterest.Periods()
	require.NoError(t, err)

	assert.Equal(t, expectedPeriods, period)

	// Test error case: present=0
	simpleInterest.present = 0
	_, err = simpleInterest.Periods()
	assert.Error(t, err)

	// Test error case: rate=0
	simpleInterest.present = 5_000
	simpleInterest.rateInterest = 0
	_, err = simpleInterest.Periods()
	assert.Error(t, err)
}

func TestPeriodsWithPresentAndFuture(t *testing.T) {
	simpleInterest := New(5_500, 5_000, 0, 0.05, Period{}) // future=5500, present=5000, rate=0.05
	expectedPeriods := 2.0

	period, err := simpleInterest.PeriodsWithPresentAndFuture()
	require.NoError(t, err)

	assert.InDelta(t, expectedPeriods, period, 1e-10)

	// Test error case: present=0
	simpleInterest.present = 0
	_, err = simpleInterest.PeriodsWithPresentAndFuture()
	assert.Error(t, err)

	// Test error case: rate=0
	simpleInterest.present = 5_000
	simpleInterest.rateInterest = 0
	_, err = simpleInterest.PeriodsWithPresentAndFuture()
	assert.Error(t, err)
}
