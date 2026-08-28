package money

import (
	"errors"
	"math"
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
)

func TestNewMoney(t *testing.T) {
	m, err := New(12345, 2, USD)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.String() != "123.45" {
		t.Errorf("expected 123.45, got %s", m.String())
	}
	if m.GetCurrency() != USD {
		t.Errorf("expected USD, got %v", m.GetCurrency())
	}

	if _, err := New(1, 20, USD); !errors.Is(err, decimal.ErrPrecOutOfRange) {
		t.Errorf("expected ErrPrecOutOfRange, got %v", err)
	}
}

func TestNewMoneyFromStringInvalid(t *testing.T) {
	if _, err := NewMoneyFromString("not-a-number", USD); err == nil {
		t.Error("expected error for invalid decimal string")
	}
}

func TestMustMoneyFromStringPanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid decimal string")
		}
	}()
	MustMoneyFromString("not-a-number", USD)
}

func TestMoneyToDecimal(t *testing.T) {
	d := MustMoneyFromString("100.50", USD).GetDecimal()
	if d.String() != "100.5" {
		t.Errorf("expected 100.5, got %s", d.String())
	}
}

func TestMoneySub(t *testing.T) {
	a := MustMoneyFromString("10", USD)
	b := MustMoneyFromString("3.5", USD)

	diff := a.Sub(b)
	if diff.String() != "6.5" {
		t.Errorf("expected 6.5, got %s", diff.String())
	}
}

func TestMoneyRoundBankString(t *testing.T) {
	m := MustMoneyFromString("1.235", USD)
	if got := m.RoundBankString(2); got != "1.24" {
		t.Errorf("expected 1.24, got %s", got)
	}
}

func TestMoneyRoundAway(t *testing.T) {
	m := MustMoneyFromString("1.21", USD)
	rounded := m.RoundAway(1)
	if rounded.String() != "1.3" {
		t.Errorf("expected 1.3, got %s", rounded.String())
	}
	if rounded.GetCurrency() != USD {
		t.Errorf("currency should be preserved")
	}
}

func TestMoneyTrunc(t *testing.T) {
	m := MustMoneyFromString("1.999", USD)
	trunc := m.Trunc(1)
	if trunc.String() != "1.9" {
		t.Errorf("expected 1.9, got %s", trunc.String())
	}
}

func TestMoneyAbs(t *testing.T) {
	m := MustMoneyFromString("-42.50", USD)
	abs := m.Abs()
	if abs.String() != "42.5" {
		t.Errorf("expected 42.5, got %s", abs.String())
	}
	if abs.GetCurrency() != USD {
		t.Errorf("currency should be preserved")
	}
}

func TestMoneyIsZero(t *testing.T) {
	if !MustMoneyFromString("0", USD).IsZero() {
		t.Error("expected 0 to be zero")
	}
	if MustMoneyFromString("0.01", USD).IsZero() {
		t.Error("expected 0.01 to not be zero")
	}
}

func TestMoneyInexactFloat64(t *testing.T) {
	m := MustMoneyFromString("19.99", USD)
	if got := m.InexactFloat64(); got != 19.99 {
		t.Errorf("expected 19.99, got %v", got)
	}
}

func TestMoneyCmp(t *testing.T) {
	a := MustMoneyFromString("5", USD)
	b := MustMoneyFromString("10", USD)

	if a.Cmp(b) >= 0 {
		t.Error("expected a < b")
	}
	if a.Cmp(a) != 0 {
		t.Error("expected a == a")
	}
}

func TestMoneyFloorCeil(t *testing.T) {
	m := MustMoneyFromString("1.5", USD)
	if got := m.Floor().String(); got != "1" {
		t.Errorf("Floor: expected 1, got %s", got)
	}
	if got := m.Ceil().String(); got != "2" {
		t.Errorf("Ceil: expected 2, got %s", got)
	}
	if m.Floor().GetCurrency() != USD || m.Ceil().GetCurrency() != USD {
		t.Error("currency should be preserved by Floor/Ceil")
	}
}

func TestMoneyComparisonOperators(t *testing.T) {
	a := MustMoneyFromString("5", USD)
	b := MustMoneyFromString("10", USD)

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
}

func TestMoneyEqualRequiresSameCurrency(t *testing.T) {
	usd := MustMoneyFromString("10", USD)
	eur := MustMoneyFromString("10", EUR)
	usd2 := MustMoneyFromString("10", USD)

	if usd.Equal(eur) {
		t.Error("expected different currencies to not be equal")
	}
	if !usd.Equal(usd2) {
		t.Error("expected same value+currency to be equal")
	}
}

func TestMoneyMarshalJSONInvalidCurrency(t *testing.T) {
	m := Money{value: decimal.One, currency: Currency(255)}
	if _, err := m.MarshalJSON(); err == nil {
		t.Error("expected error marshalling money with an unknown currency")
	}
}

func TestMoneyUnmarshalJSONPlainNumber(t *testing.T) {
	var m Money
	if err := m.UnmarshalJSON([]byte("123.45")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.String() != "123.45" {
		t.Errorf("expected 123.45, got %s", m.String())
	}
	if m.GetCurrency() != USD {
		t.Errorf("expected default currency USD, got %v", m.GetCurrency())
	}
}

func TestMoneyUnmarshalJSONWithCurrency(t *testing.T) {
	m, err := NewMoneyFromString("99.99", EUR)
	if err != nil {
		t.Fatal(err)
	}

	d, err := m.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var m2 Money
	if err := m2.UnmarshalJSON(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m2.String() != "99.99" {
		t.Errorf("expected 99.99, got %s", m.String())
	}
	if m2.GetCurrency() != EUR {
		t.Errorf("expected EUR, got %v", m.GetCurrency())
	}
}

func TestMoneyUnmarshalJSONEmptyCurrencyDefaultsToUSD(t *testing.T) {
	var m Money
	if err := m.UnmarshalJSON([]byte(`{"value":5.00,"currency":""}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.GetCurrency() != USD {
		t.Errorf("expected default currency USD, got %v", m.GetCurrency())
	}
}

func TestMoneyUnmarshalJSONInvalidCurrency(t *testing.T) {
	var m Money
	if err := m.UnmarshalJSON([]byte(`{"value":"5.00","currency":"ZZZ"}`)); err == nil {
		t.Error("expected error for unknown currency code")
	}
}

func TestMoneyUnmarshalJSONInvalidValue(t *testing.T) {
	var m Money
	if err := m.UnmarshalJSON([]byte(`{"value":"abc","currency":"USD"}`)); err == nil {
		t.Error("expected error for invalid decimal value")
	}
}

func TestMoneyUnmarshalJSONMalformed(t *testing.T) {
	var m Money
	if err := m.UnmarshalJSON([]byte(`not-json`)); err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestMoneyValue(t *testing.T) {
	m := MustMoneyFromString("42.50", USD)
	v, err := m.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "42.5" {
		t.Errorf("expected \"42.5\", got %v", v)
	}
}

func TestMoneyScan(t *testing.T) {
	tests := []struct {
		name  string
		src   any
		want  string
		error bool
	}{
		{"bytes", []byte("12.34"), "12.34", false},
		{"string", "56.78", "56.78", false},
		{"uint64", uint64(100), "100", false},
		{"int64", int64(-50), "-50", false},
		{"int", int(7), "7", false},
		{"int32", int32(-3), "-3", false},
		{"float64", 3.14, "3.14", false},
		{"nil", nil, "", true},
		{"unsupported", true, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m Money
			err := m.Scan(tt.src)
			if tt.error {
				if err == nil {
					t.Errorf("expected error scanning %#v", tt.src)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := m.String(); got != tt.want {
				t.Errorf("Scan(%#v): expected %s, got %s", tt.src, tt.want, got)
			}
		})
	}
}

func TestNewMoneyFromFloat64InvalidValue(t *testing.T) {
	if _, err := NewMoneyFromFloat64(math.NaN(), USD); !errors.Is(err, decimal.ErrInvalidFormat) {
		t.Errorf("expected ErrInvalidFormat for NaN, got %v", err)
	}
	if _, err := NewMoneyFromFloat64(math.Inf(-1), USD); !errors.Is(err, decimal.ErrInvalidFormat) {
		t.Errorf("expected ErrInvalidFormat for -Inf, got %v", err)
	}
}

func TestMustMoneyFromFloat64PanicsOnInvalidValue(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for NaN")
		}
	}()
	MustMoneyFromFloat64(math.NaN(), USD)
}

func hugeMoney(currency Currency) Money {
	return Money{
		value:    decimal.MustFromHiLo(false, ^uint64(0), ^uint64(0), 0),
		currency: currency,
	}
}

func TestMoneyAddOverflowPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Add overflow")
		}
	}()
	m := hugeMoney(USD)
	m.Add(m)
}

func TestMoneySubOverflowPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Sub overflow")
		}
	}()
	m := hugeMoney(USD)
	m.Sub(m.Neg())
}

func TestMoneyMulInt64OverflowPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on MulInt64 overflow")
		}
	}()
	m := Money{
		value:    decimal.MustFromHiLo(false, ^uint64(0)>>1, ^uint64(0), 0),
		currency: USD,
	}
	m.MulInt64(math.MaxInt64)
}

func TestMoneyDivInt64Overflow(t *testing.T) {
	// Dividing by a tiny divisor forces the result's coefficient past 128
	// bits once rescaled to maxScale digits of precision.
	m := hugeMoney(USD)
	if _, err := m.DivInt64(1); !errors.Is(err, decimal.ErrOverflow) {
		t.Errorf("expected ErrOverflow, got %v", err)
	}
}

func TestMoneyMustDivInt64(t *testing.T) {
	m := MustMoneyFromString("100", USD)
	if got := m.MustDivInt64(4).String(); got != "25" {
		t.Errorf("expected 25, got %s", got)
	}
}

func TestMoneyMustDivInt64PanicsOnDivideByZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for division by zero")
		}
	}()
	MustMoneyFromString("100", USD).MustDivInt64(0)
}

func TestMoneyMinReturnsOtherWhenSmaller(t *testing.T) {
	a := MustMoneyFromString("20", USD)
	b := MustMoneyFromString("5", USD)

	got, err := a.Min(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(b) {
		t.Errorf("expected min to be %s, got %s", b.String(), got.String())
	}
}

func TestMoneyMaxReturnsSelfWhenLarger(t *testing.T) {
	a := MustMoneyFromString("20", USD)
	b := MustMoneyFromString("5", USD)

	got, err := a.Max(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(a) {
		t.Errorf("expected max to be %s, got %s", a.String(), got.String())
	}
}

func TestMoneyUnmarshalJSONScientificNotationRejected(t *testing.T) {
	var m Money
	if err := m.UnmarshalJSON([]byte("1e2")); err == nil {
		t.Error("expected error for scientific notation, which parseDecimal doesn't support")
	}
}

func TestMoneyScanInvalidStringPropagatesParseError(t *testing.T) {
	var m Money
	if err := m.Scan("not-a-number"); err == nil {
		t.Error("expected parse error to propagate from Scan")
	}
}
