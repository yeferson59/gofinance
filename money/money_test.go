package money

import (
	"errors"
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
)

func TestSafeAddSameCurrency(t *testing.T) {
	a := MustMoneyFromFloat64(10, USD)
	b := MustMoneyFromFloat64(2.5, USD)

	sum, err := a.SafeAdd(b)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sum.String() != "12.5" {
		t.Errorf("expected 12.5, got %s", sum.String())
	}
	if sum.Currency() != USD {
		t.Errorf("expected USD currency, got %v", sum.Currency())
	}
}

func TestSafeAddCurrencyMismatch(t *testing.T) {
	a := MustMoneyFromFloat64(10, USD)
	b := MustMoneyFromFloat64(5, EUR)

	if _, err := a.SafeAdd(b); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestSafeSubSameCurrency(t *testing.T) {
	a := MustMoneyFromFloat64(10, USD)
	b := MustMoneyFromFloat64(2.5, USD)

	diff, err := a.SafeSub(b)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if diff.String() != "7.5" {
		t.Errorf("expected 7.5, got %s", diff.String())
	}
}

func TestSafeSubCurrencyMismatch(t *testing.T) {
	a := MustMoneyFromFloat64(10, USD)
	b := MustMoneyFromFloat64(5, EUR)

	if _, err := a.SafeSub(b); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestMoneyConstantsCurrency(t *testing.T) {
	if MoneyZero.Currency() != USD {
		t.Errorf("expected MoneyZero currency to be USD, got %v", MoneyZero.Currency())
	}
	if MoneyOne.Currency() != USD {
		t.Errorf("expected MoneyOne currency to be USD, got %v", MoneyOne.Currency())
	}
}

func TestGetCurrencyISOCodeSpelling(t *testing.T) {
	tests := []struct {
		currency Currency
		expected string
	}{
		{MMK, "MMK"},
		{MNT, "MNT"},
		{USD, "USD"},
	}

	for _, tt := range tests {
		code, err := tt.currency.GetCurrencyISOCode()
		if err != nil {
			t.Errorf("unexpected error for %v: %v", tt.currency, err)
			continue
		}
		if code != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, code)
		}
	}
}

func TestGetCurrencyPrecisionCode(t *testing.T) {
	tests := []struct {
		currency Currency
		expected uint8
	}{
		{USD, 2},
		{EUR, 2},
		{COP, 2},
		{JPY, 0},
		{CLP, 0},
		{KWD, 3},
		{BHD, 3},
	}

	for _, tt := range tests {
		prec, err := tt.currency.GetCurrencyPrecisionCode()
		if err != nil {
			t.Errorf("unexpected error for %v: %v", tt.currency, err)
			continue
		}
		if prec != tt.expected {
			t.Errorf("expected precision %d for %v, got %d", tt.expected, tt.currency, prec)
		}
	}

	if _, err := Currency(9999).GetCurrencyPrecisionCode(); err == nil {
		t.Error("expected error for unknown currency")
	}
}

func TestMulInt64(t *testing.T) {
	price := MustMoneyFromFloat64(19.99, USD)

	total := price.MulInt64(3)
	if total.StringFixed(2) != "59.97" {
		t.Errorf("expected 59.97, got %s", total.StringFixed(2))
	}
}

func TestDivInt64(t *testing.T) {
	total := MustMoneyFromFloat64(30, USD)

	share, err := total.DivInt64(3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if share.StringFixed(2) != "10.00" {
		t.Errorf("expected 10.00, got %s", share.StringFixed(2))
	}

	if _, err := total.DivInt64(0); !errors.Is(err, decimal.ErrDivideByZero) {
		t.Errorf("expected ErrDivideByZero, got %v", err)
	}
}

func TestMinMax(t *testing.T) {
	a := MustMoneyFromFloat64(10, USD)
	b := MustMoneyFromFloat64(20, USD)

	gotMin, err := a.Min(b)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !gotMin.Equal(a) {
		t.Errorf("expected min to be %s, got %s", a.String(), gotMin.String())
	}

	gotMax, err := a.Max(b)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !gotMax.Equal(b) {
		t.Errorf("expected max to be %s, got %s", b.String(), gotMax.String())
	}

	c := MustMoneyFromFloat64(10, EUR)
	if _, err := a.Min(c); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}
	if _, err := a.Max(c); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestIsPositiveIsNegative(t *testing.T) {
	pos := MustMoneyFromFloat64(1, USD)
	neg := MustMoneyFromFloat64(-1, USD)
	zero := MustMoneyFromFloat64(0, USD)

	if !pos.IsPositive() || pos.IsNegative() {
		t.Errorf("expected %s to be positive", pos.String())
	}
	if !neg.IsNegative() || neg.IsPositive() {
		t.Errorf("expected %s to be negative", neg.String())
	}
	if zero.IsPositive() || zero.IsNegative() {
		t.Errorf("expected %s to be neither positive nor negative", zero.String())
	}
}

func TestCurrencySymbol(t *testing.T) {
	tests := []struct {
		currency Currency
		expected string
	}{
		{USD, "$"},
		{EUR, "€"},
		{JPY, "¥"},
		{XTS, "XTS"}, // no distinct symbol, falls back to ISO code
	}

	for _, tt := range tests {
		symbol, err := tt.currency.Symbol()
		if err != nil {
			t.Errorf("unexpected error for %v: %v", tt.currency, err)
			continue
		}
		if symbol != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, symbol)
		}
	}

	if _, err := Currency(9999).Symbol(); err == nil {
		t.Error("expected error for unknown currency")
	}
}

func TestMoneyFormat(t *testing.T) {
	m := MustMoneyFromFloat64(1234.5, USD)

	got, err := m.Format()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "$1234.50" {
		t.Errorf("expected $1234.50, got %s", got)
	}
}

func TestMoneyStringMoney(t *testing.T) {
	m := MustMoneyFromFloat64(1234.5, USD)

	got, err := m.StringMoney()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "USD 1234.50" {
		t.Errorf("expected \"USD 1234.50\", got %s", got)
	}
}

func TestMoneyStringMoneyInvalidCurrency(t *testing.T) {
	m := Money{value: decimal.One, currency: Currency(9999)}

	if _, err := m.StringMoney(); err == nil {
		t.Error("expected error for unknown currency")
	}
}

func TestMoneyFormatInvalidCurrency(t *testing.T) {
	m := Money{value: decimal.One, currency: Currency(9999)}

	if _, err := m.Format(); err == nil {
		t.Error("expected error for unknown currency")
	}
}

func TestCurrencyFromISOCodeNormalization(t *testing.T) {
	tests := []struct {
		input string
		want  Currency
	}{
		{"usd", USD},
		{"  USD  ", USD},
		{"Eur", EUR},
	}

	for _, tt := range tests {
		got, err := CurrencyFromISOCode(tt.input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("expected %v for %q, got %v", tt.want, tt.input, got)
		}
	}
}

func TestCurrencyFromISOCodeEmpty(t *testing.T) {
	if _, err := CurrencyFromISOCode(""); err == nil {
		t.Error("expected error for empty ISO code")
	}
	if _, err := CurrencyFromISOCode("   "); err == nil {
		t.Error("expected error for whitespace-only ISO code")
	}
}

func TestCurrencyFromISOCodeUnknown(t *testing.T) {
	if _, err := CurrencyFromISOCode("ZZZ"); err == nil {
		t.Error("expected error for unknown ISO code")
	}
}

func TestCurrencyISOCodeRoundTrip(t *testing.T) {
	for currency, code := range currencyCode {
		parsed, err := CurrencyFromISOCode(code)
		if err != nil {
			t.Errorf("unexpected error for %s: %v", code, err)
			continue
		}
		if parsed != currency {
			t.Errorf("round trip mismatch for %s: expected %v, got %v", code, currency, parsed)
		}
	}
}

func TestTryAddCurrencyMismatch(t *testing.T) {
	a := MustMoneyFromString("10", USD)
	b := MustMoneyFromString("5", EUR)

	if _, err := a.TryAdd(b); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestTrySubCurrencyMismatch(t *testing.T) {
	a := MustMoneyFromString("10", USD)
	b := MustMoneyFromString("5", EUR)

	if _, err := a.TrySub(b); !errors.Is(err, ErrCurrencyMismatch) {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestTryAddSameCurrency(t *testing.T) {
	a := MustMoneyFromString("10.25", USD)
	b := MustMoneyFromString("5.75", USD)

	sum, err := a.TryAdd(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.String() != "16" || sum.Currency() != USD {
		t.Errorf("expected 16 USD, got %s %v", sum.String(), sum.Currency())
	}
}

func TestTrySubSameCurrency(t *testing.T) {
	a := MustMoneyFromString("10.25", USD)
	b := MustMoneyFromString("5.75", USD)

	diff, err := a.TrySub(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff.String() != "4.5" || diff.Currency() != USD {
		t.Errorf("expected 4.5 USD, got %s %v", diff.String(), diff.Currency())
	}
}

func TestAddCurrencyMismatchPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Add currency mismatch")
		}
	}()
	MustMoneyFromString("10", USD).Add(MustMoneyFromString("5", EUR))
}

func TestSubCurrencyMismatchPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on Sub currency mismatch")
		}
	}()
	MustMoneyFromString("10", USD).Sub(MustMoneyFromString("5", EUR))
}
