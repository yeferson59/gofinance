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

func TestCompositeConfigMoneySetters(t *testing.T) {
	present := money.MustMoneyFromFloat64(1000, money.USD)
	future := money.MustMoneyFromFloat64(1500, money.USD)
	rate := money.MustFromFloat64(0.05)

	config := NewComposite().
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

func TestCompositeConfigFutureSetter(t *testing.T) {
	config := NewComposite().Future(1500, money.USD)

	if config.future.Currency() != money.USD {
		t.Fatalf("expected USD currency, got %v", config.future.Currency())
	}
	if config.future.InexactFloat64() != 1500 {
		t.Errorf("expected 1500, got %v", config.future.InexactFloat64())
	}
}

func TestCompositeConfigFrequencyConvenienceMethods(t *testing.T) {
	tests := []struct {
		name     string
		config   CompositeConfig
		expected CompoundingFrequency
	}{
		{"annually", NewComposite().Annually(), Annually},
		{"quarterly", NewComposite().Quarterly(), QuarterlyOne},
		{"daily", NewComposite().Daily(), Daily},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.frequency != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, tt.config.frequency)
			}
		})
	}
}

func TestCompositeConfigMustBuild(t *testing.T) {
	ci := NewComposite().
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

func TestCompositeConfigMustBuildPanicsOnInvalidParams(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid period")
		}
	}()

	NewComposite().Periods(-1).MustBuild()
}

func TestGetEqualsRateInterestPeriodsPropagatesError(t *testing.T) {
	// A zero-value CompositeInterest has an invalid (empty) period frequency,
	// so the error must be surfaced instead of silently returning zero values.
	var ci CompositeInterest
	if _, _, err := ci.GetEqualsRateInterestPeriods(); err == nil {
		t.Error("expected error for invalid period, got nil")
	}
}
