package annuities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// newZeroRateAnnuity builds a 12-period monthly annuity at a 0% rate with the
// given payment, present and future values.
func newZeroRateAnnuity(t *testing.T, value, present, future float64) Annuity {
	t.Helper()

	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.Zero, compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	annuity, err := New(
		money.MustMoneyFromFloat64(value, money.USD),
		money.MustMoneyFromFloat64(present, money.USD),
		money.MustMoneyFromFloat64(future, money.USD),
		period,
		rateInterest,
	)
	require.NoError(t, err)

	return annuity
}

// TestZeroRateNeverPanics is the regression test for TESTING_PLAN.md §2.1: the
// payment functions called the panicking decimal helpers (MustPow, MustDiv)
// inside functions that return an error, so a 0% rate — a legitimate input,
// as in promotional financing or a loan between family members — crashed the
// caller instead of returning a value.
//
// A function whose signature promises an error must never panic on a valid
// input, so this sweeps the whole payment API and fails on any panic.
func TestZeroRateNeverPanics(t *testing.T) {
	annuity := newZeroRateAnnuity(t, 100, 12000, 12000)

	calls := map[string]func() (money.Money, error){
		"PaymentFromPresentValue":           annuity.PaymentFromPresentValue,
		"PaymentFromFutureValue":            annuity.PaymentFromFutureValue,
		"AnticipatePaymentFromPresentValue": annuity.AnticipatePaymentFromPresentValue,
		"AnticipatePaymentFromFutureValue":  annuity.AnticipatePaymentFromFutureValue,
		"Present":                           annuity.Present,
		"AnticipatePresent":                 annuity.AnticipatePresent,
		"Future":                            annuity.Future,
		"AnticipateFuture":                  annuity.AnticipateFuture,
		"FutureWithContributions":           annuity.FutureWithContributions,
		"AnticipateFutureWithContributions": annuity.AnticipateFutureWithContributions,
		"PaymentFromPresentValueDeferred":   func() (money.Money, error) { return annuity.PaymentFromPresentValueDeferred(3) },
		"PresentDeferred":                   func() (money.Money, error) { return annuity.PresentDeferred(3) },
		"AnticipatePresentDeferred":         func() (money.Money, error) { return annuity.AnticipatePresentDeferred(3) },
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, err := call()
				assert.NoError(t, err)
			})
		})
	}
}

// TestZeroRatePaymentValues checks the analytic limits the payment functions
// must return at a 0% rate: with no interest to service, the amount is simply
// split evenly across the periods.
func TestZeroRatePaymentValues(t *testing.T) {
	// PV = FV = 12000 over 12 periods -> 1000 per period, whether the
	// payment falls at the beginning or the end of the period.
	annuity := newZeroRateAnnuity(t, 0, 12000, 12000)

	payments := map[string]func() (money.Money, error){
		"PaymentFromPresentValue":           annuity.PaymentFromPresentValue,
		"PaymentFromFutureValue":            annuity.PaymentFromFutureValue,
		"AnticipatePaymentFromPresentValue": annuity.AnticipatePaymentFromPresentValue,
		"AnticipatePaymentFromFutureValue":  annuity.AnticipatePaymentFromFutureValue,
	}

	for name, payment := range payments {
		t.Run(name, func(t *testing.T) {
			got, err := payment()
			require.NoError(t, err)
			assert.InDelta(t, 1000.0, got.InexactFloat64(), 1e-9)
		})
	}
}

// TestZeroPeriodsReportsError checks the other degenerate input: an annuity
// with no periods has no payment to compute, since both the zero-rate limit
// (1/n) and the general formula divide by zero. It must be reported as an
// error, and — like the zero rate — must not panic.
func TestZeroPeriodsReportsError(t *testing.T) {
	rates := map[string]float64{"zero rate": 0, "nonzero rate": 0.01}

	for name, rate := range rates {
		t.Run(name, func(t *testing.T) {
			period, err := compoundinterest.NewPeriod(decimal.Zero, compoundinterest.Monthly)
			require.NoError(t, err)

			rateInterest, err := compoundinterest.NewRateInterest(
				decimal.MustFromFloat64(rate), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
			require.NoError(t, err)

			zero := money.MustMoneyFromFloat64(0, money.USD)
			present := money.MustMoneyFromFloat64(12000, money.USD)

			annuity, err := New(zero, present, present, period, rateInterest)
			require.NoError(t, err)

			payments := map[string]func() (money.Money, error){
				"PaymentFromPresentValue":           annuity.PaymentFromPresentValue,
				"PaymentFromFutureValue":            annuity.PaymentFromFutureValue,
				"AnticipatePaymentFromPresentValue": annuity.AnticipatePaymentFromPresentValue,
				"AnticipatePaymentFromFutureValue":  annuity.AnticipatePaymentFromFutureValue,
				"PaymentFromPresentValueDeferred":   func() (money.Money, error) { return annuity.PaymentFromPresentValueDeferred(3) },
			}

			for paymentName, payment := range payments {
				t.Run(paymentName, func(t *testing.T) {
					assert.NotPanics(t, func() {
						_, err := payment()
						assert.Error(t, err)
					})
				})
			}
		})
	}
}

// TestOverflowingTermReportsError exercises the error branch that replaced the
// panicking MustPow: a term long enough to overflow (1+i)^n must surface the
// overflow rather than crash.
func TestOverflowingTermReportsError(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(100000), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(
		decimal.One, compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	amount := money.MustMoneyFromFloat64(1000, money.USD)

	annuity, err := New(money.MustMoneyFromFloat64(0, money.USD), amount, amount, period, rateInterest)
	require.NoError(t, err)

	payments := map[string]func() (money.Money, error){
		"PaymentFromPresentValue":           annuity.PaymentFromPresentValue,
		"PaymentFromFutureValue":            annuity.PaymentFromFutureValue,
		"AnticipatePaymentFromPresentValue": annuity.AnticipatePaymentFromPresentValue,
		"AnticipatePaymentFromFutureValue":  annuity.AnticipatePaymentFromFutureValue,
	}

	for name, payment := range payments {
		t.Run(name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, err := payment()
				assert.Error(t, err)
			})
		})
	}
}

// TestZeroRateScheduleCloses checks the invariant end to end: the payment
// solved at a 0% rate must amortize the loan exactly, leaving no balance and
// charging no interest.
func TestZeroRateScheduleCloses(t *testing.T) {
	annuity := newZeroRateAnnuity(t, 0, 12000, 0)

	payment, err := annuity.PaymentFromPresentValue()
	require.NoError(t, err)

	schedule, err := BuildSchedule(
		money.MustMoneyFromFloat64(12000, money.USD),
		decimal.Zero,
		payment,
		decimal.MustFromInt64(12, 0),
	)
	require.NoError(t, err)

	last := schedule[len(schedule)-1]
	assert.InDelta(t, 0.0, last.Balance.InexactFloat64(), 1e-9)
	assert.InDelta(t, 0.0, last.SumInterest.InexactFloat64(), 1e-9)

	principal := money.MustMoneyFromFloat64(0, money.USD)
	for _, row := range schedule {
		principal = principal.Add(row.Principal)
	}

	assert.InDelta(t, 12000.0, principal.InexactFloat64(), 1e-9)
}

// TestAnticipatedRateAnnuityIsComputed guards the junction between the two
// defects in TESTING_PLAN.md: an annuity configured with an anticipated
// (discount) rate used to get a periodic rate of 0 from the silent-zero
// conversion (§2.2), which then panicked the payment functions (§2.1). Both
// halves are fixed, so the annuity must now compute the same payment as the
// equivalent ordinary rate.
func TestAnticipatedRateAnnuityIsComputed(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	// d = 1/101 is the anticipated rate equivalent to an ordinary 1%.
	anticipated, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.01).MustDiv(decimal.MustFromFloat64(1.01)),
		compoundinterest.Monthly,
		compoundinterest.RateAnticipateEffectyPeriodic,
	)
	require.NoError(t, err)

	ordinary, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.01),
		compoundinterest.Monthly,
		compoundinterest.RateEffectyPeriodic,
	)
	require.NoError(t, err)

	present := money.MustMoneyFromFloat64(10000, money.USD)
	zero := money.MustMoneyFromFloat64(0, money.USD)

	withAnticipated, err := New(zero, present, zero, period, anticipated)
	require.NoError(t, err)

	withOrdinary, err := New(zero, present, zero, period, ordinary)
	require.NoError(t, err)

	fromAnticipated, err := withAnticipated.PaymentFromPresentValue()
	require.NoError(t, err)

	fromOrdinary, err := withOrdinary.PaymentFromPresentValue()
	require.NoError(t, err)

	assert.InDelta(t, fromOrdinary.InexactFloat64(), fromAnticipated.InexactFloat64(), 1e-9)
	// $10,000 over 12 months at 1% -> 888.4879 per month.
	assert.InDelta(t, 888.4879, fromAnticipated.InexactFloat64(), 1e-4)
}
