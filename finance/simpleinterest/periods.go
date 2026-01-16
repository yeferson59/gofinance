package simpleinterest

import "errors"

// Periods calculates the periods using interest, present value, and rate.
// Formula: Periods = Interest / (Present × Rate)
// Returns an error if present or rate is zero.
// The timePeriod specifies the unit for the returned Period.
func (s SimpleInterest) Periods() (float64, error) {
	if s.present == 0 || s.rateInterest == 0 {
		return 0, errors.New("invalid present or rate interest for operation")
	}

	numberOfPeriods := (s.interest / (s.present * s.rateInterest))

	return numberOfPeriods, nil
}

// PeriodsWithPresentAndFuture calculates the periods using future, present value, and rate.
// Formula: Periods = (Future / Present - 1) / Rate
// Returns an error if present or rate is zero.
// The timePeriod specifies the unit for the returned Period.
func (s SimpleInterest) PeriodsWithPresentAndFuture() (float64, error) {
	if s.present == 0 || s.rateInterest == 0 {
		return 0, errors.New("invalid present or rate interest for operation")
	}

	numberOfPeriods := (1 / s.rateInterest) * ((s.future / s.present) - 1)

	return numberOfPeriods, nil
}
