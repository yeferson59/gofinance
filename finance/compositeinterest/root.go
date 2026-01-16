// Package compositeinterest provides functionality for compound interest calculations.
//
// This package enables financial calculations related to compound interest,
// including:
//   - Calculation of Future Value (FV)
//   - Calculation of Present Value (PV)
//   - Calculation of periodic interest rate
//   - Calculation of the number of periods
//   - Conversion between different types of rates (periodic, nominal, effective annual)
//   - Handling of anticipated rates (discount)
//
// The package supports multiple compounding frequencies:
//   - Daily (365 periods per year)
//   - Monthly (12 periods per year)
//   - Bimonthly (6 periods per year)
//   - Quarterly (4 or 3 periods per year)
//   - Semi-annually (2 periods per year)
//   - Annually (1 period per year)
//
// Basic usage example:
//
//	// Create a monthly interest rate of 1% periodic
//	rateInterest, _ := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
//
//	// Create a period of 12 months
//	period, _ := NewPeriod(12, Monthly)
//
//	// Create a compound interest object with $1000 initial capital
//	ci, _ := New(1000, 0, rateInterest, period)
//
//	// Calculate the future value
//	future, _ := ci.Future()
//	// future ≈ 1126.83
package compositeinterest

import (
	"errors"
	"math"
)

// CompoundingFrequency defines the frequency of interest compounding in a year.
// Valid values are: Daily, Monthly, Bimonthly, QuarterlyOne, QuarterlyTwo,
// SemiAnnually, Annually.
type CompoundingFrequency string

// TypeRate defines the type of interest rate to use in calculations.
// Valid values include discount rates (anticipated) and ordinary rates.
// Examples: RateEffectyPeriodic, RateEffectyNominal, RateEffectyAnnually,
// RateAnticipateEffectyPeriodic, RateAnticipateEffectyNominal, RateAnticipateEffectyAnnually.
type TypeRate string

// Period represents the number of compounding periods for a compound interest calculation.
// It stores a single period value along with its compounding frequency.
type Period struct {
	value     float64
	frequency CompoundingFrequency
}

// NewPeriod creates a new Period instance for the specified compounding frequency.
//
// Parameters:
//   - numberPeriods: The number of periods (e.g., 12 for 12 months if frequency is Monthly)
//   - compoundingFrequency: The compounding frequency (Daily, Monthly, etc.)
//
// Returns:
//   - A Period instance with the specified period value and frequency
//   - An error if the value is negative or the frequency is invalid
//
// Example:
//
//	period, err := NewPeriod(12, Monthly)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewPeriod(value float64, compoundingFrequency CompoundingFrequency) (Period, error) {
	if value < 0 {
		return Period{}, errors.New("value periods must be greater or equal to zero")
	}

	switch compoundingFrequency {
	case Daily, Monthly, Bimonthly, QuarterlyOne, QuarterlyTwo, SemiAnnually, Annually:
		return Period{
			value:     value,
			frequency: compoundingFrequency,
		}, nil
	default:
		return Period{}, errors.New("invalid compounding frequency")
	}
}

// getPeriod extracts the period value and its associated frequency from the Period structure.
//
// Returns:
//   - The numeric value of the period
//   - The corresponding compounding frequency
//   - An error if the frequency is invalid or uninitialized
func (p *Period) getPeriod() (float64, CompoundingFrequency, error) {
	// Direct lookup via frequency field
	switch p.frequency {
	case Daily, Monthly, Bimonthly, QuarterlyOne, QuarterlyTwo, SemiAnnually, Annually:
		return p.value, p.frequency, nil
	default:
		return 0, "", errors.New("failed to get valid periods")
	}
}

// RateInterest represents an interest rate with its compounding frequency and type.
// Fields:
//   - value: The value of the rate (as decimal, e.g., 0.05 for 5%)
//   - compoundingFrequency: The frequency with which interest is compounded
//   - typeRate: The type of rate (periodic, nominal, effective annual, etc.)
type RateInterest struct {
	value                float64
	compoundingFrequency CompoundingFrequency
	typeRate             TypeRate
}

// NewRateInterest creates a new RateInterest instance.
//
// Parameters:
//   - value: The numeric value of the rate (e.g., 0.05 for 5%)
//   - compoundingFrequency: The compounding frequency
//   - typeRate: The type of rate to use
//
// Returns:
//   - A RateInterest instance
//   - An error if parameters are invalid
//
// Example:
//
//	rate, err := NewRateInterest(0.12, Monthly, RateEffectyNominal)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewRateInterest(value float64, compoundingFrequency CompoundingFrequency, typeRate TypeRate) (RateInterest, error) {
	if value < 0 {
		return RateInterest{}, errors.New("invalid value for rate interest must be greater o equal to zero")
	}

	return RateInterest{
		value:                value,
		compoundingFrequency: compoundingFrequency,
		typeRate:             typeRate,
	}, nil
}

// CompositeInterest contains all the necessary parameters for compound interest calculations.
// Fields:
//   - future: The future value (if known)
//   - present: The present value or initial capital
//   - rateInterest: The interest rate to apply
//   - periods: The number of compounding periods
//
// Use the methods Future(), Present(), Interest() and Periods() to calculate unknown values.
type CompositeInterest struct {
	future       float64
	present      float64
	rateInterest RateInterest
	periods      Period
}

// New creates a new CompositeInterest instance with the specified parameters.
//
// Parameters:
//   - present: The present value or initial capital (0 if unknown)
//   - future: The future value (0 if unknown)
//   - rateInterest: The interest rate to apply
//   - periods: The number of periods
//
// Note: You must provide at least three of the four values (present, future, rateInterest, periods).
// The fourth value will be calculated using the corresponding method.
//
// Returns:
//   - A CompositeInterest instance
//   - An error if parameters are invalid
//
// Example:
//
//	ci, err := New(1000, 0, rateInterest, period)
//	if err != nil {
//	    log.Fatal(err)
//	}
func New(present, future float64, rateInterest RateInterest, periods Period) (CompositeInterest, error) {
	return CompositeInterest{
		present:      present,
		future:       future,
		rateInterest: rateInterest,
		periods:      periods,
	}, nil
}

// GetEqualsRateInterestPeriods converts all parameters to the same base for calculations.
// Specifically:
//   - Converts the rate to periodic if necessary
//   - Adjusts the periods if the compounding frequency does not match the rate
//
// Returns:
//   - The adjusted number of periods
//   - The equivalent periodic rate
//   - An error if valid values cannot be obtained
func (c CompositeInterest) GetEqualsRateInterestPeriods() (float64, float64, error) {
	periodValue, compoundingFrequency, err := c.periods.getPeriod()
	if err != nil {
		return 0, 0, nil
	}

	periodicRate := c.rateInterest.value

	if c.rateInterest.typeRate != RateEffectyPeriodic {
		periodicRate, err = c.rateInterest.RatePeriodic()
		if err != nil {
			return 0, 0, nil
		}
	}

	if compoundingFrequency != c.rateInterest.compoundingFrequency {
		periodsInMonths, err := getCompoundingFrequencytoMonths(compoundingFrequency)
		if err != nil {
			return 0, 0, err
		}

		rateFrequencyInMonths, err := getCompoundingFrequencytoMonths(c.rateInterest.compoundingFrequency)
		if err != nil {
			return 0, 0, err
		}

		periodWeight, err := getOrderTime(compoundingFrequency)
		if err != nil {
			return 0, 0, err
		}

		rateWeight, err := getOrderTime(c.rateInterest.compoundingFrequency)
		if err != nil {
			return 0, 0, err
		}

		if rateWeight > periodWeight {
			if periodsInMonths < 1 {
				return math.Round(periodsInMonths * periodValue), periodicRate, nil
			}

			return (periodValue / rateFrequencyInMonths), periodicRate, nil
		}

		return (periodsInMonths * periodValue), periodicRate, nil
	}

	return periodValue, periodicRate, nil
}
