package annuities

import (
	"math"

	"github.com/yeferson59/gofinance/money"
)

func (a Annuity) Present() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthFactor := rateInterest.Add(money.One)

	growthPower := money.MustFromFloat64(math.Pow(growthFactor.InexactFloat64(), periods.InexactFloat64()))

	numerator := money.One.Sub(money.One.MustDiv(growthPower))

	denominator := rateInterest

	quotient, err := numerator.Div(denominator)
	if err != nil {
		return money.Money{}, err
	}
	presentDecimal := a.value.Mul(quotient.ToMoney(a.value.Currency()))

	return presentDecimal, nil
}
