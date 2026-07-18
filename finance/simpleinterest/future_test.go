package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/money"
)

func TestSuccessSimpleFuture(t *testing.T) {
	num, _ := decimal.NewFromInt64(2, 0)
	periods := NewPeriod(num, Days)

	present, _ := money.New(5_000, 0, money.COP)
	interest, _ := money.New(1_000, 0, money.COP)

	simpleInterest := New(money.Money{}, present, interest, decimal.Decimal{}, periods)
	expectedFuture, _ := money.New(6_000, 0, money.COP)

	future, err := simpleInterest.Future()
	require.NoError(t, err)

	assert.True(t, future.String() != money.Money{}.String(), "future should not be 0")

	assert.Equal(t, expectedFuture.String(), future.String())
}

func TestSuccessSimpleFutureWithRateInterest(t *testing.T) {
	numPeriods, _ := decimal.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	present, _ := money.New(5_000, 0, money.COP)
	rate, _ := decimal.NewFromFloat64(0.05)

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
		expectedFuture int64
	}

	numPeriods, _ := decimal.NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)
	present, _ := money.New(5_000, 0, money.COP)

	// FV = PV × (1 + r×n) with PV = 5000 and n = 2
	testData := []DataTest{
		{
			simpleInterest: NewWithRate(present, 0.02, periods),
			expectedFuture: 5_200,
		},
		{
			simpleInterest: NewWithRate(present, 0.01, periods),
			expectedFuture: 5_100,
		},
		{
			simpleInterest: NewWithRate(present, 0.006, periods),
			expectedFuture: 5_060,
		},
		{
			simpleInterest: NewWithRate(present, 0.0001, periods),
			expectedFuture: 5_001,
		},
		{
			simpleInterest: NewWithRate(present, 0.05, periods),
			expectedFuture: 5_500,
		},
	}

	t.Run("not error", func(t *testing.T) {
		for _, data := range testData {
			value, err := data.simpleInterest.FutureWithRateInterest()
			require.NoError(t, err)

			expected, _ := money.New(data.expectedFuture, 0, money.COP)
			assert.Equal(t, expected.String(), value.String())
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
	rateDecimal, _ := decimal.NewFromFloat64(rate)
	return New(money.Money{}, present, money.Money{}, rateDecimal, periods)
}
