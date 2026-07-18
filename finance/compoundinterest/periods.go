package compoundinterest

import "github.com/yeferson59/gofinance/decimal"

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
func (c CompoundInterest) Periods() (decimal.Decimal, error) {
	if periodValue, _, err := c.periods.getPeriod(); err == nil && !periodValue.IsZero() {
		return periodValue, nil
	}

	_, periodicRate, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if c.present.IsZero() || c.future.IsZero() || periodicRate.IsZero() {
		return decimal.Decimal{}, ErrInvalidOperation
	}

	futureToPresent, err := c.future.ToDecimal().Div(c.present.ToDecimal())
	if err != nil {
		return decimal.Decimal{}, err
	}

	logarithmRatio, err := futureToPresent.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	growthFactor := periodicRate.Add(decimal.One)
	logarithmGrowth, err := growthFactor.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	numberOfPeriods, err := logarithmRatio.Div(logarithmGrowth)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return numberOfPeriods, nil
}
