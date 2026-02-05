package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestPresent(t *testing.T) {
	numPeriods, _ := money.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	interest, _ := money.New(500, 0, money.COP)
	rate, _ := money.NewFromFloat64(0.05)

	simpleInterest := New(money.Money{}, money.Money{}, interest, rate, periods)
	expectedPresent, _ := money.New(5000, 0, money.COP)

	present, err := simpleInterest.Present()
	require.NoError(t, err)

	assert.Equal(t, expectedPresent.String(), present.String())

	// Test error case: period=0
	zeroPeriod, _ := money.NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.Present()
	assert.Error(t, err)

	// Test error case: rate=0
	numPeriods2, _ := money.NewFromInt64(2, 0)
	simpleInterest.periods = NewPeriod(numPeriods2, Days)
	simpleInterest.rateInterest, _ = money.NewFromInt64(0, 0)
	_, err = simpleInterest.Present()
	assert.Error(t, err)
}

func TestPresentWithFuture(t *testing.T) {
	numPeriods, _ := money.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	future, _ := money.New(5500, 0, money.COP)
	rate, _ := money.NewFromFloat64(0.05)

	simpleInterest := New(future, money.Money{}, money.Money{}, rate, periods)
	expectedPresent, _ := money.New(5000, 0, money.COP)

	present, err := simpleInterest.PresentWithFuture()
	require.NoError(t, err)

	// Compare the string representation of Decimal values
	assert.Equal(t, expectedPresent.String(), present.String())

	// Test error case: period=0
	zeroPeriod, _ := money.NewFromInt64(0, 0)
	simpleInterest.periods = NewPeriod(zeroPeriod, Days)
	_, err = simpleInterest.PresentWithFuture()
	assert.Error(t, err)

	// Test error case: rate=0
	numPeriods2, _ := money.NewFromInt64(2, 0)
	simpleInterest.periods = NewPeriod(numPeriods2, Days)
	simpleInterest.rateInterest, _ = money.NewFromInt64(0, 0)
	_, err = simpleInterest.PresentWithFuture()
	assert.Error(t, err)
}
