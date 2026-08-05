// Package simpleinterest provides calculations for simple interest financial formulas.
// It includes functions to compute future value, present value, interest, rate, and periods
// based on the simple interest formula: Interest = Principal × Rate × Time.
package simpleinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/term"
	"github.com/yeferson59/gofinance/v2/money"
)

// Periods represents the time unit for periods (days, weeks, months, years).
// It is an alias of term.Unit, the shared vocabulary across the finance
// packages.
type Periods = term.Unit

// Period holds the value for different time periods.
// Exactly one of days, weeks, months, or years should be non-zero,
// with the periods field tracking which one is active.
type Period struct {
	days    decimal.Decimal
	weeks   decimal.Decimal
	months  decimal.Decimal
	years   decimal.Decimal
	periods Periods
}

// NewPeriod creates a new Period with the specified number and time unit.
// Valid time units are Days, Weeks, Months, Years.
// Returns an empty Period if timePeriod is invalid.
func NewPeriod(value decimal.Decimal, timePeriod Periods) Period {
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
func (p *Period) getPeriod() (decimal.Decimal, error) {
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
		return decimal.Decimal{}, errors.New("failed to get valid periods")
	}
}

// SimpleInterest holds the values for simple interest calculations.
// Fields are set via New and modified by calculation methods.
type SimpleInterest struct {
	future       money.Money
	present      money.Money
	interest     money.Money
	rateInterest decimal.Decimal
	periods      Period
}

// New creates a new SimpleInterest instance with the provided values.
// Parameters can be 0 if they will be calculated later.
// periods can be an empty Period for calculations that do not need it.
// currency returns the single currency the configured amounts are expressed
// in, ignoring the ones a given calculation does not need.
//
// A simple interest configuration carries three amounts and a caller supplies
// only some of them — a principal and a term, say, with the interest being
// what is asked for. The unset ones are the zero money.Money, which carries no
// currency, so building a result from whichever field a formula happens to
// read produced an amount in XXX whenever that field was the unset one, and
// combining two of them failed with a currency mismatch.
//
// It returns money.ErrCurrencyMismatch when two amounts that are set disagree.
func (s SimpleInterest) currency() (money.Currency, error) {
	return money.ResolveCurrency(s.future, s.present, s.interest)
}

func New(future, present, interest money.Money, rateInterest decimal.Decimal, periods Period) SimpleInterest {
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
func (s SimpleInterest) GetPeriods() (decimal.Decimal, error) {
	periods, err := s.periods.getPeriod()
	return periods, err
}
