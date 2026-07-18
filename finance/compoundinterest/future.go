package compoundinterest

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// Future calculates the future value of an investment or loan based on compound interest.
// The future value is the amount of money an investment will grow to over time.
//
// Formula: FV = PV × (1 + r)^n
// Where:
//   - FV = Future Value
//   - PV = Present Value (principal)
//   - r = Periodic interest rate
//   - n = Number of periods
//
// Returns:
//   - The future value as a Money instance
//   - An error if the calculation cannot be performed (e.g., missing present value, rate, or periods)
//
// Example:
//
//	ci, _ := New(
//	    money.MustMoneyFromFloat64(1000, USD),
//	    money.MoneyZero,
//	    rateInterest,
//	    period,
//	)
//	future, _ := ci.Future()
//	// future contains the calculated future value
func (c CompoundInterest) Future() (money.Money, error) {
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

	growthFactor := periodicRate.Add(decimal.One)
	compoundGrowth, err := growthFactor.Pow(numberOfPeriods)
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(c.present.ToDecimal().Mul(compoundGrowth), c.present.Currency()), nil
}
