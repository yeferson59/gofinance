// Package daycount implements the day-count conventions used to measure the
// time between two dates for interest accrual: how many days lie in a period
// and what fraction of a year that represents.
//
// These conventions are the shared foundation for date-based financial math
// (accrued interest, bond pricing, XIRR/XNPV): the same two dates yield
// different year fractions depending on the convention, so the choice is part
// of an instrument's terms.
//
// Supported conventions:
//
//   - Thirty360 — 30/360 (US / Bond Basis): every month is treated as 30 days
//     and every year as 360 days, with the standard end-of-month adjustments.
//   - Actual360 — Actual/360: actual calendar days over a 360-day year.
//   - Actual365Fixed — Actual/365 (Fixed): actual calendar days over a fixed
//     365-day year.
//   - ActualActualISDA — Actual/Actual (ISDA): actual days, splitting the
//     period at year boundaries so days in a leap year count over 366 and days
//     in a common year over 365.
//
// Year fractions are returned as decimal.Decimal so they compose with the rest
// of the library at full fixed-point precision.
//
// Basic usage:
//
//	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
//	end := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
//	yf, _ := daycount.YearFraction(start, end, daycount.Actual365Fixed)
//	// yf ≈ 0.4959 (181/365)
package daycount

import "errors"

// Convention identifies a day-count convention.
type Convention int

const (
	// Thirty360 is the 30/360 US (Bond Basis) convention.
	Thirty360 Convention = iota
	// Actual360 is the Actual/360 convention.
	Actual360
	// Actual365Fixed is the Actual/365 (Fixed) convention.
	Actual365Fixed
	// ActualActualISDA is the Actual/Actual (ISDA) convention.
	ActualActualISDA
)

var (
	// ErrEndBeforeStart is returned when the end date precedes the start
	// date. Day counts are defined from an earlier date to a later one.
	ErrEndBeforeStart = errors.New("daycount: end date is before start date")

	// ErrInvalidConvention is returned when an unknown Convention value is
	// passed.
	ErrInvalidConvention = errors.New("daycount: unknown convention")
)

// String returns a human-readable name for the convention.
func (c Convention) String() string {
	switch c {
	case Thirty360:
		return "30/360"
	case Actual360:
		return "Actual/360"
	case Actual365Fixed:
		return "Actual/365 Fixed"
	case ActualActualISDA:
		return "Actual/Actual ISDA"
	default:
		return "unknown"
	}
}
