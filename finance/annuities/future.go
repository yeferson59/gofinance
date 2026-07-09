package annuities

import (
	"github.com/yeferson59/gofinance/money"
)

func (a Annuity) Future() (money.Money, error) {
	// If the underlying compound interest data can already resolve a future
	// value (explicit or derivable), use it. Otherwise fall back to the
	// annuity formula based on the periodic payment.
	if future, err := a.compositeInterest.Future(); err == nil && !future.IsZero() {
		return future, nil
	}

	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthFactor := rateInterest.Add(money.One)

	growthPower := growthFactor.MustPow(periods)

	quotient, err := growthPower.Sub(money.One).Div(rateInterest)
	if err != nil {
		return money.Money{}, err
	}
	accumulationFactor := quotient

	futureDecimal := a.value.Mul(accumulationFactor.ToMoney(a.value.Currency()))

	return futureDecimal, nil
}
