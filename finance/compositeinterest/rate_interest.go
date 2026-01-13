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
	periods, _, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	rateInterest := math.Pow((c.future/c.present), (1/periods)) - 1

	return rateInterest, nil
}
