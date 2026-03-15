package compositeinterest

import (
	"math"

	"github.com/quagmt/udecimal"
	"github.com/yeferson59/gofinance/money"
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

	// Step 1: Calculate the ratio of Future to Present
	futureToPresent := c.future.MustDiv(c.present.Decimal)

	// Step 2: Calculate the reciprocal of periods (1/n)
	reciprocalPeriods := udecimal.One.MustDiv(numberOfPeriods.Decimal)

	// Step 3: Raise the ratio to the power of 1/n to get the periodic growth factor
	growthFactor := math.Pow(futureToPresent.InexactFloat64(), reciprocalPeriods.InexactFloat64())

	// Step 4: Subtract 1 to get the periodic rate
	periodicRate := growthFactor - 1

	return money.MustFromFloat64(periodicRate), nil
}
