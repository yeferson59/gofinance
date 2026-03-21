package simpleinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/money"
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

	present, err := s.interest.Div(numberOfPeriods.Mul(s.rateInterest).ToMoney())
	if err != nil {
		return money.Money{}, err
	}

	return present, nil
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

	present, err := s.future.Div(money.MoneyOne.Add(numberOfPeriods.Mul(s.rateInterest).ToMoney()))
	if err != nil {
		return money.Money{}, err
	}

	return present, nil
}
