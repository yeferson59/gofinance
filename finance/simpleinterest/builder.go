package simpleinterest

import (
	"github.com/yeferson59/gofinance/money"
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
	rate       money.Decimal
	periods    int
	periodType Periods
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
	s.present = money.MustMoneyFromFloat64(amount, currency)
	return s
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
	s.future = money.MustMoneyFromFloat64(amount, currency)
	return s
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
	s.interest = money.MustMoneyFromFloat64(amount, currency)
	return s
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
	s.rate = money.MustFromFloat64(r)
	return s
}

// AnnualRate sets the annual interest rate and automatically converts it to the periodic rate
// based on the configured period type.
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
	divisor := 12.0
	switch s.periodType {
	case Days:
		divisor = 365.0
	case Weeks:
		divisor = 52.0
	case Years:
		divisor = 1.0
	}
	s.rate = money.MustFromFloat64(r / divisor)
	return s
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
	period := NewPeriod(money.MustFromFloat64(float64(s.periods)), s.periodType)
	return New(s.future, s.present, s.interest, s.rate, period)
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
