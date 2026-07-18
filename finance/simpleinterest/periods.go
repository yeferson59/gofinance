package simpleinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/decimal"
)

// Periods calculates the number of periods using interest, present value, and rate.
// Formula: Periods = Interest / (Present × Rate)
// Returns an error if present or rate is zero.
func (s SimpleInterest) Periods() (decimal.Decimal, error) {
	if s.present.IsZero() || s.rateInterest.IsZero() {
		return decimal.Decimal{}, errors.New("invalid present or rate interest for operation")
	}

	return s.interest.ToDecimal().Div(s.present.ToDecimal().Mul(s.rateInterest))
}

// PeriodsWithPresentAndFuture calculates the number of periods using future, present value, and rate.
// Formula: Periods = (Future / Present - 1) / Rate
// Returns an error if present or rate is zero.
func (s SimpleInterest) PeriodsWithPresentAndFuture() (decimal.Decimal, error) {
	if s.present.IsZero() || s.rateInterest.IsZero() {
		return decimal.Decimal{}, errors.New("invalid present or rate interest for operation")
	}

	// Step 1: Calculate the ratio of Future to Present
	futureToPresent, err := s.future.ToDecimal().Div(s.present.ToDecimal())
	if err != nil {
		return decimal.Decimal{}, err
	}

	// Step 2: Subtract 1 and divide by the rate to get the number of periods
	return futureToPresent.Sub(decimal.One).Div(s.rateInterest)
}
