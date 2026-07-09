package money

import (
	"errors"
	"testing"
)

func TestConvertAppliesRateAndCurrency(t *testing.T) {
	usd := MustMoneyFromFloat64(100, USD)

	eur, err := usd.Convert(EUR, MustFromFloat64(0.92))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eur.Currency() != EUR {
		t.Errorf("expected EUR currency, got %v", eur.Currency())
	}
	if eur.StringFixed(2) != "92.00" {
		t.Errorf("expected 92.00, got %s", eur.StringFixed(2))
	}
}

func TestConvertRoundsToTargetPrecision(t *testing.T) {
	usd := MustMoneyFromFloat64(10, USD)

	// JPY has zero decimal places.
	jpy, err := usd.Convert(JPY, MustFromFloat64(150.4))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if jpy.StringFixed(0) != "1504" {
		t.Errorf("expected 1504, got %s", jpy.StringFixed(0))
	}
}

func TestConvertInvalidRate(t *testing.T) {
	usd := MustMoneyFromFloat64(10, USD)

	if _, err := usd.Convert(EUR, Zero); !errors.Is(err, ErrInvalidExchangeRate) {
		t.Errorf("expected ErrInvalidExchangeRate for zero rate, got %v", err)
	}

	negRate := MustFromFloat64(-1)
	if _, err := usd.Convert(EUR, negRate); !errors.Is(err, ErrInvalidExchangeRate) {
		t.Errorf("expected ErrInvalidExchangeRate for negative rate, got %v", err)
	}
}

func TestConvertUnknownTargetCurrency(t *testing.T) {
	usd := MustMoneyFromFloat64(10, USD)

	if _, err := usd.Convert(Currency(9999), One); err == nil {
		t.Error("expected error for unknown target currency")
	}
}

func TestConvertFloat64(t *testing.T) {
	usd := MustMoneyFromFloat64(50, USD)

	gbp, err := usd.ConvertFloat64(GBP, 0.79)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gbp.Currency() != GBP {
		t.Errorf("expected GBP currency, got %v", gbp.Currency())
	}
	if gbp.StringFixed(2) != "39.50" {
		t.Errorf("expected 39.50, got %s", gbp.StringFixed(2))
	}
}

func TestMustConvertPanicsOnInvalidRate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid rate")
		}
	}()

	usd := MustMoneyFromFloat64(10, USD)
	usd.MustConvert(EUR, Zero)
}
