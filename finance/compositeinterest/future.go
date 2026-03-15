package compositeinterest

import (
	"math"

	"github.com/quagmt/udecimal"
	"github.com/yeferson59/gofinance/money"
)

// Future calculates the future value using the compound interest formula: FV = PV * (1 + i)^n
// where:
//   - PV is the present value
//   - i is the periodic rate
//   - n is the number of periods
//
// Returns:
//   - The calculated future value
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ci, _ := New(1000, 0, rateInterest, period)
//	future, err := ci.Future()
//	// If all is correct, future will be the amount after the specified periods
func (c CompositeInterest) Future() (money.Money, error) {
	if !c.future.IsZero() {
		return c.future, nil
	}

	numberOfPeriods, periodicRate, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	if numberOfPeriods.IsZero() || periodicRate.IsZero() || c.present.IsZero() {
		return money.Money{}, ErrInvalidOperation
	}

	// Step 1: Calculate the growth factor (1 + rate)
	growthFactor := periodicRate.Add(udecimal.One)

	// Step 2: Raise the growth factor to the power of the number of periods
	compoundGrowth := math.Pow(growthFactor.InexactFloat64(), numberOfPeriods.InexactFloat64())

	// Step 3: Multiply the present value by the compound growth to get future value
	future := c.present.Mul(udecimal.MustFromFloat64(compoundGrowth))

	return money.Money{Decimal: future}, nil
}
