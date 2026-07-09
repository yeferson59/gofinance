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

	presentFloat := present.InexactFloat64()
	valueFloat := a.value.InexactFloat64()
	rateFloat := rateInterest.InexactFloat64()

	// Step 1: Calculate the denominator base: Present × rate
	presentTimesRate := presentFloat * rateFloat

	// Step 2: Calculate the denominator: PMT - (PV × rate)
	denominatorValue := valueFloat - presentTimesRate

	// Step 3: Calculate the ratio: PMT / (PMT - (PV × rate))
	ratio := valueFloat / denominatorValue

	// Step 4: Calculate the natural logarithm of the ratio (numerator)
	logarithmRatio := money.MustFromFloat64(ratio).MustLn().InexactFloat64()

	// Step 5: Calculate the natural logarithm of the growth factor (denominator)
	growthFactor := 1 + rateFloat
	logarithmGrowth := money.MustFromFloat64(growthFactor).MustLn().InexactFloat64()

	// Step 6: Divide to get the number of periods
	periods := logarithmRatio / logarithmGrowth

	return money.MustFromFloat64(periods), nil
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

	futureFloat := future.InexactFloat64()
	valueFloat := a.value.InexactFloat64()
	rateFloat := rateInterest.InexactFloat64()

	// Step 1: Calculate the numerator base: Future × rate
	futureTimesRate := futureFloat * rateFloat

	// Step 2: Calculate the numerator: (FV × rate) + PMT
	numeratorValue := futureTimesRate + valueFloat

	// Step 3: Calculate the natural logarithm of the numerator
	logarithmNumerator := money.MustFromFloat64(numeratorValue).MustLn().InexactFloat64()

	// Step 4: Calculate the natural logarithm of the denominator (PMT)
	logarithmDenominator := money.MustFromFloat64(valueFloat).MustLn().InexactFloat64()

	// Step 5: Calculate the numerator of the periods formula
	logarithmRatio := logarithmNumerator - logarithmDenominator

	// Step 6: Calculate the natural logarithm of the growth factor
	growthFactor := 1 + rateFloat
	logarithmGrowth := money.MustFromFloat64(growthFactor).MustLn().InexactFloat64()

	// Step 7: Divide to get the number of periods
	periods := logarithmRatio / logarithmGrowth

	return money.MustFromFloat64(periods), nil
}
