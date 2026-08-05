package depreciation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// method is a schedule-producing method taking the standard cost/salvage/life
// triple.
type method func(cost, salvage money.Money, life int) ([]Schedule, error)

// fullyDepreciatingMethods are the methods that reach the salvage value
// exactly, so their charges must sum to the depreciable base.
//
// DecliningBalance is deliberately absent: the pure declining-balance form
// approaches salvage geometrically without reaching it, which its own doc
// comment states. It gets its own weaker invariant below.
func fullyDepreciatingMethods() map[string]method {
	return map[string]method{
		"StraightLine":           StraightLine,
		"DoubleDecliningBalance": DoubleDecliningBalance,
		"SumOfYearsDigits":       SumOfYearsDigits,
	}
}

// methods is every schedule-producing method, including the declining-balance
// variants, for the invariants that hold across all of them.
func methods() map[string]method {
	all := fullyDepreciatingMethods()

	all["DecliningBalance150"] = func(c, s money.Money, l int) ([]Schedule, error) {
		return DecliningBalance(c, s, l, decimal.MustFromFloat64(1.5))
	}
	all["DecliningBalance200"] = func(c, s money.Money, l int) ([]Schedule, error) {
		return DecliningBalance(c, s, l, decimal.MustFromFloat64(2))
	}

	return all
}

// TestScheduleSumsToDepreciableBase is the invariant every method must satisfy:
// the depreciation charged over the asset's life adds up to exactly cost −
// salvage, and the closing book value is exactly the salvage.
func TestScheduleSumsToDepreciableBase(t *testing.T) {
	assets := []struct {
		name    string
		cost    float64
		salvage float64
		life    int
	}{
		{"typical", 10000, 1000, 5},
		{"no salvage", 10000, 0, 5},
		{"high salvage", 10000, 9000, 5},
		{"one year", 5000, 500, 1},
		{"long life", 100000, 5000, 20},
		{"uneven division", 10000, 1, 3},
	}

	for _, asset := range assets {
		for name, depreciate := range fullyDepreciatingMethods() {
			t.Run(asset.name+"/"+name, func(t *testing.T) {
				cost := money.MustMoneyFromFloat64(asset.cost, money.USD)
				salvage := money.MustMoneyFromFloat64(asset.salvage, money.USD)

				schedule, err := depreciate(cost, salvage, asset.life)
				require.NoError(t, err)
				require.Len(t, schedule, asset.life)

				total := money.MustMoneyFromFloat64(0, money.USD)
				for _, row := range schedule {
					total = total.Add(row.Depreciation)
				}

				assert.InDelta(t, asset.cost-asset.salvage, total.InexactFloat64(), 1e-9,
					"depreciation must sum to cost − salvage")

				last := schedule[len(schedule)-1]
				assert.InDelta(t, asset.salvage, last.BookValue.InexactFloat64(), 1e-9,
					"closing book value must equal salvage")
				assert.InDelta(t, asset.cost-asset.salvage, last.Accumulated.InexactFloat64(), 1e-9,
					"accumulated depreciation must equal the depreciable base")
			})
		}
	}
}

// TestScheduleIsMonotonic checks the structural properties of every schedule:
// years run 1..life, book value never rises and never dips below salvage,
// accumulated depreciation never falls, and no year charges a negative amount.
func TestScheduleIsMonotonic(t *testing.T) {
	cost := money.MustMoneyFromFloat64(50000, money.USD)
	salvage := money.MustMoneyFromFloat64(5000, money.USD)

	for name, method := range methods() {
		t.Run(name, func(t *testing.T) {
			schedule, err := method(cost, salvage, 8)
			require.NoError(t, err)

			previousBook := cost.InexactFloat64()
			previousAccumulated := 0.0

			for i, row := range schedule {
				assert.Equal(t, i+1, row.Year)
				assert.GreaterOrEqual(t, row.Depreciation.InexactFloat64(), 0.0,
					"year %d charges a negative amount", row.Year)
				assert.LessOrEqual(t, row.BookValue.InexactFloat64(), previousBook,
					"book value rose in year %d", row.Year)
				assert.GreaterOrEqual(t, row.BookValue.InexactFloat64(), salvage.InexactFloat64()-1e-9,
					"book value fell below salvage in year %d", row.Year)
				assert.GreaterOrEqual(t, row.Accumulated.InexactFloat64(), previousAccumulated,
					"accumulated depreciation fell in year %d", row.Year)

				previousBook = row.BookValue.InexactFloat64()
				previousAccumulated = row.Accumulated.InexactFloat64()
			}
		})
	}
}

// TestStraightLineIsLevel checks the defining property of the straight-line
// method: every year charges the same amount.
func TestStraightLineIsLevel(t *testing.T) {
	schedule, err := StraightLine(
		money.MustMoneyFromFloat64(10000, money.USD),
		money.MustMoneyFromFloat64(1000, money.USD), 5)
	require.NoError(t, err)

	for _, row := range schedule {
		assert.InDelta(t, 1800.0, row.Depreciation.InexactFloat64(), 1e-9)
	}
}

// TestSumOfYearsDigitsIsDecreasing checks that the charge falls every year and
// follows the digit weights: for a 5-year life the weights are 5/15, 4/15, …
func TestSumOfYearsDigitsIsDecreasing(t *testing.T) {
	schedule, err := SumOfYearsDigits(
		money.MustMoneyFromFloat64(10000, money.USD),
		money.MustMoneyFromFloat64(1000, money.USD), 5)
	require.NoError(t, err)

	base := 9000.0
	weights := []float64{5, 4, 3, 2, 1}

	for i, row := range schedule {
		assert.InDelta(t, base*weights[i]/15, row.Depreciation.InexactFloat64(), 1e-9)

		if i > 0 {
			assert.Less(t, row.Depreciation.InexactFloat64(), schedule[i-1].Depreciation.InexactFloat64())
		}
	}
}

// TestDoubleDecliningSwitchesToStraightLine checks the switchover the method is
// named for: it starts at twice the straight-line rate and moves to a level
// charge once that becomes the larger deduction.
func TestDoubleDecliningSwitchesToStraightLine(t *testing.T) {
	schedule, err := DoubleDecliningBalance(
		money.MustMoneyFromFloat64(10000, money.USD),
		money.MustMoneyFromFloat64(1000, money.USD), 5)
	require.NoError(t, err)

	// Year 1 at 2/5 of the opening book value.
	assert.InDelta(t, 4000.0, schedule[0].Depreciation.InexactFloat64(), 1e-9)
	// Year 2 at 2/5 of the remaining 6000.
	assert.InDelta(t, 2400.0, schedule[1].Depreciation.InexactFloat64(), 1e-9)

	// With no salvage the switchover is visible as a level tail: at 10000 over
	// 5 years the declining charges are 4000, 2400, 1440, and from year 4 the
	// straight-line charge over the remaining life (2160/2 = 1080) is larger,
	// so years 4 and 5 both charge 1080.
	noSalvage, err := DoubleDecliningBalance(
		money.MustMoneyFromFloat64(10000, money.USD),
		money.MustMoneyFromFloat64(0, money.USD), 5)
	require.NoError(t, err)

	expected := []float64{4000, 2400, 1440, 1080, 1080}
	for i, want := range expected {
		assert.InDelta(t, want, noSalvage[i].Depreciation.InexactFloat64(), 1e-9,
			"year %d", i+1)
	}
}

// TestDecliningBalanceLeavesResidual pins the documented difference between the
// pure declining-balance form and the double-declining one: approaching salvage
// geometrically never reaches it, so book value is left behind unless the clamp
// at salvage binds first.
func TestDecliningBalanceLeavesResidual(t *testing.T) {
	cost := money.MustMoneyFromFloat64(10000, money.USD)
	salvage := money.MustMoneyFromFloat64(1, money.USD)

	pure, err := DecliningBalance(cost, salvage, 3, decimal.MustFromFloat64(2))
	require.NoError(t, err)

	// 10000 -> 6666.67 -> 4444.44 -> 2962.96 charged in total 7037.04, well
	// short of the 9999 depreciable base.
	last := pure[len(pure)-1]
	assert.Greater(t, last.BookValue.InexactFloat64(), salvage.InexactFloat64(),
		"the pure form must leave book value above salvage")

	// The double-declining variant closes that gap on the same asset.
	switched, err := DoubleDecliningBalance(cost, salvage, 3)
	require.NoError(t, err)
	assert.InDelta(t, 1.0, switched[len(switched)-1].BookValue.InexactFloat64(), 1e-9)
}

// TestDecliningBalanceClampsAtSalvage checks the case where the clamp does
// bind: with a high salvage the geometric decline would undershoot it, so the
// schedule stops exactly at salvage.
func TestDecliningBalanceClampsAtSalvage(t *testing.T) {
	schedule, err := DecliningBalance(
		money.MustMoneyFromFloat64(10000, money.USD),
		money.MustMoneyFromFloat64(5000, money.USD), 5, decimal.MustFromFloat64(2))
	require.NoError(t, err)

	last := schedule[len(schedule)-1]
	assert.InDelta(t, 5000.0, last.BookValue.InexactFloat64(), 1e-9)

	total := money.MustMoneyFromFloat64(0, money.USD)
	for _, row := range schedule {
		total = total.Add(row.Depreciation)
	}

	assert.InDelta(t, 5000.0, total.InexactFloat64(), 1e-9)
}

// TestDecliningBalanceFactorOrdering checks that a faster factor front-loads
// the deduction: whatever the factor, the total is the same, but a 200%
// schedule charges more in year one than a 150% one.
func TestDecliningBalanceFactorOrdering(t *testing.T) {
	cost := money.MustMoneyFromFloat64(20000, money.USD)
	salvage := money.MustMoneyFromFloat64(2000, money.USD)

	slower, err := DecliningBalance(cost, salvage, 6, decimal.MustFromFloat64(1.5))
	require.NoError(t, err)

	faster, err := DecliningBalance(cost, salvage, 6, decimal.MustFromFloat64(2))
	require.NoError(t, err)

	assert.Greater(t, faster[0].Depreciation.InexactFloat64(), slower[0].Depreciation.InexactFloat64())
}

// TestMACRSAllRecoveryPeriods checks every supported GDS table: the schedule
// runs one year longer than the recovery period (the half-year convention) and
// the percentages recover the whole cost, since MACRS ignores salvage.
func TestMACRSAllRecoveryPeriods(t *testing.T) {
	cost := money.MustMoneyFromFloat64(100000, money.USD)

	for _, recovery := range []int{3, 5, 7, 10, 15, 20} {
		t.Run(decimal.MustFromInt64(int64(recovery), 0).String(), func(t *testing.T) {
			schedule, err := MACRS(cost, recovery)
			require.NoError(t, err)

			// The half-year convention spreads an n-year class over n+1 years.
			assert.Len(t, schedule, recovery+1)

			total := money.MustMoneyFromFloat64(0, money.USD)
			for _, row := range schedule {
				total = total.Add(row.Depreciation)
				assert.Positive(t, row.Depreciation.InexactFloat64())
			}

			assert.InDelta(t, 100000.0, total.InexactFloat64(), 1e-6,
				"MACRS recovers the full cost")

			last := schedule[len(schedule)-1]
			assert.InDelta(t, 0.0, last.BookValue.InexactFloat64(), 1e-6)
		})
	}
}

// TestMACRSUnsupportedRecovery checks the recovery periods with no GDS table.
func TestMACRSUnsupportedRecovery(t *testing.T) {
	cost := money.MustMoneyFromFloat64(100000, money.USD)

	for _, recovery := range []int{0, 1, 2, 4, 6, 8, 25, -5} {
		_, err := MACRS(cost, recovery)
		assert.ErrorIs(t, err, ErrUnsupportedRecovery)
	}
}

// TestMACRSRejectsNonPositiveCost checks the cost validation MACRS shares with
// the other methods.
func TestMACRSRejectsNonPositiveCost(t *testing.T) {
	_, err := MACRS(money.MustMoneyFromFloat64(0, money.USD), 5)
	assert.ErrorIs(t, err, ErrNonPositiveCost)

	_, err = MACRS(money.MustMoneyFromFloat64(-100, money.USD), 5)
	assert.ErrorIs(t, err, ErrNonPositiveCost)
}

// TestValidationAcrossMethods sweeps the shared validation: every method must
// reject the same invalid cost, salvage and life combinations rather than
// producing a schedule that does not add up.
func TestValidationAcrossMethods(t *testing.T) {
	invalid := []struct {
		name    string
		cost    money.Money
		salvage money.Money
		life    int
		want    error
	}{
		{"zero cost", money.MustMoneyFromFloat64(0, money.USD), money.MustMoneyFromFloat64(0, money.USD), 5, ErrNonPositiveCost},
		{"negative cost", money.MustMoneyFromFloat64(-1000, money.USD), money.MustMoneyFromFloat64(0, money.USD), 5, ErrNonPositiveCost},
		{"zero life", money.MustMoneyFromFloat64(1000, money.USD), money.MustMoneyFromFloat64(0, money.USD), 0, ErrInvalidLife},
		{"negative life", money.MustMoneyFromFloat64(1000, money.USD), money.MustMoneyFromFloat64(0, money.USD), -3, ErrInvalidLife},
		{"negative salvage", money.MustMoneyFromFloat64(1000, money.USD), money.MustMoneyFromFloat64(-100, money.USD), 5, ErrInvalidSalvage},
		{"salvage above cost", money.MustMoneyFromFloat64(1000, money.USD), money.MustMoneyFromFloat64(2000, money.USD), 5, ErrInvalidSalvage},
		{"currency mismatch", money.MustMoneyFromFloat64(1000, money.USD), money.MustMoneyFromFloat64(100, money.EUR), 5, money.ErrCurrencyMismatch},
	}

	for _, test := range invalid {
		for name, method := range methods() {
			t.Run(test.name+"/"+name, func(t *testing.T) {
				_, err := method(test.cost, test.salvage, test.life)
				assert.ErrorIs(t, err, test.want)
			})
		}
	}
}

// TestSalvageEqualToCost checks the boundary the validation allows: nothing to
// depreciate, so every year charges zero.
func TestSalvageEqualToCost(t *testing.T) {
	amount := money.MustMoneyFromFloat64(1000, money.USD)

	for name, method := range methods() {
		t.Run(name, func(t *testing.T) {
			schedule, err := method(amount, amount, 5)
			require.NoError(t, err)

			for _, row := range schedule {
				assert.InDelta(t, 0.0, row.Depreciation.InexactFloat64(), 1e-9)
				assert.InDelta(t, 1000.0, row.BookValue.InexactFloat64(), 1e-9)
			}
		})
	}
}

// TestMustVariants covers the panicking helpers, three of which had no
// coverage: they return the schedule on valid input and panic on invalid.
func TestMustVariants(t *testing.T) {
	cost := money.MustMoneyFromFloat64(10000, money.USD)
	salvage := money.MustMoneyFromFloat64(1000, money.USD)
	bad := money.MustMoneyFromFloat64(-1, money.USD)

	assert.NotPanics(t, func() {
		assert.Len(t, MustStraightLine(cost, salvage, 5), 5)
		assert.Len(t, MustDecliningBalance(cost, salvage, 5, decimal.MustFromFloat64(2)), 5)
		assert.Len(t, MustDoubleDecliningBalance(cost, salvage, 5), 5)
		assert.Len(t, MustSumOfYearsDigits(cost, salvage, 5), 5)
		assert.Len(t, MustMACRS(cost, 5), 6)
	})

	assert.Panics(t, func() { MustStraightLine(bad, salvage, 5) })
	assert.Panics(t, func() { MustDecliningBalance(bad, salvage, 5, decimal.MustFromFloat64(2)) })
	assert.Panics(t, func() { MustDoubleDecliningBalance(bad, salvage, 5) })
	assert.Panics(t, func() { MustSumOfYearsDigits(bad, salvage, 5) })
	assert.Panics(t, func() { MustMACRS(cost, 4) })
}

// TestCurrencyIsPreserved checks the schedules carry the asset's currency
// rather than defaulting to USD.
func TestCurrencyIsPreserved(t *testing.T) {
	cost := money.MustMoneyFromFloat64(10000, money.EUR)
	salvage := money.MustMoneyFromFloat64(1000, money.EUR)

	for name, method := range methods() {
		t.Run(name, func(t *testing.T) {
			schedule, err := method(cost, salvage, 5)
			require.NoError(t, err)

			for _, row := range schedule {
				assert.Equal(t, money.EUR, row.Depreciation.Currency())
				assert.Equal(t, money.EUR, row.BookValue.Currency())
				assert.Equal(t, money.EUR, row.Accumulated.Currency())
			}
		})
	}

	macrs, err := MACRS(cost, 5)
	require.NoError(t, err)
	assert.Equal(t, money.EUR, macrs[0].BookValue.Currency())
}
