package bonds

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// Setters given a float with no decimal value report it through every
// terminal method instead of panicking.
func TestInvalidFloatInputsReturnErrors(t *testing.T) {
	nan := math.NaN()

	configs := map[string]Config{
		"Face":        NewBond().Face(nan, money.USD).CouponRate(0.05).Periods(10).Yield(0.06).MarketPrice(950),
		"CouponRate":  NewBond().Face(1000, money.USD).CouponRate(math.Inf(1)).Periods(10).Yield(0.06).MarketPrice(950),
		"Yield":       NewBond().Face(1000, money.USD).CouponRate(0.05).Periods(10).Yield(nan).MarketPrice(950),
		"MarketPrice": NewBond().Face(1000, money.USD).CouponRate(0.05).Periods(10).Yield(0.06).MarketPrice(math.Inf(-1)),
	}

	for name, config := range configs {
		t.Run(name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				calls := map[string]func() error{
					"CouponPayment":    func() error { _, err := config.CouponPayment(); return err },
					"Price":            func() error { _, err := config.Price(); return err },
					"YTM":              func() error { _, err := config.YTM(); return err },
					"MacaulayDuration": func() error { _, err := config.MacaulayDuration(); return err },
					"ModifiedDuration": func() error { _, err := config.ModifiedDuration(); return err },
					"Convexity":        func() error { _, err := config.Convexity(); return err },
				}

				for method, call := range calls {
					err := call()
					assert.True(t, errors.Is(err, decimal.ErrInvalidFormat), "%s error = %v", method, err)
				}
			})
		})
	}
}

func TestDecimalSettersMatchFloatSetters(t *testing.T) {
	viaFloat := NewBond().Face(1000, money.USD).CouponRate(0.05).Periods(10).Yield(0.06)
	viaDecimal := NewBond().
		FaceMoney(money.MustMoneyFromString("1000", money.USD)).
		CouponRateDecimal(decimal.MustFromString("0.05")).
		Periods(10).
		YieldDecimal(decimal.MustFromString("0.06"))

	want, err := viaFloat.Price()
	require.NoError(t, err)

	got, err := viaDecimal.Price()
	require.NoError(t, err)

	assert.True(t, want.Equal(got), "Price: float %s, decimal %s", want, got)

	wantYTM, err := viaFloat.MarketPrice(950).YTM()
	require.NoError(t, err)

	gotYTM, err := viaDecimal.MarketPriceDecimal(decimal.MustFromString("950")).YTM()
	require.NoError(t, err)

	assert.True(t, wantYTM.Equal(gotYTM), "YTM: float %s, decimal %s", wantYTM, gotYTM)
}
