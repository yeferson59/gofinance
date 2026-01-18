package compositeinterest

import "errors"

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
func (cf CompoundingFrequency) getCompoundingFrequency() (float64, error) {
	periodsPerYear, ok := countCompoundingFrequency[cf]
	if !ok {
		return 0, errors.New("invalid value compounding frequency")
	}

	return periodsPerYear, nil
}

func (cf CompoundingFrequency) getCompoundingFrequencytoMonths() (float64, error) {
	monthsPerPeriod, ok := countCompoundingFrequencyMonths[cf]
	if !ok {
		return 0, errors.New("invalid value compounding frequency")
	}

	return monthsPerPeriod, nil
}

func (cf CompoundingFrequency) getOrderTime() (float64, error) {
	orderWeight, ok := orderTime[cf]
	if !ok {
		return 0, errors.New("invalid value compounding frequency")
	}

	return orderWeight, nil
}
