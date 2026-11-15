package simpleinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/internal/guard"

	"github.com/yeferson59/gofinance/v2/decimal"
)

func (s SimpleInterest) RateInterest() (_ decimal.Decimal, err error) {
	defer guard.Recover(&err)

	if s.err != nil {
		return decimal.Decimal{}, s.err
	}

	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return decimal.Decimal{}, errors.New("invalid present or periods for operation")
	}

	return s.interest.GetDecimal().Div(s.present.GetDecimal().Mul(numberOfPeriods))
}

func (s SimpleInterest) RateInterestWithPresentAndFuture() (_ decimal.Decimal, err error) {
	defer guard.Recover(&err)

	if s.err != nil {
		return decimal.Decimal{}, s.err
	}

	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if s.present.IsZero() || numberOfPeriods.IsZero() {
		return decimal.Decimal{}, errors.New("invalid present or periods for operation")
	}

	num, err := s.future.GetDecimal().Div(s.present.GetDecimal())
	if err != nil {
		return decimal.Decimal{}, err
	}

	return num.Sub(decimal.One).Div(numberOfPeriods)
}
