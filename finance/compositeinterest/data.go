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
