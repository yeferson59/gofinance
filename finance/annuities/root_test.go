package annuities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/finance/compoundinterest"
	"github.com/yeferson59/gofinance/money"
)

func TestAnnuityAnticipatePaymentFromPresentValue(t *testing.T) {
	// PV = 10000, i = 0.01, n = 12
	// PMT_due = PMT_ordinary / (1+i) = 888.4879 / 1.01 = 879.6910
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	payment, err := annuity.AnticipatePaymentFromPresentValue()
	require.NoError(t, err)
	assert.InDelta(t, 879.6910, payment.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityAnticipatePaymentFromFutureValue(t *testing.T) {
	// FV = 15000, i = 0.01, n = 12
	// PMT_due = PMT_ordinary / (1+i) = 1182.7318 / 1.01 = 1171.0216
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(12), compoundinterest.Monthly)
	require.NoError(t, err)

	rateInterest, err := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.12), compoundinterest.Monthly, compoundinterest.RateEffectyNominal)
	require.NoError(t, err)

	value, err := money.New(0, 2, money.USD)
	require.NoError(t, err)
	present, err := money.New(1000000, 2, money.USD)
	require.NoError(t, err)
	future, err := money.New(1500000, 2, money.USD)
	require.NoError(t, err)
	annuity, err := New(value, present, future, period, rateInterest)
	require.NoError(t, err)

	payment, err := annuity.AnticipatePaymentFromFutureValue()
	require.NoError(t, err)
	assert.InDelta(t, 1171.0216, payment.ToDecimal().InexactFloat64(), 0.01)
}

func TestAnnuityConfigAnticipatePayment(t *testing.T) {
	payment := NewAnnuity().
		Present(300000, money.USD).
		AnnualRate(0.06).
		Periods(360).
		Monthly().
		MustAnticipatePayment()

	ordinary := NewAnnuity().
		Present(300000, money.USD).
		AnnualRate(0.06).
		Periods(360).
		Monthly().
		MustPayment()

	// The anticipated payment should be smaller than the ordinary one, since
	// it's paid one period earlier: PMT_due = PMT_ordinary / (1+i).
	assert.True(t, payment.LessThan(ordinary))
}
