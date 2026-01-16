package compositeinterest

import (
	"math"
)

// Future calculates the future value using the compound interest formula: FV = PV * (1 + i)^n
// where:
//   - PV is the present value
//   - i is the periodic rate
//   - n is the number of periods
//
// Returns:
//   - The calculated future value
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ci, _ := New(1000, 0, rateInterest, period)
//	future, err := ci.Future()
//	// If all is correct, future will be the amount after the specified periods
func (c CompositeInterest) Future() (float64, error) {
	if c.future != 0 {
		return c.future, nil
	}

	numberOfPeriods, periodicRate, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	if numberOfPeriods == 0 || periodicRate == 0 || c.present == 0 {
		return 0, ErrInvalidOperation
	}

	future := c.present * math.Pow(1+periodicRate, numberOfPeriods)

	return future, nil
}
