package invariants

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/annuities"
	"github.com/yeferson59/gofinance/v2/finance/loans"
	"github.com/yeferson59/gofinance/v2/money"
)

// TestAmortizationScheduleCloses is the invariant an amortization table exists
// to satisfy: paying the solved payment for the full term must retire the
// balance exactly, the principal repaid must add up to the amount borrowed,
// and the interest column must account for the rest of what was paid.
func TestAmortizationScheduleCloses(t *testing.T) {
	cases := []struct {
		name      string
		principal float64
		rate      float64
		periods   int
	}{
		{"30-year mortgage", 300000, 0.005, 360},
		{"car loan", 25000, 0.004, 60},
		{"short term", 5000, 0.01, 12},
		{"zero rate", 12000, 0, 12},
		{"single payment", 1000, 0.02, 1},
		{"high rate", 10000, 0.05, 24},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			payment, err := newAnnuity(t, 0, test.principal, 0, test.rate, test.periods).
				PaymentFromPresentValue()
			require.NoError(t, err)

			schedule, err := annuities.BuildSchedule(
				usd(test.principal),
				decimal.MustFromFloat64(test.rate),
				payment,
				decimal.MustFromInt64(int64(test.periods), 0),
			)
			require.NoError(t, err)

			// A period-0 opening row plus one row per payment.
			require.Len(t, schedule, test.periods+1)
			assert.InDelta(t, test.principal, schedule[0].Balance.InexactFloat64(), 1e-9,
				"the opening row must carry the original balance")

			principalRepaid := usd(0)
			interestCharged := usd(0)

			for _, row := range schedule[1:] {
				principalRepaid = principalRepaid.Add(row.Principal)
				interestCharged = interestCharged.Add(row.Interest)
			}

			last := schedule[len(schedule)-1]

			assert.InDelta(t, 0.0, last.Balance.InexactFloat64(), 1e-6,
				"the balance must be retired")
			assert.InDelta(t, test.principal, principalRepaid.InexactFloat64(), 1e-6,
				"principal repaid must equal the amount borrowed")
			assert.InDelta(t, interestCharged.InexactFloat64(), last.SumInterest.InexactFloat64(), 1e-9,
				"the running interest total must match the sum of the column")

			// Every payment splits into exactly principal and interest.
			for _, row := range schedule[1:] {
				assert.InDelta(t,
					row.Payment.InexactFloat64(),
					row.Principal.InexactFloat64()+row.Interest.InexactFloat64(),
					1e-9, "period %v must split into principal and interest", row.Period)
			}
		})
	}
}

// TestScheduleBalanceDecreasesMonotonically checks the shape of a level-payment
// schedule: the balance falls every period, interest falls with it, and the
// principal portion grows.
func TestScheduleBalanceDecreasesMonotonically(t *testing.T) {
	const (
		principal = 200000.0
		rate      = 0.005
		periods   = 120
	)

	payment, err := newAnnuity(t, 0, principal, 0, rate, periods).PaymentFromPresentValue()
	require.NoError(t, err)

	schedule, err := annuities.BuildSchedule(
		usd(principal), decimal.MustFromFloat64(rate), payment, decimal.MustFromInt64(periods, 0))
	require.NoError(t, err)

	for i := 2; i < len(schedule); i++ {
		previous, current := schedule[i-1], schedule[i]

		assert.Less(t, current.Balance.InexactFloat64(), previous.Balance.InexactFloat64(),
			"balance must fall every period")
		assert.Less(t, current.Interest.InexactFloat64(), previous.Interest.InexactFloat64(),
			"interest must fall as the balance does")
		assert.Greater(t, current.Principal.InexactFloat64(), previous.Principal.InexactFloat64(),
			"the principal portion must grow")
	}
}

// TestLoanPayoffCloses checks the same closure through finance/loans, which
// amortizes with its own routine rather than annuities.BuildSchedule: the
// payments made must equal the principal plus the interest charged.
func TestLoanPayoffCloses(t *testing.T) {
	loan := loans.NewLoan().Principal(180000, money.USD).AnnualRate(0.06).Years(15).Monthly()

	payoff, err := loan.Payoff()
	require.NoError(t, err)

	assert.Equal(t, 180, payoff.Periods)
	assert.InDelta(t,
		180000+payoff.TotalInterest.InexactFloat64(),
		payoff.TotalPaid.InexactFloat64(), 0.01,
		"total paid must be the principal plus the interest")

	last := payoff.Schedule[len(payoff.Schedule)-1]
	assert.InDelta(t, 0.0, last.Balance.InexactFloat64(), 0.01)
}

// TestExtraPaymentShortensAndSaves checks the relations a prepayment must
// satisfy: it retires the loan sooner, charges less interest, and still repays
// the same principal.
func TestExtraPaymentShortensAndSaves(t *testing.T) {
	loan := loans.NewLoan().Principal(180000, money.USD).AnnualRate(0.06).Years(15).Monthly()

	savings, err := loan.ExtraPayment(300).Savings()
	require.NoError(t, err)

	assert.Positive(t, savings.PeriodsSaved, "an extra payment must shorten the loan")
	assert.Positive(t, savings.InterestSaved.InexactFloat64(), "and must save interest")

	assert.Less(t, savings.Accelerated.Periods, savings.Scheduled.Periods)
	assert.Less(t,
		savings.Accelerated.TotalInterest.InexactFloat64(),
		savings.Scheduled.TotalInterest.InexactFloat64())

	// The saving is exactly the difference between the two interest totals.
	assert.InDelta(t,
		savings.Scheduled.TotalInterest.InexactFloat64()-savings.Accelerated.TotalInterest.InexactFloat64(),
		savings.InterestSaved.InexactFloat64(), 0.01)

	// Both routes repay the same principal; only the interest differs.
	assert.InDelta(t,
		savings.Scheduled.TotalPaid.InexactFloat64()-savings.Scheduled.TotalInterest.InexactFloat64(),
		savings.Accelerated.TotalPaid.InexactFloat64()-savings.Accelerated.TotalInterest.InexactFloat64(),
		0.01)
}

// TestAllocationSumsBack is the invariant money.Allocate exists for: splitting
// an amount must neither lose nor invent money, whatever the ratios and
// whatever the currency's precision.
func TestAllocationSumsBack(t *testing.T) {
	amounts := []float64{100, 0.05, 0.10, 1000.01, -0.05, -1234.56, 0}
	ratioSets := [][]uint32{
		{1, 1, 1},
		{1, 2, 3},
		{1},
		{0, 1, 0, 1},
		{7, 11, 13, 17},
		{1, 1, 1, 1, 1, 1, 1},
	}
	currencies := []money.Currency{money.USD, money.JPY, money.EUR}

	for _, currency := range currencies {
		for _, amount := range amounts {
			for _, ratios := range ratioSets {
				original := money.MustMoneyFromFloat64(amount, currency)

				parts, err := original.Allocate(ratios...)
				require.NoError(t, err)
				require.Len(t, parts, len(ratios))

				total := money.MustMoneyFromFloat64(0, currency)
				for _, part := range parts {
					assert.Equal(t, currency, part.Currency())
					total = total.Add(part)
				}

				assert.True(t, total.Equal(original),
					"%v split %v in %v: parts sum to %v", original, ratios, currency, total)
			}
		}
	}
}

// TestAllocateEvenlySumsBack checks the same for the even split, including
// counts that do not divide the amount cleanly.
func TestAllocateEvenlySumsBack(t *testing.T) {
	for _, amount := range []float64{100, 0.01, 33.33, -10} {
		for _, parts := range []int{1, 2, 3, 7, 100} {
			original := money.MustMoneyFromFloat64(amount, money.USD)

			split, err := original.AllocateEvenly(parts)
			require.NoError(t, err)
			require.Len(t, split, parts)

			total := money.MustMoneyFromFloat64(0, money.USD)
			for _, part := range split {
				total = total.Add(part)
			}

			assert.True(t, total.Equal(original),
				"%v split evenly in %d: parts sum to %v", original, parts, total)
		}
	}
}

// TestAllocationPartsAreOrdered checks the documented remainder rule: the
// leftover smallest units go to the earliest ratios, so no part is more than
// one unit away from its neighbour's fair share.
func TestAllocationPartsAreOrdered(t *testing.T) {
	// 0.10 split three ways: 0.04, 0.03, 0.03.
	parts, err := usd(0.10).Allocate(1, 1, 1)
	require.NoError(t, err)

	assert.InDelta(t, 0.04, parts[0].InexactFloat64(), 1e-9)
	assert.InDelta(t, 0.03, parts[1].InexactFloat64(), 1e-9)
	assert.InDelta(t, 0.03, parts[2].InexactFloat64(), 1e-9)

	// A negative amount distributes the remainder the same way.
	negative, err := usd(-0.05).Allocate(1, 1, 1)
	require.NoError(t, err)

	assert.InDelta(t, -0.02, negative[0].InexactFloat64(), 1e-9)
	assert.InDelta(t, -0.02, negative[1].InexactFloat64(), 1e-9)
	assert.InDelta(t, -0.01, negative[2].InexactFloat64(), 1e-9)
}
