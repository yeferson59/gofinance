package simpleinterest

import "errors"

func (s *InterestSimple) Periods(timePeriod Periods) (*Period, error) {
	if s.present == 0 || s.rateInterest == 0 {
		return nil, errors.New("invalid present or rate interest for operation")
	}
	valuePeriod := (s.interest / (s.present * s.rateInterest))
	s.periods = NewPeriod(valuePeriod, timePeriod)

	return s.periods, nil
}

func (s *InterestSimple) PeriodsWithPresentAndFuture(timePeriod Periods) (*Period, error) {
	if s.present == 0 || s.rateInterest == 0 {
		return nil, errors.New("invalid present o rate interest for operation")
	}

	valuePeriod := (1 / s.rateInterest) * ((s.future / s.present) - 1)
	s.periods = NewPeriod(valuePeriod, timePeriod)

	return s.periods, nil
}
