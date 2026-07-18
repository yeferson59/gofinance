// Package charts: this file renders annuities.InvestmentSchedule (compound
// interest plus periodic contributions) as go-echarts line and area charts.
package charts

import (
	echartslib "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/yeferson59/gofinance/v2/finance/annuities"
)

// InvestmentXAxis returns the period indexes (0..len(schedule)-1) used as
// the X axis for every chart in this file.
func InvestmentXAxis(schedule []annuities.InvestmentSchedule) []any {
	xAxis := make([]any, 0, len(schedule))

	for i := range schedule {
		xAxis = append(xAxis, i)
	}

	return xAxis
}

// InvestmentBalance extracts the per-period balance series from schedule,
// rounding each value to its currency's own precision. It returns an error
// if the currency's precision cannot be determined.
func InvestmentBalance(schedule []annuities.InvestmentSchedule) ([]opts.LineData, error) {
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

// InvestmentChange extracts the per-period change series from schedule
// (each period's balance minus the previous one, i.e. that period's
// contribution plus interest earned), rounding to its currency's own
// precision. It returns an error if the currency's precision cannot be
// determined.
func InvestmentChange(schedule []annuities.InvestmentSchedule) ([]opts.LineData, error) {
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

// InvestmentChangePercent extracts the per-period percentage change series
// from schedule. Unlike plain compound interest, this varies over time: a
// fixed contribution is a larger fraction of a still-small balance early
// on, so it starts high and converges toward the periodic rate as the
// compounding balance grows past it.
func InvestmentChangePercent(schedule []annuities.InvestmentSchedule) []opts.LineData {
	items := make([]opts.LineData, 0, len(schedule))

	for _, s := range schedule {
		items = append(items, opts.LineData{Value: s.ChangePercent.StringFixed(changePercentPrecision)})
	}

	return items
}

// InvestmentSumContributions extracts the per-period cumulative
// contributions series from schedule, rounding each value to its
// currency's own precision. It returns an error if the currency's precision
// cannot be determined.
func InvestmentSumContributions(schedule []annuities.InvestmentSchedule) ([]opts.LineData, error) {
	items := make([]opts.LineData, 0, len(schedule))

	if len(schedule) == 0 {
		return items, nil
	}

	prec, err := schedule[0].SumContributions.Currency().GetCurrencyPrecisionCode()
	if err != nil {
		return nil, err
	}

	for _, s := range schedule {
		items = append(items, opts.LineData{Value: s.SumContributions.RoundBankString(prec)})
	}

	return items, nil
}

// InvestmentSumInterest extracts the per-period cumulative interest series
// from schedule, rounding each value to its currency's own precision. It
// returns an error if the currency's precision cannot be determined.
func InvestmentSumInterest(schedule []annuities.InvestmentSchedule) ([]opts.LineData, error) {
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

// InvestmentChart builds a line chart with the Balance, cumulative
// contributions (SumContributions), and cumulative interest (SumInterest)
// series over the schedule's periods, showing how much of the final
// balance came from contributions versus compounding. opt overrides
// DefaultChartOption when provided.
func InvestmentChart(schedule []annuities.InvestmentSchedule, opt ...ChartOption) (*echartslib.Line, error) {
	option := DefaultChartOption()
	if len(opt) > 0 {
		option = opt[0]
	}
	option.Title = "Investment Growth"

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

	xAxis := InvestmentXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	balance, err := InvestmentBalance(schedule)
	if err != nil {
		return nil, err
	}

	sumContributions, err := InvestmentSumContributions(schedule)
	if err != nil {
		return nil, err
	}

	sumInterest, err := InvestmentSumInterest(schedule)
	if err != nil {
		return nil, err
	}

	line.AddSeries("Balance", balance).
		AddSeries("Total Contributions", sumContributions).
		AddSeries("Total Interest", sumInterest)

	return line, nil
}

// InvestmentBalanceOnlyChart builds a line chart with only the Balance
// series over the schedule's periods. opt overrides DefaultChartOption when
// provided.
func InvestmentBalanceOnlyChart(schedule []annuities.InvestmentSchedule, opt ...ChartOption) (*echartslib.Line, error) {
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

	xAxis := InvestmentXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	balance, err := InvestmentBalance(schedule)
	if err != nil {
		return nil, err
	}

	line.AddSeries("Balance", balance)

	return line, nil
}

// ContributionVsInterestChart builds a stacked area chart comparing
// cumulative contributions and cumulative interest over the schedule's
// periods, mirroring StackedAreaChart's principal-vs-interest breakdown for
// loans. opt overrides DefaultChartOption when provided.
func ContributionVsInterestChart(schedule []annuities.InvestmentSchedule, opt ...ChartOption) (*echartslib.Line, error) {
	option := DefaultChartOption()
	if len(opt) > 0 {
		option = opt[0]
	}
	option.Title = "Contributions vs Interest"

	line := echartslib.NewLine()
	line.SetGlobalOptions(
		echartslib.WithInitializationOpts(opts.Initialization{Theme: option.Theme}),
		echartslib.WithTitleOpts(opts.Title{Title: option.Title}),
		echartslib.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
	)

	xAxis := InvestmentXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(
		echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}),
		echartslib.WithAreaStyleOpts(opts.AreaStyle{Opacity: opts.Float(0.5)}),
	)

	sumContributions, err := InvestmentSumContributions(schedule)
	if err != nil {
		return nil, err
	}

	sumInterest, err := InvestmentSumInterest(schedule)
	if err != nil {
		return nil, err
	}

	line.AddSeries("Total Contributions", sumContributions).
		AddSeries("Total Interest", sumInterest)

	return line, nil
}

// InvestmentChangePercentChart builds a line chart with only the
// ChangePercent series over the schedule's periods, visualizing how a fixed
// contribution's relative weight shrinks over time as the compounding
// balance grows past it. opt overrides DefaultChartOption when provided.
func InvestmentChangePercentChart(schedule []annuities.InvestmentSchedule, opt ...ChartOption) *echartslib.Line {
	option := DefaultChartOption()
	if len(opt) > 0 {
		option = opt[0]
	}
	option.Title = "Period-over-Period Growth Rate"

	line := echartslib.NewLine()
	line.SetGlobalOptions(
		echartslib.WithInitializationOpts(opts.Initialization{Theme: option.Theme}),
		echartslib.WithTitleOpts(opts.Title{Title: option.Title}),
		echartslib.WithXAxisOpts(opts.XAxis{Name: option.XAxisName}),
		echartslib.WithYAxisOpts(opts.YAxis{Name: "Growth (%)"}),
	)

	xAxis := InvestmentXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	line.AddSeries("Change %", InvestmentChangePercent(schedule))

	return line
}
