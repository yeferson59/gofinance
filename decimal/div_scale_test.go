package decimal

import "testing"

var fourE19 = MustFromString("40000000000000000000") // 4e19

// Division works at 19 fractional digits, so even an exact quotient such as
// 1/2 carries them: 0.5000000000000000000. Aligning a large integer with that
// needed 19 more digits than the integer had, so 0.5 (from Div) + 4e19
// overflowed although the sum is tiny next to the engine's ~3.4e38 range. Add
// now retries with the trailing zeros dropped.
func TestAddAfterDivKeepsRange(t *testing.T) {
	half := One.MustDiv(MustFromInt64(2, 0))

	sum, err := half.TryAdd(fourE19)
	if err != nil {
		t.Fatalf("0.5 (from Div) + 4e19: %v", err)
	}

	if got := sum.String(); got != "40000000000000000000.5" {
		t.Errorf("0.5 + 4e19 = %s", got)
	}

	diff, err := fourE19.TrySub(half)
	if err != nil {
		t.Fatalf("4e19 - 0.5 (from Div): %v", err)
	}

	if got := diff.String(); got != "39999999999999999999.5" {
		t.Errorf("4e19 - 0.5 = %s", got)
	}
}

// A quotient above ~3.4e19 can't carry 19 fractional digits; Div gives some
// up instead of overflowing.
func TestDivOfLargeValue(t *testing.T) {
	q, err := MustFromString("100000000000000000000").Div(One)
	if err != nil {
		t.Fatalf("1e20 / 1: %v", err)
	}

	if got := q.String(); got != "100000000000000000000" {
		t.Errorf("1e20 / 1 = %s", got)
	}

	q, err = MustFromString("100000000000000000000").Div64(3)
	if err != nil {
		t.Fatalf("1e20 / 3: %v", err)
	}

	if got := q.String(); got != "33333333333333333333.3333333333333333333" {
		t.Errorf("1e20 / 3 = %s", got)
	}
}

// Mul retries the same way: 1.0000000000000000000 (a quotient's shape) times
// 1e20 overflows at 19 fractional digits but not once the zeros are gone.
func TestMulAfterDivKeepsRange(t *testing.T) {
	one := MustFromInt64(3, 0).MustDiv(MustFromInt64(3, 0)) // 1 at scale 19

	p, err := one.TryMul(MustFromString("100000000000000000000"))
	if err != nil {
		t.Fatalf("1 (from Div) × 1e20: %v", err)
	}

	if got := p.String(); got != "100000000000000000000" {
		t.Errorf("1 × 1e20 = %s", got)
	}
}

// Inexact quotients are unchanged: they still carry all 19 digits.
func TestDivInexactKeepsAllDigits(t *testing.T) {
	third := One.MustDiv(MustFromInt64(3, 0))

	if got := third.String(); got != "0.3333333333333333333" {
		t.Errorf("1/3 = %s", got)
	}
}

// The trailing zeros are only dropped on the overflow path, so an ordinary
// quotient keeps the representation it always had — which is what
// MarshalBinary writes, so stored blobs don't change.
func TestDivKeepsRepresentation(t *testing.T) {
	half := One.MustDiv(MustFromInt64(2, 0))

	if half.value.scale != maxScale {
		t.Errorf("1/2 has scale %d, want %d", half.value.scale, maxScale)
	}
}
