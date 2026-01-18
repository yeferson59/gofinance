package annuities

import "math"

// Present calculates the present value of an ordinary annuity using the formula:
// PV = PMT × [(1 + i)^n - 1] / [i(1 + i)^n]
// where:
//   - PMT is the periodic payment amount
//   - i is the periodic rate
//   - n is the number of periods
//
// This represents the current value of a series of future periodic payments.
//
// Returns:
//   - The calculated present value
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(100, 0, 0, period, rate)
//	present, err := ann.Present()
//	// present is the current value of receiving $100 for n periods
func (a Annuity) Present() (float64, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	// Step 1: Calculate the growth factor (1 + rate)
	growthFactor := 1 + rateInterest

	// Step 2: Raise the growth factor to the power of periods
	growthPower := math.Pow(growthFactor, periods)

	// Step 3: Calculate the numerator: (1 + rate)^n - 1
	numerator := growthPower - 1

	// Step 4: Calculate the denominator: rate × (1 + rate)^n
	denominator := rateInterest * growthPower

	// Step 5: Calculate the present value: PMT × [numerator / denominator]
	present := a.value * (numerator / denominator)

	return present, nil
}
