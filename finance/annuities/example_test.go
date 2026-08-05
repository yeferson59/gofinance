package annuities_test

import (
	"fmt"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/annuities"
	"github.com/yeferson59/gofinance/v2/money"
)

func ExampleAnnuityConfig_Payment() {
	// The level payment that amortizes a 300,000 loan at 6% a year over 30
	// years of monthly payments.
	payment, err := annuities.NewAnnuity().
		Present(300000, money.USD).
		AnnualRate(0.06).
		Years(30).
		Monthly().
		Payment()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f a month\n", payment.InexactFloat64())
	// Output: 1798.65 a month
}

func ExampleAnnuityConfig_FutureValue() {
	// Saving 500 a month at 6% a year for ten years.
	future, err := annuities.NewAnnuity().
		Value(500, money.USD).
		AnnualRate(0.06).
		Years(10).
		Monthly().
		FutureValue()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", future.InexactFloat64())
	// Output: 81939.67
}

func ExampleAnnuityConfig_DeferredPresentValue() {
	// A pension of 500 a month for a year that only starts after three
	// months of grace: the ordinary present value pushed back three periods.
	// A grace period of zero reduces to the ordinary present value.
	immediate, err := annuities.NewAnnuity().
		Value(500, money.USD).
		AnnualRate(0.06).
		Periods(12).
		Monthly().
		Defer(0).
		DeferredPresentValue()
	if err != nil {
		panic(err)
	}

	deferred, err := annuities.NewAnnuity().
		Value(500, money.USD).
		AnnualRate(0.06).
		Periods(12).
		Monthly().
		Defer(3).
		DeferredPresentValue()
	if err != nil {
		panic(err)
	}

	fmt.Printf("immediate %.2f, deferred %.2f\n",
		immediate.InexactFloat64(), deferred.InexactFloat64())
	// Output: immediate 5809.47, deferred 5723.19
}

func ExampleAnnuityConfig_AnticipateDeferredPresentValue() {
	// The same grace period, but with each payment falling at the start of
	// its month rather than the end, so every one arrives a month earlier and
	// is discounted less.
	ordinary := annuities.NewAnnuity().
		Value(500, money.USD).
		AnnualRate(0.06).
		Periods(12).
		Monthly().
		Defer(3).
		MustDeferredPresentValue()

	due := annuities.NewAnnuity().
		Value(500, money.USD).
		AnnualRate(0.06).
		Periods(12).
		Monthly().
		Defer(3).
		MustAnticipateDeferredPresentValue()

	fmt.Printf("ordinary %.2f, due %.2f\n", ordinary.InexactFloat64(), due.InexactFloat64())
	// Output: ordinary 5723.19, due 5751.80
}

func ExampleBuildSchedule() {
	// An amortization table: the opening row carries the original balance,
	// and the last one closes at zero.
	payment := annuities.NewAnnuity().
		Present(10000, money.USD).
		Rate(0.01).
		Periods(12).
		Monthly().
		MustPayment()

	schedule, err := annuities.BuildSchedule(
		money.MustMoneyFromFloat64(10000, money.USD),
		decimal.MustFromFloat64(0.01),
		payment,
		decimal.MustFromInt64(12, 0),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d rows, opening %.2f, closing %.2f\n",
		len(schedule),
		schedule[0].Balance.InexactFloat64(),
		schedule[len(schedule)-1].Balance.InexactFloat64())
	// Output: 13 rows, opening 10000.00, closing 0.00
}
