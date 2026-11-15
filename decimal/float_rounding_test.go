package decimal

import (
	"math"
	"testing"
)

// A float64 whose shortest decimal form needs more than maxScale fractional
// digits is rounded to maxScale digits instead of being rejected. Rates such
// as 1/4100 (COP→USD) or 0.001/12 are ordinary inputs, and the binary value
// they hold is itself only an approximation, so refusing them protects no
// precision that was ever there.
func TestNewFromFloat64RoundsExcessFractionalDigits(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want string
	}{
		{"cop to usd rate", 1.0 / 4100, "0.0002439024390243902"},
		{"negative cop to usd rate", -1.0 / 4100, "-0.0002439024390243902"},
		{"small monthly rate", 0.001 / 12, "0.0000833333333333333"},
		{"rounds up to the smallest unit", 6e-20, "0.0000000000000000001"},
		{"below half the smallest unit", 1e-20, "0"},
		{"negative below half the smallest unit", -1e-20, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewFromFloat64(tt.in)
			if err != nil {
				t.Fatalf("NewFromFloat64(%v) returned error %v", tt.in, err)
			}

			if got := d.String(); got != tt.want {
				t.Fatalf("NewFromFloat64(%v) = %s, want %s", tt.in, got, tt.want)
			}
		})
	}
}

// Floats that already fit keep their shortest form, scale included, exactly
// as before: rounding only ever applies when the digits do not fit.
func TestNewFromFloat64KeepsShortestFormWhenItFits(t *testing.T) {
	tests := []struct {
		in        float64
		want      string
		wantScale uint8
	}{
		{0.1, "0.1", 1},
		{123.456, "123.456", 3},
		{0.05 / 365, "0.000136986301369863", 18},
		// Computed at run time: as a constant expression 0.1 + 0.2 is exactly 0.3.
		{addFloats(0.1, 0.2), "0.30000000000000004", 17},
	}

	for _, tt := range tests {
		d, err := NewFromFloat64(tt.in)
		if err != nil {
			t.Fatalf("NewFromFloat64(%v) returned error %v", tt.in, err)
		}

		if got := d.String(); got != tt.want {
			t.Errorf("NewFromFloat64(%v) = %s, want %s", tt.in, got, tt.want)
		}

		if d.value.scale != tt.wantScale {
			t.Errorf("NewFromFloat64(%v) scale = %d, want %d", tt.in, d.value.scale, tt.wantScale)
		}
	}
}

func addFloats(a, b float64) float64 { return a + b }

// Rounding must not mask the inputs that have no decimal value at all.
func TestNewFromFloat64StillRejectsNonFinite(t *testing.T) {
	for _, f := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := NewFromFloat64(f); err != ErrInvalidFormat {
			t.Errorf("NewFromFloat64(%v) error = %v, want ErrInvalidFormat", f, err)
		}
	}
}
