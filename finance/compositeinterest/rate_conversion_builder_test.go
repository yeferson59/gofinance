package compositeinterest

import (
	"math"
	"testing"

	"github.com/yeferson59/gofinance/decimal"
)

func assertRate(t *testing.T, got decimal.Decimal, err error, expected float64) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if math.Abs(got.InexactFloat64()-expected) > 1e-9 {
		t.Errorf("expected %.10f, got %s", expected, got.String())
	}
}

func TestRateConversionNominalToPeriodic(t *testing.T) {
	periodic, err := NewRateConversion().
		Rate(0.12).
		Nominal().
		Monthly().
		ToPeriodic()

	assertRate(t, periodic, err, 0.01)
}

func TestRateConversionPeriodicToNominal(t *testing.T) {
	nominal, err := NewRateConversion().
		Rate(0.01).
		Periodic().
		Monthly().
		ToNominal()

	assertRate(t, nominal, err, 0.12)
}

func TestRateConversionPeriodicToEffectiveAnnual(t *testing.T) {
	effective, err := NewRateConversion().
		Rate(0.01).
		Periodic().
		Monthly().
		ToEffectiveAnnual()

	assertRate(t, effective, err, math.Pow(1.01, 12)-1)
}

func TestRateConversionEffectiveAnnualToPeriodic(t *testing.T) {
	effectiveAnnual := math.Pow(1.01, 12) - 1

	periodic, err := NewRateConversion().
		Rate(effectiveAnnual).
		EffectiveAnnual().
		Monthly().
		ToPeriodic()

	assertRate(t, periodic, err, 0.01)
}

func TestRateConversionToPeriodicAt(t *testing.T) {
	quarterly, err := NewRateConversion().
		Rate(0.01).
		Periodic().
		Monthly().
		ToPeriodicAt(QuarterlyOne)

	assertRate(t, quarterly, err, math.Pow(1.01, 3)-1)
}

func TestRateConversionToNominalAt(t *testing.T) {
	nominalQuarterly, err := NewRateConversion().
		Rate(0.01).
		Periodic().
		Monthly().
		ToNominalAt(QuarterlyOne)

	assertRate(t, nominalQuarterly, err, (math.Pow(1.01, 3)-1)*4)
}

func TestRateConversionAnticipatedPeriodicToEffectiveAnnual(t *testing.T) {
	effective, err := NewRateConversion().
		Rate(0.01).
		AnticipatedPeriodic().
		Monthly().
		ToEffectiveAnnual()

	assertRate(t, effective, err, math.Pow(0.99, -12)-1)
}

func TestRateConversionPeriodicToAnticipatedPeriodic(t *testing.T) {
	anticipated, err := NewRateConversion().
		Rate(0.01).
		Periodic().
		Monthly().
		ToAnticipatedPeriodic()

	// d = 1 - (1 + i)^-1 for the same frequency
	assertRate(t, anticipated, err, 1-1/1.01)
}

func TestRateConversionBuild(t *testing.T) {
	rt, err := NewRateConversion().
		Rate(0.12).
		Nominal().
		Monthly().
		Build()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	periodic, err := rt.RatePeriodic()
	assertRate(t, periodic, err, 0.01)
}

func TestRateConversionConfigSetters(t *testing.T) {
	rate := decimal.MustFromFloat64(0.03)

	config := NewRateConversion().
		RateDecimal(rate).
		Frequency(QuarterlyOne).
		RateType(RateEffectyNominal)

	if !config.rate.Equal(rate) {
		t.Errorf("expected rate %v, got %v", rate, config.rate)
	}
	if config.frequency != QuarterlyOne {
		t.Errorf("expected QuarterlyOne, got %v", config.frequency)
	}
	if config.rateType != RateEffectyNominal {
		t.Errorf("expected RateEffectyNominal, got %v", config.rateType)
	}
}

func TestRateConversionFrequencyConvenienceMethods(t *testing.T) {
	tests := []struct {
		name     string
		config   RateConversionConfig
		expected CompoundingFrequency
	}{
		{"daily", NewRateConversion().Daily(), Daily},
		{"quarterly", NewRateConversion().Quarterly(), QuarterlyOne},
		{"semiannually", NewRateConversion().SemiAnnually(), SemiAnnually},
		{"annually", NewRateConversion().Annually(), Annually},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.frequency != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, tt.config.frequency)
			}
		})
	}
}

func TestRateConversionAnticipatedNominal(t *testing.T) {
	config := NewRateConversion().Rate(0.01).AnticipatedNominal().Monthly()

	if config.rateType != RateAnticipateEffectyNominal {
		t.Errorf("expected RateAnticipateEffectyNominal, got %v", config.rateType)
	}
	if !config.isAnticipated() {
		t.Error("expected isAnticipated() to be true")
	}
}

func TestRateConversionAnticipatedNominalRoundTrip(t *testing.T) {
	// An anticipated nominal rate of 12×(1 - 1/1.01) corresponds to an
	// ordinary periodic rate of 0.01 monthly, i.e. an ordinary nominal
	// rate of 0.12.
	anticipatedNominal := 12 * (1 - 1/1.01)

	nominal, err := NewRateConversion().
		Rate(anticipatedNominal).
		AnticipatedNominal().
		Monthly().
		ToNominal()

	assertRate(t, nominal, err, 0.12)
}

func TestRateConversionToAnticipatedNominal(t *testing.T) {
	// From an ordinary periodic rate of 0.01 monthly, the anticipated
	// nominal rate is 12×(1 - 1/1.01).
	anticipatedNominal, err := NewRateConversion().
		Rate(0.01).
		Periodic().
		Monthly().
		ToAnticipatedNominal()

	assertRate(t, anticipatedNominal, err, 12*(1-1/1.01))
}

func TestRateConversionConfigMustBuild(t *testing.T) {
	rt := NewRateConversion().
		Rate(0.12).
		Nominal().
		Monthly().
		MustBuild()

	periodic, err := rt.RatePeriodic()
	assertRate(t, periodic, err, 0.01)
}

func TestRateConversionConfigMustBuildPanicsOnInvalidRate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for negative rate")
		}
	}()

	NewRateConversion().Rate(-0.05).MustBuild()
}

func TestRateConversionNegativeRate(t *testing.T) {
	if _, err := NewRateConversion().Rate(-0.05).ToPeriodic(); err == nil {
		t.Error("expected error for negative rate, got nil")
	}
}
