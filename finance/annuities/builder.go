package annuities

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
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
	value            money.Money
	present          money.Money
	future           money.Money
	periods          int
	rate             float64
	frequency        compoundinterest.CompoundingFrequency
	rateType         compoundinterest.TypeRate
	paymentFrequency compoundinterest.CompoundingFrequency
	hasPaymentFreq   bool
	deferPeriods     int
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
		frequency: compoundinterest.Monthly,
		rateType:  compoundinterest.RateEffectyPeriodic,
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
// This is a convenience method that calculates Periods = years × periods per year.
//
// Note: Call this after setting the frequency (Monthly(), Quarterly(), etc.)
// and, for general annuities, after PaymentFrequency(), since the
// conversion uses whichever of the two is configured at the time of the
// call.
//
// Parameters:
//   - n: The number of years (e.g., 30 for a 30-year loan)
//
// Example:
//
//	.NewAnnuity().Monthly().Years(30)  // Automatically calculates 360 periods for monthly payments
func (a AnnuityConfig) Years(n int) AnnuityConfig {
	frequency := a.frequency
	if a.hasPaymentFreq {
		frequency = a.paymentFrequency
	}

	periodsPerYear := 12
	switch frequency {
	case compoundinterest.Daily:
		periodsPerYear = 365
	case compoundinterest.Bimonthly:
		periodsPerYear = 6
	case compoundinterest.Quarterly:
		periodsPerYear = 4
	case compoundinterest.FourMonthly:
		periodsPerYear = 3
	case compoundinterest.SemiAnnually:
		periodsPerYear = 2
	case compoundinterest.Annually:
		periodsPerYear = 1
	}

	a.periods = n * periodsPerYear
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
	case compoundinterest.Daily:
		divisor = 365.0
	case compoundinterest.Bimonthly:
		divisor = 6.0
	case compoundinterest.Quarterly:
		divisor = 4.0
	case compoundinterest.FourMonthly:
		divisor = 3.0
	case compoundinterest.SemiAnnually:
		divisor = 2.0
	case compoundinterest.Annually:
		divisor = 1.0
	}

	a.rate = r / divisor
	a.rateType = compoundinterest.RateEffectyPeriodic

	return a
}

// EffectiveAnnualRate sets the effective annual interest rate and automatically converts
// it to the periodic rate based on the configured compounding frequency.
//
// Use this when you have an effective annual rate (the true yearly cost/return) and need
// to convert it to the periodic rate for calculations.
//
// The conversion formula is: periodic_rate = (1 + annual_rate)^(1/periods) - 1
// For monthly: periodic_rate = (1 + annual_rate)^(1/12) - 1
//
// Parameters:
//   - r: The effective annual interest rate as a decimal (e.g., 0.2668 for 26.68% effective annual)
//
// Example:
//
//	.NewAnnuity().EffectiveAnnualRate(0.2668)  // 26.68% effective annual, converts to ~2% monthly
func (a AnnuityConfig) EffectiveAnnualRate(r float64) AnnuityConfig {
	a.rate = r
	a.rateType = compoundinterest.RateEffectyAnnually

	return a
}

// Monthly sets the payment frequency to monthly (12 payments per year).
// This is the most common frequency for mortgages and loans.
//
// Example:
//
//	.NewAnnuity().Monthly()
func (a AnnuityConfig) Monthly() AnnuityConfig {
	a.frequency = compoundinterest.Monthly

	return a
}

// Annually sets the payment frequency to annually (1 payment per year).
//
// Example:
//
//	.NewAnnuity().Annually()
func (a AnnuityConfig) Annually() AnnuityConfig {
	a.frequency = compoundinterest.Annually

	return a
}

// Quarterly sets the payment frequency to quarterly (4 payments per year).
//
// Example:
//
//	.NewAnnuity().Quarterly()
func (a AnnuityConfig) Quarterly() AnnuityConfig {
	a.frequency = compoundinterest.Quarterly

	return a
}

// PaymentFrequency sets the payment frequency for a general annuity — one
// whose payments are made more or less often than the interest compounds
// (e.g. monthly payments on a quarterly-compounded rate). Set the
// compounding frequency as usual with Monthly(), Quarterly(), etc.; that
// stays the frequency the rate is understood in, while PaymentFrequency
// controls how often Periods()/Years() payments actually occur.
//
// Without a call to PaymentFrequency, the annuity is a "simple" one: it
// pays and compounds on the same frequency.
//
// Example:
//
//	// Monthly payments against a quarterly-compounded 8% nominal rate.
//	annuity := NewAnnuity().
//	    Quarterly().
//	    AnnualRate(0.08).
//	    PaymentFrequency(compoundinterest.Monthly).
//	    Periods(24).
//	    MustBuild()
func (a AnnuityConfig) PaymentFrequency(f compoundinterest.CompoundingFrequency) AnnuityConfig {
	a.paymentFrequency = f
	a.hasPaymentFreq = true

	return a
}

// Defer sets the number of grace periods — periods with no payment — before
// the annuity's regular payments begin, turning it into a deferred annuity.
// Use DeferredPresentValue/DeferredPayment (or their Anticipate/Must
// variants) to get results that account for it; Present(), Future(), and
// Payment() ignore it.
//
// Example:
//
//	.NewAnnuity().Value(500, money.USD).AnnualRate(0.06).Periods(12).Monthly().Defer(3)
func (a AnnuityConfig) Defer(periods int) AnnuityConfig {
	a.deferPeriods = periods
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
func (a AnnuityConfig) Build() (Annuity, error) {
	rate, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(a.rate),
		a.frequency,
		a.rateType,
	)
	if err != nil {
		return Annuity{}, err
	}

	periodFrequency := a.frequency
	if a.hasPaymentFreq {
		periodFrequency = a.paymentFrequency
	}

	period, err := compoundinterest.NewPeriod(
		decimal.MustFromFloat64(float64(a.periods)),
		periodFrequency,
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
func (a AnnuityConfig) MustBuild() Annuity {
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
func (a AnnuityConfig) Payment() (money.Money, error) {
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
func (a AnnuityConfig) MustPayment() money.Money {
	m, err := a.Payment()
	if err != nil {
		panic(err)
	}

	return m
}

// AnticipatePayment is like Payment, but assumes each payment is made at the
// beginning of its period (annuity due) instead of the end (ordinary
// annuity).
//
// Example:
//
//	payment, err := NewAnnuity().
//	    Present(300000, money.USD).
//	    AnnualRate(0.06).
//	    Periods(360).
//	    Monthly().
//	    AnticipatePayment()
func (a AnnuityConfig) AnticipatePayment() (money.Money, error) {
	annuity, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}

	return annuity.AnticipatePaymentFromPresentValue()
}

// MustAnticipatePayment is like AnticipatePayment, but panics if the
// calculation fails.
func (a AnnuityConfig) MustAnticipatePayment() money.Money {
	m, err := a.AnticipatePayment()
	if err != nil {
		panic(err)
	}

	return m
}

// FutureValue builds the Annuity and calculates the future value of a
// recurring investment: the initial principal (Present) compounded over the
// term, plus the future value of equal periodic contributions (Value) made
// at the end of each period (ordinary annuity).
//
// Use this to model an investment plan such as "I have $1,000 today and add
// $100 every month".
//
// Returns:
//   - The total future value as a Money instance
//   - An error if the calculation fails
//
// Example:
//
//	future, err := NewAnnuity().
//	    Present(1000, money.USD).
//	    Value(100, money.USD).
//	    AnnualRate(0.06).
//	    Periods(12).
//	    Monthly().
//	    FutureValue()
func (a AnnuityConfig) FutureValue() (money.Money, error) {
	annuity, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}

	return annuity.FutureWithContributions()
}

// MustFutureValue is like FutureValue, but panics if the calculation fails.
func (a AnnuityConfig) MustFutureValue() money.Money {
	m, err := a.FutureValue()
	if err != nil {
		panic(err)
	}

	return m
}

// AnticipateFutureValue is like FutureValue, but assumes each contribution is
// made at the beginning of its period (annuity due) instead of at the end.
func (a AnnuityConfig) AnticipateFutureValue() (money.Money, error) {
	annuity, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}

	return annuity.AnticipateFutureWithContributions()
}

// MustAnticipateFutureValue is like AnticipateFutureValue, but panics if the
// calculation fails.
func (a AnnuityConfig) MustAnticipateFutureValue() money.Money {
	m, err := a.AnticipateFutureValue()
	if err != nil {
		panic(err)
	}

	return m
}

// DeferredPresentValue builds the Annuity and calculates the present value
// of a deferred ordinary annuity, using the grace period configured with
// Defer().
//
// Example:
//
//	present, err := NewAnnuity().
//	    Value(500, money.USD).
//	    AnnualRate(0.06).
//	    Periods(12).
//	    Monthly().
//	    Defer(3).
//	    DeferredPresentValue()
func (a AnnuityConfig) DeferredPresentValue() (money.Money, error) {
	annuity, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}

	return annuity.PresentDeferred(a.deferPeriods)
}

// MustDeferredPresentValue is like DeferredPresentValue, but panics if the
// calculation fails.
func (a AnnuityConfig) MustDeferredPresentValue() money.Money {
	m, err := a.DeferredPresentValue()
	if err != nil {
		panic(err)
	}

	return m
}

// AnticipateDeferredPresentValue is like DeferredPresentValue, but assumes
// each payment, once the grace period ends, is made at the beginning of its
// period (annuity due) instead of the end.
func (a AnnuityConfig) AnticipateDeferredPresentValue() (money.Money, error) {
	annuity, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}

	return annuity.AnticipatePresentDeferred(a.deferPeriods)
}

// MustAnticipateDeferredPresentValue is like AnticipateDeferredPresentValue,
// but panics if the calculation fails.
func (a AnnuityConfig) MustAnticipateDeferredPresentValue() money.Money {
	m, err := a.AnticipateDeferredPresentValue()
	if err != nil {
		panic(err)
	}

	return m
}

// DeferredPayment builds the Annuity and calculates the periodic payment
// for a deferred ordinary annuity that pays off the configured Present
// value, using the grace period configured with Defer().
//
// Example:
//
//	payment, err := NewAnnuity().
//	    Present(300000, money.USD).
//	    AnnualRate(0.06).
//	    Periods(360).
//	    Monthly().
//	    Defer(6).
//	    DeferredPayment()
func (a AnnuityConfig) DeferredPayment() (money.Money, error) {
	annuity, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}

	return annuity.PaymentFromPresentValueDeferred(a.deferPeriods)
}

// MustDeferredPayment is like DeferredPayment, but panics if the
// calculation fails.
func (a AnnuityConfig) MustDeferredPayment() money.Money {
	m, err := a.DeferredPayment()
	if err != nil {
		panic(err)
	}

	return m
}
