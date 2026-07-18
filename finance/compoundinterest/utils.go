package compoundinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/decimal"
)

// getOrderTime returns the temporal ordering weight for a compounding
// frequency. This is used internally to compare and convert between
// different frequencies. Periods-per-year and months-per-period conversions
// come from term.Frequency's PeriodsPerYear and MonthsPerPeriod methods.
func getOrderTime(cf CompoundingFrequency) (decimal.Decimal, error) {
	orderWeight, ok := orderTime[cf]
	if !ok {
		return decimal.Decimal{}, errors.New("invalid value compounding frequency")
	}

	return orderWeight, nil
}
