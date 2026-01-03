package simpleinterest

import (
	"testing"
)

func NewPeriodTest(t *testing.T) *Period {
	numberPeriod := 2.0
	period := NewPeriod(numberPeriod, Days)

	if period == nil {
		t.Errorf("period mustn't be %v", period)
	}

	if period.weeks != nil {
		t.Errorf("weeks periods must be %v", period)
	}

	if period.months != nil {
		t.Errorf("month periods must be %v", period)
	}

	if period.years != nil {
		t.Errorf("years periods must be %v", period)
	}

	if period.days == nil {
		t.Errorf("days periods mustn't be %v", period)
	}

	periods := *period.days

	if periods != numberPeriod {
		t.Errorf("period must be equals to %f and %f", numberPeriod, periods)
	}

	return period
}

func TestNewPeriod(t *testing.T) {
	_ = NewPeriodTest(t)
}

func TestNewPeriodWithDifferentsTimes(t *testing.T) {
	type DataTest struct {
		number  float64
		Periods Periods
	}

	test := []DataTest{
		{
			number:  2,
			Periods: Days,
		},
		{
			number:  3,
			Periods: Weeks,
		},
		{
			number:  3,
			Periods: Months,
		},
		{
			number:  4,
			Periods: Years,
		},
		{
			number:  4,
			Periods: "",
		},
	}

	t.Run("valid new correct period", func(t *testing.T) {
		for _, value := range test {
			period := NewPeriod(value.number, value.Periods)

			if value.Periods == Days {
				if period.days == nil {
					t.Error(period.days)
				}
			}

			if value.Periods == Weeks {
				if period.weeks == nil {
					t.Error(period.weeks)
				}
			}

			if value.Periods == Months {
				if period.months == nil {
					t.Error(period.months)
				}
			}

			if value.Periods == Years {
				if period.years == nil {
					t.Error(period.years)
				}
			}

			if value.Periods == "" {
				if period.days != nil && period.weeks != nil && period.months != nil && period.years != nil {
					t.Error(period.years)
				}
			}
		}
	})

	t.Run("get valid value for period", func(t *testing.T) {
		for _, value := range test {
			period := NewPeriod(value.number, value.Periods)

			valuePeriod, err := period.getPeriod()

			if value.Periods == "" {
				if valuePeriod != nil {
					t.Error(valuePeriod)
				}
				if err == nil {
					t.Error(err)
				}
			} else {
				periods := *valuePeriod
				if periods != value.number {
					t.Error(valuePeriod)
				}

				if err != nil {
					t.Error(err)
				}
			}
		}
	})
}

func TestGetPeriod(t *testing.T) {
	period := NewPeriodTest(t)
	expectedPeriod := 2.0

	numberPeriod, err := period.getPeriod()
	if err != nil {
		t.Error(err)
	}

	if expectedPeriod != *numberPeriod {
		t.Error("number period and expected period must be equals")
	}
}

func TestNewSimpleInterest(t *testing.T) {
	period := NewPeriodTest(t)

	simpleInterest := New(0, 0, 0, 0, period)

	if simpleInterest.future != 0 {
		t.Error(simpleInterest.future)
	}

	if simpleInterest.present != 0 {
		t.Error(simpleInterest.present)
	}

	if simpleInterest.interest != 0 {
		t.Error(simpleInterest.interest)
	}

	if simpleInterest.rateInterest != 0 {
		t.Error(simpleInterest.rateInterest)
	}

	if simpleInterest.periods == nil {
		t.Error(simpleInterest.periods)
	}
}

func TestGetPeriodWithSimpleInterest(t *testing.T) {
	period := NewPeriodTest(t)

	simpleInterest := New(0, 0, 0, 0, period)

	valuePeriod, err := simpleInterest.GetPeriods()

	if valuePeriod != 2 {
		t.Error(valuePeriod)
	}

	if err != nil {
		t.Error(err)
	}
}
