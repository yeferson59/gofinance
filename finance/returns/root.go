// Package returns provides return and growth metrics for investments:
// compound annual growth rate (CAGR), simple return on investment (ROI),
// holding-period return, and annualization of a total return over a number
// of periods.
//
// Every calculation runs on the decimal engine, so results keep full
// fixed-point precision instead of drifting through float64. Monetary inputs are money.Money (currency-checked where
// two amounts are compared); rates and ratios are returned as decimal.Decimal
// expressed as fractions (e.g. 0.12 for 12%).
//
// Basic usage:
//
//	begin := money.MustMoneyFromFloat64(1000, money.USD)
//	end := money.MustMoneyFromFloat64(2000, money.USD)
//	cagr, _ := returns.CAGR(begin, end, decimal.MustFromFloat64(5))
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

	// ErrInvalidInflationRate is returned when an inflation rate is not greater
	// than −1, which would make the price-level factor (1+inflation) zero or
	// negative.
	ErrInvalidInflationRate = errors.New("returns: inflation rate must be greater than -1")

	// ErrNoSubperiods is returned when a time-weighted return is asked for
	// without a single valuation subperiod to link.
	ErrNoSubperiods = errors.New("returns: at least one subperiod is required")

	// ErrNoReturns is returned when a statistic is asked for over an empty
	// series of returns.
	ErrNoReturns = errors.New("returns: at least one return is required")

	// ErrInsufficientReturns is returned when a dispersion measure that needs
	// more than one observation (a sample variance, volatility, or Sharpe
	// ratio) is given a single return.
	ErrInsufficientReturns = errors.New("returns: at least two returns are required")

	// ErrZeroVolatility is returned by the Sharpe ratio when the returns never
	// vary, leaving no risk to divide the excess return by.
	ErrZeroVolatility = errors.New("returns: volatility is zero")
)
