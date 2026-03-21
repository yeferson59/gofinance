package simpleinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/money"
)

func (s SimpleInterest) RateInterest() (money.Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Decimal{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return money.Decimal{}, errors.New("invalid present or periods for operation")
	}

	rateInterest, err := s.interest.Div(s.present.Mul(numberOfPeriods.ToMoney()))
	if err != nil {
		return money.Decimal{}, err
	}

	return rateInterest.ToDecimal(), nil
}

func (s SimpleInterest) RateInterestWithPresentAndFuture() (money.Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Decimal{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return money.Decimal{}, errors.New("invalid present or periods for operation")
	}

	num, err := s.future.Div(s.present)
	if err != nil {
		return money.Decimal{}, err
	}

	rateInterest, err := num.Sub(money.MoneyOne).Div(numberOfPeriods.ToMoney())
	if err != nil {
		return money.Decimal{}, err
	}

	return rateInterest.ToDecimal(), nil
}
