package compositeinterest

import (
	"math"

	"github.com/quagmt/udecimal"
	"github.com/yeferson59/gofinance/money"
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
func (rt RateInterest) RatePeriodic() (money.Decimal, error) {
	if rt.typeRate == RateEffectyPeriodic {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var periodicRate udecimal.Decimal

	if rt.typeRate == RateEffectyNominal {
		periodicRate = rt.value.MustDiv(compoundingPeriodsPerYear.Decimal)
	}

	if rt.typeRate == RateEffectyAnnually {
		pow := math.Pow((rt.value.Add(udecimal.One).InexactFloat64()), (udecimal.One.MustDiv(compoundingPeriodsPerYear.Decimal).InexactFloat64()))
		nominalRate := compoundingPeriodsPerYear.Mul(udecimal.MustFromFloat64(pow).Sub(udecimal.One))
		periodicRate = nominalRate.MustDiv(compoundingPeriodsPerYear.Decimal)
	}

	return money.Decimal{Decimal: periodicRate}, nil
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
func (rt RateInterest) RateNominal() (money.Decimal, error) {
	if rt.typeRate == RateEffectyNominal {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var nominalRate udecimal.Decimal

	if rt.typeRate == RateEffectyAnnually {
		// Cache the power calculation to avoid redundant computation
		pow := math.Pow((rt.value.Add(udecimal.One).InexactFloat64()), (udecimal.One.MustDiv(compoundingPeriodsPerYear.Decimal).InexactFloat64()))
		nominalRate = compoundingPeriodsPerYear.Decimal.Mul(udecimal.MustFromFloat64(pow).Sub(udecimal.One))
	}

	if rt.typeRate == RateEffectyPeriodic {
		nominalRate = rt.value.Mul(compoundingPeriodsPerYear.Decimal)
	}

	return money.Decimal{Decimal: nominalRate}, nil
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
func (rt RateInterest) RateEffectyAnnually() (money.Decimal, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var effectiveAnnualRate udecimal.Decimal

	if rt.typeRate == RateEffectyNominal {
		periodicRate := rt.value.MustDiv(compoundingPeriodsPerYear.Decimal)
		effectiveAnnualRate = udecimal.MustFromFloat64(math.Pow(periodicRate.Add(udecimal.One).InexactFloat64(), compoundingPeriodsPerYear.InexactFloat64())).Sub(udecimal.One)
	}

	if rt.typeRate == RateEffectyPeriodic {
		effectiveAnnualRate = udecimal.MustFromFloat64(math.Pow(rt.value.Add(udecimal.One).InexactFloat64(), compoundingPeriodsPerYear.InexactFloat64())).Sub(udecimal.One)
	}

	return money.Decimal{Decimal: effectiveAnnualRate}, nil
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

	exponent := udecimal.One.MustDiv(newPeriodsPerYear.Decimal).Mul(currentPeriodsPerYear.Decimal)

	newPeriodicRate := math.Pow(currentPeriodicRate.Add(udecimal.One).InexactFloat64(), exponent.InexactFloat64()) - 1

	return money.Decimal{Decimal: udecimal.MustFromFloat64(newPeriodicRate)}, nil
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

// RateAnticipateEffectyAnnually converts an anticipated (discount) rate
// to its equivalent annual.
//
// Anticipated rates are charged at the beginning of the period rather than at the end.
// They are commonly used in bill of exchange discounts and other financial instruments.
//
// Returns:
//   - The equivalent anticipated effective annual rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt RateInterest) RateAnticipateEffectyAnnually() (money.Decimal, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var effectiveAnnualRate udecimal.Decimal

	if rt.typeRate == RateAnticipateEffectyNominal {
		periodicRate := rt.value.MustDiv(compoundingPeriodsPerYear.Decimal)
		effectiveAnnualRate = udecimal.MustFromFloat64(math.Pow(udecimal.One.Sub(periodicRate).InexactFloat64(), udecimal.One.Neg().Mul(compoundingPeriodsPerYear.Decimal).InexactFloat64())).Sub(udecimal.One)
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		effectiveAnnualRate = udecimal.MustFromFloat64(math.Pow(udecimal.One.Sub(rt.value.Decimal).InexactFloat64(), udecimal.One.Neg().Mul(compoundingPeriodsPerYear.Decimal).InexactFloat64())).Sub(udecimal.One)
	}

	return money.Decimal{Decimal: effectiveAnnualRate}, nil
}

// RateAnticipateNominal converts an anticipated rate to its equivalent nominal.
//
// Returns:
//   - The equivalent anticipated nominal rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt RateInterest) RateAnticipateNominal() (money.Decimal, error) {
	if rt.typeRate == RateAnticipateEffectyNominal {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var nominalRate udecimal.Decimal

	if rt.typeRate == RateEffectyAnnually {
		effectiveAnnualRate, err := rt.RateAnticipateEffectyAnnually()
		if err != nil {
			return money.Decimal{}, err
		}

		nominalRate = compoundingPeriodsPerYear.Mul(udecimal.MustFromFloat64((1 - math.Pow(effectiveAnnualRate.Add(udecimal.One).InexactFloat64(), (udecimal.One.Neg().MustDiv(compoundingPeriodsPerYear.Decimal).InexactFloat64())))))
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		nominalRate = rt.value.Mul(compoundingPeriodsPerYear.Decimal)
	}

	return money.Decimal{Decimal: nominalRate}, nil
}

// RateAnticipatePeriodic converts an anticipated rate to its equivalent periodic.
//
// Returns:
//   - The equivalent anticipated periodic rate
//   - An error if the valid compounding frequency cannot be obtained
func (rt RateInterest) RateAnticipatePeriodic() (money.Decimal, error) {
	if rt.typeRate == RateAnticipateEffectyPeriodic {
		return rt.value, nil
	}

	compoundingPeriodsPerYear, err := rt.compoundingFrequency.getCompoundingFrequency()
	if err != nil {
		return money.Decimal{}, err
	}

	var periodicRate udecimal.Decimal

	if rt.typeRate == RateAnticipateEffectyNominal {
		periodicRate = rt.value.MustDiv(compoundingPeriodsPerYear.Decimal)
	}

	if rt.typeRate == RateEffectyAnnually {
		periodicRate = udecimal.MustFromFloat64((1 - math.Pow(rt.value.Add(udecimal.One).InexactFloat64(), (udecimal.One.Neg().MustDiv(compoundingPeriodsPerYear.Decimal).InexactFloat64()))))
	}

	return money.Decimal{Decimal: periodicRate}, nil
}

// ToAnticipateNominal converts an ordinary (vencida) rate to its equivalent anticipated nominal.
// First calculates the effective annual vencida rate, then converts it to anticipated nominal.
//
// Returns:
//   - The equivalent anticipated nominal rate
//   - An error if there are problems in the conversion
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

// ToAnticipatePeriodic converts an ordinary (vencida) rate to its equivalent anticipated periodic.
// First calculates the effective annual vencida rate, then converts it to anticipated periodic.
//
// Returns:
//   - The equivalent anticipated periodic rate
//   - An error if there are problems in the conversion
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

// ToNominal converts an anticipated rate to its equivalent ordinary (vencida) nominal.
// First calculates the anticipated effective annual rate, then converts it to ordinary nominal.
//
// Returns:
//   - The equivalent vencida nominal rate
//   - An error if there are problems in the conversion
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

// ToPeriodic converts an anticipated rate to its equivalent ordinary (vencida) periodic.
// First calculates the anticipated effective annual rate, then converts it to ordinary periodic.
//
// Returns:
//   - The equivalent vencida periodic rate
//   - An error if there are problems in the conversion
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
