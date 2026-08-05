package daycount

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var allConventions = []Convention{Thirty360, Actual360, Actual365Fixed, ActualActualISDA}

// TestThirty360EndOfFebruary is the regression test for TESTING_PLAN.md §2.5.
// The implementation documented "the standard end-of-month adjustments" but
// applied only the 31st rule, so a period starting at the end of February
// measured 182 days where the US (NASD) convention and spreadsheet DAYS360
// measure 180. The February rules are now applied.
func TestThirty360EndOfFebruary(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		// Rule 1: a start at the end of February counts as the 30th, so six
		// whole months measure 180 days.
		{"leap Feb 29 to Aug 31", date(2024, time.February, 29), date(2024, time.August, 31), 180},
		{"common Feb 28 to Aug 31", date(2023, time.February, 28), date(2023, time.August, 31), 180},
		{"leap Feb 29 to Aug 30", date(2024, time.February, 29), date(2024, time.August, 30), 180},

		// Rule 2: both ends at the end of February measure a whole year.
		{"Feb 28 to Feb 28", date(2023, time.February, 28), date(2024, time.February, 29), 360},
		{"Feb 29 to Feb 28", date(2024, time.February, 29), date(2025, time.February, 28), 360},

		// A start before the month end is untouched by the February rules.
		{"Feb 15 to Aug 15", date(2024, time.February, 15), date(2024, time.August, 15), 180},
		{"Feb 27 to Aug 31", date(2024, time.February, 27), date(2024, time.August, 31), 184},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			days, err := Days(test.start, test.end, Thirty360)
			require.NoError(t, err)
			assert.Equal(t, test.expected, days)
		})
	}
}

// TestThirty360DayThirtyOneRules covers the two adjustments that were already
// implemented, so the February work cannot regress them.
func TestThirty360DayThirtyOneRules(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		// Rule 3: a start on the 31st counts as the 30th.
		{"Jan 31 to Mar 31", date(2024, time.January, 31), date(2024, time.March, 31), 60},
		{"Jan 30 to Mar 31", date(2024, time.January, 30), date(2024, time.March, 31), 60},
		// Rule 4 does not apply when the start is before the 30th, so the end
		// keeps its 31st.
		{"Jan 15 to Mar 31", date(2024, time.January, 15), date(2024, time.March, 31), 76},
		{"Jan 1 to Jan 31", date(2024, time.January, 1), date(2024, time.January, 31), 30},
		{"whole year", date(2024, time.January, 1), date(2025, time.January, 1), 360},
		{"same day", date(2024, time.June, 15), date(2024, time.June, 15), 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			days, err := Days(test.start, test.end, Thirty360)
			require.NoError(t, err)
			assert.Equal(t, test.expected, days)
		})
	}
}

// TestActualConventionsCountRealDays checks the Actual conventions count
// calendar days and differ only in their denominator.
func TestActualConventionsCountRealDays(t *testing.T) {
	start := date(2024, time.January, 1)
	end := date(2024, time.July, 1)

	// 2024 is a leap year: Jan(31)+Feb(29)+Mar(31)+Apr(30)+May(31)+Jun(30)=182
	for _, convention := range []Convention{Actual360, Actual365Fixed, ActualActualISDA} {
		days, err := Days(start, end, convention)
		require.NoError(t, err)
		assert.Equal(t, 182, days, "%v", convention)
	}

	over360, err := YearFraction(start, end, Actual360)
	require.NoError(t, err)
	assert.InDelta(t, 182.0/360.0, over360.InexactFloat64(), 1e-12)

	over365, err := YearFraction(start, end, Actual365Fixed)
	require.NoError(t, err)
	assert.InDelta(t, 182.0/365.0, over365.InexactFloat64(), 1e-12)

	// Actual/Actual ISDA divides by the year's real length.
	isda, err := YearFraction(start, end, ActualActualISDA)
	require.NoError(t, err)
	assert.InDelta(t, 182.0/366.0, isda.InexactFloat64(), 1e-12)
}

// TestActualActualISDASpansYears checks the convention's defining behaviour:
// the period is split at calendar-year boundaries and each part divided by its
// own year length.
func TestActualActualISDASpansYears(t *testing.T) {
	// 1 Dec 2023 to 1 Feb 2024: 31 days in 2023 (365) + 31 days in 2024 (366).
	fraction, err := YearFraction(date(2023, time.December, 1), date(2024, time.February, 1), ActualActualISDA)
	require.NoError(t, err)
	assert.InDelta(t, 31.0/365.0+31.0/366.0, fraction.InexactFloat64(), 1e-12)

	// A whole leap year is exactly one.
	whole, err := YearFraction(date(2024, time.January, 1), date(2025, time.January, 1), ActualActualISDA)
	require.NoError(t, err)
	assert.InDelta(t, 1.0, whole.InexactFloat64(), 1e-12)

	// Several years, spanning two leap years, with stub periods at both ends.
	multi, err := YearFraction(date(2023, time.July, 1), date(2027, time.March, 1), ActualActualISDA)
	require.NoError(t, err)
	// 184/365 (2023 stub) + 1 + 1 + 1 (2024, 2025, 2026) + 59/365 (2027 stub)
	assert.InDelta(t, 184.0/365.0+3+59.0/365.0, multi.InexactFloat64(), 1e-12)
}

// TestLeapYearBoundaries checks the leap rule at the century marks, where the
// divisible-by-100 and divisible-by-400 cases disagree.
func TestLeapYearBoundaries(t *testing.T) {
	// 2000 is a leap year (divisible by 400), 1900 and 2100 are not.
	leap, err := YearFraction(date(2000, time.January, 1), date(2001, time.January, 1), ActualActualISDA)
	require.NoError(t, err)
	assert.InDelta(t, 1.0, leap.InexactFloat64(), 1e-12)

	days, err := Days(date(2000, time.January, 1), date(2001, time.January, 1), Actual365Fixed)
	require.NoError(t, err)
	assert.Equal(t, 366, days)

	common, err := Days(date(2100, time.January, 1), date(2101, time.January, 1), Actual365Fixed)
	require.NoError(t, err)
	assert.Equal(t, 365, common)
}

// TestNormalizationIgnoresClockAndZone checks that the time of day and the
// caller's location do not shift a day count: dates are reduced to UTC
// midnight first.
func TestNormalizationIgnoresClockAndZone(t *testing.T) {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	newYork, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)

	// The same calendar dates, expressed in three zones and at odd hours.
	start := time.Date(2024, time.March, 1, 23, 45, 0, 0, tokyo)
	end := time.Date(2024, time.September, 1, 3, 15, 0, 0, newYork)

	for _, convention := range allConventions {
		fromZoned, err := Days(start, end, convention)
		require.NoError(t, err)

		fromUTC, err := Days(date(2024, time.March, 1), date(2024, time.September, 1), convention)
		require.NoError(t, err)

		assert.Equal(t, fromUTC, fromZoned, "%v must ignore clock time and zone", convention)
	}
}

// TestEndBeforeStartIsRejected checks both entry points refuse a reversed
// range rather than returning a negative count.
func TestEndBeforeStartIsRejected(t *testing.T) {
	later := date(2024, time.June, 1)
	earlier := date(2024, time.January, 1)

	for _, convention := range allConventions {
		_, err := Days(later, earlier, convention)
		assert.ErrorIs(t, err, ErrEndBeforeStart)

		_, err = YearFraction(later, earlier, convention)
		assert.ErrorIs(t, err, ErrEndBeforeStart)
	}
}

// TestInvalidConventionIsRejected checks an unrecognised convention is
// reported instead of silently falling through to a default.
func TestInvalidConventionIsRejected(t *testing.T) {
	start := date(2024, time.January, 1)
	end := date(2024, time.June, 1)

	for _, invalid := range []Convention{Convention(-1), Convention(4), Convention(99)} {
		_, err := Days(start, end, invalid)
		assert.ErrorIs(t, err, ErrInvalidConvention)

		_, err = YearFraction(start, end, invalid)
		assert.ErrorIs(t, err, ErrInvalidConvention)
	}
}

// TestConventionStringIsDistinct checks every declared convention has a
// distinct name and an unknown one is labelled rather than printed blank.
func TestConventionStringIsDistinct(t *testing.T) {
	names := map[string]bool{}

	for _, convention := range allConventions {
		name := convention.String()
		assert.NotEmpty(t, name)
		assert.False(t, names[name], "duplicate name %q", name)
		names[name] = true
	}

	assert.NotEmpty(t, Convention(99).String())
}

// TestZeroLengthPeriod checks a period of no length measures zero under every
// convention.
func TestZeroLengthPeriod(t *testing.T) {
	day := date(2024, time.February, 29)

	for _, convention := range allConventions {
		days, err := Days(day, day, convention)
		require.NoError(t, err)
		assert.Equal(t, 0, days)

		fraction, err := YearFraction(day, day, convention)
		require.NoError(t, err)
		assert.True(t, fraction.IsZero(), "%v", convention)
	}
}
