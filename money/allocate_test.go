package money

import (
	"errors"
	"testing"
)

func TestAllocateSumsExactly(t *testing.T) {
	m := MustMoneyFromFloat64(100, USD)

	parts, err := m.Allocate(1, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}

	sum := parts[0].Add(parts[1]).Add(parts[2])
	if !sum.Equal(m) {
		t.Errorf("expected parts to sum to %s, got %s", m.String(), sum.String())
	}

	// Fowler's algorithm hands the leftover cent(s) to the earliest ratios.
	want := []string{"33.34", "33.33", "33.33"}
	for i, p := range parts {
		if got := p.StringFixed(2); got != want[i] {
			t.Errorf("part %d: expected %s, got %s", i, want[i], got)
		}
	}
}

func TestAllocateByWeightedRatios(t *testing.T) {
	m := MustMoneyFromFloat64(50, USD)

	parts, err := m.Allocate(1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum := parts[0].Add(parts[1])
	if !sum.Equal(m) {
		t.Errorf("expected parts to sum to %s, got %s", m.String(), sum.String())
	}
	if parts[0].StringFixed(2) != "12.50" {
		t.Errorf("expected 12.50, got %s", parts[0].StringFixed(2))
	}
	if parts[1].StringFixed(2) != "37.50" {
		t.Errorf("expected 37.50, got %s", parts[1].StringFixed(2))
	}
}

func TestAllocateNoRatios(t *testing.T) {
	m := MustMoneyFromFloat64(10, USD)

	if _, err := m.Allocate(); !errors.Is(err, ErrNoAllocationRatios) {
		t.Errorf("expected ErrNoAllocationRatios, got %v", err)
	}
}

func TestAllocateZeroRatioSum(t *testing.T) {
	m := MustMoneyFromFloat64(10, USD)

	if _, err := m.Allocate(0, 0); !errors.Is(err, ErrZeroAllocationRatios) {
		t.Errorf("expected ErrZeroAllocationRatios, got %v", err)
	}
}

func TestAllocateEvenly(t *testing.T) {
	m := MustMoneyFromFloat64(10, USD)

	parts, err := m.AllocateEvenly(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum := parts[0].Add(parts[1]).Add(parts[2])
	if !sum.Equal(m) {
		t.Errorf("expected parts to sum to %s, got %s", m.String(), sum.String())
	}
}

func TestAllocateEvenlyInvalidCount(t *testing.T) {
	m := MustMoneyFromFloat64(10, USD)

	if _, err := m.AllocateEvenly(0); !errors.Is(err, ErrInvalidAllocationCount) {
		t.Errorf("expected ErrInvalidAllocationCount, got %v", err)
	}
	if _, err := m.AllocateEvenly(-1); !errors.Is(err, ErrInvalidAllocationCount) {
		t.Errorf("expected ErrInvalidAllocationCount, got %v", err)
	}
}

func TestAllocateNegativeAmount(t *testing.T) {
	m := MustMoneyFromFloat64(-10, USD)

	parts, err := m.Allocate(1, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum := parts[0].Add(parts[1]).Add(parts[2])
	if !sum.Equal(m) {
		t.Errorf("expected parts to sum to %s, got %s", m.String(), sum.String())
	}
}
