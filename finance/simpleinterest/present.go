package simpleinterest

import "errors"

func (s *InterestSimple) Present() (float64, error) {
	valuePeriods, err := s.periods.GetPeriod()
	if err != nil {
		return 0, err
	}

	period := *valuePeriods

	if period == 0 || s.rateInterest == 0 {
		return 0, errors.New("invalid period or rate interest for operation")
	}

	s.present = ((s.interest) / (period * s.rateInterest))

	return s.present, nil
}

func (s *InterestSimple) PresentWithFuture() (float64, error) {
	valuePeriods, err := s.periods.GetPeriod()
	if err != nil {
		return 0, err
	}

	periods := *valuePeriods

	if (periods == 0 || s.rateInterest == 0) && s.interest == 0 {
		return 0, errors.New("invalid period, rate interest or interest for operation")
	}

	s.present = (s.future / (s.interest + (periods * s.rateInterest)))

	return s.present, nil
}
