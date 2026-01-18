package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccessSimpleFuture(t *testing.T) {
	num, _ := NewFromInt64(2, 0)
	periods := NewPeriod(num, Days)

	present, _ := NewFromInt64(5_000, 0)
	interest, _ := NewFromInt64(1_000, 0)

	simpleInterest := New(Decimal{}, present, interest, Decimal{}, periods)
	expectedFuture, _ := NewFromInt64(6_000, 0)

	future, err := simpleInterest.Future()
	require.NoError(t, err)

	assert.True(t, future.String() != Decimal{}.String(), "future should not be 0")

	assert.Equal(t, expectedFuture.String(), future.String())
}

func TestSuccessSimpleFutureWithRateInterest(t *testing.T) {
	numPeriods, _ := NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)

	present, _ := NewFromInt64(5_000, 0)
	rate, _ := NewFromFloat64(0.05)

	simpleInterest := New(Decimal{}, present, Decimal{}, rate, periods)
	expectedFuture, _ := NewFromInt64(5_500, 0)

	future, err := simpleInterest.FutureWithRateInterest()
	require.NoError(t, err)

	assert.True(t, future.String() != Decimal{}.String(), "future should not be 0")

	assert.Equal(t, expectedFuture.String(), future.String())
}

func TestSuccessComplexFutureWithRateInterest(t *testing.T) {
	type DataTest struct {
		simpleInterest SimpleInterest
	}

	numPeriods, _ := NewFromInt64(2, 0)
	periods := NewPeriod(numPeriods, Days)
	present, _ := NewFromInt64(5_000, 0)

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

			maxValue, _ := NewFromInt64(6_000, 0)
			assert.True(t, value.Decimal.Cmp(maxValue.Decimal) <= 0, "value should be <= 6000")

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
func NewWithRate(present Decimal, rate float64, periods Period) SimpleInterest {
	rateDecimal, _ := NewFromFloat64(rate)
	return New(Decimal{}, present, Decimal{}, rateDecimal, periods)
}
