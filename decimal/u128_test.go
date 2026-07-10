package decimal

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestU128AddOverflow(t *testing.T) {
	maxU128 := u128{hi: ^uint64(0), lo: ^uint64(0)}
	if _, ok := maxU128.Add(u128One); ok {
		t.Fatal("expected overflow")
	}

	sum, ok := u128FromU64(1).Add(u128FromU64(2))
	if !ok || sum != u128FromU64(3) {
		t.Fatalf("expected 3, got %+v (ok=%v)", sum, ok)
	}
}

func TestU128SubBorrow(t *testing.T) {
	if _, ok := u128FromU64(1).Sub(u128FromU64(2)); ok {
		t.Fatal("expected borrow/false")
	}

	diff, ok := u128FromU64(5).Sub(u128FromU64(3))
	if !ok || diff != u128FromU64(2) {
		t.Fatalf("expected 2, got %+v", diff)
	}
}

func TestU128Mul64Overflow(t *testing.T) {
	maxU128 := u128{hi: ^uint64(0), lo: ^uint64(0)}
	if _, ok := maxU128.Mul64(2); ok {
		t.Fatal("expected overflow")
	}

	p, ok := u128FromU64(1000).Mul64(1000)
	if !ok || p != u128FromU64(1_000_000) {
		t.Fatalf("expected 1000000, got %+v", p)
	}
}

func TestU128MulFull(t *testing.T) {
	// 2^64 * 2^64 = 2^128, i.e. carry=1, hi=0, lo=0
	a := u128{lo: 0, hi: 1} // 2^64
	prod := a.MulFull(a)
	if prod.carry != u128One || prod.hi != 0 || prod.lo != 0 {
		t.Fatalf("expected carry=1,hi=0,lo=0, got %+v", prod)
	}

	// simple case within 128 bits
	b := u128FromU64(123456789)
	prod2 := b.MulFull(b)
	want := uint64(123456789) * uint64(123456789)
	if !prod2.carry.IsZero() || prod2.hi != 0 || prod2.lo != want {
		t.Fatalf("expected lo=%d, got %+v", want, prod2)
	}
}

func TestU128QuoRem64(t *testing.T) {
	u := u128FromU64(1_000_000_007)
	q, r := u.QuoRem64(97)
	wantQ, wantR := uint64(1_000_000_007)/97, uint64(1_000_000_007)%97
	if q.lo != wantQ || q.hi != 0 || r != wantR {
		t.Fatalf("expected q=%d r=%d, got q=%+v r=%d", wantQ, wantR, q, r)
	}
}

func TestU128String(t *testing.T) {
	tests := []struct {
		u    u128
		want string
	}{
		{u128{}, "0"},
		{u128FromU64(0), "0"},
		{u128FromU64(123), "123"},
		{u128FromU64(18446744073709551615), "18446744073709551615"}, // max uint64
		{pow10[19], "10000000000000000000"},
		{pow10[38], "100000000000000000000000000000000000000"},
	}

	for _, tt := range tests {
		if got := tt.u.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestQuoRem128by128(t *testing.T) {
	// divisor requiring the full 128 bits: v = 2^70 (hi != 0)
	v := u128{hi: 1 << 6, lo: 0}
	u := v // u == v, quotient must be exactly 1
	q, r, ok := quoRem128by128(u, v)
	if !ok || q != u128One || !r.IsZero() {
		t.Fatalf("expected q=1 r=0, got q=%+v r=%+v ok=%v", q, r, ok)
	}

	// u = 3*v + 5
	five := u128FromU64(5)
	threeV, _ := v.Mul64(3)
	u2, _ := threeV.Add(five)
	q2, r2, ok2 := quoRem128by128(u2, v)
	if !ok2 || q2 != u128FromU64(3) || r2 != five {
		t.Fatalf("expected q=3 r=5, got q=%+v r=%+v", q2, r2)
	}

	if _, _, ok3 := quoRem128by128(u128One, u128Zero); ok3 {
		t.Fatal("expected division by zero to fail")
	}
}

func TestU128Cmp64(t *testing.T) {
	tests := []struct {
		name string
		u    u128
		v    uint64
		want int
	}{
		{"hi nonzero always greater", u128{hi: 1, lo: 0}, ^uint64(0), 1},
		{"lo less", u128FromU64(5), 10, -1},
		{"lo greater", u128FromU64(10), 5, 1},
		{"equal", u128FromU64(42), 42, 0},
		{"zero vs zero", u128{}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Cmp64(tt.v); got != tt.want {
				t.Errorf("Cmp64(%+v, %d) = %d, want %d", tt.u, tt.v, got, tt.want)
			}
		})
	}
}

func TestU128Add64Overflow(t *testing.T) {
	maxU128 := u128{hi: ^uint64(0), lo: ^uint64(0)}
	if _, ok := maxU128.Add64(1); ok {
		t.Fatal("expected overflow")
	}

	// carry into hi that itself doesn't overflow hi.
	u := u128{hi: 1, lo: ^uint64(0)}
	sum, ok := u.Add64(1)
	if !ok || sum != (u128{hi: 2, lo: 0}) {
		t.Fatalf("expected {hi:2,lo:0}, got %+v (ok=%v)", sum, ok)
	}
}

func TestU128Mul(t *testing.T) {
	// both operands >= 2^64: guaranteed overflow.
	a := u128{hi: 1, lo: 0}
	b := u128{hi: 1, lo: 0}
	if _, ok := a.Mul(b); ok {
		t.Fatal("expected overflow when both operands have hi != 0")
	}

	// v.hi == 0: delegates to u.Mul64(v.lo).
	u := u128FromU64(1000)
	v := u128FromU64(2000)
	prod, ok := u.Mul(v)
	if !ok || prod != u128FromU64(2_000_000) {
		t.Fatalf("expected 2000000, got %+v (ok=%v)", prod, ok)
	}

	// u.hi == 0, v.hi != 0: delegates to v.Mul64(u.lo).
	u2 := u128FromU64(5)
	v2 := u128{hi: 1, lo: 0} // 2^64
	prod2, ok2 := u2.Mul(v2)
	want2, _ := v2.Mul64(5)
	if !ok2 || prod2 != want2 {
		t.Fatalf("expected %+v, got %+v (ok=%v)", want2, prod2, ok2)
	}

	// overflow surfaces even through the v.hi==0 delegation path.
	maxU128 := u128{hi: ^uint64(0), lo: ^uint64(0)}
	if _, ok := maxU128.Mul(u128FromU64(2)); ok {
		t.Fatal("expected overflow")
	}
}

func TestU128BitLenSmallDividend(t *testing.T) {
	// Exercises u128.bitLen()'s u.hi==0 branch: only reachable through
	// quoRem128by128's general (v.hi!=0) path when the dividend itself
	// fits in 64 bits.
	u := u128FromU64(1000)
	v := u128{hi: 1, lo: 0} // 2^64, dividend is smaller: quotient 0, remainder u

	q, r, ok := quoRem128by128(u, v)
	if !ok || !q.IsZero() || r != u {
		t.Fatalf("expected q=0 r=%+v, got q=%+v r=%+v ok=%v", u, q, r, ok)
	}
}

// TestQuoRem128By128FullDivisorCrossCheck hammers the div3by2-based path
// (v.hi != 0) with full-range random operands, cross-checking against
// math/big. This covers the normalization shifts (including s == 0), the
// qhat correction loop, and the a2 == v1 estimate branch.
func TestQuoRem128By128FullDivisorCrossCheck(t *testing.T) {
	r := rand.New(rand.NewSource(5))

	for range 20000 {
		u := u128{hi: r.Uint64(), lo: r.Uint64()}
		v := u128{hi: r.Uint64(), lo: r.Uint64()}
		for v.hi == 0 {
			v.hi = r.Uint64()
		}

		q, rem, ok := quoRem128by128(u, v)
		if !ok {
			t.Fatalf("unexpected !ok for u=%+v v=%+v", u, v)
		}

		wantQ, wantR := new(big.Int).QuoRem(u128ToBig(u), u128ToBig(v), new(big.Int))
		if u128ToBig(q).Cmp(wantQ) != 0 || u128ToBig(rem).Cmp(wantR) != 0 {
			t.Fatalf("u=%+v / v=%+v: got q=%+v r=%+v, want q=%s r=%s", u, v, q, rem, wantQ, wantR)
		}
	}
}

func TestU256QuoRem128(t *testing.T) {
	// (10^38) / (10^19) == 10^19, remainder 0
	prod := pow10[19].MulFull(pow10[19]) // 10^19 * 10^19 = 10^38
	q, r, ok := prod.QuoRem128(pow10[19])
	if !ok {
		t.Fatal("expected success")
	}
	if q != pow10[19] || !r.IsZero() {
		t.Fatalf("expected q=10^19 r=0, got q=%s r=%s", q.String(), r.String())
	}

	// division by a divisor requiring full 128 bits (general binary long division path)
	v := u128{hi: 1 << 10, lo: 12345}
	prod2 := v.MulFull(u128FromU64(999999))
	q2, r2, ok2 := prod2.QuoRem128(v)
	if !ok2 || q2 != u128FromU64(999999) || !r2.IsZero() {
		t.Fatalf("expected q=999999 r=0, got q=%s r=%s ok=%v", q2.String(), r2.String(), ok2)
	}

	// overflow: quotient doesn't fit in 128 bits (dividing by 1 leaves the full 256-bit value)
	huge := pow10[38].MulFull(pow10[38])
	if _, _, ok3 := huge.QuoRem128(u128One); ok3 {
		t.Fatal("expected overflow (quotient too large for 128 bits)")
	}
}
