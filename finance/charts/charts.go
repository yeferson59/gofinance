// Package charts renders amortization schedules produced by
// finance/annuities as go-echarts line charts (balance, principal,
// interest, cumulative interest, and payment composition over time).
package charts

import (
	echartslib "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/yeferson59/gofinance/finance/annuities"
)

// ChartOption configures the appearance of a chart produced by this
// package: titles, axis labels, theme, line smoothing, and legend
// visibility.
type ChartOption struct {
	Title      string
	Subtitle   string
	Theme      string
	Smooth     bool
	ShowLegend bool
	XAxisName  string
	YAxisName  string
}

// DefaultChartOption returns the ChartOption used by every chart builder
// in this package when no explicit option is supplied.
func DefaultChartOption() ChartOption {
	return ChartOption{
		Title:      "Amortization Chart",
		Subtitle:   "Schedule visualization",
		Theme:      types.ThemeWesteros,
		Smooth:     false,
		ShowLegend: true,
		XAxisName:  "Period",
		YAxisName:  "Amount",
	}
}

// GenerateXAxis returns the period indexes (0..len(schedule)-1) used as the
// X axis for all charts in this package.
func GenerateXAxis(schedule []annuities.Schedule) []any {
	xAxis := make([]any, 0, len(schedule))

	for i := range schedule {
		xAxis = append(xAxis, i)
	}

	return xAxis
}

// Balance extracts the per-period balance series from schedule, rounding
// each value to its currency's own precision (e.g. 0 decimals for JPY, 3
// for BHD) rather than assuming two decimal places. It returns an error if
// the currency's precision cannot be determined.
func Balance(schedule []annuities.Schedule) ([]opts.LineData, error) {
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

// Principal extracts the per-period principal series from schedule,
// rounding each value to its currency's own precision. It returns an error
// if the currency's precision cannot be determined.
func Principal(schedule []annuities.Schedule) ([]opts.LineData, error) {
	items := make([]opts.LineData, 0, len(schedule))

	if len(schedule) == 0 {
		return items, nil
	}

	prec, err := schedule[0].Principal.Currency().GetCurrencyPrecisionCode()
	if err != nil {
		return nil, err
	}

	for _, s := range schedule {
		items = append(items, opts.LineData{Value: s.Principal.RoundBankString(prec)})
	}

	return items, nil
}

// Interest extracts the per-period interest series from schedule, rounding
// each value to its currency's own precision. It returns an error if the
// currency's precision cannot be determined.
func Interest(schedule []annuities.Schedule) ([]opts.LineData, error) {
	items := make([]opts.LineData, 0, len(schedule))

	if len(schedule) == 0 {
		return items, nil
	}

	prec, err := schedule[0].Interest.Currency().GetCurrencyPrecisionCode()
	if err != nil {
		return nil, err
	}

	for _, s := range schedule {
		items = append(items, opts.LineData{Value: s.Interest.RoundBankString(prec)})
	}

	return items, nil
}

// TotalInterest extracts the per-period cumulative interest series from
// schedule, rounding each value to its currency's own precision. It
// returns an error if the currency's precision cannot be determined.
func TotalInterest(schedule []annuities.Schedule) ([]opts.LineData, error) {
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

// Payment extracts the per-period payment series from schedule, rounding
// each value to its currency's own precision. It returns an error if the
// currency's precision cannot be determined.
func Payment(schedule []annuities.Schedule) ([]opts.LineData, error) {
	items := make([]opts.LineData, 0, len(schedule))

	if len(schedule) == 0 {
		return items, nil
	}

	prec, err := schedule[0].Payment.Currency().GetCurrencyPrecisionCode()
	if err != nil {
		return nil, err
	}

	for _, s := range schedule {
		items = append(items, opts.LineData{Value: s.Payment.RoundBankString(prec)})
	}

	return items, nil
}

// AmortizationChart builds a line chart with the Balance, Principal,
// Interest, Total Interest, and Payment series over the schedule's
// periods. opt overrides DefaultChartOption when provided.
func AmortizationChart(schedule []annuities.Schedule, opt ...ChartOption) (*echartslib.Line, error) {
	option := DefaultChartOption()
	if len(opt) > 0 {
		option = opt[0]
	}

	line := echartslib.NewLine()

	line.SetGlobalOptions(
		echartslib.WithInitializationOpts(opts.Initialization{Theme: option.Theme}),
		echartslib.WithTitleOpts(opts.Title{
			Title:    option.Title,
			Subtitle: option.Subtitle,
		}),
		echartslib.WithXAxisOpts(opts.XAxis{
			Name: option.XAxisName,
		}),
		echartslib.WithYAxisOpts(opts.YAxis{
			Name: option.YAxisName,
		}),
	)

	if option.ShowLegend {
		line.SetGlobalOptions(echartslib.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}))
	}

	xAxis := GenerateXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	balance, err := Balance(schedule)
	if err != nil {
		return nil, err
	}

	principal, err := Principal(schedule)
	if err != nil {
		return nil, err
	}

	interest, err := Interest(schedule)
	if err != nil {
		return nil, err
	}

	totalInterest, err := TotalInterest(schedule)
	if err != nil {
		return nil, err
	}

	payment, err := Payment(schedule)
	if err != nil {
		return nil, err
	}

	line.AddSeries("Balance", balance).
		AddSeries("Principal", principal).
		AddSeries("Interest", interest).
		AddSeries("Total Interest", totalInterest).
		AddSeries("Payment", payment)

	return line, nil
}

// BalanceOnlyChart builds a line chart with only the Balance series over
// the schedule's periods. opt overrides DefaultChartOption when provided.
func BalanceOnlyChart(schedule []annuities.Schedule, opt ...ChartOption) (*echartslib.Line, error) {
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

	xAxis := GenerateXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	balance, err := Balance(schedule)
	if err != nil {
		return nil, err
	}

	line.AddSeries("Balance", balance)

	return line, nil
}

// PrincipalVsInterestChart builds a line chart comparing the Principal and
// Interest series over the schedule's periods. opt overrides
// DefaultChartOption when provided.
func PrincipalVsInterestChart(schedule []annuities.Schedule, opt ...ChartOption) (*echartslib.Line, error) {
	option := DefaultChartOption()

	if len(opt) > 0 {
		option = opt[0]
	}

	option.Title = "Principal vs Interest"

	line := echartslib.NewLine()

	line.SetGlobalOptions(
		echartslib.WithInitializationOpts(opts.Initialization{Theme: option.Theme}),
		echartslib.WithTitleOpts(opts.Title{Title: option.Title}),
		echartslib.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
	)

	xAxis := GenerateXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	principal, err := Principal(schedule)
	if err != nil {
		return nil, err
	}

	interest, err := Interest(schedule)
	if err != nil {
		return nil, err
	}

	line.AddSeries("Principal", principal).AddSeries("Interest", interest)

	return line, nil
}

// InterestAccumulationChart builds a line chart with only the Total
// Interest series over the schedule's periods. opt overrides
// DefaultChartOption when provided.
func InterestAccumulationChart(schedule []annuities.Schedule, opt ...ChartOption) (*echartslib.Line, error) {
	option := DefaultChartOption()
	if len(opt) > 0 {
		option = opt[0]
	}

	option.Title = "Interest Accumulation"

	line := echartslib.NewLine()
	line.SetGlobalOptions(
		echartslib.WithInitializationOpts(opts.Initialization{Theme: option.Theme}),
		echartslib.WithTitleOpts(opts.Title{Title: option.Title}),
	)

	xAxis := GenerateXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}))

	totalInterest, err := TotalInterest(schedule)
	if err != nil {
		return nil, err
	}
	line.AddSeries("Total Interest", totalInterest)

	return line, nil
}

// StackedAreaChart builds an area chart showing the Principal and Interest
// composition of each payment over the schedule's periods. opt overrides
// DefaultChartOption when provided.
func StackedAreaChart(schedule []annuities.Schedule, opt ...ChartOption) (*echartslib.Line, error) {
	option := DefaultChartOption()
	if len(opt) > 0 {
		option = opt[0]
	}
	option.Title = "Payment Composition"

	line := echartslib.NewLine()
	line.SetGlobalOptions(
		echartslib.WithInitializationOpts(opts.Initialization{Theme: option.Theme}),
		echartslib.WithTitleOpts(opts.Title{Title: option.Title}),
		echartslib.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
	)

	xAxis := GenerateXAxis(schedule)
	line.SetXAxis(xAxis)
	line.SetSeriesOptions(
		echartslib.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(option.Smooth)}),
		echartslib.WithAreaStyleOpts(opts.AreaStyle{Opacity: opts.Float(0.5)}),
	)

	principal, err := Principal(schedule)
	if err != nil {
		return nil, err
	}
	interest, err := Interest(schedule)
	if err != nil {
		return nil, err
	}
	line.AddSeries("Principal", principal).
		AddSeries("Interest", interest)

	return line, nil
}
