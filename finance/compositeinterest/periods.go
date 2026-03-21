package compositeinterest

import (
	"math"

	"github.com/yeferson59/gofinance/money"
)

// Periods calculates the number of compounding periods needed to reach a target future value
// from a present value at a given interest rate.
//
// Formula: n = ln(FV/PV) / ln(1 + r)
// Where:
//   - n = Number of periods
//   - FV = Future Value
//   - PV = Present Value
//   - r = Periodic interest rate
//
// Returns:
//   - The number of periods as a Decimal
//   - An error if the calculation cannot be performed (e.g., missing values, zero rate)
//
// Example:
//
//	ci, _ := New(
//	    money.MustMoneyFromFloat64(1000, USD),
//	    money.MustMoneyFromFloat64(1500, USD),
//	    rateInterest,
//	    period,
//	)
//	periods, _ := ci.Periods()
//	// periods contains how many periods needed to reach $1500 from $1000
func (c CompositeInterest) Periods() (money.Decimal, error) {
	if periodValue, _, err := c.periods.getPeriod(); err == nil && !periodValue.IsZero() {
		return periodValue, nil
	}

	_, periodicRate, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	if c.present.IsZero() || c.future.IsZero() || periodicRate.IsZero() {
		return money.Decimal{}, ErrInvalidOperation
	}

	futureToPresent := c.future.MustDiv(c.present)

	logarithmRatio := math.Log(futureToPresent.InexactFloat64())

	growthFactor := periodicRate.Add(money.One)
	logarithmGrowth := math.Log(growthFactor.InexactFloat64())

	numberOfPeriods := money.MustFromFloat64(logarithmRatio).MustDiv(money.MustFromFloat64(logarithmGrowth))

	return numberOfPeriods, nil
}
