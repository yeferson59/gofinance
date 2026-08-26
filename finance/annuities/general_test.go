package annuities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestNewGeneralAnnuityMatchesSimpleWhenFrequenciesAgree(t *testing.T) {
	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value := money.MustMoneyFromFloat64(1000, money.USD)
	zero := money.MustMoneyFromFloat64(0, money.USD)

	general, err := NewGeneral(value, zero, zero, 12, compoundinterest.Monthly, rateInterest)
	require.NoError(t, err)

	simple, err := NewAnnuity().
		Value(1000, money.USD).
		AnnualRate(0.12).
		Periods(12).
		Monthly().
		Build()
	require.NoError(t, err)

	generalPresent, err := general.Present()
	require.NoError(t, err)
	simplePresent, err := simple.Present()
	require.NoError(t, err)

	assert.InDelta(t, simplePresent.GetDecimal().InexactFloat64(), generalPresent.GetDecimal().InexactFloat64(), 0.01)
}

func TestNewGeneralAnnuityQuarterlyRateMonthlyPayments(t *testing.T) {
	// 24 monthly payments of $100 against a quarterly-compounded 12%
	// nominal rate (1% periodic per quarter). GetEqualsRateInterestPeriods
	// rebases the 24 monthly payments onto the rate's quarterly frequency
	// (8 quarters), so the result matches a simple quarterly annuity of the
	// same payment count and rate.
	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Quarterly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value := money.MustMoneyFromFloat64(100, money.USD)
	zero := money.MustMoneyFromFloat64(0, money.USD)

	general, err := NewGeneral(value, zero, zero, 24, compoundinterest.Monthly, rateInterest)
	require.NoError(t, err)

	generalPresent, err := general.Present()
	require.NoError(t, err)

	quarterlyEquivalent, err := NewAnnuity().
		Value(100, money.USD).
		Quarterly().
		AnnualRate(0.12).
		Periods(8).
		Build()
	require.NoError(t, err)

	quarterlyPresent, err := quarterlyEquivalent.Present()
	require.NoError(t, err)

	assert.InDelta(t, quarterlyPresent.GetDecimal().InexactFloat64(), generalPresent.GetDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityConfigPaymentFrequency(t *testing.T) {
	payment, err := NewAnnuity().
		Present(300000, money.USD).
		Quarterly().
		AnnualRate(0.12).
		PaymentFrequency(compoundinterest.Monthly).
		Periods(24).
		Payment()
	require.NoError(t, err)
	assert.True(t, payment.IsPositive())
}
