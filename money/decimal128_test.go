package money

import (
	"errors"
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
