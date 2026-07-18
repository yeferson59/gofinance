package bonds

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/finance/daycount"
	"github.com/yeferson59/gofinance/money"
)

func usd(amount float64) money.Money {
	return money.MustMoneyFromFloat64(amount, money.USD)
}

func parBond() Config {
	return NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(10)
}

func TestPriceAtPar(t *testing.T) {
	// Coupon rate equals yield → price equals face.
	price, err := parBond().Yield(0.05).Price()
	require.NoError(t, err)
	assert.InDelta(t, 1000.0, price.InexactFloat64(), 0.01)
	assert.Equal(t, money.USD, price.Currency())
}

func TestPriceBelowParWhenYieldExceedsCoupon(t *testing.T) {
	// 5% coupon at a 6% yield → ≈ 957.35.
	price, err := parBond().Yield(0.06).Price()
	require.NoError(t, err)
	assert.InDelta(t, 957.35, price.InexactFloat64(), 0.01)
}

func TestYTMInvertsPrice(t *testing.T) {
	price := parBond().Yield(0.06).MustPrice()

	ytm, err := parBond().MarketPrice(price.InexactFloat64()).YTM()
	require.NoError(t, err)
	assert.InDelta(t, 0.06, ytm.InexactFloat64(), 1e-4)
}

func TestCouponPayment(t *testing.T) {
	coupon, err := parBond().CouponPayment()
	require.NoError(t, err)
	assert.InDelta(t, 25.0, coupon.InexactFloat64(), 1e-9)
}

func TestDurationOrdering(t *testing.T) {
	bond := parBond().Yield(0.05)

	macaulay, err := bond.MacaulayDuration()
	require.NoError(t, err)

	modified, err := bond.ModifiedDuration()
	require.NoError(t, err)

	convexity, err := bond.Convexity()
	require.NoError(t, err)

	// A 5-year bond's duration is under 5 years, and modified < Macaulay.
	assert.Greater(t, macaulay.InexactFloat64(), 4.0)
	assert.Less(t, macaulay.InexactFloat64(), 5.0)
	assert.Less(t, modified.InexactFloat64(), macaulay.InexactFloat64())
	assert.Greater(t, convexity.InexactFloat64(), 0.0)
}

func TestZeroCouponDurationEqualsMaturity(t *testing.T) {
	// A zero-coupon bond's Macaulay duration equals its maturity in years.
	zero := NewBond().Face(1000, money.USD).CouponRate(0).Frequency(2).Periods(10).Yield(0.05)

	macaulay, err := zero.MacaulayDuration()
	require.NoError(t, err)
	assert.InDelta(t, 5.0, macaulay.InexactFloat64(), 1e-9)
}

func TestAccruedInterest(t *testing.T) {
	// Halfway-ish through a semiannual period under 30/360: 90/180 of a 25 coupon.
	coupon := usd(25)
	last := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	settlement := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	next := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)

	accrued, err := AccruedInterest(coupon, last, settlement, next, daycount.Thirty360)
	require.NoError(t, err)
	assert.InDelta(t, 12.5, accrued.InexactFloat64(), 1e-9)
}

func TestBondErrors(t *testing.T) {
	_, err := parBond().Yield(0.05).Frequency(0).Price()
	assert.ErrorIs(t, err, ErrInvalidFrequency)

	_, err = NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(0).Yield(0.05).Price()
	assert.ErrorIs(t, err, ErrInvalidPeriods)

	_, err = parBond().MarketPrice(0).YTM()
	assert.ErrorIs(t, err, ErrNonPositivePrice)
}

func TestMustPanics(t *testing.T) {
	assert.Panics(t, func() { parBond().MarketPrice(0).MustYTM() })
	assert.NotPanics(t, func() { parBond().Yield(0.05).MustPrice() })
}
