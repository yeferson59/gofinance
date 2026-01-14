package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccessSimpleFuture(t *testing.T) {
	periods := NewPeriod(2, Days)
	simpleInterest := New(0, 5_000, 1_000, 0, periods)
	expectedFuture := 6_000.0

	future, err := simpleInterest.Future()
	require.NoError(t, err)

	assert.True(t, future != 0, "future should not be 0")

	assert.Equal(t, expectedFuture, future)
}

func TestSuccessSimpleFutureWithRateInterest(t *testing.T) {
	periods := NewPeriod(2, Days)
	simpleInterest := New(0, 5_000, 0, 0.05, periods)
	expectedFuture := 5_500.0

	future, err := simpleInterest.FutureWithRateInterest()
	require.NoError(t, err)

	assert.True(t, future != 0, "future should not be 0")

	assert.Equal(t, expectedFuture, future)
}

func TestSuccessComplexFutureWithRateInterest(t *testing.T) {
	type DataTest struct {
		simpleInterest SimpleInterest
	}

	periods := NewPeriod(2, Days)

	testData := []DataTest{
		{
			simpleInterest: New(0, 5_000, 0, 0.02, periods),
		},
		{
			simpleInterest: New(0, 5_000, 0, 0.01, periods),
		},
		{
			simpleInterest: New(0, 5_000, 0, 0.006, periods),
		},
		{
			simpleInterest: New(0, 5_000, 0, 0.0001, periods),
		},
		{
			simpleInterest: New(0, 5_000, 0, 0.05, periods),
		},
	}

	t.Run("not error", func(t *testing.T) {
		for _, data := range testData {
			value, err := data.simpleInterest.FutureWithRateInterest()

			assert.True(t, value <= 6_000.0, "value should be <= 6000")

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
