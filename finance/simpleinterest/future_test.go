package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestSuccessSimpleFuture(t *testing.T) {
	num, _ := money.NewFromInt64(2, 0)
	periods := NewPeriod(num, Days)

	present, _ := money.New(5_000, 0, money.COP)
	interest, _ := money.New(1_000, 0, money.COP)

	simpleInterest := New(money.Money{}, present, interest, money.Decimal{}, periods)
	expectedFuture, _ := money.New(6_000, 0, money.COP)

	future, err := simpleInterest.Future()
	require.NoError(t, err)

	assert.True(t, future.String() != money.Money{}.String(), "future should not be 0")

	assert.Equal(t, expectedFuture.String(), future.String())
}

func TestSuccessSimpleFutureWithRateInterest(t *testing.T) {
	numPeriods, _ := money.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	present, _ := money.New(5_000, 0, money.COP)
	rate, _ := money.NewFromFloat64(0.05)

	simpleInterest := New(money.Money{}, present, money.Money{}, rate, periods)
	expectedFuture, _ := money.New(5_500, 0, money.COP)

	future, err := simpleInterest.FutureWithRateInterest()
	require.NoError(t, err)

	assert.True(t, future.String() != money.Money{}.String(), "future should not be 0")

	assert.Equal(t, expectedFuture.String(), future.String())
}

func TestSuccessComplexFutureWithRateInterest(t *testing.T) {
	type DataTest struct {
		simpleInterest SimpleInterest
	}

	numPeriods, _ := money.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)
	present, _ := money.New(5_000, 0, money.COP)

	testData := []DataTest{
		{
			simpleInterest: NewWithRate(present, 0.02, periods),
		},
		{
			simpleInterest: NewWithRate(present, 0.01, periods),
		},
		{
			simpleInterest: NewWithRate(present, 0.006, periods),
		},
		{
			simpleInterest: NewWithRate(present, 0.0001, periods),
		},
		{
			simpleInterest: NewWithRate(present, 0.05, periods),
		},
	}

	t.Run("not error", func(t *testing.T) {
		for _, data := range testData {
			value, err := data.simpleInterest.FutureWithRateInterest()

			maxValue, _ := money.New(6_000, 0, money.COP)
			assert.True(t, value.Cmp(maxValue) <= 0, "value should be <= 6000")

			require.NoError(t, err)
		}
	})

	t.Run("error on invalid periods", func(t *testing.T) {
		for _, data := range testData {
			data.simpleInterest.periods = Period{}
			_, err := data.simpleInterest.FutureWithRateInterest()

			assert.Error(t, err)
		}
	})
}

// Helper function to create SimpleInterest with rate
func NewWithRate(present money.Money, rate float64, periods Period) SimpleInterest {
	rateDecimal, _ := money.NewFromFloat64(rate)
	return New(money.Money{}, present, money.Money{}, rateDecimal, periods)
}
