package money

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
	if x.carry.hi != 0 {
		return 192 + bits.Len64(x.carry.hi)
	}
	if x.carry.lo != 0 {
		return 128 + bits.Len64(x.carry.lo)
	}
	if x.hi != 0 {
		return 64 + bits.Len64(x.hi)
	}
	return bits.Len64(x.lo)
}

func (x u256) bitAt(i int) uint64 {
	switch {
	case i < 64:
		return (x.lo >> uint(i)) & 1
	case i < 128:
		return (x.hi >> uint(i-64)) & 1
	case i < 192:
		return (x.carry.lo >> uint(i-128)) & 1
	default:
		return (x.carry.hi >> uint(i-192)) & 1
	}
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

	return binaryDivU128(x.bitAt, x.bitLen(), v)
}
