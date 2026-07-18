package charts

import (
	"strings"
	"testing"

	echartslib "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/annuities"
	"github.com/yeferson59/gofinance/v2/money"
)

func buildTestInvestmentSchedule(t *testing.T, currency money.Currency, principalAmount, contributionAmount float64) []annuities.InvestmentSchedule {
	t.Helper()

	principal := money.MustMoneyFromFloat64(principalAmount, currency)
	contribution := money.MustMoneyFromFloat64(contributionAmount, currency)

	rows, err := annuities.BuildInvestmentSchedule(principal, contribution, decimal.MustFromFloat64(0.01), decimal.MustFromFloat64(12))
	if err != nil {
		t.Fatalf("failed to build test investment schedule: %v", err)
	}
	return rows
}

func TestInvestmentXAxis(t *testing.T) {
	schedule := buildTestInvestmentSchedule(t, money.USD, 1000, 100)

	xAxis := InvestmentXAxis(schedule)
	if len(xAxis) != len(schedule) {
		t.Fatalf("expected %d entries, got %d", len(schedule), len(xAxis))
	}
	for i, v := range xAxis {
		if v != i {
			t.Errorf("expected index %d, got %v", i, v)
		}
	}

	if empty := InvestmentXAxis(nil); len(empty) != 0 {
		t.Errorf("expected empty slice for nil schedule, got %v", empty)
	}
}

func TestInvestmentSeriesRoundToCurrencyPrecision(t *testing.T) {
	jpySchedule := buildTestInvestmentSchedule(t, money.JPY, 1000, 100)
	usdSchedule := buildTestInvestmentSchedule(t, money.USD, 1000, 100)

	seriesFns := map[string]func([]annuities.InvestmentSchedule) ([]opts.LineData, error){
		"Balance":          InvestmentBalance,
		"Change":           InvestmentChange,
		"SumContributions": InvestmentSumContributions,
		"SumInterest":      InvestmentSumInterest,
	}

	for name, fn := range seriesFns {
		jpySeries, err := fn(jpySchedule)
		if err != nil {
			t.Fatalf("%s(JPY): expected no error, got %v", name, err)
		}
		if len(jpySeries) != len(jpySchedule) {
			t.Fatalf("%s(JPY): expected %d items, got %d", name, len(jpySchedule), len(jpySeries))
		}
		for i, item := range jpySeries {
			s, ok := item.Value.(string)
			if !ok {
				t.Fatalf("%s(JPY) row %d: expected string value, got %T", name, i, item.Value)
			}
			if strings.Contains(s, ".") {
				t.Errorf("%s(JPY) row %d: expected value without decimals, got %q", name, i, s)
			}
		}

		usdSeries, err := fn(usdSchedule)
		if err != nil {
			t.Fatalf("%s(USD): expected no error, got %v", name, err)
		}
		for i, item := range usdSeries {
			s, ok := item.Value.(string)
			if !ok {
				t.Fatalf("%s(USD) row %d: expected string value, got %T", name, i, item.Value)
			}
			if idx := strings.Index(s, "."); idx == -1 || len(s)-idx-1 != 2 {
				t.Errorf("%s(USD) row %d: expected exactly 2 decimals, got %q", name, i, s)
			}
		}
	}
}

func TestInvestmentSeriesEmptySchedule(t *testing.T) {
	seriesFns := map[string]func([]annuities.InvestmentSchedule) ([]opts.LineData, error){
		"Balance":          InvestmentBalance,
		"Change":           InvestmentChange,
		"SumContributions": InvestmentSumContributions,
		"SumInterest":      InvestmentSumInterest,
	}

	for name, fn := range seriesFns {
		items, err := fn(nil)
		if err != nil {
			t.Errorf("%s(nil): expected no error, got %v", name, err)
		}
		if len(items) != 0 {
			t.Errorf("%s(nil): expected empty result, got %d items", name, len(items))
		}
	}

	if items := InvestmentChangePercent(nil); len(items) != 0 {
		t.Errorf("InvestmentChangePercent(nil): expected empty result, got %d items", len(items))
	}
}

func TestInvestmentChangePercentShrinksOverTime(t *testing.T) {
	// A fixed contribution is a bigger fraction of a still-small balance
	// early on, so ChangePercent should trend down as the compounding
	// balance grows past it (unlike plain compound interest, where it's
	// constant).
	schedule := buildTestInvestmentSchedule(t, money.USD, 0, 100)

	items := InvestmentChangePercent(schedule)
	if len(items) != len(schedule) {
		t.Fatalf("expected %d items, got %d", len(schedule), len(items))
	}

	last := len(items) - 1
	if schedule[last].ChangePercent.GreaterThanOrEqual(schedule[2].ChangePercent) {
		t.Errorf("expected ChangePercent to shrink over time, period 2 = %s, last period = %s",
			schedule[2].ChangePercent.String(), schedule[last].ChangePercent.String())
	}
}

func TestInvestmentChartBuilders(t *testing.T) {
	schedule := buildTestInvestmentSchedule(t, money.USD, 1000, 100)

	builders := map[string]func([]annuities.InvestmentSchedule, ...ChartOption) (*echartslib.Line, error){
		"InvestmentChart":             InvestmentChart,
		"InvestmentBalanceOnlyChart":  InvestmentBalanceOnlyChart,
		"ContributionVsInterestChart": ContributionVsInterestChart,
	}

	for name, build := range builders {
		line, err := build(schedule)
		if err != nil {
			t.Fatalf("%s: expected no error, got %v", name, err)
		}
		if line == nil {
			t.Fatalf("%s: expected non-nil chart", name)
		}
	}

	if line := InvestmentChangePercentChart(schedule); line == nil {
		t.Fatal("InvestmentChangePercentChart: expected non-nil chart")
	}
}

func TestInvestmentChartIncludesContributionSeries(t *testing.T) {
	schedule := buildTestInvestmentSchedule(t, money.USD, 1000, 100)

	line, err := InvestmentChart(schedule)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found := false
	for _, s := range line.MultiSeries {
		if s.Name == "Total Contributions" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected InvestmentChart to include a Total Contributions series")
	}
}
