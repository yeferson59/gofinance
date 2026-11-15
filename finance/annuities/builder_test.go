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
			rate, err := tt.config.periodicRate()
			if err != nil {
				t.Fatalf("periodicRate: %v", err)
			}

			if math.Abs(rate.InexactFloat64()-tt.expected) > 1e-9 {
				t.Errorf("expected rate %v, got %v", tt.expected, rate)
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
	if !config.rate.Equal(decimal.MustFromString("0.2668")) {
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

// The periodic rate an annual rate turns into must depend on the frequency the
// builder ends up with, not on whether AnnualRate was called before or after
// the frequency was set.
func TestAnnualRateIsOrderIndependent(t *testing.T) {
	before, err := NewAnnuity().Present(10000, money.USD).Quarterly().AnnualRate(0.08).Periods(8).Payment()
	if err != nil {
		t.Fatalf("frequency first: %v", err)
	}

	after, err := NewAnnuity().Present(10000, money.USD).AnnualRate(0.08).Quarterly().Periods(8).Payment()
	if err != nil {
		t.Fatalf("rate first: %v", err)
	}

	if !before.Equal(after) {
		t.Fatalf("payment depends on call order: frequency first %s, rate first %s", before, after)
	}

	// 10,000 at 2% per quarter over 8 quarters.
	if got := before.RoundBank(2).StringFixed(2); got != "1365.10" {
		t.Errorf("expected 1365.10, got %s", got)
	}
}

// Rate and EffectiveAnnualRate replace an earlier AnnualRate entirely, so the
// later call is never divided by the periods per year.
func TestLaterRateSetterReplacesAnnualRate(t *testing.T) {
	periodic, err := NewAnnuity().Present(10000, money.USD).Quarterly().Rate(0.02).Periods(8).Payment()
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}

	replaced, err := NewAnnuity().Present(10000, money.USD).Quarterly().AnnualRate(0.5).Rate(0.02).Periods(8).Payment()
	if err != nil {
		t.Fatalf("AnnualRate then Rate: %v", err)
	}

	if !periodic.Equal(replaced) {
		t.Errorf("Rate did not replace AnnualRate: got %s, want %s", replaced, periodic)
	}

	effective, err := NewAnnuity().Present(10000, money.USD).EffectiveAnnualRate(0.2668).Periods(12).Payment()
	if err != nil {
		t.Fatalf("EffectiveAnnualRate: %v", err)
	}

	replacedEffective, err := NewAnnuity().Present(10000, money.USD).AnnualRate(0.5).EffectiveAnnualRate(0.2668).Periods(12).Payment()
	if err != nil {
		t.Fatalf("AnnualRate then EffectiveAnnualRate: %v", err)
	}

	if !effective.Equal(replacedEffective) {
		t.Errorf("EffectiveAnnualRate did not replace AnnualRate: got %s, want %s", replacedEffective, effective)
	}
}

// Methods that return an error must report bad input through it instead of
// panicking.
func TestBuilderReportsInvalidInputAsError(t *testing.T) {
	tests := []struct {
		name   string
		config AnnuityConfig
	}{
		{"NaN present", NewAnnuity().Present(math.NaN(), money.USD).Rate(0.01).Periods(12)},
		{"infinite value", NewAnnuity().Value(math.Inf(1), money.USD).Rate(0.01).Periods(12)},
		{"infinite future", NewAnnuity().Future(math.Inf(-1), money.USD).Rate(0.01).Periods(12)},
		{"NaN rate", NewAnnuity().Present(1000, money.USD).Rate(math.NaN()).Periods(12)},
		{"NaN annual rate", NewAnnuity().Present(1000, money.USD).AnnualRate(math.NaN()).Periods(12)},
		{"NaN effective rate", NewAnnuity().Present(1000, money.USD).EffectiveAnnualRate(math.NaN()).Periods(12)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked: %v", r)
				}
			}()

			if _, err := tt.config.Build(); err == nil {
				t.Error("Build: expected an error")
			}

			if _, err := tt.config.Payment(); err == nil {
				t.Error("Payment: expected an error")
			}
		})
	}
}

// An annual rate whose periodic share has more digits than a float64 prints
// in 19 decimals used to panic inside Build.
func TestBuilderAcceptsSmallAnnualRate(t *testing.T) {
	payment, err := NewAnnuity().Present(10000, money.USD).AnnualRate(0.001).Periods(8).Monthly().Payment()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := payment.RoundBank(2).StringFixed(2); got != "1250.47" {
		t.Errorf("expected 1250.47, got %s", got)
	}
}
