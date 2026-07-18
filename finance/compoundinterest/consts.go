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
