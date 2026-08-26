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
	if firstPayment.GetCurrency() != gradient.GetCurrency() {
		return Arithmetic{}, money.ErrCurrencyMismatch
	}

	b, err := newBase(firstPayment.GetCurrency(), period, rateInterest)
	if err != nil {
		return Arithmetic{}, err
	}

	return Arithmetic{firstPayment: firstPayment, gradient: gradient, base: b}, nil
}

// half is the exact decimal 0.5, used to halve the gradient's triangular sum.
var half = decimal.MustFromFloat64(0.5)

// zeroRateValue returns the series' value when the periodic rate is zero.
// With nothing to discount or compound, the present and future values both
// collapse to the plain sum of the payments:
//
//	A×n + G×(0+1+…+(n−1)) = A×n + G×n(n−1)/2
//
// Both general formulas divide by the rate, so this limit is returned
// directly instead of failing on a legitimate input.
func (g Arithmetic) zeroRateValue(periods decimal.Decimal) money.Money {
	gradientSteps := periods.Mul(periods.Sub(decimal.One)).Mul(half)

	return g.firstPayment.MulDecimal(periods).Add(g.gradient.MulDecimal(gradientSteps))
}

// Present returns the present value of the arithmetic gradient series:
//
//	PV = A × [1-(1+i)^-n]/i + G/i × { [1-(1+i)^-n]/i - n×(1+i)^-n }
//
// where A is FirstPayment, G is Gradient, i is the periodic rate, and n is
// the number of periods.
//
// At a zero rate this reduces to the sum of the payments, A×n + G×n(n−1)/2.
func (g Arithmetic) Present() (money.Money, error) {
	periods, rateInterest, err := g.periodsAndRate()
	if err != nil {
		return money.Money{}, err
	}

	if rateInterest.IsZero() {
		return g.zeroRateValue(periods), nil
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
//
// At a zero rate this reduces to the sum of the payments, A×n + G×n(n−1)/2,
// the same as Present: with no interest there is nothing to move the money
// through time.
func (g Arithmetic) Future() (money.Money, error) {
	periods, rateInterest, err := g.periodsAndRate()
	if err != nil {
		return money.Money{}, err
	}

	if rateInterest.IsZero() {
		return g.zeroRateValue(periods), nil
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
