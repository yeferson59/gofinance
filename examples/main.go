package main

import (
	"fmt"
	"os"
	"path/filepath"

	echartslib "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/yeferson59/gofinance/finance/annuities"
	"github.com/yeferson59/gofinance/finance/charts"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/finance/simpleinterest"
	"github.com/yeferson59/gofinance/money"
)

func main() {
	compoundExample()

	annuityExample()

	simpleExample()

	chartsExample()
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
	fmt.Println("Future value:", future.StringFixed(2))
}

func annuityExample() {
	fmt.Println("\n=== Annuity Payment ===")

	payment := annuities.NewAnnuity().
		Present(300000, money.USD).
		AnnualRate(0.06).
		Periods(360).
		Monthly().
		MustPayment()

	fmt.Println("Loan: $300,000, Rate: 6%, 360 months")
	fmt.Println("Monthly payment:", payment.StringFixed(2))

	schedule, err := annuities.BuildSchedule(
		money.MustMoneyFromFloat64(300000, money.USD),
		money.MustFromFloat64(0.005),
		payment,
		money.MustFromFloat64(360),
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
	fmt.Println("Future value:", future.StringFixed(2))
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

	payment := annuities.NewAnnuity().
		Present(300000, money.USD).
		AnnualRate(0.06).
		Periods(360).
		Monthly().
		MustPayment()

	schedule, err := annuities.BuildSchedule(
		money.MustMoneyFromFloat64(300000, money.USD),
		money.MustFromFloat64(0.005),
		payment,
		money.MustFromFloat64(360),
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
	jpyPayment := annuities.NewAnnuity().
		Present(3000000, money.JPY).
		AnnualRate(0.03).
		Periods(24).
		Monthly().
		MustPayment()
	jpySchedule, err := annuities.BuildSchedule(
		money.MustMoneyFromFloat64(3000000, money.JPY),
		money.MustFromFloat64(0.0025),
		jpyPayment,
		money.MustFromFloat64(24),
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
