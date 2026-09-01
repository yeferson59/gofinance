package money

import (
	"bytes"
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/yeferson59/gofinance/v2/decimal"
)

var MoneyZero = Money{decimal.Zero, USD}
var MoneyOne = Money{decimal.One, USD}

type Money struct {
	value    decimal.Decimal
	currency Currency
}

func New(value int64, precision uint8, currency Currency) (Money, error) {
	parsedValue, err := decimal.NewFromInt64(value, precision)
	if err != nil {
		return Money{}, err
	}

	return Money{parsedValue, currency}, nil
}

func NewMoneyFromFloat64(f float64, currency Currency) (Money, error) {
	parsedValue, err := decimal.NewFromFloat64(f)
	if err != nil {
		return Money{}, err
	}

	return Money{parsedValue, currency}, nil
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

	return Money{parsedValue, currency}, nil
}

func MustMoneyFromString(s string, currency Currency) Money {
	m, err := NewMoneyFromString(s, currency)
	if err != nil {
		panic(err)
	}

	return m
}

// NewFromDecimal attaches a currency to d, turning it into a Money value.
func NewFromDecimal(d decimal.Decimal, currency Currency) Money {
	return Money{d, currency}
}

func (m Money) GetDecimal() decimal.Decimal {
	return m.value
}

func (m Money) GetCurrency() Currency {
	return m.currency
}

func (m *Money) SetCurrency(c Currency) {
	if m == nil {
		return
	}

	if !c.Valid() {
		m.currency = USD

		return
	}

	m.currency = c
}

// TryAdd returns the sum of m and other, keeping the currency. It returns
// ErrCurrencyMismatch if the operands' currencies differ, or the decimal
// engine's error on overflow.
func (m Money) TryAdd(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	v, err := m.value.TryAdd(other.value)
	if err != nil {
		return Money{}, err
	}

	return Money{v, m.currency}, nil
}

// Add returns the sum of m and other. Like the decimal engine's Add, it
// panics instead of returning an error — on a currency mismatch or on
// overflow. Use TryAdd for the error-returning form.
func (m Money) Add(other Money) Money {
	v, err := m.TryAdd(other)
	if err != nil {
		panic(err)
	}

	return v
}

// TrySub returns the difference of m and other, keeping the currency. It
// returns ErrCurrencyMismatch if the operands' currencies differ, or the
// decimal engine's error on overflow.
func (m Money) TrySub(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	v, err := m.value.TrySub(other.value)
	if err != nil {
		return Money{}, err
	}

	return Money{v, m.currency}, nil
}

// Sub returns the difference of m and other. Like the decimal engine's Sub,
// it panics instead of returning an error — on a currency mismatch or on
// overflow. Use TrySub for the error-returning form.
func (m Money) Sub(other Money) Money {
	v, err := m.TrySub(other)
	if err != nil {
		panic(err)
	}

	return v
}

func (m Money) RoundBank(prec uint8) Money {
	return Money{m.value.RoundBank(prec), m.currency}
}

func (m Money) RoundBankString(prec uint8) string {
	return m.value.RoundBank(prec).StringFixed(prec)
}

func (m Money) RoundAway(prec uint8) Money {
	return Money{m.value.RoundAway(prec), m.currency}
}

func (m Money) Trunc(prec uint8) Money {
	return Money{m.value.Trunc(prec), m.currency}
}

func (m Money) Abs() Money {
	return Money{m.value.Abs(), m.currency}
}

func (m Money) Neg() Money {
	return Money{m.value.Neg(), m.currency}
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

	return Money{m.value.Mul(factor), m.currency}
}

// MulDecimal multiplies m by a currency-less factor, such as a rate or a
// growth factor, keeping m's currency. Prefer this over Mul when the factor
// is not itself an amount of money.
func (m Money) MulDecimal(d decimal.Decimal) Money {
	return Money{m.value.Mul(d), m.currency}
}

// DivDecimal divides m by a currency-less divisor, such as a rate or a
// number of periods, keeping m's currency.
func (m Money) DivDecimal(d decimal.Decimal) (Money, error) {
	v, err := m.value.Div(d)
	if err != nil {
		return Money{}, err
	}

	return Money{v, m.currency}, nil
}

// MustDivDecimal is like DivDecimal but panics on error.
func (m Money) MustDivDecimal(d decimal.Decimal) Money {
	v, err := m.DivDecimal(d)
	if err != nil {
		panic(err)
	}

	return v
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

	return Money{v, m.currency}, nil
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

func (m Money) InexactFloat64() float64 {
	return m.value.InexactFloat64()
}

func (m Money) Cmp(other Money) int {
	return m.value.Cmp(other.value)
}

func (m Money) Floor() Money {
	return Money{m.value.Floor(), m.currency}
}

func (m Money) Ceil() Money {
	return Money{m.value.Ceil(), m.currency}
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
	Value    jsontext.Value `json:"value"`
	Currency string         `json:"currency"`
}

func (m Money) MarshalJSON() ([]byte, error) {
	isoCode, err := m.currency.GetCurrencyISOCode()
	if err != nil {
		return nil, err
	}

	return json.Marshal(moneyJSON{
		Value:    m.value.String(),
		Currency: string(isoCode[:]),
	})
}

func (m *Money) UnmarshalJSON(data []byte) error {
	v, err := jsontext.NewDecoder(bytes.NewBuffer(data)).ReadValue()
	if err != nil {
		return err
	}

	if v.Kind() == jsontext.KindNumber {
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

	t, err := jsontext.NewDecoder(bytes.NewBuffer(mj.Value)).ReadToken()
	if err != nil {
		return err
	}

	var dec decimal.Decimal
	if err := dec.UnmarshalJSON([]byte(t.String())); err != nil {
		return err
	}

	m.value = dec

	if mj.Currency == "" {
		m.currency = USD

		return nil
	}

	currency, err := GetCurrencyFromISOCode(mj.Currency)
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
	if err := m.value.Scan(src); err != nil {
		return err
	}

	m.SetCurrency(USD)

	return nil
}
