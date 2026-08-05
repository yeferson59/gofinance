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

	// With no interest the payments are neither discounted nor grown, so the
	// present value is just their sum. The general formula divides by the
	// rate, so the limit is returned directly.
	if rateInterest.IsZero() {
		return money.FromDecimal(a.value.ToDecimal().Mul(periods), a.currency), nil
	}

	growthPower, err := rateInterest.Add(decimal.One).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	quotient, err := growthPower.Sub(decimal.One).Div(rateInterest.Mul(growthPower))
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(a.value.ToDecimal().Mul(quotient), a.currency), nil
}

func (a Annuity) AnticipatePresent() (money.Money, error) {
	periods, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	// At a zero rate the payment's timing inside the period stops mattering,
	// so this matches Present: the sum of the payments.
	if rateInterest.IsZero() {
		return money.FromDecimal(a.value.ToDecimal().Mul(periods), a.currency), nil
	}

	growthPower, err := rateInterest.Add(decimal.One).Pow(periods.Sub(decimal.One))
	if err != nil {
		return money.Money{}, err
	}

	quotient, err := growthPower.Sub(decimal.One).Div(rateInterest.Mul(growthPower))
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(a.value.ToDecimal().Mul(decimal.One.Add(quotient)), a.currency), nil
}
