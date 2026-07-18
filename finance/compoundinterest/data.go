package compoundinterest

import "github.com/yeferson59/gofinance/decimal"

// orderTime assigns each compounding frequency a temporal ordering weight,
// used by the rate conversions to compare frequencies. Periods-per-year and
// months-per-period lookups live on term.Frequency itself
// (PeriodsPerYear, MonthsPerPeriod).
var orderTime = map[CompoundingFrequency]decimal.Decimal{
	Daily:        decimal.MustFromInt64(1, 0),
	Monthly:      decimal.MustFromInt64(2, 0),
	Bimonthly:    decimal.MustFromInt64(3, 0),
	Quarterly:    decimal.MustFromInt64(4, 0),
	FourMonthly:  decimal.MustFromInt64(5, 0),
	SemiAnnually: decimal.MustFromInt64(6, 0),
	Annually:     decimal.MustFromInt64(7, 0),
}
