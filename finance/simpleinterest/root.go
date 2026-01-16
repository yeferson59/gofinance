// Package simpleinterest provides calculations for simple interest financial formulas.
// It includes functions to compute future value, present value, interest, rate, and periods
// based on the simple interest formula: Interest = Principal × Rate × Time.
package simpleinterest

import "errors"

// Periods represents the time unit for periods (days, weeks, months, years).
type Periods string

// Period holds the value for different time periods.
// Exactly one of days, weeks, months, or years should be non-zero,
// with the periods field tracking which one is active.
type Period struct {
	days   float64
	weeks  float64
	months float64
	years  float64
	periods Periods // Track which period type is active for O(1) lookup
}

// NewPeriod creates a new Period with the specified number and time unit.
// Valid time units are Days, Weeks, Months, Years.
// Returns an empty Period if timePeriod is invalid.
func NewPeriod(numberPeriods float64, timePeriod Periods) Period {
	switch timePeriod {
	case Days:
		return Period{
			days:    numberPeriods,
			periods: Days,
		}
	case Weeks:
		return Period{
			weeks:   numberPeriods,
			periods: Weeks,
		}
	case Months:
		return Period{
			months:  numberPeriods,
			periods: Months,
		}
	case Years:
		return Period{
			years:   numberPeriods,
			periods: Years,
		}
	default:
		return Period{}
	}
}

// getPeriod returns the period value and an error if no valid period is set.
// Uses O(1) lookup via the periods field (a Periods type indicator).
func (p *Period) getPeriod() (float64, error) {
	switch p.periods {
	case Days:
		return p.days, nil
	case Weeks:
		return p.weeks, nil
	case Months:
		return p.months, nil
	case Years:
		return p.years, nil
	default:
		return 0, errors.New("failed get valid periods")
	}
}

// SimpleInterest holds the values for simple interest calculations.
// Fields are set via New and modified by calculation methods.
type SimpleInterest struct {
	future       float64
	present      float64
	interest     float64
	rateInterest float64
	periods      Period
}

// New creates a new SimpleInterest instance with the provided values.
// Parameters can be 0 if they will be calculated later.
// periods can be nil for some calculations.
func New(future, present, interest, rateInterest float64, periods Period) SimpleInterest {
	return SimpleInterest{
		future:       future,
		present:      present,
		interest:     interest,
		rateInterest: rateInterest,
		periods:      periods,
	}
}

// GetPeriods returns the period value from the associated Period.
// Returns an error if periods is invalid.
func (s SimpleInterest) GetPeriods() (float64, error) {
	periods, err := s.periods.getPeriod()
	return periods, err
}
