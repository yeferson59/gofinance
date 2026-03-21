package compositeinterest

import (
	"github.com/yeferson59/gofinance/money"
)

// CompositeConfig is a builder for creating CompositeInterest instances using a fluent API.
// It allows you to chain method calls to set the required parameters before building.
//
// Use NewComposite() to create a new builder instance, then chain the necessary methods,
// and finally call Build() or MustBuild() to create the CompositeInterest.
//
// Example:
//
//	ci := NewComposite().
//	    Present(1000, money.USD).
//	    Rate(0.05).
//	    Periods(12).
//	    Monthly().
//	    MustBuild()
type CompositeConfig struct {
	present   money.Money
	future    money.Money
	rate      money.Decimal
	periods   int
	frequency CompoundingFrequency
	rateType  TypeRate
}

// NewComposite creates a new CompositeConfig builder instance.
// The builder starts with zero values for all parameters.
// You must set at least the present or future value, rate, periods, and frequency
// before calling Build() or MustBuild().
func NewComposite() CompositeConfig {
	return CompositeConfig{}
}

// Present sets the present value (initial capital/principal) for the compound interest calculation.
// This is the amount of money you have now, which will grow over time.
//
// Parameters:
//   - amount: The monetary amount as a float64 (e.g., 1000.00 for $1000)
//   - currency: The currency code (e.g., money.USD, money.EUR)
//
// Example:
//
//	.NewComposite().Present(1000, money.USD)
func (c CompositeConfig) Present(amount float64, currency money.Currency) CompositeConfig {
	c.present = money.MustMoneyFromFloat64(amount, currency)
	return c
}

// PresentMoney sets the present value using an existing Money instance.
// Use this when you already have a Money object.
//
// Parameters:
//   - m: An existing Money instance representing the present value
func (c CompositeConfig) PresentMoney(m money.Money) CompositeConfig {
	c.present = m
	return c
}

// Future sets the future value (target amount) for the compound interest calculation.
// This is the amount of money you want to have in the future.
//
// Parameters:
//   - amount: The monetary amount as a float64 (e.g., 1500.00 for $1500)
//   - currency: The currency code (e.g., money.USD, money.EUR)
//
// Example:
//
//	.NewComposite().Future(1500, money.USD)
func (c CompositeConfig) Future(amount float64, currency money.Currency) CompositeConfig {
	c.future = money.MustMoneyFromFloat64(amount, currency)
	return c
}

// FutureMoney sets the future value using an existing Money instance.
// Use this when you already have a Money object.
//
// Parameters:
//   - m: An existing Money instance representing the future value
func (c CompositeConfig) FutureMoney(m money.Money) CompositeConfig {
	c.future = m
	return c
}

// Rate sets the interest rate for the compound interest calculation.
// The rate should be expressed as a decimal (e.g., 0.05 for 5%).
//
// Note: This sets the rate in the compounding frequency you specify.
// For annual rates that need conversion, use RateType to set the appropriate rate type.
//
// Parameters:
//   - rate: The interest rate as a decimal (e.g., 0.05 for 5%, 0.10 for 10%)
//
// Example:
//
//	.NewComposite().Rate(0.05)  // 5% interest rate
func (c CompositeConfig) Rate(rate float64) CompositeConfig {
	c.rate = money.MustFromFloat64(rate)
	return c
}

// RateMoney sets the interest rate using an existing Decimal instance.
// Use this when you already have a Decimal value.
//
// Parameters:
//   - r: An existing Decimal instance representing the interest rate
func (c CompositeConfig) RateMoney(r money.Decimal) CompositeConfig {
	c.rate = r
	return c
}

// Periods sets the number of compounding periods for the calculation.
// The interpretation of this number depends on the frequency (e.g., 12 months if Monthly).
//
// Parameters:
//   - n: The number of periods (e.g., 12 for one year with Monthly frequency)
//
// Example:
//
//	.NewComposite().Periods(12)  // 12 months, years, etc. based on frequency
func (c CompositeConfig) Periods(n int) CompositeConfig {
	c.periods = n
	return c
}

// Frequency sets the compounding frequency using a CompoundingFrequency value.
//
// Parameters:
//   - f: The compounding frequency (Daily, Monthly, Bimonthly, QuarterlyOne, QuarterlyTwo, SemiAnnually, Annually)
//
// Example:
//
//	.NewComposite().Frequency(Monthly)
func (c CompositeConfig) Frequency(f CompoundingFrequency) CompositeConfig {
	c.frequency = f
	return c
}

// RateType sets the type of interest rate being used.
//
// Parameters:
//   - t: The rate type (RateEffectyPeriodic, RateEffectyNominal, RateEffectyAnnually, RateAnticipateEffectyPeriodic, etc.)
//
// Example:
//
//	.NewComposite().RateType(RateEffectyPeriodic)
func (c CompositeConfig) RateType(t TypeRate) CompositeConfig {
	c.rateType = t
	return c
}

// Monthly sets the compounding frequency to monthly (12 periods per year).
// This is a convenience method equivalent to calling Frequency(Monthly).
//
// Use this when calculating investments or loans with monthly compounding.
//
// Example:
//
//	.NewComposite().Monthly()
func (c CompositeConfig) Monthly() CompositeConfig {
	c.frequency = Monthly
	return c
}

// Annually sets the compounding frequency to annually (1 period per year).
// This is a convenience method equivalent to calling Frequency(Annually).
//
// Use this for annual compounding calculations.
//
// Example:
//
//	.NewComposite().Annually()
func (c CompositeConfig) Annually() CompositeConfig {
	c.frequency = Annually
	return c
}

// Quarterly sets the compounding frequency to quarterly (4 periods per year).
// This is a convenience method equivalent to calling Frequency(QuarterlyOne).
//
// Use this for quarterly compounding calculations.
//
// Example:
//
//	.NewComposite().Quarterly()
func (c CompositeConfig) Quarterly() CompositeConfig {
	c.frequency = QuarterlyOne
	return c
}

// Daily sets the compounding frequency to daily (365 periods per year).
// This is a convenience method equivalent to calling Frequency(Daily).
//
// Use this for daily compounding calculations.
//
// Example:
//
//	.NewComposite().Daily()
func (c CompositeConfig) Daily() CompositeConfig {
	c.frequency = Daily
	return c
}

// Build creates and returns a CompositeInterest instance based on the configured parameters.
// It validates the parameters and returns an error if they are invalid.
//
// Returns:
//   - A CompositeInterest instance with all configured parameters
//   - An error if validation fails (invalid frequency, negative values, etc.)
//
// Example:
//
//	ci, err := NewComposite().
//	    Present(1000, money.USD).
//	    Rate(0.05).
//	    Periods(12).
//	    Monthly().
//	    Build()
func (c *CompositeConfig) Build() (CompositeInterest, error) {
	period, err := NewPeriod(money.MustFromFloat64(float64(c.periods)), c.frequency)
	if err != nil {
		return CompositeInterest{}, err
	}

	rate, err := NewRateInterest(c.rate, c.frequency, c.rateType)
	if err != nil {
		return CompositeInterest{}, err
	}

	return New(c.present, c.future, rate, period)
}

// MustBuild creates and returns a CompositeInterest instance based on the configured parameters.
// Unlike Build(), this method panics if validation fails.
//
// Use this for quick prototyping or when you are certain the parameters are valid.
//
// Returns:
//   - A CompositeInterest instance with all configured parameters
//
// Panics:
//   - If validation fails (invalid frequency, negative values, etc.)
//
// Example:
//
//	ci := NewComposite().
//	    Present(1000, money.USD).
//	    Rate(0.05).
//	    Periods(12).
//	    Monthly().
//	    MustBuild()
func (c *CompositeConfig) MustBuild() CompositeInterest {
	ci, err := c.Build()
	if err != nil {
		panic(err)
	}
	return ci
}
