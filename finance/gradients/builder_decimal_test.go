package gradients

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func usd(s string) money.Money { return money.MustMoneyFromString(s, money.USD) }

func TestArithmeticDecimalSettersMatchFloatSetters(t *testing.T) {
	want, err := NewArithmeticSeries().FirstPayment(1000, money.USD).Gradient(100, money.USD).AnnualRate(0.12).Periods(10).Present()
	require.NoError(t, err)

	got, err := NewArithmeticSeries().
		FirstPaymentMoney(usd("1000")).
		GradientMoney(usd("100")).
		AnnualRateDecimal(decimal.MustFromString("0.12")).
		Periods(10).
		Present()
	require.NoError(t, err)

	assert.True(t, want.Equal(got), "float %s, decimal %s", want, got)

	periodic, err := NewArithmeticSeries().FirstPaymentMoney(usd("1000")).GradientMoney(usd("100")).RateDecimal(decimal.MustFromString("0.01")).Periods(10).Future()
	require.NoError(t, err)

	periodicFloat, err := NewArithmeticSeries().FirstPayment(1000, money.USD).Gradient(100, money.USD).Rate(0.01).Periods(10).Future()
	require.NoError(t, err)

	assert.True(t, periodic.Equal(periodicFloat), "RateDecimal %s, Rate %s", periodic, periodicFloat)
}

func TestGeometricDecimalSettersMatchFloatSetters(t *testing.T) {
	want, err := NewGeometricSeries().FirstPayment(1000, money.USD).GrowthRate(0.05).AnnualRate(0.12).Periods(10).Present()
	require.NoError(t, err)

	got, err := NewGeometricSeries().
		FirstPaymentMoney(usd("1000")).
		GrowthRateDecimal(decimal.MustFromString("0.05")).
		AnnualRateDecimal(decimal.MustFromString("0.12")).
		Periods(10).
		Present()
	require.NoError(t, err)

	assert.True(t, want.Equal(got), "float %s, decimal %s", want, got)

	periodic, err := NewGeometricSeries().FirstPaymentMoney(usd("1000")).GrowthRateDecimal(decimal.MustFromString("0.05")).RateDecimal(decimal.MustFromString("0.01")).Periods(10).Future()
	require.NoError(t, err)

	periodicFloat, err := NewGeometricSeries().FirstPayment(1000, money.USD).GrowthRate(0.05).Rate(0.01).Periods(10).Future()
	require.NoError(t, err)

	assert.True(t, periodic.Equal(periodicFloat), "RateDecimal %s, Rate %s", periodic, periodicFloat)
}
