package annuities

import (
	"github.com/yeferson59/gofinance/money"
)

func (a Annuity) Present() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthFactor := rateInterest.Add(money.One)

	growthPower := growthFactor.MustPow(periods)

	numerator := money.One.Sub(money.One.MustDiv(growthPower))

	denominator := rateInterest

	quotient, err := numerator.Div(denominator)
	if err != nil {
		return money.Money{}, err
	}
	presentDecimal := a.value.Mul(quotient.ToMoney(a.value.Currency()))

	return presentDecimal, nil
}
