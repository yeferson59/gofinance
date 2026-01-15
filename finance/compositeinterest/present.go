package compositeinterest

import (
	"math"
)

// Present calculates the present value (discount) using the formula: PV = FV / (1 + i)^n
// where:
//   - FV is the future value
//   - i is the periodic rate
//   - n is the number of periods
//
// This method is the inverse of Future() and is used to determine how much money
// needs to be invested today to achieve a specific future value.
//
// Returns:
//   - The calculated present value
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ci, _ := New(0, 1126.83, rateInterest, period)
//	present, err := ci.Present()
//	// present will be approximately 1000
func (c CompositeInterest) Present() (float64, error) {
	if c.present != 0 {
		return c.present, nil
	}

	periods, rateInterest, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	if c.future == 0 || rateInterest == 0 || periods == 0 {
		return 0, ErrInvalidOperation
	}

	present := (c.future / math.Pow(1+rateInterest, periods))

	return present, nil
}
