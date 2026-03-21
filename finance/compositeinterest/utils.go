package compositeinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/money"
)

// getCompoundingFrequency gets the compounding factor for a given frequency.
// Returns the number of times interest compounds in one year.
//
// Parameters:
//   - compoundingFrequency: The compounding frequency
//
// Returns:
//   - The compounding factor (number of periods per year)
//   - An error if the frequency is invalid
//
// Example:
//
//	factor, err := getCompoundingFrequency(Monthly)
//	// factor is 12 (12 months in a year)
func (cf CompoundingFrequency) getCompoundingFrequency() (money.Decimal, error) {
	periodsPerYear, ok := countCompoundingFrequency[cf]
	if !ok {
		return money.Decimal{}, errors.New("invalid value compounding frequency")
	}

	return periodsPerYear, nil
}

// getCompoundingFrequencytoMonths converts a compounding frequency to the equivalent
// number of months per compounding period.
//
// Returns:
//   - The number of months per compounding period as a Decimal
//   - An error if the frequency is invalid
//
// Example:
//
//	months, _ := getCompoundingFrequencytoMonths(QuarterlyOne)
//	// months is 3 (3 months per quarter)
func (cf CompoundingFrequency) getCompoundingFrequencytoMonths() (money.Decimal, error) {
	monthsPerPeriod, ok := countCompoundingFrequencyMonths[cf]
	if !ok {
		return money.Decimal{}, errors.New("invalid value compounding frequency")
	}

	return monthsPerPeriod, nil
}

// getOrderTime returns the temporal ordering weight for a compounding frequency.
// This is used internally to compare and convert between different frequencies.
//
// Returns:
//   - The order weight as a Decimal
//   - An error if the frequency is invalid
func (cf CompoundingFrequency) getOrderTime() (money.Decimal, error) {
	orderWeight, ok := orderTime[cf]
	if !ok {
		return money.Decimal{}, errors.New("invalid value compounding frequency")
	}

	return orderWeight, nil
}
