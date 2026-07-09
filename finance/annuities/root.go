// Package annuities provides functionality for annuity calculations.
package annuities

import (
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

type Annuity struct {
	value             money.Money
	compositeInterest compositeinterest.CompositeInterest
}

func New(value, present, future money.Money, period compositeinterest.Period, rateInterest compositeinterest.RateInterest) (Annuity, error) {
	ci, err := compositeinterest.New(present, future, rateInterest, period)
	if err != nil {
		return Annuity{}, err
	}

	return Annuity{
		value:             value,
		compositeInterest: ci,
	}, nil
}

func (a Annuity) PaymentFromPresentValue() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return money.Money{}, err
	}

	growthFactor := rateInterest.Add(money.One)

	growthPower := growthFactor.MustPow(periods)

	numerator := rateInterest.Mul(growthPower)

	denominator := growthPower.Sub(money.One)

	annuity := present.Mul(numerator.MustDiv(denominator).ToMoney(present.Currency()))

	return annuity, nil
}

func (a Annuity) PaymentFromFutureValue() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return money.Money{}, err
	}

	growthFactor := rateInterest.Add(money.One)

	growthPower := growthFactor.MustPow(periods)

	denominator := growthPower.Sub(money.One)

	annuity := future.Mul(rateInterest.MustDiv(denominator).ToMoney(future.Currency()))

	return annuity, nil
}
