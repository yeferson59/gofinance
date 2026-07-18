package compoundinterest

import "github.com/yeferson59/gofinance/decimal"

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
func (rt RateInterest) RatePeriodic() (decimal.Decimal, error) {
	if rt.typeRate == RateEffectyPeriodic {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return decimal.Decimal{}, err
	}

	if rt.typeRate == RateEffectyNominal {
		periodicRate, err := rt.value.Div(compoundingPeriodsPerYear)
		if err != nil {
			return decimal.Decimal{}, err
		}

		return periodicRate, nil
	}

	if rt.typeRate == RateEffectyAnnually {
		div, err := decimal.One.Div(compoundingPeriodsPerYear)
		if err != nil {
			return decimal.Decimal{}, err
		}

		pow, err := rt.value.Add(decimal.One).Pow(div)
		if err != nil {
			return decimal.Decimal{}, err
		}

		periodicRate, err := compoundingPeriodsPerYear.Mul(pow.Sub(decimal.One)).Div(compoundingPeriodsPerYear)
		if err != nil {
			return decimal.Decimal{}, err
		}

		return periodicRate, nil
	}

	return decimal.Decimal{}, nil
}

// RateNominal converts the current interest rate to a nominal rate.
// A nominal rate is the stated rate before considering compounding frequency.
//
// Returns:
//   - The nominal interest rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateNominal() (decimal.Decimal, error) {
	if rt.typeRate == RateEffectyNominal {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return decimal.Decimal{}, err
	}

	var nominalRate decimal.Decimal

	if rt.typeRate == RateEffectyAnnually {
		pow := rt.value.Add(decimal.One).MustPow(decimal.One.MustDiv(compoundingPeriodsPerYear))
		nominalRate = compoundingPeriodsPerYear.Mul(pow.Sub(decimal.One))
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
func (rt RateInterest) RateEffectyAnnually() (decimal.Decimal, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return decimal.Decimal{}, err
	}

	var effectiveAnnualRate decimal.Decimal

	if rt.typeRate == RateEffectyNominal {
		periodicRate := rt.value.MustDiv(compoundingPeriodsPerYear)
		effectiveAnnualRate = periodicRate.Add(decimal.One).MustPow(compoundingPeriodsPerYear).Sub(decimal.One)
	}

	if rt.typeRate == RateEffectyPeriodic {
		effectiveAnnualRate = rt.value.Add(decimal.One).MustPow(compoundingPeriodsPerYear).Sub(decimal.One)
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
func (rt RateInterest) RatePeriodicToPeriodic(newCompoundingFrequency CompoundingFrequency) (decimal.Decimal, error) {
	currentPeriodicRate, err := rt.RatePeriodic()
	if err != nil {
		return decimal.Decimal{}, err
	}

	newPeriodsPerYear, err := newCompoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return decimal.Decimal{}, err
	}

	currentPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return decimal.Decimal{}, err
	}

	exponent := decimal.One.MustDiv(newPeriodsPerYear).Mul(currentPeriodsPerYear)

	newPeriodicRate := currentPeriodicRate.Add(decimal.One).MustPow(exponent).Sub(decimal.One)

	return newPeriodicRate, nil
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
func (rt RateInterest) RateNominalToNominal(newCompoundingFrequency CompoundingFrequency) (decimal.Decimal, error) {
	newPeriodicRate, err := rt.RatePeriodicToPeriodic(newCompoundingFrequency)
	if err != nil {
		return decimal.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(newPeriodicRate, newCompoundingFrequency, RateEffectyPeriodic)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return newRateInterest.RateNominal()
}

// RateAnticipateEffectyAnnually converts an anticipated (discount) rate to an effective annual rate.
// Anticipated rates are used in discount instruments where interest is deducted upfront.
//
// Returns:
//   - The effective annual rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateAnticipateEffectyAnnually() (decimal.Decimal, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return decimal.Decimal{}, err
	}

	var effectiveAnnualRate decimal.Decimal

	if rt.typeRate == RateAnticipateEffectyNominal {
		periodicRate := rt.value.MustDiv(compoundingPeriodsPerYear)
		effectiveAnnualRate = decimal.One.Sub(periodicRate).MustPow(decimal.One.Neg().Mul(compoundingPeriodsPerYear)).Sub(decimal.One)
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		effectiveAnnualRate = decimal.One.Sub(rt.value).MustPow(decimal.One.Neg().Mul(compoundingPeriodsPerYear)).Sub(decimal.One)
	}

	return effectiveAnnualRate, nil
}

// RateAnticipateNominal converts the current rate to an anticipated (discount) nominal rate.
//
// Returns:
//   - The anticipated nominal rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateAnticipateNominal() (decimal.Decimal, error) {
	if rt.typeRate == RateAnticipateEffectyNominal {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return decimal.Decimal{}, err
	}

	var nominalRate decimal.Decimal

	if rt.typeRate == RateEffectyAnnually {
		effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
		if err != nil {
			return decimal.Decimal{}, err
		}

		pow := effectiveAnnualRate.Add(decimal.One).MustPow(decimal.One.Neg().MustDiv(compoundingPeriodsPerYear))
		nominalRate = compoundingPeriodsPerYear.Mul(decimal.One.Sub(pow))
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
func (rt RateInterest) RateAnticipatePeriodic() (decimal.Decimal, error) {
	if rt.typeRate == RateAnticipateEffectyPeriodic {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return decimal.Decimal{}, err
	}

	var periodicRate decimal.Decimal

	if rt.typeRate == RateAnticipateEffectyNominal {
		periodicRate = rt.value.MustDiv(compoundingPeriodsPerYear)
	}

	if rt.typeRate == RateEffectyAnnually {
		pow := rt.value.Add(decimal.One).MustPow(decimal.One.Neg().MustDiv(compoundingPeriodsPerYear))
		periodicRate = decimal.One.Sub(pow)
	}

	return periodicRate, nil
}

// ToAnticipateNominal converts an effective or periodic rate to an anticipated nominal rate.
// This is useful for discount instruments like Treasury bills.
//
// Returns:
//   - The anticipated nominal rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToAnticipateNominal() (decimal.Decimal, error) {
	effectiveAnnualRate, err := rt.RateEffectyAnnually()
	if err != nil {
		return decimal.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return newRateInterest.RateAnticipateNominal()
}

// ToAnticipatePeriodic converts an effective or periodic rate to an anticipated periodic rate.
//
// Returns:
//   - The anticipated periodic rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToAnticipatePeriodic() (decimal.Decimal, error) {
	effectiveAnnualRate, err := rt.RateEffectyAnnually()
	if err != nil {
		return decimal.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return newRateInterest.RateAnticipatePeriodic()
}

// ToNominal converts an anticipated rate to a nominal rate.
//
// Returns:
//   - The nominal rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToNominal() (decimal.Decimal, error) {
	effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
	if err != nil {
		return decimal.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return newRateInterest.RateNominal()
}

// ToPeriodic converts an anticipated rate to a periodic rate.
//
// Returns:
//   - The periodic rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToPeriodic() (decimal.Decimal, error) {
	effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
	if err != nil {
		return decimal.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return newRateInterest.RatePeriodic()
}
