package simpleinterest

import (
	"errors"
	"math"
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// invalidInputConfigs returns builders that each received one float with no
// decimal value, alongside otherwise valid terms.
func invalidInputConfigs() map[string]SimpleConfig {
	nan := math.NaN()

	return map[string]SimpleConfig{
		"present":     NewSimple().Present(nan, money.USD).Future(6000, money.USD).AnnualRate(0.12).Periods(18),
		"future":      NewSimple().Present(5000, money.USD).Future(math.Inf(1), money.USD).AnnualRate(0.12).Periods(18),
		"interest":    NewSimple().Present(5000, money.USD).Interest(math.Inf(-1), money.USD).AnnualRate(0.12).Periods(18),
		"rate":        NewSimple().Present(5000, money.USD).Future(6000, money.USD).Rate(nan).Periods(18),
		"annual rate": NewSimple().Present(5000, money.USD).Future(6000, money.USD).AnnualRate(math.Inf(1)).Periods(18),
	}
}

// Setters given a NaN or infinite float must not panic: the builder's
// error-returning shortcuts report the problem instead.
func TestBuilderReportsInvalidInputAsError(t *testing.T) {
	for name, config := range invalidInputConfigs() {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked: %v", r)
				}
			}()

			if _, err := config.FutureValue(); !errors.Is(err, decimal.ErrInvalidFormat) {
				t.Errorf("FutureValue error = %v, want ErrInvalidFormat", err)
			}

			if _, err := config.PresentValue(); !errors.Is(err, decimal.ErrInvalidFormat) {
				t.Errorf("PresentValue error = %v, want ErrInvalidFormat", err)
			}
		})
	}
}

// Build has no error to return, so the SimpleInterest it produces carries the
// input error and every calculation on it reports it rather than computing
// from a zeroed field.
func TestBuiltValueCarriesInputError(t *testing.T) {
	for name, config := range invalidInputConfigs() {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked: %v", r)
				}
			}()

			si := config.Build()

			calculations := map[string]func() error{
				"Future":                           func() error { _, err := si.Future(); return err },
				"FutureWithRateInterest":           func() error { _, err := si.FutureWithRateInterest(); return err },
				"Present":                          func() error { _, err := si.Present(); return err },
				"PresentWithFuture":                func() error { _, err := si.PresentWithFuture(); return err },
				"Interest":                         func() error { _, err := si.Interest(); return err },
				"InterestWithPresentAndFuture":     func() error { _, err := si.InterestWithPresentAndFuture(); return err },
				"RateInterest":                     func() error { _, err := si.RateInterest(); return err },
				"RateInterestWithPresentAndFuture": func() error { _, err := si.RateInterestWithPresentAndFuture(); return err },
				"Periods":                          func() error { _, err := si.Periods(); return err },
				"PeriodsWithPresentAndFuture":      func() error { _, err := si.PeriodsWithPresentAndFuture(); return err },
				"GetPeriods":                       func() error { _, err := si.GetPeriods(); return err },
			}

			for method, calc := range calculations {
				if err := calc(); !errors.Is(err, decimal.ErrInvalidFormat) {
					t.Errorf("%s error = %v, want ErrInvalidFormat", method, err)
				}
			}
		})
	}
}

// The first invalid input is the one reported, even if later setters
// overwrite the field it was given for.
func TestBuilderKeepsFirstInputError(t *testing.T) {
	_, err := NewSimple().
		Present(math.NaN(), money.USD).
		Present(5000, money.USD).
		AnnualRate(0.12).
		Periods(18).
		FutureValue()
	if !errors.Is(err, decimal.ErrInvalidFormat) {
		t.Fatalf("FutureValue error = %v, want ErrInvalidFormat", err)
	}
}

// A valid builder is unaffected: the same terms still compute.
func TestBuilderValidInputStillComputes(t *testing.T) {
	future, err := NewSimple().Present(5000, money.USD).AnnualRate(0.12).Periods(18).Months().FutureValue()
	if err != nil {
		t.Fatalf("FutureValue: %v", err)
	}

	if got := future.StringFixed(2); got != "5900.00" {
		t.Errorf("expected 5900.00, got %s", got)
	}
}

// A rate large enough to overflow the growth factor is reported as an error
// by the error-returning shortcut instead of panicking inside the decimal
// engine's Mul.
func TestBuilderReportsGrowthOverflow(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()

	_, err := NewSimple().Present(5000, money.USD).AnnualRate(1e30).Periods(1_000_000_000).Years().FutureValue()
	if !errors.Is(err, decimal.ErrOverflow) {
		t.Fatalf("FutureValue error = %v, want ErrOverflow", err)
	}
}
