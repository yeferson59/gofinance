package money

import (
	"errors"
	"math"
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

func TestConvertOverflow(t *testing.T) {
	huge := Money{
		value:    decimal128{coef: u128{hi: ^uint64(0) >> 1, lo: ^uint64(0)}, scale: 0},
		currency: USD,
	}

	if _, err := huge.Convert(EUR, MustFromFloat64(1e30)); !errors.Is(err, ErrOverflow) {
		t.Errorf("expected ErrOverflow, got %v", err)
	}
}

func TestMustConvertSuccess(t *testing.T) {
	usd := MustMoneyFromFloat64(100, USD)

	eur := usd.MustConvert(EUR, MustFromFloat64(0.92))
	if eur.Currency() != EUR {
		t.Errorf("expected EUR, got %v", eur.Currency())
	}
	if eur.StringFixed(2) != "92.00" {
		t.Errorf("expected 92.00, got %s", eur.StringFixed(2))
	}
}

func TestConvertFloat64InvalidRate(t *testing.T) {
	usd := MustMoneyFromFloat64(10, USD)

	if _, err := usd.ConvertFloat64(EUR, math.NaN()); !errors.Is(err, ErrInvalidFormat) {
		t.Errorf("expected ErrInvalidFormat for NaN rate, got %v", err)
	}
	if _, err := usd.ConvertFloat64(EUR, math.Inf(1)); !errors.Is(err, ErrInvalidFormat) {
		t.Errorf("expected ErrInvalidFormat for +Inf rate, got %v", err)
	}
}

func TestMustConvertFloat64Success(t *testing.T) {
	usd := MustMoneyFromFloat64(50, USD)

	gbp := usd.MustConvertFloat64(GBP, 0.79)
	if gbp.Currency() != GBP {
		t.Errorf("expected GBP, got %v", gbp.Currency())
	}
	if gbp.StringFixed(2) != "39.50" {
		t.Errorf("expected 39.50, got %s", gbp.StringFixed(2))
	}
}

func TestMustConvertFloat64PanicsOnInvalidRate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid rate")
		}
	}()

	usd := MustMoneyFromFloat64(10, USD)
	usd.MustConvertFloat64(EUR, math.NaN())
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
