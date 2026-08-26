package invariants

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/annuities"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/finance/returns"
	"github.com/yeferson59/gofinance/v2/finance/simpleinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// The types in this library that carry several amounts let a caller set only
// the ones a given calculation needs: an annuity described by a present value
// alone has no periodic payment, a simple interest configuration solved for
// its interest has no interest yet. The amounts left out are the zero
// money.Money, which carries money.XXX — the ISO code for "no currency".
//
// Deriving a result's currency from one particular field therefore produced
// XXX whenever that field was the unset one, and combining two such results
// panicked with a currency mismatch. It shipped in annuities, simpleinterest
// and returns.HoldingPeriodReturn before this sweep found it.
//
// The rule these tests hold every such type to:
//
//	a partially configured value computes in the currency the caller set,
//	never in XXX, and never panics.
//
// The sweep enumerates the partial configurations rather than sampling them,
// because the defect lived precisely in the combinations no single-case test
// visited.

const sweepCurrency = money.EUR

// set and unset are the two states each optional amount can be in.
func set(amount float64) money.Money {
	return money.MustMoneyFromFloat64(amount, sweepCurrency)
}

func unset() money.Money {
	return money.Money{}
}

// assertResolved checks one calculation of a partially configured value: it
// must not panic, and whatever it returns must carry the caller's currency.
func assertResolved(t *testing.T, name string, compute func() (money.Money, error)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		assert.NotPanics(t, func() {
			result, err := compute()
			if err != nil {
				// Refusing to compute from amounts that are not there is
				// fine; carrying the wrong currency is not.
				return
			}

			assert.Equal(t, sweepCurrency, result.GetCurrency(),
				"%s returned %v in %v", name, result, result.GetCurrency())
		})
	})
}

// TestAnnuityPartialConfigurations sweeps which of the annuity's three amounts
// are set.
func TestAnnuityPartialConfigurations(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromInt64(12, 0), compoundinterest.Monthly)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.01), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	configurations := []struct {
		name                   string
		value, present, future money.Money
	}{
		{"payment only", set(100), unset(), unset()},
		{"present only", unset(), set(1000), unset()},
		{"future only", unset(), unset(), set(2000)},
		{"payment and present", set(100), set(1000), unset()},
		{"payment and future", set(100), unset(), set(2000)},
		{"present and future", unset(), set(1000), set(2000)},
		{"all three", set(100), set(1000), set(2000)},
	}

	for _, configuration := range configurations {
		t.Run(configuration.name, func(t *testing.T) {
			annuity, err := annuities.New(
				configuration.value, configuration.present, configuration.future, period, rate)
			require.NoError(t, err)

			assertResolved(t, "Present", annuity.Present)
			assertResolved(t, "AnticipatePresent", annuity.AnticipatePresent)
			assertResolved(t, "Future", annuity.Future)
			assertResolved(t, "AnticipateFuture", annuity.AnticipateFuture)
			assertResolved(t, "PrincipalFuture", annuity.PrincipalFuture)
			assertResolved(t, "FutureWithContributions", annuity.FutureWithContributions)
			assertResolved(t, "AnticipateFutureWithContributions", annuity.AnticipateFutureWithContributions)
			assertResolved(t, "PaymentFromPresentValue", annuity.PaymentFromPresentValue)
			assertResolved(t, "PaymentFromFutureValue", annuity.PaymentFromFutureValue)
			assertResolved(t, "PresentDeferred", func() (money.Money, error) {
				return annuity.PresentDeferred(3)
			})
		})
	}
}

// TestSimpleInterestPartialConfigurations sweeps which of the simple interest
// configuration's three amounts are set.
func TestSimpleInterestPartialConfigurations(t *testing.T) {
	period := simpleinterest.NewPeriod(decimal.MustFromInt64(12, 0), simpleinterest.Months)
	rate := decimal.MustFromFloat64(0.01)

	configurations := []struct {
		name                      string
		future, present, interest money.Money
	}{
		{"present only", unset(), set(1000), unset()},
		{"future only", set(2000), unset(), unset()},
		{"interest only", unset(), unset(), set(120)},
		{"present and interest", unset(), set(1000), set(120)},
		{"future and present", set(2000), set(1000), unset()},
		{"future and interest", set(2000), unset(), set(120)},
		{"all three", set(2000), set(1000), set(120)},
	}

	for _, configuration := range configurations {
		t.Run(configuration.name, func(t *testing.T) {
			simple := simpleinterest.New(
				configuration.future, configuration.present, configuration.interest, rate, period)

			assertResolved(t, "Future", simple.Future)
			assertResolved(t, "FutureWithRateInterest", simple.FutureWithRateInterest)
			assertResolved(t, "Present", simple.Present)
			assertResolved(t, "PresentWithFuture", simple.PresentWithFuture)
			assertResolved(t, "Interest", simple.Interest)
			assertResolved(t, "InterestWithPresentAndFuture", simple.InterestWithPresentAndFuture)
		})
	}
}

// TestCompoundInterestPartialConfigurations sweeps the two amounts a compound
// interest configuration carries. It came through the sweep clean; this keeps
// it that way.
func TestCompoundInterestPartialConfigurations(t *testing.T) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromInt64(12, 0), compoundinterest.Monthly)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.01), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	configurations := []struct {
		name            string
		present, future money.Money
	}{
		{"present only", set(1000), unset()},
		{"future only", unset(), set(2000)},
		{"both", set(1000), set(2000)},
	}

	for _, configuration := range configurations {
		t.Run(configuration.name, func(t *testing.T) {
			compound, err := compoundinterest.New(
				configuration.present, configuration.future, rate, period)
			require.NoError(t, err)

			assertResolved(t, "Present", compound.Present)
			assertResolved(t, "Future", compound.Future)
		})
	}
}

// TestOptionalAmountsMayBeOmitted covers the metrics that take an amount a
// caller can reasonably leave out — a share that paid no dividend, a period
// with no external flow — where an unset money.Money used to be rejected as a
// currency mismatch rather than read as zero.
func TestOptionalAmountsMayBeOmitted(t *testing.T) {
	initial := set(1000)
	final := set(1100)

	omitted, err := returns.HoldingPeriodReturn(initial, final, unset())
	require.NoError(t, err)

	explicit, err := returns.HoldingPeriodReturn(initial, final, set(0))
	require.NoError(t, err)

	assert.InDelta(t, explicit.InexactFloat64(), omitted.InexactFloat64(), 1e-12)
	assert.InDelta(t, 0.10, omitted.InexactFloat64(), 1e-12)

	// A period with no external flow is the same as one flowing zero.
	withUnset, err := returns.MoneyWeightedReturn(initial, []money.Money{unset()}, final)
	require.NoError(t, err)

	withZero, err := returns.MoneyWeightedReturn(initial, []money.Money{set(0)}, final)
	require.NoError(t, err)

	assert.InDelta(t, withZero.InexactFloat64(), withUnset.InexactFloat64(), 1e-12)
}

// TestMixedCurrenciesAreRejectedEverywhere checks the other half of the rule:
// resolving across the fields must not paper over amounts that genuinely
// disagree.
func TestMixedCurrenciesAreRejectedEverywhere(t *testing.T) {
	dollars := money.MustMoneyFromFloat64(1000, money.USD)
	euros := money.MustMoneyFromFloat64(1000, money.EUR)

	period, err := compoundinterest.NewPeriod(decimal.MustFromInt64(12, 0), compoundinterest.Monthly)
	require.NoError(t, err)

	rate, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.01), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	_, err = annuities.New(dollars, euros, unset(), period, rate)
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	simple := simpleinterest.New(dollars, euros, unset(),
		decimal.MustFromFloat64(0.01), simpleinterest.NewPeriod(decimal.MustFromInt64(12, 0), simpleinterest.Months))

	_, err = simple.Future()
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	_, err = simple.Interest()
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	_, err = returns.HoldingPeriodReturn(dollars, euros, unset())
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)
}

// TestResolveCurrency covers the shared helper directly, including the
// all-unset case that has no currency to resolve.
func TestResolveCurrency(t *testing.T) {
	resolved, err := money.ResolveCurrency(unset(), set(100), unset())
	require.NoError(t, err)
	assert.Equal(t, sweepCurrency, resolved)

	resolved, err = money.ResolveCurrency(unset(), unset())
	require.NoError(t, err)
	assert.Equal(t, money.XXX, resolved)

	resolved, err = money.ResolveCurrency()
	require.NoError(t, err)
	assert.Equal(t, money.XXX, resolved)

	_, err = money.ResolveCurrency(
		money.MustMoneyFromFloat64(1, money.USD),
		money.MustMoneyFromFloat64(1, money.JPY))
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)
}
