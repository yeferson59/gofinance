package simpleinterest

import "errors"

// Periods calculates the periods using interest, present value, and rate.
// Formula: Periods = Interest / (Present × Rate)
// Returns an error if present or rate is zero.
// The timePeriod specifies the unit for the returned Period.
func (s SimpleInterest) Periods(timePeriod Periods) (Period, error) {
	if s.present == 0 || s.rateInterest == 0 {
		return Period{}, errors.New("invalid present or rate interest for operation")
	}
	valuePeriod := (s.interest / (s.present * s.rateInterest))
	periods := NewPeriod(valuePeriod, timePeriod)

	return periods, nil
}

// PeriodsWithPresentAndFuture calculates the periods using future, present value, and rate.
// Formula: Periods = (Future / Present - 1) / Rate
// Returns an error if present or rate is zero.
// The timePeriod specifies the unit for the returned Period.
func (s SimpleInterest) PeriodsWithPresentAndFuture(timePeriod Periods) (Period, error) {
	if s.present == 0 || s.rateInterest == 0 {
		return Period{}, errors.New("invalid present or rate interest for operation")
	}

	valuePeriod := (1 / s.rateInterest) * ((s.future / s.present) - 1)
	periods := NewPeriod(valuePeriod, timePeriod)

	return periods, nil
}
