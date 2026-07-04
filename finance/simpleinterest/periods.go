package simpleinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/money"
)

// Periods calculates the number of periods using interest, present value, and rate.
// Formula: Periods = Interest / (Present × Rate)
// Returns an error if present or rate is zero.
func (s SimpleInterest) Periods() (money.Decimal, error) {
	if s.present.IsZero() || s.rateInterest.IsZero() {
		return money.Decimal{}, errors.New("invalid present or rate interest for operation")
	}

	numberOfPeriods, err := s.interest.Div(s.present.Mul(s.rateInterest.ToMoney()))
	if err != nil {
		return money.Decimal{}, err
	}

	return numberOfPeriods.ToDecimal(), nil
}

// PeriodsWithPresentAndFuture calculates the number of periods using future, present value, and rate.
// Formula: Periods = (Future / Present - 1) / Rate
// Returns an error if present or rate is zero.
func (s SimpleInterest) PeriodsWithPresentAndFuture() (money.Decimal, error) {
	if s.present.IsZero() || s.rateInterest.IsZero() {
		return money.Decimal{}, errors.New("invalid present or rate interest for operation")
	}

	// Step 1: Calculate the ratio of Future to Present
	futureToPresent, err := s.future.Div(s.present)
	if err != nil {
		return money.Decimal{}, err
	}

	// Step 2: Subtract 1 from the ratio (Future/Present - 1)
	ratioMinusOne := futureToPresent.Sub(money.One.ToMoney())

	// Step 3: Divide the result by the rate to get the number of periods
	numberOfPeriods, err := ratioMinusOne.Div(s.rateInterest.ToMoney())
	if err != nil {
		return money.Decimal{}, err
	}

	return numberOfPeriods.ToDecimal(), nil
}
