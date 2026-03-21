package compositeinterest

import (
	"math"

	"github.com/yeferson59/gofinance/money"
)

func (c CompositeInterest) Interest() (money.Decimal, error) {
	if !c.rateInterest.value.IsZero() {
		return c.rateInterest.value, nil
	}

	numberOfPeriods, _, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	if c.future.IsZero() || c.present.IsZero() || numberOfPeriods.IsZero() {
		return money.Decimal{}, ErrInvalidOperation
	}

	futureToPresent := c.future.MustDiv(c.present)

	reciprocalPeriods := money.One.MustDiv(numberOfPeriods)

	growthFactor := math.Pow(futureToPresent.InexactFloat64(), reciprocalPeriods.InexactFloat64())

	periodicRate := growthFactor - 1

	return money.MustFromFloat64(periodicRate), nil
}
