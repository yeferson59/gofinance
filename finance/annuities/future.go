package annuities

import (
	"math"

	"github.com/quagmt/udecimal"
	"github.com/yeferson59/gofinance/money"
)

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
//   - The calculated future value as money.Money
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ann, _ := New(100, 0, 0, period, rate)
//	future, err := ann.Future()
//	// future is the total accumulated value after all periods
func (a Annuity) Future() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	// Step 1: Calculate the growth factor (1 + rate)
	growthFactor := rateInterest.Add(udecimal.One)

	// Step 2: Raise the growth factor to the power of periods
	growthPower := udecimal.MustFromFloat64(math.Pow(growthFactor.InexactFloat64(), periods.InexactFloat64()))

	// Step 3: Calculate the accumulation factor: ((1 + rate)^n - 1) / rate
	quotient, err := growthPower.Sub(udecimal.One).Div(rateInterest.Decimal)
	if err != nil {
		return money.Money{}, err
	}
	accumulationFactor := quotient

	// Step 4: Multiply the periodic payment by the accumulation factor
	futureDecimal := a.value.Decimal.Mul(accumulationFactor)

	return money.Money{Decimal: futureDecimal}, nil
}
