package gradients

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// Geometric represents a series of n periodic payments that grow (or,
// with a negative GrowthRate, shrink) by a constant percentage each
// period: FirstPayment, FirstPayment×(1+GrowthRate),
// FirstPayment×(1+GrowthRate)², ..., FirstPayment×(1+GrowthRate)ⁿ⁻¹. Each
// payment is assumed to fall at the end of its period.
type Geometric struct {
	firstPayment money.Money
	growthRate   decimal.Decimal
	base
}

// NewGeometric creates a Geometric gradient series.
//
// Parameters:
//   - firstPayment: the payment made at the end of the first period
//   - growthRate: the constant percentage each subsequent payment grows by, as a decimal (e.g. 0.08 for 8%; negative for a shrinking series)
//   - period: the number of periods in the series
//   - rateInterest: the periodic interest rate
//
// Example:
//
//	period, _ := compoundinterest.NewPeriod(decimal.MustFromFloat64(5), compoundinterest.Annually)
//	rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.1), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
//	series, _ := NewGeometric(
//	    money.MustMoneyFromFloat64(1000, money.USD),
//	    decimal.MustFromFloat64(0.08),
//	    period, rate)
func NewGeometric(firstPayment money.Money, growthRate decimal.Decimal, period compoundinterest.Period, rateInterest compoundinterest.RateInterest) (Geometric, error) {
	b, err := newBase(firstPayment.GetCurrency(), period, rateInterest)
	if err != nil {
		return Geometric{}, err
	}

	return Geometric{firstPayment: firstPayment, growthRate: growthRate, base: b}, nil
}

// Present returns the present value of the geometric gradient series.
//
// When GrowthRate equals the periodic rate i, the ratio (1+g)/(1+i)
// collapses to 1, so the usual formula divides by zero; the series is
// discounted at a flat rate instead:
//
//	PV = A × n / (1+i)
//
// Otherwise:
//
//	PV = A × [1 - ((1+g)/(1+i))^n] / (i-g)
func (g Geometric) Present() (money.Money, error) {
	periods, rateInterest, err := g.periodsAndRate()
	if err != nil {
		return money.Money{}, err
	}

	if g.growthRate.Equal(rateInterest) {
		factor, err := periods.Div(decimal.One.Add(rateInterest))
		if err != nil {
			return money.Money{}, err
		}

		return g.firstPayment.MulDecimal(factor), nil
	}

	ratio, err := decimal.One.Add(g.growthRate).Div(decimal.One.Add(rateInterest))
	if err != nil {
		return money.Money{}, err
	}

	ratioPower, err := ratio.Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	factor, err := decimal.One.Sub(ratioPower).Div(rateInterest.Sub(g.growthRate))
	if err != nil {
		return money.Money{}, err
	}

	return g.firstPayment.MulDecimal(factor), nil
}

// Future returns the future value of the geometric gradient series.
//
// When GrowthRate equals the periodic rate i:
//
//	FV = A × n × (1+i)^(n-1)
//
// Otherwise:
//
//	FV = A × [(1+i)^n - (1+g)^n] / (i-g)
func (g Geometric) Future() (money.Money, error) {
	periods, rateInterest, err := g.periodsAndRate()
	if err != nil {
		return money.Money{}, err
	}

	if g.growthRate.Equal(rateInterest) {
		growthPower, err := decimal.One.Add(rateInterest).Pow(periods.Sub(decimal.One))
		if err != nil {
			return money.Money{}, err
		}

		return g.firstPayment.MulDecimal(periods.Mul(growthPower)), nil
	}

	rateGrowth, err := decimal.One.Add(rateInterest).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	gradientGrowth, err := decimal.One.Add(g.growthRate).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	factor, err := rateGrowth.Sub(gradientGrowth).Div(rateInterest.Sub(g.growthRate))
	if err != nil {
		return money.Money{}, err
	}

	return g.firstPayment.MulDecimal(factor), nil
}
