package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriods(t *testing.T) {
	present, _ := NewFromInt64(5_000, 0)
	interest, _ := NewFromInt64(500, 0)
	rate, _ := NewFromFloat64(0.05)

	simpleInterest := New(Decimal{}, present, interest, rate, Period{})
	expectedPeriods, _ := NewFromFloat64(500 / (5_000 * 0.05))

	period, err := simpleInterest.Periods()
	require.NoError(t, err)

	assert.Equal(t, expectedPeriods.String(), period.String())

	// Test error case: present=0
	zero, _ := NewFromInt64(0, 0)
	simpleInterest.present = zero
	_, err = simpleInterest.Periods()
	assert.Error(t, err)

	// Test error case: rate=0
	simpleInterest.present = present
	simpleInterest.rateInterest, _ = NewFromInt64(0, 0)
	_, err = simpleInterest.Periods()
	assert.Error(t, err)
}

func TestPeriodsWithPresentAndFuture(t *testing.T) {
	future, _ := NewFromInt64(5_500, 0)
	present, _ := NewFromInt64(5_000, 0)
	rate, _ := NewFromFloat64(0.05)

	simpleInterest := New(future, present, Decimal{}, rate, Period{})
	expectedPeriods, _ := NewFromFloat64(2.0)

	period, err := simpleInterest.PeriodsWithPresentAndFuture()
	require.NoError(t, err)

	assert.Equal(t, expectedPeriods.String(), period.String())

	// Test error case: present=0
	zero, _ := NewFromInt64(0, 0)
	simpleInterest.present = zero
	_, err = simpleInterest.PeriodsWithPresentAndFuture()
	assert.Error(t, err)

	// Test error case: rate=0
	simpleInterest.present = present
	simpleInterest.rateInterest, _ = NewFromInt64(0, 0)
	_, err = simpleInterest.PeriodsWithPresentAndFuture()
	assert.Error(t, err)
}
