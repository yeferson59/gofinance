package compositeinterest

import (
	"math"
)

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

	if rt.typeRate == RateEffectyAnnualy {
		rateNominal := value * (math.Pow((1+rt.value), (1/value)) - 1)
		ratePeriodic = rateNominal / value
	}

	return ratePeriodic, nil
}

func (rt *RateInterest) RateNominal() (float64, error) {
	if rt.typeRate == RateEffectyNominal {
		return rt.value, nil
	}

	value, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var rateNominal float64

	if rt.typeRate == RateEffectyAnnualy {
		rateNominal = value * (math.Pow((1+rt.value), (1/value)) - 1)
	}

	if rt.typeRate == RateEffectyPeriodic {
		rateNominal = rt.value * value
	}

	return rateNominal, nil
}

func (rt *RateInterest) RateEffectyAnnualy() (float64, error) {
	if rt.typeRate == RateEffectyAnnualy {
		return rt.value, nil
	}

	value, err := getCompoundingFrequency(rt.compoundingFrequency)
	if err != nil {
		return 0, err
	}

	var rateEffectyAnnualy float64

	if rt.typeRate == RateEffectyNominal {
		ratePeriodic := rt.value / value
		rateEffectyAnnualy = math.Pow(1+ratePeriodic, value) - 1
	}

	if rt.typeRate == RateEffectyPeriodic {
		rateEffectyAnnualy = math.Pow(1+rt.value, value) - 1
	}

	return rateEffectyAnnualy, nil
}

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
