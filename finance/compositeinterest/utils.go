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
func getCompoundingFrequency(compoundingFrequency CompoundingFrequency) (float64, error) {
	periodsPerYear, ok := countCompoundingFrequency[compoundingFrequency]
	if !ok {
		return 0, errors.New("invalid value compounding frequency")
	}

	return periodsPerYear, nil
}

func getCompoundingFrequencytoMonths(compoundingFrequency CompoundingFrequency) (float64, error) {
	monthsPerPeriod, ok := countCompoundingFrequencyMonths[compoundingFrequency]
	if !ok {
		return 0, errors.New("invalid value compounding frequency")
	}

	return monthsPerPeriod, nil
}

func getOrderTime(compoundingFrequency CompoundingFrequency) (float64, error) {
	orderWeight, ok := orderTime[compoundingFrequency]
	if !ok {
		return 0, errors.New("invalid value compounding frequency")
	}

	return orderWeight, nil
}
