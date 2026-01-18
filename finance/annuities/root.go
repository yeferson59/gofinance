// Package annuities provides functionality for annuity calculations.
//
// This package enables financial calculations related to ordinary annuities,
// including:
//   - Calculation of periodic payments from present value (loan amortization)
//   - Calculation of periodic payments from future value (savings accumulation)
//   - Calculation of present value of an annuity
//   - Calculation of future value of an annuity
//   - Calculation of the number of periods needed
//
// An ordinary annuity is a series of equal periodic payments made at the end of each period.
// The package supports multiple compounding frequencies:
//   - Daily (365 periods per year)
//   - Monthly (12 periods per year)
//   - Bimonthly (6 periods per year)
//   - Quarterly (4 or 3 periods per year)
//   - Semi-annually (2 periods per year)
//   - Annually (1 period per year)
//
// Basic usage example:
//
//	// Create a monthly interest rate of 1% periodic
//	rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
//
//	// Create a period of 12 months
//	period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)
//
//	// Create an annuity with $100 monthly payment
//	ann, _ := New(100, 0, 0, period, rate)
//
//	// Calculate the present value
//	present, _ := ann.Present()
//	// present ≈ 1125.51
package annuities

import (
	"math"

	"github.com/yeferson59/gofinance/finance/compositeinterest"
)

// Annuity represents an ordinary annuity with periodic equal payments.
// It stores a periodic payment amount and the underlying composite interest calculation.
type Annuity struct {
	// value is the periodic payment amount for the annuity
	value float64
	// compositeInterest holds the underlying composite interest calculation
	compositeInterest compositeinterest.CompositeInterest
}

// New creates a new Annuity instance.
//
// Parameters:
//   - value: The periodic payment amount
//   - present: The present value (for loan calculations) or 0
//   - future: The future value (for savings calculations) or 0
//   - period: The period structure containing the number of periods and frequency
//   - rateInterest: The rate interest structure containing the rate and frequency
//
// Returns:
//   - An Annuity instance
//   - An error if the composite interest creation fails
//
// Example:
//
//	rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
//	period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)
//	ann, err := New(100, 0, 0, period, rate)
//	if err != nil {
//	    log.Fatal(err)
//	}
func New(value, present, future float64, period compositeinterest.Period, rateInterest compositeinterest.RateInterest) (Annuity, error) {
	compositeinterest, err := compositeinterest.New(present, future, rateInterest, period)
	if err != nil {
		return Annuity{}, err
	}

	return Annuity{
		value:             value,
		compositeInterest: compositeinterest,
	}, nil
}

// PaymentFromPresentValue calculates the periodic payment needed to amortize a present value (loan payment).
// Uses the formula: PMT = PV × [i(1 + i)^n] / [(1 + i)^n - 1]
// where:
//   - PV is the present value
//   - i is the periodic rate
//   - n is the number of periods
//
// Returns:
//   - The calculated periodic payment amount
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(0, 10000, 0, period, rate) // Loan of $10,000
//	payment, err := ann.PaymentFromPresentValue()
//	// payment is the monthly amount needed to pay off the loan
func (a Annuity) PaymentFromPresentValue() (float64, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return 0, err
	}

	// Step 1: Calculate the growth factor (1 + rate)
	growthFactor := 1 + rateInterest

	// Step 2: Raise the growth factor to the power of periods
	growthPower := math.Pow(growthFactor, periods)

	// Step 3: Calculate the numerator: rate × (1 + rate)^n
	numerator := rateInterest * growthPower

	// Step 4: Calculate the denominator: (1 + rate)^n - 1
	denominator := growthPower - 1

	// Step 5: Calculate the annuity payment: PV × [numerator / denominator]
	annuity := present * (numerator / denominator)

	return annuity, nil
}

// PaymentFromFutureValue calculates the periodic payment needed to accumulate a future value (savings payment).
// Uses the formula: PMT = FV × [i / ((1 + i)^n - 1)]
// where:
//   - FV is the future value
//   - i is the periodic rate
//   - n is the number of periods
//
// Returns:
//   - The calculated periodic payment amount
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(0, 0, 10000, period, rate) // Goal: accumulate $10,000
//	payment, err := ann.PaymentFromFutureValue()
//	// payment is the monthly amount needed to reach the goal
func (a Annuity) PaymentFromFutureValue() (float64, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return 0, err
	}

	// Step 1: Calculate the growth factor (1 + rate)
	growthFactor := 1 + rateInterest

	// Step 2: Raise the growth factor to the power of periods
	growthPower := math.Pow(growthFactor, periods)

	// Step 3: Calculate the denominator: (1 + rate)^n - 1
	denominator := growthPower - 1

	// Step 4: Calculate the annuity payment: FV × [rate / denominator]
	annuity := future * (rateInterest / denominator)

	return annuity, nil
}
