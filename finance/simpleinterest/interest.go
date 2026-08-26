package simpleinterest

import (
	"github.com/yeferson59/gofinance/v2/money"
)

// Interest calculates the interest using present value, rate, and periods.
// Formula: Interest = Present × Periods × Rate
//
// It returns an error if the periods are invalid, and
// money.ErrCurrencyMismatch if the configured amounts are not all in one
// currency.
func (s SimpleInterest) Interest() (money.Money, error) {
	currency, err := s.currency()
	if err != nil {
		return money.Money{}, err
	}

	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	interest, err := s.present.GetDecimal().TryMul(numberOfPeriods.Mul(s.rateInterest))
	if err != nil {
		return money.Money{}, err
	}

	return money.NewFromDecimal(interest, currency), nil
}

// InterestWithPresentAndFuture calculates the interest as the difference
// between the two known amounts: Interest = Future − Present.
//
// It returns money.ErrCurrencyMismatch if the configured amounts are not all
// in one currency.
func (s SimpleInterest) InterestWithPresentAndFuture() (money.Money, error) {
	currency, err := s.currency()
	if err != nil {
		return money.Money{}, err
	}

	interest, err := s.future.GetDecimal().TrySub(s.present.GetDecimal())
	if err != nil {
		return money.Money{}, err
	}

	return money.NewFromDecimal(interest, currency), nil
}
