package money

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/quagmt/udecimal"
)

// ErrCurrencyMismatch is returned by operations that require both operands
// to share the same currency, such as SafeAdd and SafeSub.
var ErrCurrencyMismatch = errors.New("money: currency mismatch")

var MoneyZero = Money{value: udecimal.Zero, currency: USD}
var MoneyOne = Money{value: udecimal.One, currency: USD}

type Money struct {
	value    udecimal.Decimal
	currency Currency
}

func New(value int64, precision uint8, currency Currency) (Money, error) {
	parsedValue, err := udecimal.NewFromInt64(value, precision)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    parsedValue,
		currency: currency,
	}, nil
}

func NewMoneyFromFloat64(f float64, currency Currency) (Money, error) {
	parsedValue, err := udecimal.NewFromFloat64(f)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    parsedValue,
		currency: currency,
	}, nil
}

func MustMoneyFromFloat64(f float64, currency Currency) Money {
	return Money{
		value:    udecimal.MustFromFloat64(f),
		currency: currency,
	}
}

func NewMoneyFromString(s string, currency Currency) (Money, error) {
	parsedValue, err := udecimal.Parse(s)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    parsedValue,
		currency: currency,
	}, nil
}

func MustMoneyFromString(s string, currency Currency) Money {
	return Money{
		value:    udecimal.MustParse(s),
		currency: currency,
	}
}

func (m Money) ToDecimal() Decimal {
	return Decimal{m.value}
}

func NewMoneyFromUDecimal(d udecimal.Decimal, currency Currency) Money {
	return Money{
		value:    d,
		currency: currency,
	}
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
		value:    m.value.RoundAwayFromZero(prec),
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
	return Money{
		value:    m.value.MustDiv(other.value),
		currency: m.currency,
	}
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
	return m.value.Value()
}

func (m *Money) Scan(src any) error {
	return m.value.Scan(src)
}
