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

	pow := math.Pow(1+rateInterest, periods)
	present := a.value * ((pow - 1) / (rateInterest * pow))

	return present, nil
}
