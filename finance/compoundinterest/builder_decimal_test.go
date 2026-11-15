package compoundinterest

import (
	"errors"
	"math"
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestRateDecimalMatchesRate(t *testing.T) {
	want, err := NewCompound().Present(1000, money.USD).Rate(0.01).Periods(12).Monthly().MustBuild().Future()
	if err != nil {
		t.Fatal(err)
	}

	got, err := NewCompound().PresentMoney(money.MustMoneyFromString("1000", money.USD)).RateDecimal(decimal.MustFromString("0.01")).Periods(12).Monthly().MustBuild().Future()
	if err != nil {
		t.Fatal(err)
	}

	if !want.Equal(got) {
		t.Errorf("Rate gave %s, RateDecimal gave %s", want, got)
	}
}

// Setters given a float with no decimal value report it through Build instead
// of panicking.
func TestCompoundBuilderReportsInvalidInputAsError(t *testing.T) {
	configs := map[string]CompoundConfig{
		"present": NewCompound().Present(math.NaN(), money.USD).Rate(0.01).Periods(12),
		"future":  NewCompound().Future(math.Inf(1), money.USD).Rate(0.01).Periods(12),
		"rate":    NewCompound().Present(1000, money.USD).Rate(math.NaN()).Periods(12),
	}

	for name, config := range configs {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked: %v", r)
				}
			}()

			if _, err := config.Build(); !errors.Is(err, decimal.ErrInvalidFormat) {
				t.Errorf("Build error = %v, want ErrInvalidFormat", err)
			}
		})
	}
}

func TestRateConversionReportsInvalidInputAsError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()

	if _, err := NewRateConversion().Rate(math.NaN()).ToEffectiveAnnual(); !errors.Is(err, decimal.ErrInvalidFormat) {
		t.Errorf("ToEffectiveAnnual error = %v, want ErrInvalidFormat", err)
	}
}

// The correctly spelled rate types are the same values as the old ones.
func TestRateTypeAliases(t *testing.T) {
	pairs := map[TypeRate]TypeRate{
		RateEffectivePeriodic:           RateEffectyPeriodic,
		RateEffectiveAnnually:           RateEffectyAnnually,
		RateEffectiveNominal:            RateEffectyNominal,
		RateAnticipateEffectivePeriodic: RateAnticipateEffectyPeriodic,
		RateAnticipateEffectiveNominal:  RateAnticipateEffectyNominal,
	}

	for current, deprecated := range pairs {
		if current != deprecated {
			t.Errorf("%q != %q", current, deprecated)
		}
	}
}
