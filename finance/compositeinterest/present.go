package compositeinterest

import (
	"math"

	"github.com/yeferson59/gofinance/money"
)

// Present calculates the present value (principal) needed to achieve a target future value.
// The present value is the current amount of money needed to reach a future goal given compound interest.
//
// Formula: PV = FV / (1 + r)^n
// Where:
//   - PV = Present Value
//   - FV = Future Value (target amount)
//   - r = Periodic interest rate
//   - n = Number of periods
//
// Returns:
//   - The present value as a Money instance
//   - An error if the calculation cannot be performed (e.g., missing future value, rate, or periods)
//
// Example:
//
//	ci, _ := New(
//	    money.MoneyZero,
//	    money.MustMoneyFromFloat64(1500, USD),
//	    rateInterest,
//	    period,
//	)
//	present, _ := ci.Present()
//	// present contains the amount needed today to reach $1500
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

	growthFactor := periodicRate.Add(money.One)
	discountFactor := math.Pow(growthFactor.InexactFloat64(), numberOfPeriods.InexactFloat64())

	present := c.future.MustDiv(money.MustMoneyFromFloat64(discountFactor, c.future.Currency()))

	return present, nil
}
