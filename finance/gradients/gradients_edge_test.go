package gradients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// newSeries builds the period and rate shared by the tests below.
func newSeries(t *testing.T, periods, rate float64) (compoundinterest.Period, compoundinterest.RateInterest) {
	t.Helper()

	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(periods), compoundinterest.Annually)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(rate), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	return period, rateInterest
}

// TestArithmeticZeroRate is the regression test for TESTING_PLAN.md §2.3:
// with a 0% rate the general formula divides by the rate, so both Present and
// Future used to fail on a legitimate input. The limit exists and is the plain
// sum of the payments.
func TestArithmeticZeroRate(t *testing.T) {
	// A = 1000, G = 100, n = 5 -> payments 1000, 1100, 1200, 1300, 1400.
	// Sum = 1000×5 + 100×(5×4/2) = 5000 + 1000 = 6000.
	period, rate := newSeries(t, 5, 0)

	series, err := NewArithmetic(
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(100, money.USD),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 6000.0, present.InexactFloat64(), 1e-9)

	future, err := series.Future()
	require.NoError(t, err)
	assert.InDelta(t, 6000.0, future.InexactFloat64(), 1e-9)
}

// TestArithmeticZeroRateDecreasingSeries checks the same limit with a negative
// gradient, where the payments shrink each period.
func TestArithmeticZeroRateDecreasingSeries(t *testing.T) {
	// A = 1000, G = -100, n = 5 -> 1000, 900, 800, 700, 600 = 4000.
	period, rate := newSeries(t, 5, 0)

	series, err := NewArithmetic(
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(-100, money.USD),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 4000.0, present.InexactFloat64(), 1e-9)
}

// TestGeometricZeroRate pins the zero-rate behaviour of the geometric series.
// Unlike the arithmetic one it never divided by the rate — the denominator is
// (i−g) — so this guards the existing behaviour rather than fixing it.
func TestGeometricZeroRate(t *testing.T) {
	// A = 1000, g = 8%, i = 0, n = 5 -> 1000×(1.08^5 − 1)/0.08 = 5866.6010
	period, rate := newSeries(t, 5, 0)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.MustFromFloat64(0.08),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 5866.6010, present.InexactFloat64(), 0.001)

	future, err := series.Future()
	require.NoError(t, err)
	assert.InDelta(t, 5866.6010, future.InexactFloat64(), 0.001)
}

// TestGeometricZeroRateAndZeroGrowth covers the corner where both the rate and
// the growth are zero, so the g == i branch is taken with i = 0.
func TestGeometricZeroRateAndZeroGrowth(t *testing.T) {
	period, rate := newSeries(t, 5, 0)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.Zero,
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 5000.0, present.InexactFloat64(), 1e-9)

	future, err := series.Future()
	require.NoError(t, err)
	assert.InDelta(t, 5000.0, future.InexactFloat64(), 1e-9)
}

// TestArithmeticMatchesDirectDiscounting checks the closed-form present value
// against the definition: discounting each payment one by one. The formula is
// a summation identity, so for a small series the two must agree exactly.
func TestArithmeticMatchesDirectDiscounting(t *testing.T) {
	const (
		first     = 1000.0
		gradient  = 100.0
		rateValue = 0.10
		periods   = 5
	)

	period, rate := newSeries(t, periods, rateValue)

	series, err := NewArithmetic(
		money.MustMoneyFromFloat64(first, money.USD),
		money.MustMoneyFromFloat64(gradient, money.USD),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)

	// PV = Σ (A + t×G) / (1+i)^(t+1) for t = 0 … n−1
	expected := 0.0
	discount := 1.0

	for t := range periods {
		discount *= 1 + rateValue
		expected += (first + float64(t)*gradient) / discount
	}

	assert.InDelta(t, expected, present.InexactFloat64(), 1e-6)
}

// TestGeometricMatchesDirectDiscounting is the same check for the geometric
// series: each payment grows by a constant percentage.
func TestGeometricMatchesDirectDiscounting(t *testing.T) {
	const (
		first     = 1000.0
		growth    = 0.08
		rateValue = 0.10
		periods   = 5
	)

	period, rate := newSeries(t, periods, rateValue)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(first, money.USD),
		decimal.MustFromFloat64(growth),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)

	// PV = Σ A×(1+g)^t / (1+i)^(t+1) for t = 0 … n−1
	expected := 0.0
	payment := first
	discount := 1.0

	for range periods {
		discount *= 1 + rateValue
		expected += payment / discount
		payment *= 1 + growth
	}

	assert.InDelta(t, expected, present.InexactFloat64(), 1e-6)
}

// TestGeometricGrowthAboveRate covers a series growing faster than the
// discount rate: the present value stays finite for a finite term, and each
// payment is worth more than the last in today's money.
func TestGeometricGrowthAboveRate(t *testing.T) {
	// A = 1000, g = 15%, i = 10%, n = 5
	// (1.15/1.10)^5 = (23/22)^5 = 6436343/5153632 = 1.2488945
	// PV = 1000 × (1 − 1.2488945) / (0.10 − 0.15) = 4977.8905
	period, rate := newSeries(t, 5, 0.10)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.MustFromFloat64(0.15),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 4977.8905, present.InexactFloat64(), 0.001)

	// Growing faster than the discount rate must beat a level annuity of the
	// same first payment.
	level, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.Zero,
		period, rate)
	require.NoError(t, err)

	levelPresent, err := level.Present()
	require.NoError(t, err)
	assert.Greater(t, present.InexactFloat64(), levelPresent.InexactFloat64())
}

// TestGeometricShrinkingSeries covers a negative growth rate.
func TestGeometricShrinkingSeries(t *testing.T) {
	// A = 1000, g = −10%, i = 10%, n = 5
	// (0.90/1.10)^5 = (9/11)^5 = 59049/161051 = 0.3666478
	// PV = 1000 × (1 − 0.3666478) / (0.10 + 0.10) = 3166.7608
	period, rate := newSeries(t, 5, 0.10)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.MustFromFloat64(-0.10),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.InDelta(t, 3166.7608, present.InexactFloat64(), 0.001)
}

// TestSinglePeriodSeries checks the degenerate one-payment series: with a
// single payment the gradient and growth never apply, so both series are just
// that payment discounted one period.
func TestSinglePeriodSeries(t *testing.T) {
	period, rate := newSeries(t, 1, 0.10)

	arithmetic, err := NewArithmetic(
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(500, money.USD),
		period, rate)
	require.NoError(t, err)

	arithmeticPresent, err := arithmetic.Present()
	require.NoError(t, err)
	assert.InDelta(t, 1000/1.10, arithmeticPresent.InexactFloat64(), 1e-6)

	geometric, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.MustFromFloat64(0.50),
		period, rate)
	require.NoError(t, err)

	geometricPresent, err := geometric.Present()
	require.NoError(t, err)
	assert.InDelta(t, 1000/1.10, geometricPresent.InexactFloat64(), 1e-6)
}

// TestGeometricCurrencyIsPreserved checks the result carries the first
// payment's currency rather than defaulting to USD.
func TestGeometricCurrencyIsPreserved(t *testing.T) {
	period, rate := newSeries(t, 5, 0.10)

	series, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.EUR),
		decimal.MustFromFloat64(0.08),
		period, rate)
	require.NoError(t, err)

	present, err := series.Present()
	require.NoError(t, err)
	assert.Equal(t, money.EUR, present.GetCurrency())

	future, err := series.Future()
	require.NoError(t, err)
	assert.Equal(t, money.EUR, future.GetCurrency())
}

// TestBuilderAnnualRateSpreadsAcrossFrequency covers AnnualRate, which divides
// the annual rate by the configured frequency's periods per year. The
// frequency must be set before the rate for the division to use it.
func TestBuilderAnnualRateSpreadsAcrossFrequency(t *testing.T) {
	tests := []struct {
		name     string
		build    func(float64) ArithmeticConfig
		expected float64
	}{
		{"monthly", func(r float64) ArithmeticConfig {
			return NewArithmeticSeries().Monthly().AnnualRate(r)
		}, 0.12 / 12},
		{"quarterly", func(r float64) ArithmeticConfig {
			return NewArithmeticSeries().Quarterly().AnnualRate(r)
		}, 0.12 / 4},
		{"annually", func(r float64) ArithmeticConfig {
			return NewArithmeticSeries().Annually().AnnualRate(r)
		}, 0.12},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// One period at the periodic rate: PV = A/(1+i), so the rate the
			// builder derived is observable in the result.
			present, err := test.build(0.12).
				FirstPayment(1000, money.USD).
				Gradient(0, money.USD).
				Periods(1).
				Present()
			require.NoError(t, err)

			assert.InDelta(t, 1000/(1+test.expected), present.InexactFloat64(), 1e-6)
		})
	}
}

// TestBuilderAnnualRateIsOrderIndependent checks the frequency is honoured
// whichever order the builder methods are called in.
//
// AnnualRate used to divide the moment it was called, reading whatever
// frequency had been set so far — Monthly by default. Setting the rate before
// the frequency therefore produced a monthly rate silently, so
// AnnualRate(0.12).Annually() charged 1% a year instead of 12%, with no error.
func TestBuilderAnnualRateIsOrderIndependent(t *testing.T) {
	pairs := []struct {
		name   string
		before func() (money.Money, error)
		after  func() (money.Money, error)
	}{
		{
			"arithmetic annually",
			func() (money.Money, error) {
				return NewArithmeticSeries().Annually().AnnualRate(0.12).
					FirstPayment(1000, money.USD).Gradient(0, money.USD).Periods(5).Present()
			},
			func() (money.Money, error) {
				return NewArithmeticSeries().AnnualRate(0.12).Annually().
					FirstPayment(1000, money.USD).Gradient(0, money.USD).Periods(5).Present()
			},
		},
		{
			"geometric quarterly",
			func() (money.Money, error) {
				return NewGeometricSeries().Quarterly().AnnualRate(0.12).
					FirstPayment(1000, money.USD).GrowthRate(0).Periods(5).Present()
			},
			func() (money.Money, error) {
				return NewGeometricSeries().AnnualRate(0.12).Quarterly().
					FirstPayment(1000, money.USD).GrowthRate(0).Periods(5).Present()
			},
		},
	}

	for _, pair := range pairs {
		t.Run(pair.name, func(t *testing.T) {
			before, err := pair.before()
			require.NoError(t, err)

			after, err := pair.after()
			require.NoError(t, err)

			assert.InDelta(t, before.InexactFloat64(), after.InexactFloat64(), 1e-9,
				"the method order changed the result")
		})
	}

	// The annual rate must actually be divided: 12% a year over quarterly
	// periods is 3% each, so one period discounts by 1.03.
	present, err := NewGeometricSeries().AnnualRate(0.12).Quarterly().
		FirstPayment(1000, money.USD).GrowthRate(0).Periods(1).Present()
	require.NoError(t, err)
	assert.InDelta(t, 1000/1.03, present.InexactFloat64(), 1e-6)
}

// TestRateOverridesAnnualRate checks the two rate setters replace each other
// rather than compounding.
func TestRateOverridesAnnualRate(t *testing.T) {
	// A periodic rate set last is used as given, not divided again.
	present, err := NewArithmeticSeries().Annually().AnnualRate(0.12).Rate(0.05).
		FirstPayment(1000, money.USD).Gradient(0, money.USD).Periods(1).Present()
	require.NoError(t, err)
	assert.InDelta(t, 1000/1.05, present.InexactFloat64(), 1e-6)
}

// TestGeometricBuilderAnnualRate covers the same conversion on the geometric
// builder.
func TestGeometricBuilderAnnualRate(t *testing.T) {
	present, err := NewGeometricSeries().
		Quarterly().
		AnnualRate(0.12).
		FirstPayment(1000, money.USD).
		GrowthRate(0).
		Periods(1).
		Present()
	require.NoError(t, err)

	assert.InDelta(t, 1000/1.03, present.InexactFloat64(), 1e-6)
}

// TestBuilderMustVariants covers the panicking builder helpers: they return
// the value on success and panic on a failing configuration.
func TestBuilderMustVariants(t *testing.T) {
	arithmetic := NewArithmeticSeries().
		FirstPayment(1000, money.USD).
		Gradient(100, money.USD).
		Rate(0.10).
		Periods(5).
		Annually()

	assert.NotPanics(t, func() {
		series := arithmetic.MustBuild()

		present, err := series.Present()
		require.NoError(t, err)
		assert.InDelta(t, 4476.9669, present.InexactFloat64(), 0.01)
	})

	assert.InDelta(t, 4476.9669, arithmetic.MustPresent().InexactFloat64(), 0.01)
	assert.InDelta(t, 7210.20, arithmetic.MustFuture().InexactFloat64(), 0.01)

	geometric := NewGeometricSeries().
		FirstPayment(1000, money.USD).
		GrowthRate(0.08).
		Rate(0.10).
		Periods(5).
		Annually()

	assert.NotPanics(t, func() { geometric.MustBuild() })
	assert.Positive(t, geometric.MustPresent().InexactFloat64())
	assert.Positive(t, geometric.MustFuture().InexactFloat64())

	// A currency mismatch fails Build, so the Must variants must panic.
	mismatched := NewArithmeticSeries().
		FirstPayment(1000, money.USD).
		Gradient(100, money.EUR).
		Rate(0.10).
		Periods(5).
		Annually()

	assert.Panics(t, func() { mismatched.MustBuild() })
	assert.Panics(t, func() { mismatched.MustPresent() })
	assert.Panics(t, func() { mismatched.MustFuture() })
}

// TestBuilderPropagatesCurrencyMismatch checks the non-panicking builder path
// surfaces the same error.
func TestBuilderPropagatesCurrencyMismatch(t *testing.T) {
	_, err := NewArithmeticSeries().
		FirstPayment(1000, money.USD).
		Gradient(100, money.EUR).
		Rate(0.10).
		Periods(5).
		Annually().
		Present()
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	_, err = NewArithmeticSeries().
		FirstPayment(1000, money.USD).
		Gradient(100, money.EUR).
		Rate(0.10).
		Periods(5).
		Annually().
		Future()
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)
}

// TestBuilderRejectsNegativeRate covers the error path shared by every builder
// entry point: a negative rate fails when the rate is constructed, before any
// series exists.
func TestBuilderRejectsNegativeRate(t *testing.T) {
	arithmetic := NewArithmeticSeries().
		FirstPayment(1000, money.USD).
		Gradient(100, money.USD).
		Rate(-0.10).
		Periods(5).
		Annually()

	_, err := arithmetic.Build()
	require.Error(t, err)

	_, err = arithmetic.Present()
	require.Error(t, err)

	_, err = arithmetic.Future()
	require.Error(t, err)

	assert.Panics(t, func() { arithmetic.MustBuild() })
	assert.Panics(t, func() { arithmetic.MustPresent() })
	assert.Panics(t, func() { arithmetic.MustFuture() })

	geometric := NewGeometricSeries().
		FirstPayment(1000, money.USD).
		GrowthRate(0.08).
		Rate(-0.10).
		Periods(5).
		Monthly()

	_, err = geometric.Build()
	require.Error(t, err)

	_, err = geometric.Present()
	require.Error(t, err)

	_, err = geometric.Future()
	require.Error(t, err)

	assert.Panics(t, func() { geometric.MustBuild() })
	assert.Panics(t, func() { geometric.MustPresent() })
	assert.Panics(t, func() { geometric.MustFuture() })
}

// TestOverflowingTermReportsError exercises the error branches inside Present
// and Future: a term long enough to overflow (1+i)^n must surface the overflow
// instead of panicking or returning a bogus amount.
func TestOverflowingTermReportsError(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(100000), compoundinterest.Annually)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(decimal.One, compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	arithmetic, err := NewArithmetic(
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(100, money.USD),
		period, rate)
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		_, err := arithmetic.Present()
		assert.Error(t, err)

		_, err = arithmetic.Future()
		assert.Error(t, err)
	})

	// Growing far faster than the discount rate, so the ratio (1+g)/(1+i)
	// exceeds one and its power overflows too.
	geometric, err := NewGeometric(
		money.MustMoneyFromFloat64(1000, money.USD),
		decimal.MustFromFloat64(2),
		period, rate)
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		_, err := geometric.Present()
		assert.Error(t, err)

		_, err = geometric.Future()
		assert.Error(t, err)
	})
}

// TestGeometricCurrencyMismatchFromGrowth checks that a geometric series needs
// no currency check on its growth rate — it is a plain decimal — so a series
// in any currency builds cleanly.
func TestGeometricCurrencyMismatchFromGrowth(t *testing.T) {
	_, err := NewGeometricSeries().
		FirstPayment(1000, money.JPY).
		GrowthRate(0.08).
		Rate(0.10).
		Periods(5).
		Annually().
		Build()
	require.NoError(t, err)
}
