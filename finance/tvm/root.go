// Package tvm is a general time-value-of-money solver in the style of a
// financial calculator (HP-12C, TI BA II). It relates the five standard TVM
// variables and solves for any one of them given the other four:
//
//   - N    — number of periods
//   - Rate — interest rate per period (a fraction, e.g. 0.005 for 0.5%/month)
//   - PV   — present value
//   - PMT  — payment made each period
//   - FV   — future value
//
// They satisfy the cash-flow equation
//
//	PV·(1+i)ᴺ + PMT·(1 + i·type)·((1+i)ᴺ − 1)/i + FV = 0
//
// where type is 0 for an ordinary annuity (payments at the end of each
// period, the default) and 1 for an annuity due (payments at the beginning).
// Follow the sign convention of financial calculators: money you receive is
// positive and money you pay out is negative, so a solved payment or future
// value typically carries the opposite sign of the present value.
//
// Everything runs on the decimal engine, so the closed-form solves are exact
// to the engine's precision and the iterative rate solve converges against
// exact residuals. Values are decimal.Decimal; the fluent builder accepts
// float64 for convenience.
//
// Example — monthly payment on a $300,000 loan at 6%/yr over 30 years:
//
//	pmt := tvm.NewTVM().
//	    PV(300000).
//	    Rate(0.06 / 12).
//	    N(360).
//	    MustSolvePMT()
//	// pmt ≈ -1798.65 (a cash outflow)
package tvm

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/decimal"
)

var (
	// ErrInvalidRate is returned when a rate is not greater than −1, which
	// makes the growth factor (1+rate) zero or negative.
	ErrInvalidRate = errors.New("tvm: rate must be greater than -1")

	// ErrNoSolution is returned when the inputs are degenerate and no finite
	// value of the requested variable satisfies the equation (for example,
	// solving for PMT with zero periods, or for N with a zero payment at a
	// zero rate).
	ErrNoSolution = errors.New("tvm: no solution for the given inputs")

	// ErrNoConvergence is returned by SolveRate when no rate satisfying the
	// equation can be bracketed.
	ErrNoConvergence = errors.New("tvm: rate did not converge")
)

// Config is a fluent builder holding the known TVM variables. Create one
// with NewTVM, set the four you know, then call the Solve* method for the
// fifth.
type Config struct {
	n    decimal.Decimal
	rate decimal.Decimal
	pv   decimal.Decimal
	pmt  decimal.Decimal
	fv   decimal.Decimal
	due  bool

	// err holds the first invalid input a setter received (a NaN or infinite
	// float, say). Setters can't return an error without breaking the chain,
	// so it is kept here and reported by every Solve method.
	err error
}

// withErr records err as the builder's error unless one is already held, so
// the first problem is the one reported. A nil err leaves t unchanged.
func (t Config) withErr(err error) Config {
	if t.err == nil {
		t.err = err
	}

	return t
}

// NewTVM returns a Config with every variable initialized to zero and an
// ordinary-annuity (end-of-period) payment timing.
func NewTVM() Config {
	return Config{
		n:    decimal.Zero,
		rate: decimal.Zero,
		pv:   decimal.Zero,
		pmt:  decimal.Zero,
		fv:   decimal.Zero,
	}
}

// N sets the number of periods.
func (t Config) N(n float64) Config {
	d, err := decimal.NewFromFloat64(n)

	return t.withErr(err).NDecimal(d)
}

// NDecimal is like N but takes the number of periods as a decimal.Decimal, so a
// value a float64 can't hold exactly keeps every digit.
func (t Config) NDecimal(n decimal.Decimal) Config {
	t.n = n
	return t
}

// Rate sets the interest rate per period as a fraction (e.g. 0.005 for
// 0.5% per period).
func (t Config) Rate(rate float64) Config {
	d, err := decimal.NewFromFloat64(rate)

	return t.withErr(err).RateDecimal(d)
}

// RateDecimal is like Rate but takes the rate per period as a decimal.Decimal, so a
// value a float64 can't hold exactly keeps every digit.
func (t Config) RateDecimal(rate decimal.Decimal) Config {
	t.rate = rate
	return t
}

// PV sets the present value.
func (t Config) PV(pv float64) Config {
	d, err := decimal.NewFromFloat64(pv)

	return t.withErr(err).PVDecimal(d)
}

// PVDecimal is like PV but takes the present value as a decimal.Decimal, so a
// value a float64 can't hold exactly keeps every digit.
func (t Config) PVDecimal(pv decimal.Decimal) Config {
	t.pv = pv
	return t
}

// PMT sets the payment made each period.
func (t Config) PMT(pmt float64) Config {
	d, err := decimal.NewFromFloat64(pmt)

	return t.withErr(err).PMTDecimal(d)
}

// PMTDecimal is like PMT but takes the payment as a decimal.Decimal, so a
// value a float64 can't hold exactly keeps every digit.
func (t Config) PMTDecimal(pmt decimal.Decimal) Config {
	t.pmt = pmt
	return t
}

// FV sets the future value.
func (t Config) FV(fv float64) Config {
	d, err := decimal.NewFromFloat64(fv)

	return t.withErr(err).FVDecimal(d)
}

// FVDecimal is like FV but takes the future value as a decimal.Decimal, so a
// value a float64 can't hold exactly keeps every digit.
func (t Config) FVDecimal(fv decimal.Decimal) Config {
	t.fv = fv
	return t
}

// Due switches payment timing to an annuity due (payments at the beginning of
// each period).
func (t Config) Due() Config {
	t.due = true
	return t
}

// Ordinary switches payment timing to an ordinary annuity (payments at the end
// of each period). This is the default.
func (t Config) Ordinary() Config {
	t.due = false
	return t
}
