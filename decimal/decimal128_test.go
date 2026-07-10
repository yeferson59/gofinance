package decimal

import (
	"errors"
	"math"
	"testing"
)

func mustParseDec(t *testing.T, s string) decimal128 {
	t.Helper()
	d, err := parseDecimal(s)
	if err != nil {
		t.Fatalf("parseDecimal(%q): unexpected error: %v", s, err)
	}
	return d
}

func TestParseDecimalValid(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"0", "0"},
		{"-0", "0"},
		{"123", "123"},
		{"123.456", "123.456"},
		{"-123.456", "-123.456"},
		{"+123.456", "123.456"},
		{"0.001", "0.001"},
		{"1000000", "1000000"},
		{"0.0000000000000000001", "0.0000000000000000001"}, // 19 frac digits, exactly maxScale
	}

	for _, tt := range tests {
		d := mustParseDec(t, tt.input)
		if got := d.String(); got != tt.want {
			t.Errorf("parseDecimal(%q).String() = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseDecimalInvalid(t *testing.T) {
	tests := []string{
		"", "abc", ".", "-.", "+.", ".5", "-.5", "123.", "-123.",
		"1.2.3", "12a3", "1 2", "-", "+",
	}

	for _, tt := range tests {
		if _, err := parseDecimal(tt); err == nil {
			t.Errorf("parseDecimal(%q): expected error, got nil", tt)
		}
	}

	// more than maxScale (19) fractional digits
	if _, err := parseDecimal("0.00000000000000000001"); !errors.Is(err, ErrPrecOutOfRange) {
		t.Errorf("expected ErrPrecOutOfRange, got %v", err)
	}
}

func TestDecimal128AddSubDifferentScales(t *testing.T) {
	a := mustParseDec(t, "1.1")
	b := mustParseDec(t, "2.22")

	sum, err := a.Add(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := sum.String(); got != "3.32" {
		t.Errorf("1.1+2.22 = %s, want 3.32", got)
	}

	diff, err := a.Sub(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := diff.String(); got != "-1.12" {
		t.Errorf("1.1-2.22 = %s, want -1.12", got)
	}
}

func TestDecimal128AddSignCancellation(t *testing.T) {
	a := mustParseDec(t, "5")
	b := mustParseDec(t, "-5")
	sum, err := a.Add(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sum.IsZero() || sum.neg {
		t.Errorf("5+(-5) should be canonical zero, got %+v", sum)
	}
}

func TestDecimal128MulPrecisionCap(t *testing.T) {
	a := mustParseDec(t, "0.1") // scale 1
	b := mustParseDec(t, "0.1") // scale 1
	prod, err := a.Mul(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := prod.String(); got != "0.01" {
		t.Errorf("0.1*0.1 = %s, want 0.01", got)
	}

	// force scale sum > maxScale (19): two operands each with scale 15 => sum 30 > 19, truncated
	c := mustParseDec(t, "1.000000000000001") // scale 15
	d, err := c.Mul(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.scale != maxScale {
		t.Errorf("expected result scale capped at %d, got %d", maxScale, d.scale)
	}
}

func TestDecimal128DivPrecision(t *testing.T) {
	a := mustParseDec(t, "10")
	b := mustParseDec(t, "3")
	q, err := a.Div(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "3.3333333333333333333" // 19 threes, truncated not rounded
	if got := q.String(); got != want {
		t.Errorf("10/3 = %s, want %s", got, want)
	}

	if q.scale != maxScale {
		t.Errorf("expected scale=%d, got %d", maxScale, q.scale)
	}
}

func TestDecimal128DivByZero(t *testing.T) {
	a := mustParseDec(t, "1")
	if _, err := a.Div(decZero); !errors.Is(err, ErrDivideByZero) {
		t.Errorf("expected ErrDivideByZero, got %v", err)
	}
	if _, err := a.Div64(0); !errors.Is(err, ErrDivideByZero) {
		t.Errorf("expected ErrDivideByZero, got %v", err)
	}
}

func TestDecimal128Rounding(t *testing.T) {
	tests := []struct {
		name string
		fn   func(d decimal128, prec uint8) decimal128
		in   string
		prec uint8
		want string
	}{
		{"bank-down", decimal128.RoundBank, "1.12345", 4, "1.1234"},
		{"bank-tie-even-up", decimal128.RoundBank, "1.12335", 4, "1.1234"},
		{"bank-tie-half", decimal128.RoundBank, "1.5", 0, "2"},
		{"bank-tie-half-neg", decimal128.RoundBank, "-1.5", 0, "-2"},
		{"bank-tie-even-down", decimal128.RoundBank, "2.5", 0, "2"},
		{"away-up", decimal128.RoundAway, "1.12", 1, "1.2"},
		{"away-tie", decimal128.RoundAway, "1.15", 1, "1.2"},
		{"away-neg", decimal128.RoundAway, "-1.12", 1, "-1.2"},
		{"haz-up", decimal128.RoundHAZ, "1.12345", 4, "1.1235"},
		{"haz-down", decimal128.RoundHAZ, "1.12335", 4, "1.1234"},
		{"haz-tie", decimal128.RoundHAZ, "1.5", 0, "2"},
		{"haz-tie-neg", decimal128.RoundHAZ, "-1.5", 0, "-2"},
		{"htz-down", decimal128.RoundHTZ, "1.12345", 4, "1.1234"},
		{"htz-tie", decimal128.RoundHTZ, "1.5", 0, "1"},
		{"htz-tie-neg", decimal128.RoundHTZ, "-1.5", 0, "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := mustParseDec(t, tt.in)
			got := tt.fn(d, tt.prec).String()
			if got != tt.want {
				t.Errorf("%s(%s, %d) = %s, want %s", tt.name, tt.in, tt.prec, got, tt.want)
			}
		})
	}
}

func TestDecimal128FloorCeil(t *testing.T) {
	tests := []struct {
		in          string
		floor, ceil string
	}{
		{"1.5", "1", "2"},
		{"-1.5", "-2", "-1"},
		{"123.99", "123", "124"},
		{"123.01", "123", "124"},
		{"5", "5", "5"},
	}

	for _, tt := range tests {
		d := mustParseDec(t, tt.in)
		if got := d.Floor().String(); got != tt.floor {
			t.Errorf("Floor(%s) = %s, want %s", tt.in, got, tt.floor)
		}
		if got := d.Ceil().String(); got != tt.ceil {
			t.Errorf("Ceil(%s) = %s, want %s", tt.in, got, tt.ceil)
		}
	}
}

func TestDecimal128CmpAcrossScales(t *testing.T) {
	a := mustParseDec(t, "1.1")
	b := mustParseDec(t, "1.10")
	if a.Cmp(b) != 0 {
		t.Errorf("1.1 should equal 1.10")
	}

	c := mustParseDec(t, "1.09")
	if a.Cmp(c) <= 0 {
		t.Errorf("1.1 should be greater than 1.09")
	}

	neg := mustParseDec(t, "-1")
	pos := mustParseDec(t, "0.5")
	if neg.Cmp(pos) >= 0 {
		t.Errorf("-1 should be less than 0.5")
	}
}

func TestDecimal128StringFixed(t *testing.T) {
	tests := []struct {
		in   string
		prec uint8
		want string
	}{
		{"1.23", 4, "1.2300"},
		{"-1.23", 4, "-1.2300"},
		{"5", 2, "5.00"},
		{"5.123", 2, "5.123"}, // prec < stored precision: stays unchanged
		{"0", 5, "0.00000"},
	}

	for _, tt := range tests {
		d := mustParseDec(t, tt.in)
		if got := d.StringFixed(tt.prec); got != tt.want {
			t.Errorf("StringFixed(%s, %d) = %s, want %s", tt.in, tt.prec, got, tt.want)
		}
	}
}

func TestDecimal128TrailingZeroTrim(t *testing.T) {
	// 1.20 constructed via HiLo (coef=120, scale=2) should print as "1.2"
	d, err := decFromUint64(120, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := d.String(); got != "1.2" {
		t.Errorf("String() = %s, want 1.2", got)
	}
}

func TestDecimal128Int64(t *testing.T) {
	d := mustParseDec(t, "123.999")
	v, err := d.Int64()
	if err != nil || v != 123 {
		t.Errorf("Int64() = %d, %v, want 123, nil", v, err)
	}

	neg := mustParseDec(t, "-123.999")
	v2, err := neg.Int64()
	if err != nil || v2 != -123 {
		t.Errorf("Int64() = %d, %v, want -123, nil", v2, err)
	}
}

func TestDecimal128QuoRemMod(t *testing.T) {
	a := mustParseDec(t, "10")
	b := mustParseDec(t, "3")
	q, r, err := a.QuoRem(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.String() != "3" || r.String() != "1" {
		t.Errorf("QuoRem(10,3) = q=%s r=%s, want q=3 r=1", q.String(), r.String())
	}

	m, err := mustParseDec(t, "-10").Mod(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.String() != "-1" {
		t.Errorf("Mod(-10,3) = %s, want -1", m.String())
	}
}

func TestDecimal128OverflowOnParse(t *testing.T) {
	// 39-digit integer part overflows a 128-bit coefficient (max ~3.4x10^38).
	big39 := "999999999999999999999999999999999999999"
	if _, err := parseDecimal(big39); !errors.Is(err, ErrOverflow) {
		t.Errorf("expected ErrOverflow for %d-digit number, got %v", len(big39), err)
	}
}

func TestParseDecimalOverflowAtExactBoundary(t *testing.T) {
	// 2^128-1 is the largest representable u128 coefficient: it must parse.
	maxU128 := "340282366920938463463374607431768211455"
	if _, err := parseDecimal(maxU128); err != nil {
		t.Fatalf("expected max u128 to parse cleanly, got %v", err)
	}

	// One more (2^128) overflows exactly on the final digit's Add64, after
	// the preceding Mul64(10) itself stayed just within range.
	overBoundary := "340282366920938463463374607431768211456"
	if _, err := parseDecimal(overBoundary); !errors.Is(err, ErrOverflow) {
		t.Errorf("expected ErrOverflow at 2^128, got %v", err)
	}
}

func TestDecFromInt64MinInt64(t *testing.T) {
	// math.MinInt64 has no positive counterpart representable in int64, so
	// decFromInt64 special-cases it to avoid overflowing during negation.
	d, err := decFromInt64(math.MinInt64, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := d.String(); got != "-9223372036854775808" {
		t.Errorf("expected -9223372036854775808, got %s", got)
	}
}

func TestDecimal128NegZeroUnchanged(t *testing.T) {
	if got := decZero.Neg(); got != decZero {
		t.Errorf("expected Neg(0) to stay canonical zero, got %+v", got)
	}
}

func TestDecimal128CmpSignBranches(t *testing.T) {
	pos := mustParseDec(t, "5")
	neg := mustParseDec(t, "-3")

	if got := pos.Cmp(neg); got != 1 {
		t.Errorf("expected positive > negative to return 1, got %d", got)
	}
	if got := neg.Cmp(pos); got != -1 {
		t.Errorf("expected negative < positive to return -1, got %d", got)
	}
}

func TestDecimal128CmpMagnitudeScaleAlignmentOverflow(t *testing.T) {
	huge := decimal128{coef: u128{hi: ^uint64(0), lo: ^uint64(0)}, scale: 0}
	tiny := decimal128{coef: u128FromU64(1), scale: 1}

	// huge, once scaled up to tiny's scale, doesn't fit in 128 bits, but
	// it's obviously still larger than any 128-bit-coefficient value.
	if got := huge.cmpMagnitude(tiny); got != 1 {
		t.Errorf("expected 1 when scaling the larger operand overflows, got %d", got)
	}

	hugeLowScale := decimal128{coef: u128{hi: ^uint64(0), lo: ^uint64(0)}, scale: 0}
	tinyHighScale := decimal128{coef: u128FromU64(1), scale: 1}

	// mirror case: now a has the larger scale, so cmpMagnitude tries to
	// scale up b (the huge operand) instead, and that overflows too.
	if got := tinyHighScale.cmpMagnitude(hugeLowScale); got != -1 {
		t.Errorf("expected -1 when scaling the other operand overflows, got %d", got)
	}
}

func TestDecimal128AddScaleAlignmentOverflow(t *testing.T) {
	huge := decimal128{coef: u128{hi: ^uint64(0), lo: ^uint64(0)}, scale: 0}
	tiny := decimal128{coef: u128FromU64(1), scale: 1}

	if _, err := huge.Add(tiny); !errors.Is(err, ErrOverflow) {
		t.Errorf("expected ErrOverflow scaling a huge low-scale operand up, got %v", err)
	}
	if _, err := tiny.Add(huge); !errors.Is(err, ErrOverflow) {
		t.Errorf("expected ErrOverflow scaling a huge high-scale operand up, got %v", err)
	}
}

func TestDecimal128MulZeroOperand(t *testing.T) {
	five := mustParseDec(t, "5")

	prod, err := five.Mul(decZero)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !prod.IsZero() {
		t.Errorf("expected 5*0 to be zero, got %s", prod.String())
	}

	prod2, err := decZero.Mul(five)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !prod2.IsZero() {
		t.Errorf("expected 0*5 to be zero, got %s", prod2.String())
	}
}

func TestDecimal128MulOverflowSurvivesScaleReduction(t *testing.T) {
	// Both operands near the 128-bit ceiling, each at scale 19: the raw
	// 256-bit product, even after dividing out the excess scale (19), is
	// still far too large to fit in 128 bits.
	maxCoef := u128{hi: ^uint64(0), lo: ^uint64(0)}
	a := decimal128{coef: maxCoef, scale: maxScale}
	b := decimal128{coef: maxCoef, scale: maxScale}

	if _, err := a.Mul(b); !errors.Is(err, ErrOverflow) {
		t.Errorf("expected ErrOverflow, got %v", err)
	}
}

func TestDecimal128DivZeroDividend(t *testing.T) {
	q, err := decZero.Div(mustParseDec(t, "5"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !q.IsZero() {
		t.Errorf("expected 0/5 to be zero, got %s", q.String())
	}
}

func TestDecimal128Div64ZeroDividend(t *testing.T) {
	q, err := decZero.Div64(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !q.IsZero() {
		t.Errorf("expected 0/5 to be zero, got %s", q.String())
	}
}

func TestDecimal128Div64Overflow(t *testing.T) {
	huge := decimal128{coef: u128{hi: ^uint64(0), lo: ^uint64(0)}, scale: 0}
	if _, err := huge.Div64(1); !errors.Is(err, ErrOverflow) {
		t.Errorf("expected ErrOverflow, got %v", err)
	}
}

func TestDecimal128QuoRemZeroDividend(t *testing.T) {
	q, r, err := decZero.QuoRem(mustParseDec(t, "5"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !q.IsZero() || !r.IsZero() {
		t.Errorf("expected q=0 r=0, got q=%s r=%s", q.String(), r.String())
	}
}

func TestDecimal128QuoRemUsesLargerScale(t *testing.T) {
	// b's scale (1) exceeds a's scale (0), so QuoRem must align to it.
	a := mustParseDec(t, "10")
	b := mustParseDec(t, "3.5")

	q, r, err := a.QuoRem(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.String() != "2" || r.String() != "3" {
		t.Errorf("QuoRem(10, 3.5) = q=%s r=%s, want q=2 r=3", q.String(), r.String())
	}
}

func TestDecimal128QuoRemScaleAlignmentOverflow(t *testing.T) {
	huge := decimal128{coef: u128{hi: ^uint64(0), lo: ^uint64(0)}, scale: 0}
	scaled := decimal128{coef: u128FromU64(1), scale: 1}

	if _, _, err := huge.QuoRem(scaled); !errors.Is(err, ErrOverflow) {
		t.Errorf("expected ErrOverflow scaling up a's huge coefficient, got %v", err)
	}
	if _, _, err := scaled.QuoRem(huge); !errors.Is(err, ErrOverflow) {
		t.Errorf("expected ErrOverflow scaling up b's huge coefficient, got %v", err)
	}
}

func TestDecimal128RoundingNoOpWhenPrecGreaterOrEqualScale(t *testing.T) {
	d := mustParseDec(t, "1.23")

	for _, tt := range []struct {
		name string
		fn   func(decimal128, uint8) decimal128
	}{
		{"RoundAway", decimal128.RoundAway},
		{"RoundHAZ", decimal128.RoundHAZ},
		{"RoundHTZ", decimal128.RoundHTZ},
		{"Trunc", decimal128.Trunc},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(d, 5).String(); got != "1.23" {
				t.Errorf("%s(1.23, 5) = %s, want unchanged 1.23", tt.name, got)
			}
			if got := tt.fn(d, 2).String(); got != "1.23" {
				t.Errorf("%s(1.23, 2) = %s, want unchanged 1.23", tt.name, got)
			}
		})
	}
}

func TestDecimal128RoundHTZAwayFromZeroPastHalf(t *testing.T) {
	// remainder (0.99) is strictly greater than half (0.5): even
	// half-toward-zero rounds away from zero here.
	d := mustParseDec(t, "1.99")
	if got := d.RoundHTZ(0).String(); got != "2" {
		t.Errorf("RoundHTZ(1.99, 0) = %s, want 2", got)
	}
}

func TestDecimal128Int64Overflow(t *testing.T) {
	huge := decimal128{coef: u128{hi: 1, lo: 0}, scale: 0} // >= 2^64, doesn't fit int64
	if _, err := huge.Int64(); !errors.Is(err, ErrIntPartOverflow) {
		t.Errorf("expected ErrIntPartOverflow, got %v", err)
	}

	tooBig := decimal128{coef: u128FromU64(uint64(math.MaxInt64) + 1), scale: 0}
	if _, err := tooBig.Int64(); !errors.Is(err, ErrIntPartOverflow) {
		t.Errorf("expected ErrIntPartOverflow, got %v", err)
	}
}

func TestDecimal128StringFixedClampsExcessivePrecision(t *testing.T) {
	d := mustParseDec(t, "1.5")
	if got := d.StringFixed(200); got != d.StringFixed(maxScale) {
		t.Errorf("StringFixed(200) = %s, want same as StringFixed(maxScale) = %s", got, d.StringFixed(maxScale))
	}
}

func TestDecimal128RescaleOverflowKeepsOriginal(t *testing.T) {
	// A coefficient already near the 128-bit ceiling can't be scaled up
	// any further: rescale must fall back to the untouched (trimmed)
	// value rather than panicking or silently corrupting it.
	huge := decimal128{coef: u128{hi: ^uint64(0), lo: ^uint64(0)}, scale: 0}

	got := huge.StringFixed(maxScale)
	if got != huge.String() {
		t.Errorf("StringFixed on an unscalable huge value = %s, want unchanged %s", got, huge.String())
	}
}
