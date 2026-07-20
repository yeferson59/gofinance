// Package gradients provides present and future value calculations for
// gradient payment series: sequences of periodic payments that change from
// one period to the next by a constant amount (Arithmetic) or by a
// constant percentage (Geometric), rather than staying fixed like an
// ordinary annuity.
//
// Basic usage example:
//
//	period, _ := compoundinterest.NewPeriod(decimal.MustFromFloat64(10), compoundinterest.Annually)
//	rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.1), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
//
//	// A series starting at $1000 and growing by $100 every year.
//	arithmetic, _ := NewArithmetic(
//	    money.MustMoneyFromFloat64(1000, money.USD),
//	    money.MustMoneyFromFloat64(100, money.USD),
//	    period, rate)
//	present, _ := arithmetic.Present()
//
//	// A series starting at $1000 and growing 8% every year.
//	geometric, _ := NewGeometric(
//	    money.MustMoneyFromFloat64(1000, money.USD),
//	    decimal.MustFromFloat64(0.08),
//	    period, rate)
//	future, _ := geometric.Future()
package gradients

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// base holds the period/rate configuration shared by Arithmetic and
// Geometric. It wraps compoundinterest.CompoundInterest purely to reuse
// GetEqualsRateInterestPeriods, which reconciles the rate's compounding
// frequency with the period's before every calculation.
type base struct {
	compoundInterest compoundinterest.CompoundInterest
}

func newBase(currency money.Currency, period compoundinterest.Period, rateInterest compoundinterest.RateInterest) (base, error) {
	zero := money.MustMoneyFromFloat64(0, currency)

	ci, err := compoundinterest.New(zero, zero, rateInterest, period)
	if err != nil {
		return base{}, err
	}

	return base{compoundInterest: ci}, nil
}

// periodsAndRate returns the number of periods and the periodic rate used
// by every gradient calculation.
func (b base) periodsAndRate() (decimal.Decimal, decimal.Decimal, error) {
	return b.compoundInterest.GetEqualsRateInterestPeriods()
}
