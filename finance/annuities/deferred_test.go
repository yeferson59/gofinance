package annuities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestAnnuityPresentDeferred(t *testing.T) {
	// PMT = 1000, 1% monthly, n = 12, deferred 3 periods.
	// PV_ordinary = 1000 × [1 - 1.01^-12] / 0.01 = 11255.0775
	// PV_deferred = PV_ordinary / 1.01^3 = 10924.07
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.01), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value := money.MustMoneyFromFloat64(1000, money.USD)
	zero := money.MustMoneyFromFloat64(0, money.USD)

	annuity, err := New(value, zero, zero, period, rateInterest)
	require.NoError(t, err)

	present, err := annuity.PresentDeferred(3)
	require.NoError(t, err)
	assert.InDelta(t, 10924.07, present.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityPresentDeferredZeroReducesToPresent(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.01), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value := money.MustMoneyFromFloat64(1000, money.USD)
	zero := money.MustMoneyFromFloat64(0, money.USD)

	annuity, err := New(value, zero, zero, period, rateInterest)
	require.NoError(t, err)

	deferred, err := annuity.PresentDeferred(0)
	require.NoError(t, err)

	ordinary, err := annuity.Present()
	require.NoError(t, err)

	assert.InDelta(t, ordinary.ToDecimal().InexactFloat64(), deferred.ToDecimal().InexactFloat64(), 0.0001)
}

func TestAnnuityAnticipatePresentDeferred(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.01), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	value := money.MustMoneyFromFloat64(1000, money.USD)
	zero := money.MustMoneyFromFloat64(0, money.USD)

	annuity, err := New(value, zero, zero, period, rateInterest)
	require.NoError(t, err)

	deferredDue, err := annuity.AnticipatePresentDeferred(3)
	require.NoError(t, err)

	ordinaryDue, err := annuity.AnticipatePresent()
	require.NoError(t, err)

	growth := 1.01 * 1.01 * 1.01
	assert.InDelta(t, ordinaryDue.ToDecimal().InexactFloat64()/growth, deferredDue.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityPaymentFromPresentValueDeferred(t *testing.T) {
	// Present = 10924.07, 1% monthly, n = 12, deferred 3 periods should
	// recover PMT ≈ 1000 (inverse of TestAnnuityPresentDeferred).
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.01), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	present := money.MustMoneyFromFloat64(10924.07, money.USD)
	zero := money.MustMoneyFromFloat64(0, money.USD)

	annuity, err := New(zero, present, zero, period, rateInterest)
	require.NoError(t, err)

	payment, err := annuity.PaymentFromPresentValueDeferred(3)
	require.NoError(t, err)
	assert.InDelta(t, 1000, payment.ToDecimal().InexactFloat64(), 0.5)
}

func TestAnnuityPaymentFromPresentValueDeferredZeroReducesToOrdinary(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.01), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	present := money.MustMoneyFromFloat64(11255.0775, money.USD)
	zero := money.MustMoneyFromFloat64(0, money.USD)

	annuity, err := New(zero, present, zero, period, rateInterest)
	require.NoError(t, err)

	deferred, err := annuity.PaymentFromPresentValueDeferred(0)
	require.NoError(t, err)

	ordinary, err := annuity.PaymentFromPresentValue()
	require.NoError(t, err)

	assert.InDelta(t, ordinary.ToDecimal().InexactFloat64(), deferred.ToDecimal().InexactFloat64(), 0.0001)
}

func TestAnnuityConfigDeferredPresentValue(t *testing.T) {
	present, err := NewAnnuity().
		Value(1000, money.USD).
		Rate(0.01).
		Periods(12).
		Monthly().
		Defer(3).
		DeferredPresentValue()
	require.NoError(t, err)
	assert.InDelta(t, 10924.07, present.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityConfigDeferredPayment(t *testing.T) {
	payment := NewAnnuity().
		Present(10924.07, money.USD).
		Rate(0.01).
		Periods(12).
		Monthly().
		Defer(3).
		MustDeferredPayment()
	assert.InDelta(t, 1000, payment.ToDecimal().InexactFloat64(), 0.5)
}
