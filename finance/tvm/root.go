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

	"github.com/yeferson59/gofinance/decimal"
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
	t.n = decimal.MustFromFloat64(n)
	return t
}

// Rate sets the interest rate per period as a fraction (e.g. 0.005 for
// 0.5% per period).
func (t Config) Rate(rate float64) Config {
	t.rate = decimal.MustFromFloat64(rate)
	return t
}

// PV sets the present value.
func (t Config) PV(pv float64) Config {
	t.pv = decimal.MustFromFloat64(pv)
	return t
}

// PMT sets the payment made each period.
func (t Config) PMT(pmt float64) Config {
	t.pmt = decimal.MustFromFloat64(pmt)
	return t
}

// FV sets the future value.
func (t Config) FV(fv float64) Config {
	t.fv = decimal.MustFromFloat64(fv)
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
