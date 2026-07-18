package daycount

import (
	"time"

	"github.com/yeferson59/gofinance/v2/decimal"
)

// normalize strips the clock time and time zone, reducing t to a UTC
// midnight so that day differences are exact whole numbers regardless of the
// caller's location or daylight-saving transitions.
func normalize(t time.Time) time.Time {
	y, m, d := t.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// actualDays returns the number of calendar days from start to end, both
// already normalized to UTC midnight.
func actualDays(start, end time.Time) int {
	return int(end.Sub(start).Hours()/24 + 0.5)
}

func isLeap(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

func daysInYear(year int) int {
	if isLeap(year) {
		return 366
	}

	return 365
}

// fraction returns num/den as a decimal.Decimal.
func fraction(num, den int) (decimal.Decimal, error) {
	n, err := decimal.NewFromInt64(int64(num), 0)
	if err != nil {
		return decimal.Decimal{}, err
	}

	d, err := decimal.NewFromInt64(int64(den), 0)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return n.Div(d)
}

// thirty360Days returns the 30/360 (US Bond Basis) day count between two
// normalized dates, applying the standard end-of-month adjustments.
func thirty360Days(start, end time.Time) int {
	y1, m1, d1 := start.Date()
	y2, m2, d2 := end.Date()

	if d1 == 31 {
		d1 = 30
	}

	if d2 == 31 && d1 == 30 {
		d2 = 30
	}

	return 360*(y2-y1) + 30*(int(m2)-int(m1)) + (d2 - d1)
}

// Days returns the day-count numerator between start and end under the given
// convention: the 360-based count for Thirty360 and the actual calendar-day
// count for the Actual conventions. The dates are reduced to UTC midnight
// before counting.
//
// It returns ErrEndBeforeStart if end precedes start and ErrInvalidConvention
// for an unknown convention.
func Days(start, end time.Time, conv Convention) (int, error) {
	s := normalize(start)
	e := normalize(end)

	if e.Before(s) {
		return 0, ErrEndBeforeStart
	}

	switch conv {
	case Thirty360:
		return thirty360Days(s, e), nil
	case Actual360, Actual365Fixed, ActualActualISDA:
		return actualDays(s, e), nil
	default:
		return 0, ErrInvalidConvention
	}
}

// YearFraction returns the fraction of a year between start and end under the
// given convention, as a decimal.Decimal. The dates are reduced to UTC midnight
// before measuring.
//
// It returns ErrEndBeforeStart if end precedes start and ErrInvalidConvention
// for an unknown convention.
func YearFraction(start, end time.Time, conv Convention) (decimal.Decimal, error) {
	s := normalize(start)
	e := normalize(end)

	if e.Before(s) {
		return decimal.Decimal{}, ErrEndBeforeStart
	}

	switch conv {
	case Thirty360:
		return fraction(thirty360Days(s, e), 360)
	case Actual360:
		return fraction(actualDays(s, e), 360)
	case Actual365Fixed:
		return fraction(actualDays(s, e), 365)
	case ActualActualISDA:
		return actualActualISDA(s, e)
	default:
		return decimal.Decimal{}, ErrInvalidConvention
	}
}

// actualActualISDA computes the Actual/Actual (ISDA) year fraction, counting
// days that fall in a leap year over 366 and days in a common year over 365,
// splitting the period at calendar-year boundaries.
func actualActualISDA(start, end time.Time) (decimal.Decimal, error) {
	startYear := start.Year()
	endYear := end.Year()

	if startYear == endYear {
		return fraction(actualDays(start, end), daysInYear(startYear))
	}

	// Days from start up to the first day of the next year.
	nextYearStart := time.Date(startYear+1, 1, 1, 0, 0, 0, 0, time.UTC)
	firstPart, err := fraction(actualDays(start, nextYearStart), daysInYear(startYear))
	if err != nil {
		return decimal.Decimal{}, err
	}

	// Days from the first day of the end year up to end.
	endYearStart := time.Date(endYear, 1, 1, 0, 0, 0, 0, time.UTC)
	lastPart, err := fraction(actualDays(endYearStart, end), daysInYear(endYear))
	if err != nil {
		return decimal.Decimal{}, err
	}

	total := firstPart.Add(lastPart)

	if wholeYears := endYear - startYear - 1; wholeYears > 0 {
		years, err := decimal.NewFromInt64(int64(wholeYears), 0)
		if err != nil {
			return decimal.Decimal{}, err
		}

		total = total.Add(years)
	}

	return total, nil
}
