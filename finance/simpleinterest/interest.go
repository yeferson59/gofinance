package simpleinterest

// Interest calculates the interest using present value, rate, and periods.
// Formula: Interest = Present × Periods × Rate
// Returns an error if periods is invalid.
func (s SimpleInterest) Interest() (float64, error) {
	valuePeriods, err := s.periods.getPeriod()
	if err != nil {
		return 0, err
	}

	interest := s.present * valuePeriods * s.rateInterest

	return interest, nil
}
