package gradients

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// Arithmetic represents a series of n periodic payments that increase (or,
// with a negative Gradient, decrease) by a constant amount each period:
// FirstPayment, FirstPayment+Gradient, FirstPayment+2×Gradient, ...,
// FirstPayment+(n-1)×Gradient. Each payment is assumed to fall at the end
// of its period.
type Arithmetic struct {
	firstPayment money.Money
	gradient     money.Money
	base
}

// NewArithmetic creates an Arithmetic gradient series.
//
// Parameters:
//   - firstPayment: the payment made at the end of the first period
//   - gradient: the constant amount added to each subsequent payment (negative for a decreasing series); must share firstPayment's currency
//   - period: the number of periods in the series
//   - rateInterest: the periodic interest rate
//
// Example:
//
//	period, _ := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
//	rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.1), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
//	series, _ := NewArithmetic(
//	    money.MustMoneyFromFloat64(1000, money.USD),
//	    money.MustMoneyFromFloat64(100, money.USD),
//	    period, rate)
func NewArithmetic(firstPayment, gradient money.Money, period compoundinterest.Period, rateInterest compoundinterest.RateInterest) (Arithmetic, error) {
	if firstPayment.Currency() != gradient.Currency() {
		return Arithmetic{}, money.ErrCurrencyMismatch
	}

	b, err := newBase(firstPayment.Currency(), period, rateInterest)
	if err != nil {
		return Arithmetic{}, err
	}

	return Arithmetic{firstPayment: firstPayment, gradient: gradient, base: b}, nil
}

// Present returns the present value of the arithmetic gradient series:
//
//	PV = A × [1-(1+i)^-n]/i + G/i × { [1-(1+i)^-n]/i - n×(1+i)^-n }
//
// where A is FirstPayment, G is Gradient, i is the periodic rate, and n is
// the number of periods.
func (g Arithmetic) Present() (money.Money, error) {
	periods, rateInterest, err := g.periodsAndRate()
	if err != nil {
		return money.Money{}, err
	}

	growthPower, err := decimal.One.Add(rateInterest).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	discountFactor, err := decimal.One.Div(growthPower)
	if err != nil {
		return money.Money{}, err
	}

	annuityFactor, err := decimal.One.Sub(discountFactor).Div(rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	gradientFactor, err := annuityFactor.Sub(periods.Mul(discountFactor)).Div(rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	return g.firstPayment.MulDecimal(annuityFactor).Add(g.gradient.MulDecimal(gradientFactor)), nil
}

// Future returns the future value of the arithmetic gradient series:
//
//	FV = A × [(1+i)^n-1]/i + G/i × { [(1+i)^n-1]/i - n }
func (g Arithmetic) Future() (money.Money, error) {
	periods, rateInterest, err := g.periodsAndRate()
	if err != nil {
		return money.Money{}, err
	}

	growthPower, err := decimal.One.Add(rateInterest).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	annuityFactor, err := growthPower.Sub(decimal.One).Div(rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	gradientFactor, err := annuityFactor.Sub(periods).Div(rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	return g.firstPayment.MulDecimal(annuityFactor).Add(g.gradient.MulDecimal(gradientFactor)), nil
}
