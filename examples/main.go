package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	echartslib "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/yeferson59/gofinance/finance/annuities"
	"github.com/yeferson59/gofinance/finance/charts"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/finance/daycount"
	dcf "github.com/yeferson59/gofinance/finance/investment"
	"github.com/yeferson59/gofinance/finance/returns"
	"github.com/yeferson59/gofinance/finance/simpleinterest"
	"github.com/yeferson59/gofinance/finance/tvm"
	"github.com/yeferson59/gofinance/money"
)

func main() {
	compoundExample()

	annuityExample()

	simpleExample()

	chartsExample()

	growthExample()

	investmentExample()

	returnsExample()

	dcfExample()

	tvmExample()

	daycountExample()
}

// returnsExample shows the finance/returns metrics: CAGR, ROI, and
// annualizing a cumulative return.
func returnsExample() {
	fmt.Println("\n=== Returns ===")

	begin := money.MustMoneyFromFloat64(1000, money.USD)
	end := money.MustMoneyFromFloat64(2000, money.USD)

	cagr := returns.MustCAGR(begin, end, money.MustFromFloat64(5))
	roi := returns.MustROI(begin, end)

	fmt.Println("$1,000 → $2,000 over 5 years")
	fmt.Println("CAGR:", cagr.RoundBank(4).StringFixed(4))
	fmt.Println("ROI:", roi.RoundBank(4).StringFixed(4))
}

// dcfExample shows the finance/investment discounted cash-flow metrics: net
// present value and internal rate of return over a stream of cash flows.
func dcfExample() {
	fmt.Println("\n=== Investment (NPV / IRR) ===")

	flows := []money.Money{
		money.MustMoneyFromFloat64(-1000, money.USD), // outflow today
		money.MustMoneyFromFloat64(400, money.USD),
		money.MustMoneyFromFloat64(400, money.USD),
		money.MustMoneyFromFloat64(400, money.USD),
	}

	npv := dcf.MustNPV(money.MustFromFloat64(0.10), flows)
	irr := dcf.MustIRR(flows)

	fmt.Println("Cash flows: -1000, 400, 400, 400")
	fmt.Println("NPV @ 10%:", npv.RoundBankString(2))
	fmt.Println("IRR:", irr.RoundBank(4).StringFixed(4))
}

// tvmExample shows the finance/tvm solver recovering the monthly payment of a
// loan from the other four time-value-of-money variables.
func tvmExample() {
	fmt.Println("\n=== Time Value of Money ===")

	payment := tvm.NewTVM().
		PV(300000).
		Rate(0.06 / 12).
		N(360).
		MustSolvePMT()

	fmt.Println("$300,000 loan at 6%/yr over 360 months")
	fmt.Println("Monthly payment:", payment.RoundBank(2).StringFixed(2))
}

// daycountExample shows the finance/daycount year fractions between two dates
// under different conventions.
func daycountExample() {
	fmt.Println("\n=== Day Count ===")

	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	for _, conv := range []daycount.Convention{
		daycount.Thirty360,
		daycount.Actual365Fixed,
		daycount.ActualActualISDA,
	} {
		yf, err := daycount.YearFraction(start, end, conv)
		if err != nil {
			fmt.Println("daycount error:", err)
			return
		}

		fmt.Printf("%-20s %s\n", conv, yf.RoundBank(6).StringFixed(6))
	}
}

func compoundExample() {
	fmt.Println("=== Composite Interest ===")

	ci := compositeinterest.NewComposite().
		Present(1000, money.USD).
		Rate(0.05).
		Periods(12).
		Monthly().
		RateType(compositeinterest.RateEffectyPeriodic).
		MustBuild()

	future, _ := ci.Future()
	fmt.Println("Present: $1000, Rate: 5%, 12 months")
	fmt.Println("Future value:", future.RoundBankString(2))
}

func annuityExample() {
	fmt.Println("\n=== Annuity Payment ===")

	// Compute the monthly periodic rate once and reuse it for both the
	// payment and the schedule below, instead of duplicating "0.06/12" as a
	// separate hardcoded literal — a drift between the two would silently
	// produce negative amortization (payment too small to cover interest).
	const annualRate = 0.06
	const periods = 360
	periodicRate := annualRate / 12

	payment := annuities.NewAnnuity().
		Present(300000, money.USD).
		Rate(periodicRate).
		Periods(periods).
		Monthly().
		MustPayment()

	fmt.Println("Loan: $300,000, Rate: 6%, 360 months")
	fmt.Println("Monthly payment:", payment.RoundBankString(2))

	schedule, err := annuities.BuildSchedule(
		money.MustMoneyFromFloat64(300000, money.USD),
		money.MustFromFloat64(periodicRate),
		payment,
		money.MustFromFloat64(periods),
	)
	if err != nil {
		fmt.Println("schedule error:", err)
		return
	}
	fmt.Println("Schedule rows:", len(schedule))
}

func simpleExample() {
	fmt.Println("\n=== Simple Interest ===")

	future, _ := simpleinterest.NewSimple().
		Present(5000, money.USD).
		AnnualRate(0.12).
		Periods(18).
		Months().
		FutureValue()

	fmt.Println("Present: $5000, Rate: 12%, 18 months")
	fmt.Println("Future value:", future.RoundBankString(2))
}

// chartsExample renders a variety of amortization charts using the
// finance/charts package and saves each one as a standalone HTML file
// under examples/output/. Open the files in a browser to view them.
func chartsExample() {
	fmt.Println("\n=== Charts ===")

	outDir := "output"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Println("failed to create output dir:", err)
		return
	}

	// Compute the monthly periodic rate once and reuse it for both the
	// payment and the schedule below, instead of duplicating "0.06/12" as a
	// separate hardcoded literal — a drift between the two would silently
	// produce negative amortization (payment too small to cover interest).
	const chartAnnualRate = 0.06
	const chartPeriods = 360
	chartPeriodicRate := chartAnnualRate / 12

	payment := annuities.NewAnnuity().
		Present(300000, money.USD).
		Rate(chartPeriodicRate).
		Periods(chartPeriods).
		Monthly().
		MustPayment()

	schedule, err := annuities.BuildSchedule(
		money.MustMoneyFromFloat64(300000, money.USD),
		money.MustFromFloat64(chartPeriodicRate),
		payment,
		money.MustFromFloat64(chartPeriods),
	)
	if err != nil {
		fmt.Println("schedule error:", err)
		return
	}

	// Default options: every series on one chart.
	full, err := charts.AmortizationChart(schedule)
	if err != nil {
		fmt.Println("AmortizationChart error:", err)
		return
	}
	saveChart(outDir, "amortization.html", full)

	// A single-series chart, useful for a compact dashboard widget.
	balanceOnly, err := charts.BalanceOnlyChart(schedule)
	if err != nil {
		fmt.Println("BalanceOnlyChart error:", err)
		return
	}
	saveChart(outDir, "balance-only.html", balanceOnly)

	// Comparing principal vs. interest per period.
	principalVsInterest, err := charts.PrincipalVsInterestChart(schedule)
	if err != nil {
		fmt.Println("PrincipalVsInterestChart error:", err)
		return
	}
	saveChart(outDir, "principal-vs-interest.html", principalVsInterest)

	// Custom ChartOption: different theme and smoothed lines.
	customOption := charts.DefaultChartOption()
	customOption.Theme = types.ThemeMacarons
	customOption.Smooth = true
	interestAccumulation, err := charts.InterestAccumulationChart(schedule, customOption)
	if err != nil {
		fmt.Println("InterestAccumulationChart error:", err)
		return
	}
	saveChart(outDir, "interest-accumulation.html", interestAccumulation)

	// Payment composition as a stacked area chart.
	stackedArea, err := charts.StackedAreaChart(schedule)
	if err != nil {
		fmt.Println("StackedAreaChart error:", err)
		return
	}
	saveChart(outDir, "payment-composition.html", stackedArea)

	// Non-2-decimal currency (JPY has 0 decimal places): each series is
	// rounded to JPY's own precision instead of assuming 2 decimals.
	const jpyAnnualRate = 0.03
	const jpyPeriods = 24
	jpyPeriodicRate := jpyAnnualRate / 12

	jpyPayment := annuities.NewAnnuity().
		Present(3000000, money.JPY).
		Rate(jpyPeriodicRate).
		Periods(jpyPeriods).
		Monthly().
		MustPayment()
	jpySchedule, err := annuities.BuildSchedule(
		money.MustMoneyFromFloat64(3000000, money.JPY),
		money.MustFromFloat64(jpyPeriodicRate),
		jpyPayment,
		money.MustFromFloat64(jpyPeriods),
	)
	if err != nil {
		fmt.Println("JPY schedule error:", err)
		return
	}
	jpyChart, err := charts.AmortizationChart(jpySchedule)
	if err != nil {
		fmt.Println("JPY AmortizationChart error:", err)
		return
	}
	saveChart(outDir, "amortization-jpy.html", jpyChart)

	fmt.Println("Charts saved to", outDir+"/")
}

// growthExample builds a plain compound interest growth schedule (a lump
// sum with no periodic contributions) and renders it with the
// finance/charts package as standalone HTML files under examples/output/.
func growthExample() {
	fmt.Println("\n=== Compound Interest Growth ===")

	outDir := "output"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Println("failed to create output dir:", err)
		return
	}

	present := money.MustMoneyFromFloat64(10000, money.USD)
	rate := money.MustFromFloat64(0.01)
	periods := money.MustFromFloat64(24)

	schedule, err := compositeinterest.BuildGrowthSchedule(present, rate, periods)
	if err != nil {
		fmt.Println("growth schedule error:", err)
		return
	}

	last := schedule[len(schedule)-1]
	fmt.Println("Present: $10,000, Rate: 1% monthly, 24 months")
	fmt.Println("Final balance:", last.Balance.RoundBankString(2))
	fmt.Println("Total interest earned:", last.SumInterest.RoundBankString(2))

	growth, err := charts.GrowthChart(schedule)
	if err != nil {
		fmt.Println("GrowthChart error:", err)
		return
	}
	saveChart(outDir, "growth.html", growth)

	balanceOnly, err := charts.GrowthBalanceOnlyChart(schedule)
	if err != nil {
		fmt.Println("GrowthBalanceOnlyChart error:", err)
		return
	}
	saveChart(outDir, "growth-balance-only.html", balanceOnly)

	// The interest earned per period, growing even though the rate itself
	// is constant — "interest on interest".
	change, err := charts.GrowthChangeChart(schedule)
	if err != nil {
		fmt.Println("GrowthChangeChart error:", err)
		return
	}
	saveChart(outDir, "growth-change.html", change)

	fmt.Println("Charts saved to", outDir+"/")
}

// investmentExample builds an investment growth schedule (compound interest
// plus a fixed contribution every period), for both ordinary (end of
// period) and anticipated (start of period) contribution timing, and
// renders them with the finance/charts package as standalone HTML files
// under examples/output/.
func investmentExample() {
	fmt.Println("\n=== Investment With Contributions ===")

	outDir := "output"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Println("failed to create output dir:", err)
		return
	}

	principal := money.MustMoneyFromFloat64(1000, money.USD)
	contribution := money.MustMoneyFromFloat64(100, money.USD)
	rate := money.MustFromFloat64(0.01)
	periods := money.MustFromFloat64(24)

	total := annuities.NewAnnuity().
		Present(1000, money.USD).
		Value(100, money.USD).
		Rate(0.01).
		Periods(24).
		Monthly().
		MustFutureValue()

	fmt.Println("Principal: $1,000, Contribution: $100/month, Rate: 1% monthly, 24 months")
	fmt.Println("Future value:", total.RoundBankString(2))

	schedule, err := annuities.BuildInvestmentSchedule(principal, contribution, rate, periods)
	if err != nil {
		fmt.Println("investment schedule error:", err)
		return
	}

	investment, err := charts.InvestmentChart(schedule)
	if err != nil {
		fmt.Println("InvestmentChart error:", err)
		return
	}
	saveChart(outDir, "investment.html", investment)

	balanceOnly, err := charts.InvestmentBalanceOnlyChart(schedule)
	if err != nil {
		fmt.Println("InvestmentBalanceOnlyChart error:", err)
		return
	}
	saveChart(outDir, "investment-balance-only.html", balanceOnly)

	// Cumulative contributions vs. cumulative interest, stacked.
	contributionVsInterest, err := charts.ContributionVsInterestChart(schedule)
	if err != nil {
		fmt.Println("ContributionVsInterestChart error:", err)
		return
	}
	saveChart(outDir, "investment-contribution-vs-interest.html", contributionVsInterest)

	// How the fixed contribution's relative weight shrinks over time as
	// the compounding balance grows past it.
	changePercent := charts.InvestmentChangePercentChart(schedule)
	saveChart(outDir, "investment-change-percent.html", changePercent)

	// Anticipated (due): contributions made at the start of each period
	// instead of the end, so they also earn interest in their own first
	// period — the balance ends up slightly higher than the ordinary case.
	anticipatedSchedule, err := annuities.BuildAnticipateInvestmentSchedule(principal, contribution, rate, periods)
	if err != nil {
		fmt.Println("anticipated investment schedule error:", err)
		return
	}

	anticipatedOption := charts.DefaultChartOption()
	anticipatedOption.Title = "Investment Growth (Anticipated)"
	anticipatedOption.Subtitle = "Contributions at the start of each period"

	anticipatedChart, err := charts.InvestmentChart(anticipatedSchedule, anticipatedOption)
	if err != nil {
		fmt.Println("anticipated InvestmentChart error:", err)
		return
	}
	saveChart(outDir, "investment-anticipated.html", anticipatedChart)

	fmt.Println("Charts saved to", outDir+"/")
}

func saveChart(dir, filename string, line *echartslib.Line) {
	path := filepath.Join(dir, filename)
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("failed to create", path, ":", err)
		return
	}
	defer f.Close()

	if err := line.Render(f); err != nil {
		fmt.Println("failed to render", path, ":", err)
		return
	}
	fmt.Println("wrote", path)
}
