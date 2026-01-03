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

	if rt.typeRate == RateEffectyAnnually {
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

	if rt.typeRate == RateEffectyAnnually {
		rateNominal = value * (math.Pow((1+rt.value), (1/value)) - 1)
	}

	if rt.typeRate == RateEffectyPeriodic {
		rateNominal = rt.value * value
	}

	return rateNominal, nil
}

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
