package tvm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
)

// TestOverflowNeverPanics is the regression test for the panics the solver
// carried: PV·(1+i)ᴺ and PMT·coef were formed with the panicking Mul, so a
// large principal or payment crashed functions that return an error.
func TestOverflowNeverPanics(t *testing.T) {
	solvers := map[string]func() (decimal.Decimal, error){
		"SolveFV from PV":  NewTVM().PV(1e35).Rate(0.10).N(100).SolveFV,
		"SolveFV from PMT": NewTVM().PMT(1e35).Rate(0.10).N(100).SolveFV,
		"SolvePV":          NewTVM().FV(1e35).Rate(0.10).N(100).SolvePV,
		"SolvePMT":         NewTVM().PV(1e35).FV(1e35).Rate(0.10).N(100).SolvePMT,
		"SolveRate":        NewTVM().PV(1e35).PMT(-1e34).N(100).SolveRate,
	}

	for name, solve := range solvers {
		t.Run(name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, err := solve()
				assert.Error(t, err)
			})
		})
	}
}

// TestSolveRateLongTerm checks the rate solver on realistic long schedules,
// where the candidate sweep's high end overflows the growth factor and those
// candidates have to be skipped rather than abandoning the search.
func TestSolveRateLongTerm(t *testing.T) {
	for _, periods := range []float64{120, 360, 600} {
		t.Run(decimal.MustFromFloat64(periods).String(), func(t *testing.T) {
			payment := NewTVM().PV(300000).Rate(0.005).N(periods).MustSolvePMT()

			rate, err := NewTVM().PV(300000).PMT(payment.InexactFloat64()).N(periods).SolveRate()
			require.NoError(t, err)
			assert.InDelta(t, 0.005, rate.InexactFloat64(), 1e-6)
		})
	}
}

// TestFiveUnknownsRoundTrip is the solver's central invariant: fix any four of
// the five variables, solve the fifth, and feeding it back must reproduce each
// of the others.
func TestFiveUnknownsRoundTrip(t *testing.T) {
	const (
		pv      = 300000.0
		rate    = 0.005
		periods = 360.0
	)

	payment := NewTVM().PV(pv).Rate(rate).N(periods).MustSolvePMT()

	recoveredPV, err := NewTVM().PMT(payment.InexactFloat64()).Rate(rate).N(periods).SolvePV()
	require.NoError(t, err)
	assert.InDelta(t, pv, recoveredPV.InexactFloat64(), 0.01)

	recoveredRate, err := NewTVM().PV(pv).PMT(payment.InexactFloat64()).N(periods).SolveRate()
	require.NoError(t, err)
	assert.InDelta(t, rate, recoveredRate.InexactFloat64(), 1e-6)

	recoveredN, err := NewTVM().PV(pv).PMT(payment.InexactFloat64()).Rate(rate).SolveN()
	require.NoError(t, err)
	assert.InDelta(t, periods, recoveredN.InexactFloat64(), 1e-6)

	// The future value of a fully amortized loan is zero.
	futureValue, err := NewTVM().PV(pv).PMT(payment.InexactFloat64()).Rate(rate).N(periods).SolveFV()
	require.NoError(t, err)
	assert.InDelta(t, 0.0, futureValue.InexactFloat64(), 0.01)
}

// TestDueVersusOrdinary checks the annuity-due timing against the ordinary
// one across all five solvers, and that Ordinary restores the default.
func TestDueVersusOrdinary(t *testing.T) {
	const (
		pv      = 100000.0
		rate    = 0.01
		periods = 60.0
	)

	ordinary, err := NewTVM().PV(pv).Rate(rate).N(periods).SolvePMT()
	require.NoError(t, err)

	due, err := NewTVM().PV(pv).Rate(rate).N(periods).Due().SolvePMT()
	require.NoError(t, err)

	// Paying at the start of each period services less interest, so the
	// payment is the ordinary one divided by (1+i). Both are negative under
	// the sign convention, so the due payment is smaller in magnitude.
	assert.InDelta(t, ordinary.InexactFloat64()/(1+rate), due.InexactFloat64(), 1e-6)
	assert.Less(t, -due.InexactFloat64(), -ordinary.InexactFloat64())

	// Ordinary() must undo Due().
	restored, err := NewTVM().PV(pv).Rate(rate).N(periods).Due().Ordinary().SolvePMT()
	require.NoError(t, err)
	assert.InDelta(t, ordinary.InexactFloat64(), restored.InexactFloat64(), 1e-9)

	// The same relation holds on the future value of a savings plan.
	ordinaryFV, err := NewTVM().PMT(-100).Rate(rate).N(periods).SolveFV()
	require.NoError(t, err)

	dueFV, err := NewTVM().PMT(-100).Rate(rate).N(periods).Due().SolveFV()
	require.NoError(t, err)

	assert.InDelta(t, ordinaryFV.InexactFloat64()*(1+rate), dueFV.InexactFloat64(), 1e-6)
}

// TestDueRoundTrips checks each solver inverts correctly under annuity-due
// timing too, not just the ordinary default.
func TestDueRoundTrips(t *testing.T) {
	const (
		pv      = 50000.0
		rate    = 0.008
		periods = 120.0
	)

	payment := NewTVM().PV(pv).Rate(rate).N(periods).Due().MustSolvePMT()

	recoveredPV, err := NewTVM().PMT(payment.InexactFloat64()).Rate(rate).N(periods).Due().SolvePV()
	require.NoError(t, err)
	assert.InDelta(t, pv, recoveredPV.InexactFloat64(), 0.01)

	recoveredN, err := NewTVM().PV(pv).PMT(payment.InexactFloat64()).Rate(rate).Due().SolveN()
	require.NoError(t, err)
	assert.InDelta(t, periods, recoveredN.InexactFloat64(), 1e-6)

	recoveredRate, err := NewTVM().PV(pv).PMT(payment.InexactFloat64()).N(periods).Due().SolveRate()
	require.NoError(t, err)
	assert.InDelta(t, rate, recoveredRate.InexactFloat64(), 1e-6)
}

// TestZeroRateAcrossSolvers checks the degenerate but legitimate 0% case in
// every solver: with no interest the money neither grows nor discounts.
func TestZeroRateAcrossSolvers(t *testing.T) {
	// Borrowing 1200 and repaying 100 for 12 periods clears it exactly.
	futureValue, err := NewTVM().PV(-1200).PMT(100).Rate(0).N(12).SolveFV()
	require.NoError(t, err)
	assert.InDelta(t, 0.0, futureValue.InexactFloat64(), 1e-9)

	presentValue, err := NewTVM().PMT(100).FV(0).Rate(0).N(12).SolvePV()
	require.NoError(t, err)
	assert.InDelta(t, -1200.0, presentValue.InexactFloat64(), 1e-9)

	payment, err := NewTVM().PV(1200).FV(0).Rate(0).N(12).SolvePMT()
	require.NoError(t, err)
	assert.InDelta(t, -100.0, payment.InexactFloat64(), 1e-9)

	periods, err := NewTVM().PV(-1200).PMT(100).FV(0).Rate(0).SolveN()
	require.NoError(t, err)
	assert.InDelta(t, 12.0, periods.InexactFloat64(), 1e-9)

	// At a zero rate the due and ordinary payments coincide.
	duePayment, err := NewTVM().PV(1200).FV(0).Rate(0).N(12).Due().SolvePMT()
	require.NoError(t, err)
	assert.InDelta(t, -100.0, duePayment.InexactFloat64(), 1e-9)
}

// TestSolveNUnrepayableBalance covers a payment smaller than the periodic
// interest: the balance grows forever, so no term repays it.
func TestSolveNUnrepayableBalance(t *testing.T) {
	_, err := NewTVM().PV(1000).PMT(-1).Rate(0.10).SolveN()
	assert.ErrorIs(t, err, ErrNoSolution)
}

// TestSolveNReturnsNegativeTermForPastBalance documents current behaviour
// rather than endorsing it: when the payment points the same way as the
// balance, nothing repays anything going forward, but the equation still has a
// root at a negative term — the point in the past where the balance would have
// been zero. SolveN returns it without comment.
//
// Whether a solver should hand back a negative number of periods, or report
// that no future term exists, is an open question recorded in TESTING_PLAN.md.
// This test pins the behaviour so a deliberate change is visible.
func TestSolveNReturnsNegativeTermForPastBalance(t *testing.T) {
	periods, err := NewTVM().PV(1000).PMT(100).Rate(0.10).SolveN()
	require.NoError(t, err)
	assert.Negative(t, periods.InexactFloat64())
}

// TestInvalidRateBound checks the documented domain: 1+rate must be positive.
func TestInvalidRateBound(t *testing.T) {
	solvers := map[string]func() (decimal.Decimal, error){
		"SolveFV":  NewTVM().PV(1000).Rate(-1).N(10).SolveFV,
		"SolvePV":  NewTVM().FV(1000).Rate(-1).N(10).SolvePV,
		"SolvePMT": NewTVM().PV(1000).Rate(-1).N(10).SolvePMT,
		"SolveN":   NewTVM().PV(1000).PMT(-100).Rate(-1.5).SolveN,
	}

	for name, solve := range solvers {
		t.Run(name, func(t *testing.T) {
			_, err := solve()
			assert.ErrorIs(t, err, ErrInvalidRate)
		})
	}
}

// TestSolveRateNoConvergence checks the exhaustion path: cash flows that never
// balance at any candidate rate must be reported, not approximated.
func TestSolveRateNoConvergence(t *testing.T) {
	// All flows point the same way, so the residual never crosses zero.
	_, err := NewTVM().PV(1000).PMT(100).FV(1000).N(10).SolveRate()
	assert.ErrorIs(t, err, ErrNoConvergence)
}

// TestNegativeRate covers a genuinely negative periodic rate, which the solver
// accepts as long as 1+rate stays positive.
func TestNegativeRate(t *testing.T) {
	// Money shrinking 1% per period: 1000 becomes 1000×0.99^10.
	futureValue, err := NewTVM().PV(-1000).Rate(-0.01).N(10).SolveFV()
	require.NoError(t, err)
	assert.InDelta(t, 904.38, futureValue.InexactFloat64(), 0.01)
}

// TestMustVariants covers the panicking helpers, several of which had no
// coverage: they return the value on success and panic on a failing setup.
func TestMustVariants(t *testing.T) {
	valid := NewTVM().PV(300000).Rate(0.005).N(360)

	assert.NotPanics(t, func() {
		assert.NotZero(t, valid.MustSolvePMT().InexactFloat64())
		assert.NotZero(t, valid.MustSolveFV().InexactFloat64())
	})

	assert.NotPanics(t, func() {
		payment := valid.MustSolvePMT()
		assert.NotZero(t, NewTVM().PMT(payment.InexactFloat64()).Rate(0.005).N(360).MustSolvePV().InexactFloat64())
		assert.NotZero(t, NewTVM().PV(300000).PMT(payment.InexactFloat64()).N(360).MustSolveRate().InexactFloat64())
		assert.NotZero(t, NewTVM().PV(300000).PMT(payment.InexactFloat64()).Rate(0.005).MustSolveN().InexactFloat64())
	})

	invalid := NewTVM().PV(1000).Rate(-1).N(10)

	assert.Panics(t, func() { invalid.MustSolveFV() })
	assert.Panics(t, func() { invalid.MustSolvePV() })
	assert.Panics(t, func() { invalid.MustSolvePMT() })
	assert.Panics(t, func() { NewTVM().PV(1000).PMT(0).Rate(0).MustSolveN() })
	assert.Panics(t, func() { NewTVM().PV(1000).PMT(100).FV(1000).N(10).MustSolveRate() })
}

// TestSingleAndZeroPeriods covers the shortest terms: one period behaves like
// simple interest, and zero periods leaves the present value untouched.
func TestSingleAndZeroPeriods(t *testing.T) {
	futureValue, err := NewTVM().PV(-1000).Rate(0.05).N(1).SolveFV()
	require.NoError(t, err)
	assert.InDelta(t, 1050.0, futureValue.InexactFloat64(), 1e-9)

	unchanged, err := NewTVM().PV(-1000).Rate(0.05).N(0).SolveFV()
	require.NoError(t, err)
	assert.InDelta(t, 1000.0, unchanged.InexactFloat64(), 1e-9)
}
