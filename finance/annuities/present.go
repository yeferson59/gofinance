package annuities

import (
	"math"

	"github.com/quagmt/udecimal"
	"github.com/yeferson59/gofinance/money"
)

// Present calculates the present value of an ordinary annuity using the formula:
// PV = PMT × [1 - (1 + i)^-n] / i
// where:
//   - PMT is the periodic payment amount
//   - i is the periodic rate
//   - n is the number of periods
//
// This represents the current value of a series of future periodic payments.
//
// Returns:
//   - The calculated present value as money.Decimal
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(100, 0, 0, period, rate)
//	present, err := ann.Present()
//	// present is the current value of receiving $100 for n periods
func (a Annuity) Present() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	// Step 1: Calculate the growth factor (1 + rate)
	growthFactor := rateInterest.Add(udecimal.One)

	// Step 2: Raise the growth factor to the power of periods
	growthPower := udecimal.MustFromFloat64(math.Pow(growthFactor.InexactFloat64(), periods.InexactFloat64()))

	// Step 3: Calculate the numerator: 1 - (1 + rate)^-n
	numerator := udecimal.One.Sub(udecimal.One.MustDiv(growthPower))

	// Step 4: Calculate the denominator: rate
	denominator := rateInterest.Decimal

	// Step 5: Calculate the present value: PMT × [numerator / denominator]
	quotient, err := numerator.Div(denominator)
	if err != nil {
		return money.Money{}, err
	}
	presentDecimal := a.value.Decimal.Mul(quotient)

	return money.Money{Decimal: presentDecimal}, nil
}
