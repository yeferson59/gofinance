package compositeinterest

import (
	"github.com/yeferson59/gofinance/money"
)

// RatePeriodic converts the current interest rate to a periodic (per-period) rate.
// This is useful when you have a nominal or effective annual rate and need the actual
// rate applied to each compounding period.
//
// Returns:
//   - The periodic interest rate as a Decimal
//   - An error if the conversion fails
//
// Example:
//
//	rate, _ := NewRateInterest(0.12, Monthly, RateEffectyNominal)
//	periodicRate, _ := rate.RatePeriodic()  // Converts 12% nominal to ~1% monthly
func (rt RateInterest) RatePeriodic() (money.Decimal, error) {
	if rt.typeRate == RateEffectyPeriodic {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var periodicRate money.Decimal

	if rt.typeRate == RateEffectyNominal {
		periodicRate = rt.value.MustDiv(compoundingPeriodsPerYear)
	}

	if rt.typeRate == RateEffectyAnnually {
		pow := rt.value.Add(money.One).MustPow(money.One.MustDiv(compoundingPeriodsPerYear))
		nominalRate := compoundingPeriodsPerYear.Mul(pow.Sub(money.One))
		periodicRate = nominalRate.MustDiv(compoundingPeriodsPerYear)
	}

	return periodicRate, nil
}

// RateNominal converts the current interest rate to a nominal rate.
// A nominal rate is the stated rate before considering compounding frequency.
//
// Returns:
//   - The nominal interest rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateNominal() (money.Decimal, error) {
	if rt.typeRate == RateEffectyNominal {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var nominalRate money.Decimal

	if rt.typeRate == RateEffectyAnnually {
		pow := rt.value.Add(money.One).MustPow(money.One.MustDiv(compoundingPeriodsPerYear))
		nominalRate = compoundingPeriodsPerYear.Mul(pow.Sub(money.One))
	}

	if rt.typeRate == RateEffectyPeriodic {
		nominalRate = rt.value.Mul(compoundingPeriodsPerYear)
	}

	return nominalRate, nil
}

// RateEffectyAnnually converts the current interest rate to an effective annual rate.
// The effective annual rate accounts for compounding and represents the true yearly return.
//
// Returns:
//   - The effective annual interest rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateEffectyAnnually() (money.Decimal, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var effectiveAnnualRate money.Decimal

	if rt.typeRate == RateEffectyNominal {
		periodicRate := rt.value.MustDiv(compoundingPeriodsPerYear)
		effectiveAnnualRate = periodicRate.Add(money.One).MustPow(compoundingPeriodsPerYear).Sub(money.One)
	}

	if rt.typeRate == RateEffectyPeriodic {
		effectiveAnnualRate = rt.value.Add(money.One).MustPow(compoundingPeriodsPerYear).Sub(money.One)
	}

	return effectiveAnnualRate, nil
}

// RatePeriodicToPeriodic converts a periodic rate from one compounding frequency to another.
// For example, converts a monthly rate to a quarterly rate.
//
// Parameters:
//   - newCompoundingFrequency: The target compounding frequency
//
// Returns:
//   - The converted periodic rate for the new frequency as a Decimal
//   - An error if the conversion fails
//
// Example:
//
//	rate, _ := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
//	quarterlyRate, _ := rate.RatePeriodicToPeriodic(QuarterlyOne)
func (rt RateInterest) RatePeriodicToPeriodic(newCompoundingFrequency CompoundingFrequency) (money.Decimal, error) {
	currentPeriodicRate, err := rt.RatePeriodic()
	if err != nil {
		return money.Decimal{}, err
	}

	newPeriodsPerYear, err := newCompoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	currentPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	exponent := money.One.MustDiv(newPeriodsPerYear).Mul(currentPeriodsPerYear)

	newPeriodicRate := currentPeriodicRate.Add(money.One).MustPow(exponent).InexactFloat64() - 1

	return money.MustFromFloat64(newPeriodicRate), nil
}

// RateNominalToNominal converts a nominal rate from one compounding frequency to another.
// For example, converts a monthly nominal rate to a quarterly nominal rate.
//
// Parameters:
//   - newCompoundingFrequency: The target compounding frequency
//
// Returns:
//   - The converted nominal rate for the new frequency as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateNominalToNominal(newCompoundingFrequency CompoundingFrequency) (money.Decimal, error) {
	newPeriodicRate, err := rt.RatePeriodicToPeriodic(newCompoundingFrequency)
	if err != nil {
		return money.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(newPeriodicRate, newCompoundingFrequency, RateEffectyPeriodic)
	if err != nil {
		return money.Decimal{}, err
	}

	return newRateInterest.RateNominal()
}

// RateAnticipateEffectyAnnually converts an anticipated (discount) rate to an effective annual rate.
// Anticipated rates are used in discount instruments where interest is deducted upfront.
//
// Returns:
//   - The effective annual rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateAnticipateEffectyAnnually() (money.Decimal, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var effectiveAnnualRate money.Decimal

	if rt.typeRate == RateAnticipateEffectyNominal {
		periodicRate := rt.value.MustDiv(compoundingPeriodsPerYear)
		effectiveAnnualRate = money.One.Sub(periodicRate).MustPow(money.One.Neg().Mul(compoundingPeriodsPerYear)).Sub(money.One)
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		effectiveAnnualRate = money.One.Sub(rt.value).MustPow(money.One.Neg().Mul(compoundingPeriodsPerYear)).Sub(money.One)
	}

	return effectiveAnnualRate, nil
}

// RateAnticipateNominal converts the current rate to an anticipated (discount) nominal rate.
//
// Returns:
//   - The anticipated nominal rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateAnticipateNominal() (money.Decimal, error) {
	if rt.typeRate == RateAnticipateEffectyNominal {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var nominalRate money.Decimal

	if rt.typeRate == RateEffectyAnnually {
		effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
		if err != nil {
			return money.Decimal{}, err
		}

		pow := effectiveAnnualRate.Add(money.One).MustPow(money.One.Neg().MustDiv(compoundingPeriodsPerYear)).InexactFloat64()
		nominalRate = compoundingPeriodsPerYear.Mul(money.MustFromFloat64(1 - pow))
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		nominalRate = rt.value.Mul(compoundingPeriodsPerYear)
	}

	return nominalRate, nil
}

// RateAnticipatePeriodic converts the current rate to an anticipated (discount) periodic rate.
//
// Returns:
//   - The anticipated periodic rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateAnticipatePeriodic() (money.Decimal, error) {
	if rt.typeRate == RateAnticipateEffectyPeriodic {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var periodicRate money.Decimal

	if rt.typeRate == RateAnticipateEffectyNominal {
		periodicRate = rt.value.MustDiv(compoundingPeriodsPerYear)
	}

	if rt.typeRate == RateEffectyAnnually {
		pow := rt.value.Add(money.One).MustPow(money.One.Neg().MustDiv(compoundingPeriodsPerYear)).InexactFloat64()
		periodicRate = money.MustFromFloat64(1 - pow)
	}

	return periodicRate, nil
}

// ToAnticipateNominal converts an effective or periodic rate to an anticipated nominal rate.
// This is useful for discount instruments like Treasury bills.
//
// Returns:
//   - The anticipated nominal rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToAnticipateNominal() (money.Decimal, error) {
	effectiveAnnualRate, err := rt.RateEffectyAnnually()
	if err != nil {
		return money.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return money.Decimal{}, err
	}

	return newRateInterest.RateAnticipateNominal()
}

// ToAnticipatePeriodic converts an effective or periodic rate to an anticipated periodic rate.
//
// Returns:
//   - The anticipated periodic rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToAnticipatePeriodic() (money.Decimal, error) {
	effectiveAnnualRate, err := rt.RateEffectyAnnually()
	if err != nil {
		return money.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return money.Decimal{}, err
	}

	return newRateInterest.RateAnticipatePeriodic()
}

// ToNominal converts an anticipated rate to a nominal rate.
//
// Returns:
//   - The nominal rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToNominal() (money.Decimal, error) {
	effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
	if err != nil {
		return money.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return money.Decimal{}, err
	}

	return newRateInterest.RateNominal()
}

// ToPeriodic converts an anticipated rate to a periodic rate.
//
// Returns:
//   - The periodic rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToPeriodic() (money.Decimal, error) {
	effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
	if err != nil {
		return money.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return money.Decimal{}, err
	}

	return newRateInterest.RatePeriodic()
}
