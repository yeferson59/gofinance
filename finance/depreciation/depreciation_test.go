package depreciation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func usd(amount float64) money.Money {
	return money.MustMoneyFromFloat64(amount, money.USD)
}

// totals returns the sum of depreciation and the final book value of a schedule.
func totals(rows []Schedule) (float64, float64) {
	var sum float64
	for _, r := range rows {
		sum += r.Depreciation.InexactFloat64()
	}

	return sum, rows[len(rows)-1].BookValue.InexactFloat64()
}

func TestStraightLine(t *testing.T) {
	rows, err := StraightLine(usd(10000), usd(1000), 5)
	require.NoError(t, err)
	require.Len(t, rows, 5)

	assert.InDelta(t, 1800.0, rows[0].Depreciation.InexactFloat64(), 1e-9)

	sum, book := totals(rows)
	assert.InDelta(t, 9000.0, sum, 1e-9)
	assert.InDelta(t, 1000.0, book, 1e-9)
}

func TestSumOfYearsDigits(t *testing.T) {
	rows, err := SumOfYearsDigits(usd(10000), usd(1000), 5)
	require.NoError(t, err)

	// First year weight 5/15 of the 9000 base = 3000.
	assert.InDelta(t, 3000.0, rows[0].Depreciation.InexactFloat64(), 1e-9)

	sum, book := totals(rows)
	assert.InDelta(t, 9000.0, sum, 1e-9)
	assert.InDelta(t, 1000.0, book, 1e-9)
}

func TestDoubleDecliningBalance(t *testing.T) {
	rows, err := DoubleDecliningBalance(usd(10000), usd(1000), 5)
	require.NoError(t, err)

	// Rate 2/5 = 40% of 10000 in year one.
	assert.InDelta(t, 4000.0, rows[0].Depreciation.InexactFloat64(), 1e-9)

	sum, book := totals(rows)
	assert.InDelta(t, 9000.0, sum, 1e-9)
	// Switch to straight-line brings the book down to exactly salvage.
	assert.InDelta(t, 1000.0, book, 1e-9)
}

func TestDecliningBalanceNeverBelowSalvage(t *testing.T) {
	rows, err := DecliningBalance(usd(10000), usd(1000), 5, decimal.MustFromFloat64(2))
	require.NoError(t, err)

	assert.InDelta(t, 4000.0, rows[0].Depreciation.InexactFloat64(), 1e-9)

	for _, r := range rows {
		assert.GreaterOrEqual(t, r.BookValue.InexactFloat64(), 1000.0-1e-9)
	}
}

func TestMACRS(t *testing.T) {
	rows, err := MACRS(usd(10000), 5)
	require.NoError(t, err)
	require.Len(t, rows, 6)

	// 5-year GDS: 20% in year one.
	assert.InDelta(t, 2000.0, rows[0].Depreciation.InexactFloat64(), 1e-9)

	sum, book := totals(rows)
	// MACRS ignores salvage and recovers the full cost.
	assert.InDelta(t, 10000.0, sum, 1e-9)
	assert.InDelta(t, 0.0, book, 1e-9)
}

func TestMACRSUnsupported(t *testing.T) {
	_, err := MACRS(usd(10000), 4)
	assert.ErrorIs(t, err, ErrUnsupportedRecovery)
}

func TestValidationErrors(t *testing.T) {
	_, err := StraightLine(usd(10000), money.MustMoneyFromFloat64(1000, money.EUR), 5)
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	_, err = StraightLine(usd(0), usd(0), 5)
	assert.ErrorIs(t, err, ErrNonPositiveCost)

	_, err = StraightLine(usd(10000), usd(1000), 0)
	assert.ErrorIs(t, err, ErrInvalidLife)

	_, err = StraightLine(usd(10000), usd(12000), 5)
	assert.ErrorIs(t, err, ErrInvalidSalvage)

	_, err = StraightLine(usd(10000), usd(-100), 5)
	assert.ErrorIs(t, err, ErrInvalidSalvage)
}

func TestMustPanics(t *testing.T) {
	assert.Panics(t, func() { MustMACRS(usd(10000), 4) })
	assert.NotPanics(t, func() { MustStraightLine(usd(10000), usd(1000), 5) })
}
