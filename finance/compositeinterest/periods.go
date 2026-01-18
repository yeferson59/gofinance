package compositeinterest

import (
	"math"
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
func (c CompositeInterest) Periods() (float64, error) {
	if periodValue, _, err := c.periods.getPeriod(); err == nil && periodValue != 0 {
		return periodValue, nil
	}

	_, periodicRate, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	if c.present == 0 || c.future == 0 || periodicRate == 0 {
		return 0, ErrInvalidOperation
	}

	// Step 1: Calculate the ratio of Future to Present
	futureToPresent := c.future / c.present

	// Step 2: Calculate the natural logarithm of the ratio (numerator)
	logarithmRatio := math.Log(futureToPresent)

	// Step 3: Calculate the natural logarithm of the growth factor (denominator)
	growthFactor := 1 + periodicRate
	logarithmGrowth := math.Log(growthFactor)

	// Step 4: Divide to get the number of periods
	numberOfPeriods := logarithmRatio / logarithmGrowth

	return numberOfPeriods, nil
}
