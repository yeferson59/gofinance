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
