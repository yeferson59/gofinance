package depreciation

import "github.com/yeferson59/gofinance/v2/decimal"

// macrsTable returns the MACRS GDS half-year-convention depreciation rates (as
// fractions of the original cost) for the given recovery period, and whether a
// table exists for it. The rates are the IRS-published percentages; the final
// year is adjusted by MACRS so the schedule fully recovers the cost.
func macrsTable(recovery int) ([]decimal.Decimal, bool) {
	strs, ok := macrsPercents[recovery]
	if !ok {
		return nil, false
	}

	rates := make([]decimal.Decimal, len(strs))
	for i, s := range strs {
		rates[i] = decimal.MustFromString(s)
	}

	return rates, true
}

// macrsPercents holds the GDS half-year-convention rates as exact decimal
// fractions, keyed by recovery period in years.
var macrsPercents = map[int][]string{
	3: {"0.3333", "0.4445", "0.1481", "0.0741"},
	5: {"0.2000", "0.3200", "0.1920", "0.1152", "0.1152", "0.0576"},
	7: {"0.1429", "0.2449", "0.1749", "0.1249", "0.0893", "0.0892", "0.0893", "0.0446"},
	10: {
		"0.1000", "0.1800", "0.1440", "0.1152", "0.0922", "0.0737",
		"0.0655", "0.0655", "0.0656", "0.0655", "0.0328",
	},
	15: {
		"0.0500", "0.0950", "0.0855", "0.0770", "0.0693", "0.0623",
		"0.0590", "0.0590", "0.0591", "0.0590", "0.0591", "0.0590",
		"0.0591", "0.0590", "0.0591", "0.0295",
	},
	20: {
		"0.03750", "0.07219", "0.06677", "0.06177", "0.05713", "0.05285",
		"0.04888", "0.04522", "0.04462", "0.04461", "0.04462", "0.04461",
		"0.04462", "0.04461", "0.04462", "0.04461", "0.04462", "0.04461",
		"0.04462", "0.04461", "0.02231",
	},
}
