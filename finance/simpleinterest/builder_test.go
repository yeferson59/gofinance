package simpleinterest

import (
	"math"
	"testing"

	"github.com/yeferson59/gofinance/v2/money"
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

func TestAnnualRateRespectsPeriodType(t *testing.T) {
	tests := []struct {
		name     string
		config   SimpleConfig
		expected float64
	}{
		{"months", NewSimple().Months().AnnualRate(0.12), 0.01},
		{"days", NewSimple().Days().AnnualRate(0.365), 0.001},
		{"weeks", NewSimple().Weeks().AnnualRate(0.52), 0.01},
		{"years", NewSimple().Years().AnnualRate(0.05), 0.05},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.rate.InexactFloat64(); math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("expected rate %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestSimpleConfigRateSetter(t *testing.T) {
	config := NewSimple().Rate(0.01)

	if config.rate.InexactFloat64() != 0.01 {
		t.Errorf("expected 0.01, got %v", config.rate.InexactFloat64())
	}
}

func TestSimpleConfigInterestSetter(t *testing.T) {
	config := NewSimple().Interest(900, money.USD)

	if config.interest.Currency() != money.USD {
		t.Fatalf("expected USD currency, got %v", config.interest.Currency())
	}
	if config.interest.InexactFloat64() != 900 {
		t.Errorf("expected 900, got %v", config.interest.InexactFloat64())
	}
}

func TestSimpleConfigPeriodTypeSetter(t *testing.T) {
	config := NewSimple().PeriodType(Weeks)

	if config.periodType != Weeks {
		t.Errorf("expected Weeks, got %v", config.periodType)
	}
}

func TestSimpleConfigPeriodTypeConvenienceMethods(t *testing.T) {
	tests := []struct {
		name     string
		config   SimpleConfig
		expected Periods
	}{
		{"years", NewSimple().Years(), Years},
		{"days", NewSimple().Days(), Days},
		{"weeks", NewSimple().Weeks(), Weeks},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.periodType != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, tt.config.periodType)
			}
		})
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
