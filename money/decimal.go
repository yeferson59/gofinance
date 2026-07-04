package money

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/quagmt/udecimal"
)

var Zero = Decimal{value: udecimal.Zero}
var One = Decimal{value: udecimal.One}

type Decimal struct {
	value udecimal.Decimal
}

func NewFromFloat64(f float64) (Decimal, error) {
	decimal, err := udecimal.NewFromFloat64(f)

	return Decimal{decimal}, err
}

func NewFromInt64(coef int64, prec uint8) (Decimal, error) {
	decimal, err := udecimal.NewFromInt64(coef, prec)

	return Decimal{decimal}, err
}

func NewFromUint64(coef uint64, prec uint8) (Decimal, error) {
	decimal, err := udecimal.NewFromUint64(coef, prec)

	return Decimal{decimal}, err
}

func NewFromHiLo(neg bool, hi, lo uint64, prec uint8) (Decimal, error) {
	decimal, err := udecimal.NewFromHiLo(neg, hi, lo, prec)

	return Decimal{decimal}, err
}

func MustFromFloat64(f float64) Decimal {
	return Decimal{udecimal.MustFromFloat64(f)}
}

func MustFromInt64(coef int64, prec uint8) Decimal {
	return Decimal{udecimal.MustFromInt64(coef, prec)}
}

func MustFromUint64(coef uint64, prec uint8) Decimal {
	return Decimal{udecimal.MustFromUint64(coef, prec)}
}

func MustFromHiLo(neg bool, hi, lo uint64, prec uint8) Decimal {
	d, err := udecimal.NewFromHiLo(neg, hi, lo, prec)
	if err != nil {
		log.Fatal("invalid value from hilo")
	}

	return Decimal{d}
}

func NewFromString(s string) (Decimal, error) {
	decimal, err := udecimal.Parse(s)

	return Decimal{decimal}, err
}

func MustFromString(s string) Decimal {
	return Decimal{udecimal.MustParse(s)}
}

func (d Decimal) ToMoney(currency ...Currency) Money {
	if len(currency) != 0 {
		return Money{value: d.value, currency: currency[0]}
	}

	return Money{value: d.value}
}

func NewDecimalFromUDecimal(d udecimal.Decimal) Decimal {
	return Decimal{d}
}

func (d Decimal) Add(other Decimal) Decimal {
	return Decimal{d.value.Add(other.value)}
}

func (d Decimal) Sub(other Decimal) Decimal {
	return Decimal{d.value.Sub(other.value)}
}

func (d Decimal) Mul(other Decimal) Decimal {
	return Decimal{d.value.Mul(other.value)}
}

func (d Decimal) Div(other Decimal) (Decimal, error) {
	div, err := d.value.Div(other.value)

	return Decimal{div}, err
}

func (d Decimal) MustDiv(other Decimal) Decimal {
	return Decimal{d.value.MustDiv(other.value)}
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
	return Decimal{d.value.RoundAwayFromZero(prec)}
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

func parseDecimalJSON(data []byte) (udecimal.Decimal, error) {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return udecimal.Parse(s)
	}

	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		return udecimal.Parse(num.String())
	}

	return udecimal.Decimal{}, errors.New("invalid decimal JSON")
}
