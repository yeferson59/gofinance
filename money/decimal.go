package money

import (
	"github.com/yeferson59/gofinance/v2/decimal"
)

// FromDecimal attaches a currency to d, turning it into a Money value.
// It replaces the former Decimal.ToMoney method.
func FromDecimal(d decimal.Decimal, currency Currency) Money {
	return Money{d, currency}
}
