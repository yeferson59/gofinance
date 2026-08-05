package invariants

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/bonds"
	"github.com/yeferson59/gofinance/v2/finance/investment"
	"github.com/yeferson59/gofinance/v2/finance/loans"
	"github.com/yeferson59/gofinance/v2/finance/returns"
	"github.com/yeferson59/gofinance/v2/finance/tvm"
	"github.com/yeferson59/gofinance/v2/money"
)

// The sweeps below feed inputs that are valid by type but stress the decimal
// engine's range: huge and tiny amounts, absurd terms, rates at the edge of the
// domain. The invariant under test is narrow and deliberate:
//
//	a function that returns an error must never panic.
//
// It says nothing about the fluent builders' float64 setters — Face, PMT,
// ExtraPayment and friends. Those return a Config, not an error, so a value
// the decimal engine cannot represent has no channel but a panic, which is the
// documented contract of the Must-style constructors behind them. The sweeps
// therefore only feed setters values that are exactly representable, and hold
// the terminal, error-returning calls to the no-panic rule.
//
// This exists because Phase 2 found four separate panics of exactly this shape
// (TESTING_PLAN.md §2.8), all from the panicking decimal helpers used inside
// error-returning code. Rather than wait to trip over the next one, it hunts
// for the pattern across the packages — and found four more.
var extremeAmounts = []float64{
	0, 1e-8, 1e15, 1e30, 1e35, -1e30,
}

var extremeRates = []float64{
	0, 1e-12, 0.5, 10, 1000, -0.9999,
}

var extremeTerms = []int{1, 2, 1200, 6000}

// TestLoansNeverPanic sweeps the loan API. finance/loans had not been audited
// for the panicking-helper pattern.
func TestLoansNeverPanic(t *testing.T) {
	for _, principal := range extremeAmounts {
		for _, rate := range extremeRates {
			for _, periods := range extremeTerms {
				loan := loans.NewLoan().
					Principal(principal, money.USD).
					Rate(rate).
					Periods(periods).
					Monthly()

				assert.NotPanics(t, func() {
					_, _ = loan.Payment()
					_, _ = loan.PeriodicRate()
					_, _ = loan.NumberOfPayments()
					_, _ = loan.EffectiveAnnualRate()
					_, _ = loan.NetProceeds()
					_, _ = loan.APR()
					_, _ = loan.PeriodicAPR()
					_, _ = loan.EffectiveAPR()
					_, _ = loan.Payoff()
				}, "principal %v rate %v periods %d", principal, rate, periods)

				// The setters take float64 and panic on a value the decimal
				// engine cannot represent, so feed them exact powers of ten
				// rather than a float division that lands on a 30-digit
				// fraction. The rule under test is about the terminal calls.
				assert.NotPanics(t, func() {
					_, _ = loan.Fees(1000).APR()
					_, _ = loan.ExtraPayment(100).Savings()
				}, "with fees: principal %v rate %v periods %d", principal, rate, periods)
			}
		}
	}
}

// TestLoanComparisonNeverPanics sweeps the refinance comparison, which runs two
// amortizations and an NPV over their difference.
func TestLoanComparisonNeverPanics(t *testing.T) {
	for _, rate := range extremeRates {
		for _, periods := range []int{1, 360, 6000} {
			current := loans.NewLoan().Principal(200000, money.USD).Rate(rate).Periods(periods).Monthly()
			offer := loans.NewLoan().Principal(200000, money.USD).Rate(rate / 2).Periods(periods).Monthly().Fees(3000)

			assert.NotPanics(t, func() {
				_, _ = loans.Compare(current, offer)
			}, "rate %v periods %d", rate, periods)
		}
	}
}

// TestReturnsNeverPanic sweeps the return and risk metrics with extreme series,
// including the degenerate ones: empty, single-element, and all-identical.
func TestReturnsNeverPanic(t *testing.T) {
	series := [][]decimal.Decimal{
		{},
		{decimal.MustFromFloat64(0.05)},
		{decimal.Zero, decimal.Zero, decimal.Zero},
		{decimal.MustFromFloat64(-1), decimal.MustFromFloat64(-1)},
		{decimal.MustFromFloat64(1e30), decimal.MustFromFloat64(-1e30)},
		{decimal.MustFromFloat64(0.01), decimal.MustFromFloat64(-0.02), decimal.MustFromFloat64(0.03)},
	}

	for _, rates := range series {
		assert.NotPanics(t, func() {
			_, _ = returns.Mean(rates)
			_, _ = returns.Variance(rates)
			_, _ = returns.PopulationVariance(rates)
			_, _ = returns.Volatility(rates)
			_, _ = returns.PopulationVolatility(rates)
			_, _ = returns.ChainReturns(rates)
		}, "series %v", rates)

		// AnnualizedVolatility scales an already-computed volatility, so feed
		// it both a sane figure and the extremes of the range.
		for _, volatility := range []float64{0, 0.2, 1e30} {
			for _, periodsPerYear := range []float64{0, 1, 252, 1e12} {
				assert.NotPanics(t, func() {
					_, _ = returns.AnnualizedVolatility(
						decimal.MustFromFloat64(volatility), decimal.MustFromFloat64(periodsPerYear))
				}, "volatility %v periodsPerYear %v", volatility, periodsPerYear)
			}
		}

		for _, periodsPerYear := range []float64{0, 1, 252, 1e12} {
			assert.NotPanics(t, func() {
				_, _ = returns.SharpeRatio(rates, decimal.Zero)
				_, _ = returns.AnnualizedSharpeRatio(rates, decimal.Zero, decimal.MustFromFloat64(periodsPerYear))
			}, "series %v periodsPerYear %v", rates, periodsPerYear)
		}
	}
}

// TestReturnMetricsNeverPanic sweeps the point-in-time return metrics with
// extreme amounts and horizons, including a total loss and a zero starting
// value.
func TestReturnMetricsNeverPanic(t *testing.T) {
	for _, begin := range extremeAmounts {
		for _, end := range extremeAmounts {
			for _, years := range []float64{0, 1e-9, 1, 1000} {
				assert.NotPanics(t, func() {
					_, _ = returns.CAGR(money.MustMoneyFromFloat64(begin, money.USD),
						money.MustMoneyFromFloat64(end, money.USD), decimal.MustFromFloat64(years))
					_, _ = returns.ROI(money.MustMoneyFromFloat64(begin, money.USD),
						money.MustMoneyFromFloat64(end, money.USD))
					_, _ = returns.HoldingPeriodReturn(money.MustMoneyFromFloat64(begin, money.USD),
						money.MustMoneyFromFloat64(end, money.USD), money.MustMoneyFromFloat64(0, money.USD))
					_, _ = returns.Annualized(decimal.MustFromFloat64(0.1), decimal.MustFromFloat64(years))
				}, "begin %v end %v years %v", begin, end, years)
			}
		}
	}
}

// TestInflationMetricsNeverPanic sweeps the inflation adjustments, including
// the degenerate −100% inflation that makes the deflator zero.
func TestInflationMetricsNeverPanic(t *testing.T) {
	for _, amount := range extremeAmounts {
		for _, inflation := range []float64{0, -1, -0.5, 0.02, 1e6} {
			for _, years := range []float64{0, 1, 500} {
				assert.NotPanics(t, func() {
					_, _ = returns.RealValue(money.MustMoneyFromFloat64(amount, money.USD),
						decimal.MustFromFloat64(inflation), decimal.MustFromFloat64(years))
					_, _ = returns.NominalValue(money.MustMoneyFromFloat64(amount, money.USD),
						decimal.MustFromFloat64(inflation), decimal.MustFromFloat64(years))
					_, _ = returns.RealRate(decimal.MustFromFloat64(0.08), decimal.MustFromFloat64(inflation))
				}, "amount %v inflation %v years %v", amount, inflation, years)
			}
		}
	}
}

// TestSolversNeverPanic re-sweeps the packages Phase 2 fixed, so the panics
// found there cannot come back through a different door.
func TestSolversNeverPanic(t *testing.T) {
	for _, amount := range extremeAmounts {
		for _, rate := range extremeRates {
			for _, periods := range extremeTerms {
				assert.NotPanics(t, func() {
					config := tvm.NewTVM().PV(amount).PMT(-100).
						Rate(rate).N(float64(periods))

					_, _ = config.SolveFV()
					_, _ = config.SolvePV()
					_, _ = config.SolvePMT()
					_, _ = config.SolveN()
					_, _ = config.SolveRate()
					_, _ = config.Due().SolveFV()
				}, "tvm: amount %v rate %v periods %d", amount, rate, periods)

				assert.NotPanics(t, func() {
					bond := bonds.NewBond().Face(amount, money.USD).CouponRate(rate).
						Frequency(2).Periods(periods).Yield(rate)

					_, _ = bond.Price()
					_, _ = bond.CouponPayment()
					_, _ = bond.MacaulayDuration()
					_, _ = bond.ModifiedDuration()
					_, _ = bond.Convexity()
					_, _ = bond.MarketPrice(amount).YTM()
				}, "bonds: amount %v rate %v periods %d", amount, rate, periods)
			}
		}
	}
}

// TestInvestmentNeverPanics sweeps the cash-flow metrics with long series and
// extreme rates.
func TestInvestmentNeverPanics(t *testing.T) {
	for _, count := range []int{1, 2, 400, 2000} {
		for _, amount := range []float64{1, 1e15, 1e30} {
			flows := make([]money.Money, 0, count+1)
			flows = append(flows, money.MustMoneyFromFloat64(-amount, money.USD))

			for range count {
				flows = append(flows, money.MustMoneyFromFloat64(amount/10, money.USD))
			}

			assert.NotPanics(t, func() {
				_, _ = investment.IRR(flows)

				for _, rate := range extremeRates {
					_, _ = investment.NPV(decimal.MustFromFloat64(rate), flows)
				}
			}, "count %d amount %v", count, amount)
		}
	}
}

// TestDatedInvestmentNeverPanics sweeps the date-based metrics, where the year
// offsets can be large enough to overflow the discount factor.
func TestDatedInvestmentNeverPanics(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, spanDays := range []int{1, 365, 365 * 100, 365 * 3000} {
		for _, amount := range []float64{1, 1e30} {
			flows := []investment.DatedCashFlow{
				{Date: base, Amount: money.MustMoneyFromFloat64(-amount, money.USD)},
				{Date: base.AddDate(0, 0, spanDays/2), Amount: money.MustMoneyFromFloat64(amount/2, money.USD)},
				{Date: base.AddDate(0, 0, spanDays), Amount: money.MustMoneyFromFloat64(amount, money.USD)},
			}

			assert.NotPanics(t, func() {
				_, _ = investment.XIRR(flows)

				for _, rate := range extremeRates {
					_, _ = investment.XNPV(decimal.MustFromFloat64(rate), flows)
				}
			}, "span %d days, amount %v", spanDays, amount)
		}
	}
}
