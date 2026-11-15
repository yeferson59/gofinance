package decimal

import (
	"errors"
	"math"
	"testing"
)

// Int64 must accept every value NewFromInt64 does, including the one int64
// whose magnitude has no positive counterpart.
func TestInt64RoundTripsMinInt64(t *testing.T) {
	d := MustFromInt64(math.MinInt64, 0)

	got, err := d.Int64()
	if err != nil {
		t.Fatalf("Int64() error = %v", err)
	}

	if got != math.MinInt64 {
		t.Fatalf("Int64() = %d, want %d", got, int64(math.MinInt64))
	}

	// One past it still overflows.
	below, err := d.TrySub(One)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := below.Int64(); !errors.Is(err, ErrIntPartOverflow) {
		t.Errorf("Int64() of MinInt64-1 error = %v, want ErrIntPartOverflow", err)
	}
}

// StringFixed rounds half away from zero when the value carries more digits
// than requested, and pads with zeros when it carries fewer.
func TestStringFixedRounds(t *testing.T) {
	tests := []struct {
		in   string
		prec uint8
		want string
	}{
		{"1.239", 2, "1.24"},
		{"5.123", 2, "5.12"},
		{"5.125", 2, "5.13"},
		{"-5.125", 2, "-5.13"},
		{"10.005", 2, "10.01"},
		{"0.004", 2, "0.00"},
		{"-0.004", 2, "0.00"},
		{"9.999", 2, "10.00"},
		{"2.5", 0, "3"},
		{"1.2", 4, "1.2000"},
	}

	for _, tt := range tests {
		if got := MustFromString(tt.in).StringFixed(tt.prec); got != tt.want {
			t.Errorf("StringFixed(%s, %d) = %s, want %s", tt.in, tt.prec, got, tt.want)
		}
	}
}

// StringFixedBank is the banker's-rounding counterpart: ties go to the even
// digit.
func TestStringFixedBankRoundsHalfToEven(t *testing.T) {
	tests := []struct {
		in   string
		prec uint8
		want string
	}{
		{"5.125", 2, "5.12"},
		{"5.135", 2, "5.14"},
		{"-5.125", 2, "-5.12"},
		{"5.1251", 2, "5.13"},
		{"2.5", 0, "2"},
		{"1.2", 4, "1.2000"},
	}

	for _, tt := range tests {
		if got := MustFromString(tt.in).StringFixedBank(tt.prec); got != tt.want {
			t.Errorf("StringFixedBank(%s, %d) = %s, want %s", tt.in, tt.prec, got, tt.want)
		}
	}
}

// Scan accepts every Go numeric kind a caller might hand it directly, and a
// float32 converts from its own shortest form rather than its float64
// widening (0.1f is not 0.10000000149011612).
func TestScanAcceptsAllNumericKinds(t *testing.T) {
	tests := []struct {
		src  any
		want string
	}{
		{float32(0.1), "0.1"},
		{float32(-2.5), "-2.5"},
		{int8(-8), "-8"},
		{int16(-16), "-16"},
		{uint(7), "7"},
		{uint8(8), "8"},
		{uint16(16), "16"},
		{uint32(32), "32"},
	}

	for _, tt := range tests {
		var d Decimal
		if err := d.Scan(tt.src); err != nil {
			t.Errorf("Scan(%T %v) error = %v", tt.src, tt.src, err)
			continue
		}

		if got := d.String(); got != tt.want {
			t.Errorf("Scan(%T %v) = %s, want %s", tt.src, tt.src, got, tt.want)
		}
	}
}

// NULL has no Decimal value. Scan reports it with a sentinel callers can test
// for, and leaves the receiver untouched.
func TestScanNullReturnsSentinel(t *testing.T) {
	d := MustFromString("1.5")

	if err := d.Scan(nil); !errors.Is(err, ErrScanNull) {
		t.Fatalf("Scan(nil) error = %v, want ErrScanNull", err)
	}

	if got := d.String(); got != "1.5" {
		t.Errorf("Scan(nil) changed the receiver to %s", got)
	}
}
