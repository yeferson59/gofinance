package annuities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// TestFutureFromPresentValueAlone is the regression test for the panic a user
// hit computing the future value of an annuity described by a rate, a term and
// a present value — with no periodic payment.
//
// Every result took its currency from the periodic payment, and an annuity
// with no payment leaves that field as the zero money.Money, which carries
// money.XXX ("no currency"). FutureWithContributions then added contributions
// denominated in XXX to a principal denominated in USD, and Money.Add panics
// on a currency mismatch. The annuity now resolves one currency from whichever
// of its amounts are set.
func TestFutureFromPresentValueAlone(t *testing.T) {
	for _, currency := range []money.Currency{money.USD, money.EUR, money.JPY} {
		t.Run(currency.String(), func(t *testing.T) {
			config := NewAnnuity().
				Present(1000, currency).
				AnnualRate(0.12).
				Periods(12).
				Monthly()

			assert.NotPanics(t, func() {
				future, err := config.FutureValue()
				require.NoError(t, err)

				// 1000 × 1.01^12 = 1126.8250
				assert.InDelta(t, 1126.8250, future.InexactFloat64(), 0.0001)
				assert.Equal(t, currency, future.Currency())
			})

			assert.NotPanics(t, func() {
				future, err := config.AnticipateFutureValue()
				require.NoError(t, err)
				assert.Equal(t, currency, future.Currency())
			})
		})
	}
}

// TestPresentFromFutureValueAlone covers the mirror case: an annuity described
// by a future value alone, where the payment is the unset field.
func TestPresentFromFutureValueAlone(t *testing.T) {
	config := NewAnnuity().
		Future(2000, money.EUR).
		AnnualRate(0.12).
		Periods(12).
		Monthly()

	assert.NotPanics(t, func() {
		present, err := config.PresentValue()
		require.NoError(t, err)
		assert.Equal(t, money.EUR, present.Currency())
	})
}

// TestResultsCarryTheConfiguredCurrency sweeps the whole calculation API of a
// partially configured annuity, checking every result comes back in the
// currency the caller set rather than in XXX.
func TestResultsCarryTheConfiguredCurrency(t *testing.T) {
	annuity, err := New(
		money.Money{}, // no periodic payment
		money.MustMoneyFromFloat64(1000, money.EUR),
		money.Money{},
		mustPeriod(t, 12),
		mustRate(t, 0.01),
	)
	require.NoError(t, err)

	results := map[string]func() (money.Money, error){
		"Future":                            annuity.Future,
		"AnticipateFuture":                  annuity.AnticipateFuture,
		"PrincipalFuture":                   annuity.PrincipalFuture,
		"FutureWithContributions":           annuity.FutureWithContributions,
		"AnticipateFutureWithContributions": annuity.AnticipateFutureWithContributions,
		"Present":                           annuity.Present,
		"AnticipatePresent":                 annuity.AnticipatePresent,
	}

	for name, compute := range results {
		t.Run(name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				result, err := compute()
				require.NoError(t, err)
				assert.Equal(t, money.EUR, result.Currency(),
					"%s returned %v", name, result.Currency())
			})
		})
	}
}

// TestMixedCurrenciesAreRejected checks an annuity cannot be built from
// amounts in different currencies, which has no meaning and used to surface
// much later as a panic deep in a calculation.
func TestMixedCurrenciesAreRejected(t *testing.T) {
	mixtures := []struct {
		name                   string
		value, present, future money.Money
	}{
		{
			"payment and present disagree",
			money.MustMoneyFromFloat64(100, money.EUR),
			money.MustMoneyFromFloat64(1000, money.USD),
			money.Money{},
		},
		{
			"present and future disagree",
			money.Money{},
			money.MustMoneyFromFloat64(1000, money.USD),
			money.MustMoneyFromFloat64(2000, money.JPY),
		},
		{
			"payment and future disagree",
			money.MustMoneyFromFloat64(100, money.GBP),
			money.Money{},
			money.MustMoneyFromFloat64(2000, money.USD),
		},
	}

	for _, mixture := range mixtures {
		t.Run(mixture.name, func(t *testing.T) {
			_, err := New(mixture.value, mixture.present, mixture.future,
				mustPeriod(t, 12), mustRate(t, 0.01))
			assert.ErrorIs(t, err, money.ErrCurrencyMismatch)
		})
	}

	// The builder surfaces it the same way.
	_, err := NewAnnuity().
		Present(1000, money.USD).
		Value(100, money.EUR).
		AnnualRate(0.12).
		Periods(12).
		Monthly().
		FutureValue()
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)
}

// TestAllZeroAmountsResolveToNoCurrency checks the degenerate configuration:
// with nothing set there is no currency to resolve, and the results say so
// rather than inventing one.
func TestAllZeroAmountsResolveToNoCurrency(t *testing.T) {
	annuity, err := New(money.Money{}, money.Money{}, money.Money{},
		mustPeriod(t, 12), mustRate(t, 0.01))
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		future, err := annuity.Future()
		require.NoError(t, err)
		assert.Equal(t, money.XXX, future.Currency())
		assert.True(t, future.IsZero())
	})
}

// mustPeriod and mustRate build the monthly period and periodic rate the tests
// above share.
func mustPeriod(t *testing.T, periods float64) compoundinterest.Period {
	t.Helper()

	period, err := compoundinterest.NewPeriod(
		decimal.MustFromFloat64(periods), compoundinterest.Monthly)
	require.NoError(t, err)

	return period
}

func mustRate(t *testing.T, rate float64) compoundinterest.RateInterest {
	t.Helper()

	rateInterest, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(rate), compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
	require.NoError(t, err)

	return rateInterest
}
