package gradients

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// seriesConfig holds the period/rate parameters shared by ArithmeticConfig
// and GeometricConfig.
type seriesConfig struct {
	firstPayment money.Money
	periods      int
	rate         decimal.Decimal
	annual       bool
	frequency    compoundinterest.CompoundingFrequency
	rateType     compoundinterest.TypeRate

	// err holds the first invalid input a setter received (a NaN or infinite
	// float, say). Setters can't return an error without breaking the chain,
	// so it is kept here and reported by Build.
	err error
}

// record keeps err as the builder's error unless one is already held, so the
// first problem is the one reported. A nil err changes nothing.
func (c *seriesConfig) record(err error) {
	if c.err == nil {
		c.err = err
	}
}

// moneyInput converts a float amount for a setter, recording a conversion
// failure instead of panicking.
func (c *seriesConfig) moneyInput(amount float64, currency money.Currency) money.Money {
	m, err := money.NewMoneyFromFloat64(amount, currency)
	c.record(err)

	return m
}

// decimalInput converts a float rate for a setter, recording a conversion
// failure instead of panicking.
func (c *seriesConfig) decimalInput(r float64) decimal.Decimal {
	d, err := decimal.NewFromFloat64(r)
	c.record(err)

	return d
}

func newSeriesConfig() seriesConfig {
	return seriesConfig{
		frequency: compoundinterest.Monthly,
		rateType:  compoundinterest.RateEffectyPeriodic,
	}
}

// periodicRate returns the rate per period, dividing an annual rate by the
// configured frequency's periods per year. The division happens here rather
// than in AnnualRate so it uses the frequency the builder ends up with,
// whatever order the methods were called in.
//
// The periods per year come from the shared finance/term vocabulary rather
// than a table of this package's own, so the two cannot disagree.
func (c seriesConfig) periodicRate() (decimal.Decimal, error) {
	rate := c.rate
	if !c.annual {
		return rate, nil
	}

	periodsPerYear, err := c.frequency.PeriodsPerYear()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return rate.Div(periodsPerYear)
}

func (c seriesConfig) buildPeriodAndRate() (compoundinterest.Period, compoundinterest.RateInterest, error) {
	if c.err != nil {
		return compoundinterest.Period{}, compoundinterest.RateInterest{}, c.err
	}

	periodic, err := c.periodicRate()
	if err != nil {
		return compoundinterest.Period{}, compoundinterest.RateInterest{}, err
	}

	rateInterest, err := compoundinterest.NewRateInterest(periodic, c.frequency, c.rateType)
	if err != nil {
		return compoundinterest.Period{}, compoundinterest.RateInterest{}, err
	}

	periods, err := decimal.NewFromInt64(int64(c.periods), 0)
	if err != nil {
		return compoundinterest.Period{}, compoundinterest.RateInterest{}, err
	}

	period, err := compoundinterest.NewPeriod(periods, c.frequency)
	if err != nil {
		return compoundinterest.Period{}, compoundinterest.RateInterest{}, err
	}

	return period, rateInterest, nil
}

// ArithmeticConfig is a fluent builder for Arithmetic gradient series.
//
// Example:
//
//	present, err := NewArithmeticSeries().
//	    FirstPayment(1000, money.USD).
//	    Gradient(100, money.USD).
//	    AnnualRate(0.12).
//	    Periods(10).
//	    Monthly().
//	    Present()
type ArithmeticConfig struct {
	seriesConfig
	gradient money.Money
}

// NewArithmeticSeries creates a new ArithmeticConfig builder with Monthly
// frequency and a periodic rate type as defaults.
func NewArithmeticSeries() ArithmeticConfig {
	return ArithmeticConfig{seriesConfig: newSeriesConfig()}
}

// FirstPayment sets the payment made at the end of the first period.
func (a ArithmeticConfig) FirstPayment(amount float64, currency money.Currency) ArithmeticConfig {
	a.firstPayment = a.moneyInput(amount, currency)
	return a
}

// Gradient sets the constant amount added to each subsequent payment
// (negative for a decreasing series).
func (a ArithmeticConfig) Gradient(amount float64, currency money.Currency) ArithmeticConfig {
	a.gradient = a.moneyInput(amount, currency)
	return a
}

// Periods sets the total number of payments in the series.
func (a ArithmeticConfig) Periods(n int) ArithmeticConfig {
	a.periods = n
	return a
}

// Rate sets the periodic interest rate directly. It replaces any rate set
// with AnnualRate.
func (a ArithmeticConfig) Rate(r float64) ArithmeticConfig {
	a.rate = a.decimalInput(r)
	a.annual = false

	return a
}

// AnnualRate sets the annual interest rate, which is divided by the
// configured frequency's periods per year to obtain the periodic rate,
// whatever order the builder methods are called in. It replaces any rate set
// with Rate.
func (a ArithmeticConfig) AnnualRate(r float64) ArithmeticConfig {
	a.rate = a.decimalInput(r)
	a.annual = true
	a.rateType = compoundinterest.RateEffectyPeriodic

	return a
}

// Monthly sets the compounding/payment frequency to monthly.
func (a ArithmeticConfig) Monthly() ArithmeticConfig {
	a.frequency = compoundinterest.Monthly
	return a
}

// Quarterly sets the compounding/payment frequency to quarterly.
func (a ArithmeticConfig) Quarterly() ArithmeticConfig {
	a.frequency = compoundinterest.Quarterly
	return a
}

// Annually sets the compounding/payment frequency to annually.
func (a ArithmeticConfig) Annually() ArithmeticConfig {
	a.frequency = compoundinterest.Annually
	return a
}

// Build creates the Arithmetic series from the configured parameters.
func (a ArithmeticConfig) Build() (Arithmetic, error) {
	period, rateInterest, err := a.buildPeriodAndRate()
	if err != nil {
		return Arithmetic{}, err
	}

	return NewArithmetic(a.firstPayment, a.gradient, period, rateInterest)
}

// MustBuild is like Build, but panics if the calculation fails.
func (a ArithmeticConfig) MustBuild() Arithmetic {
	g, err := a.Build()
	if err != nil {
		panic(err)
	}

	return g
}

// Present builds the series and returns its present value.
func (a ArithmeticConfig) Present() (money.Money, error) {
	g, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}

	return g.Present()
}

// MustPresent is like Present, but panics if the calculation fails.
func (a ArithmeticConfig) MustPresent() money.Money {
	m, err := a.Present()
	if err != nil {
		panic(err)
	}

	return m
}

// Future builds the series and returns its future value.
func (a ArithmeticConfig) Future() (money.Money, error) {
	g, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}

	return g.Future()
}

// MustFuture is like Future, but panics if the calculation fails.
func (a ArithmeticConfig) MustFuture() money.Money {
	m, err := a.Future()
	if err != nil {
		panic(err)
	}

	return m
}

// GeometricConfig is a fluent builder for Geometric gradient series.
//
// Example:
//
//	present, err := NewGeometricSeries().
//	    FirstPayment(1000, money.USD).
//	    GrowthRate(0.08).
//	    AnnualRate(0.12).
//	    Periods(10).
//	    Annually().
//	    Present()
type GeometricConfig struct {
	seriesConfig
	growthRate decimal.Decimal
}

// NewGeometricSeries creates a new GeometricConfig builder with Monthly
// frequency and a periodic rate type as defaults.
func NewGeometricSeries() GeometricConfig {
	return GeometricConfig{seriesConfig: newSeriesConfig()}
}

// FirstPayment sets the payment made at the end of the first period.
func (g GeometricConfig) FirstPayment(amount float64, currency money.Currency) GeometricConfig {
	g.firstPayment = g.moneyInput(amount, currency)
	return g
}

// GrowthRate sets the constant percentage each subsequent payment grows by,
// as a decimal (e.g. 0.08 for 8%; negative for a shrinking series).
func (g GeometricConfig) GrowthRate(r float64) GeometricConfig {
	g.growthRate = g.decimalInput(r)
	return g
}

// Periods sets the total number of payments in the series.
func (g GeometricConfig) Periods(n int) GeometricConfig {
	g.periods = n
	return g
}

// Rate sets the periodic interest rate directly. It replaces any rate set
// with AnnualRate.
func (g GeometricConfig) Rate(r float64) GeometricConfig {
	g.rate = g.decimalInput(r)
	g.annual = false

	return g
}

// AnnualRate sets the annual interest rate, which is divided by the
// configured frequency's periods per year to obtain the periodic rate,
// whatever order the builder methods are called in. It replaces any rate set
// with Rate.
func (g GeometricConfig) AnnualRate(r float64) GeometricConfig {
	g.rate = g.decimalInput(r)
	g.annual = true
	g.rateType = compoundinterest.RateEffectyPeriodic

	return g
}

// Monthly sets the compounding/payment frequency to monthly.
func (g GeometricConfig) Monthly() GeometricConfig {
	g.frequency = compoundinterest.Monthly
	return g
}

// Quarterly sets the compounding/payment frequency to quarterly.
func (g GeometricConfig) Quarterly() GeometricConfig {
	g.frequency = compoundinterest.Quarterly
	return g
}

// Annually sets the compounding/payment frequency to annually.
func (g GeometricConfig) Annually() GeometricConfig {
	g.frequency = compoundinterest.Annually
	return g
}

// Build creates the Geometric series from the configured parameters.
func (g GeometricConfig) Build() (Geometric, error) {
	period, rateInterest, err := g.buildPeriodAndRate()
	if err != nil {
		return Geometric{}, err
	}

	return NewGeometric(g.firstPayment, g.growthRate, period, rateInterest)
}

// MustBuild is like Build, but panics if the calculation fails.
func (g GeometricConfig) MustBuild() Geometric {
	series, err := g.Build()
	if err != nil {
		panic(err)
	}

	return series
}

// Present builds the series and returns its present value.
func (g GeometricConfig) Present() (money.Money, error) {
	series, err := g.Build()
	if err != nil {
		return money.Money{}, err
	}

	return series.Present()
}

// MustPresent is like Present, but panics if the calculation fails.
func (g GeometricConfig) MustPresent() money.Money {
	m, err := g.Present()
	if err != nil {
		panic(err)
	}

	return m
}

// Future builds the series and returns its future value.
func (g GeometricConfig) Future() (money.Money, error) {
	series, err := g.Build()
	if err != nil {
		return money.Money{}, err
	}

	return series.Future()
}

// MustFuture is like Future, but panics if the calculation fails.
func (g GeometricConfig) MustFuture() money.Money {
	m, err := g.Future()
	if err != nil {
		panic(err)
	}

	return m
}
