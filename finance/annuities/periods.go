package annuities

import "math"

// PeriodsWithPresent calculates the number of periods needed for periodic payments to reach
// a specific present value using the formula:
// n = ln(PMT / (PMT - PV × i)) / ln(1 + i)
// where:
//   - PMT is the periodic payment amount
//   - PV is the present value
//   - i is the periodic rate
//   - n is the number of periods
//
// This is useful for determining how long it will take to pay off a loan.
//
// Returns:
//   - The calculated number of periods
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(500, 5000, 0, period, rate)
//	periods, err := ann.PeriodsWithPresent()
//	// periods is how many payment periods needed to pay off $5,000 with $500 payments
func (a Annuity) PeriodsWithPresent() (float64, error) {
	_, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return 0, err
	}

	logBase := math.Log(1 + rateInterest)
	periods := (math.Log(a.value/(a.value-(present*rateInterest))) / logBase)

	return periods, nil
}

// PeriodsWithFuture calculates the number of periods needed for periodic payments to reach
// a specific future value using the formula:
// n = ln((FV × i + PMT) / PMT) / ln(1 + i)
// where:
//   - FV is the future value (goal amount)
//   - PMT is the periodic payment amount
//   - i is the periodic rate
//   - n is the number of periods
//
// This is useful for determining how long it will take to accumulate a target savings amount.
//
// Returns:
//   - The calculated number of periods
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(500, 0, 10000, period, rate)
//	periods, err := ann.PeriodsWithFuture()
//	// periods is how many payment periods needed to accumulate $10,000 with $500 payments
func (a Annuity) PeriodsWithFuture() (float64, error) {
	_, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return 0, err
	}

	logBase := math.Log(1 + rateInterest)
	periods := ((math.Log((rateInterest*future)+a.value) - math.Log(a.value)) / logBase)

	return periods, nil
}
