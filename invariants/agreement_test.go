package invariants

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/annuities"
	"github.com/yeferson59/gofinance/v2/finance/bonds"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/finance/gradients"
	"github.com/yeferson59/gofinance/v2/finance/investment"
	"github.com/yeferson59/gofinance/v2/finance/loans"
	"github.com/yeferson59/gofinance/v2/finance/returns"
	"github.com/yeferson59/gofinance/v2/finance/tvm"
	"github.com/yeferson59/gofinance/v2/money"
)

func usd(amount float64) money.Money {
	return money.MustMoneyFromFloat64(amount, money.USD)
}

// newAnnuity builds a monthly annuity with a periodic rate, the configuration
// most of the cross-checks below need.
func newAnnuity(t *testing.T, payment, present, future, rate float64, periods int) annuities.Annuity {
	t.Helper()

	period, err := compoundinterest.NewPeriod(
		decimal.MustFromFloat64(float64(periods)), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(rate), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	annuity, err := annuities.New(usd(payment), usd(present), usd(future), period, rateInterest)
	require.NoError(t, err)

	return annuity
}

// TestLoanPaymentAgreesAcrossPackages checks the library's three routes to the
// same number. A level-payment loan can be priced through the general
// time-value solver, through the annuity formula, or through the loan builder;
// all three must return the same payment, or one of them is wrong.
func TestLoanPaymentAgreesAcrossPackages(t *testing.T) {
	cases := []struct {
		name      string
		principal float64
		rate      float64
		periods   int
	}{
		{"30-year mortgage", 300000, 0.005, 360},
		{"car loan", 25000, 0.004, 60},
		{"short term", 5000, 0.01, 12},
		{"single payment", 1000, 0.02, 1},
		{"low rate", 100000, 0.0001, 240},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fromTVM, err := tvm.NewTVM().PV(test.principal).Rate(test.rate).N(float64(test.periods)).SolvePMT()
			require.NoError(t, err)

			fromAnnuities, err := newAnnuity(t, 0, test.principal, 0, test.rate, test.periods).
				PaymentFromPresentValue()
			require.NoError(t, err)

			fromLoans, err := loans.NewLoan().
				Principal(test.principal, money.USD).
				Rate(test.rate).
				Periods(test.periods).
				Payment()
			require.NoError(t, err)

			// tvm returns the payment as an outflow; the other two as a
			// positive amount owed.
			assert.InDelta(t, -fromTVM.InexactFloat64(), fromAnnuities.InexactFloat64(), 0.01,
				"tvm and annuities must agree")
			assert.InDelta(t, fromAnnuities.InexactFloat64(), fromLoans.InexactFloat64(), 0.01,
				"annuities and loans must agree")
		})
	}
}

// TestAnnuityPresentValueIsNPVOfItsPayments checks the annuity formula against
// the definition it is a closed form of: discounting the payments one by one
// with finance/investment must give the same present value.
func TestAnnuityPresentValueIsNPVOfItsPayments(t *testing.T) {
	const (
		payment = 500.0
		rate    = 0.01
		periods = 24
	)

	fromAnnuities, err := newAnnuity(t, payment, 0, 0, rate, periods).Present()
	require.NoError(t, err)

	// NPV treats index 0 as undiscounted, so an ordinary annuity's flows start
	// with a zero at t = 0.
	flows := make([]money.Money, 0, periods+1)
	flows = append(flows, usd(0))

	for range periods {
		flows = append(flows, usd(payment))
	}

	fromInvestment, err := investment.NPV(decimal.MustFromFloat64(rate), flows)
	require.NoError(t, err)

	assert.InDelta(t, fromInvestment.InexactFloat64(), fromAnnuities.InexactFloat64(), 0.01)
}

// TestBondPriceIsNPVOfItsCashFlows checks finance/bonds against
// finance/investment: a bond's clean price is by definition the present value
// of its coupons and redemption at the per-period yield.
func TestBondPriceIsNPVOfItsCashFlows(t *testing.T) {
	const (
		face       = 1000.0
		couponRate = 0.06
		frequency  = 2
		periods    = 20
		yield      = 0.08
	)

	bond := bonds.NewBond().Face(face, money.USD).CouponRate(couponRate).
		Frequency(frequency).Periods(periods).Yield(yield)

	price, err := bond.Price()
	require.NoError(t, err)

	coupon, err := bond.CouponPayment()
	require.NoError(t, err)

	flows := make([]money.Money, 0, periods+1)
	flows = append(flows, usd(0))

	for period := 1; period <= periods; period++ {
		amount := coupon.InexactFloat64()
		if period == periods {
			amount += face
		}

		flows = append(flows, usd(amount))
	}

	npv, err := investment.NPV(decimal.MustFromFloat64(yield/frequency), flows)
	require.NoError(t, err)

	assert.InDelta(t, npv.InexactFloat64(), price.InexactFloat64(), 0.01)
}

// TestBondYTMIsIRROfItsCashFlows checks the other direction: buying the bond at
// its price and collecting its cash flows has an internal rate of return equal
// to the per-period yield.
func TestBondYTMIsIRROfItsCashFlows(t *testing.T) {
	const (
		face       = 1000.0
		couponRate = 0.05
		frequency  = 2
		periods    = 10
		yield      = 0.07
	)

	bond := bonds.NewBond().Face(face, money.USD).CouponRate(couponRate).
		Frequency(frequency).Periods(periods).Yield(yield)

	price, err := bond.Price()
	require.NoError(t, err)

	coupon, err := bond.CouponPayment()
	require.NoError(t, err)

	flows := make([]money.Money, 0, periods+1)
	flows = append(flows, usd(-price.InexactFloat64()))

	for period := 1; period <= periods; period++ {
		amount := coupon.InexactFloat64()
		if period == periods {
			amount += face
		}

		flows = append(flows, usd(amount))
	}

	irr, err := investment.IRR(flows)
	require.NoError(t, err)

	assert.InDelta(t, yield/frequency, irr.InexactFloat64(), 1e-6)
}

// TestGradientWithoutGrowthIsAnOrdinaryAnnuity checks finance/gradients
// degenerates into finance/annuities: a series that never changes is a level
// annuity, whether the change is expressed as a constant amount or a constant
// percentage.
func TestGradientWithoutGrowthIsAnOrdinaryAnnuity(t *testing.T) {
	const (
		payment = 1000.0
		rate    = 0.08
		periods = 10
	)

	period, err := compoundinterest.NewPeriod(
		decimal.MustFromFloat64(periods), compoundinterest.Annually)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(rate), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	annuity, err := annuities.New(usd(payment), usd(0), usd(0), period, rateInterest)
	require.NoError(t, err)

	annuityPresent, err := annuity.Present()
	require.NoError(t, err)

	annuityFuture, err := annuity.Future()
	require.NoError(t, err)

	arithmetic, err := gradients.NewArithmetic(usd(payment), usd(0), period, rateInterest)
	require.NoError(t, err)

	arithmeticPresent, err := arithmetic.Present()
	require.NoError(t, err)

	arithmeticFuture, err := arithmetic.Future()
	require.NoError(t, err)

	assert.InDelta(t, annuityPresent.InexactFloat64(), arithmeticPresent.InexactFloat64(), 0.01)
	assert.InDelta(t, annuityFuture.InexactFloat64(), arithmeticFuture.InexactFloat64(), 0.01)

	geometric, err := gradients.NewGeometric(usd(payment), decimal.Zero, period, rateInterest)
	require.NoError(t, err)

	geometricPresent, err := geometric.Present()
	require.NoError(t, err)

	geometricFuture, err := geometric.Future()
	require.NoError(t, err)

	assert.InDelta(t, annuityPresent.InexactFloat64(), geometricPresent.InexactFloat64(), 0.01)
	assert.InDelta(t, annuityFuture.InexactFloat64(), geometricFuture.InexactFloat64(), 0.01)
}

// TestGradientMatchesNPVOfItsPayments checks both gradient series against
// finance/investment, discounting the payments they describe one by one.
func TestGradientMatchesNPVOfItsPayments(t *testing.T) {
	const (
		first   = 1000.0
		rate    = 0.10
		periods = 8
	)

	period, err := compoundinterest.NewPeriod(
		decimal.MustFromFloat64(periods), compoundinterest.Annually)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(rate), compoundinterest.Annually, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	t.Run("arithmetic", func(t *testing.T) {
		const gradient = 150.0

		series, err := gradients.NewArithmetic(usd(first), usd(gradient), period, rateInterest)
		require.NoError(t, err)

		present, err := series.Present()
		require.NoError(t, err)

		flows := []money.Money{usd(0)}
		for step := range periods {
			flows = append(flows, usd(first+float64(step)*gradient))
		}

		npv, err := investment.NPV(decimal.MustFromFloat64(rate), flows)
		require.NoError(t, err)

		assert.InDelta(t, npv.InexactFloat64(), present.InexactFloat64(), 0.01)
	})

	t.Run("geometric", func(t *testing.T) {
		const growth = 0.06

		series, err := gradients.NewGeometric(usd(first), decimal.MustFromFloat64(growth), period, rateInterest)
		require.NoError(t, err)

		present, err := series.Present()
		require.NoError(t, err)

		flows := []money.Money{usd(0)}
		payment := first

		for range periods {
			flows = append(flows, usd(payment))
			payment *= 1 + growth
		}

		npv, err := investment.NPV(decimal.MustFromFloat64(rate), flows)
		require.NoError(t, err)

		assert.InDelta(t, npv.InexactFloat64(), present.InexactFloat64(), 0.01)
	})
}

// TestMoneyWeightedReturnIsIRR checks finance/returns delegates correctly: the
// money-weighted return of a portfolio is the internal rate of return of the
// investor's own cash flows.
//
// The two packages use opposite sign conventions, which is the point of
// checking them together. returns takes contributions as positive — the
// investor's point of view — and negates them internally; investment.IRR takes
// raw cash flows, where money leaving the investor is negative.
func TestMoneyWeightedReturnIsIRR(t *testing.T) {
	// Put in 10,000, add 2,000, withdraw 1,000, nothing, end worth 13,000.
	initial := usd(10000)
	interim := []money.Money{usd(2000), usd(-1000), usd(0)}
	final := usd(13000)

	moneyWeighted, err := returns.MoneyWeightedReturn(initial, interim, final)
	require.NoError(t, err)

	// The same story as raw cash flows.
	irr, err := investment.IRR([]money.Money{
		usd(-10000), usd(-2000), usd(1000), usd(0), usd(13000),
	})
	require.NoError(t, err)

	assert.InDelta(t, irr.InexactFloat64(), moneyWeighted.InexactFloat64(), 1e-9)
	assert.InDelta(t, 0.0426876, moneyWeighted.InexactFloat64(), 1e-6)
}

// TestCompoundFutureMatchesTVM checks finance/compoundinterest against
// finance/tvm on the one thing they both compute: growing a lump sum.
func TestCompoundFutureMatchesTVM(t *testing.T) {
	const (
		present = 10000.0
		rate    = 0.01
		periods = 36
	)

	period, err := compoundinterest.NewPeriod(
		decimal.MustFromFloat64(periods), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(rate), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	compound, err := compoundinterest.New(usd(present), usd(0), rateInterest, period)
	require.NoError(t, err)

	fromCompound, err := compound.Future()
	require.NoError(t, err)

	fromTVM, err := tvm.NewTVM().PV(-present).Rate(rate).N(periods).SolveFV()
	require.NoError(t, err)

	assert.InDelta(t, fromTVM.InexactFloat64(), fromCompound.InexactFloat64(), 0.01)
}

// TestLoanAPRExceedsNoteRateWithFees checks the relation between the two rates
// finance/loans reports: fees reduce the cash received without reducing the
// payments, so the APR must exceed the note rate, and must equal it when there
// are no fees.
func TestLoanAPRExceedsNoteRateWithFees(t *testing.T) {
	base := loans.NewLoan().Principal(200000, money.USD).AnnualRate(0.06).Years(30).Monthly()

	withoutFees, err := base.APR()
	require.NoError(t, err)
	assert.InDelta(t, 0.06, withoutFees.InexactFloat64(), 1e-6,
		"with no fees the APR is the note rate")

	withFees, err := base.Fees(5000).APR()
	require.NoError(t, err)
	assert.Greater(t, withFees.InexactFloat64(), withoutFees.InexactFloat64())

	// A bigger fee costs more, so the APR must rise with it.
	withMoreFees, err := base.Fees(10000).APR()
	require.NoError(t, err)
	assert.Greater(t, withMoreFees.InexactFloat64(), withFees.InexactFloat64())
}
