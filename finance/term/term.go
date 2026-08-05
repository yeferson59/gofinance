// Package term defines the shared time vocabulary used across the finance
// packages, so that they interoperate on one set of types instead of each
// declaring its own:
//
//   - Unit is the calendar unit a duration is expressed in (days, weeks,
//     months, years), used by finance/simpleinterest.
//   - Frequency is how many times per year an event occurs — a compounding
//     or payment cadence — used by finance/compoundinterest and
//     finance/annuities.
//
// Day-count conventions (30/360, Actual/365, …) answer a different
// question — how a date range maps to a year fraction — and remain in
// finance/daycount.
package term

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/decimal"
)

// ErrInvalidUnit is returned when a Unit is not one of the declared
// constants.
var ErrInvalidUnit = errors.New("term: invalid time unit")

// ErrInvalidFrequency is returned when a Frequency is not one of the
// declared constants.
var ErrInvalidFrequency = errors.New("term: invalid frequency")

// Unit is the calendar unit a duration is expressed in.
type Unit string

const (
	Days   Unit = "days"
	Weeks  Unit = "weeks"
	Months Unit = "months"
	Years  Unit = "years"
)

// Valid reports whether u is one of the declared unit constants.
func (u Unit) Valid() bool {
	switch u {
	case Days, Weeks, Months, Years:
		return true
	default:
		return false
	}
}

// Frequency is the number of times per year an event occurs, such as a
// compounding or payment cadence.
type Frequency string

const (
	Daily        Frequency = "daily"        // 365 times per year
	Monthly      Frequency = "monthly"      // 12 times per year
	Bimonthly    Frequency = "bimonthly"    // 6 times per year
	Quarterly    Frequency = "quarterly"    // 4 times per year
	FourMonthly  Frequency = "fourMonthly"  // 3 times per year (every four months)
	SemiAnnually Frequency = "semiAnnually" // 2 times per year
	Annually     Frequency = "annually"     // 1 time per year
)

// Valid reports whether f is one of the declared frequency constants.
func (f Frequency) Valid() bool {
	_, err := f.PeriodsPerYear()
	return err == nil
}

// PeriodsPerYear returns how many periods of frequency f fit in one year
// (e.g. 12 for Monthly). It returns ErrInvalidFrequency for an unknown
// frequency.
func (f Frequency) PeriodsPerYear() (decimal.Decimal, error) {
	switch f {
	case Daily:
		return decimal.MustFromInt64(365, 0), nil
	case Monthly:
		return decimal.MustFromInt64(12, 0), nil
	case Bimonthly:
		return decimal.MustFromInt64(6, 0), nil
	case Quarterly:
		return decimal.MustFromInt64(4, 0), nil
	case FourMonthly:
		return decimal.MustFromInt64(3, 0), nil
	case SemiAnnually:
		return decimal.MustFromInt64(2, 0), nil
	case Annually:
		return decimal.MustFromInt64(1, 0), nil
	default:
		return decimal.Decimal{}, ErrInvalidFrequency
	}
}

// dailyMonthsPerPeriod is 12/365, the length of one daily period in months.
// It is derived from Daily's own 365 periods per year rather than written out
// as a constant, so the two methods cannot drift apart.
var dailyMonthsPerPeriod = decimal.MustFromInt64(12, 0).MustDiv(decimal.MustFromInt64(365, 0))

// MonthsPerPeriod returns the length of one period of frequency f expressed
// in months (e.g. 3 for Quarterly). It returns ErrInvalidFrequency for an
// unknown frequency.
//
// The value is always 12 / PeriodsPerYear, so for Daily it is 12/365 ≈
// 0.032876712, not 1/30.
func (f Frequency) MonthsPerPeriod() (decimal.Decimal, error) {
	switch f {
	case Daily:
		return dailyMonthsPerPeriod, nil
	case Monthly:
		return decimal.MustFromInt64(1, 0), nil
	case Bimonthly:
		return decimal.MustFromInt64(2, 0), nil
	case Quarterly:
		return decimal.MustFromInt64(3, 0), nil
	case FourMonthly:
		return decimal.MustFromInt64(4, 0), nil
	case SemiAnnually:
		return decimal.MustFromInt64(6, 0), nil
	case Annually:
		return decimal.MustFromInt64(12, 0), nil
	default:
		return decimal.Decimal{}, ErrInvalidFrequency
	}
}
