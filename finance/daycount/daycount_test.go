package daycount

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestDaysThirty360(t *testing.T) {
	// 30/360: six whole months = 180 days.
	n, err := Days(date(2026, 1, 1), date(2026, 7, 1), Thirty360)
	require.NoError(t, err)
	assert.Equal(t, 180, n)
}

func TestDaysThirty360EndOfMonth(t *testing.T) {
	// Jan 30 → Mar 31: d1=30, d2=31→30, so 30*(3-1)+(30-30) = 60.
	n, err := Days(date(2026, 1, 30), date(2026, 3, 31), Thirty360)
	require.NoError(t, err)
	assert.Equal(t, 60, n)
}

func TestDaysActual(t *testing.T) {
	// Jan 1 → Jul 1 2026 (common year): 181 actual days.
	n, err := Days(date(2026, 1, 1), date(2026, 7, 1), Actual365Fixed)
	require.NoError(t, err)
	assert.Equal(t, 181, n)
}

func TestYearFractionThirty360(t *testing.T) {
	yf, err := YearFraction(date(2026, 1, 1), date(2026, 7, 1), Thirty360)
	require.NoError(t, err)
	assert.InDelta(t, 0.5, yf.InexactFloat64(), 1e-12)
}

func TestYearFractionActual365Fixed(t *testing.T) {
	yf, err := YearFraction(date(2026, 1, 1), date(2026, 7, 1), Actual365Fixed)
	require.NoError(t, err)
	assert.InDelta(t, 181.0/365.0, yf.InexactFloat64(), 1e-9)
}

func TestYearFractionActual360(t *testing.T) {
	yf, err := YearFraction(date(2026, 1, 1), date(2026, 7, 1), Actual360)
	require.NoError(t, err)
	assert.InDelta(t, 181.0/360.0, yf.InexactFloat64(), 1e-9)
}

func TestYearFractionActualActualSameYear(t *testing.T) {
	// 2026 is a common year → denominator 365.
	yf, err := YearFraction(date(2026, 1, 1), date(2026, 7, 1), ActualActualISDA)
	require.NoError(t, err)
	assert.InDelta(t, 181.0/365.0, yf.InexactFloat64(), 1e-9)
}

func TestYearFractionActualActualLeapYear(t *testing.T) {
	// 2024 is a leap year → denominator 366 for 182 days.
	yf, err := YearFraction(date(2024, 1, 1), date(2024, 7, 1), ActualActualISDA)
	require.NoError(t, err)
	assert.InDelta(t, 182.0/366.0, yf.InexactFloat64(), 1e-9)
}

func TestYearFractionActualActualAcrossYears(t *testing.T) {
	// Jul 1 2023 → Jul 1 2024 spans a common and a leap year:
	// 184/365 + 182/366 ≈ 1.001377.
	yf, err := YearFraction(date(2023, 7, 1), date(2024, 7, 1), ActualActualISDA)
	require.NoError(t, err)
	assert.InDelta(t, 184.0/365.0+182.0/366.0, yf.InexactFloat64(), 1e-9)
}

func TestActualActualMultipleWholeYears(t *testing.T) {
	// A full three common years is exactly 3.0.
	yf, err := YearFraction(date(2021, 1, 1), date(2024, 1, 1), ActualActualISDA)
	require.NoError(t, err)
	assert.InDelta(t, 3.0, yf.InexactFloat64(), 1e-9)
}

func TestEndBeforeStart(t *testing.T) {
	_, err := Days(date(2026, 7, 1), date(2026, 1, 1), Actual365Fixed)
	assert.ErrorIs(t, err, ErrEndBeforeStart)

	_, err = YearFraction(date(2026, 7, 1), date(2026, 1, 1), Actual365Fixed)
	assert.ErrorIs(t, err, ErrEndBeforeStart)
}

func TestInvalidConvention(t *testing.T) {
	_, err := YearFraction(date(2026, 1, 1), date(2026, 7, 1), Convention(99))
	assert.ErrorIs(t, err, ErrInvalidConvention)
}

func TestConventionString(t *testing.T) {
	assert.Equal(t, "30/360", Thirty360.String())
	assert.Equal(t, "Actual/Actual ISDA", ActualActualISDA.String())
}

func TestNormalizeIgnoresClockAndZone(t *testing.T) {
	// Same calendar day, different clock/zone → zero-length period.
	loc := time.FixedZone("X", 5*3600)
	start := time.Date(2026, 3, 15, 23, 30, 0, 0, loc)
	end := time.Date(2026, 3, 15, 1, 0, 0, 0, time.UTC)
	n, err := Days(start, end, Actual365Fixed)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
}
