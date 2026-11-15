package compoundinterest

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/internal/guard"
	"github.com/yeferson59/gofinance/v2/money"
)

// CompoundConfig is a builder for creating CompoundInterest instances using a fluent API.
// It allows you to chain method calls to set the required parameters before building.
//
// Use NewCompound() to create a new builder instance, then chain the necessary methods,
// and finally call Build() or MustBuild() to create the CompoundInterest.
//
// Example:
//
//	ci := NewCompound().
//	    Present(1000, money.USD).
//	    Rate(0.05).
//	    Periods(12).
//	    Monthly().
//	    MustBuild()
type CompoundConfig struct {
	present   money.Money
	future    money.Money
	rate      decimal.Decimal
	periods   int
	frequency CompoundingFrequency
	rateType  TypeRate

	// err holds the first invalid input a setter received (a NaN or infinite
	// float, say). Setters can't return an error without breaking the chain,
	// so it is kept here and reported by Build.
	err error
}

// withErr records err as the builder's error unless one is already held, so
// the first problem is the one reported. A nil err leaves c unchanged.
func (c CompoundConfig) withErr(err error) CompoundConfig {
	if c.err == nil {
		c.err = err
	}

	return c
}

// NewCompound creates a new CompoundConfig builder instance.
// The default rate type is RateEffectivePeriodic; use RateType() to change it.
// You must set at least the present or future value, rate, periods, and frequency
// before calling Build() or MustBuild().
func NewCompound() CompoundConfig {
	return CompoundConfig{
		rateType: RateEffectivePeriodic,
	}
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
//	.NewCompound().Present(1000, money.USD)
func (c CompoundConfig) Present(amount float64, currency money.Currency) CompoundConfig {
	m, err := money.NewMoneyFromFloat64(amount, currency)

	return c.withErr(err).PresentMoney(m)
}

// PresentMoney sets the present value using an existing Money instance.
// Use this when you already have a Money object.
//
// Parameters:
//   - m: An existing Money instance representing the present value
func (c CompoundConfig) PresentMoney(m money.Money) CompoundConfig {
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
//	.NewCompound().Future(1500, money.USD)
func (c CompoundConfig) Future(amount float64, currency money.Currency) CompoundConfig {
	m, err := money.NewMoneyFromFloat64(amount, currency)

	return c.withErr(err).FutureMoney(m)
}

// FutureMoney sets the future value using an existing Money instance.
// Use this when you already have a Money object.
//
// Parameters:
//   - m: An existing Money instance representing the future value
func (c CompoundConfig) FutureMoney(m money.Money) CompoundConfig {
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
//	.NewCompound().Rate(0.05)  // 5% interest rate
func (c CompoundConfig) Rate(rate float64) CompoundConfig {
	d, err := decimal.NewFromFloat64(rate)

	return c.withErr(err).RateDecimal(d)
}

// RateDecimal is like Rate but takes the rate as a decimal.Decimal, so a rate
// a float64 can't hold exactly keeps every digit.
func (c CompoundConfig) RateDecimal(r decimal.Decimal) CompoundConfig {
	c.rate = r
	return c
}

// RateMoney sets the interest rate from a Decimal.
//
// Deprecated: the name is misleading — a rate is not an amount of money. Use
// RateDecimal, which does the same.
func (c CompoundConfig) RateMoney(r decimal.Decimal) CompoundConfig {
	return c.RateDecimal(r)
}

// Periods sets the number of compounding periods for the calculation.
// The interpretation of this number depends on the frequency (e.g., 12 months if Monthly).
//
// Parameters:
//   - n: The number of periods (e.g., 12 for one year with Monthly frequency)
//
// Example:
//
//	.NewCompound().Periods(12)  // 12 months, years, etc. based on frequency
func (c CompoundConfig) Periods(n int) CompoundConfig {
	c.periods = n
	return c
}

// Frequency sets the compounding frequency using a CompoundingFrequency value.
//
// Parameters:
//   - f: The compounding frequency (Daily, Monthly, Bimonthly, Quarterly, FourMonthly, SemiAnnually, Annually)
//
// Example:
//
//	.NewCompound().Frequency(Monthly)
func (c CompoundConfig) Frequency(f CompoundingFrequency) CompoundConfig {
	c.frequency = f
	return c
}

// RateType sets the type of interest rate being used.
//
// Parameters:
//   - t: The rate type (RateEffectivePeriodic, RateEffectiveNominal, RateEffectiveAnnually, RateAnticipateEffectivePeriodic, etc.)
//
// Example:
//
//	.NewCompound().RateType(RateEffectivePeriodic)
func (c CompoundConfig) RateType(t TypeRate) CompoundConfig {
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
//	.NewCompound().Monthly()
func (c CompoundConfig) Monthly() CompoundConfig {
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
//	.NewCompound().Annually()
func (c CompoundConfig) Annually() CompoundConfig {
	c.frequency = Annually
	return c
}

// Quarterly sets the compounding frequency to quarterly (4 periods per year).
// This is a convenience method equivalent to calling Frequency(Quarterly).
//
// Use this for quarterly compounding calculations.
//
// Example:
//
//	.NewCompound().Quarterly()
func (c CompoundConfig) Quarterly() CompoundConfig {
	c.frequency = Quarterly
	return c
}

// Daily sets the compounding frequency to daily (365 periods per year).
// This is a convenience method equivalent to calling Frequency(Daily).
//
// Use this for daily compounding calculations.
//
// Example:
//
//	.NewCompound().Daily()
func (c CompoundConfig) Daily() CompoundConfig {
	c.frequency = Daily
	return c
}

// Build creates and returns a CompoundInterest instance based on the configured parameters.
// It validates the parameters and returns an error if they are invalid.
//
// Returns:
//   - A CompoundInterest instance with all configured parameters
//   - An error if validation fails (invalid frequency, negative values, etc.)
//
// Example:
//
//	ci, err := NewCompound().
//	    Present(1000, money.USD).
//	    Rate(0.05).
//	    Periods(12).
//	    Monthly().
//	    Build()
func (c CompoundConfig) Build() (_ CompoundInterest, err error) {
	defer guard.Recover(&err)

	if c.err != nil {
		return CompoundInterest{}, c.err
	}

	periods, err := decimal.NewFromInt64(int64(c.periods), 0)
	if err != nil {
		return CompoundInterest{}, err
	}

	period, err := NewPeriod(periods, c.frequency)
	if err != nil {
		return CompoundInterest{}, err
	}

	rate, err := NewRateInterest(c.rate, c.frequency, c.rateType)
	if err != nil {
		return CompoundInterest{}, err
	}

	return New(c.present, c.future, rate, period)
}

// MustBuild creates and returns a CompoundInterest instance based on the configured parameters.
// Unlike Build(), this method panics if validation fails.
//
// Use this for quick prototyping or when you are certain the parameters are valid.
//
// Returns:
//   - A CompoundInterest instance with all configured parameters
//
// Panics:
//   - If validation fails (invalid frequency, negative values, etc.)
//
// Example:
//
//	ci := NewCompound().
//	    Present(1000, money.USD).
//	    Rate(0.05).
//	    Periods(12).
//	    Monthly().
//	    MustBuild()
func (c CompoundConfig) MustBuild() CompoundInterest {
	ci, err := c.Build()
	if err != nil {
		panic(err)
	}
	return ci
}
