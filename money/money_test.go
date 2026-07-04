package money

import (
	"errors"
	"testing"
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
