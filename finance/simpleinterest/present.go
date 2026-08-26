package simpleinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// Present calculates the present value using interest, rate, and periods.
// Formula: Present = Interest / (Periods × Rate)
//
// It returns an error if periods or rate is zero, and
// money.ErrCurrencyMismatch if the configured amounts are not all in one
// currency.
func (s SimpleInterest) Present() (money.Money, error) {
	currency, err := s.currency()
	if err != nil {
		return money.Money{}, err
	}

	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	if numberOfPeriods.IsZero() || s.rateInterest.IsZero() {
		return money.Money{}, errors.New("invalid period or rate interest for operation")
	}

	present, err := s.interest.GetDecimal().Div(numberOfPeriods.Mul(s.rateInterest))
	if err != nil {
		return money.Money{}, err
	}

	return money.NewFromDecimal(present, currency), nil
}

// PresentWithFuture calculates the present value using future value, rate, and periods.
// Formula: Present = Future / (1 + Periods × Rate)
//
// It returns an error if periods or rate is zero, and
// money.ErrCurrencyMismatch if the configured amounts are not all in one
// currency.
func (s SimpleInterest) PresentWithFuture() (money.Money, error) {
	currency, err := s.currency()
	if err != nil {
		return money.Money{}, err
	}

	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	if numberOfPeriods.IsZero() || s.rateInterest.IsZero() {
		return money.Money{}, errors.New("invalid period or rate interest for operation")
	}

	present, err := s.future.GetDecimal().Div(decimal.One.Add(numberOfPeriods.Mul(s.rateInterest)))
	if err != nil {
		return money.Money{}, err
	}

	return money.NewFromDecimal(present, currency), nil
}
