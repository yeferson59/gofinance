package compositeinterest

import (
	"math"

	"github.com/quagmt/udecimal"
	"github.com/yeferson59/gofinance/money"
)

// Periods calculates the number of periods needed using the formula: n = ln(FV/PV) / ln(1 + i)
// where:
//   - FV is the future value
//   - PV is the present value
//   - i is the periodic rate
//
// This method determines how many periods are needed for a capital to grow
// from a present value to a future value with a given rate.
//
// Returns:
//   - The calculated number of periods
//   - An error if there are problems obtaining the valid rate or a mathematical error occurs
//
// Example:
//
//	ci, _ := New(1000, 1126.83, rateInterest, period)
//	periods, err := ci.Periods()
//	// periods will be approximately 12 (for the example with monthly rate)
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

	// Step 1: Calculate the ratio of Future to Present
	futureToPresent := c.future.MustDiv(c.present.Decimal)

	// Step 2: Calculate the natural logarithm of the ratio (numerator)
	logarithmRatio := math.Log(futureToPresent.InexactFloat64())

	// Step 3: Calculate the natural logarithm of the growth factor (denominator)
	growthFactor := periodicRate.Add(udecimal.One)
	logarithmGrowth := math.Log(growthFactor.InexactFloat64())

	// Step 4: Divide to get the number of periods
	numberOfPeriods := money.MustFromFloat64(logarithmRatio).MustDiv(money.MustFromFloat64(logarithmGrowth).Decimal)

	return money.Decimal{Decimal: numberOfPeriods}, nil
}
