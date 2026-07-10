package annuities

import (
	"github.com/yeferson59/gofinance/money"
)

func (a Annuity) Present() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthPower, err := rateInterest.Add(money.One).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	quotient, err := growthPower.Sub(money.One).Div(rateInterest.Mul(growthPower))
	if err != nil {
		return money.Money{}, err
	}

	return a.value.ToDecimal().Mul(quotient).ToMoney(a.value.Currency()), nil
}

func (a Annuity) AnticipatePresent() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthPower, err := rateInterest.Add(money.One).Pow(periods.Sub(money.One))
	if err != nil {
		return money.Money{}, err
	}

	quotient, err := growthPower.Sub(money.One).Div(rateInterest.Mul(growthPower))
	if err != nil {
		return money.Money{}, err
	}

	return a.value.ToDecimal().Mul(money.One.Add(quotient)).ToMoney(a.value.Currency()), nil
}
