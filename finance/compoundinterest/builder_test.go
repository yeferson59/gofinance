package compoundinterest

import (
	"math"
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestBuilderChainedBuildDefaultsToPeriodicRate(t *testing.T) {
	ci, err := NewCompound().
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

func TestCompoundConfigMoneySetters(t *testing.T) {
	present := money.MustMoneyFromFloat64(1000, money.USD)
	future := money.MustMoneyFromFloat64(1500, money.USD)
	rate := decimal.MustFromFloat64(0.05)

	config := NewCompound().
		PresentMoney(present).
		FutureMoney(future).
		RateMoney(rate).
		Frequency(Monthly).
		RateType(RateEffectyNominal)

	if !config.present.Equal(present) {
		t.Errorf("expected present %v, got %v", present, config.present)
	}
	if !config.future.Equal(future) {
		t.Errorf("expected future %v, got %v", future, config.future)
	}
	if !config.rate.Equal(rate) {
		t.Errorf("expected rate %v, got %v", rate, config.rate)
	}
	if config.frequency != Monthly {
		t.Errorf("expected Monthly, got %v", config.frequency)
	}
	if config.rateType != RateEffectyNominal {
		t.Errorf("expected RateEffectyNominal, got %v", config.rateType)
	}
}

func TestCompoundConfigFutureSetter(t *testing.T) {
	config := NewCompound().Future(1500, money.USD)

	if config.future.GetCurrency() != money.USD {
		t.Fatalf("expected USD currency, got %v", config.future.GetCurrency())
	}
	if config.future.InexactFloat64() != 1500 {
		t.Errorf("expected 1500, got %v", config.future.InexactFloat64())
	}
}

func TestCompoundConfigFrequencyConvenienceMethods(t *testing.T) {
	tests := []struct {
		name     string
		config   CompoundConfig
		expected CompoundingFrequency
	}{
		{"annually", NewCompound().Annually(), Annually},
		{"quarterly", NewCompound().Quarterly(), Quarterly},
		{"daily", NewCompound().Daily(), Daily},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.frequency != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, tt.config.frequency)
			}
		})
	}
}

func TestCompoundConfigMustBuild(t *testing.T) {
	ci := NewCompound().
		Present(1000, money.USD).
		Rate(0.05).
		Periods(12).
		Monthly().
		MustBuild()

	future, err := ci.Future()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := 1000 * math.Pow(1.05, 12)
	if math.Abs(future.InexactFloat64()-expected) > 0.01 {
		t.Errorf("expected %.2f, got %s", expected, future.StringFixed(2))
	}
}

func TestCompoundConfigMustBuildPanicsOnInvalidParams(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid period")
		}
	}()

	NewCompound().Periods(-1).MustBuild()
}

func TestGetEqualsRateInterestPeriodsPropagatesError(t *testing.T) {
	// A zero-value CompoundInterest has an invalid (empty) period frequency,
	// so the error must be surfaced instead of silently returning zero values.
	var ci CompoundInterest
	if _, _, err := ci.GetEqualsRateInterestPeriods(); err == nil {
		t.Error("expected error for invalid period, got nil")
	}
}
