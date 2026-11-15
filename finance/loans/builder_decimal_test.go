package loans

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestDecimalSettersMatchFloatSetters(t *testing.T) {
	viaFloat := mortgage().Fees(3500).ExtraPayment(200)
	viaDecimal := NewLoan().
		PrincipalMoney(money.MustMoneyFromString("250000", money.USD)).
		AnnualRateDecimal(decimal.MustFromString("0.065")).
		Years(30).
		Monthly().
		FeesDecimal(decimal.MustFromString("3500")).
		ExtraPaymentDecimal(decimal.MustFromString("200"))

	for name, pair := range map[string][2]func() (money.Money, error){
		"Payment":     {viaFloat.Payment, viaDecimal.Payment},
		"NetProceeds": {viaFloat.NetProceeds, viaDecimal.NetProceeds},
	} {
		want, err := pair[0]()
		require.NoError(t, err, name)

		got, err := pair[1]()
		require.NoError(t, err, name)

		assert.True(t, want.Equal(got), "%s: float %s, decimal %s", name, want, got)
	}

	wantAPR, err := viaFloat.APR()
	require.NoError(t, err)

	gotAPR, err := viaDecimal.APR()
	require.NoError(t, err)

	assert.True(t, wantAPR.Equal(gotAPR), "APR: float %s, decimal %s", wantAPR, gotAPR)

	periodic, err := NewLoan().PrincipalMoney(money.MustMoneyFromString("1000", money.USD)).RateDecimal(decimal.MustFromString("0.01")).Periods(12).Payment()
	require.NoError(t, err)

	periodicFloat, err := NewLoan().Principal(1000, money.USD).Rate(0.01).Periods(12).Payment()
	require.NoError(t, err)

	assert.True(t, periodic.Equal(periodicFloat), "RateDecimal %s, Rate %s", periodic, periodicFloat)
}
