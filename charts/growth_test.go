package charts

import (
	"strings"
	"testing"

	echartslib "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

func buildTestGrowthSchedule(t *testing.T, currency money.Currency, presentAmount float64) []compoundinterest.GrowthSchedule {
	t.Helper()

	present := money.MustMoneyFromFloat64(presentAmount, currency)

	rows, err := compoundinterest.BuildGrowthSchedule(present, decimal.MustFromFloat64(0.01), decimal.MustFromFloat64(12))
	if err != nil {
		t.Fatalf("failed to build test growth schedule: %v", err)
	}
	return rows
}

func TestGrowthXAxis(t *testing.T) {
	schedule := buildTestGrowthSchedule(t, money.USD, 1000)

	xAxis := GrowthXAxis(schedule)
	if len(xAxis) != len(schedule) {
		t.Fatalf("expected %d entries, got %d", len(schedule), len(xAxis))
	}
	for i, v := range xAxis {
		if v != i {
			t.Errorf("expected index %d, got %v", i, v)
		}
	}

	if empty := GrowthXAxis(nil); len(empty) != 0 {
		t.Errorf("expected empty slice for nil schedule, got %v", empty)
	}
}

func TestGrowthSeriesRoundToCurrencyPrecision(t *testing.T) {
	jpySchedule := buildTestGrowthSchedule(t, money.JPY, 1000)
	usdSchedule := buildTestGrowthSchedule(t, money.USD, 1000)

	seriesFns := map[string]func([]compoundinterest.GrowthSchedule) ([]opts.LineData, error){
		"Balance":     GrowthBalance,
		"Change":      GrowthChange,
		"SumInterest": GrowthSumInterest,
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

func TestGrowthSeriesEmptySchedule(t *testing.T) {
	seriesFns := map[string]func([]compoundinterest.GrowthSchedule) ([]opts.LineData, error){
		"Balance":     GrowthBalance,
		"Change":      GrowthChange,
		"SumInterest": GrowthSumInterest,
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

	if items := GrowthChangePercent(nil); len(items) != 0 {
		t.Errorf("GrowthChangePercent(nil): expected empty result, got %d items", len(items))
	}
}

func TestGrowthChangePercentIsConstant(t *testing.T) {
	// Plain compound interest has a constant period-over-period percentage
	// change (it's always the configured rate), unlike an investment with
	// contributions. Compared numerically rather than by rendered string,
	// since money rounded to cents can make the last few digits of the
	// division wobble (e.g. 0.0099999999999999999 vs 0.01) without the
	// value actually differing at any meaningful precision.
	schedule := buildTestGrowthSchedule(t, money.USD, 1000)

	items := GrowthChangePercent(schedule)
	if len(items) != len(schedule) {
		t.Fatalf("expected %d items, got %d", len(schedule), len(items))
	}

	const tolerance = 1e-9
	reference := schedule[1].ChangePercent.InexactFloat64()
	for i := 2; i < len(schedule); i++ {
		got := schedule[i].ChangePercent.InexactFloat64()
		if diff := got - reference; diff < -tolerance || diff > tolerance {
			t.Errorf("expected constant ChangePercent, row %d = %v, row 1 = %v", i, got, reference)
		}
	}
}

func TestGrowthChartBuilders(t *testing.T) {
	schedule := buildTestGrowthSchedule(t, money.USD, 1000)

	builders := map[string]func([]compoundinterest.GrowthSchedule, ...ChartOption) (*echartslib.Line, error){
		"GrowthChart":            GrowthChart,
		"GrowthBalanceOnlyChart": GrowthBalanceOnlyChart,
		"GrowthChangeChart":      GrowthChangeChart,
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
