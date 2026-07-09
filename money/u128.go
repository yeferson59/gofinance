package money

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

// bitAt returns the i-th bit (0 = LSB) of u.
func (u u128) bitAt(i int) uint64 {
	if i >= 64 {
		return (u.hi >> uint(i-64)) & 1
	}
	return (u.lo >> uint(i)) & 1
}

// setBit sets the i-th bit (0 = LSB) of u to 1.
func (u *u128) setBit(i int) {
	if i >= 64 {
		u.hi |= 1 << uint(i-64)
	} else {
		u.lo |= 1 << uint(i)
	}
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
	// v needs the full 128 bits, so the quotient is guaranteed to fit
	// comfortably within 128 bits too.
	q, r, _ = binaryDivU128(u.bitAt, u.bitLen(), v)
	return q, r, true
}

// binaryDivU128 performs schoolbook binary long division of an arbitrary
// width dividend (described by bitAt/bitLen) by the 128-bit divisor v,
// producing a 128-bit quotient. ok is false if the true quotient doesn't
// fit in 128 bits.
func binaryDivU128(bitAt func(i int) uint64, bitLen int, v u128) (q, r u128, ok bool) {
	if v.IsZero() {
		return u128{}, u128{}, false
	}

	var rem u128
	var remCarry uint64
	var quo u128

	for i := bitLen - 1; i >= 0; i-- {
		remCarry = (remCarry << 1) | (rem.hi >> 63)
		rem.hi = (rem.hi << 1) | (rem.lo >> 63)
		rem.lo = (rem.lo << 1) | bitAt(i)

		if remCarry != 0 || rem.Cmp(v) >= 0 {
			lo, b1 := bits.Sub64(rem.lo, v.lo, 0)
			hi, _ := bits.Sub64(rem.hi, v.hi, b1)
			rem = u128{hi: hi, lo: lo}
			remCarry = 0

			if i >= 128 {
				return u128{}, u128{}, false
			}
			quo.setBit(i)
		}
	}

	return quo, rem, true
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
