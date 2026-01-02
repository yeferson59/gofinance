package compositeinterest

import (
	"errors"
)

type CompoundingFrequency string
type TypeRate string

type Period struct {
	daily        *float64
	monthly      *float64
	bimonthly    *float64
	quarterlyOne *float64
	quarterlyTwo *float64
	semiAnnually *float64
	annually     *float64
}

func NewPeriod(numberPeriods float64, compoundingFrequency CompoundingFrequency) (*Period, error) {
	switch compoundingFrequency {
	case Daily:
		return &Period{
			daily: &numberPeriods,
		}, nil
	case Monthly:
		return &Period{
			monthly: &numberPeriods,
		}, nil
	case Bimonthly:
		return &Period{
			bimonthly: &numberPeriods,
		}, nil
	case QuarterlyOne:
		return &Period{
			quarterlyOne: &numberPeriods,
		}, nil
	case QuarterlyTwo:
		return &Period{
			quarterlyTwo: &numberPeriods,
		}, nil
	case SemiAnnually:
		return &Period{
			semiAnnually: &numberPeriods,
		}, nil
	case Annually:
		return &Period{
			annually: &numberPeriods,
		}, nil
	default:
		return &Period{}, nil
	}
}

func (p *Period) getPeriod() (float64, CompoundingFrequency, error) {
	if p.daily != nil {
		return *p.daily, Daily, nil
	}

	if p.monthly != nil {
		return *p.monthly, Monthly, nil
	}

	if p.bimonthly != nil {
		return *p.bimonthly, Bimonthly, nil
	}

	if p.quarterlyOne != nil {
		return *p.quarterlyOne, QuarterlyOne, nil
	}

	if p.quarterlyTwo != nil {
		return *p.quarterlyTwo, QuarterlyTwo, nil
	}

	if p.semiAnnually != nil {
		return *p.semiAnnually, SemiAnnually, nil
	}

	if p.annually != nil {
		return *p.annually, Annually, nil
	}

	return 0, "", errors.New("failed get valid periods")
}

type RateInterest struct {
	value                float64
	compoundingFrequency CompoundingFrequency
	typeRate             TypeRate
}

func NewRateInterest(value float64, compoundingFrequency CompoundingFrequency, typeRate TypeRate) (*RateInterest, error) {
	return &RateInterest{
		value:                value,
		compoundingFrequency: compoundingFrequency,
		typeRate:             typeRate,
	}, nil
}

type CompositeInterest struct {
	future       float64
	present      float64
	rateInterest *RateInterest
	periods      *Period
}

func New(present, future float64, rateInterest *RateInterest, periods *Period) (*CompositeInterest, error) {
	return &CompositeInterest{
		present:      present,
		future:       future,
		rateInterest: rateInterest,
		periods:      periods,
	}, nil
}

func (c *CompositeInterest) getEqualsRateInterestPeriods() (float64, float64, error) {
	valuePeriod, compoundingFrequency, err := c.periods.getPeriod()
	if err != nil {
		return 0, 0, nil
	}
	rateInterest := c.rateInterest.value

	if c.rateInterest.typeRate != RateEffectyPeriodic {
		rateInterest, err = c.rateInterest.RatePeriodic()
		if err != nil {
			return 0, 0, nil
		}
	}

	if compoundingFrequency != c.rateInterest.compoundingFrequency {
		valueInterest, err := getCompoundingFrequency(c.rateInterest.compoundingFrequency)
		if err != nil {
			return 0, 0, nil
		}

		valuePeriod = valuePeriod * valueInterest
	}

	return valuePeriod, rateInterest, nil
}
