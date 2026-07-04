package annuities

import (
	"math"
	"testing"

	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

func TestYearsRespectsFrequency(t *testing.T) {
	tests := []struct {
		name     string
		config   AnnuityConfig
		expected int
	}{
		{"monthly", NewAnnuity().Monthly().Years(30), 360},
		{"quarterly", NewAnnuity().Quarterly().Years(30), 120},
		{"annually", NewAnnuity().Annually().Years(30), 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.periods != tt.expected {
				t.Errorf("expected %d periods, got %d", tt.expected, tt.config.periods)
			}
		})
	}
}

func TestFutureFromPaymentsOnly(t *testing.T) {
	rate, err := compositeinterest.NewRateInterest(
		money.MustFromFloat64(0.01),
		compositeinterest.Monthly,
		compositeinterest.RateEffectyPeriodic,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	period, err := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	ann, err := New(money.MustMoneyFromFloat64(500, money.USD), money.MoneyZero, money.MoneyZero, period, rate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	future, err := ann.Future()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// FV = PMT × ((1 + r)^n - 1) / r
	expected := 500 * (math.Pow(1.01, 12) - 1) / 0.01
	if math.Abs(future.InexactFloat64()-expected) > 0.01 {
		t.Errorf("expected %.2f, got %s", expected, future.StringFixed(2))
	}
}
