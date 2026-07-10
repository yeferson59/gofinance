package decimal

import (
	"math/bits"
	"strconv"
)

// u128 is an unsigned 128-bit integer represented as two 64-bit words.
// value = hi*2^64 + lo
type u128 struct {
	hi uint64
	lo uint64
}

var (
	u128Zero = u128{}
	u128One  = u128{lo: 1}
)

// pow10 holds precomputed powers of ten, pow10[i] == 10^i, for i in [0, 38].
// 10^38 is the largest power of ten that still fits in 128 bits.
var pow10 = computePow10Table()

func computePow10Table() [39]u128 {
	var t [39]u128
	t[0] = u128One
	for i := 1; i < len(t); i++ {
		hi, lo := bits.Mul64(t[i-1].lo, 10)
		// t[i-1].hi is small enough (< 10^20) that hi*10 never overflows
		// a uint64 for any i <= 38.
		hi += t[i-1].hi * 10
		t[i] = u128{hi: hi, lo: lo}
	}
	return t
}

func u128FromU64(v uint64) u128 {
	return u128{lo: v}
}

func (u u128) IsZero() bool {
	return u.hi == 0 && u.lo == 0
}

// Cmp compares u and v, returning -1, 0 or 1.
func (u u128) Cmp(v u128) int {
	switch {
	case u.hi < v.hi:
		return -1
	case u.hi > v.hi:
		return 1
	case u.lo < v.lo:
		return -1
	case u.lo > v.lo:
		return 1
	default:
		return 0
	}
}

// Cmp64 compares u against the 64-bit value v.
func (u u128) Cmp64(v uint64) int {
	if u.hi != 0 {
		return 1
	}
	switch {
	case u.lo < v:
		return -1
	case u.lo > v:
		return 1
	default:
		return 0
	}
}

// Add returns u+v. ok is false if the sum overflows 128 bits.
func (u u128) Add(v u128) (r u128, ok bool) {
	lo, carry := bits.Add64(u.lo, v.lo, 0)
	hi, carry := bits.Add64(u.hi, v.hi, carry)
	if carry != 0 {
		return u128{}, false
	}
	return u128{hi: hi, lo: lo}, true
}

// Add64 returns u+v where v is a 64-bit value.
func (u u128) Add64(v uint64) (r u128, ok bool) {
	lo, carry := bits.Add64(u.lo, v, 0)
	hi, carry := bits.Add64(u.hi, 0, carry)
	if carry != 0 {
		return u128{}, false
	}
	return u128{hi: hi, lo: lo}, true
}

// Sub returns u-v, assuming u >= v. ok is false if u < v (borrow out).
func (u u128) Sub(v u128) (r u128, ok bool) {
	lo, borrow := bits.Sub64(u.lo, v.lo, 0)
	hi, borrow := bits.Sub64(u.hi, v.hi, borrow)
	if borrow != 0 {
		return u128{}, false
	}
	return u128{hi: hi, lo: lo}, true
}

// Mul64 returns u*v where v is a 64-bit value. ok is false on overflow.
func (u u128) Mul64(v uint64) (r u128, ok bool) {
	hi, lo := bits.Mul64(u.lo, v)
	if u.hi == 0 {
		return u128{hi: hi, lo: lo}, true
	}

	p0, p1 := bits.Mul64(u.hi, v)
	hi, c := bits.Add64(hi, p1, 0)
	if p0 != 0 || c != 0 {
		return u128{}, false
	}
	return u128{hi: hi, lo: lo}, true
}

// Mul returns u*v. ok is false on overflow (result doesn't fit in 128 bits).
func (u u128) Mul(v u128) (r u128, ok bool) {
	if u.hi != 0 && v.hi != 0 {
		// both operands are >= 2^64, so the product is guaranteed >= 2^128.
		return u128{}, false
	}
	if v.hi == 0 {
		return u.Mul64(v.lo)
	}
	return v.Mul64(u.lo)
}

// MulFull returns the full, non-truncated 256-bit product of u*v.
func (u u128) MulFull(v u128) u256 {
	hi, lo := bits.Mul64(u.lo, v.lo)
	p0, p1 := bits.Mul64(u.hi, v.lo)
	p2, p3 := bits.Mul64(u.lo, v.hi)

	hi, c0 := bits.Add64(hi, p1, 0)
	hi, c1 := bits.Add64(hi, p3, 0)
	carryLo := c0 + c1

	e0, e1 := bits.Mul64(u.hi, v.hi)
	d, d0 := bits.Add64(p0, p2, 0)
	d, d1 := bits.Add64(d, carryLo, 0)
	upperLo, upperHi := bits.Add64(d, e1, 0)

	return u256{
		lo: lo,
		hi: hi,
		carry: u128{
			// e0+d0+d1+upperHi can't overflow a uint64: the true upper bound
			// of the product's top word is bounded by 10^38's own magnitude.
			hi: e0 + d0 + d1 + upperHi,
			lo: upperLo,
		},
	}
}

// QuoRem64 returns q = u/v and r = u%v for a 64-bit divisor v (v != 0).
func (u u128) QuoRem64(v uint64) (q u128, r uint64) {
	if u.hi < v {
		q.lo, r = bits.Div64(u.hi, u.lo, v)
		return q, r
	}
	q.hi, r = bits.Div64(0, u.hi, v)
	q.lo, r = bits.Div64(r, u.lo, v)
	return q, r
}

func (u u128) bitLen() int {
	if u.hi != 0 {
		return 64 + bits.Len64(u.hi)
	}
	return bits.Len64(u.lo)
}

// decimalDigits returns the number of base-10 digits of u (1 for zero).
func (u u128) decimalDigits() int {
	if u.IsZero() {
		return 1
	}

	// 1233/4096 approximates log10(2); the estimate is exact or one low.
	d := (u.bitLen() * 1233) >> 12
	if d < len(pow10) && u.Cmp(pow10[d]) >= 0 {
		d++
	}

	return d
}

// bitAt returns the i-th bit (0 = LSB) of u.
func (u u128) bitAt(i int) uint64 {
	if i >= 64 {
		return (u.hi >> uint(i-64)) & 1
	}
	return (u.lo >> uint(i)) & 1
}

// quoRem128by128 returns q = u/v and r = u%v for a 128-bit divisor v.
// ok is false if v is zero.
func quoRem128by128(u, v u128) (q, r u128, ok bool) {
	if v.IsZero() {
		return u128{}, u128{}, false
	}
	if v.hi == 0 {
		qq, rr := u.QuoRem64(v.lo)
		return qq, u128{lo: rr}, true
	}

	// v uses the full 128 bits, so the quotient fits in a single word
	// (u < 2^128 <= v*2^64): one normalized 3-by-2 division step suffices.
	s := uint(bits.LeadingZeros64(v.hi))
	v1 := v.hi<<s | v.lo>>1>>(63-s)
	v0 := v.lo << s

	d2 := u.hi >> 1 >> (63 - s)
	d1 := u.hi<<s | u.lo>>1>>(63-s)
	d0 := u.lo << s

	qw, r1, r0 := div3by2(d2, d1, d0, v1, v0)

	return u128{lo: qw}, u128{hi: r1 >> s, lo: r0>>s | r1<<1<<(63-s)}, true
}

// div3by2 divides the three-word value [a2 a1 a0] by the normalized
// two-word divisor [v1 v0] (top bit of v1 set), assuming [a2 a1] < [v1 v0]
// so the quotient fits in a single word. It returns the quotient word and
// the two-word remainder. This is the base-2^64 schoolbook division step
// of Knuth's algorithm D.
func div3by2(a2, a1, a0, v1, v0 uint64) (q, r1, r0 uint64) {
	var (
		qhat, rhat uint64
		rhatBig    bool // rhat >= 2^64: no further correction can apply
	)

	if a2 == v1 {
		// bits.Div64 requires a2 < v1; here the estimate is the maximum
		// word value instead.
		qhat = ^uint64(0)

		var c uint64
		rhat, c = bits.Add64(a1, v1, 0)
		rhatBig = c != 0
	} else {
		qhat, rhat = bits.Div64(a2, a1, v1)
	}

	// Lower qhat until qhat*[v1 v0] <= [a2 a1 a0], comparing via
	// qhat*v0 <= rhat*2^64 + a0 with rhat = [a2 a1] - qhat*v1. With a
	// normalized divisor the estimate overshoots by at most 2.
	for !rhatBig {
		ph, pl := bits.Mul64(qhat, v0)
		if ph < rhat || (ph == rhat && pl <= a0) {
			break
		}

		qhat--

		var c uint64
		rhat, c = bits.Add64(rhat, v1, 0)
		rhatBig = c != 0
	}

	// remainder = [a2 a1 a0] - qhat*[v1 v0]. The correction loop above
	// guarantees it fits in two words, so working modulo 2^128 (dropping
	// the top words of both operands) is exact.
	ph0, pl0 := bits.Mul64(qhat, v0)
	_, pl1 := bits.Mul64(qhat, v1)
	p1, _ := bits.Add64(ph0, pl1, 0)

	r0, b := bits.Sub64(a0, pl0, 0)
	r1, _ = bits.Sub64(a1, p1, b)

	return qhat, r1, r0
}

// String returns the base-10 representation of u.
func (u u128) String() string {
	if u.IsZero() {
		return "0"
	}

	// Break u into 19-digit, base-10^19 chunks (least significant first).
	var chunks []uint64
	for !u.IsZero() {
		var r uint64
		u, r = u.QuoRem64(1e19)
		chunks = append(chunks, r)
	}

	last := len(chunks) - 1
	buf := strconv.AppendUint(nil, chunks[last], 10)
	for i := last - 1; i >= 0; i-- {
		s := strconv.FormatUint(chunks[i], 10)
		for pad := 19 - len(s); pad > 0; pad-- {
			buf = append(buf, '0')
		}
		buf = append(buf, s...)
	}
	return string(buf)
}
