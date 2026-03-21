package compositeinterest

import (
	"math"

	"github.com/yeferson59/gofinance/money"
)

func (c CompositeInterest) Present() (money.Money, error) {
	if !c.present.IsZero() {
		return c.present, nil
	}

	numberOfPeriods, periodicRate, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	if c.future.IsZero() || periodicRate.IsZero() || numberOfPeriods.IsZero() {
		return money.Money{}, ErrInvalidOperation
	}

	growthFactor := periodicRate.Add(money.One)
	discountFactor := math.Pow(growthFactor.InexactFloat64(), numberOfPeriods.InexactFloat64())

	present := c.future.MustDiv(money.MustMoneyFromFloat64(discountFactor, c.future.Currency()))

	return present, nil
}
