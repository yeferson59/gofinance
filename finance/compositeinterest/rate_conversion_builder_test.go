package compositeinterest

import (
	"math"
	"testing"

	"github.com/yeferson59/gofinance/money"
)

func assertRate(t *testing.T, got money.Decimal, err error, expected float64) {
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

func TestRateConversionNegativeRate(t *testing.T) {
	if _, err := NewRateConversion().Rate(-0.05).ToPeriodic(); err == nil {
		t.Error("expected error for negative rate, got nil")
	}
}
