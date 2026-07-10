package money

import "testing"

func TestNewMoneyFromFloat64Success(t *testing.T) {
	m, err := NewMoneyFromFloat64(100.50, USD)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m.Currency() != USD {
		t.Errorf("expected USD, got %v", m.Currency())
	}
}

func TestNewMoneyFromStringSuccess(t *testing.T) {
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

func TestDecimalToMoney(t *testing.T) {
	d := MustFromString("100.50")

	withoutCurrency := d.ToMoney()
	if withoutCurrency.Currency() != USD {
		t.Errorf("expected USD (default currency) when omitted, got %v", withoutCurrency.Currency())
	}
	if withoutCurrency.String() != "100.5" {
		t.Errorf("expected 100.5, got %s", withoutCurrency.String())
	}

	withCurrency := d.ToMoney(EUR)
	if withCurrency.Currency() != EUR {
		t.Errorf("expected EUR, got %v", withCurrency.Currency())
	}
}
