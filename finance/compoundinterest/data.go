package compoundinterest

import "github.com/yeferson59/gofinance/decimal"

// countCompoundingFrequency maps each compounding frequency with the number
// of times it compounds in a year. This map is used to convert between
// different compounding frequencies.
//
// Example:
//   - Daily: 365 (daily compounding)
//   - Monthly: 12 (monthly compounding)
//   - Annually: 1 (annual compounding)
var countCompoundingFrequency = map[CompoundingFrequency]decimal.Decimal{
	Daily:        decimal.MustFromInt64(365, 0),
	Monthly:      decimal.MustFromInt64(12, 0),
	Bimonthly:    decimal.MustFromInt64(6, 0),
	QuarterlyOne: decimal.MustFromInt64(4, 0),
	QuarterlyTwo: decimal.MustFromInt64(3, 0),
	SemiAnnually: decimal.MustFromInt64(2, 0),
	Annually:     decimal.MustFromInt64(1, 0),
}

var countCompoundingFrequencyMonths = map[CompoundingFrequency]decimal.Decimal{
	Daily:        decimal.MustFromFloat64(0.03333333),
	Monthly:      decimal.MustFromInt64(1, 0),
	Bimonthly:    decimal.MustFromInt64(2, 0),
	QuarterlyOne: decimal.MustFromInt64(3, 0),
	QuarterlyTwo: decimal.MustFromInt64(4, 0),
	SemiAnnually: decimal.MustFromInt64(6, 0),
	Annually:     decimal.MustFromInt64(12, 0),
}

var orderTime = map[CompoundingFrequency]decimal.Decimal{
	Daily:        decimal.MustFromInt64(1, 0),
	Monthly:      decimal.MustFromInt64(2, 0),
	Bimonthly:    decimal.MustFromInt64(3, 0),
	QuarterlyOne: decimal.MustFromInt64(4, 0),
	QuarterlyTwo: decimal.MustFromInt64(5, 0),
	SemiAnnually: decimal.MustFromInt64(6, 0),
	Annually:     decimal.MustFromInt64(7, 0),
}
