package money

import (
	"github.com/yeferson59/gofinance/decimal"
)

// Decimal is money's general-purpose numeric type, a thin wrapper around
// decimal.Decimal. It exists alongside Money so that rates, factors, and
// other currency-less quantities (e.g. an exchange rate or a periodic
// interest rate) can be modeled distinctly from an amount of money.
type Decimal struct {
	value decimal.Decimal
}

var Zero = Decimal{decimal.Zero}
var One = Decimal{decimal.One}

func NewFromFloat64(f float64) (Decimal, error) {
	d, err := decimal.NewFromFloat64(f)

	return Decimal{d}, err
}

func NewFromInt64(coef int64, prec uint8) (Decimal, error) {
	d, err := decimal.NewFromInt64(coef, prec)

	return Decimal{d}, err
}

func NewFromUint64(coef uint64, prec uint8) (Decimal, error) {
	d, err := decimal.NewFromUint64(coef, prec)

	return Decimal{d}, err
}

func NewFromHiLo(neg bool, hi, lo uint64, prec uint8) (Decimal, error) {
	d, err := decimal.NewFromHiLo(neg, hi, lo, prec)

	return Decimal{d}, err
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
	d, err := decimal.NewFromString(s)

	return Decimal{d}, err
}

func MustFromString(s string) Decimal {
	d, err := NewFromString(s)
	if err != nil {
		panic(err)
	}

	return d
}

// ToMoney attaches a currency to d, turning it into a Money value. It
// defaults to USD when currency is omitted.
func (d Decimal) ToMoney(currency ...Currency) Money {
	if len(currency) != 0 {
		return Money{value: d.value, currency: currency[0]}
	}

	return Money{value: d.value, currency: USD}
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
	v, err := d.value.Div(other.value)

	return Decimal{v}, err
}

func (d Decimal) MustDiv(other Decimal) Decimal {
	v, err := d.Div(other)
	if err != nil {
		panic(err)
	}

	return v
}

func (d Decimal) Pow(exponent Decimal) (Decimal, error) {
	v, err := d.value.Pow(exponent.value)

	return Decimal{v}, err
}

func (d Decimal) MustPow(exponent Decimal) Decimal {
	v, err := d.Pow(exponent)
	if err != nil {
		panic(err)
	}

	return v
}

func (d Decimal) Sqrt() (Decimal, error) {
	v, err := d.value.Sqrt()

	return Decimal{v}, err
}

func (d Decimal) MustSqrt() Decimal {
	v, err := d.Sqrt()
	if err != nil {
		panic(err)
	}

	return v
}

func (d Decimal) Ln() (Decimal, error) {
	v, err := d.value.Ln()

	return Decimal{v}, err
}

func (d Decimal) MustLn() Decimal {
	v, err := d.Ln()
	if err != nil {
		panic(err)
	}

	return v
}

func (d Decimal) Log10() (Decimal, error) {
	v, err := d.value.Log10()

	return Decimal{v}, err
}

func (d Decimal) MustLog10() Decimal {
	v, err := d.Log10()
	if err != nil {
		panic(err)
	}

	return v
}

func (d Decimal) Log2() (Decimal, error) {
	v, err := d.value.Log2()

	return Decimal{v}, err
}

func (d Decimal) MustLog2() Decimal {
	v, err := d.Log2()
	if err != nil {
		panic(err)
	}

	return v
}

func (d Decimal) Log(base Decimal) (Decimal, error) {
	v, err := d.value.Log(base.value)

	return Decimal{v}, err
}

func (d Decimal) MustLog(base Decimal) Decimal {
	v, err := d.Log(base)
	if err != nil {
		panic(err)
	}

	return v
}

func (d Decimal) Mod(other Decimal) (Decimal, error) {
	v, err := d.value.Mod(other.value)

	return Decimal{v}, err
}

func (d Decimal) Div64(other uint64) (Decimal, error) {
	v, err := d.value.Div64(other)

	return Decimal{v}, err
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
	return d.value.Float64()
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
	return d.value.MarshalJSON()
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	return d.value.UnmarshalJSON(data)
}
