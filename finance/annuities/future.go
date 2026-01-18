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

	future := a.value * ((math.Pow(1+rateInterest, periods) - 1) / rateInterest)

	return future, nil
}
