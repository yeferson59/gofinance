package simpleinterest

import "errors"

type Periods string

const (
	Days   Periods = "days"
	Weeks  Periods = "weeks"
	Months Periods = "months"
	Years  Periods = "years"
)

type Period struct {
	days   *float64
	weeks  *float64
	months *float64
	years  *float64
}

func NewPeriod(numberPeriods float64, timePeriod Periods) *Period {
	switch timePeriod {
	case Days:
		return &Period{
			days: &numberPeriods,
		}
	case Weeks:
		return &Period{
			weeks: &numberPeriods,
		}
	case Months:
		return &Period{
			months: &numberPeriods,
		}
	case Years:
		return &Period{
			years: &numberPeriods,
		}
	default:
		return &Period{}
	}
}

func (p *Period) getPeriod() (*float64, error) {
	if p.days != nil {
		return p.days, nil
	}

	if p.months != nil {
		return p.months, nil
	}

	if p.weeks != nil {
		return p.weeks, nil
	}

	if p.years != nil {
		return p.years, nil
	}

	return nil, errors.New("failed get valid periods")
}

type InterestSimple struct {
	future       float64
	present      float64
	interest     float64
	rateInterest float64
	periods      *Period
}

func New(future, present, interest, rateInterest float64, periods *Period) *InterestSimple {
	return &InterestSimple{
		future:       future,
		present:      present,
		interest:     interest,
		rateInterest: rateInterest,
		periods:      periods,
	}
}

func (s *InterestSimple) GetPeriods() (float64, error) {
	periods, err := s.periods.getPeriod()
	return *periods, err
}
