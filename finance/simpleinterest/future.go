package simpleinterest

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// Future returns the future value as the principal plus the interest already
// known: Future = Present + Interest.
//
// It returns money.ErrCurrencyMismatch if the configured amounts are not all
// in one currency.
func (s SimpleInterest) Future() (money.Money, error) {
	currency, err := s.currency()
	if err != nil {
		return money.Money{}, err
	}

	total, err := s.present.GetDecimal().TryAdd(s.interest.GetDecimal())
	if err != nil {
		return money.Money{}, err
	}

	return money.NewFromDecimal(total, currency), nil
}

// FutureWithRateInterest returns the future value from the principal, rate and
// term: Future = Present × (1 + Periods × Rate).
//
// It returns money.ErrCurrencyMismatch if the configured amounts are not all
// in one currency.
func (s SimpleInterest) FutureWithRateInterest() (money.Money, error) {
	currency, err := s.currency()
	if err != nil {
		return money.Money{}, err
	}

	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	growth := numberOfPeriods.Mul(s.rateInterest).Add(decimal.One)

	future, err := s.present.GetDecimal().TryMul(growth)
	if err != nil {
		return money.Money{}, err
	}

	return money.NewFromDecimal(future, currency), nil
}
