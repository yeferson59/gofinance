package compositeinterest

import "errors"

func getCompoundingFrequency(compoundingFrequency CompoundingFrequency) (float64, error) {
	value, ok := countCompoundingFrequency[compoundingFrequency]
	if !ok {
		return 0, errors.New("Invalid value compounding frequency")
	}

	return value, nil
}
