package compositeinterest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestBuildGrowthSchedule(t *testing.T) {
	// present = 1000, rate = 0.01, n = 3
	// period 1: interest = 10.00, balance = 1010.00
	// period 2: interest = 10.10, balance = 1020.10
	// period 3: interest = 10.20 (10.201 truncates to cents), balance = 1030.30
	present := money.MustMoneyFromFloat64(1000, money.USD)
	rate := money.MustFromFloat64(0.01)
	nper := money.MustFromFloat64(3)

	rows, err := BuildGrowthSchedule(present, rate, nper)
	require.NoError(t, err)
	require.Len(t, rows, 4)

	assert.True(t, rows[0].Balance.Equal(present))
	assert.True(t, rows[0].Change.IsZero())
	assert.Equal(t, 0.0, rows[0].ChangePercent.InexactFloat64())

	assert.InDelta(t, 1010.00, rows[1].Balance.ToDecimal().InexactFloat64(), 0.01)
	assert.InDelta(t, 0.01, rows[1].ChangePercent.InexactFloat64(), 0.0001)

	assert.InDelta(t, 1020.10, rows[2].Balance.ToDecimal().InexactFloat64(), 0.01)
	assert.InDelta(t, 0.01, rows[2].ChangePercent.InexactFloat64(), 0.0001)

	assert.InDelta(t, 1030.30, rows[3].Balance.ToDecimal().InexactFloat64(), 0.01)
	assert.InDelta(t, 0.01, rows[3].ChangePercent.InexactFloat64(), 0.0001)

	// The interest earned each period grows because it's computed on an
	// ever-larger balance, even though the percentage change stays constant.
	assert.True(t, rows[2].Change.GreaterThan(rows[1].Change))
	assert.True(t, rows[3].Change.GreaterThan(rows[2].Change))

	// SumInterest is cumulative and must match total growth over the period.
	totalGrowth := rows[3].Balance.Sub(present)
	assert.True(t, rows[3].SumInterest.Equal(totalGrowth))
}

func TestBuildGrowthScheduleInvalidPeriods(t *testing.T) {
	present := money.MustMoneyFromFloat64(1000, money.USD)
	rate := money.MustFromFloat64(0.01)

	tests := []money.Decimal{
		money.MustFromFloat64(0),
		money.MustFromFloat64(-3),
	}

	for _, nper := range tests {
		_, err := BuildGrowthSchedule(present, rate, nper)
		assert.True(t, errors.Is(err, ErrInvalidPeriods))
	}
}
