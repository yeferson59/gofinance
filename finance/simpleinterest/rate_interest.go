package simpleinterest

import "errors"

// RateInterest calculates the interest rate using interest, present value, and periods.
// Formula: Rate = Interest / (Present × Periods)
// Returns an error if present or periods is zero.
func (s SimpleInterest) RateInterest() (float64, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return 0, err
	}

	if s.present == 0 || numberOfPeriods == 0 {
		return 0, errors.New("invalid present or periods for operation")
	}

	rateInterest := (s.interest / (s.present * numberOfPeriods))

	return rateInterest, nil
}

// RateInterestWithPresentAndFuture calculates the interest rate using future, present value, and periods.
// Formula: Rate = ((Future / Present) - 1) / Periods
// Returns an error if present or periods is zero.
func (s SimpleInterest) RateInterestWithPresentAndFuture() (float64, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return 0, err
	}

	if s.present == 0 || numberOfPeriods == 0 {
		return 0, errors.New("invalid present or periods for operation")
	}

	rateInterest := (((s.future / s.present) - 1) / numberOfPeriods)

	return rateInterest, nil
}
