package money

import "errors"

// ErrInvalidExchangeRate is returned by Convert when given a rate that is
// zero or negative.
var ErrInvalidExchangeRate = errors.New("money: exchange rate must be positive")

// Convert converts m into the given target currency using rate, the
// number of units of to that one unit of m's currency is worth
// (e.g. amountIn(to) = amountIn(m.Currency()) * rate). The result is
// rounded (half to even) to the target currency's standard precision.
//
//	usd := money.MustMoneyFromFloat64(100, money.USD)
//	eur, err := usd.Convert(money.EUR, money.MustFromFloat64(0.92))
func (m Money) Convert(to Currency, rate Decimal) (Money, error) {
	if !rate.IsPos() {
		return Money{}, ErrInvalidExchangeRate
	}

	v, err := m.value.TryMul(rate)
	if err != nil {
		return Money{}, err
	}

	prec, err := to.GetCurrencyPrecisionCode()
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    v.RoundBank(prec),
		currency: to,
	}, nil
}

// MustConvert is like Convert but panics on error.
func (m Money) MustConvert(to Currency, rate Decimal) Money {
	converted, err := m.Convert(to, rate)
	if err != nil {
		panic(err)
	}

	return converted
}

// ConvertFloat64 is like Convert but takes rate as a float64, for
// convenience when the rate isn't already a Decimal.
func (m Money) ConvertFloat64(to Currency, rate float64) (Money, error) {
	rateDec, err := NewFromFloat64(rate)
	if err != nil {
		return Money{}, err
	}

	return m.Convert(to, rateDec)
}

// MustConvertFloat64 is like ConvertFloat64 but panics on error.
func (m Money) MustConvertFloat64(to Currency, rate float64) Money {
	converted, err := m.ConvertFloat64(to, rate)
	if err != nil {
		panic(err)
	}

	return converted
}
