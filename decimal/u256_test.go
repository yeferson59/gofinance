package decimal

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestU256IsZero(t *testing.T) {
	tests := []struct {
		name string
		x    u256
		want bool
	}{
		{"zero value", u256{}, true},
		{"lo set", u256{lo: 1}, false},
		{"hi set", u256{hi: 1}, false},
		{"carry.lo set", u256{carry: u128{lo: 1}}, false},
		{"carry.hi set", u256{carry: u128{hi: 1}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestU256BitLen(t *testing.T) {
	tests := []struct {
		name string
		x    u256
		want int
	}{
		{"zero", u256{}, 0},
		{"lo only", u256{lo: 5}, 3}, // 0b101
		{"lo full", u256{lo: ^uint64(0)}, 64},
		{"hi only", u256{hi: 1}, 65}, // 64 + 1
		{"hi full", u256{hi: ^uint64(0)}, 128},
		{"carry.lo only", u256{carry: u128{lo: 1}}, 129}, // 128 + 1
		{"carry.lo full", u256{carry: u128{lo: ^uint64(0)}}, 192},
		{"carry.hi only", u256{carry: u128{hi: 1}}, 193}, // 192 + 1
		{"carry.hi full", u256{carry: u128{hi: ^uint64(0)}}, 256},
		// carry.hi takes priority over carry.lo/hi/lo even when all are set
		{"all set", u256{hi: ^uint64(0), lo: ^uint64(0), carry: u128{hi: 1, lo: ^uint64(0)}}, 193},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.bitLen(); got != tt.want {
				t.Errorf("bitLen() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestU256BitAt(t *testing.T) {
	x := u256{
		lo:    1 << 3,                        // bit 3
		hi:    1 << 5,                        // bit 64+5 = 69
		carry: u128{lo: 1 << 7, hi: 1 << 11}, // bits 128+7=135, 192+11=203
	}

	tests := []struct {
		i    int
		want uint64
	}{
		{0, 0},
		{3, 1},
		{63, 0},
		{64, 0},
		{69, 1},
		{127, 0},
		{128, 0},
		{135, 1},
		{191, 0},
		{192, 0},
		{203, 1},
		{255, 0},
	}

	for _, tt := range tests {
		if got := x.bitAt(tt.i); got != tt.want {
			t.Errorf("bitAt(%d) = %d, want %d", tt.i, got, tt.want)
		}
	}
}

func TestU256QuoRem128ZeroDivisor(t *testing.T) {
	x := u256{lo: 100}
	if _, _, ok := x.QuoRem128(u128Zero); ok {
		t.Fatal("expected division by zero to fail")
	}
}

func TestU256QuoRem128DelegatesWhenCarryZero(t *testing.T) {
	// carry == 0: x fits in 128 bits, QuoRem128 must agree exactly with
	// quoRem128by128 on the low 128 bits, for both a 64-bit and a full
	// 128-bit divisor.
	x := u256{hi: 0, lo: 1000}
	v := u128FromU64(7)

	q, r, ok := x.QuoRem128(v)
	wantQ, wantR, wantOK := quoRem128by128(u128{lo: 1000}, v)
	if ok != wantOK || q != wantQ || r != wantR {
		t.Fatalf("QuoRem128 = q=%+v r=%+v ok=%v, want q=%+v r=%+v ok=%v", q, r, ok, wantQ, wantR, wantOK)
	}
	if !ok || q != u128FromU64(142) || r != u128FromU64(6) {
		t.Fatalf("1000/7 = q=%+v r=%+v, want q=142 r=6", q, r)
	}
}

// u128ToBig converts a u128 to a big.Int for cross-checking.
func u128ToBig(u u128) *big.Int {
	b := new(big.Int).Lsh(new(big.Int).SetUint64(u.hi), 64)
	return b.Or(b, new(big.Int).SetUint64(u.lo))
}

// u128FromBig converts a non-negative big.Int (< 2^128) to a u128.
func u128FromBig(t *testing.T, b *big.Int) u128 {
	t.Helper()
	if b.Sign() < 0 || b.BitLen() > 128 {
		t.Fatalf("u128FromBig: %s doesn't fit in 128 bits", b.String())
	}
	mask := new(big.Int).SetUint64(^uint64(0))
	lo := new(big.Int).And(b, mask).Uint64()
	hi := new(big.Int).And(new(big.Int).Rsh(b, 64), mask).Uint64()
	return u128{hi: hi, lo: lo}
}

// u256FromBig converts a non-negative big.Int (< 2^256) to a u256.
func u256FromBig(t *testing.T, b *big.Int) u256 {
	t.Helper()
	if b.Sign() < 0 || b.BitLen() > 256 {
		t.Fatalf("u256FromBig: %s doesn't fit in 256 bits", b.String())
	}
	mask := new(big.Int).SetUint64(^uint64(0))
	lo := new(big.Int).And(b, mask).Uint64()
	rest := new(big.Int).Rsh(b, 64)
	hi := new(big.Int).And(rest, mask).Uint64()
	rest = new(big.Int).Rsh(rest, 64)
	carryLo := new(big.Int).And(rest, mask).Uint64()
	carryHi := new(big.Int).Rsh(rest, 64).Uint64()
	return u256{hi: hi, lo: lo, carry: u128{hi: carryHi, lo: carryLo}}
}

// TestU256QuoRem128FastPath exercises the v.hi==0 && carry.hi==0 fast
// division path (bits.Div64-based) with random q*v+r reconstructions,
// where q spans the full 128 bits and v is a random nonzero uint64. This
// guarantees carry.hi==0 and carry.lo<v.lo, per the invariant that
// q<2^128 and v<2^64 implies q*v+r < v*2^128.
func TestU256QuoRem128FastPath(t *testing.T) {
	r := rand.New(rand.NewSource(1))

	for i := 0; i < 2000; i++ {
		vLo := r.Uint64()
		if vLo == 0 {
			vLo = 1
		}
		v := u128FromU64(vLo)

		q := u128{hi: r.Uint64(), lo: r.Uint64()}
		rem := r.Uint64() % vLo

		qBig := u128ToBig(q)
		vBig := u128ToBig(v)
		xBig := new(big.Int).Mul(qBig, vBig)
		xBig.Add(xBig, new(big.Int).SetUint64(rem))

		x := u256FromBig(t, xBig)
		if x.carry.hi != 0 {
			t.Fatalf("precondition violated: carry.hi = %d, want 0", x.carry.hi)
		}
		if x.carry.lo >= v.lo {
			t.Fatalf("precondition violated: carry.lo=%d >= v.lo=%d", x.carry.lo, v.lo)
		}

		gotQ, gotR, ok := x.QuoRem128(v)
		if !ok {
			t.Fatalf("case %d: expected ok=true for q=%+v v=%+v rem=%d", i, q, v, rem)
		}
		if gotQ != q {
			t.Fatalf("case %d: quotient = %+v, want %+v (v=%+v rem=%d)", i, gotQ, q, v, rem)
		}
		if gotR != (u128{lo: rem}) {
			t.Fatalf("case %d: remainder = %+v, want {lo:%d}", i, gotR, rem)
		}
	}
}

// TestU256QuoRem128GeneralPath exercises the binaryDivU128 path taken when
// v.hi != 0 (a full 128-bit divisor), again reconstructing x = q*v+r so
// the expected result is known exactly.
func TestU256QuoRem128GeneralPath(t *testing.T) {
	r := rand.New(rand.NewSource(2))

	for i := 0; i < 2000; i++ {
		v := u128{hi: r.Uint64()%1000 + 1, lo: r.Uint64()} // v.hi != 0
		q := u128{hi: r.Uint64() % 1000, lo: r.Uint64()}   // keep q small enough that q*v's carry stays below v

		vBig := u128ToBig(v)
		remBig := new(big.Int).Mod(big.NewInt(int64(r.Uint32())), vBig)
		xBig := new(big.Int).Mul(u128ToBig(q), vBig)
		xBig.Add(xBig, remBig)

		if xBig.BitLen() > 256 {
			continue
		}
		x := u256FromBig(t, xBig)

		gotQ, gotR, ok := x.QuoRem128(v)
		if !ok {
			t.Fatalf("case %d: expected ok=true for q=%+v v=%+v", i, q, v)
		}
		if gotQ != q {
			t.Fatalf("case %d: quotient = %+v, want %+v", i, gotQ, q)
		}
		wantR := u128FromBig(t, remBig)
		if gotR != wantR {
			t.Fatalf("case %d: remainder = %+v, want %+v", i, gotR, wantR)
		}
	}
}

func TestU256QuoRem128OverflowQuotientTooBig(t *testing.T) {
	// x = 2^255 (needs carry.hi's top bit), v = 1: true quotient is 2^255,
	// which doesn't fit in 128 bits.
	x := u256{carry: u128{hi: 1 << 63}}
	if _, _, ok := x.QuoRem128(u128One); ok {
		t.Fatal("expected overflow (quotient doesn't fit in 128 bits)")
	}

	// carry >= v with v.hi != 0 also must overflow: quotient can't fit.
	v := u128{hi: 1, lo: 0}
	x2 := u256{carry: u128{hi: 1, lo: 0}} // carry == v
	if _, _, ok := x2.QuoRem128(v); ok {
		t.Fatal("expected overflow when carry >= v")
	}
}

func TestU256QuoRem128FastPathOverflow(t *testing.T) {
	// v.hi==0 and x.carry.hi==0 (fast-path shape), but x.carry.lo >= v.lo:
	// the quotient would need more than 128 bits, so it must fail even
	// though the divisor is only 64 bits wide.
	v := u128FromU64(5)
	x := u256{carry: u128{lo: 10}}
	if _, _, ok := x.QuoRem128(v); ok {
		t.Fatal("expected overflow when carry.lo >= v.lo")
	}
}

func TestU256QuoRem128ExactSquareRoot(t *testing.T) {
	// (10^38) / (10^19) == 10^19 exactly, remainder 0 — a concrete sanity
	// check alongside the randomized cross-checks above.
	prod := pow10[19].MulFull(pow10[19])
	q, r, ok := prod.QuoRem128(pow10[19])
	if !ok || q != pow10[19] || !r.IsZero() {
		t.Fatalf("10^38 / 10^19 = q=%s r=%s ok=%v, want q=%s r=0 ok=true", q.String(), r.String(), ok, pow10[19].String())
	}
}
