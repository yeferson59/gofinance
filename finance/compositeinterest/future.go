package compositeinterest

import (
	"math"

	"github.com/yeferson59/gofinance/money"
)

func (c CompositeInterest) Future() (money.Money, error) {
	if !c.future.IsZero() {
		return c.future, nil
	}

	numberOfPeriods, periodicRate, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	if numberOfPeriods.IsZero() || periodicRate.IsZero() || c.present.IsZero() {
		return money.Money{}, ErrInvalidOperation
	}

	growthFactor := periodicRate.Add(money.One)
	compoundGrowth := math.Pow(growthFactor.InexactFloat64(), numberOfPeriods.InexactFloat64())

	future := c.present.Mul(money.MustMoneyFromFloat64(compoundGrowth, c.present.Currency()))

	return future, nil
}
