package simpleinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// Present calculates the present value using interest, rate, and periods.
// Formula: Present = Interest / (Periods × Rate)
// Returns an error if periods or rate is zero.
func (s SimpleInterest) Present() (money.Money, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	if numberOfPeriods.IsZero() || s.rateInterest.IsZero() {
		return money.Money{}, errors.New("invalid period or rate interest for operation")
	}

	return s.interest.DivDecimal(numberOfPeriods.Mul(s.rateInterest))
}

// PresentWithFuture calculates the present value using future value, rate, and periods.
// Formula: Present = Future / (1 + Periods × Rate)
// Returns an error if periods or rate is zero.
func (s SimpleInterest) PresentWithFuture() (money.Money, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	if numberOfPeriods.IsZero() || s.rateInterest.IsZero() {
		return money.Money{}, errors.New("invalid period or rate interest for operation")
	}

	return s.future.DivDecimal(decimal.One.Add(numberOfPeriods.Mul(s.rateInterest)))
}
