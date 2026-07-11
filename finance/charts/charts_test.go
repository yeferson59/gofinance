package charts

import (
	"strings"
	"testing"

	echartslib "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/yeferson59/gofinance/finance/annuities"
	"github.com/yeferson59/gofinance/money"
)

func buildTestSchedule(t *testing.T, currency money.Currency, pvAmount, paymentAmount float64) []annuities.Schedule {
	t.Helper()

	pv := money.MustMoneyFromFloat64(pvAmount, currency)
	payment := money.MustMoneyFromFloat64(paymentAmount, currency)

	rows, err := annuities.BuildSchedule(pv, money.MustFromFloat64(0.01), payment, money.MustFromFloat64(12))
	if err != nil {
		t.Fatalf("failed to build test schedule: %v", err)
	}
	return rows
}

func TestDefaultChartOption(t *testing.T) {
	opt := DefaultChartOption()

	if opt.Title != "Amortization Chart" {
		t.Errorf("expected title %q, got %q", "Amortization Chart", opt.Title)
	}
	if !opt.ShowLegend {
		t.Error("expected ShowLegend to be true by default")
	}
	if opt.Smooth {
		t.Error("expected Smooth to be false by default")
	}
	if opt.XAxisName != "Period" || opt.YAxisName != "Amount" {
		t.Errorf("unexpected axis names: %q / %q", opt.XAxisName, opt.YAxisName)
	}
}

func TestGenerateXAxis(t *testing.T) {
	schedule := buildTestSchedule(t, money.USD, 1000, 90)

	xAxis := GenerateXAxis(schedule)
	if len(xAxis) != len(schedule) {
		t.Fatalf("expected %d entries, got %d", len(schedule), len(xAxis))
	}
	for i, v := range xAxis {
		if v != i {
			t.Errorf("expected index %d, got %v", i, v)
		}
	}

	if empty := GenerateXAxis(nil); len(empty) != 0 {
		t.Errorf("expected empty slice for nil schedule, got %v", empty)
	}
}

func TestSeriesRoundToCurrencyPrecision(t *testing.T) {
	// JPY has zero decimal places, unlike the hardcoded 2 this package
	// used to assume.
	jpySchedule := buildTestSchedule(t, money.JPY, 1000, 90)
	// USD has two decimal places, the previous hardcoded default.
	usdSchedule := buildTestSchedule(t, money.USD, 1000, 90)

	seriesFns := map[string]func([]annuities.Schedule) ([]opts.LineData, error){
		"Balance":       Balance,
		"Principal":     Principal,
		"Interest":      Interest,
		"TotalInterest": TotalInterest,
		"Payment":       Payment,
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

func TestSeriesEmptySchedule(t *testing.T) {
	seriesFns := map[string]func([]annuities.Schedule) ([]opts.LineData, error){
		"Balance":       Balance,
		"Principal":     Principal,
		"Interest":      Interest,
		"TotalInterest": TotalInterest,
		"Payment":       Payment,
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
}

func TestChartBuilders(t *testing.T) {
	schedule := buildTestSchedule(t, money.USD, 1000, 90)

	builders := map[string]func([]annuities.Schedule, ...ChartOption) (*echartslib.Line, error){
		"AmortizationChart":         AmortizationChart,
		"BalanceOnlyChart":          BalanceOnlyChart,
		"PrincipalVsInterestChart":  PrincipalVsInterestChart,
		"InterestAccumulationChart": InterestAccumulationChart,
		"StackedAreaChart":          StackedAreaChart,
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
}

func TestAmortizationChartIncludesPaymentSeries(t *testing.T) {
	schedule := buildTestSchedule(t, money.USD, 1000, 90)

	line, err := AmortizationChart(schedule)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found := false
	for _, s := range line.MultiSeries {
		if s.Name == "Payment" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected AmortizationChart to include a Payment series")
	}
}
