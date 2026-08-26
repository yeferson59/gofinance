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
			if got := tt.config.periodicRate().InexactFloat64(); math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("expected rate %v, got %v", tt.expected, got)
			}
		})
	}
}

// TestAnnualRateIsOrderIndependent checks that the period type is honoured
// whichever order the builder methods are called in.
//
// AnnualRate used to divide the moment it was called, reading whatever period
// type had been set so far — Months by default. Setting the rate before the
// period type therefore produced a monthly rate silently, so
// AnnualRate(0.06).Years() charged 0.5% a year instead of 6%, with no error.
// The conversion now happens at Build, from the period type the builder ends
// up with.
func TestAnnualRateIsOrderIndependent(t *testing.T) {
	pairs := []struct {
		name     string
		before   SimpleConfig
		after    SimpleConfig
		expected float64
	}{
		{
			"years", NewSimple().Years().AnnualRate(0.06), NewSimple().AnnualRate(0.06).Years(), 0.06,
		},
		{
			"days", NewSimple().Days().AnnualRate(0.365), NewSimple().AnnualRate(0.365).Days(), 0.001,
		},
		{
			"weeks", NewSimple().Weeks().AnnualRate(0.52), NewSimple().AnnualRate(0.52).Weeks(), 0.01,
		},
		{
			"months", NewSimple().Months().AnnualRate(0.12), NewSimple().AnnualRate(0.12).Months(), 0.01,
		},
	}

	for _, pair := range pairs {
		t.Run(pair.name, func(t *testing.T) {
			before := pair.before.periodicRate().InexactFloat64()
			after := pair.after.periodicRate().InexactFloat64()

			if math.Abs(before-after) > 1e-9 {
				t.Errorf("order changed the rate: %v before the period type, %v after", before, after)
			}

			if math.Abs(before-pair.expected) > 1e-9 {
				t.Errorf("expected rate %v, got %v", pair.expected, before)
			}
		})
	}
}

// TestRateOverridesAnnualRate checks the two setters replace each other rather
// than compounding.
func TestRateOverridesAnnualRate(t *testing.T) {
	// A periodic rate set last must be used as given, not divided again.
	periodic := NewSimple().Years().AnnualRate(0.12).Rate(0.02).periodicRate()
	if got := periodic.InexactFloat64(); math.Abs(got-0.02) > 1e-9 {
		t.Errorf("expected 0.02, got %v", got)
	}

	// And an annual rate set last must be converted.
	annual := NewSimple().Months().Rate(0.02).AnnualRate(0.12).periodicRate()
	if got := annual.InexactFloat64(); math.Abs(got-0.01) > 1e-9 {
		t.Errorf("expected 0.01, got %v", got)
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

	if config.interest.GetCurrency() != money.USD {
		t.Fatalf("expected USD currency, got %v", config.interest.GetCurrency())
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
