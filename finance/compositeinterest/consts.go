package compositeinterest

const (
	Daily        CompoundingFrequency = "daily"
	Monthly      CompoundingFrequency = "monthly"
	Bimonthly    CompoundingFrequency = "Bimontly"
	QuarterlyOne CompoundingFrequency = "quarterlyOne"
	QuarterlyTwo CompoundingFrequency = "quarterlyTwo"
	SemiAnnually CompoundingFrequency = "semiAnnually"
	Annually     CompoundingFrequency = "annually"
)

const (
	RateEffectyPeriodic           TypeRate = "periodic"
	RateEffectyAnnually           TypeRate = "annual"
	RateEffectyNominal            TypeRate = "nominal"
	RateAnticipateEffectyPeriodic TypeRate = "anticipatePeriodic"
	RateAnticipateEffectyAnnually TypeRate = "anticipateAnnual"
	RateAnticipateEffectyNominal  TypeRate = "anticipateNominal"
)
