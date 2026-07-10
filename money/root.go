package money

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/yeferson59/gofinance/decimal"
)

// ErrCurrencyMismatch is returned by operations that require both operands
// to share the same currency, such as SafeAdd and SafeSub.
var ErrCurrencyMismatch = errors.New("money: currency mismatch")

var MoneyZero = Money{value: decimal.Zero, currency: USD}
var MoneyOne = Money{value: decimal.One, currency: USD}

type Money struct {
	value    decimal.Decimal
	currency Currency
}

func New(value int64, precision uint8, currency Currency) (Money, error) {
	parsedValue, err := decimal.NewFromInt64(value, precision)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    parsedValue,
		currency: currency,
	}, nil
}

func NewMoneyFromFloat64(f float64, currency Currency) (Money, error) {
	parsedValue, err := decimal.NewFromFloat64(f)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    parsedValue,
		currency: currency,
	}, nil
}

func MustMoneyFromFloat64(f float64, currency Currency) Money {
	m, err := NewMoneyFromFloat64(f, currency)
	if err != nil {
		panic(err)
	}

	return m
}

func NewMoneyFromString(s string, currency Currency) (Money, error) {
	parsedValue, err := decimal.NewFromString(s)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    parsedValue,
		currency: currency,
	}, nil
}

func MustMoneyFromString(s string, currency Currency) Money {
	m, err := NewMoneyFromString(s, currency)
	if err != nil {
		panic(err)
	}

	return m
}

func (m Money) ToDecimal() Decimal {
	return Decimal{m.value}
}

func (m Money) Add(other Money) Money {
	return Money{
		value:    m.value.Add(other.value),
		currency: m.currency,
	}
}

// SafeAdd returns the sum of m and other.
// It returns ErrCurrencyMismatch if the operands have different currencies.
func (m Money) SafeAdd(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	return m.Add(other), nil
}

func (m Money) Mul(other Money) Money {
	return Money{
		value:    m.value.Mul(other.value),
		currency: m.currency,
	}
}

func (m Money) Sub(other Money) Money {
	return Money{
		value:    m.value.Sub(other.value),
		currency: m.currency,
	}
}

// SafeSub returns the difference of m and other.
// It returns ErrCurrencyMismatch if the operands have different currencies.
func (m Money) SafeSub(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	return m.Sub(other), nil
}

func (m Money) Currency() Currency {
	return m.currency
}

func (m Money) RoundBank(prec uint8) Money {
	return Money{
		value:    m.value.RoundBank(prec),
		currency: m.currency,
	}
}

func (m Money) RoundBankString(prec uint8) string {
	return m.value.RoundBank(prec).StringFixed(prec)
}

func (m Money) RoundAway(prec uint8) Money {
	return Money{
		value:    m.value.RoundAway(prec),
		currency: m.currency,
	}
}

func (m Money) Trunc(prec uint8) Money {
	return Money{
		value:    m.value.Trunc(prec),
		currency: m.currency,
	}
}

func (m Money) Abs() Money {
	return Money{
		value:    m.value.Abs(),
		currency: m.currency,
	}
}

func (m Money) Neg() Money {
	return Money{
		value:    m.value.Neg(),
		currency: m.currency,
	}
}

func (m Money) IsZero() bool {
	return m.value.IsZero()
}

// IsPositive reports whether m is strictly greater than zero.
func (m Money) IsPositive() bool {
	return m.value.IsPos()
}

// IsNegative reports whether m is strictly less than zero.
func (m Money) IsNegative() bool {
	return m.value.IsNeg()
}

// MulInt64 multiplies m by a plain integer factor, such as a quantity
// (e.g. unit price * quantity).
func (m Money) MulInt64(n int64) Money {
	factor, err := decimal.NewFromInt64(n, 0)
	if err != nil {
		panic(err)
	}

	return Money{
		value:    m.value.Mul(factor),
		currency: m.currency,
	}
}

// DivInt64 divides m by a plain integer divisor.
func (m Money) DivInt64(n int64) (Money, error) {
	divisor, err := decimal.NewFromInt64(n, 0)
	if err != nil {
		return Money{}, err
	}

	v, err := m.value.Div(divisor)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    v,
		currency: m.currency,
	}, nil
}

// MustDivInt64 is like DivInt64 but panics on error.
func (m Money) MustDivInt64(n int64) Money {
	v, err := m.DivInt64(n)
	if err != nil {
		panic(err)
	}

	return v
}

// Min returns the smaller of m and other.
// It returns ErrCurrencyMismatch if the operands have different currencies.
func (m Money) Min(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	if m.value.Cmp(other.value) <= 0 {
		return m, nil
	}

	return other, nil
}

// Max returns the larger of m and other.
// It returns ErrCurrencyMismatch if the operands have different currencies.
func (m Money) Max(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	if m.value.Cmp(other.value) >= 0 {
		return m, nil
	}

	return other, nil
}

func (m Money) Div(other Money) (Money, error) {
	div, err := m.value.Div(other.value)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    div,
		currency: m.currency,
	}, nil
}

func (m Money) MustDiv(other Money) Money {
	div, err := m.Div(other)
	if err != nil {
		panic(err)
	}

	return div
}

func (m Money) InexactFloat64() float64 {
	return m.value.InexactFloat64()
}

func (m Money) Cmp(other Money) int {
	return m.value.Cmp(other.value)
}

func (m Money) Floor() Money {
	return Money{
		value:    m.value.Floor(),
		currency: m.currency,
	}
}

func (m Money) Ceil() Money {
	return Money{
		value:    m.value.Ceil(),
		currency: m.currency,
	}
}

func (m Money) String() string {
	return m.value.String()
}

func (m Money) StringFixed(prec uint8) string {
	return m.value.StringFixed(prec)
}

func (m Money) GreaterThan(other Money) bool {
	return m.value.GreaterThan(other.value)
}

func (m Money) GreaterThanOrEqual(other Money) bool {
	return m.value.GreaterThanOrEqual(other.value)
}

func (m Money) LessThan(other Money) bool {
	return m.value.LessThan(other.value)
}

func (m Money) LessThanOrEqual(other Money) bool {
	return m.value.LessThanOrEqual(other.value)
}

func (m Money) Equal(other Money) bool {
	return m.value.Equal(other.value) && m.currency == other.currency
}

type moneyJSON struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type moneyJSONRaw struct {
	Value    json.RawMessage `json:"value"`
	Currency string          `json:"currency"`
}

func (m Money) MarshalJSON() ([]byte, error) {
	isoCode, err := m.currency.GetCurrencyISOCode()
	if err != nil {
		return nil, err
	}

	return json.Marshal(moneyJSON{
		Value:    m.value.String(),
		Currency: isoCode,
	})
}

func (m *Money) UnmarshalJSON(data []byte) error {

	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		var dec decimal.Decimal
		if err := dec.UnmarshalJSON(data); err != nil {
			return err
		}

		m.value = dec
		m.currency = USD

		return nil
	}

	var mj moneyJSONRaw
	if err := json.Unmarshal(data, &mj); err != nil {
		return err
	}

	var dec decimal.Decimal
	if err := dec.UnmarshalJSON(mj.Value); err != nil {
		return err
	}

	m.value = dec

	if mj.Currency == "" {
		m.currency = USD

		return nil
	}

	currency, err := CurrencyFromISOCode(mj.Currency)
	if err != nil {
		return err
	}

	m.currency = currency

	return nil
}

func (m Money) Value() (driver.Value, error) {
	return m.value.String(), nil
}

func (m *Money) Scan(src any) error {
	var (
		dec decimal.Decimal
		err error
	)

	switch v := src.(type) {
	case []byte:
		dec, err = decimal.NewFromString(string(v))
	case string:
		dec, err = decimal.NewFromString(v)
	case uint64:
		dec, err = decimal.NewFromUint64(v, 0)
	case int64:
		dec, err = decimal.NewFromInt64(v, 0)
	case int:
		dec, err = decimal.NewFromInt64(int64(v), 0)
	case int32:
		dec, err = decimal.NewFromInt64(int64(v), 0)
	case float64:
		dec, err = decimal.NewFromFloat64(v)
	case nil:
		err = fmt.Errorf("money: can't scan nil to Money")
	default:
		err = fmt.Errorf("money: can't scan %T to Money", src)
	}

	if err != nil {
		return err
	}

	m.value = dec

	return nil
}
