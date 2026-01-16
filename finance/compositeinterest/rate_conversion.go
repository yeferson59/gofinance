package compositeinterest

import (
	"math"
)

// RatePeriodic converts the interest rate to periodic type.
// If the rate is already periodic, returns the value unchanged.
// If it is nominal, divides by the compounding frequency.
// If it is effective annual, first converts to nominal and then to periodic.
//
// The periodic rate is the one applied in each individual compounding period.
//
// Returns:
//   - The equivalent periodic rate
//   - An error if the valid compounding frequency cannot be obtained
//
// Example:
//
// rate, _ := NewRateInterest(0.12, Monthly, RateEffectyNominal)
// periodic, err := rate.RatePeriodic()
// // periodic is 0.01 (1% monthly)
func (rt RateInterest) RatePeriodic() (float64, error) {
	if rt.typeRate == RateEffectyPeriodic {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var periodicRate float64

	if rt.typeRate == RateEffectyNominal {
		periodicRate = rt.value / compoundingPeriodsPerYear
	}

	if rt.typeRate == RateEffectyAnnually {
		// Cache the power calculation to avoid redundant computation
		pow := math.Pow((1 + rt.value), (1 / compoundingPeriodsPerYear))
		nominalRate := compoundingPeriodsPerYear * (pow - 1)
		periodicRate = nominalRate / compoundingPeriodsPerYear
	}

	return periodicRate, nil
}

// RateNominal converts the interest rate to nominal type.
// If the rate is already nominal, returns the value unchanged.
// If it is effective annual, performs the conversion from effective to nominal rate.
// If it is periodic, multiplies by the compounding frequency.
//
// The nominal rate is an annual rate that does not consider compounding.
// It is divided by the frequency to obtain the periodic rate.
//
// Returns:
//   - The equivalent annual nominal rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt RateInterest) RateNominal() (float64, error) {
	if rt.typeRate == RateEffectyNominal {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var nominalRate float64

	if rt.typeRate == RateEffectyAnnually {
		// Cache the power calculation to avoid redundant computation
		pow := math.Pow((1 + rt.value), (1 / compoundingPeriodsPerYear))
		nominalRate = compoundingPeriodsPerYear * (pow - 1)
	}

	if rt.typeRate == RateEffectyPeriodic {
		nominalRate = rt.value * compoundingPeriodsPerYear
	}

	return nominalRate, nil
}

// RateEffectyAnnually converts the interest rate to effective annual rate.
// If the rate is already effective annual, returns the value unchanged.
// The effective annual rate is the real rate earned considering compounding.
//
// Example: A nominal rate of 12% compounded monthly equals approximately
// an effective annual rate of 12.68%.
//
// Returns:
//   - The equivalent effective annual rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt RateInterest) RateEffectyAnnually() (float64, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var effectiveAnnualRate float64

	if rt.typeRate == RateEffectyNominal {
		periodicRate := rt.value / compoundingPeriodsPerYear
		effectiveAnnualRate = math.Pow(1+periodicRate, compoundingPeriodsPerYear) - 1
	}

	if rt.typeRate == RateEffectyPeriodic {
		effectiveAnnualRate = math.Pow(1+rt.value, compoundingPeriodsPerYear) - 1
	}

	return effectiveAnnualRate, nil
}

// RatePeriodicToPeriodic converts a periodic rate to another periodic rate
// with a different compounding frequency.
//
// Example: Convert a monthly rate of 1% to its equivalent quarterly rate.
// If you have 1% monthly, the equivalent quarterly would be higher because 3 months
// of monthly compounding equals 1 quarterly period.
//
// Parameters:
//   - newCompoundingFrequency: The desired new compounding frequency
//
// Returns:
//   - The equivalent periodic rate in the new frequency
//   - An error if the valid compounding frequency cannot be obtained
func (rt RateInterest) RatePeriodicToPeriodic(newCompoundingFrequency CompoundingFrequency) (float64, error) {
	currentPeriodicRate, err := rt.RatePeriodic()
	if err != nil {
		return 0, err
	}

	newPeriodsPerYear, err := getCompoundingFrequency(newCompoundingFrequency)
	if err != nil {
		return 0, err
	}

	currentPeriodsPerYear, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	exponent := (1 / newPeriodsPerYear) * currentPeriodsPerYear

	newPeriodicRate := math.Pow(1+currentPeriodicRate, exponent) - 1

	return newPeriodicRate, nil
}

// RateNominalToNominal converts a nominal rate to another nominal rate
// with a different compounding frequency.
//
// Parameters:
//   - newCompoundingFrequency: The desired new compounding frequency
//
// Returns:
//   - The equivalent nominal annual rate in the new frequency
//   - An error if there are problems in the conversion
func (rt RateInterest) RateNominalToNominal(newCompoundingFrequency CompoundingFrequency) (float64, error) {
	newPeriodicRate, err := rt.RatePeriodicToPeriodic(newCompoundingFrequency)
	if err != nil {
		return 0, err
	}

	newRateInterest, err := NewRateInterest(newPeriodicRate, newCompoundingFrequency, RateEffectyPeriodic)
	if err != nil {
		return 0, err
	}

	return newRateInterest.RateNominal()
}

// RateAnticipateEffectyAnnually converts an anticipated (discount) rate
// to its equivalent annual.
//
// Anticipated rates are charged at the beginning of the period rather than at the end.
// They are commonly used in bill of exchange discounts and other financial instruments.
//
// Returns:
//   - The equivalent anticipated effective annual rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt RateInterest) RateAnticipateEffectyAnnually() (float64, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var effectiveAnnualRate float64

	if rt.typeRate == RateAnticipateEffectyNominal {
		periodicRate := rt.value / compoundingPeriodsPerYear
		effectiveAnnualRate = math.Pow(1-periodicRate, -compoundingPeriodsPerYear) - 1
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		effectiveAnnualRate = math.Pow(1-rt.value, -compoundingPeriodsPerYear) - 1
	}

	return effectiveAnnualRate, nil
}

// RateAnticipateNominal converts an anticipated rate to its equivalent nominal.
//
// Returns:
//   - The equivalent anticipated nominal rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt RateInterest) RateAnticipateNominal() (float64, error) {
	if rt.typeRate == RateAnticipateEffectyNominal {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var nominalRate float64

	if rt.typeRate == RateEffectyAnnually {
		effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
		if err != nil {
			return 0, err
		}

		nominalRate = compoundingPeriodsPerYear * (1 - math.Pow(1+effectiveAnnualRate, (-1/compoundingPeriodsPerYear)))
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		nominalRate = rt.value * compoundingPeriodsPerYear
	}

	return nominalRate, nil
}

// RateAnticipatePeriodic converts an anticipated rate to its equivalent periodic.
//
// Returns:
//   - The equivalent anticipated periodic rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt RateInterest) RateAnticipatePeriodic() (float64, error) {
	if rt.typeRate == RateAnticipateEffectyPeriodic {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var periodicRate float64

	if rt.typeRate == RateAnticipateEffectyNominal {
		periodicRate = rt.value / compoundingPeriodsPerYear
	}

	if rt.typeRate == RateEffectyAnnually {
		periodicRate = (1 - math.Pow(1+rt.value, (-1/compoundingPeriodsPerYear)))
	}

	return periodicRate, nil
}

// ToAnticipateNominal converts an ordinary (vencida) rate to its equivalent anticipated nominal.
// First calculates the effective annual vencida rate, then converts it to anticipated nominal.
//
// Returns:
//   - The equivalent anticipated nominal rate
//   - An error if there are problems in the conversion
func (rt RateInterest) ToAnticipateNominal() (float64, error) {
	effectiveAnnualRate, err := rt.RateEffectyAnnually()
	if err != nil {
		return 0, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return 0, err
	}

	return newRateInterest.RateAnticipateNominal()
}

// ToAnticipatePeriodic converts an ordinary (vencida) rate to its equivalent anticipated periodic.
// First calculates the effective annual vencida rate, then converts it to anticipated periodic.
//
// Returns:
//   - The equivalent anticipated periodic rate
//   - An error if there are problems in the conversion
func (rt RateInterest) ToAnticipatePeriodic() (float64, error) {
	effectiveAnnualRate, err := rt.RateEffectyAnnually()
	if err != nil {
		return 0, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return 0, err
	}

	return newRateInterest.RateAnticipatePeriodic()
}

// ToNominal converts an anticipated rate to its equivalent ordinary (vencida) nominal.
// First calculates the anticipated effective annual rate, then converts it to ordinary nominal.
//
// Returns:
//   - The equivalent vencida nominal rate
//   - An error if there are problems in the conversion
func (rt RateInterest) ToNominal() (float64, error) {
	effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
	if err != nil {
		return 0, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return 0, err
	}

	return newRateInterest.RateNominal()
}

// ToPeriodic converts an anticipated rate to its equivalent ordinary (vencida) periodic.
// First calculates the anticipated effective annual rate, then converts it to ordinary periodic.
//
// Returns:
//   - The equivalent vencida periodic rate
//   - An error if there are problems in the conversion
func (rt RateInterest) ToPeriodic() (float64, error) {
	effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
	if err != nil {
		return 0, err
	}

	newRateInterest, err := NewRateInterest(effectiveAnnualRate, rt.compoundingFrequency, RateEffectyAnnually)
	if err != nil {
		return 0, err
	}

	return newRateInterest.RatePeriodic()
}
