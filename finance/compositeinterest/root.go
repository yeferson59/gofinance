package compositeinterest

type CompoundingFrequency string
type TypeRate string

const (
	Daily        CompoundingFrequency = "daily"
	Monthly      CompoundingFrequency = "monthly"
	Bimonthly    CompoundingFrequency = "Bimontly"
	QuarterlyOne CompoundingFrequency = "quarterlyOne"
	QuarterlyTwo CompoundingFrequency = "quarterlyTwo"
	SemiAnnually CompoundingFrequency = "semiAnnually"
	Annually     CompoundingFrequency = "annualy"
)

const (
	RateEffectyPeriodic TypeRate = "periodic"
	RateEffectyAnnualy  TypeRate = "annual"
	RateEffectyNominal  TypeRate = "nominal"
)

var countCompoundingFrequency = map[CompoundingFrequency]float64{
	Daily:        365,
	Monthly:      12,
	Bimonthly:    6,
	QuarterlyOne: 4,
	QuarterlyTwo: 3,
	SemiAnnually: 2,
	Annually:     1,
}

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
