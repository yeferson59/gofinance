package simpleinterest

import "errors"

// Present calculates the present value using interest, rate, and periods.
// Formula: Present = Interest / (Periods × Rate)
// Returns an error if periods or rate is zero.
func (s *SimpleInterest) Present() (float64, error) {
	valuePeriods, err := s.periods.getPeriod()
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

// PresentWithFuture calculates the present value using future value, rate, and periods.
// Formula: Present = Future / (1 + Periods × Rate)
// Returns an error if periods or rate is zero.
func (s *SimpleInterest) PresentWithFuture() (float64, error) {
	valuePeriods, err := s.periods.getPeriod()
	if err != nil {
		return 0, err
	}

	periods := *valuePeriods

	if periods == 0 || s.rateInterest == 0 {
		return 0, errors.New("invalid period or rate interest for operation")
	}

	s.present = (s.future / (1 + (periods * s.rateInterest)))

	return s.present, nil
}
