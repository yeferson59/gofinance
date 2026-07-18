package tvm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSolvePMTLoan(t *testing.T) {
	// $300,000 loan at 0.5%/month over 360 months.
	pmt, err := NewTVM().PV(300000).Rate(0.005).N(360).SolvePMT()
	require.NoError(t, err)
	assert.InDelta(t, -1798.65, pmt.InexactFloat64(), 0.01)
}

func TestSolveFVSavings(t *testing.T) {
	// Depositing $100/month at 1%/month for 12 months.
	fv, err := NewTVM().PMT(-100).Rate(0.01).N(12).SolveFV()
	require.NoError(t, err)
	assert.InDelta(t, 1268.25, fv.InexactFloat64(), 0.01)
}

func TestSolvePVRoundTrip(t *testing.T) {
	// The PV that produces the loan payment must be the original principal.
	loan := NewTVM().PV(300000).Rate(0.005).N(360)
	pmt := loan.MustSolvePMT()

	pv, err := NewTVM().PMT(pmt.InexactFloat64()).Rate(0.005).N(360).SolvePV()
	require.NoError(t, err)
	assert.InDelta(t, 300000, pv.InexactFloat64(), 0.01)
}

func TestSolveN(t *testing.T) {
	// How many $100 monthly deposits at 1%/month reach the saved FV?
	fv := NewTVM().PMT(-100).Rate(0.01).N(12).MustSolveFV()

	n, err := NewTVM().PMT(-100).Rate(0.01).FV(fv.InexactFloat64()).SolveN()
	require.NoError(t, err)
	assert.InDelta(t, 12.0, n.InexactFloat64(), 1e-6)
}

func TestSolveRate(t *testing.T) {
	// Recover the loan's periodic rate from its payment.
	pmt := NewTVM().PV(300000).Rate(0.005).N(360).MustSolvePMT()

	rate, err := NewTVM().PV(300000).PMT(pmt.InexactFloat64()).N(360).SolveRate()
	require.NoError(t, err)
	assert.InDelta(t, 0.005, rate.InexactFloat64(), 1e-6)
}

func TestZeroRateFV(t *testing.T) {
	// At a 0% rate, borrowing 1000 and repaying 100 for 10 periods clears it.
	fv, err := NewTVM().PV(-1000).PMT(100).N(10).Rate(0).SolveFV()
	require.NoError(t, err)
	assert.InDelta(t, 0.0, fv.InexactFloat64(), 1e-9)
}

func TestZeroRateN(t *testing.T) {
	n, err := NewTVM().PV(-1000).PMT(100).FV(0).Rate(0).SolveN()
	require.NoError(t, err)
	assert.InDelta(t, 10.0, n.InexactFloat64(), 1e-9)
}

func TestAnnuityDue(t *testing.T) {
	ordinary := NewTVM().PMT(-100).Rate(0.01).N(12).MustSolveFV()
	due := NewTVM().PMT(-100).Rate(0.01).N(12).Due().MustSolveFV()

	// An annuity due earns one extra period of interest: FV_due = FV_ord·(1+i).
	assert.InDelta(t, ordinary.InexactFloat64()*1.01, due.InexactFloat64(), 0.01)
	assert.True(t, due.GreaterThan(ordinary))
}

func TestSolvePMTNoSolution(t *testing.T) {
	_, err := NewTVM().PV(1000).Rate(0.01).N(0).SolvePMT()
	assert.ErrorIs(t, err, ErrNoSolution)
}

func TestSolveNNoSolution(t *testing.T) {
	// Zero payment at a zero rate: nothing changes the balance.
	_, err := NewTVM().PV(-1000).PMT(0).Rate(0).SolveN()
	assert.ErrorIs(t, err, ErrNoSolution)
}

func TestInvalidRate(t *testing.T) {
	_, err := NewTVM().PV(1000).Rate(-2).N(12).SolveFV()
	assert.ErrorIs(t, err, ErrInvalidRate)
}
