package decimal

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// maxScale is the maximum number of digits decimal128 can keep after the
// decimal point. It's chosen so that a power of ten up to 10^maxScale still
// fits in a uint64 (10^19 < 2^64 <= 10^20), which keeps scale-alignment
// (Add/Sub/Cmp/rounding) on the fast, allocation-free uint64 division path.
const maxScale uint8 = 19

var (
	ErrOverflow        = errors.New("money: numeric overflow")
	ErrDivideByZero    = errors.New("money: division by zero")
	ErrEmptyString     = errors.New("money: can't parse empty string")
	ErrInvalidFormat   = errors.New("money: invalid decimal format")
	ErrPrecOutOfRange  = errors.New("money: precision out of range, maximum is 19 digits after the decimal point")
	ErrIntPartOverflow = errors.New("money: integer part is too large to fit in int64")
)

// decimal128 is a fixed-point decimal number backed by a 128-bit unsigned
// coefficient, a sign, and a scale (number of digits after the decimal
// point). It's the engine behind both Decimal and Money.
//
// value = (-1)^neg * coef / 10^scale
//
// The zero value represents 0 and is always canonical: a zero coefficient
// always has neg == false and scale == 0.
type decimal128 struct {
	coef  u128
	neg   bool
	scale uint8
}

var (
	decZero = decimal128{}
	decOne  = decimal128{coef: u128One}
)

// newDec builds a canonical decimal128, collapsing every zero-coefficient
// value to the canonical zero regardless of the requested sign/scale.
func newDec(neg bool, coef u128, scale uint8) decimal128 {
	if coef.IsZero() {
		return decimal128{}
	}

	return decimal128{coef: coef, neg: neg, scale: scale}
}

// mulPow10 returns coef*10^exp. exp must be <= maxScale, which guarantees
// pow10[exp] fits in a uint64 and this stays on the fast multiplication
// path (no 256-bit intermediate needed).
func mulPow10(coef u128, exp uint8) (u128, bool) {
	if exp == 0 {
		return coef, true
	}

	return coef.Mul64(pow10[exp].lo)
}

func decFromUint64(coef uint64, scale uint8) (decimal128, error) {
	if scale > maxScale {
		return decimal128{}, ErrPrecOutOfRange
	}

	return newDec(false, u128FromU64(coef), scale), nil
}

func decFromInt64(coef int64, scale uint8) (decimal128, error) {
	if scale > maxScale {
		return decimal128{}, ErrPrecOutOfRange
	}

	var (
		neg bool
		mag uint64
	)

	switch {
	case coef < 0 && coef == math.MinInt64:
		neg = true
		mag = uint64(math.MaxInt64) + 1
	case coef < 0:
		neg = true
		mag = uint64(-coef)
	default:
		mag = uint64(coef)
	}

	return newDec(neg, u128FromU64(mag), scale), nil
}

func decFromHiLo(neg bool, hi, lo uint64, scale uint8) (decimal128, error) {
	if scale > maxScale {
		return decimal128{}, ErrPrecOutOfRange
	}

	return newDec(neg, u128{hi: hi, lo: lo}, scale), nil
}

func decFromFloat64(f float64) (decimal128, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return decimal128{}, ErrInvalidFormat
	}

	s := strconv.FormatFloat(f, 'f', -1, 64)

	return parseDecimal(s)
}

// parseDecimal parses a string of the form [+-]digits[.digits] into a
// decimal128. It doesn't support scientific notation.
func parseDecimal(s string) (decimal128, error) {
	if len(s) == 0 {
		return decimal128{}, ErrEmptyString
	}

	var i int
	var neg bool

	switch s[0] {
	case '-':
		neg = true
		i++
	case '+':
		i++
	}

	if i == len(s) || s[i] == '.' {
		return decimal128{}, ErrInvalidFormat
	}

	var (
		coef       u128
		seenDigit  bool
		seenDot    bool
		fracDigits int
		ok         bool
	)

	for ; i < len(s); i++ {
		c := s[i]
		if c == '.' {
			if seenDot {
				return decimal128{}, ErrInvalidFormat
			}

			seenDot = true

			continue
		}

		if c < '0' || c > '9' {
			return decimal128{}, ErrInvalidFormat
		}

		seenDigit = true
		if seenDot {

			fracDigits++

			if fracDigits > int(maxScale) {
				return decimal128{}, ErrPrecOutOfRange
			}
		}

		coef, ok = coef.Mul64(10)
		if !ok {
			return decimal128{}, ErrOverflow
		}

		coef, ok = coef.Add64(uint64(c - '0'))
		if !ok {
			return decimal128{}, ErrOverflow
		}
	}

	if !seenDigit || (seenDot && fracDigits == 0) {
		return decimal128{}, ErrInvalidFormat
	}

	return newDec(neg, coef, uint8(fracDigits)), nil
}

func (a decimal128) IsZero() bool {
	return a.coef.IsZero()
}

func (a decimal128) IsNeg() bool {
	return a.neg && !a.coef.IsZero()
}

func (a decimal128) IsPos() bool {
	return !a.neg && !a.coef.IsZero()
}

func (a decimal128) Sign() int {
	if a.coef.IsZero() {
		return 0
	}

	if a.neg {
		return -1
	}

	return 1
}

func (a decimal128) Neg() decimal128 {
	if a.coef.IsZero() {
		return a
	}

	return decimal128{coef: a.coef, neg: !a.neg, scale: a.scale}
}

func (a decimal128) Abs() decimal128 {
	return decimal128{coef: a.coef, neg: false, scale: a.scale}
}

// cmpMagnitude compares |a| and |b|, ignoring sign.
func (a decimal128) cmpMagnitude(b decimal128) int {
	switch {
	case a.scale == b.scale:
		return a.coef.Cmp(b.coef)
	case a.scale < b.scale:
		scaled, ok := mulPow10(a.coef, b.scale-a.scale)
		if !ok {
			// a, once aligned to b's scale, doesn't fit in 128 bits, so it
			// must be larger than any 128-bit b.coef.
			return 1
		}
		return scaled.Cmp(b.coef)
	default:
		scaled, ok := mulPow10(b.coef, a.scale-b.scale)
		if !ok {
			return -1
		}

		return a.coef.Cmp(scaled)
	}
}

func (a decimal128) Cmp(b decimal128) int {
	if a.neg && !b.neg {
		return -1
	}

	if !a.neg && b.neg {
		return 1
	}

	cmp := a.cmpMagnitude(b)
	if a.neg {
		return -cmp
	}

	return cmp
}

func (a decimal128) Equal(b decimal128) bool              { return a.Cmp(b) == 0 }
func (a decimal128) LessThan(b decimal128) bool           { return a.Cmp(b) < 0 }
func (a decimal128) LessThanOrEqual(b decimal128) bool    { return a.Cmp(b) <= 0 }
func (a decimal128) GreaterThan(b decimal128) bool        { return a.Cmp(b) > 0 }
func (a decimal128) GreaterThanOrEqual(b decimal128) bool { return a.Cmp(b) >= 0 }

// Add returns a+b.
func (a decimal128) Add(b decimal128) (decimal128, error) {
	scale := a.scale
	ac, bc := a.coef, b.coef
	var ok bool

	switch {
	case a.scale < b.scale:
		scale = b.scale
		if ac, ok = mulPow10(ac, b.scale-a.scale); !ok {
			return decimal128{}, ErrOverflow
		}
	case a.scale > b.scale:
		if bc, ok = mulPow10(bc, a.scale-b.scale); !ok {
			return decimal128{}, ErrOverflow
		}
	}

	if a.neg == b.neg {
		sum, ok := ac.Add(bc)
		if !ok {
			return decimal128{}, ErrOverflow
		}

		return newDec(a.neg, sum, scale), nil
	}

	switch ac.Cmp(bc) {
	case 0:
		return decimal128{}, nil
	case 1:
		d, _ := ac.Sub(bc)

		return newDec(a.neg, d, scale), nil
	default:
		d, _ := bc.Sub(ac)

		return newDec(b.neg, d, scale), nil
	}
}

// Sub returns a-b.
func (a decimal128) Sub(b decimal128) (decimal128, error) {
	return a.Add(b.Neg())
}

// Mul returns a*b, truncated to at most maxScale digits after the decimal
// point (no rounding is applied when truncating excess precision).
func (a decimal128) Mul(b decimal128) (decimal128, error) {
	if a.coef.IsZero() || b.coef.IsZero() {
		return decimal128{}, nil
	}

	neg := a.neg != b.neg
	totalScale := int(a.scale) + int(b.scale)
	prod := a.coef.MulFull(b.coef)

	if totalScale <= int(maxScale) {
		if !prod.carry.IsZero() {
			return decimal128{}, ErrOverflow
		}
		//nolint:gosec // totalScale <= maxScale (19)
		return newDec(neg, u128{hi: prod.hi, lo: prod.lo}, uint8(totalScale)), nil
	}

	excess := totalScale - int(maxScale)
	q, _, ok := prod.QuoRem128(pow10[excess])
	if !ok {
		return decimal128{}, ErrOverflow
	}

	return newDec(neg, q, maxScale), nil
}

// Div returns a/b with the result always expressed to maxScale digits
// after the decimal point (truncated, not rounded).
func (a decimal128) Div(b decimal128) (decimal128, error) {
	if b.coef.IsZero() {
		return decimal128{}, ErrDivideByZero
	}

	if a.coef.IsZero() {
		return decimal128{}, nil
	}

	neg := a.neg != b.neg
	factor := int(maxScale) - int(a.scale) + int(b.scale)

	num := a.coef.MulFull(pow10[factor])
	q, _, ok := num.QuoRem128(b.coef)
	if !ok {
		return decimal128{}, ErrOverflow
	}

	return newDec(neg, q, maxScale), nil
}

// Div64 returns a/v with the result always expressed to maxScale digits
// after the decimal point (truncated, not rounded).
func (a decimal128) Div64(v uint64) (decimal128, error) {
	if v == 0 {
		return decimal128{}, ErrDivideByZero
	}

	if a.coef.IsZero() {
		return decimal128{}, nil
	}

	num := a.coef.MulFull(pow10[maxScale-a.scale])
	q, _, ok := num.QuoRem128(u128FromU64(v))
	if !ok {
		return decimal128{}, ErrOverflow
	}

	return newDec(a.neg, q, maxScale), nil
}

// QuoRem returns q = trunc(a/b) (scale 0) and r = a - q*b (same sign as a,
// scale = max(a.scale, b.scale)), similar to C's fmod.
func (a decimal128) QuoRem(b decimal128) (q, r decimal128, err error) {
	if b.coef.IsZero() {
		return decimal128{}, decimal128{}, ErrDivideByZero
	}

	if a.coef.IsZero() {
		return decimal128{}, decimal128{}, nil
	}

	factor := max(b.scale, a.scale)

	ac, ok := mulPow10(a.coef, factor-a.scale)
	if !ok {
		return decimal128{}, decimal128{}, ErrOverflow
	}

	bc, ok := mulPow10(b.coef, factor-b.scale)
	if !ok {
		return decimal128{}, decimal128{}, ErrOverflow
	}

	qc, rc, ok := quoRem128by128(ac, bc)
	if !ok {
		return decimal128{}, decimal128{}, ErrOverflow
	}

	q = newDec(a.neg != b.neg, qc, 0)
	r = newDec(a.neg, rc, factor)

	return q, r, nil
}

func (a decimal128) Mod(b decimal128) (decimal128, error) {
	_, r, err := a.QuoRem(b)

	return r, err
}

// RoundBank rounds a to prec digits after the decimal point using
// round-half-to-even (banker's rounding).
func (a decimal128) RoundBank(prec uint8) decimal128 {
	if prec >= a.scale {
		return a
	}

	factor := pow10[a.scale-prec].lo
	q, r := a.coef.QuoRem64(factor)
	half := factor / 2

	if half < r || (half == r && q.lo%2 == 1) {
		q, _ = q.Add64(1)
	}

	return newDec(a.neg, q, prec)
}

// RoundAway rounds a to prec digits, always rounding away from zero when
// there's any remainder (a.k.a. ROUND_UP in other libraries).
func (a decimal128) RoundAway(prec uint8) decimal128 {
	if prec >= a.scale {
		return a
	}

	factor := pow10[a.scale-prec].lo
	q, r := a.coef.QuoRem64(factor)
	if r != 0 {
		q, _ = q.Add64(1)
	}

	return newDec(a.neg, q, prec)
}

// RoundHAZ rounds a to prec digits using half-away-from-zero.
func (a decimal128) RoundHAZ(prec uint8) decimal128 {
	if prec >= a.scale {
		return a
	}

	factor := pow10[a.scale-prec].lo
	q, r := a.coef.QuoRem64(factor)
	half := factor / 2

	if r >= half {
		q, _ = q.Add64(1)
	}
	return newDec(a.neg, q, prec)
}

// RoundHTZ rounds a to prec digits using half-toward-zero.
func (a decimal128) RoundHTZ(prec uint8) decimal128 {
	if prec >= a.scale {
		return a
	}

	factor := pow10[a.scale-prec].lo
	q, r := a.coef.QuoRem64(factor)
	half := factor / 2

	if r > half {
		q, _ = q.Add64(1)
	}
	return newDec(a.neg, q, prec)
}

// Trunc truncates a to prec digits after the decimal point (no rounding).
func (a decimal128) Trunc(prec uint8) decimal128 {
	if prec >= a.scale {
		return a
	}

	factor := pow10[a.scale-prec].lo
	q, _ := a.coef.QuoRem64(factor)

	return newDec(a.neg, q, prec)
}

// Floor returns the largest integer value <= a.
func (a decimal128) Floor() decimal128 {
	if a.scale == 0 {
		return a
	}

	factor := pow10[a.scale].lo
	q, r := a.coef.QuoRem64(factor)
	if a.neg && r != 0 {
		q, _ = q.Add64(1)
	}

	return newDec(a.neg, q, 0)
}

// Ceil returns the smallest integer value >= a.
func (a decimal128) Ceil() decimal128 {
	if a.scale == 0 {
		return a
	}

	factor := pow10[a.scale].lo
	q, r := a.coef.QuoRem64(factor)
	if !a.neg && r != 0 {
		q, _ = q.Add64(1)
	}

	return newDec(a.neg, q, 0)
}

// Int64 returns the integer part of a, truncated toward zero.
func (a decimal128) Int64() (int64, error) {
	t := a.Trunc(0)
	if t.coef.hi != 0 || t.coef.lo > uint64(math.MaxInt64) {
		return 0, ErrIntPartOverflow
	}

	v := int64(t.coef.lo)
	if t.neg {
		v = -v
	}

	return v, nil
}

func (a decimal128) InexactFloat64() float64 {
	f, _ := strconv.ParseFloat(a.String(), 64)

	return f
}

// trimTrailingZeros drops trailing zero digits from the fractional part,
// reducing scale accordingly.
func (a decimal128) trimTrailingZeros() decimal128 {
	if a.coef.IsZero() || a.scale == 0 {
		return a
	}

	coef := a.coef
	scale := a.scale
	for scale > 0 {
		q, r := coef.QuoRem64(10)
		if r != 0 {
			break
		}

		coef = q
		scale--
	}

	return decimal128{coef: coef, neg: a.neg, scale: scale}
}

// rescale returns a with scale == prec, only ever increasing precision
// (mirrors Decimal.StringFixed semantics: it never drops significant
// digits, it only pads with zeros).
func (a decimal128) rescale(prec uint8) decimal128 {
	dTrim := a.trimTrailingZeros()
	if prec > maxScale {
		prec = maxScale
	}

	if prec <= dTrim.scale {
		return dTrim
	}

	diff := prec - dTrim.scale
	coef, ok := dTrim.coef.Mul64(pow10[diff].lo)
	if !ok {
		return dTrim
	}

	// Deliberately bypass newDec's zero-canonicalization: this is the one
	// place where a zero value must keep a non-zero scale, so that
	// StringFixed on zero pads out, e.g. 0.StringFixed(5) == "0.00000".
	return decimal128{coef: coef, neg: dTrim.neg, scale: prec}
}

func (a decimal128) String() string {
	if a.coef.IsZero() {
		return "0"
	}

	return a.formatString(a.scale, true)
}

// StringFixed returns the string representation of a with fixed prec.
// Trailing zeros are not removed; if a already has more digits than prec,
// it's returned unchanged (at its own precision).
func (a decimal128) StringFixed(prec uint8) string {
	d := a.rescale(prec)
	if prec < d.scale {
		return d.String()
	}

	return d.formatString(d.scale, false)
}

func (a decimal128) formatString(scale uint8, trim bool) string {
	digits := a.coef.String()
	intLen := len(digits) - int(scale)

	var intPart, fracPart string
	if intLen <= 0 {
		intPart = "0"
		fracPart = strings.Repeat("0", -intLen) + digits
	} else {
		intPart = digits[:intLen]
		fracPart = digits[intLen:]
	}

	if trim {
		fracPart = strings.TrimRight(fracPart, "0")
	}

	var b strings.Builder

	if a.neg {
		b.WriteByte('-')
	}

	b.WriteString(intPart)

	if len(fracPart) > 0 {
		b.WriteByte('.')
		b.WriteString(fracPart)
	}

	return b.String()
}
