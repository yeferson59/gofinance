package investment

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// Perpetuity returns the present value of a level perpetuity — a fixed payment
// received every period forever — discounted at the given periodic rate:
//
//	PV = payment / rate
//
// rate must be strictly positive. The result carries the payment's currency.
//
// It returns ErrNonPositiveRate if rate is not positive.
func Perpetuity(payment money.Money, rate decimal.Decimal) (money.Money, error) {
	if !rate.IsPos() {
		return money.Money{}, ErrNonPositiveRate
	}

	value, err := payment.GetDecimal().Div(rate)
	if err != nil {
		return money.Money{}, err
	}

	return money.NewFromDecimal(value, payment.GetCurrency()), nil
}

// MustPerpetuity is like Perpetuity but panics on error.
func MustPerpetuity(payment money.Money, rate decimal.Decimal) money.Money {
	m, err := Perpetuity(payment, rate)
	if err != nil {
		panic(err)
	}

	return m
}

// GrowingPerpetuity returns the present value of a perpetuity whose payment
// grows at a constant rate each period (the Gordon growth model):
//
//	PV = payment / (rate − growth)
//
// payment is the cash flow one period from now. rate must exceed growth for the
// series to converge. The result carries the payment's currency.
//
// It returns ErrRateBelowGrowth if rate ≤ growth.
func GrowingPerpetuity(payment money.Money, rate, growth decimal.Decimal) (money.Money, error) {
	if rate.LessThanOrEqual(growth) {
		return money.Money{}, ErrRateBelowGrowth
	}

	value, err := payment.GetDecimal().Div(rate.Sub(growth))
	if err != nil {
		return money.Money{}, err
	}

	return money.NewFromDecimal(value, payment.GetCurrency()), nil
}

// MustGrowingPerpetuity is like GrowingPerpetuity but panics on error.
func MustGrowingPerpetuity(payment money.Money, rate, growth decimal.Decimal) money.Money {
	m, err := GrowingPerpetuity(payment, rate, growth)
	if err != nil {
		panic(err)
	}

	return m
}
