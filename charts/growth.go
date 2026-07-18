// Package charts: this file renders compoundinterest.GrowthSchedule
// (plain compound interest growth, no periodic contributions) as
// go-echarts line charts.
package charts

import (
	echartslib "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
)

// changePercentPrecision is the number of decimal digits used to render
// ChangePercent series, which are dimensionless rates rather than currency
// amounts and so have no natural currency precision to round to.
const changePercentPrecision = 6

// GrowthXAxis returns the period indexes (0..len(schedule)-1) used as the X
// axis for every chart in this file.
func GrowthXAxis(schedule []compoundinterest.GrowthSchedule) []any {
	xAxis := make([]any, 0, len(schedule))

	for i := range schedule {
		xAxis = append(xAxis, i)
	}

	return xAxis
}

// GrowthBalance extracts the per-period balance series from schedule,
// rounding each value to its currency's own precision. It returns an error
// if the currency's precision cannot be determined.
func GrowthBalance(schedule []compoundinterest.GrowthSchedule) ([]opts.LineData, error) {
	items := make([]opts.LineData, 0, len(schedule))

	if len(schedule) == 0 {
		return items, nil
	}

	prec, err := schedule[0].Balance.Currency().GetCurrencyPrecisionCode()
	if err != nil {
		return nil, err
	}

	for _, s := range schedule {
		items = append(items, opts.LineData{Value: s.Balance.RoundBankString(prec)})
	}

	return items, nil
}

// GrowthChange extracts the per-period change series from schedule (each
// period's balance minus the previous one, i.e. the interest earned that
// period), rounding to its currency's own precision. It returns an error if
// the currency's precision cannot be determined.
func GrowthChange(schedule []compoundinterest.GrowthSchedule) ([]opts.LineData, error) {
	items := make([]opts.LineData, 0, len(schedule))

	if len(schedule) == 0 {
		return items, nil
	}

	prec, err := schedule[0].Balance.Currency().GetCurrencyPrecisionCode()
	if err != nil {
		return nil, err
	}

	for _, s := range schedule {
		items = append(items, opts.LineData{Value: s.Change.RoundBankString(prec)})
	}

	return items, nil
}

// GrowthChangePercent extracts the per-period percentage change series from
// schedule. For plain compound interest this is constant across periods
// (it's always the configured rate), unlike the equivalent series for an
// investment with contributions.
func GrowthChangePercent(schedule []compoundinterest.GrowthSchedule) []opts.LineData {
	items := make([]opts.LineData, 0, len(schedule))

	for _, s := range schedule {
		items = append(items, opts.LineData{Value: s.ChangePercent.StringFixed(changePercentPrecision)})
	}

	return items
}

// GrowthSumInterest extracts the per-period cumulative interest series from
// schedule, rounding each value to its currency's own precision. It returns
// an error if the currency's precision cannot be determined.
func GrowthSumInterest(schedule []compoundinterest.GrowthSchedule) ([]opts.LineData, error) {
	items := make([]opts.LineData, 0, len(schedule))

	if len(schedule) == 0 {
		return items, nil
	}

	prec, err := schedule[0].SumInterest.Currency().GetCurrencyPrecisionCode()
	if err != nil {
		return nil, err
	}

	for _, s := range schedule {
		items = append(items, opts.LineData{Value: s.SumInterest.RoundBankString(prec)})
	}

	return items, nil
}

// GrowthChart builds a line chart with the Balance and cumulative interest
// (SumInterest) series over the schedule's periods. opt overrides
// DefaultChartOption when provided.
func GrowthChart(schedule []compoundinterest.GrowthSchedule, opt ...ChartOption) (*echartslib.Line, error) {
	option := DefaultChartOption()
	if len(opt) > 0 {
		option = opt[0]
	}
	option.Title = "Compound Interest Growth"

	line := echartslib.NewLine()
	line.SetGlobalOptions(
		echartslib.WithInitializationOpts(opts.Initialization{Theme: option.Theme}),
		echartslib.WithTitleOpts(opts.Title{
			Title:    option.Title,
			Subtitle: option.Subtitle,
		}),
		echartslib.WithXAxisOpts(opts.XAxis{Name: option.XAxisName}),
		echartslib.WithYAxisOpts(opts.YAxis{Name: option.YAxisName}),
	)

	if option.ShowLegend {
		line.SetGlobalOptions(echartslib.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}))
	}

	xAxis := GrowthXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	balance, err := GrowthBalance(schedule)
	if err != nil {
		return nil, err
	}

	sumInterest, err := GrowthSumInterest(schedule)
	if err != nil {
		return nil, err
	}

	line.AddSeries("Balance", balance).
		AddSeries("Total Interest", sumInterest)

	return line, nil
}

// GrowthBalanceOnlyChart builds a line chart with only the Balance series
// over the schedule's periods. opt overrides DefaultChartOption when
// provided.
func GrowthBalanceOnlyChart(schedule []compoundinterest.GrowthSchedule, opt ...ChartOption) (*echartslib.Line, error) {
	option := DefaultChartOption()
	if len(opt) > 0 {
		option = opt[0]
	}
	option.Title = "Balance Over Time"

	line := echartslib.NewLine()
	line.SetGlobalOptions(
		echartslib.WithInitializationOpts(opts.Initialization{Theme: option.Theme}),
		echartslib.WithTitleOpts(opts.Title{Title: option.Title}),
		echartslib.WithXAxisOpts(opts.XAxis{Name: option.XAxisName}),
		echartslib.WithYAxisOpts(opts.YAxis{Name: option.YAxisName}),
	)

	xAxis := GrowthXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	balance, err := GrowthBalance(schedule)
	if err != nil {
		return nil, err
	}

	line.AddSeries("Balance", balance)

	return line, nil
}

// GrowthChangeChart builds a line chart with only the per-period Change
// series (the interest earned each period, which grows over time even
// though the rate is constant) over the schedule's periods. opt overrides
// DefaultChartOption when provided.
func GrowthChangeChart(schedule []compoundinterest.GrowthSchedule, opt ...ChartOption) (*echartslib.Line, error) {
	option := DefaultChartOption()
	if len(opt) > 0 {
		option = opt[0]
	}
	option.Title = "Interest Earned Per Period"

	line := echartslib.NewLine()
	line.SetGlobalOptions(
		echartslib.WithInitializationOpts(opts.Initialization{Theme: option.Theme}),
		echartslib.WithTitleOpts(opts.Title{Title: option.Title}),
		echartslib.WithXAxisOpts(opts.XAxis{Name: option.XAxisName}),
		echartslib.WithYAxisOpts(opts.YAxis{Name: option.YAxisName}),
	)

	xAxis := GrowthXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	change, err := GrowthChange(schedule)
	if err != nil {
		return nil, err
	}

	line.AddSeries("Interest Earned", change)

	return line, nil
}
