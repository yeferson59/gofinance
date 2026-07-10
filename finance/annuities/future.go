package annuities

import (
	"github.com/yeferson59/gofinance/money"
)

func (a Annuity) Future() (money.Money, error) {
	if future, err := a.compositeInterest.Future(); err == nil && !future.IsZero() {
		return future, nil
	}

	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthPower, err := money.One.Add(rateInterest).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	result, err := growthPower.Sub(money.One).Div(rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	return a.value.ToDecimal().Mul(result).ToMoney(a.value.Currency()), nil
}

func (a Annuity) AnticipateFuture() (money.Money, error) {
	if future, err := a.compositeInterest.Future(); err == nil && !future.IsZero() {
		return future, nil
	}

	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	numerator, err := rateInterest.Add(money.One).Pow(money.One.Add(periods))
	if err != nil {
		return money.Money{}, err
	}

	result, err := numerator.Div(rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	return result.Sub(money.One).Mul(a.value.ToDecimal()).ToMoney(a.value.Currency()), nil
}
