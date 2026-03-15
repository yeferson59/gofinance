package compositeinterest

import (
	"math"

	"github.com/quagmt/udecimal"
	"github.com/yeferson59/gofinance/money"
)

// Present calculates the present value (discount) using the formula: PV = FV / (1 + i)^n
// where:
//   - FV is the future value
//   - i is the periodic rate
//   - n is the number of periods
//
// This method is the inverse of Future() and is used to determine how much money
// needs to be invested today to achieve a specific future value.
//
// Returns:
//   - The calculated present value
//   - An error if there are problems obtaining valid rate or period values
//
// Example:
//
//	ci, _ := New(0, 1126.83, rateInterest, period)
//	present, err := ci.Present()
//	// present will be approximately 1000
func (c CompositeInterest) Present() (money.Money, error) {
	if !c.present.IsZero() {
		return c.present, nil
	}

	numberOfPeriods, periodicRate, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	if c.future.IsZero() || periodicRate.IsZero() || numberOfPeriods.IsZero() {
		return money.Money{}, ErrInvalidOperation
	}

	// Step 1: Calculate the growth factor (1 + rate)
	growthFactor := periodicRate.Add(udecimal.One)

	// Step 2: Raise the growth factor to the power of the number of periods
	discountFactor := math.Pow(growthFactor.InexactFloat64(), numberOfPeriods.InexactFloat64())

	// Step 3: Divide the future value by the discount factor to get present value
	present := c.future.MustDiv(udecimal.MustFromFloat64(discountFactor))

	return money.Money{Decimal: present}, nil
}
