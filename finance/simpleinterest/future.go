package simpleinterest

// Future calculates the future value using present value and interest.
// Formula: Future = Present + Interest
// Returns an error if calculation fails, but currently always succeeds.
func (s SimpleInterest) Future() (float64, error) {
	future := s.present + s.interest

	return future, nil
}

// FutureWithRateInterest calculates the future value using present value, rate, and periods.
// Formula: Future = Present × (1 + Periods × Rate)
// Returns an error if periods is invalid.
func (s SimpleInterest) FutureWithRateInterest() (float64, error) {
	valuePeriods, err := s.periods.getPeriod()
	if err != nil {
		return 0, err
	}

	future := s.present * (1 + (valuePeriods * s.rateInterest))

	return future, nil
}
