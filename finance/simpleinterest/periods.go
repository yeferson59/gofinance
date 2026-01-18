package simpleinterest

import (
	"errors"

	"github.com/quagmt/udecimal"
)

// Periods calculates the periods using interest, present value, and rate.
// Formula: Periods = Interest / (Present × Rate)
// Returns an error if present or rate is zero.
// The timePeriod specifies the unit for the returned Period.
func (s SimpleInterest) Periods() (Decimal, error) {
	if s.present.IsZero() || s.rateInterest.IsZero() {
		return Decimal{}, errors.New("invalid present or rate interest for operation")
	}

	numberOfPeriods, err := s.interest.Div(s.present.Mul(s.rateInterest.Decimal))
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{numberOfPeriods}, nil
}

// PeriodsWithPresentAndFuture calculates the periods using future, present value, and rate.
// Formula: Periods = (Future / Present - 1) / Rate
// Returns an error if present or rate is zero.
// The timePeriod specifies the unit for the returned Period.
func (s SimpleInterest) PeriodsWithPresentAndFuture() (Decimal, error) {
	if s.present.IsZero() || s.rateInterest.IsZero() {
		return Decimal{}, errors.New("invalid present or rate interest for operation")
	}

	// Step 1: Calculate the ratio of Future to Present
	futureToPresent, err := s.future.Div(s.present.Decimal)
	if err != nil {
		return Decimal{}, err
	}

	// Step 2: Subtract 1 from the ratio (Future/Present - 1)
	ratioMinusOne := futureToPresent.Sub(udecimal.One)

	// Step 3: Divide the result by the rate to get the number of periods
	numberOfPeriods, err := ratioMinusOne.Div(s.rateInterest.Decimal)
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{numberOfPeriods}, nil
}
