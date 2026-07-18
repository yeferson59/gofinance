package returns

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// RealValue converts a nominal amount into its real (inflation-adjusted) value
// after the given number of periods, i.e. its purchasing power in today's
// money:
//
//	real = nominal / (1 + inflation)^periods
//
// inflation is the per-period inflation rate as a fraction and must be greater
// than −1. The result carries the nominal amount's currency.
//
// It returns ErrInvalidInflationRate if 1+inflation is not positive.
func RealValue(nominal money.Money, inflation, periods decimal.Decimal) (money.Money, error) {
	factor, err := priceLevelFactor(inflation, periods)
	if err != nil {
		return money.Money{}, err
	}

	realAmount, err := nominal.ToDecimal().Div(factor)
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(realAmount, nominal.Currency()), nil
}

// MustRealValue is like RealValue but panics on error.
func MustRealValue(nominal money.Money, inflation, periods decimal.Decimal) money.Money {
	m, err := RealValue(nominal, inflation, periods)
	if err != nil {
		panic(err)
	}

	return m
}

// NominalValue converts a real amount into the nominal amount needed after the
// given number of periods to preserve the same purchasing power:
//
//	nominal = real × (1 + inflation)^periods
//
// inflation is the per-period inflation rate as a fraction and must be greater
// than −1. The result carries the real amount's currency.
//
// It returns ErrInvalidInflationRate if 1+inflation is not positive.
func NominalValue(realAmount money.Money, inflation, periods decimal.Decimal) (money.Money, error) {
	factor, err := priceLevelFactor(inflation, periods)
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(realAmount.ToDecimal().Mul(factor), realAmount.Currency()), nil
}

// MustNominalValue is like NominalValue but panics on error.
func MustNominalValue(realAmount money.Money, inflation, periods decimal.Decimal) money.Money {
	m, err := NominalValue(realAmount, inflation, periods)
	if err != nil {
		panic(err)
	}

	return m
}

// RealRate returns the inflation-adjusted (real) rate of return from a nominal
// rate and an inflation rate, via the Fisher equation:
//
//	real = (1 + nominal) / (1 + inflation) − 1
//
// Both rates are fractions per the same period; inflation must be greater than
// −1. The result is a decimal.Decimal fraction.
//
// It returns ErrInvalidInflationRate if 1+inflation is not positive.
func RealRate(nominalRate, inflation decimal.Decimal) (decimal.Decimal, error) {
	onePlusInflation := decimal.One.Add(inflation)
	if !onePlusInflation.IsPos() {
		return decimal.Decimal{}, ErrInvalidInflationRate
	}

	ratio, err := decimal.One.Add(nominalRate).Div(onePlusInflation)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return ratio.Sub(decimal.One), nil
}

// MustRealRate is like RealRate but panics on error.
func MustRealRate(nominalRate, inflation decimal.Decimal) decimal.Decimal {
	d, err := RealRate(nominalRate, inflation)
	if err != nil {
		panic(err)
	}

	return d
}

// priceLevelFactor returns (1 + inflation)^periods, the cumulative price-level
// growth used to convert between nominal and real values. It returns
// ErrInvalidInflationRate if 1+inflation is not positive.
func priceLevelFactor(inflation, periods decimal.Decimal) (decimal.Decimal, error) {
	base := decimal.One.Add(inflation)
	if !base.IsPos() {
		return decimal.Decimal{}, ErrInvalidInflationRate
	}

	return base.Pow(periods)
}
