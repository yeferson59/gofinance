package money

import (
	"encoding/json"
	"math"
)

var Zero = Decimal{value: decZero}
var One = Decimal{value: decOne}

// Decimal is a fixed-point decimal number with up to 19 digits of
// precision after the decimal point, backed by a 128-bit coefficient.
// It performs arithmetic without heap allocations or external
// dependencies.
type Decimal struct {
	value decimal128
}

func NewFromFloat64(f float64) (Decimal, error) {
	decimal, err := decFromFloat64(f)

	return Decimal{decimal}, err
}

func NewFromInt64(coef int64, prec uint8) (Decimal, error) {
	decimal, err := decFromInt64(coef, prec)

	return Decimal{decimal}, err
}

func NewFromUint64(coef uint64, prec uint8) (Decimal, error) {
	decimal, err := decFromUint64(coef, prec)

	return Decimal{decimal}, err
}

func NewFromHiLo(neg bool, hi, lo uint64, prec uint8) (Decimal, error) {
	decimal, err := decFromHiLo(neg, hi, lo, prec)

	return Decimal{decimal}, err
}

func MustFromFloat64(f float64) Decimal {
	d, err := NewFromFloat64(f)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFromInt64(coef int64, prec uint8) Decimal {
	d, err := NewFromInt64(coef, prec)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFromUint64(coef uint64, prec uint8) Decimal {
	d, err := NewFromUint64(coef, prec)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFromHiLo(neg bool, hi, lo uint64, prec uint8) Decimal {
	d, err := NewFromHiLo(neg, hi, lo, prec)
	if err != nil {
		panic(err)
	}

	return d
}

func NewFromString(s string) (Decimal, error) {
	decimal, err := parseDecimal(s)

	return Decimal{decimal}, err
}

func MustFromString(s string) Decimal {
	d, err := NewFromString(s)
	if err != nil {
		panic(err)
	}

	return d
}

func (d Decimal) ToMoney(currency ...Currency) Money {
	if len(currency) != 0 {
		return Money{value: d.value, currency: currency[0]}
	}

	return Money{value: d.value}
}

func (d Decimal) Add(other Decimal) Decimal {
	v, err := d.value.Add(other.value)
	if err != nil {
		panic(err)
	}

	return Decimal{v}
}

func (d Decimal) Sub(other Decimal) Decimal {
	v, err := d.value.Sub(other.value)
	if err != nil {
		panic(err)
	}

	return Decimal{v}
}

func (d Decimal) Mul(other Decimal) Decimal {
	v, err := d.value.Mul(other.value)
	if err != nil {
		panic(err)
	}

	return Decimal{v}
}

func (d Decimal) Div(other Decimal) (Decimal, error) {
	div, err := d.value.Div(other.value)

	return Decimal{div}, err
}

func (d Decimal) MustDiv(other Decimal) Decimal {
	div, err := d.Div(other)
	if err != nil {
		panic(err)
	}

	return div
}

// Pow returns d raised to the power of exponent. decimal128 has no native
// exponentiation, so the computation is done in float64 via math.Pow and
// the result converted back to a Decimal. It returns an error if the
// result isn't a finite number (e.g. a negative base with a fractional
// exponent, which is undefined).
func (d Decimal) Pow(exponent Decimal) (Decimal, error) {
	result := math.Pow(d.InexactFloat64(), exponent.InexactFloat64())

	return NewFromFloat64(result)
}

// MustPow is like Pow but panics on error.
func (d Decimal) MustPow(exponent Decimal) Decimal {
	result, err := d.Pow(exponent)
	if err != nil {
		panic(err)
	}

	return result
}

func (d Decimal) Mod(other Decimal) (Decimal, error) {
	mod, err := d.value.Mod(other.value)

	return Decimal{mod}, err
}

func (d Decimal) Div64(other uint64) (Decimal, error) {
	div, err := d.value.Div64(other)

	return Decimal{div}, err
}

func (d Decimal) InexactFloat64() float64 {
	return d.value.InexactFloat64()
}

func (d Decimal) RoundBank(prec uint8) Decimal {
	return Decimal{d.value.RoundBank(prec)}
}

func (d Decimal) RoundAway(prec uint8) Decimal {
	return Decimal{d.value.RoundAway(prec)}
}

func (d Decimal) RoundHAZ(prec uint8) Decimal {
	return Decimal{d.value.RoundHAZ(prec)}
}

func (d Decimal) RoundHTZ(prec uint8) Decimal {
	return Decimal{d.value.RoundHTZ(prec)}
}

func (d Decimal) Trunc(prec uint8) Decimal {
	return Decimal{d.value.Trunc(prec)}
}

func (d Decimal) Abs() Decimal {
	return Decimal{d.value.Abs()}
}

func (d Decimal) Floor() Decimal {
	return Decimal{d.value.Floor()}
}

func (d Decimal) Ceil() Decimal {
	return Decimal{d.value.Ceil()}
}

func (d Decimal) Neg() Decimal {
	return Decimal{d.value.Neg()}
}

func (d Decimal) Sign() int {
	return d.value.Sign()
}

func (d Decimal) Cmp(other Decimal) int {
	return d.value.Cmp(other.value)
}

func (d Decimal) String() string {
	return d.value.String()
}

func (d Decimal) StringFixed(prec uint8) string {
	return d.value.StringFixed(prec)
}

func (d Decimal) Float64() (float64, error) {
	return d.value.InexactFloat64(), nil
}

func (d Decimal) Int64() (int64, error) {
	return d.value.Int64()
}

func (d Decimal) IsZero() bool {
	return d.value.IsZero()
}

func (d Decimal) IsNeg() bool {
	return d.value.IsNeg()
}

func (d Decimal) IsPos() bool {
	return d.value.IsPos()
}

func (d Decimal) GreaterThan(other Decimal) bool {
	return d.value.GreaterThan(other.value)
}

func (d Decimal) GreaterThanOrEqual(other Decimal) bool {
	return d.value.GreaterThanOrEqual(other.value)
}

func (d Decimal) LessThan(other Decimal) bool {
	return d.value.LessThan(other.value)
}

func (d Decimal) LessThanOrEqual(other Decimal) bool {
	return d.value.LessThanOrEqual(other.value)
}

func (d Decimal) Equal(other Decimal) bool {
	return d.value.Equal(other.value)
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.value.String())
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	dec, err := parseDecimalJSON(data)
	if err != nil {
		return err
	}

	d.value = dec

	return nil
}

func parseDecimalJSON(data []byte) (decimal128, error) {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return parseDecimal(s)
	}

	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		return parseDecimal(num.String())
	}

	return decimal128{}, ErrInvalidFormat
}
