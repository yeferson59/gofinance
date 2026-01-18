package compositeinterest

import (
	"math"
)

// Interest calculates the periodic interest rate using the formula: i = (FV/PV)^(1/n) - 1
// where:
//   - FV is the future value
//   - PV is the present value
//   - n is the number of periods
//
// This method determines what periodic rate was applied between the specified present and future values.
//
// Returns:
//   - The periodic rate as decimal (e.g., 0.01 for 1%)
//   - An error if there are problems obtaining a valid period
//
// Example:
//
//	ci, _ := New(1000, 1126.83, rateInterest, period)
//	rate, err := ci.Interest()
//	// rate will be approximately 0.01 (1% periodic)
func (c CompositeInterest) Interest() (float64, error) {
	if c.rateInterest.value != 0 {
		return c.rateInterest.value, nil
	}

	numberOfPeriods, _, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	if c.future == 0 || c.present == 0 || numberOfPeriods == 0 {
		return 0, ErrInvalidOperation
	}

	// Step 1: Calculate the ratio of Future to Present
	futureToPresent := c.future / c.present

	// Step 2: Calculate the reciprocal of periods (1/n)
	reciprocalPeriods := 1 / numberOfPeriods

	// Step 3: Raise the ratio to the power of 1/n to get the periodic growth factor
	growthFactor := math.Pow(futureToPresent, reciprocalPeriods)

	// Step 4: Subtract 1 to get the periodic rate
	periodicRate := growthFactor - 1

	return periodicRate, nil
}
