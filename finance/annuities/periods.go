package annuities

import "github.com/yeferson59/gofinance/v2/decimal"

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
//   - The calculated number of periods as decimal.Decimal
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(500, 5000, 0, period, rate)
//	periods, err := ann.PeriodsWithPresent()
//	// periods is how many payment periods needed to pay off $5,000 with $500 payments
func (a Annuity) PeriodsWithPresent() (decimal.Decimal, error) {
	_, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return decimal.Decimal{}, err
	}

	present, err := a.compoundInterest.Present()
	if err != nil {
		return decimal.Decimal{}, err
	}

	ratio, err := a.value.GetDecimal().Div(a.value.GetDecimal().Sub(present.GetDecimal().Mul(rateInterest)))
	if err != nil {
		return decimal.Decimal{}, err
	}

	logarithmRatio, err := ratio.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	logarithmGrowth, err := decimal.One.Add(rateInterest).Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	periods, err := logarithmRatio.Div(logarithmGrowth)
	if err != nil {
		return decimal.Decimal{}, err
	}

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
//   - The calculated number of periods as decimal.Decimal
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(500, 0, 10000, period, rate)
//	periods, err := ann.PeriodsWithFuture()
//	// periods is how many payment periods needed to accumulate $10,000 with $500 payments
func (a Annuity) PeriodsWithFuture() (decimal.Decimal, error) {
	_, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return decimal.Decimal{}, err
	}

	future, err := a.compoundInterest.Future()
	if err != nil {
		return decimal.Decimal{}, err
	}

	ratio, err := future.GetDecimal().Mul(rateInterest).Add(a.value.GetDecimal()).Div(a.value.GetDecimal())
	if err != nil {
		return decimal.Decimal{}, err
	}

	logarithmRatio, err := ratio.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	logarithmGrowth, err := decimal.One.Add(rateInterest).Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	periods, err := logarithmRatio.Div(logarithmGrowth)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return periods, nil
}

// AnticipatePeriodsWithPresent is like PeriodsWithPresent, but assumes each
// payment is made at the beginning of its period (annuity due) instead of
// the end (ordinary annuity).
//
// Since PV_due = PV_ordinary × (1+i) for the same payment and number of
// periods, dividing the present value by (1+i) first reduces this to the
// same ordinary formula: n = ln(PMT / (PMT - [PV/(1+i)] × i)) / ln(1+i)
func (a Annuity) AnticipatePeriodsWithPresent() (decimal.Decimal, error) {
	_, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return decimal.Decimal{}, err
	}

	present, err := a.compoundInterest.Present()
	if err != nil {
		return decimal.Decimal{}, err
	}

	growthFactor := rateInterest.Add(decimal.One)

	equivalentPresent, err := present.GetDecimal().Div(growthFactor)
	if err != nil {
		return decimal.Decimal{}, err
	}

	ratio, err := a.value.GetDecimal().Div(a.value.GetDecimal().Sub(equivalentPresent.Mul(rateInterest)))
	if err != nil {
		return decimal.Decimal{}, err
	}

	logarithmRatio, err := ratio.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	logarithmGrowth, err := growthFactor.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	periods, err := logarithmRatio.Div(logarithmGrowth)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return periods, nil
}

// AnticipatePeriodsWithFuture is like PeriodsWithFuture, but assumes each
// payment is made at the beginning of its period (annuity due) instead of
// the end (ordinary annuity).
//
// Since FV_due = FV_ordinary × (1+i) for the same payment and number of
// periods, dividing the future value by (1+i) first reduces this to the
// same ordinary formula: n = ln(([FV/(1+i)] × i + PMT) / PMT) / ln(1+i)
func (a Annuity) AnticipatePeriodsWithFuture() (decimal.Decimal, error) {
	_, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return decimal.Decimal{}, err
	}

	future, err := a.compoundInterest.Future()
	if err != nil {
		return decimal.Decimal{}, err
	}

	growthFactor := rateInterest.Add(decimal.One)

	equivalentFuture, err := future.GetDecimal().Div(growthFactor)
	if err != nil {
		return decimal.Decimal{}, err
	}

	ratio, err := equivalentFuture.Mul(rateInterest).Add(a.value.GetDecimal()).Div(a.value.GetDecimal())
	if err != nil {
		return decimal.Decimal{}, err
	}

	logarithmRatio, err := ratio.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	logarithmGrowth, err := growthFactor.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	periods, err := logarithmRatio.Div(logarithmGrowth)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return periods, nil
}
