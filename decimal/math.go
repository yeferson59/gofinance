package decimal

import "math/bits"

// This file implements Pow, Ln, Log2, Log10 and Log natively on
// decimal128. All intermediate arithmetic runs in "fp120" binary fixed
// point: a u128 magnitude scaled by 2^120, plus an explicit sign where one
// is needed. fp120 resolves 2^-120 ≈ 7.5e-37 — about 36 decimal digits —
// so results come out correctly rounded to decimal128's 19 digits, and
// reducing a 256-bit product back to fp120 is a bit shift instead of a
// long division.
//
// Domain sizing: the integer part of an fp120 value spans up to 2^8 = 256,
// which covers every magnitude these kernels see (|ln x| <= ln(2^128) ≈
// 88.7 and |log2 x| < 129 for any positive decimal128).

// The constants below are round(c * 2^120), generated with math/big at
// 500-bit precision and cross-checked against math/big in the tests.
var (
	fpOne     = u128{hi: 1 << 56}                                    // 1.0
	fpTwo     = u128{hi: 1 << 57}                                    // 2.0
	ln2Fp     = u128{hi: 0x00b17217f7d1cf79, lo: 0xabc9e3b39803f2f7} // ln(2)
	ln10Fp    = u128{hi: 0x024d763776aaa2b0, lo: 0x5ba95b58ae0b4c29} // ln(10)
	invLn2Fp  = u128{hi: 0x0171547652b82fe1, lo: 0x777d0ffda0d23a7d} // 1/ln(2)
	invLn10Fp = u128{hi: 0x006f2dec549b9438, lo: 0xca9aadd557d699ee} // 1/ln(10)
)

// shr120Round returns round(x / 2^120), i.e. it rescales a 256-bit product
// of two fp120 values back to fp120, rounding half up. ok is false if the
// result doesn't fit in 128 bits.
func (x u256) shr120Round() (u128, bool) {
	// Add 2^119 (bit 55 of the second word) so the truncating shift rounds.
	hi, c := bits.Add64(x.hi, 1<<55, 0)
	clo, c := bits.Add64(x.carry.lo, 0, c)
	chi, c := bits.Add64(x.carry.hi, 0, c)

	if c != 0 || chi>>56 != 0 {
		return u128{}, false
	}

	return u128{
		hi: clo>>56 | chi<<8,
		lo: hi>>56 | clo<<8,
	}, true
}

// shl120 returns x * 2^120 as a 256-bit value.
func shl120(x u128) u256 {
	return u256{
		hi: x.lo << 56,
		carry: u128{
			lo: x.lo>>8 | x.hi<<56,
			hi: x.hi >> 8,
		},
	}
}

// inc256 returns x+1, ignoring overflow past 256 bits.
func inc256(x u256) u256 {
	lo, c := bits.Add64(x.lo, 1, 0)
	hi, c := bits.Add64(x.hi, 0, c)
	clo, c := bits.Add64(x.carry.lo, 0, c)
	chi, _ := bits.Add64(x.carry.hi, 0, c)

	return u256{lo: lo, hi: hi, carry: u128{lo: clo, hi: chi}}
}

// fpMul returns round(a*b / 2^120), the fp120 product. Callers must
// guarantee the result fits in 128 bits.
func fpMul(a, b u128) u128 {
	r, _ := a.MulFull(b).shr120Round()

	return r
}

// divRound returns round(x/v), rounding half up. ok is false if the
// quotient doesn't fit in 128 bits (or v is zero).
func divRound(x u256, v u128) (u128, bool) {
	q, r, ok := x.QuoRem128(v)
	if !ok {
		return u128{}, false
	}

	if double, ok2 := r.Add(r); !ok2 || double.Cmp(v) >= 0 {
		if q, ok = q.Add64(1); !ok {
			return u128{}, false
		}
	}

	return q, true
}

// lnFp computes ln(a) for a > 0 as an fp120 magnitude plus sign, accurate
// to within a few ulps of 2^-120.
//
// Argument reduction: writing a = m * 10^e10 with m in [1, 10) and then
// halving m into [1, 2),
//
//	ln(a) = e10*ln(10) + k2*ln(2) + 2*atanh(z), z = (m'-1)/(m'+1)
//
// with z in [0, 1/3), so the atanh series gains a factor >= 9 per term and
// converges to full fp120 precision in at most ~38 terms.
func (a decimal128) lnFp() (mag u128, neg bool) {
	d := a.coef.decimalDigits()
	e10 := d - 1 - int(a.scale)

	// m = coef * 2^120 / 10^(d-1), in [2^120, 10*2^120). 10^(d-1) can need
	// up to 38 digits, so divide in one or two 64-bit steps (each power
	// fits a uint64), rounding on the last one.
	x := shl120(a.coef)
	rest := d - 1
	if rest > int(maxScale) {
		x, _ = x.quoRem64(pow10[maxScale].lo)
		rest -= int(maxScale)
	}

	q, r := x.quoRem64(pow10[rest].lo)
	m := u128{hi: q.hi, lo: q.lo}
	if v := pow10[rest].lo; r >= v-v/2 {
		m, _ = m.Add64(1)
	}

	// Halve m into [1, 2), folding the halvings into k2 (at most 3).
	k2 := 0
	for m.Cmp(fpTwo) >= 0 {
		m, _ = m.Add64(1) // round the halving
		m = u128{hi: m.hi >> 1, lo: m.lo>>1 | m.hi<<63}
		k2++
	}

	// z = (m-1)/(m+1) in [0, 1/3).
	num, _ := m.Sub(fpOne)
	den, _ := m.Add(fpOne)
	z, _ := divRound(shl120(num), den)

	// atanh series: z + z^3/3 + z^5/5 + ...
	sum := z
	zsq := fpMul(z, z)
	term := z
	for k := uint64(3); !term.IsZero(); k += 2 {
		term = fpMul(term, zsq)
		t, _ := term.QuoRem64(k)
		if t.IsZero() {
			break
		}
		sum, _ = sum.Add(t)
	}

	// pos = k2*ln(2) + 2*atanh(z), always >= 0.
	pos := u128{hi: sum.hi<<1 | sum.lo>>63, lo: sum.lo << 1}
	if k2 > 0 {
		kl, _ := ln2Fp.Mul64(uint64(k2))
		pos, _ = pos.Add(kl)
	}

	if e10 >= 0 {
		el, _ := ln10Fp.Mul64(uint64(e10))
		mag, _ = pos.Add(el)

		return mag, false
	}

	el, _ := ln10Fp.Mul64(uint64(-e10))
	if el.Cmp(pos) >= 0 {
		mag, _ = el.Sub(pos)

		return mag, true
	}

	mag, _ = pos.Sub(el)

	return mag, false
}

// fpToDec converts an fp120 magnitude plus sign to a decimal128 rounded
// (half away from zero) to maxScale fractional digits.
func fpToDec(mag u128, neg bool) decimal128 {
	coef, _ := mag.MulFull(pow10[maxScale]).shr120Round()

	return newDec(neg, coef, maxScale)
}

// expFpToDec computes e^x for x = ±mag/2^120, rounded to the densest scale
// decimal128 can represent the result at. mag must be below 2^127.
//
// Argument reduction: x = k10*ln(10) + k2*ln(2) + s with s in [0, ln 2),
// so e^x = 10^k10 * 2^k2 * e^s and the Taylor series for e^s converges to
// full fp120 precision in at most ~32 terms.
func expFpToDec(mag u128, neg bool) (decimal128, error) {
	// k10 = floor(x / ln(10)), r = x - k10*ln(10) in [0, ln 10).
	t := fpMul(mag, invLn10Fp)
	est := int(t.hi >> 56)

	var (
		k10 int
		r   u128
	)

	if !neg {
		k10 = est
		kl, _ := ln10Fp.Mul64(uint64(est))

		var ok bool
		if r, ok = mag.Sub(kl); !ok {
			k10--
			kl, _ = kl.Sub(ln10Fp)
			r, _ = mag.Sub(kl)
		}
	} else {
		kl, _ := ln10Fp.Mul64(uint64(est))
		if diff, ok := kl.Sub(mag); ok {
			k10, r = -est, diff
		} else {
			k10 = -(est + 1)
			kl, _ = kl.Add(ln10Fp)
			r, _ = kl.Sub(mag)
		}
	}

	for r.Cmp(ln10Fp) >= 0 {
		r, _ = r.Sub(ln10Fp)
		k10++
	}

	// The result is 10^k10 * e^r with e^r in [1, 10).
	if k10 > 38 {
		return decimal128{}, ErrOverflow
	}

	if k10 < -(int(maxScale) + 1) {
		// Below half of decimal128's smallest positive value: rounds to 0.
		return decimal128{}, nil
	}

	// Split off k2 = floor(r / ln(2)) <= 3.
	k2 := uint(0)
	for r.Cmp(ln2Fp) >= 0 {
		r, _ = r.Sub(ln2Fp)
		k2++
	}

	// Taylor series: e^s = 1 + s + s^2/2! + ...
	sum, _ := fpOne.Add(r)
	term := r
	for k := uint64(2); !term.IsZero(); k++ {
		term = fpMul(term, r)
		term, _ = term.QuoRem64(k)
		if term.IsZero() {
			break
		}
		sum, _ = sum.Add(term)
	}

	e := u128{hi: sum.hi<<k2 | sum.lo>>(64-k2), lo: sum.lo << k2}

	// Choose the densest representable scale: sc = 19 while the whole
	// coefficient (k10+sc+1 digits) fits, sliding down to 0 as the result
	// approaches decimal128's upper range.
	sc := max(min(37-k10, int(maxScale)), 0)

	var (
		coef u128
		ok   bool
	)

	if p := k10 + sc; p >= 0 {
		coef, ok = e.MulFull(pow10[p]).shr120Round()
	} else {
		// Only k10 == -20, sc == 19: coef = round(e / 10) may still round
		// up to decimal128's smallest positive value.
		div, _ := fpOne.Mul64(10)
		coef, ok = divRound(u256{hi: e.hi, lo: e.lo}, div)
	}

	if !ok {
		return decimal128{}, ErrOverflow
	}

	return newDec(false, coef, uint8(sc)), nil
}

// Ln returns the natural logarithm of a, rounded half away from zero to
// maxScale digits. It errors if a <= 0.
func (a decimal128) Ln() (decimal128, error) {
	if a.neg || a.coef.IsZero() {
		return decimal128{}, ErrLogNonPositive
	}

	mag, neg := a.lnFp()

	return fpToDec(mag, neg), nil
}

// Log2 returns the base-2 logarithm of a, computed as ln(a)/ln(2).
func (a decimal128) Log2() (decimal128, error) {
	if a.neg || a.coef.IsZero() {
		return decimal128{}, ErrLogNonPositive
	}

	mag, neg := a.lnFp()

	return fpToDec(fpMul(mag, invLn2Fp), neg), nil
}

// Log10 returns the base-10 logarithm of a, computed as ln(a)/ln(10).
func (a decimal128) Log10() (decimal128, error) {
	if a.neg || a.coef.IsZero() {
		return decimal128{}, ErrLogNonPositive
	}

	mag, neg := a.lnFp()

	return fpToDec(fpMul(mag, invLn10Fp), neg), nil
}

// Log returns the base-base logarithm of a, computed as ln(a)/ln(base).
// It errors if a or base are non-positive, or base is 1.
func (a decimal128) Log(base decimal128) (decimal128, error) {
	if a.neg || a.coef.IsZero() || base.neg || base.coef.IsZero() {
		return decimal128{}, ErrLogNonPositive
	}

	bMag, bNeg := base.lnFp()
	if bMag.IsZero() {
		return decimal128{}, ErrDivideByZero
	}

	aMag, aNeg := a.lnFp()
	if aMag.IsZero() {
		return decimal128{}, nil
	}

	neg := aNeg != bNeg

	// coef = round(aMag * 10^19 / bMag) at scale 19; for extreme ratios
	// (base within 1e-19 of 1) that can exceed 128 bits, so fall back to
	// scale 0, where the quotient always fits.
	if coef, ok := divRound(aMag.MulFull(pow10[maxScale]), bMag); ok {
		return newDec(neg, coef, maxScale), nil
	}

	coef, _ := divRound(u256{hi: aMag.hi, lo: aMag.lo}, bMag)

	return newDec(neg, coef, 0), nil
}

// isqrt256 returns floor(sqrt(n)) for n > 0, using Newton's integer
// iteration x <- (x + n/x)/2 from an initial upper bound. Starting at
// 2^ceil(bitLen/2) >= sqrt(n), the sequence decreases monotonically to
// floor(sqrt(n)) and converges quadratically (~8 division steps).
func isqrt256(n u256) u128 {
	x := u128{hi: ^uint64(0), lo: ^uint64(0)}
	if shift := (n.bitLen() + 1) / 2; shift < 128 {
		if shift >= 64 {
			x = u128{hi: 1 << (shift - 64)}
		} else {
			x = u128{lo: 1 << shift}
		}
	}

	for {
		// x >= floor(sqrt(n)) throughout, so q = n/x <= sqrt(n) always
		// fits in 128 bits.
		q, _, _ := n.QuoRem128(x)

		// y = (x + q) / 2, keeping the 129-bit sum's carry.
		lo, c := bits.Add64(x.lo, q.lo, 0)
		hi, c := bits.Add64(x.hi, q.hi, c)
		y := u128{hi: c<<63 | hi>>1, lo: hi<<63 | lo>>1}

		if y.Cmp(x) >= 0 {
			return x
		}

		x = y
	}
}

// Sqrt returns the square root of a, correctly rounded to the nearest
// maxScale-digit value (exact ties are impossible: (s+0.5)^2 is never an
// integer). It errors if a < 0.
func (a decimal128) Sqrt() (decimal128, error) {
	if a.coef.IsZero() {
		return decimal128{}, nil
	}

	if a.neg {
		return decimal128{}, ErrSqrtNegative
	}

	// The result at scale 19 has coefficient round(sqrt(coef*10^(38-scale))).
	// The radicand n < 2^128 * 10^38 < 2^255 always fits in 256 bits, and
	// even the maximum input (2^128-1 at scale 0) roots to ~1.85e19, whose
	// scale-19 coefficient ~1.85e38 still fits in 128 bits.
	n := a.coef.MulFull(pow10[38-a.scale])
	s := isqrt256(n)

	// Round to nearest: sqrt(n) > s+0.5 iff n > s*(s+1).
	if sInc, _ := s.Add64(1); n.cmp(s.MulFull(sInc)) > 0 {
		s = sInc
	}

	return newDec(false, s, maxScale), nil
}

// fd is a small floating-decimal value used by powInt: value = coef *
// 10^exp, with coef kept to at most 38 significant digits.
type fd struct {
	coef u128
	exp  int
}

// ord returns the base-10 order of magnitude: 10^(ord-1) <= value < 10^ord.
func (x fd) ord() int {
	return x.exp + x.coef.decimalDigits()
}

// fdMul multiplies two floating-decimal values, rounding (half up) to at
// most 38 significant digits. Products that fit in 38 digits are exact.
func fdMul(x, y fd) fd {
	p := x.coef.MulFull(y.coef)
	exp := x.exp + y.exp

	// The digit-sum estimate can run one high, so reduce down to at most
	// 39 digits first (rounding on the last step) and cap to 38 after.
	drop := x.coef.decimalDigits() + y.coef.decimalDigits() - 39
	for drop > 0 {
		step := min(drop, int(maxScale))

		var r uint64
		p, r = p.quoRem64(pow10[step].lo)
		exp += step
		drop -= step

		if v := pow10[step].lo; drop == 0 && r >= v-v/2 {
			p = inc256(p)
		}
	}

	coef := u128{hi: p.hi, lo: p.lo}
	if !p.carry.IsZero() || coef.Cmp(pow10[38]) >= 0 {
		var r uint64
		p, r = p.quoRem64(10)
		exp++

		coef = u128{hi: p.hi, lo: p.lo}
		if r >= 5 {
			coef, _ = coef.Add64(1)
		}
	}

	return fd{coef: coef, exp: exp}
}

// fdInv returns 1/x to 38 significant digits (x > 0).
func fdInv(x fd) fd {
	// Normalize coef into [10^37, 10^38) so the quotient below always
	// carries 38 significant digits.
	k := 38 - x.coef.decimalDigits()
	coef := x.coef
	if k > 0 {
		coef, _ = coef.Mul(pow10[k])
	}

	num := pow10[38].MulFull(pow10[37]) // 10^75 as a u256
	q, _ := divRound(num, coef)

	return fd{coef: q, exp: -75 - (x.exp - k)}
}

// toDec converts a floating-decimal value to a decimal128 with sign neg,
// rounded half away from zero to the densest representable scale.
func (x fd) toDec(neg bool) (decimal128, error) {
	if x.exp >= 0 {
		if x.exp > 38 {
			return decimal128{}, ErrOverflow
		}

		coef, ok := x.coef.Mul(pow10[x.exp])
		if !ok {
			return decimal128{}, ErrOverflow
		}

		return newDec(neg, coef, 0), nil
	}

	sc := min(-x.exp, int(maxScale))
	drop := -x.exp - sc
	if drop > 38 {
		// Far below the smallest representable value: rounds to zero.
		return decimal128{}, nil
	}

	coef := x.coef
	for drop > 0 {
		step := min(drop, int(maxScale))

		var r uint64
		coef, r = coef.QuoRem64(pow10[step].lo)
		drop -= step

		if v := pow10[step].lo; drop == 0 && r >= v-v/2 {
			coef, _ = coef.Add64(1)
		}
	}

	//nolint:gosec // 0 <= sc <= 19
	return newDec(neg, coef, uint8(sc)), nil
}

// powInt computes a^n for a > 0 with integer exponent ±n (invert selects
// the negative sign), via binary exponentiation on 38-significant-digit
// floating-decimal intermediates: exact whenever the true result fits in
// 38 digits, accurate to ~1 ulp of the returned scale otherwise.
func (a decimal128) powInt(n u128, invert bool) (decimal128, error) {
	res := fd{coef: u128One}
	cur := fd{coef: a.coef, exp: -int(a.scale)}
	bitLen := n.bitLen()

	// outcome resolves an out-of-range magnitude: every pending factor
	// pushes the value further in the same direction (a > 1 only grows,
	// a < 1 only shrinks), so the final result is already decided.
	outcome := func(grew bool) (decimal128, error) {
		if grew != invert {
			return decimal128{}, ErrOverflow
		}

		return decimal128{}, nil
	}

	for i := range bitLen {
		if n.bitAt(i) == 1 {
			res = fdMul(res, cur)
			if o := res.ord(); o > 45 || o < -45 {
				return outcome(o > 45)
			}
		}

		if i+1 < bitLen {
			cur = fdMul(cur, cur)
			// The top bit of n is still pending, so cur alone decides.
			if o := cur.ord(); o > 90 || o < -90 {
				return outcome(o > 90)
			}
		}
	}

	if invert {
		res = fdInv(res)
	}

	return res.toDec(false)
}

// Pow returns a^e:
//   - a^0 = 1 (including 0^0), 0^e = 0 for e > 0, 0^e errors for e < 0.
//   - Integer exponents run through binary exponentiation on 38-digit
//     intermediates (exact whenever the result fits in 38 digits); for a
//     negative base the result's sign follows the exponent's parity.
//   - Fractional exponents are computed as exp(e * ln|a|) in fp120 and
//     require a > 0, otherwise ErrPowNegBase is returned.
//   - Results above decimal128's range return ErrOverflow; results below
//     half its smallest positive value round to 0.
func (a decimal128) Pow(e decimal128) (decimal128, error) {
	if e.coef.IsZero() {
		return decOne, nil
	}

	if a.coef.IsZero() {
		if e.neg {
			return decimal128{}, ErrDivideByZero
		}

		return decimal128{}, nil
	}

	base := a
	negResult := false
	et := e.trimTrailingZeros()

	if a.neg {
		if et.scale != 0 {
			return decimal128{}, ErrPowNegBase
		}

		negResult = et.coef.lo&1 == 1
		base = a.Abs()
	}

	if et.scale == 0 {
		res, err := base.powInt(et.coef, e.neg)
		if err != nil {
			return decimal128{}, err
		}

		if negResult {
			res = res.Neg()
		}

		return res, nil
	}

	lnMag, lnNeg := base.lnFp()
	if lnMag.IsZero() {
		// base == 1: the power is exactly 1.
		return decOne, nil
	}

	// x = e * ln|a| = lnMag * e.coef / 10^e.scale, rounded to fp120.
	q, rem := lnMag.MulFull(e.coef).quoRem64(pow10[e.scale].lo)
	if v := pow10[e.scale].lo; rem >= v-v/2 {
		q = inc256(q)
	}

	// |x| >= 128 always overflows (e^128 > 2^128) or rounds to zero.
	if !q.carry.IsZero() || q.hi>>63 != 0 {
		if lnNeg != e.neg {
			return decimal128{}, nil
		}

		return decimal128{}, ErrOverflow
	}

	// A negative base never reaches this path, so no sign to apply.
	return expFpToDec(u128{hi: q.hi, lo: q.lo}, lnNeg != e.neg)
}
