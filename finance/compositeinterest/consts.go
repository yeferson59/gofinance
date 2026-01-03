package compositeinterest

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
	RateEffectyPeriodic           TypeRate = "periodic"
	RateEffectyAnnualy            TypeRate = "annual"
	RateEffectyNominal            TypeRate = "nominal"
	RateAnticipateEffectyPeriodic TypeRate = "anticipatePeriodic"
	RateAnticipateEffectyAnnualy  TypeRate = "anticipateAnnual"
	RateAnticipateEffectyNominal  TypeRate = "anticipateNominal"
)
