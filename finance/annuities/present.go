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

	growthPower := money.MustFromFloat64(math.Pow(rateInterest.Add(money.One).InexactFloat64(), periods.InexactFloat64()))

	quotient, err := growthPower.Sub(money.One).Div(rateInterest.Mul(growthPower))
	if err != nil {
		return money.Money{}, err
	}

	present := a.value.Mul(quotient.ToMoney(a.value.Currency()))

	return present, nil
}

func (a Annuity) AnticipatePresent() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthPower := money.MustFromFloat64(math.Pow(rateInterest.Add(money.One).InexactFloat64(), periods.Sub(money.One).InexactFloat64()))

	quotient, err := growthPower.Sub(money.One).Div(rateInterest.Mul(growthPower))
	if err != nil {
		return money.Money{}, err
	}

	present := a.value.ToDecimal().Mul(money.One.Add(quotient))

	return present.ToMoney(a.value.Currency()), nil
}
