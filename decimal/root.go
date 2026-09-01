package decimal

import (
	"bytes"
	"database/sql/driver"
	"encoding/json/jsontext"
	"errors"
	"fmt"
	"strings"
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

func NewFromString(s string) (Decimal, error) {
	decimal, err := parseDecimal(s)

	return Decimal{decimal}, err
}

// newFraction divides num by den, wrapping the result in a Decimal. It
// returns ErrDivideByZero if den is zero.
func newFraction(num, den decimal128) (Decimal, error) {
	result, err := num.Div(den)

	return Decimal{result}, err
}

// NewFractionFloat64 returns the Decimal value of numerator / denominator.
// It returns ErrDivideByZero if denominator is zero.
func NewFractionFloat64(numerator float64, denominator float64) (Decimal, error) {
	num, err := decFromFloat64(numerator)
	if err != nil {
		return Decimal{}, err
	}

	den, err := decFromFloat64(denominator)
	if err != nil {
		return Decimal{}, err
	}

	return newFraction(num, den)
}

// NewFractionInt64 returns the Decimal value of numerator / denominator,
// with each operand scaled by its own precision (e.g. numerator=150,
// numeratorPrec=2 represents 1.50). It returns ErrDivideByZero if
// denominator is zero.
func NewFractionInt64(numerator int64, numeratorPrec uint8, denominator int64, denominatorPrec uint8) (Decimal, error) {
	num, err := decFromInt64(numerator, numeratorPrec)
	if err != nil {
		return Decimal{}, err
	}

	den, err := decFromInt64(denominator, denominatorPrec)
	if err != nil {
		return Decimal{}, err
	}

	return newFraction(num, den)
}

// NewFractionString parses a "numerator/denominator" string (e.g. "1/3")
// and returns the resulting Decimal value. Leading/trailing whitespace
// around the whole string and around each operand is ignored. It returns
// ErrSymbolFraction if frac doesn't contain exactly one '/', and
// ErrDivideByZero if the denominator is zero.
func NewFractionString(frac string) (Decimal, error) {
	frac = strings.TrimSpace(frac)

	sep := "/"

	if strings.Count(frac, sep) != 1 {
		return Decimal{}, ErrSymbolFraction
	}

	split := strings.SplitN(frac, sep, 2)

	num, err := parseDecimal(strings.TrimSpace(split[0]))
	if err != nil {
		return Decimal{}, err
	}

	den, err := parseDecimal(strings.TrimSpace(split[1]))
	if err != nil {
		return Decimal{}, err
	}

	return newFraction(num, den)
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

func MustFromString(s string) Decimal {
	d, err := NewFromString(s)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFractionFloat64(numerator float64, denominator float64) Decimal {
	d, err := NewFractionFloat64(numerator, denominator)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFractionInt64(numerator int64, numeratorPrec uint8, denominator int64, denominatorPrec uint8) Decimal {
	d, err := NewFractionInt64(numerator, numeratorPrec, denominator, denominatorPrec)
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

// Sqrt returns the square root of d, computed directly with Newton's
// integer iteration on the exact 256-bit radicand, so the result is always
// the 19-fractional-digit value nearest to the true root (exact for
// perfect squares, e.g. Sqrt(2.25) = 1.5). It returns ErrSqrtNegative if
// d is negative. For other roots, use Pow (e.g. d.Pow(1/3) for cube roots).
func (d Decimal) Sqrt() (Decimal, error) {
	v, err := d.value.Sqrt()

	return Decimal{v}, err
}

// MustSqrt is like Sqrt but panics on error.
func (d Decimal) MustSqrt() Decimal {
	v, err := d.Sqrt()
	if err != nil {
		panic(err)
	}

	return v
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
	v, err := jsontext.NewDecoder(bytes.NewBufferString(d.value.String())).ReadValue()
	if err != nil {
		return nil, err
	}

	return v.MarshalJSON()
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	dec, err := parseDecimalJSON(data)
	if err != nil {
		return err
	}

	d.value = dec

	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}

func parseDecimalJSON(data []byte) (decimal128, error) {
	t, err := jsontext.NewDecoder(bytes.NewReader(data)).ReadToken()
	if err != nil {
		return decimal128{}, err
	}

	if t.Kind() != jsontext.KindNumber {
		return decimal128{}, ErrInvalidFormat
	}

	return parseDecimal(t.String())
}

func (d *Decimal) Scan(src any) error {
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
		err = errors.New("can't scan nil")
	default:
		err = fmt.Errorf("can't scan %T", src)
	}

	if err != nil {
		return err
	}

	d.value = dec

	return nil
}
