package annuities

import (
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

// AnnuityConfig is a builder for creating Annuity instances using a fluent API.
// It allows you to chain method calls to configure annuity parameters before building.
//
// Annuities are a series of equal payments made at regular intervals.
// This builder simplifies the creation of annuity calculations for loans, mortgages,
// retirement plans, and other financial instruments with regular payments.
//
// Use NewAnnuity() to create a new builder instance, configure the parameters,
// and then call Build(), MustBuild(), Payment(), or MustPayment().
//
// Example:
//
//	payment := NewAnnuity().
//	    Present(300000, money.USD).
//	    AnnualRate(0.06).
//	    Periods(360).
//	    Monthly().
//	    MustPayment()
type AnnuityConfig struct {
	value     money.Money
	present   money.Money
	future    money.Money
	periods   int
	rate      float64
	frequency compositeinterest.CompoundingFrequency
	rateType  compositeinterest.TypeRate
}

// NewAnnuity creates a new AnnuityConfig builder instance with default values.
// Default frequency is Monthly and default rate type is RateEffectyPeriodic.
//
// After creating the builder, chain the necessary methods to configure:
//   - Present or Future value (the loan amount or savings goal)
//   - Rate using Rate() for periodic rate or AnnualRate() for annual rate
//   - Periods using Periods() or Years()
//   - Frequency using Monthly(), Quarterly(), Annually(), or Daily()
//
// Example:
//
//	annuity := NewAnnuity().
//	    Present(200000, money.USD).
//	    AnnualRate(0.05).
//	    Periods(240).
//	    Monthly().
func NewAnnuity() AnnuityConfig {
	return AnnuityConfig{
		frequency: compositeinterest.Monthly,
		rateType:  compositeinterest.RateEffectyPeriodic,
	}
}

// Value sets the periodic payment amount for the annuity.
// This is the fixed payment amount made at each interval.
//
// Use this when calculating the future or present value of a known payment amount.
//
// Parameters:
//   - amount: The payment amount as a float64 (e.g., 500.00 for $500 per period)
//   - currency: The currency code (e.g., money.USD, money.EUR)
//
// Example:
//
//	.NewAnnuity().Value(500, money.USD)
func (a AnnuityConfig) Value(amount float64, currency money.Currency) AnnuityConfig {
	a.value = money.MustMoneyFromFloat64(amount, currency)
	return a
}

// Present sets the present value (current loan amount or initial investment).
// For a loan, this is the amount borrowed.
// For a savings plan, this is the current savings amount.
//
// Parameters:
//   - amount: The monetary amount as a float64 (e.g., 300000.00 for a $300,000 loan)
//   - currency: The currency code (e.g., money.USD, money.EUR)
//
// Example:
//
//	.NewAnnuity().Present(300000, money.USD)  // $300,000 mortgage
func (a AnnuityConfig) Present(amount float64, currency money.Currency) AnnuityConfig {
	a.present = money.MustMoneyFromFloat64(amount, currency)
	return a
}

// Future sets the future value (target amount to accumulate).
// Use this when calculating payments needed to reach a savings goal.
//
// Parameters:
//   - amount: The monetary amount as a float64 (e.g., 1000000.00 for $1,000,000)
//   - currency: The currency code (e.g., money.USD, money.EUR)
//
// Example:
//
//	.NewAnnuity().Future(1000000, money.USD)  // Save $1,000,000
func (a AnnuityConfig) Future(amount float64, currency money.Currency) AnnuityConfig {
	a.future = money.MustMoneyFromFloat64(amount, currency)
	return a
}

// Periods sets the total number of payment periods for the annuity.
//
// Parameters:
//   - n: The number of periods (e.g., 360 for a 30-year mortgage with monthly payments)
//
// Example:
//
//	.NewAnnuity().Periods(360)  // 360 monthly payments (30 years)
func (a AnnuityConfig) Periods(n int) AnnuityConfig {
	a.periods = n
	return a
}

// Years sets the total duration in years and converts it to periods based on the frequency.
// This is a convenience method that calculates Periods = years × frequency.
//
// Parameters:
//   - n: The number of years (e.g., 30 for a 30-year loan)
//
// Example:
//
//	.NewAnnuity().Years(30)  // Automatically calculates 360 periods for monthly payments
func (a AnnuityConfig) Years(n int) AnnuityConfig {
	a.periods = n * 12
	return a
}

// Rate sets the periodic interest rate (rate for each compounding period).
// The rate should be expressed as a decimal (e.g., 0.005 for 0.5% per month).
//
// Note: Use this when you already have the periodic rate.
// For annual rates, use AnnualRate() which automatically converts to periodic rate.
//
// Parameters:
//   - r: The periodic interest rate as a decimal (e.g., 0.005 for 0.5%, 0.004167 for 0.4167%)
//
// Example:
//
//	.NewAnnuity().Rate(0.005)  // 0.5% monthly rate
func (a AnnuityConfig) Rate(r float64) AnnuityConfig {
	a.rate = r
	return a
}

// AnnualRate sets the annual interest rate and automatically converts it to the periodic rate
// based on the configured compounding frequency.
//
// This is the recommended method for setting interest rates when you have an annual rate
// (which is the most common way rates are quoted).
//
// The conversion formula depends on the frequency:
//   - Monthly: annual_rate / 12
//   - Quarterly: annual_rate / 4
//   - Daily: annual_rate / 365
//   - etc.
//
// Parameters:
//   - r: The annual interest rate as a decimal (e.g., 0.06 for 6% annual rate)
//
// Example:
//
//	.NewAnnuity().AnnualRate(0.06)  // 6% annual rate, converts to ~0.5% monthly
func (a AnnuityConfig) AnnualRate(r float64) AnnuityConfig {
	divisor := 12.0
	switch a.frequency {
	case compositeinterest.Daily:
		divisor = 365.0
	case compositeinterest.Bimonthly:
		divisor = 6.0
	case compositeinterest.QuarterlyOne, compositeinterest.QuarterlyTwo:
		divisor = 4.0
	case compositeinterest.SemiAnnually:
		divisor = 2.0
	case compositeinterest.Annually:
		divisor = 1.0
	}
	a.rate = r / divisor
	return a
}

// Monthly sets the payment frequency to monthly (12 payments per year).
// This is the most common frequency for mortgages and loans.
//
// Example:
//
//	.NewAnnuity().Monthly()
func (a AnnuityConfig) Monthly() AnnuityConfig {
	a.frequency = compositeinterest.Monthly

	return a
}

// Annually sets the payment frequency to annually (1 payment per year).
//
// Example:
//
//	.NewAnnuity().Annually()
func (a AnnuityConfig) Annually() AnnuityConfig {
	a.frequency = compositeinterest.Annually

	return a
}

// Quarterly sets the payment frequency to quarterly (4 payments per year).
//
// Example:
//
//	.NewAnnuity().Quarterly()
func (a AnnuityConfig) Quarterly() AnnuityConfig {
	a.frequency = compositeinterest.QuarterlyOne

	return a
}

// Build creates and returns an Annuity instance based on the configured parameters.
// It validates the parameters and returns an error if they are invalid.
//
// Returns:
//   - An Annuity instance with all configured parameters
//   - An error if validation fails (invalid frequency, negative values, etc.)
//
// Example:
//
//	annuity, err := NewAnnuity().
//	    Present(300000, money.USD).
//	    AnnualRate(0.06).
//	    Periods(360).
//	    Monthly().
//	    Build()
func (a *AnnuityConfig) Build() (Annuity, error) {
	rate, err := compositeinterest.NewRateInterest(
		money.MustFromFloat64(a.rate),
		a.frequency,
		a.rateType,
	)
	if err != nil {
		return Annuity{}, err
	}

	period, err := compositeinterest.NewPeriod(
		money.MustFromFloat64(float64(a.periods)),
		a.frequency,
	)
	if err != nil {
		return Annuity{}, err
	}

	return New(a.value, a.present, a.future, period, rate)
}

// MustBuild creates and returns an Annuity instance based on the configured parameters.
// Unlike Build(), this method panics if validation fails.
//
// Use this for quick prototyping or when you are certain the parameters are valid.
//
// Returns:
//   - An Annuity instance with all configured parameters
//
// Panics:
//   - If validation fails (invalid frequency, negative values, etc.)
//
// Example:
//
//	annuity := NewAnnuity().
//	    Present(300000, money.USD).
//	    AnnualRate(0.06).
//	    Periods(360).
//	    Monthly().
//	    MustBuild()
func (a *AnnuityConfig) MustBuild() Annuity {
	annuity, err := a.Build()
	if err != nil {
		panic(err)
	}

	return annuity
}

// Payment calculates the periodic payment amount for an annuity based on the present value.
// This is the fixed payment required to fully amortize the loan over the specified periods.
//
// This method should be used when you know the loan amount (present value) and want to
// calculate the required payment.
//
// Returns:
//   - The periodic payment amount as a Money instance
//   - An error if the calculation fails
//
// Example:
//
//	payment, err := NewAnnuity().
//	    Present(300000, money.USD).
//	    AnnualRate(0.06).
//	    Periods(360).
//	    Monthly().
//	    Payment()
func (a *AnnuityConfig) Payment() (money.Money, error) {
	annuity, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}

	return annuity.PaymentFromPresentValue()
}

// MustPayment calculates the periodic payment amount for an annuity based on the present value.
// Unlike Payment(), this method panics if the calculation fails.
//
// Use this for quick prototyping or when you are certain the parameters are valid.
//
// Returns:
//   - The periodic payment amount as a Money instance
//
// Panics:
//   - If the calculation fails
//
// Example:
//
//	payment := NewAnnuity().
//	    Present(300000, money.USD).
//	    AnnualRate(0.06).
//	    Periods(360).
//	    Monthly().
//	    MustPayment()
func (a *AnnuityConfig) MustPayment() money.Money {
	m, err := a.Payment()
	if err != nil {
		panic(err)
	}

	return m
}
