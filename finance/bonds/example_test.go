package bonds_test

import (
	"fmt"
	"time"

	"github.com/yeferson59/gofinance/v2/finance/bonds"
	"github.com/yeferson59/gofinance/v2/finance/daycount"
	"github.com/yeferson59/gofinance/v2/money"
)

func ExampleConfig_Price() {
	// A 5-year 5% semiannual bond priced at a 6% yield trades below par,
	// since the yield exceeds the coupon.
	price, err := bonds.NewBond().
		Face(1000, money.USD).
		CouponRate(0.05).
		Frequency(2).
		Periods(10).
		Yield(0.06).
		Price()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", price.InexactFloat64())
	// Output: 957.35
}

func ExampleConfig_YTM() {
	// The yield implied by an observed market price.
	yield, err := bonds.NewBond().
		Face(1000, money.USD).
		CouponRate(0.05).
		Frequency(2).
		Periods(10).
		MarketPrice(957.35).
		YTM()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f\n", yield.InexactFloat64())
	// Output: 0.0600
}

func ExampleConfig_MacaulayDuration() {
	bond := bonds.NewBond().
		Face(1000, money.USD).
		CouponRate(0.05).
		Frequency(2).
		Periods(20).
		Yield(0.05)

	macaulay, err := bond.MacaulayDuration()
	if err != nil {
		panic(err)
	}

	modified, err := bond.ModifiedDuration()
	if err != nil {
		panic(err)
	}

	// Coupons arrive before maturity, so the duration is shorter than the
	// ten-year term.
	fmt.Printf("Macaulay %.4f years, modified %.4f years\n",
		macaulay.InexactFloat64(), modified.InexactFloat64())
	// Output: Macaulay 7.9894 years, modified 7.7946 years
}

func ExampleAccruedInterest() {
	// A buyer settling three months into a six-month coupon period owes the
	// seller half the coupon.
	lastCoupon := time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)
	settlement := time.Date(2024, time.April, 15, 0, 0, 0, 0, time.UTC)
	nextCoupon := time.Date(2024, time.July, 15, 0, 0, 0, 0, time.UTC)

	accrued, err := bonds.AccruedInterest(
		money.MustMoneyFromFloat64(25, money.USD),
		lastCoupon, settlement, nextCoupon,
		daycount.Thirty360,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", accrued.InexactFloat64())
	// Output: 12.50
}
