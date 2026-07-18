// Package bonds prices fixed-coupon bonds and derives their yield and interest
// rate risk measures: clean price from yield, yield to maturity from price,
// Macaulay and modified duration, convexity, and accrued interest.
//
// A bond is described on a per-period basis: a face (par) value, an annual
// coupon rate, a coupon frequency (coupons per year), and a number of coupon
// periods to maturity. Given an annual yield to maturity the package prices the
// bond; given a price it solves for the yield. All math runs on the decimal
// engine; monetary amounts carry their currency via the money package.
//
// Example — a 5-year 5% semiannual bond priced at a 6% yield:
//
//	price := bonds.NewBond().
//	    Face(1000, money.USD).
//	    CouponRate(0.05).
//	    Frequency(2).
//	    Periods(10).
//	    Yield(0.06).
//	    MustPrice()
//	// price ≈ 957.35 (below par, since the yield exceeds the coupon)
package bonds

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

var (
	// ErrInvalidFrequency is returned when the coupon frequency is not at
	// least one coupon per year.
	ErrInvalidFrequency = errors.New("bonds: frequency must be at least 1")

	// ErrInvalidPeriods is returned when the number of coupon periods to
	// maturity is not positive.
	ErrInvalidPeriods = errors.New("bonds: periods must be positive")

	// ErrInvalidYield is returned when a yield makes the per-period discount
	// factor (1 + yield/frequency) zero or negative.
	ErrInvalidYield = errors.New("bonds: yield must be greater than -frequency")

	// ErrNonPositivePrice is returned by YTM when the target price is not
	// positive.
	ErrNonPositivePrice = errors.New("bonds: price must be positive")

	// ErrNoConvergence is returned by YTM when no yield reproducing the price
	// can be bracketed.
	ErrNoConvergence = errors.New("bonds: yield did not converge")
)

// Config is a fluent builder describing a fixed-coupon bond. Create one with
// NewBond, set its terms, then call Price, YTM, or a risk measure.
type Config struct {
	face       money.Money
	couponRate decimal.Decimal
	freq       int
	periods    int
	yield      decimal.Decimal
	price      decimal.Decimal
}

// NewBond returns a Config defaulting to semiannual coupons (frequency 2) with
// every numeric term set to zero.
func NewBond() Config {
	return Config{
		couponRate: decimal.Zero,
		freq:       2,
		yield:      decimal.Zero,
		price:      decimal.Zero,
	}
}

// Face sets the face (par) value and its currency.
func (b Config) Face(amount float64, currency money.Currency) Config {
	b.face = money.MustMoneyFromFloat64(amount, currency)
	return b
}

// CouponRate sets the annual coupon rate as a fraction (e.g. 0.05 for 5%).
func (b Config) CouponRate(rate float64) Config {
	b.couponRate = decimal.MustFromFloat64(rate)
	return b
}

// Frequency sets the number of coupons paid per year (e.g. 2 for semiannual).
func (b Config) Frequency(f int) Config {
	b.freq = f
	return b
}

// Periods sets the number of coupon periods remaining to maturity.
func (b Config) Periods(n int) Config {
	b.periods = n
	return b
}

// Yield sets the annual yield to maturity as a fraction (e.g. 0.06 for 6%),
// used by Price and the risk measures.
func (b Config) Yield(y float64) Config {
	b.yield = decimal.MustFromFloat64(y)
	return b
}

// MarketPrice sets the observed clean price, used by YTM to solve for the
// yield.
func (b Config) MarketPrice(p float64) Config {
	b.price = decimal.MustFromFloat64(p)
	return b
}

// CouponPayment returns the cash coupon paid each period, face × couponRate /
// frequency, in the face value's currency.
//
// It returns ErrInvalidFrequency if the frequency is not at least 1.
func (b Config) CouponPayment() (money.Money, error) {
	if b.freq < 1 {
		return money.Money{}, ErrInvalidFrequency
	}

	coupon, err := b.face.ToDecimal().Mul(b.couponRate).Div(decimal.MustFromInt64(int64(b.freq), 0))
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(coupon, b.face.Currency()), nil
}
