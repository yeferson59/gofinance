package compositeinterest

import (
	"math"

	"github.com/yeferson59/gofinance/money"
)

func (c CompositeInterest) Periods() (money.Decimal, error) {
	if periodValue, _, err := c.periods.getPeriod(); err == nil && !periodValue.IsZero() {
		return periodValue, nil
	}

	_, periodicRate, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	if c.present.IsZero() || c.future.IsZero() || periodicRate.IsZero() {
		return money.Decimal{}, ErrInvalidOperation
	}

	futureToPresent := c.future.MustDiv(c.present)

	logarithmRatio := math.Log(futureToPresent.InexactFloat64())

	growthFactor := periodicRate.Add(money.One)
	logarithmGrowth := math.Log(growthFactor.InexactFloat64())

	numberOfPeriods := money.MustFromFloat64(logarithmRatio).MustDiv(money.MustFromFloat64(logarithmGrowth))

	return numberOfPeriods, nil
}
