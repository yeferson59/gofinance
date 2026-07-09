package money

import "testing"

func TestU128AddOverflow(t *testing.T) {
	max := u128{hi: ^uint64(0), lo: ^uint64(0)}
	if _, ok := max.Add(u128One); ok {
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
	max := u128{hi: ^uint64(0), lo: ^uint64(0)}
	if _, ok := max.Mul64(2); ok {
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
