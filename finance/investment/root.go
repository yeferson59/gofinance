// Package investment provides discounted cash-flow metrics for evaluating a
// stream of periodic cash flows: net present value (NPV) and internal rate of
// return (IRR).
//
// Cash flows are a slice of money.Money indexed by period, all in the same
// currency: index 0 is the flow today (period 0, not discounted), index 1 the
// flow one period later, and so on. By convention an outflow (an investment)
// is negative and an inflow (a return) is positive.
//
// All discounting runs on the decimal engine through the money package, so
// results keep full fixed-point precision. Rates are money.Decimal fractions
// per period (e.g. 0.01 for 1% per period).
//
// Basic usage:
//
//	flows := []money.Money{
//	    money.MustMoneyFromFloat64(-1000, money.USD), // invest today
//	    money.MustMoneyFromFloat64(400, money.USD),   // returns each period
//	    money.MustMoneyFromFloat64(400, money.USD),
//	    money.MustMoneyFromFloat64(400, money.USD),
//	}
//	npv, _ := investment.NPV(money.MustFromFloat64(0.10), flows)
//	irr, _ := investment.IRR(flows)
package investment

import (
	"errors"

	"github.com/yeferson59/gofinance/money"
)

var (
	// ErrNoCashFlows is returned when an empty cash-flow slice is passed.
	ErrNoCashFlows = errors.New("investment: at least one cash flow is required")

	// ErrInvalidRate is returned when a discount rate is not greater than −1,
	// which would make the discount factor (1+rate) zero or negative.
	ErrInvalidRate = errors.New("investment: rate must be greater than -1")

	// ErrNoSignChange is returned by IRR when the cash flows never change
	// sign, so no internal rate of return exists.
	ErrNoSignChange = errors.New("investment: cash flows must contain at least one sign change for an IRR to exist")

	// ErrNoConvergence is returned by IRR when neither the Newton iteration
	// nor the bracketed search locate a root within the iteration budget.
	ErrNoConvergence = errors.New("investment: IRR did not converge")
)

// decimalFlows validates that cashFlows is non-empty and single-currency, and
// returns the amounts as decimals alongside their shared currency.
func decimalFlows(cashFlows []money.Money) ([]money.Decimal, money.Currency, error) {
	if len(cashFlows) == 0 {
		return nil, 0, ErrNoCashFlows
	}

	currency := cashFlows[0].Currency()
	amounts := make([]money.Decimal, len(cashFlows))

	for i, cf := range cashFlows {
		if cf.Currency() != currency {
			return nil, 0, money.ErrCurrencyMismatch
		}

		amounts[i] = cf.ToDecimal()
	}

	return amounts, currency, nil
}
