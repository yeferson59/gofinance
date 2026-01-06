package compositeinterest

import (
	"errors"
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
func (c *CompositeInterest) Periods() (float64, error) {
	_, rateInterest, err := c.getEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	if c.present == 0 {
		return 0, errors.New("invalid present for operation")
	}

	periods := (math.Log((c.future / c.present)) / math.Log(1+rateInterest))

	return periods, nil
}
