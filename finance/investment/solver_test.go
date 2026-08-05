package investment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// longFlows builds an outflow followed by count level inflows.
func longFlows(outflow, inflow float64, count int) []money.Money {
	out := make([]money.Money, 0, count+1)
	out = append(out, usd(outflow))

	for range count {
		out = append(out, usd(inflow))
	}

	return out
}

// TestIRRLongSeriesUsesBisection covers the solver's fallback path, which had
// never been executed by a test (TESTING_PLAN.md §1). A 400-period series
// pushes Newton–Raphson out of its domain, so the bracketed bisection search
// takes over.
//
// Exercising it surfaced two defects, both fixed: the running sum and the
// discount factor were accumulated with the panicking Add/Mul helpers, so a
// long series crashed a function that returns an error; and the scan aborted
// when its *first* candidate rate failed to evaluate, even though every later
// candidate was skipped on failure.
func TestIRRLongSeriesUsesBisection(t *testing.T) {
	// −100,000 followed by 400 payments of 1,000: 400,000 repaid on 100,000.
	flows := longFlows(-100000, 1000, 400)

	irr, err := IRR(flows)
	require.NoError(t, err)

	// The rate must zero out the NPV, whichever path found it.
	npv, err := NPV(irr, flows)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, npv.InexactFloat64(), 0.01)

	// Sanity bounds: repaying 4× over 400 periods is a small positive rate.
	assert.Positive(t, irr.InexactFloat64())
	assert.Less(t, irr.InexactFloat64(), 0.05)
}

// TestIRRLongSeriesNeverPanics is the regression test for the panic the
// bisection path exposed: a series long enough to overflow the discount factor
// must report the overflow, not crash.
func TestIRRLongSeriesNeverPanics(t *testing.T) {
	assert.NotPanics(t, func() {
		_, _ = IRR(longFlows(-100000, 1000, 2000))
	})

	assert.NotPanics(t, func() {
		_, _ = NPV(decimal.MustFromFloat64(0.5), longFlows(-100000, 1000, 2000))
	})
}

// TestNPVOverflowReportsError checks that an overflowing discount factor comes
// back as an error rather than a panic or a silent wrong figure.
func TestNPVOverflowReportsError(t *testing.T) {
	_, err := NPV(decimal.MustFromFloat64(0.5), longFlows(-100000, 1000, 2000))
	require.Error(t, err)
}

// TestIRRMultipleSignChanges documents the classic limitation: with more than
// one sign change several rates can zero the NPV, and IRR returns one of them
// without signalling the ambiguity (TESTING_PLAN.md §2.6). The contract is that
// whatever it returns is genuinely a root.
func TestIRRMultipleSignChanges(t *testing.T) {
	// −1000, +6000, −11000, +6000 has roots at 0%, 100% and 200%.
	flows := flows(-1000, 6000, -11000, 6000)

	irr, err := IRR(flows)
	require.NoError(t, err)

	npv, err := NPV(irr, flows)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, npv.InexactFloat64(), 0.01)
}

// TestIRRNoRealRoot checks the case where the flows change sign but no rate
// zeroes the NPV: the scan must exhaust its candidates and report
// ErrNoConvergence rather than returning a bogus rate.
func TestIRRNoRealRoot(t *testing.T) {
	_, err := IRR(flows(-1, 3, -3))
	assert.ErrorIs(t, err, ErrNoConvergence)
}

// TestIRRExtremeRates covers the ends of the candidate range: a near-total loss
// and a very large multiple.
func TestIRRExtremeRates(t *testing.T) {
	tests := []struct {
		name     string
		flows    []money.Money
		expected float64
		delta    float64
	}{
		// Recovering 10 on 1000 in one period is a 99% loss.
		{"near total loss", flows(-1000, 10), -0.99, 1e-6},
		// Recovering 1000 on 1 in one period is 99900%.
		{"thousandfold return", flows(-1, 1000), 999, 1e-6},
		// A doubling every period.
		{"doubling", flows(-100, 200), 1.0, 1e-9},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			irr, err := IRR(test.flows)
			require.NoError(t, err)
			assert.InDelta(t, test.expected, irr.InexactFloat64(), test.delta)
		})
	}
}

// TestIRRSignChangeDirections checks both orientations: an investment (outflow
// first) and a financing (inflow first) must both resolve.
func TestIRRSignChangeDirections(t *testing.T) {
	investment, err := IRR(flows(-1000, 1100))
	require.NoError(t, err)
	assert.InDelta(t, 0.10, investment.InexactFloat64(), 1e-9)

	// Borrowing 1000 and repaying 1100 carries the same 10% rate.
	financing, err := IRR(flows(1000, -1100))
	require.NoError(t, err)
	assert.InDelta(t, 0.10, financing.InexactFloat64(), 1e-9)
}

// TestIRRZeroFlowsBetween checks that interior zero flows neither break the
// sign-change detection nor shift the root.
func TestIRRZeroFlowsBetween(t *testing.T) {
	// −100 now, nothing for two periods, 121 at t = 3 is not 10%: the money
	// compounds over three periods, so the rate is 121^(1/3)/100^(1/3) − 1.
	irr, err := IRR(flows(-100, 0, 0, 121))
	require.NoError(t, err)

	npv, err := NPV(irr, flows(-100, 0, 0, 121))
	require.NoError(t, err)
	assert.InDelta(t, 0.0, npv.InexactFloat64(), 1e-6)
}

// TestNPVIRRRoundTrip is the headline invariant: discounting a series at its
// own IRR must produce a zero net present value, for any series with a root.
func TestNPVIRRRoundTrip(t *testing.T) {
	series := [][]money.Money{
		flows(-1000, 400, 400, 400),
		flows(-5000, 1000, 1500, 2000, 2500),
		flows(-100, 110),
		flows(-250, 100, 100, 100, 50),
		flows(-10000, 2000, 2000, 2000, 2000, 2000, 2000),
	}

	for i, cashFlows := range series {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			irr, err := IRR(cashFlows)
			require.NoError(t, err)

			npv, err := NPV(irr, cashFlows)
			require.NoError(t, err)
			assert.InDelta(t, 0.0, npv.InexactFloat64(), 0.01)
		})
	}
}

// TestNPVMonotonicInRate checks a structural property of a conventional series
// (one outflow then inflows): its present value falls as the discount rate
// rises.
func TestNPVMonotonicInRate(t *testing.T) {
	cashFlows := flows(-1000, 400, 400, 400)

	previous := 0.0

	for i, rate := range []float64{0, 0.05, 0.10, 0.20, 0.50} {
		npv, err := NPV(decimal.MustFromFloat64(rate), cashFlows)
		require.NoError(t, err)

		if i > 0 {
			assert.Less(t, npv.InexactFloat64(), previous,
				"NPV must fall as the discount rate rises")
		}

		previous = npv.InexactFloat64()
	}
}

// dated builds a DatedCashFlow series from year offsets measured in days.
func dated(base time.Time, offsets []int, amounts []float64) []DatedCashFlow {
	out := make([]DatedCashFlow, len(offsets))
	for i := range offsets {
		out[i] = DatedCashFlow{
			Date:   base.AddDate(0, 0, offsets[i]),
			Amount: usd(amounts[i]),
		}
	}

	return out
}

// TestXIRRMatchesIRROnAnnualDates checks the two solvers agree: flows spaced
// exactly 365 days apart are an annual series, so XIRR must find the same rate
// IRR does.
func TestXIRRMatchesIRROnAnnualDates(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	amounts := []float64{-1000, 400, 400, 400}

	xirr, err := XIRR(dated(base, []int{0, 365, 730, 1095}, amounts))
	require.NoError(t, err)

	irr, err := IRR(flows(amounts...))
	require.NoError(t, err)

	assert.InDelta(t, irr.InexactFloat64(), xirr.InexactFloat64(), 1e-6)
}

// TestXNPVIRRRoundTrip is the dated counterpart of the NPV/IRR invariant.
func TestXNPVIRRRoundTrip(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	flows := dated(base, []int{0, 90, 200, 400, 900}, []float64{-10000, 2000, 3000, 4000, 3000})

	xirr, err := XIRR(flows)
	require.NoError(t, err)

	xnpv, err := XNPV(xirr, flows)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, xnpv.InexactFloat64(), 0.01)
}

// TestXIRRLeapYear checks a series spanning 29 February: the Actual/365 basis
// counts the real days, so the extra one must show up in the year fraction.
func TestXIRRLeapYear(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	flows := dated(base, []int{0, 366}, []float64{-1000, 1100})

	xirr, err := XIRR(flows)
	require.NoError(t, err)

	// 366/365 of a year at the solved rate must return exactly 1100.
	xnpv, err := XNPV(xirr, flows)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, xnpv.InexactFloat64(), 0.01)

	// Slightly under 10%, since the money worked for a day more than a year.
	assert.Less(t, xirr.InexactFloat64(), 0.10)
	assert.Greater(t, xirr.InexactFloat64(), 0.099)
}

// TestXIRRSameDayFlows checks that several flows on one date are summed rather
// than rejected or double-counted.
func TestXIRRSameDayFlows(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	split, err := XIRR(dated(base, []int{0, 0, 365}, []float64{-600, -400, 1100}))
	require.NoError(t, err)

	combined, err := XIRR(dated(base, []int{0, 365}, []float64{-1000, 1100}))
	require.NoError(t, err)

	assert.InDelta(t, combined.InexactFloat64(), split.InexactFloat64(), 1e-9)
}

// TestXNPVSingleFlow checks the degenerate one-flow series: with nothing to
// discount against, the present value is the flow itself.
func TestXNPVSingleFlow(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	xnpv, err := XNPV(decimal.MustFromFloat64(0.10), dated(base, []int{0}, []float64{-1000}))
	require.NoError(t, err)
	assert.InDelta(t, -1000.0, xnpv.InexactFloat64(), 1e-9)
}

// TestXIRRUnorderedDates checks that a series listed out of chronological order
// is rejected, since the first flow defines the base date.
func TestXIRRUnorderedDates(t *testing.T) {
	base := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	out := []DatedCashFlow{
		{Date: base, Amount: usd(-1000)},
		{Date: base.AddDate(0, 0, -30), Amount: usd(1100)},
	}

	_, err := XIRR(out)
	assert.ErrorIs(t, err, ErrDatesBeforeBase)

	_, err = XNPV(decimal.MustFromFloat64(0.1), out)
	assert.ErrorIs(t, err, ErrDatesBeforeBase)
}

// TestXIRRNoSignChange checks the dated solver reports the same missing-root
// condition as IRR.
func TestXIRRNoSignChange(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := XIRR(dated(base, []int{0, 365}, []float64{-1000, -500}))
	assert.ErrorIs(t, err, ErrNoSignChange)
}

// TestXIRRUsesBisection covers the dated solver's fallback, the other path
// that had never run (TESTING_PLAN.md §1). Recovering 5 on 1000 after a year
// puts the root at −99.5%; Newton's first step from a 10% guess overshoots far
// below −1, leaving the domain, so the bracketed scan takes over.
func TestXIRRUsesBisection(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	flows := dated(base, []int{0, 365}, []float64{-1000, 5})

	xirr, err := XIRR(flows)
	require.NoError(t, err)

	// −1000 + 5/(1+r) = 0 -> 1+r = 0.005 -> r = −99.5%
	assert.InDelta(t, -0.995, xirr.InexactFloat64(), 1e-6)

	xnpv, err := XNPV(xirr, flows)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, xnpv.InexactFloat64(), 0.01)
}

// TestXIRRBeyondCandidateRange checks the scan's exhaustion path: doubling
// your money in a single day annualizes to a rate far above the largest
// candidate, so no bracket exists and the failure is reported.
func TestXIRRBeyondCandidateRange(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	_, err := XIRR(dated(base, []int{0, 1}, []float64{-1000, 2000}))
	assert.ErrorIs(t, err, ErrNoConvergence)
}

// TestPerpetuityEdgeCases covers the boundaries of the Gordon model: the rate
// must be positive and must exceed the growth, otherwise the series diverges
// and no finite value exists.
func TestPerpetuityEdgeCases(t *testing.T) {
	payment := usd(100)

	_, err := Perpetuity(payment, decimal.Zero)
	assert.ErrorIs(t, err, ErrNonPositiveRate)

	_, err = Perpetuity(payment, decimal.MustFromFloat64(-0.05))
	assert.ErrorIs(t, err, ErrNonPositiveRate)

	// Growth equal to the rate diverges.
	_, err = GrowingPerpetuity(payment, decimal.MustFromFloat64(0.05), decimal.MustFromFloat64(0.05))
	assert.ErrorIs(t, err, ErrRateBelowGrowth)

	// Growth above the rate diverges too.
	_, err = GrowingPerpetuity(payment, decimal.MustFromFloat64(0.05), decimal.MustFromFloat64(0.08))
	assert.ErrorIs(t, err, ErrRateBelowGrowth)

	// Negative growth is a shrinking but convergent series.
	value, err := GrowingPerpetuity(payment, decimal.MustFromFloat64(0.05), decimal.MustFromFloat64(-0.02))
	require.NoError(t, err)
	// PV = 100 / (0.05 + 0.02) = 1428.5714
	assert.InDelta(t, 1428.5714, value.InexactFloat64(), 0.001)

	assert.Panics(t, func() { MustGrowingPerpetuity(payment, decimal.MustFromFloat64(0.05), decimal.MustFromFloat64(0.08)) })
	assert.NotPanics(t, func() { MustGrowingPerpetuity(payment, decimal.MustFromFloat64(0.08), decimal.MustFromFloat64(0.05)) })
}

// TestPerpetuityMatchesLongAnnuity checks the perpetuity against its own
// definition: it is the limit of an ever-longer level annuity, so a very long
// NPV must approach it.
func TestPerpetuityMatchesLongAnnuity(t *testing.T) {
	rate := decimal.MustFromFloat64(0.10)

	perpetual, err := Perpetuity(usd(100), rate)
	require.NoError(t, err)
	assert.InDelta(t, 1000.0, perpetual.InexactFloat64(), 1e-9)

	// 300 periods of 100 discounted at 10% is within a cent of the limit.
	long := make([]money.Money, 0, 301)
	long = append(long, usd(0))

	for range 300 {
		long = append(long, usd(100))
	}

	npv, err := NPV(rate, long)
	require.NoError(t, err)
	assert.InDelta(t, perpetual.InexactFloat64(), npv.InexactFloat64(), 0.01)
}
