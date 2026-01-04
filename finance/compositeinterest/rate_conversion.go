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
func (rt *RateInterest) RatePeriodic() (float64, error) {
	if rt.typeRate == RateEffectyPeriodic {
		return rt.value, nil
	}

	value, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var ratePeriodic float64

	if rt.typeRate == RateEffectyNominal {
		ratePeriodic = rt.value / value
	}

	if rt.typeRate == RateEffectyAnnually {
		rateNominal := value * (math.Pow((1+rt.value), (1/value)) - 1)
		ratePeriodic = rateNominal / value
	}

	return ratePeriodic, nil
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
func (rt *RateInterest) RateNominal() (float64, error) {
	if rt.typeRate == RateEffectyNominal {
		return rt.value, nil
	}

	value, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var rateNominal float64

	if rt.typeRate == RateEffectyAnnually {
		rateNominal = value * (math.Pow((1+rt.value), (1/value)) - 1)
	}

	if rt.typeRate == RateEffectyPeriodic {
		rateNominal = rt.value * value
	}

	return rateNominal, nil
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
func (rt *RateInterest) RateEffectyAnnually() (float64, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	value, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var rateEffectyAnnually float64

	if rt.typeRate == RateEffectyNominal {
		ratePeriodic := rt.value / value
		rateEffectyAnnually = math.Pow(1+ratePeriodic, value) - 1
	}

	if rt.typeRate == RateEffectyPeriodic {
		rateEffectyAnnually = math.Pow(1+rt.value, value) - 1
	}

	return rateEffectyAnnually, nil
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
func (rt *RateInterest) RatePeriodicToPeriodic(newCompoundingFrequency CompoundingFrequency) (float64, error) {
	currentRatePeriodic, err := rt.RatePeriodic()
	if err != nil {
		return 0, err
	}

	newValue, err := getCompoundingFrequency(newCompoundingFrequency)
	if err != nil {
		return 0, err
	}

	value, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	exp := (1 / newValue) * value

	newRatePeriodic := math.Pow(1+currentRatePeriodic, exp) - 1

	return newRatePeriodic, nil
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
func (rt *RateInterest) RateNominalToNominal(newCompoundingFrequency CompoundingFrequency) (float64, error) {
	newPeriodic, err := rt.RatePeriodicToPeriodic(newCompoundingFrequency)
	if err != nil {
		return 0, err
	}

	oldCompoundingFrequency, oldValue, oldTypeRate := rt.compoundingFrequency, rt.value, rt.typeRate

	rt.compoundingFrequency = newCompoundingFrequency
	rt.typeRate = RateEffectyPeriodic
	rt.value = newPeriodic

	newNominal, err := rt.RateNominal()

	rt.compoundingFrequency = oldCompoundingFrequency
	rt.typeRate = oldTypeRate
	rt.value = oldValue

	return newNominal, err
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
func (rt *RateInterest) RateAnticipateEffectyAnnually() (float64, error) {
	if rt.typeRate == RateAnticipateEffectyAnnually {
		return rt.value, nil
	}

	value, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var rateEffectyAnnually float64

	if rt.typeRate == RateAnticipateEffectyNominal {
		ratePeriodic := rt.value / value
		rateEffectyAnnually = math.Pow(1-ratePeriodic, -value) - 1
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		rateEffectyAnnually = math.Pow(1-rt.value, -value) - 1
	}

	return rateEffectyAnnually, nil
}

// RateAnticipateNominal converts an anticipated rate to its equivalent nominal.
//
// Returns:
//   - The equivalent anticipated nominal rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt *RateInterest) RateAnticipateNominal() (float64, error) {
	if rt.typeRate == RateAnticipateEffectyNominal {
		return rt.value, nil
	}

	value, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var rateNominal float64

	if rt.typeRate == RateAnticipateEffectyAnnually {
		rateEffectyAnnually, err := rt.RateAnticipateEffectyAnnually()
		if err != nil {
			return 0, err
		}

		rateNominal = value * (1 - math.Pow(1+rateEffectyAnnually, (-1/value)))
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		rateNominal = rt.value * value
	}

	return rateNominal, nil
}

// RateAnticipatePeriodic converts an anticipated rate to its equivalent periodic.
//
// Returns:
//   - The equivalent anticipated periodic rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt *RateInterest) RateAnticipatePeriodic() (float64, error) {
	if rt.typeRate == RateAnticipateEffectyPeriodic {
		return rt.value, nil
	}

	value, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var ratePeriodic float64

	if rt.typeRate == RateAnticipateEffectyNominal {
		ratePeriodic = rt.value / value
	}

	if rt.typeRate == RateAnticipateEffectyAnnually {
		ratePeriodic = (1 - math.Pow(1+rt.value, (-1/value)))
	}

	return ratePeriodic, nil
}

// ToAnticipateNominal converts an ordinary (vencida) rate to its equivalent anticipated nominal.
// First calculates the effective annual vencida rate, then converts it to anticipated nominal.
//
// Returns:
//   - The equivalent anticipated nominal rate
//   - An error if there are problems in the conversion
func (rt *RateInterest) ToAnticipateNominal() (float64, error) {
	effectyAnnually, err := rt.RateEffectyAnnually()
	if err != nil {
		return 0, err
	}

	oldTypeRate, oldValue := rt.typeRate, rt.value

	rt.typeRate = RateAnticipateEffectyAnnually
	rt.value = effectyAnnually

	rateNominal, err := rt.RateAnticipateNominal()

	rt.typeRate = oldTypeRate
	rt.value = oldValue

	return rateNominal, err
}

// ToAnticipatePeriodic converts an ordinary (vencida) rate to its equivalent anticipated periodic.
// First calculates the effective annual vencida rate, then converts it to anticipated periodic.
//
// Returns:
//   - The equivalent anticipated periodic rate
//   - An error if there are problems in the conversion
func (rt *RateInterest) ToAnticipatePeriodic() (float64, error) {
	effectyAnnually, err := rt.RateEffectyAnnually()
	if err != nil {
		return 0, err
	}

	oldTypeRate, oldValue := rt.typeRate, rt.value

	rt.typeRate = RateAnticipateEffectyAnnually
	rt.value = effectyAnnually

	ratePeriodic, err := rt.RateAnticipatePeriodic()

	rt.typeRate = oldTypeRate
	rt.value = oldValue

	return ratePeriodic, err
}

// ToNominal converts an anticipated rate to its equivalent ordinary (vencida) nominal.
// First calculates the anticipated effective annual rate, then converts it to ordinary nominal.
//
// Returns:
//   - The equivalent vencida nominal rate
//   - An error if there are problems in the conversion
func (rt *RateInterest) ToNominal() (float64, error) {
	effectyAnnually, err := rt.RateAnticipateEffectyAnnually()
	if err != nil {
		return 0, err
	}

	oldTypeRate, oldValue := rt.typeRate, rt.value

	rt.typeRate = RateEffectyAnnually
	rt.value = effectyAnnually

	rateNominal, err := rt.RateNominal()

	rt.typeRate = oldTypeRate
	rt.value = oldValue

	return rateNominal, err
}

// ToPeriodic converts an anticipated rate to its equivalent ordinary (vencida) periodic.
// First calculates the anticipated effective annual rate, then converts it to ordinary periodic.
//
// Returns:
//   - The equivalent vencida periodic rate
//   - An error if there are problems in the conversion
func (rt *RateInterest) ToPeriodic() (float64, error) {
	effectyAnnually, err := rt.RateAnticipateEffectyAnnually()
	if err != nil {
		return 0, err
	}

	oldTypeRate, oldValue := rt.typeRate, rt.value

	rt.typeRate = RateEffectyAnnually
	rt.value = effectyAnnually

	ratePeriodic, err := rt.RatePeriodic()

	rt.typeRate = oldTypeRate
	rt.value = oldValue

	return ratePeriodic, err
}
