package simpleinterest

func (s *InterestSimple) Interest() (float64, error) {
	valuePeriods, err := s.periods.getPeriod()
	if err != nil {
		return 0, err
	}

	s.interest = s.present * (*valuePeriods) * s.rateInterest

	return s.interest, nil
}
