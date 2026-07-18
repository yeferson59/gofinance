package annuities

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func (a Annuity) Present() (money.Money, error) {
	periods, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthPower, err := rateInterest.Add(decimal.One).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	quotient, err := growthPower.Sub(decimal.One).Div(rateInterest.Mul(growthPower))
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(a.value.ToDecimal().Mul(quotient), a.value.Currency()), nil
}

func (a Annuity) AnticipatePresent() (money.Money, error) {
	periods, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthPower, err := rateInterest.Add(decimal.One).Pow(periods.Sub(decimal.One))
	if err != nil {
		return money.Money{}, err
	}

	quotient, err := growthPower.Sub(decimal.One).Div(rateInterest.Mul(growthPower))
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(a.value.ToDecimal().Mul(decimal.One.Add(quotient)), a.value.Currency()), nil
}
