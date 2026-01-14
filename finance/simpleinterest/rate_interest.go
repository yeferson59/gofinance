package simpleinterest

import "errors"

// RateInterest calculates the interest rate using interest, present value, and periods.
// Formula: Rate = Interest / (Present × Periods)
// Returns an error if present or periods is zero.
func (s *SimpleInterest) RateInterest() (float64, error) {
	valuePeriods, err := s.periods.getPeriod()
	if err != nil {
		return 0, err
	}

	if s.present == 0 || valuePeriods == 0 {
		return 0, errors.New("invalid present or periods for operation")
	}

	rateInterest := (s.interest / (s.present * valuePeriods))

	return rateInterest, nil
}

// RateInterestWithPresentAndFuture calculates the interest rate using future, present value, and periods.
// Formula: Rate = ((Future / Present) - 1) / Periods
// Returns an error if present or periods is zero.
func (s *SimpleInterest) RateInterestWithPresentAndFuture() (float64, error) {
	valuePeriods, err := s.periods.getPeriod()
	if err != nil {
		return 0, err
	}

	periods := valuePeriods

	if s.present == 0 || periods == 0 {
		return 0, errors.New("invalid present or periods for operation")
	}

	rateInterest := (((s.future / s.present) - 1) / periods)

	return rateInterest, nil
}
