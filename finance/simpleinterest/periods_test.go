package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestPeriods(t *testing.T) {
	present, _ := money.New(5_000, 0, money.COP)
	interest, _ := money.New(500, 0, money.COP)
	rate, _ := money.NewFromFloat64(0.05)

	simpleInterest := New(money.Money{}, present, interest, rate, Period{})
	expectedPeriods, _ := money.NewFromFloat64(500 / (5_000 * 0.05))

	period, err := simpleInterest.Periods()
	require.NoError(t, err)

	assert.Equal(t, expectedPeriods.String(), period.String())

	// Test error case: present=0
	zero, _ := money.New(0, 0, money.COP)
	simpleInterest.present = zero
	_, err = simpleInterest.Periods()
	assert.Error(t, err)

	// Test error case: rate=0
	simpleInterest.present = present
	simpleInterest.rateInterest, _ = money.NewFromInt64(0, 0)
	_, err = simpleInterest.Periods()
	assert.Error(t, err)
}

func TestPeriodsWithPresentAndFuture(t *testing.T) {
	future, _ := money.New(5_500, 0, money.COP)
	present, _ := money.New(5_000, 0, money.COP)
	rate, _ := money.NewFromFloat64(0.05)

	simpleInterest := New(future, present, money.Money{}, rate, Period{})
	expectedPeriods, _ := money.NewFromFloat64(2.0)

	period, err := simpleInterest.PeriodsWithPresentAndFuture()
	require.NoError(t, err)

	assert.Equal(t, expectedPeriods.String(), period.String())

	// Test error case: present=0
	zero, _ := money.New(0, 0, money.COP)
	simpleInterest.present = zero
	_, err = simpleInterest.PeriodsWithPresentAndFuture()
	assert.Error(t, err)

	// Test error case: rate=0
	simpleInterest.present = present
	simpleInterest.rateInterest, _ = money.NewFromInt64(0, 0)
	_, err = simpleInterest.PeriodsWithPresentAndFuture()
	assert.Error(t, err)
}
