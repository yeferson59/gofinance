package simpleinterest

func (s *SimpleInterest) Future() (float64, error) {
	s.future = s.present + s.interest
	return s.future, nil
}

func (s *SimpleInterest) FutureWithRateInterest() (float64, error) {
	valuePeriods, err := s.periods.getPeriod()
	if err != nil {
		return 0, err
	}
	s.future = s.present * (1 + ((*valuePeriods) * s.rateInterest))

	return s.future, nil
}
