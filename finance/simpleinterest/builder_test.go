package simpleinterest

import (
	"testing"

	"github.com/yeferson59/gofinance/money"
)

func TestBuilderChainedBuild(t *testing.T) {
	si := NewSimple().
		Present(5000, money.USD).
		AnnualRate(0.12).
		Periods(18).
		Months().
		Build()

	future, err := si.FutureWithRateInterest()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if future.StringFixed(2) != "5900.00" {
		t.Errorf("expected 5900.00, got %s", future.StringFixed(2))
	}
}

func TestBuilderFutureValue(t *testing.T) {
	future, err := NewSimple().
		Present(5000, money.USD).
		AnnualRate(0.12).
		Periods(18).
		Months().
		FutureValue()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if future.StringFixed(2) != "5900.00" {
		t.Errorf("expected 5900.00, got %s", future.StringFixed(2))
	}
}

func TestBuilderPresentValue(t *testing.T) {
	present, err := NewSimple().
		Future(5900, money.USD).
		AnnualRate(0.12).
		Periods(18).
		Months().
		PresentValue()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if present.StringFixed(2) != "5000.00" {
		t.Errorf("expected 5000.00, got %s", present.StringFixed(2))
	}
}
