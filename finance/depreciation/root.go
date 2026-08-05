// Package depreciation computes asset depreciation schedules using the common
// accounting methods: straight-line, declining balance (including double
// declining balance with a switch to straight-line), sum-of-the-years'-digits,
// and MACRS (the US tax system's GDS tables).
//
// Every amount is a money.Money on the decimal engine, so schedules stay exact.
// StraightLine, DoubleDecliningBalance and SumOfYearsDigits depreciate down to
// exactly the salvage value, so their charges sum to the depreciable base
// (cost − salvage); MACRS recovers the full cost, since it ignores salvage.
//
// DecliningBalance is the exception: the pure declining-balance form approaches
// salvage geometrically without reaching it, so unless the clamp at salvage
// binds it leaves some book value behind and its charges sum to less than the
// depreciable base. DoubleDecliningBalance exists precisely because it adds the
// straight-line switchover that closes that gap.
//
// Basic usage:
//
//	cost := money.MustMoneyFromFloat64(10000, money.USD)
//	salvage := money.MustMoneyFromFloat64(1000, money.USD)
//	rows, _ := depreciation.StraightLine(cost, salvage, 5)
//	// rows[0].Depreciation == 1800, rows[4].BookValue == 1000
package depreciation

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/money"
)

var (
	// ErrNonPositiveCost is returned when the asset cost is not strictly
	// positive.
	ErrNonPositiveCost = errors.New("depreciation: cost must be positive")

	// ErrInvalidLife is returned when the useful life is not a positive number
	// of years.
	ErrInvalidLife = errors.New("depreciation: useful life must be at least 1 year")

	// ErrInvalidSalvage is returned when the salvage value is negative or
	// exceeds the cost.
	ErrInvalidSalvage = errors.New("depreciation: salvage must be between zero and cost")

	// ErrUnsupportedRecovery is returned by MACRS when no GDS table exists for
	// the requested recovery period.
	ErrUnsupportedRecovery = errors.New("depreciation: unsupported MACRS recovery period")
)

// Schedule is one year of a depreciation schedule: the year's depreciation
// expense, the depreciation accumulated through that year, and the asset's
// remaining book value at year end.
type Schedule struct {
	Year         int
	Depreciation money.Money
	Accumulated  money.Money
	BookValue    money.Money
}

// validate checks the shared cost/salvage/life inputs and returns the two
// amounts' currency.
func validate(cost, salvage money.Money, life int) error {
	if cost.Currency() != salvage.Currency() {
		return money.ErrCurrencyMismatch
	}

	if !cost.IsPositive() {
		return ErrNonPositiveCost
	}

	if life < 1 {
		return ErrInvalidLife
	}

	if salvage.IsNegative() || salvage.GreaterThan(cost) {
		return ErrInvalidSalvage
	}

	return nil
}

// appendRow accumulates one year of depreciation onto the running totals and
// appends the resulting Schedule row. It returns the updated accumulated and
// book-value amounts.
func appendRow(rows []Schedule, year int, depr, accumulated, book money.Money) ([]Schedule, money.Money, money.Money) {
	accumulated = accumulated.Add(depr)
	book = book.Sub(depr)

	return append(rows, Schedule{
		Year:         year,
		Depreciation: depr,
		Accumulated:  accumulated,
		BookValue:    book,
	}), accumulated, book
}
