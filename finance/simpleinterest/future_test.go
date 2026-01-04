package simpleinterest

import "testing"

func TestSuccessSimpleFuture(t *testing.T) {
	periods := NewPeriod(2, Days)
	simpleInterest := New(0, 5_000, 1_000, 0, periods)
	expectedFuture := 6_000.0

	future, err := simpleInterest.Future()
	if err != nil {
		t.Error(err)
	}

	if future == 0 {
		t.Error(future)
	}

	if future != expectedFuture {
		t.Error(future)
	}
}

func TestSuccessSimpleFutureWithRateInterest(t *testing.T) {
	periods := NewPeriod(2, Days)
	simpleInterest := New(0, 5_000, 0, 0.05, periods)
	expectedFuture := 5_500.0

	future, err := simpleInterest.FutureWithRateInterest()
	if err != nil {
		t.Error(err)
	}

	if future == 0 {
		t.Error(future)
	}

	if future != expectedFuture {
		t.Error(future)
	}
}

func TestSuccessComplexFutureWithRateInterest(t *testing.T) {
	type DataTest struct {
		simpleInterest *SimpleInterest
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

			if value > 6_000.0 {
				t.Error(value)
			}

			if err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("not error", func(t *testing.T) {
		for _, data := range testData {
			data.simpleInterest.periods = &Period{}
			_, err := data.simpleInterest.FutureWithRateInterest()

			if err == nil {
				t.Error(err)
			}
		}
	})
}
