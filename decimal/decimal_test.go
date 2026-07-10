package decimal

import (
	"errors"
	"math"
	"testing"
)

func TestNewFromFloat64(t *testing.T) {
	d, err := NewFromFloat64(123.456)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if d.String() != "123.456" {
		t.Errorf("expected 123.456, got %s", d.String())
	}
}

func TestDecimalPow(t *testing.T) {
	tests := []struct {
		base     float64
		exponent float64
		expected string
	}{
		{2, 10, "1024"},
		{1.005, 12, "1.0616778118644976"},
		{4, 0.5, "2"},
		{5, 0, "1"},
	}

	for _, tt := range tests {
		base := MustFromFloat64(tt.base)
		exponent := MustFromFloat64(tt.exponent)

		result, err := base.Pow(exponent)
		if err != nil {
			t.Fatalf("unexpected error for %v^%v: %v", tt.base, tt.exponent, err)
		}
		if result.String() != tt.expected {
			t.Errorf("expected %v^%v = %s, got %s", tt.base, tt.exponent, tt.expected, result.String())
		}
	}
}

func TestDecimalPowInvalidResult(t *testing.T) {
	base := MustFromFloat64(-2)
	exponent := MustFromFloat64(0.5)

	if _, err := base.Pow(exponent); err == nil {
		t.Error("expected error for negative base with fractional exponent")
	}
}

func TestDecimalMustPow(t *testing.T) {
	result := MustFromFloat64(2).MustPow(MustFromFloat64(8))
	if result.String() != "256" {
		t.Errorf("expected 256, got %s", result.String())
	}
}

func TestDecimalMustPowPanicsOnInvalidResult(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for negative base with fractional exponent")
		}
	}()

	MustFromFloat64(-2).MustPow(MustFromFloat64(0.5))
}

func TestDecimalLn(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{1, "0"},
		{math.E, "1"},
		{2, "0.6931471805599453"},
	}

	for _, tt := range tests {
		result, err := MustFromFloat64(tt.input).Ln()
		if err != nil {
			t.Fatalf("unexpected error for Ln(%v): %v", tt.input, err)
		}
		if result.String() != tt.expected {
			t.Errorf("expected Ln(%v) = %s, got %s", tt.input, tt.expected, result.String())
		}
	}
}

func TestDecimalLnInvalidResult(t *testing.T) {
	if _, err := MustFromFloat64(0).Ln(); err == nil {
		t.Error("expected error for Ln(0)")
	}
	if _, err := MustFromFloat64(-1).Ln(); err == nil {
		t.Error("expected error for Ln(-1)")
	}
}

func TestDecimalLog10(t *testing.T) {
	result, err := MustFromFloat64(1000).Log10()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.String() != "3" {
		t.Errorf("expected 3, got %s", result.String())
	}
}

func TestDecimalLog2(t *testing.T) {
	result, err := MustFromFloat64(8).Log2()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.String() != "3" {
		t.Errorf("expected 3, got %s", result.String())
	}
}

func TestDecimalLog(t *testing.T) {
	tests := []struct {
		x        float64
		base     float64
		expected string
	}{
		{8, 2, "3"},
		{100, 10, "2"},
	}

	for _, tt := range tests {
		result, err := MustFromFloat64(tt.x).Log(MustFromFloat64(tt.base))
		if err != nil {
			t.Fatalf("unexpected error for Log(%v, base %v): %v", tt.x, tt.base, err)
		}
		if result.String() != tt.expected {
			t.Errorf("expected Log(%v, base %v) = %s, got %s", tt.x, tt.base, tt.expected, result.String())
		}
	}
}

func TestDecimalLogInvalidBase(t *testing.T) {
	if _, err := MustFromFloat64(10).Log(MustFromFloat64(1)); !errors.Is(err, ErrDivideByZero) {
		t.Errorf("expected ErrDivideByZero for base 1, got %v", err)
	}
}

func TestDecimalMustLnPanicsOnInvalidResult(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for Ln(0)")
		}
	}()

	MustFromFloat64(0).MustLn()
}

func TestNewFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"123.456", "123.456", false},
		{"-123.456", "-123.456", false},
		{"0.001", "0.001", false},
		{"1000000", "1000000", false},
		{"", "", true},
		{"abc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d, err := NewFromString(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if d.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, d.String())
			}
		})
	}
}

func TestMustFromString(t *testing.T) {
	d := MustFromString("99.99")
	if d.String() != "99.99" {
		t.Errorf("expected 99.99, got %s", d.String())
	}
}

func TestDecimalRoundBank(t *testing.T) {
	d := MustFromString("1.235")
	rounded := d.RoundBank(2)
	if rounded.String() != "1.24" {
		t.Errorf("expected 1.24, got %s", rounded.String())
	}
}

func TestDecimalTrunc(t *testing.T) {
	d := MustFromString("1.999")
	trunc := d.Trunc(1)
	if trunc.String() != "1.9" {
		t.Errorf("expected 1.9, got %s", trunc.String())
	}
}

func TestDecimalAbs(t *testing.T) {
	d := MustFromString("-123.45")
	abs := d.Abs()
	if abs.String() != "123.45" {
		t.Errorf("expected 123.45, got %s", abs.String())
	}
}

func TestDecimalFloor(t *testing.T) {
	d := MustFromString("123.99")
	f := d.Floor()
	if f.String() != "123" {
		t.Errorf("expected 123, got %s", f.String())
	}
}

func TestDecimalCeil(t *testing.T) {
	d := MustFromString("123.01")
	c := d.Ceil()
	if c.String() != "124" {
		t.Errorf("expected 124, got %s", c.String())
	}
}

func TestDecimalNeg(t *testing.T) {
	d := MustFromString("123.45")
	neg := d.Neg()
	if neg.String() != "-123.45" {
		t.Errorf("expected -123.45, got %s", neg.String())
	}
}

func TestDecimalStringFixed(t *testing.T) {
	d := MustFromString("1.2")
	if d.StringFixed(2) != "1.20" {
		t.Errorf("expected 1.20, got %s", d.StringFixed(2))
	}
}

func TestDecimalJSON(t *testing.T) {
	d := MustFromString("123.45")
	data, err := d.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var d2 Decimal
	if err := d2.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !d.Equal(d2) {
		t.Errorf("expected %v, got %v", d, d2)
	}
}

func TestDecimalJSONNumber(t *testing.T) {
	var d Decimal
	if err := d.UnmarshalJSON([]byte("123.45")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.String() != "123.45" {
		t.Errorf("expected 123.45, got %s", d.String())
	}
}

func TestDecimalJSONInvalid(t *testing.T) {
	var d Decimal
	if err := d.UnmarshalJSON([]byte("true")); err == nil {
		t.Error("expected error unmarshalling a JSON boolean")
	}
	if err := d.UnmarshalJSON([]byte(`"abc"`)); err == nil {
		t.Error("expected error unmarshalling an invalid decimal string")
	}
	if err := d.UnmarshalJSON([]byte(`not-json`)); err == nil {
		t.Error("expected error unmarshalling malformed JSON")
	}
}

func TestNewFromInt64(t *testing.T) {
	d, err := NewFromInt64(-12345, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.String() != "-123.45" {
		t.Errorf("expected -123.45, got %s", d.String())
	}

	if _, err := NewFromInt64(1, 20); !errors.Is(err, ErrPrecOutOfRange) {
		t.Errorf("expected ErrPrecOutOfRange, got %v", err)
	}
}

func TestMustFromInt64(t *testing.T) {
	d := MustFromInt64(500, 2)
	if d.String() != "5" {
		t.Errorf("expected 5, got %s", d.String())
	}
}

func TestMustFromInt64PanicsOnInvalidPrecision(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for precision out of range")
		}
	}()
	MustFromInt64(1, 20)
}

func TestNewFromUint64(t *testing.T) {
	d, err := NewFromUint64(12345, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.String() != "123.45" {
		t.Errorf("expected 123.45, got %s", d.String())
	}

	if _, err := NewFromUint64(1, 20); !errors.Is(err, ErrPrecOutOfRange) {
		t.Errorf("expected ErrPrecOutOfRange, got %v", err)
	}
}

func TestMustFromUint64(t *testing.T) {
	d := MustFromUint64(999, 0)
	if d.String() != "999" {
		t.Errorf("expected 999, got %s", d.String())
	}
}

func TestMustFromUint64PanicsOnInvalidPrecision(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for precision out of range")
		}
	}()
	MustFromUint64(1, 20)
}

func TestNewFromHiLo(t *testing.T) {
	// hi=0, lo=123456 with scale 3 represents 123.456
	d, err := NewFromHiLo(false, 0, 123456, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.String() != "123.456" {
		t.Errorf("expected 123.456, got %s", d.String())
	}

	neg, err := NewFromHiLo(true, 0, 5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if neg.String() != "-5" {
		t.Errorf("expected -5, got %s", neg.String())
	}

	if _, err := NewFromHiLo(false, 0, 1, 20); !errors.Is(err, ErrPrecOutOfRange) {
		t.Errorf("expected ErrPrecOutOfRange, got %v", err)
	}
}

func TestMustFromHiLo(t *testing.T) {
	d := MustFromHiLo(false, 0, 42, 0)
	if d.String() != "42" {
		t.Errorf("expected 42, got %s", d.String())
	}
}

func TestMustFromHiLoPanicsOnInvalidPrecision(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for precision out of range")
		}
	}()
	MustFromHiLo(false, 0, 1, 20)
}

func TestMustFromStringPanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid decimal string")
		}
	}()
	MustFromString("not-a-number")
}

func TestDecimalArithmetic(t *testing.T) {
	a := MustFromString("10.5")
	b := MustFromString("3")

	if got := a.Add(b).String(); got != "13.5" {
		t.Errorf("Add: expected 13.5, got %s", got)
	}
	if got := a.Sub(b).String(); got != "7.5" {
		t.Errorf("Sub: expected 7.5, got %s", got)
	}
	if got := a.Mul(b).String(); got != "31.5" {
		t.Errorf("Mul: expected 31.5, got %s", got)
	}

	quotient, err := a.Div(b)
	if err != nil {
		t.Fatalf("Div: unexpected error: %v", err)
	}
	if got := quotient.String(); got != "3.5" {
		t.Errorf("Div: expected 3.5, got %s", got)
	}

	if _, err := a.Div(Zero); !errors.Is(err, ErrDivideByZero) {
		t.Errorf("Div by zero: expected ErrDivideByZero, got %v", err)
	}

	if got := a.MustDiv(b).String(); got != "3.5" {
		t.Errorf("MustDiv: expected 3.5, got %s", got)
	}
}

func TestDecimalMustDivPanicsOnDivideByZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for division by zero")
		}
	}()
	One.MustDiv(Zero)
}

func TestDecimalMod(t *testing.T) {
	a := MustFromString("10")
	b := MustFromString("3")

	m, err := a.Mod(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.String() != "1" {
		t.Errorf("expected 1, got %s", m.String())
	}

	if _, err := a.Mod(Zero); !errors.Is(err, ErrDivideByZero) {
		t.Errorf("expected ErrDivideByZero, got %v", err)
	}
}

func TestDecimalDiv64(t *testing.T) {
	a := MustFromString("10")

	q, err := a.Div64(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.String() != "2.5" {
		t.Errorf("expected 2.5, got %s", q.String())
	}

	if _, err := a.Div64(0); !errors.Is(err, ErrDivideByZero) {
		t.Errorf("expected ErrDivideByZero, got %v", err)
	}
}

func TestDecimalSign(t *testing.T) {
	if MustFromString("5").Sign() != 1 {
		t.Error("expected sign 1 for positive")
	}
	if MustFromString("-5").Sign() != -1 {
		t.Error("expected sign -1 for negative")
	}
	if Zero.Sign() != 0 {
		t.Error("expected sign 0 for zero")
	}
}

func TestDecimalCmp(t *testing.T) {
	a := MustFromString("5")
	b := MustFromString("10")

	if a.Cmp(b) >= 0 {
		t.Error("expected a < b")
	}
	if b.Cmp(a) <= 0 {
		t.Error("expected b > a")
	}
	if a.Cmp(a) != 0 {
		t.Error("expected a == a")
	}
}

func TestDecimalComparisonOperators(t *testing.T) {
	a := MustFromString("5")
	b := MustFromString("10")

	if !a.LessThan(b) || b.LessThan(a) {
		t.Error("LessThan failed")
	}
	if !a.LessThanOrEqual(a) || !a.LessThanOrEqual(b) || b.LessThanOrEqual(a) {
		t.Error("LessThanOrEqual failed")
	}
	if !b.GreaterThan(a) || a.GreaterThan(b) {
		t.Error("GreaterThan failed")
	}
	if !a.GreaterThanOrEqual(a) || !b.GreaterThanOrEqual(a) || a.GreaterThanOrEqual(b) {
		t.Error("GreaterThanOrEqual failed")
	}
	if !a.Equal(MustFromString("5.0")) {
		t.Error("Equal failed for equivalent values with different scale")
	}
}

func TestDecimalFloat64(t *testing.T) {
	d := MustFromString("123.456")
	f, err := d.Float64()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != 123.456 {
		t.Errorf("expected 123.456, got %v", f)
	}
}

func TestDecimalInt64(t *testing.T) {
	d := MustFromString("123.999")
	v, err := d.Int64()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 123 {
		t.Errorf("expected 123, got %d", v)
	}

	neg := MustFromString("-42.5")
	v2, err := neg.Int64()
	if err != nil || v2 != -42 {
		t.Errorf("expected -42, got %d, err=%v", v2, err)
	}
}

func TestDecimalIsZeroIsNegIsPos(t *testing.T) {
	if !Zero.IsZero() || Zero.IsNeg() || Zero.IsPos() {
		t.Error("Zero should be zero, not negative or positive")
	}

	pos := MustFromString("1")
	if pos.IsZero() || pos.IsNeg() || !pos.IsPos() {
		t.Error("expected positive decimal to report IsPos")
	}

	neg := MustFromString("-1")
	if neg.IsZero() || !neg.IsNeg() || neg.IsPos() {
		t.Error("expected negative decimal to report IsNeg")
	}
}

func TestDecimalRoundAwayHAZHTZ(t *testing.T) {
	d := MustFromString("1.5")

	if got := d.RoundAway(0).String(); got != "2" {
		t.Errorf("RoundAway: expected 2, got %s", got)
	}
	if got := d.RoundHAZ(0).String(); got != "2" {
		t.Errorf("RoundHAZ: expected 2, got %s", got)
	}
	if got := d.RoundHTZ(0).String(); got != "1" {
		t.Errorf("RoundHTZ: expected 1, got %s", got)
	}
}

func TestDecimalMustLog10PanicsOnInvalidResult(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for Log10(0)")
		}
	}()
	Zero.MustLog10()
}

func TestDecimalMustLog2PanicsOnInvalidResult(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for Log2(-1)")
		}
	}()
	MustFromString("-1").MustLog2()
}

func TestDecimalMustFromFloat64PanicsOnNaN(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for NaN")
		}
	}()
	MustFromFloat64(math.NaN())
}

func TestDecimalMustLnSuccess(t *testing.T) {
	if got := MustFromFloat64(math.E).MustLn().String(); got != "1" {
		t.Errorf("expected 1, got %s", got)
	}
}

func TestDecimalMustLog10Success(t *testing.T) {
	if got := MustFromFloat64(1000).MustLog10().String(); got != "3" {
		t.Errorf("expected 3, got %s", got)
	}
}

func TestDecimalMustLog2Success(t *testing.T) {
	if got := MustFromFloat64(8).MustLog2().String(); got != "3" {
		t.Errorf("expected 3, got %s", got)
	}
}

func TestDecimalMustLogSuccess(t *testing.T) {
	if got := MustFromFloat64(8).MustLog(MustFromFloat64(2)).String(); got != "3" {
		t.Errorf("expected 3, got %s", got)
	}
}

func TestDecimalAddOverflowPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Add overflow")
		}
	}()

	huge := Decimal{decimal128{coef: u128{hi: ^uint64(0), lo: ^uint64(0)}, scale: 0}}
	huge.Add(huge)
}

func TestDecimalSubOverflowPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Sub overflow")
		}
	}()

	huge := Decimal{decimal128{coef: u128{hi: ^uint64(0), lo: ^uint64(0)}, scale: 0}}
	huge.Sub(huge.Neg())
}

func TestDecimalMulOverflowPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Mul overflow")
		}
	}()

	huge := Decimal{decimal128{coef: u128{hi: ^uint64(0) >> 1, lo: ^uint64(0)}, scale: 0}}
	huge.Mul(huge)
}

func TestDecimalMustLogPanicsOnInvalidBase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for Log with base 1")
		}
	}()
	MustFromFloat64(10).MustLog(MustFromFloat64(1))
}
