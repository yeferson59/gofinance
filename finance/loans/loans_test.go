package loans

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/finance/annuities"
	"github.com/yeferson59/gofinance/v2/money"
)

// mortgage is the running example of these tests: $250,000 over 30 years at a
// 6.5% nominal annual rate, paid monthly.
func mortgage() Config {
	return NewLoan().
		Principal(250000, money.USD).
		AnnualRate(0.065).
		Years(30).
		Monthly()
}

func TestPayment(t *testing.T) {
	payment, err := mortgage().Payment()
	require.NoError(t, err)
	assert.InDelta(t, 1580.170059, payment.InexactFloat64(), 1e-6)
}

func TestPaymentZeroRate(t *testing.T) {
	// With no interest the principal is just split evenly across the term.
	payment, err := NewLoan().Principal(1200, money.USD).Rate(0).Periods(12).Payment()
	require.NoError(t, err)
	assert.InDelta(t, 100.0, payment.InexactFloat64(), 1e-9)
}

func TestPeriodsTakePrecedenceOverYears(t *testing.T) {
	byYears := mortgage().MustPayment()
	byPeriods := mortgage().Periods(360).Years(15).MustPayment()

	assert.True(t, byYears.Equal(byPeriods))
}

func TestBuilderOrderDoesNotMatter(t *testing.T) {
	// AnnualRate is divided by the payment frequency whichever order the two
	// are set in.
	rateFirst := NewLoan().Principal(1000, money.USD).AnnualRate(0.12).Quarterly().Periods(4)
	freqFirst := NewLoan().Principal(1000, money.USD).Quarterly().AnnualRate(0.12).Periods(4)

	assert.InDelta(t, 0.03, rateFirst.MustPeriodicRate().InexactFloat64(), 1e-12)
	assert.True(t, rateFirst.MustPayment().Equal(freqFirst.MustPayment()))
}

func TestNumberOfPayments(t *testing.T) {
	n, err := mortgage().NumberOfPayments()
	require.NoError(t, err)
	assert.Equal(t, 360, n)

	n, err = NewLoan().Principal(1000, money.USD).Rate(0.01).Years(5).Quarterly().NumberOfPayments()
	require.NoError(t, err)
	assert.Equal(t, 20, n)
}

func TestEffectiveAnnualRate(t *testing.T) {
	// (1 + 0.065/12)^12 − 1.
	ear, err := mortgage().EffectiveAnnualRate()
	require.NoError(t, err)
	assert.InDelta(t, 0.06697185, ear.InexactFloat64(), 1e-8)
}

func TestNetProceeds(t *testing.T) {
	net, err := mortgage().Fees(3500).NetProceeds()
	require.NoError(t, err)
	assert.InDelta(t, 246500.0, net.InexactFloat64(), 1e-9)
	assert.Equal(t, money.USD, net.GetCurrency())
}

func TestAPRWithoutFeesIsTheNoteRate(t *testing.T) {
	loan := mortgage()

	// No finance charges, so the borrower gets the full principal and the APR
	// is the note rate itself — returned exactly, not bisected towards it.
	assert.True(t, loan.MustPeriodicAPR().Equal(loan.MustPeriodicRate()))
	assert.InDelta(t, 0.065, loan.MustAPR().InexactFloat64(), 1e-15)
}

func TestAPRWithFees(t *testing.T) {
	loan := mortgage().Fees(3500)

	apr, err := loan.APR()
	require.NoError(t, err)
	assert.InDelta(t, 0.06636063, apr.InexactFloat64(), 1e-8)

	// Fees can only make the loan dearer than its note rate.
	assert.Greater(t, apr.InexactFloat64(), 0.065)
}

func TestEffectiveAPR(t *testing.T) {
	loan := mortgage().Fees(3500)

	effective, err := loan.EffectiveAPR()
	require.NoError(t, err)
	assert.InDelta(t, 0.06841668, effective.InexactFloat64(), 1e-8)

	// Compounding the fee-inclusive rate beats compounding the note rate.
	assert.Greater(t, effective.InexactFloat64(), loan.MustEffectiveAnnualRate().InexactFloat64())
	assert.Greater(t, effective.InexactFloat64(), loan.MustAPR().InexactFloat64())
}

func TestAPRRecoversTheNetProceeds(t *testing.T) {
	loan := mortgage().Fees(3500)

	// Discounting the scheduled payments at the periodic APR must give back
	// exactly the cash the borrower received.
	factor, err := annuityFactor(loan.MustPeriodicAPR(), 360)
	require.NoError(t, err)

	// The rate is bisected to 1e-10, which over 360 periods leaves well under a
	// cent of the quarter-million discounted.
	present := loan.MustPayment().MulDecimal(factor)
	assert.InDelta(t, loan.MustNetProceeds().InexactFloat64(), present.InexactFloat64(), 0.01)
}

func TestAPRFeeErrors(t *testing.T) {
	_, err := mortgage().Fees(-100).APR()
	assert.ErrorIs(t, err, ErrNegativeFees)

	_, err = mortgage().Fees(250000).APR()
	assert.ErrorIs(t, err, ErrFeesExceedPrincipal)
}

func TestTermErrors(t *testing.T) {
	_, err := NewLoan().Principal(1000, money.USD).Rate(0.01).Payment()
	assert.ErrorIs(t, err, ErrInvalidPeriods)

	_, err = NewLoan().Principal(1000, money.USD).Rate(0.01).Periods(12).PaymentsPerYear(0).Payment()
	assert.ErrorIs(t, err, ErrInvalidFrequency)

	_, err = NewLoan().Principal(-1000, money.USD).Rate(0.01).Periods(12).Payment()
	assert.ErrorIs(t, err, ErrNonPositivePrincipal)

	_, err = NewLoan().Principal(1000, money.USD).Rate(-1.5).Periods(12).Payment()
	assert.ErrorIs(t, err, ErrInvalidRate)
}

func TestPayoffRetiresTheLoan(t *testing.T) {
	payoff, err := mortgage().Payoff()
	require.NoError(t, err)

	assert.Equal(t, 360, payoff.Periods)
	require.Len(t, payoff.Schedule, 361)

	// The opening row carries the original balance and the last one clears it.
	assert.InDelta(t, 250000.0, payoff.Schedule[0].Balance.InexactFloat64(), 1e-9)
	assert.InDelta(t, 0.0, payoff.Schedule[360].Balance.InexactFloat64(), 1e-9)

	assert.InDelta(t, 318861.221144, payoff.TotalInterest.InexactFloat64(), 1e-5)
	assert.InDelta(t, 568861.221144, payoff.TotalPaid.InexactFloat64(), 1e-5)

	// Principal repaid plus interest is what the borrower handed over.
	assert.InDelta(t,
		payoff.TotalPaid.InexactFloat64(),
		payoff.TotalInterest.InexactFloat64()+250000.0,
		1e-5)
}

func TestPayoffWithExtraPayment(t *testing.T) {
	payoff, err := mortgage().ExtraPayment(200).Payoff()
	require.NoError(t, err)

	assert.Equal(t, 265, payoff.Periods)
	assert.InDelta(t, 221243.104175, payoff.TotalInterest.InexactFloat64(), 1e-5)

	// Every regular payment carries the extra; the last one only settles what
	// is left, so it is smaller.
	assert.InDelta(t, 1780.170059, payoff.Schedule[1].Payment.InexactFloat64(), 1e-6)
	assert.Less(t, payoff.FinalPayment.InexactFloat64(), 1780.170059)
	assert.InDelta(t, 0.0, payoff.Schedule[265].Balance.InexactFloat64(), 1e-9)
}

func TestSavings(t *testing.T) {
	savings, err := mortgage().ExtraPayment(200).Savings()
	require.NoError(t, err)

	assert.Equal(t, 360, savings.Scheduled.Periods)
	assert.Equal(t, 265, savings.Accelerated.Periods)
	assert.Equal(t, 95, savings.PeriodsSaved)
	assert.InDelta(t, 97618.116969, savings.InterestSaved.InexactFloat64(), 1e-5)
}

func TestSavingsWithoutExtraPaymentAreZero(t *testing.T) {
	savings, err := mortgage().Savings()
	require.NoError(t, err)

	assert.Equal(t, 0, savings.PeriodsSaved)
	assert.True(t, savings.InterestSaved.IsZero())
}

func TestPayoffNegativeExtra(t *testing.T) {
	_, err := mortgage().ExtraPayment(-50).Payoff()
	assert.ErrorIs(t, err, ErrNegativeExtra)
}

func TestPayoffScheduleExportsAsCSV(t *testing.T) {
	// The rows are annuities.Schedule values, so the existing exporter takes
	// them unchanged.
	payoff := NewLoan().Principal(1000, money.USD).Rate(0.01).Periods(3).MustPayoff()

	var buf bytes.Buffer

	headers := []string{"period", "balance", "payment", "interest", "sum_interest", "principal"}
	require.NoError(t, annuities.WriteCSVTo(&buf, headers, payoff.Schedule))

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 5) // header + opening row + three payments
	assert.Equal(t, strings.Join(headers, ","), lines[0])
}

// refinanceOffers returns a $220,000 balance with 25 years left at 6.5%, and an
// offer to replace it at 5.25% for $4,000 of closing costs.
func refinanceOffers() (Config, Config) {
	current := NewLoan().Principal(220000, money.USD).AnnualRate(0.065).Years(25).Monthly()
	offer := NewLoan().Principal(220000, money.USD).AnnualRate(0.0525).Years(25).Monthly().Fees(4000)

	return current, offer
}

func TestCompare(t *testing.T) {
	comparison, err := Compare(refinanceOffers())
	require.NoError(t, err)

	assert.InDelta(t, 1485.455755, comparison.CurrentPayment.InexactFloat64(), 1e-6)
	assert.InDelta(t, 1318.344973, comparison.OfferPayment.InexactFloat64(), 1e-6)
	assert.InDelta(t, 167.110782, comparison.PaymentSavings.InexactFloat64(), 1e-6)
	assert.InDelta(t, 4000.0, comparison.ClosingCosts.InexactFloat64(), 1e-9)

	// 4000 / 167.11 = 23.9 periods of savings, so the 24th is the first that
	// leaves the borrower ahead.
	assert.Equal(t, 24, comparison.BreakEvenPeriods)

	assert.InDelta(t, 50133.234476, comparison.InterestSaved.InexactFloat64(), 1e-5)
	assert.Positive(t, comparison.NetPresentValue.InexactFloat64())
}

func TestCompareFreeOfferBreaksEvenImmediately(t *testing.T) {
	current, offer := refinanceOffers()

	comparison, err := Compare(current, offer.Fees(0))
	require.NoError(t, err)

	assert.Equal(t, 0, comparison.BreakEvenPeriods)
	assert.True(t, comparison.ClosingCosts.IsZero())
}

func TestCompareShorterTermCostsMoreEachPeriod(t *testing.T) {
	current, _ := refinanceOffers()

	// A 15-year replacement costs more every month, and only turns a profit
	// once it is paid off and the old loan would still have ten years to run.
	offer := NewLoan().Principal(220000, money.USD).AnnualRate(0.0525).Years(15).Monthly().Fees(4000)

	comparison, err := Compare(current, offer)
	require.NoError(t, err)

	assert.Negative(t, comparison.PaymentSavings.InexactFloat64())
	assert.Equal(t, 217, comparison.BreakEvenPeriods)
	assert.InDelta(t, 127301.151173, comparison.InterestSaved.InexactFloat64(), 1e-5)
}

func TestCompareRejectsWorseOffer(t *testing.T) {
	current, _ := refinanceOffers()
	worse := NewLoan().Principal(220000, money.USD).AnnualRate(0.08).Years(25).Monthly().Fees(4000)

	_, err := Compare(current, worse)
	assert.ErrorIs(t, err, ErrNoBreakEven)
}

func TestCompareMismatchErrors(t *testing.T) {
	current, offer := refinanceOffers()

	_, err := Compare(current, offer.Principal(220000, money.EUR))
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	_, err = Compare(current, offer.Quarterly().Years(25))
	assert.ErrorIs(t, err, ErrFrequencyMismatch)
}

func TestMustHelpersPanic(t *testing.T) {
	broken := NewLoan().Principal(1000, money.USD).Rate(0.01)

	assert.Panics(t, func() { broken.MustPayment() })
	assert.Panics(t, func() { broken.MustPeriodicRate() })
	assert.Panics(t, func() { broken.MustAPR() })
	assert.Panics(t, func() { broken.MustPeriodicAPR() })
	assert.Panics(t, func() { broken.MustEffectiveAPR() })
	assert.Panics(t, func() { broken.MustEffectiveAnnualRate() })
	assert.Panics(t, func() { broken.MustPayoff() })
	assert.Panics(t, func() { broken.MustSavings() })
	assert.Panics(t, func() { MustCompare(broken, broken) })

	// NetProceeds only looks at the fees, so it takes fees that swallow the
	// principal to make it fail.
	assert.Panics(t, func() { mortgage().Fees(250000).MustNetProceeds() })
}

func TestMustHelpersReturnValues(t *testing.T) {
	loan := mortgage().Fees(3500).ExtraPayment(200)

	assert.InDelta(t, 1580.170059, loan.MustPayment().InexactFloat64(), 1e-6)
	assert.InDelta(t, 246500.0, loan.MustNetProceeds().InexactFloat64(), 1e-9)
	assert.Equal(t, 265, loan.MustPayoff().Periods)
	assert.Equal(t, 95, loan.MustSavings().PeriodsSaved)

	current, offer := refinanceOffers()
	assert.Equal(t, 24, MustCompare(current, offer).BreakEvenPeriods)
}
