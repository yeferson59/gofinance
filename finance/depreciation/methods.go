package depreciation

import "github.com/yeferson59/gofinance/money"

// StraightLine depreciates the asset evenly over its useful life:
//
//	annual = (cost − salvage) / life
//
// The final year absorbs any rounding remainder so the ending book value is
// exactly the salvage value.
//
// It returns money.ErrCurrencyMismatch, ErrNonPositiveCost, ErrInvalidLife, or
// ErrInvalidSalvage on invalid input.
func StraightLine(cost, salvage money.Money, life int) ([]Schedule, error) {
	if err := validate(cost, salvage, life); err != nil {
		return nil, err
	}

	base := cost.Sub(salvage)

	annualDec, err := base.ToDecimal().Div(money.MustFromInt64(int64(life), 0))
	if err != nil {
		return nil, err
	}

	annual := annualDec.ToMoney(cost.Currency())

	rows := make([]Schedule, 0, life)
	accumulated := cost.Sub(cost)
	book := cost

	for year := 1; year <= life; year++ {
		depr := annual
		if year == life {
			depr = book.Sub(salvage)
		}

		rows, accumulated, book = appendRow(rows, year, depr, accumulated, book)
	}

	return rows, nil
}

// MustStraightLine is like StraightLine but panics on error.
func MustStraightLine(cost, salvage money.Money, life int) []Schedule {
	rows, err := StraightLine(cost, salvage, life)
	if err != nil {
		panic(err)
	}

	return rows
}

// DecliningBalance applies a constant rate of factor/life to the declining book
// value each year (factor is 2 for double declining, 1.5 for 150%, and so on).
// Depreciation is capped so the book value never falls below salvage; this pure
// form does not switch to straight-line, so some book value above salvage may
// remain at the end. For the full depreciation that tax and most textbooks use,
// see DoubleDecliningBalance.
//
// It returns money.ErrCurrencyMismatch, ErrNonPositiveCost, ErrInvalidLife, or
// ErrInvalidSalvage on invalid input.
func DecliningBalance(cost, salvage money.Money, life int, factor money.Decimal) ([]Schedule, error) {
	if err := validate(cost, salvage, life); err != nil {
		return nil, err
	}

	rate, err := factor.Div(money.MustFromInt64(int64(life), 0))
	if err != nil {
		return nil, err
	}

	rows := make([]Schedule, 0, life)
	accumulated := cost.Sub(cost)
	book := cost

	for year := 1; year <= life; year++ {
		depr := book.ToDecimal().Mul(rate).ToMoney(cost.Currency())

		if maxDepr := book.Sub(salvage); depr.GreaterThan(maxDepr) {
			depr = maxDepr
		}

		rows, accumulated, book = appendRow(rows, year, depr, accumulated, book)
	}

	return rows, nil
}

// MustDecliningBalance is like DecliningBalance but panics on error.
func MustDecliningBalance(cost, salvage money.Money, life int, factor money.Decimal) []Schedule {
	rows, err := DecliningBalance(cost, salvage, life, factor)
	if err != nil {
		panic(err)
	}

	return rows
}

// DoubleDecliningBalance applies a 200% declining-balance rate but switches to
// straight-line over the remaining life whenever that yields a larger
// deduction, so the asset depreciates down to exactly its salvage value by the
// end of its life — the standard double-declining-balance convention.
//
// It returns money.ErrCurrencyMismatch, ErrNonPositiveCost, ErrInvalidLife, or
// ErrInvalidSalvage on invalid input.
func DoubleDecliningBalance(cost, salvage money.Money, life int) ([]Schedule, error) {
	if err := validate(cost, salvage, life); err != nil {
		return nil, err
	}

	rate, err := money.MustFromInt64(2, 0).Div(money.MustFromInt64(int64(life), 0))
	if err != nil {
		return nil, err
	}

	rows := make([]Schedule, 0, life)
	accumulated := cost.Sub(cost)
	book := cost

	for year := 1; year <= life; year++ {
		declining := book.ToDecimal().Mul(rate).ToMoney(cost.Currency())

		remainingLife := money.MustFromInt64(int64(life-year+1), 0)

		straightDec, err := book.Sub(salvage).ToDecimal().Div(remainingLife)
		if err != nil {
			return nil, err
		}

		straight := straightDec.ToMoney(cost.Currency())

		depr := declining
		if straight.GreaterThan(declining) {
			depr = straight
		}

		if maxDepr := book.Sub(salvage); depr.GreaterThan(maxDepr) {
			depr = maxDepr
		}

		rows, accumulated, book = appendRow(rows, year, depr, accumulated, book)
	}

	return rows, nil
}

// MustDoubleDecliningBalance is like DoubleDecliningBalance but panics on error.
func MustDoubleDecliningBalance(cost, salvage money.Money, life int) []Schedule {
	rows, err := DoubleDecliningBalance(cost, salvage, life)
	if err != nil {
		panic(err)
	}

	return rows
}

// SumOfYearsDigits accelerates depreciation by weighting each year by its
// remaining life over the sum of the years' digits:
//
//	yearₖ = (cost − salvage) × (life − k + 1) / (1 + 2 + … + life)
//
// The final year absorbs any rounding remainder so the ending book value is
// exactly the salvage value.
//
// It returns money.ErrCurrencyMismatch, ErrNonPositiveCost, ErrInvalidLife, or
// ErrInvalidSalvage on invalid input.
func SumOfYearsDigits(cost, salvage money.Money, life int) ([]Schedule, error) {
	if err := validate(cost, salvage, life); err != nil {
		return nil, err
	}

	base := cost.Sub(salvage)
	digitsSum := money.MustFromInt64(int64(life*(life+1)/2), 0)

	rows := make([]Schedule, 0, life)
	accumulated := cost.Sub(cost)
	book := cost

	for year := 1; year <= life; year++ {
		depr := book.Sub(salvage)

		if year < life {
			weight, err := money.MustFromInt64(int64(life-year+1), 0).Div(digitsSum)
			if err != nil {
				return nil, err
			}

			depr = base.ToDecimal().Mul(weight).ToMoney(cost.Currency())
		}

		rows, accumulated, book = appendRow(rows, year, depr, accumulated, book)
	}

	return rows, nil
}

// MustSumOfYearsDigits is like SumOfYearsDigits but panics on error.
func MustSumOfYearsDigits(cost, salvage money.Money, life int) []Schedule {
	rows, err := SumOfYearsDigits(cost, salvage, life)
	if err != nil {
		panic(err)
	}

	return rows
}

// MACRS builds the depreciation schedule under the US Modified Accelerated Cost
// Recovery System (General Depreciation System, half-year convention) for the
// given recovery period in years. Supported periods are 3, 5, 7, 10, 15, and
// 20 years. MACRS ignores salvage value and depreciates the full cost to zero;
// the final year absorbs any rounding remainder.
//
// It returns ErrNonPositiveCost for a non-positive cost and
// ErrUnsupportedRecovery for a recovery period without a GDS table.
func MACRS(cost money.Money, recovery int) ([]Schedule, error) {
	if !cost.IsPositive() {
		return nil, ErrNonPositiveCost
	}

	percentages, ok := macrsTable(recovery)
	if !ok {
		return nil, ErrUnsupportedRecovery
	}

	rows := make([]Schedule, 0, len(percentages))
	accumulated := cost.Sub(cost)
	book := cost

	for i, pct := range percentages {
		depr := book
		if i < len(percentages)-1 {
			depr = cost.ToDecimal().Mul(pct).ToMoney(cost.Currency())
		}

		rows, accumulated, book = appendRow(rows, i+1, depr, accumulated, book)
	}

	return rows, nil
}

// MustMACRS is like MACRS but panics on error.
func MustMACRS(cost money.Money, recovery int) []Schedule {
	rows, err := MACRS(cost, recovery)
	if err != nil {
		panic(err)
	}

	return rows
}
