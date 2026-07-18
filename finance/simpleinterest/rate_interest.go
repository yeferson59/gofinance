package simpleinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/decimal"
)

func (s SimpleInterest) RateInterest() (decimal.Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return decimal.Decimal{}, errors.New("invalid present or periods for operation")
	}

	return s.interest.ToDecimal().Div(s.present.ToDecimal().Mul(numberOfPeriods))
}

func (s SimpleInterest) RateInterestWithPresentAndFuture() (decimal.Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return decimal.Decimal{}, errors.New("invalid present or periods for operation")
	}

	num, err := s.future.ToDecimal().Div(s.present.ToDecimal())
	if err != nil {
		return decimal.Decimal{}, err
	}

	return num.Sub(decimal.One).Div(numberOfPeriods)
}
