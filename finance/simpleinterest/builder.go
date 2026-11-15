package simpleinterest

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// SimpleConfig is a builder for creating SimpleInterest instances using a fluent API.
// It allows you to chain method calls to configure simple interest parameters before building.
//
// Simple interest is calculated only on the original principal amount.
// Unlike compound interest, it does not include interest on accumulated interest.
//
// Use NewSimple() to create a new builder instance, configure the parameters,
// and then call Build(), FutureValue(), or PresentValue().
//
// Example:
//
//	si := NewSimple().
//	    Present(5000, money.USD).
//	    AnnualRate(0.12).
//	    Periods(18).
//	    Months().
//	    Build()
type SimpleConfig struct {
	present    money.Money
	future     money.Money
	interest   money.Money
	rate       decimal.Decimal
	annual     bool
	periods    int
	periodType Periods

	// err holds the first invalid input a setter received (a NaN or infinite
	// float, say). Setters can't return an error without breaking the chain,
	// and Build has no error to return, so it travels into the built
	// SimpleInterest, whose calculations report it.
	err error
}

// withErr records err as the builder's error unless one is already held, so
// the first problem is the one reported. A nil err leaves s unchanged.
func (s SimpleConfig) withErr(err error) SimpleConfig {
	if s.err == nil {
		s.err = err
	}

	return s
}

// periodsPerYear returns how many periods of the configured type fit in a
// year, the divisor that turns an annual rate into a periodic one.
func (s SimpleConfig) periodsPerYear() int64 {
	switch s.periodType {
	case Days:
		return 365
	case Weeks:
		return 52
	case Years:
		return 1
	default:
		return 12
	}
}

// periodicRate returns the rate per period, converting from the annual rate
// when one was set. The conversion happens here rather than in AnnualRate so
// it uses the period type the builder ends up with, whatever order the methods
// were called in.
func (s SimpleConfig) periodicRate() (decimal.Decimal, error) {
	if !s.annual {
		return s.rate, nil
	}

	periodsPerYear, err := decimal.NewFromInt64(s.periodsPerYear(), 0)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return s.rate.Div(periodsPerYear)
}

// NewSimple creates a new SimpleConfig builder instance with default values.
// Default period type is Months.
//
// After creating the builder, chain the necessary methods to configure:
//   - Present value (the principal amount)
//   - Rate using Rate() for periodic rate or AnnualRate() for annual rate
//   - Periods using Periods() with period type (Months, Years, Days, Weeks)
//
// Example:
//
//	si := NewSimple().
//	    Present(5000, money.USD).
//	    AnnualRate(0.12).
//	    Periods(18).
//	    Months()
func NewSimple() SimpleConfig {
	return SimpleConfig{
		periodType: Months,
	}
}

// Present sets the present value (principal/initial amount) for the simple interest calculation.
// This is the original amount of money before interest is applied.
//
// Parameters:
//   - amount: The monetary amount as a float64 (e.g., 5000.00 for $5000)
//   - currency: The currency code (e.g., money.USD, money.EUR)
//
// Example:
//
//	.NewSimple().Present(5000, money.USD)  // $5,000 principal
func (s SimpleConfig) Present(amount float64, currency money.Currency) SimpleConfig {
	m, err := money.NewMoneyFromFloat64(amount, currency)
	s.present = m

	return s.withErr(err)
}

// Future sets the future value (target amount) for the simple interest calculation.
// Use this when you know the target amount and want to calculate the required present value.
//
// Parameters:
//   - amount: The monetary amount as a float64 (e.g., 6000.00 for $6000)
//   - currency: The currency code (e.g., money.USD, money.EUR)
//
// Example:
//
//	.NewSimple().Future(6000, money.USD)  // $6,000 target amount
func (s SimpleConfig) Future(amount float64, currency money.Currency) SimpleConfig {
	m, err := money.NewMoneyFromFloat64(amount, currency)
	s.future = m

	return s.withErr(err)
}

// Interest sets the interest amount for the simple interest calculation.
// Use this when you know the interest amount and want to calculate present or future value.
//
// Parameters:
//   - amount: The monetary amount as a float64 (e.g., 900.00 for $900 interest)
//   - currency: The currency code (e.g., money.USD, money.EUR)
//
// Example:
//
//	.NewSimple().Interest(900, money.USD)  // $900 interest
func (s SimpleConfig) Interest(amount float64, currency money.Currency) SimpleConfig {
	m, err := money.NewMoneyFromFloat64(amount, currency)
	s.interest = m

	return s.withErr(err)
}

// Rate sets the periodic interest rate for the simple interest calculation.
// The rate should be expressed as a decimal (e.g., 0.01 for 1%).
//
// Note: Use this when you already have the periodic rate.
// For annual rates, use AnnualRate() which automatically converts to periodic rate.
//
// Parameters:
//   - r: The periodic interest rate as a decimal (e.g., 0.01 for 1%, 0.001667 for 0.1667% monthly)
//
// Example:
//
//	.NewSimple().Rate(0.01)  // 1% monthly rate
func (s SimpleConfig) Rate(r float64) SimpleConfig {
	rate, err := decimal.NewFromFloat64(r)
	s.rate = rate
	s.annual = false

	return s.withErr(err)
}

// AnnualRate sets the annual interest rate, which is divided by the number of
// periods of the configured type that fit in a year to obtain the periodic
// rate, whatever order the builder methods are called in. It replaces any rate
// set with Rate.
//
// This is the recommended method for setting interest rates when you have an annual rate.
//
// The conversion formula depends on the period type:
//   - Months: annual_rate / 12
//   - Days: annual_rate / 365
//   - Weeks: annual_rate / 52
//   - Years: annual_rate / 1
//
// Parameters:
//   - r: The annual interest rate as a decimal (e.g., 0.12 for 12% annual rate)
//
// Example:
//
//	.NewSimple().AnnualRate(0.12)  // 12% annual rate, converts to 1% monthly
func (s SimpleConfig) AnnualRate(r float64) SimpleConfig {
	rate, err := decimal.NewFromFloat64(r)
	s.rate = rate
	s.annual = true

	return s.withErr(err)
}

// Periods sets the number of time periods for the simple interest calculation.
//
// Parameters:
//   - n: The number of periods (e.g., 18 for 18 months with Months period type)
//
// Example:
//
//	.NewSimple().Periods(18)  // 18 periods (months, days, etc. based on period type)
func (s SimpleConfig) Periods(n int) SimpleConfig {
	s.periods = n
	return s
}

// PeriodType sets the type of time period using a Periods value.
//
// Parameters:
//   - p: The period type (Days, Weeks, Months, Years)
//
// Example:
//
//	.NewSimple().PeriodType(Months)
func (s SimpleConfig) PeriodType(p Periods) SimpleConfig {
	s.periodType = p
	return s
}

// Months sets the period type to months.
// This is a convenience method equivalent to calling PeriodType(Months).
//
// Example:
//
//	.NewSimple().Months()
func (s SimpleConfig) Months() SimpleConfig {
	s.periodType = Months
	return s
}

// Years sets the period type to years.
// This is a convenience method equivalent to calling PeriodType(Years).
//
// Example:
//
//	.NewSimple().Years()
func (s SimpleConfig) Years() SimpleConfig {
	s.periodType = Years
	return s
}

// Days sets the period type to days.
// This is a convenience method equivalent to calling PeriodType(Days).
//
// Example:
//
//	.NewSimple().Days()
func (s SimpleConfig) Days() SimpleConfig {
	s.periodType = Days
	return s
}

// Weeks sets the period type to weeks.
// This is a convenience method equivalent to calling PeriodType(Weeks).
//
// Example:
//
//	.NewSimple().Weeks()
func (s SimpleConfig) Weeks() SimpleConfig {
	s.periodType = Weeks
	return s
}

// Build creates and returns a SimpleInterest instance based on the configured parameters.
//
// Build has no error to return. If a setter received an invalid input (NaN or
// ±Inf), or the annual rate can't be converted to a periodic one, the error is
// carried by the returned SimpleInterest and every calculation on it returns
// that error.
//
// Returns:
//   - A SimpleInterest instance with all configured parameters
//
// Example:
//
//	si := NewSimple().
//	    Present(5000, money.USD).
//	    AnnualRate(0.12).
//	    Periods(18).
//	    Months().
//	    Build()
//
//	future, _ := si.Future()
func (s SimpleConfig) Build() SimpleInterest {
	periods, err := decimal.NewFromInt64(int64(s.periods), 0)
	s = s.withErr(err)

	rate, err := s.periodicRate()
	s = s.withErr(err)

	si := New(s.future, s.present, s.interest, rate, NewPeriod(periods, s.periodType))
	si.err = s.err

	return si
}

// FutureValue calculates the future value based on the present value, rate, and periods.
// This is the total amount (principal + interest) after the specified time.
//
// Formula: Future = Present × (1 + Rate × Periods)
//
// Returns:
//   - The future value as a Money instance
//   - An error if the calculation fails
//
// Example:
//
//	future, err := NewSimple().
//	    Present(5000, money.USD).
//	    AnnualRate(0.12).
//	    Periods(18).
//	    Months().
//	    FutureValue()
func (s SimpleConfig) FutureValue() (money.Money, error) {
	si := s.Build()
	return si.FutureWithRateInterest()
}

// PresentValue calculates the present value based on the future value, rate, and periods.
// This is the amount needed today to reach the configured future value.
//
// Formula: Present = Future / (1 + Rate × Periods)
//
// Returns:
//   - The present value as a Money instance
//   - An error if the calculation fails
//
// Example:
//
//	present, err := NewSimple().
//	    Future(5900, money.USD).
//	    AnnualRate(0.12).
//	    Periods(18).
//	    Months().
//	    PresentValue()
func (s SimpleConfig) PresentValue() (money.Money, error) {
	si := s.Build()
	return si.PresentWithFuture()
}
