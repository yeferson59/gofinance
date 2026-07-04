package compositeinterest

import (
	"math"
	"testing"

	"github.com/yeferson59/gofinance/money"
)

func TestBuilderChainedBuildDefaultsToPeriodicRate(t *testing.T) {
	ci, err := NewComposite().
		Present(1000, money.USD).
		Rate(0.05).
		Periods(12).
		Monthly().
		Build()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	future, err := ci.Future()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := 1000 * math.Pow(1.05, 12)
	if math.Abs(future.InexactFloat64()-expected) > 0.01 {
		t.Errorf("expected %.2f, got %s", expected, future.StringFixed(2))
	}
}

func TestGetEqualsRateInterestPeriodsPropagatesError(t *testing.T) {
	// A zero-value CompositeInterest has an invalid (empty) period frequency,
	// so the error must be surfaced instead of silently returning zero values.
	var ci CompositeInterest
	if _, _, err := ci.GetEqualsRateInterestPeriods(); err == nil {
		t.Error("expected error for invalid period, got nil")
	}
}
