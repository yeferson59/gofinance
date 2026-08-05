package term

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
)

var allFrequencies = []Frequency{
	Daily, Monthly, Bimonthly, Quarterly, FourMonthly, SemiAnnually, Annually,
}

var allUnits = []Unit{Days, Weeks, Months, Years}

// TestMonthsPerPeriodMatchesPeriodsPerYear is the regression test for
// TESTING_PLAN.md §2.4. The two methods describe the same thing from opposite
// ends, so they must satisfy
//
//	MonthsPerPeriod × PeriodsPerYear = 12
//
// for every frequency. Daily broke it: PeriodsPerYear returned 365 while
// MonthsPerPeriod returned a hand-written 0.03333333 (1/30, a month of thirty
// days), a 1.4% discrepancy that leaked into any calculation mixing the two.
// The existing table test happened to omit Daily, so nothing caught it.
func TestMonthsPerPeriodMatchesPeriodsPerYear(t *testing.T) {
	twelve := decimal.MustFromInt64(12, 0)

	for _, frequency := range allFrequencies {
		t.Run(string(frequency), func(t *testing.T) {
			periodsPerYear, err := frequency.PeriodsPerYear()
			require.NoError(t, err)

			monthsPerPeriod, err := frequency.MonthsPerPeriod()
			require.NoError(t, err)

			product := monthsPerPeriod.Mul(periodsPerYear)
			assert.InDelta(t, twelve.InexactFloat64(), product.InexactFloat64(), 1e-15,
				"months per period × periods per year must be 12")
		})
	}
}

// TestDailyMonthsPerPeriod pins the corrected value.
func TestDailyMonthsPerPeriod(t *testing.T) {
	months, err := Daily.MonthsPerPeriod()
	require.NoError(t, err)

	// 12/365 = 0.0328767123287671232876712328767…
	assert.InDelta(t, 12.0/365.0, months.InexactFloat64(), 1e-15)
}

// TestFrequencyValidCoversEveryConstant checks Valid agrees with
// PeriodsPerYear across the declared set and rejects anything else.
func TestFrequencyValidCoversEveryConstant(t *testing.T) {
	for _, frequency := range allFrequencies {
		assert.True(t, frequency.Valid(), "%q must be valid", frequency)
	}

	for _, invalid := range []Frequency{"", "weekly", "MONTHLY", "annual"} {
		assert.False(t, invalid.Valid(), "%q must be invalid", invalid)

		_, err := invalid.PeriodsPerYear()
		assert.ErrorIs(t, err, ErrInvalidFrequency)

		_, err = invalid.MonthsPerPeriod()
		assert.ErrorIs(t, err, ErrInvalidFrequency)
	}
}

// TestUnitValidCoversEveryConstant checks the same for the calendar units.
func TestUnitValidCoversEveryConstant(t *testing.T) {
	for _, unit := range allUnits {
		assert.True(t, unit.Valid(), "%q must be valid", unit)
	}

	for _, invalid := range []Unit{"", "day", "DAYS", "fortnights"} {
		assert.False(t, invalid.Valid(), "%q must be invalid", invalid)
	}
}

// TestPeriodsPerYearOrdering checks the frequencies are ordered as their names
// claim: more frequent compounding means more periods and shorter ones.
func TestPeriodsPerYearOrdering(t *testing.T) {
	ordered := []Frequency{Annually, SemiAnnually, FourMonthly, Quarterly, Bimonthly, Monthly, Daily}

	previousPeriods := 0.0

	for i, frequency := range ordered {
		periodsPerYear, err := frequency.PeriodsPerYear()
		require.NoError(t, err)

		if i > 0 {
			assert.Greater(t, periodsPerYear.InexactFloat64(), previousPeriods,
				"%q must have more periods per year than the previous frequency", frequency)
		}

		previousPeriods = periodsPerYear.InexactFloat64()
	}
}
