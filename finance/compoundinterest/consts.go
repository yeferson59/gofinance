package compoundinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/finance/term"
)

// Compounding frequencies of interest per year, re-exported from the shared
// finance/term vocabulary. These constants define how many times interest
// compounds in a one-year period.
const (
	Daily        = term.Daily        // 365 times per year
	Monthly      = term.Monthly      // 12 times per year
	Bimonthly    = term.Bimonthly    // 6 times per year
	Quarterly    = term.Quarterly    // 4 times per year
	FourMonthly  = term.FourMonthly  // 3 times per year (every four months)
	SemiAnnually = term.SemiAnnually // 2 times per year
	Annually     = term.Annually     // 1 time per year
)

// QuarterlyOne is the historical name for quarterly compounding.
//
// Deprecated: use Quarterly (4 compounding periods per year).
const QuarterlyOne = term.Quarterly

// QuarterlyTwo is the historical name for every-four-months compounding.
//
// Deprecated: use FourMonthly (3 compounding periods per year).
const QuarterlyTwo = term.FourMonthly

// Types of interest rates.
// Ordinary rates (charged at the end of the period):
//   - RateEffectivePeriodic: Periodic rate (the one that compounds each period)
//   - RateEffectiveNominal: Nominal annual rate
//   - RateEffectiveAnnually: Effective annual rate
//
// Anticipated/discount rates (charged at the beginning of the period):
//   - RateAnticipateEffectivePeriodic: Anticipated periodic rate
//   - RateAnticipateEffectiveNominal: Anticipated nominal rate
const (
	RateEffectivePeriodic           TypeRate = "periodic"
	RateEffectiveAnnually           TypeRate = "annual"
	RateEffectiveNominal            TypeRate = "nominal"
	RateAnticipateEffectivePeriodic TypeRate = "anticipatePeriodic"
	RateAnticipateEffectiveNominal  TypeRate = "anticipateNominal"
)

// The original, misspelled names of the rate types, kept so existing code
// compiles. They are the same values as the correctly spelled constants.
const (
	// Deprecated: use RateEffectivePeriodic.
	RateEffectyPeriodic = RateEffectivePeriodic
	// Deprecated: use RateEffectiveAnnually.
	RateEffectyAnnually = RateEffectiveAnnually
	// Deprecated: use RateEffectiveNominal.
	RateEffectyNominal = RateEffectiveNominal
	// Deprecated: use RateAnticipateEffectivePeriodic.
	RateAnticipateEffectyPeriodic = RateAnticipateEffectivePeriodic
	// Deprecated: use RateAnticipateEffectiveNominal.
	RateAnticipateEffectyNominal = RateAnticipateEffectiveNominal
)

var (
	ErrInvalidOperation = errors.New("invalid values for operation")

	// ErrInvalidTypeRate is returned by the rate conversions when the rate's
	// TypeRate is not one of the declared constants, so no conversion is
	// defined for it.
	ErrInvalidTypeRate = errors.New("compoundinterest: invalid rate type")

	// ErrInvalidAnticipatedRate is returned when an anticipated (discount)
	// periodic rate is 100% or more. Such a rate would discount a period's
	// whole value or beyond, so it has no finite ordinary equivalent.
	ErrInvalidAnticipatedRate = errors.New("compoundinterest: anticipated periodic rate must be less than 1")
)
