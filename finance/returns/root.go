// Package returns provides return and growth metrics for investments:
// compound annual growth rate (CAGR), simple return on investment (ROI),
// holding-period return, and annualization of a total return over a number
// of periods.
//
// Every calculation runs on the decimal engine through the money package,
// so results keep money's full fixed-point precision instead of drifting
// through float64. Monetary inputs are money.Money (currency-checked where
// two amounts are compared); rates and ratios are returned as money.Decimal
// expressed as fractions (e.g. 0.12 for 12%).
//
// Basic usage:
//
//	begin := money.MustMoneyFromFloat64(1000, money.USD)
//	end := money.MustMoneyFromFloat64(2000, money.USD)
//	cagr, _ := returns.CAGR(begin, end, money.MustFromFloat64(5))
//	// cagr ≈ 0.1487 (14.87% per year)
package returns

import "errors"

var (
	// ErrNonPositiveValue is returned when a base amount that must be
	// strictly positive (such as the initial value of a CAGR or ROI) is
	// zero or negative, which would make the ratio undefined.
	ErrNonPositiveValue = errors.New("returns: base value must be positive")

	// ErrNonPositivePeriods is returned when the number of periods used to
	// annualize a return is zero or negative.
	ErrNonPositivePeriods = errors.New("returns: periods must be positive")
)
