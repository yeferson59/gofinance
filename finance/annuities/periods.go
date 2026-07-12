package annuities

import (
	"github.com/yeferson59/gofinance/money"
)

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
//   - The calculated number of periods as money.Decimal
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(500, 5000, 0, period, rate)
//	periods, err := ann.PeriodsWithPresent()
//	// periods is how many payment periods needed to pay off $5,000 with $500 payments
func (a Annuity) PeriodsWithPresent() (money.Decimal, error) {
	_, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return money.Decimal{}, err
	}

	ratio, err := a.value.ToDecimal().Div(a.value.ToDecimal().Sub(present.ToDecimal().Mul(rateInterest)))
	if err != nil {
		return money.Decimal{}, err
	}

	logarithmRatio, err := ratio.Ln()
	if err != nil {
		return money.Decimal{}, err
	}

	logarithmGrowth, err := money.One.Add(rateInterest).Ln()
	if err != nil {
		return money.Decimal{}, err
	}

	periods, err := logarithmRatio.Div(logarithmGrowth)
	if err != nil {
		return money.Decimal{}, err
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
//   - The calculated number of periods as money.Decimal
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(500, 0, 10000, period, rate)
//	periods, err := ann.PeriodsWithFuture()
//	// periods is how many payment periods needed to accumulate $10,000 with $500 payments
func (a Annuity) PeriodsWithFuture() (money.Decimal, error) {
	_, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return money.Decimal{}, err
	}

	ratio, err := future.ToDecimal().Mul(rateInterest).Add(a.value.ToDecimal()).Div(a.value.ToDecimal())
	if err != nil {
		return money.Decimal{}, err
	}

	logarithmRatio, err := ratio.Ln()
	if err != nil {
		return money.Decimal{}, err
	}

	logarithmGrowth, err := money.One.Add(rateInterest).Ln()
	if err != nil {
		return money.Decimal{}, err
	}

	periods, err := logarithmRatio.Div(logarithmGrowth)
	if err != nil {
		return money.Decimal{}, err
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
func (a Annuity) AnticipatePeriodsWithPresent() (money.Decimal, error) {
	_, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return money.Decimal{}, err
	}

	growthFactor := rateInterest.Add(money.One)

	equivalentPresent, err := present.ToDecimal().Div(growthFactor)
	if err != nil {
		return money.Decimal{}, err
	}

	ratio, err := a.value.ToDecimal().Div(a.value.ToDecimal().Sub(equivalentPresent.Mul(rateInterest)))
	if err != nil {
		return money.Decimal{}, err
	}

	logarithmRatio, err := ratio.Ln()
	if err != nil {
		return money.Decimal{}, err
	}

	logarithmGrowth, err := growthFactor.Ln()
	if err != nil {
		return money.Decimal{}, err
	}

	periods, err := logarithmRatio.Div(logarithmGrowth)
	if err != nil {
		return money.Decimal{}, err
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
func (a Annuity) AnticipatePeriodsWithFuture() (money.Decimal, error) {
	_, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return money.Decimal{}, err
	}

	growthFactor := rateInterest.Add(money.One)

	equivalentFuture, err := future.ToDecimal().Div(growthFactor)
	if err != nil {
		return money.Decimal{}, err
	}

	ratio, err := equivalentFuture.Mul(rateInterest).Add(a.value.ToDecimal()).Div(a.value.ToDecimal())
	if err != nil {
		return money.Decimal{}, err
	}

	logarithmRatio, err := ratio.Ln()
	if err != nil {
		return money.Decimal{}, err
	}

	logarithmGrowth, err := growthFactor.Ln()
	if err != nil {
		return money.Decimal{}, err
	}

	periods, err := logarithmRatio.Div(logarithmGrowth)
	if err != nil {
		return money.Decimal{}, err
	}

	return periods, nil
}
