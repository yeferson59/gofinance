package compositeinterest

import "errors"

// Compounding frequencies of interest per year.
// These constants define how many times interest compounds in a one-year period.
const (
	Daily        CompoundingFrequency = "daily"        // 365 times per year
	Monthly      CompoundingFrequency = "monthly"      // 12 times per year
	Bimonthly    CompoundingFrequency = "bimonthly"    // 6 times per year
	QuarterlyOne CompoundingFrequency = "quarterlyOne" // 4 times per year
	QuarterlyTwo CompoundingFrequency = "quarterlyTwo" // 3 times per year
	SemiAnnually CompoundingFrequency = "semiAnnually" // 2 times per year
	Annually     CompoundingFrequency = "annually"     // 1 time per year
)

// Types of interest rates.
// Ordinary rates (charged at the end of the period):
//   - RateEffectyPeriodic: Periodic rate (the one that compounds each period)
//   - RateEffectyNominal: Nominal annual rate
//   - RateEffectyAnnually: Effective annual rate
//
// Anticipated/discount rates (charged at the beginning of the period):
//   - RateAnticipateEffectyPeriodic: Anticipated periodic rate
//   - RateAnticipateEffectyNominal: Anticipated nominal rate
const (
	RateEffectyPeriodic           TypeRate = "periodic"
	RateEffectyAnnually           TypeRate = "annual"
	RateEffectyNominal            TypeRate = "nominal"
	RateAnticipateEffectyPeriodic TypeRate = "anticipatePeriodic"
	RateAnticipateEffectyNominal  TypeRate = "anticipateNominal"
)

var (
	ErrInvalidOperation = errors.New("invalid values for operation")
)
