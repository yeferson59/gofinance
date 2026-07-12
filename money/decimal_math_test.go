package money

import (
	"math"
	"testing"
)

func TestDecimalConstructors(t *testing.T) {
	if d, err := NewFromUint64(100, 2); err != nil || d.String() != "1" {
		t.Errorf("NewFromUint64(100,2): expected 1, got %s (err=%v)", d.String(), err)
	}

	if d, err := NewFromHiLo(false, 0, 12345, 2); err != nil || d.String() != "123.45" {
		t.Errorf("NewFromHiLo: expected 123.45, got %s (err=%v)", d.String(), err)
	}

	if d := MustFromUint64(500, 2); d.String() != "5" {
		t.Errorf("MustFromUint64: expected 5, got %s", d.String())
	}

	if d := MustFromHiLo(false, 0, 12345, 2); d.String() != "123.45" {
		t.Errorf("MustFromHiLo: expected 123.45, got %s", d.String())
	}
}

func TestDecimalConstructorPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustFromUint64: expected panic on invalid precision")
		}
	}()
	MustFromUint64(1, 200)
}

func TestDecimalConstructorHiLoPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustFromHiLo: expected panic on invalid precision")
		}
	}()
	MustFromHiLo(false, 0, 1, 200)
}

func TestDecimalSqrt(t *testing.T) {
	four := MustFromInt64(4, 0)

	sqrt, err := four.Sqrt()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sqrt.String() != "2" {
		t.Errorf("expected 2, got %s", sqrt.String())
	}

	if MustFromInt64(9, 0).MustSqrt().String() != "3" {
		t.Errorf("MustSqrt: expected 3, got %s", MustFromInt64(9, 0).MustSqrt().String())
	}
}

func TestDecimalSqrtPanicsOnNegative(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustSqrt: expected panic for negative input")
		}
	}()
	MustFromInt64(-4, 0).MustSqrt()
}

func TestDecimalLn(t *testing.T) {
	one := MustFromInt64(1, 0)

	ln, err := one.Ln()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ln.IsZero() {
		t.Errorf("expected ln(1) == 0, got %s", ln.String())
	}

	if got := one.MustLn(); !got.IsZero() {
		t.Errorf("MustLn: expected 0, got %s", got.String())
	}
}

func TestDecimalLog10(t *testing.T) {
	hundred := MustFromInt64(100, 0)

	log10, err := hundred.Log10()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log10.String() != "2" {
		t.Errorf("expected 2, got %s", log10.String())
	}

	if got := hundred.MustLog10(); got.String() != "2" {
		t.Errorf("MustLog10: expected 2, got %s", got.String())
	}
}

func TestDecimalLog2(t *testing.T) {
	eight := MustFromInt64(8, 0)

	log2, err := eight.Log2()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log2.String() != "3" {
		t.Errorf("expected 3, got %s", log2.String())
	}

	if got := eight.MustLog2(); got.String() != "3" {
		t.Errorf("MustLog2: expected 3, got %s", got.String())
	}
}

func TestDecimalLog(t *testing.T) {
	eight := MustFromInt64(8, 0)
	base := MustFromInt64(2, 0)

	log, err := eight.Log(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log.String() != "3" {
		t.Errorf("expected 3, got %s", log.String())
	}

	if got := eight.MustLog(base); got.String() != "3" {
		t.Errorf("MustLog: expected 3, got %s", got.String())
	}
}

func TestDecimalLogPanicsOnInvalidBase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLog: expected panic for base 1")
		}
	}()
	MustFromInt64(8, 0).MustLog(MustFromInt64(1, 0))
}

func TestDecimalMod(t *testing.T) {
	ten := MustFromInt64(10, 0)
	three := MustFromInt64(3, 0)

	mod, err := ten.Mod(three)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mod.String() != "1" {
		t.Errorf("expected 1, got %s", mod.String())
	}
}

func TestDecimalModByZeroErrors(t *testing.T) {
	ten := MustFromInt64(10, 0)

	if _, err := ten.Mod(Zero); err == nil {
		t.Error("expected error for mod by zero")
	}
}

func TestDecimalDiv64(t *testing.T) {
	ten := MustFromInt64(10, 0)

	half, err := ten.Div64(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if half.String() != "5" {
		t.Errorf("expected 5, got %s", half.String())
	}
}

func TestDecimalDiv64ByZeroErrors(t *testing.T) {
	ten := MustFromInt64(10, 0)

	if _, err := ten.Div64(0); err == nil {
		t.Error("expected error for division by zero")
	}
}

func TestDecimalRounding(t *testing.T) {
	d := MustFromString("1.005")

	if got := d.RoundBank(2).String(); got != "1" && got != "1.00" {
		t.Errorf("RoundBank: unexpected result %s", got)
	}
	if got := d.RoundAway(2).StringFixed(2); got != "1.01" {
		t.Errorf("RoundAway: expected 1.01, got %s", got)
	}
	if got := d.RoundHAZ(2).StringFixed(2); got != "1.01" {
		t.Errorf("RoundHAZ: expected 1.01, got %s", got)
	}
	if got := d.RoundHTZ(2).StringFixed(2); got != "1.00" {
		t.Errorf("RoundHTZ: expected 1.00, got %s", got)
	}
	if got := d.Trunc(2).StringFixed(2); got != "1.00" {
		t.Errorf("Trunc: expected 1.00, got %s", got)
	}
}

func TestDecimalFloorCeil(t *testing.T) {
	d := MustFromString("1.5")

	if got := d.Floor().String(); got != "1" {
		t.Errorf("Floor: expected 1, got %s", got)
	}
	if got := d.Ceil().String(); got != "2" {
		t.Errorf("Ceil: expected 2, got %s", got)
	}
}

func TestDecimalSignAndCmp(t *testing.T) {
	pos := MustFromInt64(5, 0)
	neg := MustFromInt64(-5, 0)

	if pos.Sign() != 1 {
		t.Errorf("expected sign 1, got %d", pos.Sign())
	}
	if neg.Sign() != -1 {
		t.Errorf("expected sign -1, got %d", neg.Sign())
	}
	if Zero.Sign() != 0 {
		t.Errorf("expected sign 0, got %d", Zero.Sign())
	}

	if pos.Cmp(neg) <= 0 {
		t.Errorf("expected pos > neg")
	}
	if neg.Cmp(pos) >= 0 {
		t.Errorf("expected neg < pos")
	}
	if pos.Cmp(pos) != 0 {
		t.Errorf("expected equal comparison to be 0")
	}
}

func TestDecimalFloat64(t *testing.T) {
	d := MustFromString("3.14")

	f, err := d.Float64()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != 3.14 {
		t.Errorf("expected 3.14, got %v", f)
	}
}

func TestDecimalLessThanOrEqual(t *testing.T) {
	one := MustFromInt64(1, 0)
	two := MustFromInt64(2, 0)

	if !one.LessThanOrEqual(two) {
		t.Error("expected 1 <= 2")
	}
	if !one.LessThanOrEqual(one) {
		t.Error("expected 1 <= 1")
	}
	if two.LessThanOrEqual(one) {
		t.Error("expected 2 <= 1 to be false")
	}
}

func TestDecimalEqual(t *testing.T) {
	a := MustFromString("1.50")
	b := MustFromString("1.5")

	if !a.Equal(b) {
		t.Error("expected 1.50 to equal 1.5")
	}
	if a.Equal(MustFromInt64(2, 0)) {
		t.Error("expected 1.50 to not equal 2")
	}
}

func TestDecimalMarshalUnmarshalJSON(t *testing.T) {
	d := MustFromString("42.5")

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

func TestDecimalUnmarshalJSONInvalid(t *testing.T) {
	var d Decimal
	if err := d.UnmarshalJSON([]byte("not-json")); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDecimalMustDivPanicsOnDivideByZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for division by zero")
		}
	}()
	MustFromInt64(1, 0).MustDiv(Zero)
}

func TestDecimalMustPowPanicsOnInvalidResult(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for negative base with fractional exponent")
		}
	}()
	MustFromInt64(-2, 0).MustPow(MustFromString("0.5"))
}

func TestDecimalMustLnPanicsOnNonPositive(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for non-positive input")
		}
	}()
	MustFromInt64(-1, 0).MustLn()
}

func TestDecimalMustLog10PanicsOnNonPositive(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for non-positive input")
		}
	}()
	MustFromInt64(-1, 0).MustLog10()
}

func TestDecimalMustLog2PanicsOnNonPositive(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for non-positive input")
		}
	}()
	MustFromInt64(-1, 0).MustLog2()
}

func TestDecimalMustFromFloat64PanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for NaN")
		}
	}()
	MustFromFloat64(math.NaN())
}

func TestDecimalMustFromInt64PanicsOnInvalidPrecision(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for out-of-range precision")
		}
	}()
	MustFromInt64(1, 200)
}

func TestDecimalMustFromStringPanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid string")
		}
	}()
	MustFromString("not-a-number")
}
