package compositeinterest

import (
	"math"

	"github.com/yeferson59/gofinance/money"
)

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
		pow := math.Pow((rt.value.Add(money.One).InexactFloat64()), (money.One.MustDiv(compoundingPeriodsPerYear).InexactFloat64()))
		nominalRate := compoundingPeriodsPerYear.Mul(money.MustFromFloat64(pow).Sub(money.One))
		periodicRate = nominalRate.MustDiv(compoundingPeriodsPerYear)
	}

	return periodicRate, nil
}

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
		pow := math.Pow((rt.value.Add(money.One).InexactFloat64()), (money.One.MustDiv(compoundingPeriodsPerYear).InexactFloat64()))
		nominalRate = compoundingPeriodsPerYear.Mul(money.MustFromFloat64(pow).Sub(money.One))
	}

	if rt.typeRate == RateEffectyPeriodic {
		nominalRate = rt.value.Mul(compoundingPeriodsPerYear)
	}

	return nominalRate, nil
}

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
		effectiveAnnualRate = money.MustFromFloat64(math.Pow(periodicRate.Add(money.One).InexactFloat64(), compoundingPeriodsPerYear.InexactFloat64())).Sub(money.One)
	}

	if rt.typeRate == RateEffectyPeriodic {
		effectiveAnnualRate = money.MustFromFloat64(math.Pow(rt.value.Add(money.One).InexactFloat64(), compoundingPeriodsPerYear.InexactFloat64())).Sub(money.One)
	}

	return effectiveAnnualRate, nil
}

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

	newPeriodicRate := math.Pow(currentPeriodicRate.Add(money.One).InexactFloat64(), exponent.InexactFloat64()) - 1

	return money.MustFromFloat64(newPeriodicRate), nil
}

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
		effectiveAnnualRate = money.MustFromFloat64(math.Pow(money.One.Sub(periodicRate).InexactFloat64(), money.One.Neg().Mul(compoundingPeriodsPerYear).InexactFloat64())).Sub(money.One)
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		effectiveAnnualRate = money.MustFromFloat64(math.Pow(money.One.Sub(rt.value).InexactFloat64(), money.One.Neg().Mul(compoundingPeriodsPerYear).InexactFloat64())).Sub(money.One)
	}

	return effectiveAnnualRate, nil
}

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

		nominalRate = compoundingPeriodsPerYear.Mul(money.MustFromFloat64((1 - math.Pow(effectiveAnnualRate.Add(money.One).InexactFloat64(), (money.One.Neg().MustDiv(compoundingPeriodsPerYear).InexactFloat64())))))
	}

	if rt.typeRate == RateAnticipateEffectyPeriodic {
		nominalRate = rt.value.Mul(compoundingPeriodsPerYear)
	}

	return nominalRate, nil
}

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
		periodicRate = money.MustFromFloat64((1 - math.Pow(rt.value.Add(money.One).InexactFloat64(), (money.One.Neg().MustDiv(compoundingPeriodsPerYear).InexactFloat64()))))
	}

	return periodicRate, nil
}

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
