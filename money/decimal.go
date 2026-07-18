package money

import (
	"github.com/yeferson59/gofinance/decimal"
)

// Decimal is an alias of decimal.Decimal, kept so that code written against
// earlier versions of this package (and the finance packages' previous
// signatures) keeps compiling. New code should import and use
// decimal.Decimal directly for rates, factors, and other currency-less
// quantities, and reserve this package for Money and Currency.
type Decimal = decimal.Decimal

// Zero is the decimal value 0.
//
// Deprecated: use decimal.Zero.
var Zero = decimal.Zero

// One is the decimal value 1.
//
// Deprecated: use decimal.One.
var One = decimal.One

// FromDecimal attaches a currency to d, turning it into a Money value.
// It replaces the former Decimal.ToMoney method.
func FromDecimal(d Decimal, currency Currency) Money {
	return Money{value: d, currency: currency}
}

// Deprecated: use decimal.NewFromFloat64.
func NewFromFloat64(f float64) (Decimal, error) { return decimal.NewFromFloat64(f) }

// Deprecated: use decimal.NewFromInt64.
func NewFromInt64(coef int64, prec uint8) (Decimal, error) { return decimal.NewFromInt64(coef, prec) }

// Deprecated: use decimal.NewFromUint64.
func NewFromUint64(coef uint64, prec uint8) (Decimal, error) {
	return decimal.NewFromUint64(coef, prec)
}

// Deprecated: use decimal.NewFromHiLo.
func NewFromHiLo(neg bool, hi, lo uint64, prec uint8) (Decimal, error) {
	return decimal.NewFromHiLo(neg, hi, lo, prec)
}

// Deprecated: use decimal.NewFromString.
func NewFromString(s string) (Decimal, error) { return decimal.NewFromString(s) }

// Deprecated: use decimal.MustFromFloat64.
func MustFromFloat64(f float64) Decimal { return decimal.MustFromFloat64(f) }

// Deprecated: use decimal.MustFromInt64.
func MustFromInt64(coef int64, prec uint8) Decimal { return decimal.MustFromInt64(coef, prec) }

// Deprecated: use decimal.MustFromUint64.
func MustFromUint64(coef uint64, prec uint8) Decimal { return decimal.MustFromUint64(coef, prec) }

// Deprecated: use decimal.MustFromHiLo.
func MustFromHiLo(neg bool, hi, lo uint64, prec uint8) Decimal {
	return decimal.MustFromHiLo(neg, hi, lo, prec)
}

// Deprecated: use decimal.MustFromString.
func MustFromString(s string) Decimal { return decimal.MustFromString(s) }
