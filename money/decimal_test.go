package money

import (
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

func TestNewMoneyFromFloat64(t *testing.T) {
	m, err := NewMoneyFromFloat64(100.50, USD)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m.Currency() != USD {
		t.Errorf("expected USD, got %v", m.Currency())
	}
}

func TestNewMoneyFromString(t *testing.T) {
	m, err := NewMoneyFromString("999.99", EUR)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m.Currency() != EUR {
		t.Errorf("expected EUR, got %v", m.Currency())
	}
}

func TestMoneyRoundBank(t *testing.T) {
	m, _ := NewMoneyFromString("1.235", USD)
	rounded := m.RoundBank(2)
	if rounded.Currency() != USD {
		t.Errorf("currency should be preserved, got %v", rounded.Currency())
	}
	if rounded.StringFixed(2) != "1.24" {
		t.Errorf("expected 1.24, got %s", rounded.StringFixed(2))
	}
}

func TestMoneyNeg(t *testing.T) {
	m, _ := NewMoneyFromString("100.00", USD)
	neg := m.Neg()
	if neg.Currency() != USD {
		t.Errorf("currency should be preserved")
	}
	if neg.String() != "-100" {
		t.Errorf("expected -100, got %s", neg.String())
	}
}

func TestMoneyJSON(t *testing.T) {
	m, _ := NewMoneyFromString("123.45", EUR)

	data, err := m.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var m2 Money
	if err := m2.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !m.Equal(m2) {
		t.Errorf("expected %v, got %v", m, m2)
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
