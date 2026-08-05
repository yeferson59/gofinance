package annuities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// TestBuildScheduleRejectsFractionalPeriods covers the resolution of
// TESTING_PLAN.md §2.7. A schedule has one row per payment, so a fractional
// term has no meaning; it used to be truncated silently, answering a different
// question than the one asked.
func TestBuildScheduleRejectsFractionalPeriods(t *testing.T) {
	principal := money.MustMoneyFromFloat64(1000, money.USD)
	payment := money.MustMoneyFromFloat64(100, money.USD)
	rate := decimal.MustFromFloat64(0.01)

	for _, nper := range []float64{5.7, 5.5, 0.5, 12.000001, -3.5} {
		_, err := BuildSchedule(principal, rate, payment, decimal.MustFromFloat64(nper))
		assert.ErrorIs(t, err, ErrInvalidPeriods, "nper %v must be rejected", nper)
	}
}

// TestBuildScheduleAcceptsWholePeriods checks the rejection above does not
// catch a legitimate whole term, however it was constructed.
func TestBuildScheduleAcceptsWholePeriods(t *testing.T) {
	principal := money.MustMoneyFromFloat64(1000, money.USD)
	payment := money.MustMoneyFromFloat64(100, money.USD)
	rate := decimal.MustFromFloat64(0.01)

	wholeTerms := []decimal.Decimal{
		decimal.MustFromInt64(12, 0),
		decimal.MustFromFloat64(12),
		// A whole number carrying trailing zeros is still whole.
		decimal.MustFromString("12.000"),
	}

	for _, nper := range wholeTerms {
		schedule, err := BuildSchedule(principal, rate, payment, nper)
		require.NoError(t, err, "nper %v must be accepted", nper)
		assert.Len(t, schedule, 13)
	}
}

// TestBuilderPresentValue covers the builder methods added to close the
// symmetry gap in TESTING_PLAN.md §2.9: the builder offered FutureValue and
// DeferredPresentValue but no plain PresentValue, so callers had to route
// through Defer(0).
func TestBuilderPresentValue(t *testing.T) {
	config := NewAnnuity().
		Value(500, money.USD).
		AnnualRate(0.06).
		Periods(12).
		Monthly()

	present, err := config.PresentValue()
	require.NoError(t, err)

	// The value the old workaround produced.
	viaDefer, err := config.Defer(0).DeferredPresentValue()
	require.NoError(t, err)

	assert.InDelta(t, viaDefer.InexactFloat64(), present.InexactFloat64(), 1e-9)

	// And it agrees with the annuity's own method.
	annuity, err := config.Build()
	require.NoError(t, err)

	direct, err := annuity.Present()
	require.NoError(t, err)

	assert.InDelta(t, direct.InexactFloat64(), present.InexactFloat64(), 1e-9)

	assert.NotPanics(t, func() {
		assert.InDelta(t, present.InexactFloat64(), config.MustPresentValue().InexactFloat64(), 1e-9)
	})
}

// TestBuilderAnticipatePresentValue checks the annuity-due counterpart: paying
// at the start of each period is worth more today than paying at the end.
func TestBuilderAnticipatePresentValue(t *testing.T) {
	config := NewAnnuity().
		Value(500, money.USD).
		AnnualRate(0.06).
		Periods(12).
		Monthly()

	ordinary, err := config.PresentValue()
	require.NoError(t, err)

	due, err := config.AnticipatePresentValue()
	require.NoError(t, err)

	assert.Greater(t, due.InexactFloat64(), ordinary.InexactFloat64())

	assert.NotPanics(t, func() {
		assert.InDelta(t, due.InexactFloat64(), config.MustAnticipatePresentValue().InexactFloat64(), 1e-9)
	})
}

// TestBuilderPresentValuePropagatesErrors checks the new methods report a bad
// configuration rather than producing a figure, and that the Must variants
// panic on it.
func TestBuilderPresentValuePropagatesErrors(t *testing.T) {
	// A negative rate fails when the rate is constructed.
	broken := NewAnnuity().
		Value(500, money.USD).
		Rate(-0.5).
		Periods(12).
		Monthly()

	_, err := broken.PresentValue()
	assert.Error(t, err)

	_, err = broken.AnticipatePresentValue()
	assert.Error(t, err)

	assert.Panics(t, func() { broken.MustPresentValue() })
	assert.Panics(t, func() { broken.MustAnticipatePresentValue() })
}
