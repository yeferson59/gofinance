package money

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ErrCurrencyMismatch is returned by operations that require both operands
// to share the same currency, such as SafeAdd and SafeSub.
var ErrCurrencyMismatch = errors.New("money: currency mismatch")

var MoneyZero = Money{value: decZero, currency: USD}
var MoneyOne = Money{value: decOne, currency: USD}

type Money struct {
	value    decimal128
	currency Currency
}

func New(value int64, precision uint8, currency Currency) (Money, error) {
	parsedValue, err := decFromInt64(value, precision)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    parsedValue,
		currency: currency,
	}, nil
}

func NewMoneyFromFloat64(f float64, currency Currency) (Money, error) {
	parsedValue, err := decFromFloat64(f)
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
	parsedValue, err := parseDecimal(s)
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
	v, err := m.value.Add(other.value)
	if err != nil {
		panic(err)
	}

	return Money{
		value:    v,
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
	v, err := m.value.Mul(other.value)
	if err != nil {
		panic(err)
	}

	return Money{
		value:    v,
		currency: m.currency,
	}
}

func (m Money) Sub(other Money) Money {
	v, err := m.value.Sub(other.value)
	if err != nil {
		panic(err)
	}

	return Money{
		value:    v,
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
		dec, err := parseDecimalJSON(data)
		if err != nil {
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

	dec, err := parseDecimalJSON(mj.Value)
	if err != nil {
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
		dec decimal128
		err error
	)

	switch v := src.(type) {
	case []byte:
		dec, err = parseDecimal(string(v))
	case string:
		dec, err = parseDecimal(v)
	case uint64:
		dec, err = decFromUint64(v, 0)
	case int64:
		dec, err = decFromInt64(v, 0)
	case int:
		dec, err = decFromInt64(int64(v), 0)
	case int32:
		dec, err = decFromInt64(int64(v), 0)
	case float64:
		dec, err = decFromFloat64(v)
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
