package compositeinterest

// countCompoundingFrequency maps each compounding frequency with the number
// of times it compounds in a year. This map is used to convert between
// different compounding frequencies.
//
// Example:
//   - Daily: 365 (daily compounding)
//   - Monthly: 12 (monthly compounding)
//   - Annually: 1 (annual compounding)
var countCompoundingFrequency = map[CompoundingFrequency]float64{
	Daily:        365,
	Monthly:      12,
	Bimonthly:    6,
	QuarterlyOne: 4,
	QuarterlyTwo: 3,
	SemiAnnually: 2,
	Annually:     1,
}

var countCompoundingFrequencyMonths = map[CompoundingFrequency]float64{
	Daily:        0.03333333,
	Monthly:      1,
	Bimonthly:    2,
	QuarterlyOne: 3,
	QuarterlyTwo: 4,
	SemiAnnually: 6,
	Annually:     12,
}

var orderTime = map[CompoundingFrequency]float64{
	Daily:        1,
	Monthly:      2,
	Bimonthly:    3,
	QuarterlyOne: 4,
	QuarterlyTwo: 5,
	SemiAnnually: 6,
	Annually:     7,
}
