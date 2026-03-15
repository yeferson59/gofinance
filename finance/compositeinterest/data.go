package compositeinterest

import "github.com/yeferson59/gofinance/money"

// countCompoundingFrequency maps each compounding frequency with the number
// of times it compounds in a year. This map is used to convert between
// different compounding frequencies.
//
// Example:
//   - Daily: 365 (daily compounding)
//   - Monthly: 12 (monthly compounding)
//   - Annually: 1 (annual compounding)
var countCompoundingFrequency = map[CompoundingFrequency]money.Decimal{
	Daily:        money.MustFromInt64(365, 0),
	Monthly:      money.MustFromInt64(12, 0),
	Bimonthly:    money.MustFromInt64(6, 0),
	QuarterlyOne: money.MustFromInt64(4, 0),
	QuarterlyTwo: money.MustFromInt64(3, 0),
	SemiAnnually: money.MustFromInt64(2, 0),
	Annually:     money.MustFromInt64(1, 0),
}

var countCompoundingFrequencyMonths = map[CompoundingFrequency]money.Decimal{
	Daily:        money.MustFromFloat64(0.03333333),
	Monthly:      money.MustFromInt64(1, 0),
	Bimonthly:    money.MustFromInt64(2, 0),
	QuarterlyOne: money.MustFromInt64(3, 0),
	QuarterlyTwo: money.MustFromInt64(4, 0),
	SemiAnnually: money.MustFromInt64(6, 0),
	Annually:     money.MustFromInt64(12, 0),
}

var orderTime = map[CompoundingFrequency]money.Decimal{
	Daily:        money.MustFromInt64(1, 0),
	Monthly:      money.MustFromInt64(2, 0),
	Bimonthly:    money.MustFromInt64(3, 0),
	QuarterlyOne: money.MustFromInt64(4, 0),
	QuarterlyTwo: money.MustFromInt64(5, 0),
	SemiAnnually: money.MustFromInt64(6, 0),
	Annually:     money.MustFromInt64(7, 0),
}
