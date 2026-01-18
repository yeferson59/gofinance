package annuities

import "math"

// Future calculates the future value of an ordinary annuity using the formula:
// FV = PMT × [((1 + i)^n - 1) / i]
// where:
//   - PMT is the periodic payment amount
//   - i is the periodic rate
//   - n is the number of periods
//
// This represents the accumulated value of periodic payments with compound interest.
//
// Returns:
//   - The calculated future value
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(100, 0, 0, period, rate)
//	future, err := ann.Future()
//	// future is the total accumulated value after all periods
func (a Annuity) Future() (float64, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	// Step 1: Calculate the growth factor (1 + rate)
	growthFactor := 1 + rateInterest

	// Step 2: Raise the growth factor to the power of periods
	growthPower := math.Pow(growthFactor, periods)

	// Step 3: Calculate the accumulation factor: ((1 + rate)^n - 1) / rate
	accumulationFactor := (growthPower - 1) / rateInterest

	// Step 4: Multiply the periodic payment by the accumulation factor
	future := a.value * accumulationFactor

	return future, nil
}
