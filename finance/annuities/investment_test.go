package annuities

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestBuildInvestmentSchedule(t *testing.T) {
	// principal = 0, contribution = 100, rate = 0.01, n = 3
	// period 1: balance = 100.00, change = 100.00 (previous was 0, so
	//           ChangePercent is left at zero)
	// period 2: interest = 1.00, balance = 201.00, change% = 1.01
	// period 3: interest = 2.01, balance = 303.01, change% = 0.507512
	principal := money.MustMoneyFromFloat64(0, money.USD)
	contribution := money.MustMoneyFromFloat64(100, money.USD)
	rate := decimal.MustFromFloat64(0.01)
	nper := decimal.MustFromFloat64(3)

	rows, err := BuildInvestmentSchedule(principal, contribution, rate, nper)
	require.NoError(t, err)
	require.Len(t, rows, 4)

	assert.True(t, rows[0].Balance.IsZero())

	assert.InDelta(t, 100.00, rows[1].Balance.GetDecimal().InexactFloat64(), 0.01)
	assert.InDelta(t, 100.00, rows[1].Change.GetDecimal().InexactFloat64(), 0.01)
	assert.Equal(t, 0.0, rows[1].ChangePercent.InexactFloat64())

	assert.InDelta(t, 201.00, rows[2].Balance.GetDecimal().InexactFloat64(), 0.01)
	assert.InDelta(t, 1.01, rows[2].ChangePercent.InexactFloat64(), 0.0001)

	assert.InDelta(t, 303.01, rows[3].Balance.GetDecimal().InexactFloat64(), 0.01)
	assert.InDelta(t, 0.507512, rows[3].ChangePercent.InexactFloat64(), 0.0001)

	// The contribution's relative weight shrinks over time as the balance
	// compounds past it.
	assert.True(t, rows[2].ChangePercent.GreaterThan(rows[3].ChangePercent))

	assert.InDelta(t, 300.00, rows[3].SumContributions.GetDecimal().InexactFloat64(), 0.01)
	assert.InDelta(t, 3.01, rows[3].SumInterest.GetDecimal().InexactFloat64(), 0.01)
}

func TestBuildAnticipateInvestmentSchedule(t *testing.T) {
	// Same inputs as TestBuildInvestmentSchedule, but each contribution is
	// made at the start of its period, so it also earns interest in its own
	// first period:
	// period 1: balance = 101.00
	// period 2: balance = 203.01
	// period 3: balance = 306.0401
	principal := money.MustMoneyFromFloat64(0, money.USD)
	contribution := money.MustMoneyFromFloat64(100, money.USD)
	rate := decimal.MustFromFloat64(0.01)
	nper := decimal.MustFromFloat64(3)

	rows, err := BuildAnticipateInvestmentSchedule(principal, contribution, rate, nper)
	require.NoError(t, err)
	require.Len(t, rows, 4)

	assert.InDelta(t, 101.00, rows[1].Balance.GetDecimal().InexactFloat64(), 0.01)
	assert.InDelta(t, 203.01, rows[2].Balance.GetDecimal().InexactFloat64(), 0.01)
	assert.InDelta(t, 306.04, rows[3].Balance.GetDecimal().InexactFloat64(), 0.01)

	// Each period's balance is at least as large as the ordinary annuity's,
	// since contributions here earn interest a period earlier.
	ordinaryRows, err := BuildInvestmentSchedule(principal, contribution, rate, nper)
	require.NoError(t, err)
	for i := range rows {
		assert.True(t, rows[i].Balance.GreaterThanOrEqual(ordinaryRows[i].Balance))
	}
}

func TestBuildInvestmentScheduleWithPrincipal(t *testing.T) {
	// principal = 1000, contribution = 100, rate = 0.01, n = 2: verifies
	// the schedule's final balance matches FutureWithContributions.
	principal := money.MustMoneyFromFloat64(1000, money.USD)
	contribution := money.MustMoneyFromFloat64(100, money.USD)
	rate := decimal.MustFromFloat64(0.01)
	nper := decimal.MustFromFloat64(2)

	rows, err := BuildInvestmentSchedule(principal, contribution, rate, nper)
	require.NoError(t, err)
	require.Len(t, rows, 3)

	total := NewAnnuity().
		Present(1000, money.USD).
		Value(100, money.USD).
		Rate(0.01).
		Periods(2).
		Monthly().
		MustFutureValue()

	assert.InDelta(t, total.GetDecimal().InexactFloat64(), rows[2].Balance.GetDecimal().InexactFloat64(), 0.01)
}

func TestBuildInvestmentScheduleCurrencyMismatch(t *testing.T) {
	principal := money.MustMoneyFromFloat64(1000, money.USD)
	contribution := money.MustMoneyFromFloat64(100, money.EUR)

	_, err := BuildInvestmentSchedule(principal, contribution, decimal.MustFromFloat64(0.01), decimal.MustFromFloat64(12))
	assert.True(t, errors.Is(err, money.ErrCurrencyMismatch))
}

func TestBuildInvestmentScheduleInvalidPeriods(t *testing.T) {
	principal := money.MustMoneyFromFloat64(1000, money.USD)
	contribution := money.MustMoneyFromFloat64(100, money.USD)

	tests := []decimal.Decimal{
		decimal.MustFromFloat64(0),
		decimal.MustFromFloat64(-3),
	}

	for _, nper := range tests {
		_, err := BuildInvestmentSchedule(principal, contribution, decimal.MustFromFloat64(0.01), nper)
		assert.True(t, errors.Is(err, ErrInvalidPeriods))
	}
}
