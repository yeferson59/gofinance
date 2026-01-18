package simpleinterest

import (
	"errors"

	"github.com/quagmt/udecimal"
)

// Present calculates the present value using interest, rate, and periods.
// Formula: Present = Interest / (Periods × Rate)
// Returns an error if periods or rate is zero.
func (s SimpleInterest) Present() (Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return Decimal{}, err
	}

	if numberOfPeriods.IsZero() || s.rateInterest.IsZero() {
		return Decimal{}, errors.New("invalid period or rate interest for operation")
	}

	present, err := s.interest.Div(numberOfPeriods.Mul(s.rateInterest.Decimal))
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{present}, nil
}

// PresentWithFuture calculates the present value using future value, rate, and periods.
// Formula: Present = Future / (1 + Periods × Rate)
// Returns an error if periods or rate is zero.
func (s SimpleInterest) PresentWithFuture() (Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return Decimal{}, err
	}

	if numberOfPeriods.IsZero() || s.rateInterest.IsZero() {
		return Decimal{}, errors.New("invalid period or rate interest for operation")
	}

	present, err := s.future.Div(udecimal.One.Add(numberOfPeriods.Mul(s.rateInterest.Decimal)))
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{present}, nil
}
