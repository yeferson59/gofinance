package decimal

import (
	"math/big"
	"math/rand"
	"testing"
)

// randDecimalString generates a pseudo-random decimal string with up to
// intDigits integer digits and fracDigits fractional digits.
func randDecimalString(r *rand.Rand, intDigits, fracDigits int) string {
	var s []byte
	if r.Intn(2) == 0 {
		s = append(s, '-')
	}
	n := 1 + r.Intn(intDigits)
	for i := 0; i < n; i++ {
		if i == 0 {
			s = append(s, byte('1'+r.Intn(9)))
		} else {
			s = append(s, byte('0'+r.Intn(10)))
		}
	}
	if fracDigits > 0 && r.Intn(2) == 0 {
		f := 1 + r.Intn(fracDigits)
		s = append(s, '.')
		for i := 0; i < f; i++ {
			s = append(s, byte('0'+r.Intn(10)))
		}
	}
	return string(s)
}

func toRat(t *testing.T, s string) *big.Rat {
	t.Helper()
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		t.Fatalf("big.Rat couldn't parse %q", s)
	}
	return r
}

// ratString formats a big.Rat truncated (toward zero) to prec decimal
// digits, matching decimal128's truncating (non-rounding) division/mul
// semantics.
func ratString(r *big.Rat, prec int) string {
	neg := r.Sign() < 0
	abs := new(big.Rat).Abs(r)

	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(prec)), nil)
	scaled := new(big.Int).Mul(abs.Num(), scale)
	scaled.Quo(scaled, abs.Denom())

	s := scaled.String()
	for len(s) <= prec {
		s = "0" + s
	}
	intPart := s[:len(s)-prec]
	fracPart := s[len(s)-prec:]

	out := intPart
	i := len(fracPart)
	for i > 0 && fracPart[i-1] == '0' {
		i--
	}
	if i > 0 {
		out += "." + fracPart[:i]
	}
	if out == "0" {
		neg = false
	}
	if neg {
		out = "-" + out
	}
	return out
}

func TestDecimal128CrossCheckAddSubMul(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	for i := 0; i < 2000; i++ {
		as := randDecimalString(r, 8, 8)
		bs := randDecimalString(r, 8, 8)

		a, err := parseDecimal(as)
		if err != nil {
			t.Fatalf("parseDecimal(%q): %v", as, err)
		}
		b, err := parseDecimal(bs)
		if err != nil {
			t.Fatalf("parseDecimal(%q): %v", bs, err)
		}

		ra, rb := toRat(t, as), toRat(t, bs)
		addSubScale := int(a.scale)
		if int(b.scale) > addSubScale {
			addSubScale = int(b.scale)
		}

		if sum, err := a.Add(b); err != nil {
			t.Fatalf("Add(%q,%q): %v", as, bs, err)
		} else {
			want := new(big.Rat).Add(ra, rb)
			if got, wantStr := sum.String(), ratString(want, addSubScale); got != wantStr {
				t.Fatalf("%s + %s = %s, want %s", as, bs, got, wantStr)
			}
		}

		if diff, err := a.Sub(b); err != nil {
			t.Fatalf("Sub(%q,%q): %v", as, bs, err)
		} else {
			want := new(big.Rat).Sub(ra, rb)
			if got, wantStr := diff.String(), ratString(want, addSubScale); got != wantStr {
				t.Fatalf("%s - %s = %s, want %s", as, bs, got, wantStr)
			}
		}

		if prod, err := a.Mul(b); err != nil {
			t.Fatalf("Mul(%q,%q): %v", as, bs, err)
		} else {
			want := new(big.Rat).Mul(ra, rb)
			totalScale := int(a.scale) + int(b.scale)
			if totalScale > int(maxScale) {
				totalScale = int(maxScale)
			}
			if got, wantStr := prod.String(), ratString(want, totalScale); got != wantStr {
				t.Fatalf("%s * %s = %s, want %s (scale=%d)", as, bs, got, wantStr, totalScale)
			}
		}
	}
}

func TestDecimal128CrossCheckDiv(t *testing.T) {
	r := rand.New(rand.NewSource(7))

	for i := 0; i < 2000; i++ {
		as := randDecimalString(r, 6, 6)
		bs := randDecimalString(r, 6, 6)

		a, err := parseDecimal(as)
		if err != nil {
			t.Fatalf("parseDecimal(%q): %v", as, err)
		}
		b, err := parseDecimal(bs)
		if err != nil {
			t.Fatalf("parseDecimal(%q): %v", bs, err)
		}
		if b.IsZero() {
			continue
		}

		q, err := a.Div(b)
		if err != nil {
			// overflow is acceptable for pathological magnitude combinations;
			// it should never happen for these small test values though.
			t.Fatalf("Div(%q,%q): %v", as, bs, err)
		}

		ra, rb := toRat(t, as), toRat(t, bs)
		want := new(big.Rat).Quo(ra, rb)
		wantStr := ratString(want, int(maxScale))
		if got := q.String(); got != wantStr {
			t.Fatalf("%s / %s = %s, want %s", as, bs, got, wantStr)
		}
	}
}

func TestDecimal128CrossCheckRoundTrip(t *testing.T) {
	r := rand.New(rand.NewSource(99))

	for i := 0; i < 500; i++ {
		s := randDecimalString(r, 10, 15)
		d, err := parseDecimal(s)
		if err != nil {
			t.Fatalf("parseDecimal(%q): %v", s, err)
		}

		// String() must round-trip through parseDecimal to the same value.
		d2, err := parseDecimal(d.String())
		if err != nil {
			t.Fatalf("re-parse %q: %v", d.String(), err)
		}
		if d.Cmp(d2) != 0 {
			t.Fatalf("round-trip mismatch: %q -> %q -> %q", s, d.String(), d2.String())
		}

		want := toRat(t, s)
		got := toRat(t, d.String())
		if want.Cmp(got) != 0 {
			t.Fatalf("value mismatch: parsed %q as %q, big.Rat disagrees (want %s got %s)", s, d.String(), want.String(), got.String())
		}
	}
}

func TestDecimal128CrossCheckExample(t *testing.T) {
	// sanity check that ratString/formatting helpers agree with a known case
	a := mustParseDec(t, "12345.6789")
	b := mustParseDec(t, "0.0001")
	sum, err := a.Add(b)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := sum.String(), "12345.679"; got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}
