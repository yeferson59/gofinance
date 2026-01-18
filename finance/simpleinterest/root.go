// Package simpleinterest provides calculations for simple interest financial formulas.
// It includes functions to compute future value, present value, interest, rate, and periods
// based on the simple interest formula: Interest = Principal × Rate × Time.
package simpleinterest

import (
	"errors"

	"github.com/quagmt/udecimal"
)

// Periods represents the time unit for periods (days, weeks, months, years).
type Periods string

// Period holds the value for different time periods.
// Exactly one of days, weeks, months, or years should be non-zero,
// with the periods field tracking which one is active.
type Period struct {
	days    Decimal
	weeks   Decimal
	months  Decimal
	years   Decimal
	periods Periods
}

// NewPeriod creates a new Period with the specified number and time unit.
// Valid time units are Days, Weeks, Months, Years.
// Returns an empty Period if timePeriod is invalid.
func NewPeriod(value Decimal, timePeriod Periods) Period {
	switch timePeriod {
	case Days:
		return Period{
			days:    value,
			periods: Days,
		}
	case Weeks:
		return Period{
			weeks:   value,
			periods: Weeks,
		}
	case Months:
		return Period{
			months:  value,
			periods: Months,
		}
	case Years:
		return Period{
			years:   value,
			periods: Years,
		}
	default:
		return Period{}
	}
}

// getPeriod returns the period value and an error if no valid period is set.
// Uses O(1) lookup via the periods field (a Periods type indicator).
func (p *Period) getPeriod() (Decimal, error) {
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
		return Decimal{}, errors.New("failed get valid periods")
	}
}

// SimpleInterest holds the values for simple interest calculations.
// Fields are set via New and modified by calculation methods.
type SimpleInterest struct {
	future       Decimal
	present      Decimal
	interest     Decimal
	rateInterest Decimal
	periods      Period
}

// New creates a new SimpleInterest instance with the provided values.
// Parameters can be 0 if they will be calculated later.
// periods can be nil for some calculations.
func New(future, present, interest, rateInterest Decimal, periods Period) SimpleInterest {
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
func (s SimpleInterest) GetPeriods() (Decimal, error) {
	periods, err := s.periods.getPeriod()
	return periods, err
}

type Decimal struct {
	udecimal.Decimal
}

func NewFromFloat64(f float64) (Decimal, error) {
	decimal, err := udecimal.NewFromFloat64(f)
	return Decimal{
		decimal,
	}, err
}

func NewFromInt64(coef int64, prec uint8) (Decimal, error) {
	decimal, err := udecimal.NewFromInt64(coef, prec)
	return Decimal{
		decimal,
	}, err
}
