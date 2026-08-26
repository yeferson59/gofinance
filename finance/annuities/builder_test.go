package annuities

import (
	"math"
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// withFrequency directly sets AnnuityConfig's private frequency field, since
// (unlike CompoundConfig) AnnuityConfig has no public Frequency() setter —
// only the Monthly/Quarterly/Annually convenience methods.
func withFrequency(a AnnuityConfig, f compoundinterest.CompoundingFrequency) AnnuityConfig {
	a.frequency = f
	return a
}

func TestYearsRespectsFrequency(t *testing.T) {
	tests := []struct {
		name     string
		config   AnnuityConfig
		expected int
	}{
		{"monthly", NewAnnuity().Monthly().Years(30), 360},
		{"quarterly", NewAnnuity().Quarterly().Years(30), 120},
		{"annually", NewAnnuity().Annually().Years(30), 30},
		{"daily", withFrequency(NewAnnuity(), compoundinterest.Daily).Years(1), 365},
		{"bimonthly", withFrequency(NewAnnuity(), compoundinterest.Bimonthly).Years(1), 6},
		{"quarterlyTwo", withFrequency(NewAnnuity(), compoundinterest.FourMonthly).Years(1), 3},
		{"semiAnnually", withFrequency(NewAnnuity(), compoundinterest.SemiAnnually).Years(1), 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.periods != tt.expected {
				t.Errorf("expected %d periods, got %d", tt.expected, tt.config.periods)
			}
		})
	}
}

func TestAnnualRateRespectsFrequency(t *testing.T) {
	tests := []struct {
		name     string
		config   AnnuityConfig
		expected float64
	}{
		{"monthly", NewAnnuity().Monthly().AnnualRate(0.12), 0.01},
		{"daily", withFrequency(NewAnnuity(), compoundinterest.Daily).AnnualRate(0.365), 0.001},
		{"bimonthly", withFrequency(NewAnnuity(), compoundinterest.Bimonthly).AnnualRate(0.06), 0.01},
		{"quarterlyOne", withFrequency(NewAnnuity(), compoundinterest.Quarterly).AnnualRate(0.04), 0.01},
		{"quarterlyTwo", withFrequency(NewAnnuity(), compoundinterest.FourMonthly).AnnualRate(0.03), 0.01},
		{"semiAnnually", withFrequency(NewAnnuity(), compoundinterest.SemiAnnually).AnnualRate(0.02), 0.01},
		{"annually", withFrequency(NewAnnuity(), compoundinterest.Annually).AnnualRate(0.01), 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if math.Abs(tt.config.rate-tt.expected) > 1e-9 {
				t.Errorf("expected rate %v, got %v", tt.expected, tt.config.rate)
			}
		})
	}
}

func TestAnnuityConfigFutureSetter(t *testing.T) {
	config := NewAnnuity().Future(50000, money.USD)

	if config.future.GetCurrency() != money.USD {
		t.Fatalf("expected USD currency, got %v", config.future.GetCurrency())
	}
	if config.future.InexactFloat64() != 50000 {
		t.Errorf("expected 50000, got %v", config.future.InexactFloat64())
	}
}

func TestAnnuityConfigEffectiveAnnualRate(t *testing.T) {
	config := NewAnnuity().Monthly().EffectiveAnnualRate(0.2668)

	if config.rateType != compoundinterest.RateEffectyAnnually {
		t.Errorf("expected RateEffectyAnnually, got %v", config.rateType)
	}
	if config.rate != 0.2668 {
		t.Errorf("expected 0.2668, got %v", config.rate)
	}

	payment, err := config.
		Present(300000, money.USD).
		Periods(360).
		Payment()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !payment.IsPositive() {
		t.Errorf("expected a positive payment, got %s", payment.String())
	}
}

func TestAnnuityConfigMustBuild(t *testing.T) {
	annuity := NewAnnuity().
		Present(300000, money.USD).
		AnnualRate(0.06).
		Periods(360).
		Monthly().
		MustBuild()

	payment, err := annuity.PaymentFromPresentValue()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !payment.IsPositive() {
		t.Errorf("expected a positive payment, got %s", payment.String())
	}
}

func TestAnnuityConfigMustBuildPanicsOnInvalidParams(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid frequency/rate combination")
		}
	}()

	NewAnnuity().Periods(-1).MustBuild()
}

func TestFutureFromPaymentsOnly(t *testing.T) {
	rate, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.01),
		compoundinterest.Monthly,
		compoundinterest.RateEffectyPeriodic,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
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
