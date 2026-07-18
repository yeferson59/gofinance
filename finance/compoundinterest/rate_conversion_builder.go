package compoundinterest

import "github.com/yeferson59/gofinance/decimal"

// RateConversionConfig is a builder for converting interest rates between
// types (periodic, nominal, effective annual, anticipated) and compounding
// frequencies using a fluent API.
//
// Use NewRateConversion() to create a new builder instance, configure the
// source rate, its type, and its compounding frequency, and then call one of
// the terminal conversion methods (ToPeriodic, ToNominal, ToEffectiveAnnual,
// ToAnticipatedPeriodic, ToAnticipatedNominal, ToPeriodicAt, ToNominalAt) or
// Build() to obtain the underlying RateInterest.
//
// Example:
//
//	// Convert a 12% nominal annual rate compounded monthly to its periodic rate
//	periodic, err := NewRateConversion().
//	    Rate(0.12).
//	    Nominal().
//	    Monthly().
//	    ToPeriodic()
//	// periodic is 0.01 (1% monthly)
type RateConversionConfig struct {
	rate      decimal.Decimal
	frequency CompoundingFrequency
	rateType  TypeRate
}

// NewRateConversion creates a new RateConversionConfig builder instance.
// Defaults: frequency is Monthly and rate type is RateEffectyPeriodic.
//
// Configure the source rate with Rate() or RateDecimal(), its type with
// Periodic(), Nominal(), EffectiveAnnual(), AnticipatedPeriodic(),
// AnticipatedNominal(), or RateType(), and its compounding frequency with
// Monthly(), Quarterly(), Annually(), Daily(), or Frequency().
func NewRateConversion() RateConversionConfig {
	return RateConversionConfig{
		frequency: Monthly,
		rateType:  RateEffectyPeriodic,
	}
}

// Rate sets the source interest rate as a float64 decimal (e.g., 0.12 for 12%).
//
// Example:
//
//	.NewRateConversion().Rate(0.12)
func (r RateConversionConfig) Rate(rate float64) RateConversionConfig {
	r.rate = decimal.MustFromFloat64(rate)
	return r
}

// RateDecimal sets the source interest rate using an existing Decimal instance.
// Use this when you already have a Decimal value.
func (r RateConversionConfig) RateDecimal(rate decimal.Decimal) RateConversionConfig {
	r.rate = rate
	return r
}

// Frequency sets the compounding frequency of the source rate.
//
// Parameters:
//   - f: The compounding frequency (Daily, Monthly, Bimonthly, Quarterly,
//     FourMonthly, SemiAnnually, Annually)
func (r RateConversionConfig) Frequency(f CompoundingFrequency) RateConversionConfig {
	r.frequency = f
	return r
}

// Daily sets the compounding frequency of the source rate to daily (365 periods per year).
func (r RateConversionConfig) Daily() RateConversionConfig {
	r.frequency = Daily
	return r
}

// Monthly sets the compounding frequency of the source rate to monthly (12 periods per year).
func (r RateConversionConfig) Monthly() RateConversionConfig {
	r.frequency = Monthly
	return r
}

// Quarterly sets the compounding frequency of the source rate to quarterly (4 periods per year).
func (r RateConversionConfig) Quarterly() RateConversionConfig {
	r.frequency = Quarterly
	return r
}

// SemiAnnually sets the compounding frequency of the source rate to semi-annually (2 periods per year).
func (r RateConversionConfig) SemiAnnually() RateConversionConfig {
	r.frequency = SemiAnnually
	return r
}

// Annually sets the compounding frequency of the source rate to annually (1 period per year).
func (r RateConversionConfig) Annually() RateConversionConfig {
	r.frequency = Annually
	return r
}

// RateType sets the type of the source rate.
//
// Parameters:
//   - t: The rate type (RateEffectyPeriodic, RateEffectyNominal,
//     RateEffectyAnnually, RateAnticipateEffectyPeriodic, RateAnticipateEffectyNominal)
func (r RateConversionConfig) RateType(t TypeRate) RateConversionConfig {
	r.rateType = t
	return r
}

// Periodic marks the source rate as a periodic (per-period) rate.
// Equivalent to RateType(RateEffectyPeriodic).
func (r RateConversionConfig) Periodic() RateConversionConfig {
	r.rateType = RateEffectyPeriodic
	return r
}

// Nominal marks the source rate as a nominal annual rate.
// Equivalent to RateType(RateEffectyNominal).
func (r RateConversionConfig) Nominal() RateConversionConfig {
	r.rateType = RateEffectyNominal
	return r
}

// EffectiveAnnual marks the source rate as an effective annual rate.
// Equivalent to RateType(RateEffectyAnnually).
func (r RateConversionConfig) EffectiveAnnual() RateConversionConfig {
	r.rateType = RateEffectyAnnually
	return r
}

// AnticipatedPeriodic marks the source rate as an anticipated (discount) periodic rate.
// Equivalent to RateType(RateAnticipateEffectyPeriodic).
func (r RateConversionConfig) AnticipatedPeriodic() RateConversionConfig {
	r.rateType = RateAnticipateEffectyPeriodic
	return r
}

// AnticipatedNominal marks the source rate as an anticipated (discount) nominal rate.
// Equivalent to RateType(RateAnticipateEffectyNominal).
func (r RateConversionConfig) AnticipatedNominal() RateConversionConfig {
	r.rateType = RateAnticipateEffectyNominal
	return r
}

// isAnticipated reports whether the configured source rate is an anticipated
// (discount) rate.
func (r RateConversionConfig) isAnticipated() bool {
	return r.rateType == RateAnticipateEffectyPeriodic || r.rateType == RateAnticipateEffectyNominal
}

// Build creates and returns the RateInterest configured in the builder.
//
// Returns:
//   - A RateInterest instance with the configured rate, frequency, and type
//   - An error if validation fails (e.g., negative rate)
func (r RateConversionConfig) Build() (RateInterest, error) {
	return NewRateInterest(r.rate, r.frequency, r.rateType)
}

// MustBuild creates and returns the RateInterest configured in the builder.
// Unlike Build(), this method panics if validation fails.
func (r RateConversionConfig) MustBuild() RateInterest {
	rt, err := r.Build()
	if err != nil {
		panic(err)
	}
	return rt
}

// ToPeriodic converts the configured rate to a periodic (per-period) rate.
// Anticipated source rates are converted through their effective annual equivalent.
//
// Example:
//
//	periodic, err := NewRateConversion().Rate(0.12).Nominal().Monthly().ToPeriodic()
//	// periodic is 0.01 (1% monthly)
func (r RateConversionConfig) ToPeriodic() (decimal.Decimal, error) {
	rt, err := r.Build()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if r.isAnticipated() {
		return rt.ToPeriodic()
	}

	return rt.RatePeriodic()
}

// ToNominal converts the configured rate to a nominal annual rate.
// Anticipated source rates are converted through their effective annual equivalent.
func (r RateConversionConfig) ToNominal() (decimal.Decimal, error) {
	rt, err := r.Build()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if r.isAnticipated() {
		return rt.ToNominal()
	}

	return rt.RateNominal()
}

// ToEffectiveAnnual converts the configured rate to an effective annual rate.
func (r RateConversionConfig) ToEffectiveAnnual() (decimal.Decimal, error) {
	rt, err := r.Build()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if r.isAnticipated() {
		return rt.RateAnticipateEffectyAnnually()
	}

	return rt.RateEffectyAnnually()
}

// ToAnticipatedPeriodic converts the configured rate to an anticipated
// (discount) periodic rate.
func (r RateConversionConfig) ToAnticipatedPeriodic() (decimal.Decimal, error) {
	rt, err := r.Build()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if r.isAnticipated() {
		return rt.RateAnticipatePeriodic()
	}

	return rt.ToAnticipatePeriodic()
}

// ToAnticipatedNominal converts the configured rate to an anticipated
// (discount) nominal rate.
func (r RateConversionConfig) ToAnticipatedNominal() (decimal.Decimal, error) {
	rt, err := r.Build()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if r.isAnticipated() {
		return rt.RateAnticipateNominal()
	}

	return rt.ToAnticipateNominal()
}

// ToPeriodicAt converts the configured rate to the equivalent periodic rate
// at a different compounding frequency.
//
// Example:
//
//	// Convert a 1% monthly periodic rate to its quarterly equivalent
//	quarterly, err := NewRateConversion().Rate(0.01).Periodic().Monthly().ToPeriodicAt(Quarterly)
func (r RateConversionConfig) ToPeriodicAt(newFrequency CompoundingFrequency) (decimal.Decimal, error) {
	rt, err := r.Build()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return rt.RatePeriodicToPeriodic(newFrequency)
}

// ToNominalAt converts the configured rate to the equivalent nominal rate
// at a different compounding frequency.
func (r RateConversionConfig) ToNominalAt(newFrequency CompoundingFrequency) (decimal.Decimal, error) {
	rt, err := r.Build()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return rt.RateNominalToNominal(newFrequency)
}
