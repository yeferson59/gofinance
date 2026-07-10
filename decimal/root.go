package decimal

import (
	"encoding/json"
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

// TryAdd is like Add but returns an error instead of panicking on overflow.
func (d Decimal) TryAdd(other Decimal) (Decimal, error) {
	v, err := d.value.Add(other.value)

	return Decimal{v}, err
}

func (d Decimal) Add(other Decimal) Decimal {
	v, err := d.TryAdd(other)
	if err != nil {
		panic(err)
	}

	return v
}

// TrySub is like Sub but returns an error instead of panicking on overflow.
func (d Decimal) TrySub(other Decimal) (Decimal, error) {
	v, err := d.value.Sub(other.value)

	return Decimal{v}, err
}

func (d Decimal) Sub(other Decimal) Decimal {
	v, err := d.TrySub(other)
	if err != nil {
		panic(err)
	}

	return v
}

// TryMul is like Mul but returns an error instead of panicking on overflow.
func (d Decimal) TryMul(other Decimal) (Decimal, error) {
	v, err := d.value.Mul(other.value)

	return Decimal{v}, err
}

func (d Decimal) Mul(other Decimal) Decimal {
	v, err := d.TryMul(other)
	if err != nil {
		panic(err)
	}

	return v
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

// Pow returns d raised to the power of exponent, computed natively on the
// decimal engine with 120-bit binary fixed-point internals, so the result
// is accurate to Decimal's full 19-digit precision (rounded half away from
// zero). Conventions: d^0 = 1 (including 0^0); 0^exponent = 0 for a
// positive exponent and ErrDivideByZero for a negative one; a negative
// base requires an integer exponent (the result's sign follows its
// parity), otherwise ErrPowNegBase is returned. Results beyond Decimal's
// range return ErrOverflow; results too small to represent round to zero.
func (d Decimal) Pow(exponent Decimal) (Decimal, error) {
	v, err := d.value.Pow(exponent.value)

	return Decimal{v}, err
}

// MustPow is like Pow but panics on error.
func (d Decimal) MustPow(exponent Decimal) Decimal {
	result, err := d.Pow(exponent)
	if err != nil {
		panic(err)
	}

	return result
}

// Ln returns the natural logarithm (base e) of d, accurate to Decimal's
// full 19-digit precision (rounded half away from zero). It returns
// ErrLogNonPositive if d is zero or negative, since the natural logarithm
// is undefined there.
func (d Decimal) Ln() (Decimal, error) {
	v, err := d.value.Ln()

	return Decimal{v}, err
}

// MustLn is like Ln but panics on error.
func (d Decimal) MustLn() Decimal {
	result, err := d.Ln()
	if err != nil {
		panic(err)
	}

	return result
}

// Log10 returns the base-10 logarithm of d, accurate to Decimal's full
// 19-digit precision. It returns ErrLogNonPositive if d isn't strictly
// positive.
func (d Decimal) Log10() (Decimal, error) {
	v, err := d.value.Log10()

	return Decimal{v}, err
}

// MustLog10 is like Log10 but panics on error.
func (d Decimal) MustLog10() Decimal {
	result, err := d.Log10()
	if err != nil {
		panic(err)
	}

	return result
}

// Log2 returns the base-2 logarithm of d, accurate to Decimal's full
// 19-digit precision. It returns ErrLogNonPositive if d isn't strictly
// positive.
func (d Decimal) Log2() (Decimal, error) {
	v, err := d.value.Log2()

	return Decimal{v}, err
}

// MustLog2 is like Log2 but panics on error.
func (d Decimal) MustLog2() Decimal {
	result, err := d.Log2()
	if err != nil {
		panic(err)
	}

	return result
}

// Log returns the logarithm of d in the given base, computed as
// Ln(d) / Ln(base) and accurate to Decimal's full 19-digit precision. It
// returns ErrLogNonPositive if d or base aren't strictly positive, and
// ErrDivideByZero if base equals 1 (which makes the logarithm undefined).
func (d Decimal) Log(base Decimal) (Decimal, error) {
	v, err := d.value.Log(base.value)

	return Decimal{v}, err
}

// MustLog is like Log but panics on error.
func (d Decimal) MustLog(base Decimal) Decimal {
	result, err := d.Log(base)
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
