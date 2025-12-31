package simpleinterest

import "errors"

func (s *InterestSimple) RateInterest() (float64, error) {
	valuePeriods, err := s.periods.GetPeriod()
	if err != nil {
		return 0, err
	}

	periods := *valuePeriods

	if s.present == 0 || periods == 0 {
		return 0, errors.New("invalid present or periods for operation")
	}

	s.rateInterest = (s.interest / (s.present * periods))

	return s.rateInterest, nil
}

func (s *InterestSimple) RateInterestWithPresentAndFuture() (float64, error) {
	valuePeriods, err := s.periods.GetPeriod()
	if err != nil {
		return 0, err
	}

	periods := *valuePeriods

	if s.present == 0 || periods == 0 {
		return 0, errors.New("invalid present or periods for operation")
	}

	s.rateInterest = (((s.future / s.present) - 1) / periods)

	return s.rateInterest, nil
}
