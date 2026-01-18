package simpleinterest

import (
	"errors"

	"github.com/quagmt/udecimal"
)

// RateInterest calculates the interest rate using interest, present value, and periods.
// Formula: Rate = Interest / (Present × Periods)
// Returns an error if present or periods is zero.
func (s SimpleInterest) RateInterest() (Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return Decimal{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return Decimal{}, errors.New("invalid present or periods for operation")
	}

	rateInterest, err := (s.interest.Div(s.present.Mul(numberOfPeriods.Decimal)))
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{rateInterest}, nil
}

// RateInterestWithPresentAndFuture calculates the interest rate using future, present value, and periods.
// Formula: Rate = ((Future / Present) - 1) / Periods
// Returns an error if present or periods is zero.
func (s SimpleInterest) RateInterestWithPresentAndFuture() (Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return Decimal{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return Decimal{}, errors.New("invalid present or periods for operation")
	}

	num, err := s.future.Div(s.present.Decimal)
	if err != nil {
		return Decimal{}, err
	}

	rateInterest, err := num.Sub(udecimal.One).Div(numberOfPeriods.Decimal)
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{rateInterest}, nil
}
