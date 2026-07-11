package decimal

import "math/bits"

// u256 is an unsigned 256-bit integer, used only as an intermediate value
// when multiplying/dividing 128-bit coefficients so that no precision is
// lost before the result is rescaled back down to 128 bits.
//
// value = carry*2^128 + hi*2^64 + lo
type u256 struct {
	hi, lo uint64
	carry  u128
}

func (x u256) IsZero() bool {
	return x.hi == 0 && x.lo == 0 && x.carry.IsZero()
}

func (x u256) bitLen() int {
	switch {
	case x.carry.hi != 0:
		return 192 + bits.Len64(x.carry.hi)
	case x.carry.lo != 0:
		return 128 + bits.Len64(x.carry.lo)
	case x.hi != 0:
		return 64 + bits.Len64(x.hi)
	default:
		return bits.Len64(x.lo)
	}
}

// cmp compares x and y, returning -1, 0 or 1.
func (x u256) cmp(y u256) int {
	if c := x.carry.Cmp(y.carry); c != 0 {
		return c
	}

	return (u128{hi: x.hi, lo: x.lo}).Cmp(u128{hi: y.hi, lo: y.lo})
}

// quoRem64 returns q = x/v and r = x%v for a 64-bit divisor v (v != 0),
// keeping the full 256-bit quotient.
func (x u256) quoRem64(v uint64) (q u256, r uint64) {
	q.carry.hi, r = bits.Div64(0, x.carry.hi, v)
	q.carry.lo, r = bits.Div64(r, x.carry.lo, v)
	q.hi, r = bits.Div64(r, x.hi, v)
	q.lo, r = bits.Div64(r, x.lo, v)

	return q, r
}

// QuoRem128 returns q = x/v and r = x%v. ok is false if v is zero or the
// quotient doesn't fit in 128 bits.
func (x u256) QuoRem128(v u128) (q, r u128, ok bool) {
	if v.IsZero() {
		return u128{}, u128{}, false
	}

	if x.carry.IsZero() {
		return quoRem128by128(u128{hi: x.hi, lo: x.lo}, v)
	}

	if v.hi == 0 && x.carry.hi == 0 {
		if x.carry.lo >= v.lo {
			// quotient would need more than 128 bits
			return u128{}, u128{}, false
		}

		hiQ, rem1 := bits.Div64(x.carry.lo, x.hi, v.lo)
		loQ, rem2 := bits.Div64(rem1, x.lo, v.lo)
		return u128{hi: hiQ, lo: loQ}, u128{lo: rem2}, true
	}

	if x.carry.Cmp(v) >= 0 {
		// quotient can't fit in 128 bits
		return u128{}, u128{}, false
	}

	q, r = x.quoRemKnuth(v)

	return q, r, true
}

// quoRemKnuth divides x by a full-width divisor v (v.hi != 0) with two
// normalized 3-by-2 word division steps (Knuth's algorithm D in base
// 2^64). It requires x.carry < v, which guarantees the quotient fits in
// 128 bits.
func (x u256) quoRemKnuth(v u128) (q, r u128) {
	s := uint(bits.LeadingZeros64(v.hi))
	v1 := v.hi<<s | v.lo>>1>>(63-s)
	v0 := v.lo << s

	// x << s spread over four words: the would-be fifth word is always
	// zero, because x.carry < v keeps x.carry.hi below 2^(64-s).
	d3 := x.carry.hi<<s | x.carry.lo>>1>>(63-s)
	d2 := x.carry.lo<<s | x.hi>>1>>(63-s)
	d1 := x.hi<<s | x.lo>>1>>(63-s)
	d0 := x.lo << s

	q1, rh, rl := div3by2(d3, d2, d1, v1, v0)
	q0, rh, rl := div3by2(rh, rl, d0, v1, v0)

	return u128{hi: q1, lo: q0}, u128{hi: rh >> s, lo: rl>>s | rh<<1<<(63-s)}
}
