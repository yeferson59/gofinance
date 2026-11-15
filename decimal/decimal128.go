package decimal

import (
	"math"
	"math/bits"
	"strconv"
	"strings"
)

// maxScale is the maximum number of digits decimal128 can keep after the
// decimal point. It's chosen so that a power of ten up to 10^maxScale still
// fits in a uint64 (10^19 < 2^64 <= 10^20), which keeps scale-alignment
// (Add/Sub/Cmp/rounding) on the fast, allocation-free uint64 division path.
const maxScale uint8 = 19

var (
	decZero = decimal128{}
	decOne  = decimal128{coef: u128One}
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
	return decFromFloat(f, 64)
}

// decFromFloat converts a float of the given bit size (32 or 64) held in f.
// The bit size picks which shortest form is used, so a float32 0.1 converts
// to 0.1 rather than to its float64 widening, 0.10000000149011612.
func decFromFloat(f float64, bitSize int) (decimal128, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return decimal128{}, ErrInvalidFormat
	}

	// The shortest form that round-trips is preferred, so 0.1 stays 0.1
	// rather than its full binary expansion.
	s := strconv.FormatFloat(f, 'f', -1, bitSize)
	if dot := strings.IndexByte(s, '.'); dot < 0 || len(s)-dot-1 <= int(maxScale) {
		return parseDecimal(s)
	}

	// The shortest form needs more fractional digits than a decimal128 keeps
	// (e.g. 1/4100 or 0.001/12). The float is itself only an approximation,
	// so rejecting it would protect no real precision: round its exact binary
	// value to maxScale digits instead, and drop the zeros that leaves.
	d, err := parseDecimal(strconv.FormatFloat(f, 'f', int(maxScale), bitSize))
	if err != nil {
		return decimal128{}, err
	}

	return d.trimTrailingZeros(), nil
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
//
// Aligning the operands' scales multiplies one coefficient by a power of ten,
// which can overflow even when the sum itself is small — a quotient such as
// 1/2 carries 19 fractional digits (0.5000000000000000000), and lining a
// large integer up with it needs 19 more digits than the integer has. On
// overflow Add therefore retries once with both operands' trailing zeros
// dropped; the retry costs nothing unless the first attempt overflowed.
func (a decimal128) Add(b decimal128) (decimal128, error) {
	return a.add(b, true)
}

// addOverflow is add's overflow path: retry once with trailing zeros dropped
// when retry allows and trimming changed something.
func (a decimal128) addOverflow(b decimal128, retry bool) (decimal128, error) {
	if !retry {
		return decimal128{}, ErrOverflow
	}

	if at, bt, trimmed := trimBoth(a, b); trimmed {
		return at.add(bt, false)
	}

	return decimal128{}, ErrOverflow
}

// trimBoth drops trailing fractional zeros from a and b, reporting whether
// either changed: if neither did, retrying an overflowed operation on them
// cannot help.
func trimBoth(a, b decimal128) (decimal128, decimal128, bool) {
	at, bt := a.trimTrailingZeros(), b.trimTrailingZeros()

	return at, bt, at.scale != a.scale || bt.scale != b.scale
}

// add returns a+b; on overflow it defers to addOverflow, which retries only
// when retry is true so the retry itself cannot recurse.
func (a decimal128) add(b decimal128, retry bool) (decimal128, error) {
	scale := a.scale
	ac, bc := a.coef, b.coef
	var ok bool

	switch {
	case a.scale < b.scale:
		scale = b.scale
		if ac, ok = mulPow10(ac, b.scale-a.scale); !ok {
			return a.addOverflow(b, retry)
		}
	case a.scale > b.scale:
		if bc, ok = mulPow10(bc, a.scale-b.scale); !ok {
			return a.addOverflow(b, retry)
		}
	}

	if a.neg == b.neg {
		sum, ok := ac.Add(bc)
		if !ok {
			return a.addOverflow(b, retry)
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
// point (no rounding is applied when truncating excess precision). Like Add,
// it retries once with trailing zeros dropped if the product overflows.
func (a decimal128) Mul(b decimal128) (decimal128, error) {
	return a.mul(b, true)
}

// mulOverflow is mul's overflow path, the counterpart of addOverflow.
func (a decimal128) mulOverflow(b decimal128, retry bool) (decimal128, error) {
	if !retry {
		return decimal128{}, ErrOverflow
	}

	if at, bt, trimmed := trimBoth(a, b); trimmed {
		return at.mul(bt, false)
	}

	return decimal128{}, ErrOverflow
}

// mul returns a*b; on overflow it defers to mulOverflow.
func (a decimal128) mul(b decimal128, retry bool) (decimal128, error) {
	if a.coef.IsZero() || b.coef.IsZero() {
		return decimal128{}, nil
	}

	neg := a.neg != b.neg
	totalScale := int(a.scale) + int(b.scale)
	prod := a.coef.MulFull(b.coef)

	if totalScale <= int(maxScale) {
		if !prod.carry.IsZero() {
			return a.mulOverflow(b, retry)
		}
		//nolint:gosec // totalScale <= maxScale (19)
		return newDec(neg, u128{hi: prod.hi, lo: prod.lo}, uint8(totalScale)), nil
	}

	excess := totalScale - int(maxScale)
	q, _, ok := prod.QuoRem128(pow10[excess])
	if !ok {
		return a.mulOverflow(b, retry)
	}

	return newDec(neg, q, maxScale), nil
}

// quoDensest returns trunc(coef/10^coefScale ÷ den/10^denScale) as a
// coefficient at the largest scale in [0, maxScale] whose coefficient fits in
// 128 bits, together with that scale. ok is false only when even the integer
// part doesn't fit.
//
// Working at maxScale is the norm and costs one division. A quotient too
// large for 19 fractional digits — 1e20/1, say — gives up fractional digits
// one at a time rather than overflowing, the same "densest representable
// scale" rule Pow follows. The search stops at coefScale-denScale, where the
// numerator needs no scaling at all and the quotient (at most coef) is
// guaranteed to fit.
func quoDensest(coef u128, coefScale uint8, den u128, denScale uint8) (u128, uint8, bool) {
	minScale := max(int(coefScale)-int(denScale), 0)

	for sc := int(maxScale); sc >= minScale; sc-- {
		factor := sc - int(coefScale) + int(denScale) // in [0, 38]

		if q, _, ok := coef.MulFull(pow10[factor]).QuoRem128(den); ok {
			return q, uint8(sc), true
		}
	}

	return u128{}, 0, false
}

// Div returns a/b truncated to maxScale digits after the decimal point, or to
// fewer when the quotient is too large to carry all of them.
func (a decimal128) Div(b decimal128) (decimal128, error) {
	if b.coef.IsZero() {
		return decimal128{}, ErrDivideByZero
	}

	if a.coef.IsZero() {
		return decimal128{}, nil
	}

	// The common case: one division at maxScale. quoDensest is only for a
	// quotient too large to carry every fractional digit.
	factor := int(maxScale) - int(a.scale) + int(b.scale)
	if q, _, ok := a.coef.MulFull(pow10[factor]).QuoRem128(b.coef); ok {
		return newDec(a.neg != b.neg, q, maxScale), nil
	}

	q, scale, ok := quoDensest(a.coef, a.scale, b.coef, b.scale)
	if !ok {
		return decimal128{}, ErrOverflow
	}

	return newDec(a.neg != b.neg, q, scale), nil
}

// Div64 returns a/v with the same precision rules as Div. Dividing by an
// integer never grows the magnitude, so it cannot overflow.
func (a decimal128) Div64(v uint64) (decimal128, error) {
	if v == 0 {
		return decimal128{}, ErrDivideByZero
	}

	if a.coef.IsZero() {
		return decimal128{}, nil
	}

	if q, _, ok := a.coef.MulFull(pow10[maxScale-a.scale]).QuoRem128(u128FromU64(v)); ok {
		return newDec(a.neg, q, maxScale), nil
	}

	q, scale, ok := quoDensest(a.coef, a.scale, u128FromU64(v), 0)
	if !ok {
		return decimal128{}, ErrOverflow
	}

	return newDec(a.neg, q, scale), nil
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
	if t.coef.hi != 0 {
		return 0, ErrIntPartOverflow
	}

	// A negative value reaches one further than a positive one: MinInt64's
	// magnitude, 2^63, has no positive int64 counterpart.
	if t.neg {
		switch {
		case t.coef.lo == uint64(math.MaxInt64)+1:
			return math.MinInt64, nil
		case t.coef.lo > uint64(math.MaxInt64):
			return 0, ErrIntPartOverflow
		default:
			return -int64(t.coef.lo), nil
		}
	}

	if t.coef.lo > uint64(math.MaxInt64) {
		return 0, ErrIntPartOverflow
	}

	return int64(t.coef.lo), nil
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

	// Most coefficients — any inexact quotient — end in no zero at all, so
	// test divisibility by 10 without dividing: it needs an even coefficient
	// that is also a multiple of 5, and because 2^64 ≡ 1 (mod 5), the
	// coefficient's residue mod 5 is that of hi + lo. The compiler turns a
	// constant % 5 into multiplications, not a hardware division.
	if a.coef.lo&1 != 0 || (a.coef.hi%5+a.coef.lo%5)%5 != 0 {
		return a
	}

	// 10^k divides the coefficient only if 2^k does, so its trailing zero
	// bits bound how many decimal zeros it can end in.
	tz := bits.TrailingZeros64(a.coef.lo)
	if a.coef.lo == 0 {
		tz = 64 + bits.TrailingZeros64(a.coef.hi)
	}

	// Binary search for the largest k ≤ min(tz, scale) with 10^k dividing the
	// coefficient: divisibility by 10^k implies it for every smaller power,
	// so at most five divisions settle it instead of one per digit.
	lo, hi := 0, min(tz, int(a.scale))
	coef := a.coef

	for lo < hi {
		mid := (lo + hi + 1) / 2

		if q, r := a.coef.QuoRem64(pow10[mid].lo); r == 0 {
			lo, coef = mid, q
		} else {
			hi = mid - 1
		}
	}

	return decimal128{coef: coef, neg: a.neg, scale: a.scale - uint8(lo)}
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

// StringFixed returns a formatted with exactly prec fractional digits,
// rounding half away from zero when a carries more and padding with zeros
// when it carries fewer. prec is capped at maxScale.
func (a decimal128) StringFixed(prec uint8) string {
	return a.stringFixedWith(prec, decimal128.RoundHAZ)
}

// StringFixedBank is like StringFixed but rounds half to even.
func (a decimal128) StringFixedBank(prec uint8) string {
	return a.stringFixedWith(prec, decimal128.RoundBank)
}

// stringFixedWith formats a at exactly prec fractional digits, using round to
// drop any digits beyond prec.
func (a decimal128) stringFixedWith(prec uint8, round func(decimal128, uint8) decimal128) string {
	prec = min(prec, maxScale)
	if prec < a.scale {
		a = round(a, prec)
	}

	// After rounding a never carries more than prec digits, so rescale only
	// pads. It returns fewer digits only for a coefficient too close to the
	// 128-bit ceiling to pad, which is then printed at its own precision.
	d := a.rescale(prec)

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
