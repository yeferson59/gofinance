package compositeinterest

import (
	"github.com/yeferson59/gofinance/money"
)

// Interest calculates the periodic interest rate given present value, future value, and number of periods.
// This method derives the interest rate from the compound interest formula.
//
// Formula: r = (FV/PV)^(1/n) - 1
// Where:
//   - r = Periodic interest rate
//   - FV = Future Value
//   - PV = Present Value
//   - n = Number of periods
//
// Returns:
//   - The periodic interest rate as a Decimal
//   - An error if the calculation cannot be performed (e.g., missing values)
//
// Example:
//
//	ci, _ := New(
//	    money.MustMoneyFromFloat64(1000, USD),
//	    money.MustMoneyFromFloat64(1268, USD),
//	    rateInterest,
//	    period,
//	)
//	rate, _ := ci.Interest()
//	// rate contains the calculated periodic interest rate
func (c CompositeInterest) Interest() (money.Decimal, error) {
	if !c.rateInterest.value.IsZero() {
		return c.rateInterest.value, nil
	}

	numberOfPeriods, _, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	if c.future.IsZero() || c.present.IsZero() || numberOfPeriods.IsZero() {
		return money.Decimal{}, ErrInvalidOperation
	}

	futureToPresent, err := c.future.Div(c.present)
	if err != nil {
		return money.Decimal{}, err
	}

	reciprocalPeriods, err := money.One.Div(numberOfPeriods)
	if err != nil {
		return money.Decimal{}, err
	}

	growthFactor, err := futureToPresent.ToDecimal().Pow(reciprocalPeriods)
	if err != nil {
		return money.Decimal{}, err
	}

	return growthFactor.Sub(money.One), nil
}
