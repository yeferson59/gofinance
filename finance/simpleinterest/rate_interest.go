package simpleinterest

import (
	"errors"

	"github.com/quagmt/udecimal"
	"github.com/yeferson59/gofinance/money"
)

// RateInterest calculates the interest rate using interest, present value, and periods.
// Formula: Rate = Interest / (Present × Periods)
// Returns an error if present or periods is zero.
func (s SimpleInterest) RateInterest() (money.Money, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return money.Money{}, errors.New("invalid present or periods for operation")
	}

	rateInterest, err := (s.interest.Div(s.present.Mul(numberOfPeriods.Decimal)))
	if err != nil {
		return money.Money{}, err
	}

	return money.Money{Decimal: rateInterest}, nil
}

// RateInterestWithPresentAndFuture calculates the interest rate using future, present value, and periods.
// Formula: Rate = ((Future / Present) - 1) / Periods
// Returns an error if present or periods is zero.
func (s SimpleInterest) RateInterestWithPresentAndFuture() (money.Money, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return money.Money{}, errors.New("invalid present or periods for operation")
	}

	num, err := s.future.Div(s.present.Decimal)
	if err != nil {
		return money.Money{}, err
	}

	rateInterest, err := num.Sub(udecimal.One).Div(numberOfPeriods.Decimal)
	if err != nil {
		return money.Money{}, err
	}

	return money.Money{Decimal: rateInterest}, nil
}
