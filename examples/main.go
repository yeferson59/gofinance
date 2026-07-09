package main

import (
	"fmt"

	"github.com/yeferson59/gofinance/finance/annuities"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/finance/simpleinterest"
	"github.com/yeferson59/gofinance/money"
)

func main() {
	compoundExample()

	annuityExample()

	simpleExample()
}

func compoundExample() {
	fmt.Println("=== Composite Interest ===")

	ci := compositeinterest.NewComposite().
		Present(1000, money.USD).
		Rate(0.05).
		Periods(12).
		Monthly().
		RateType(compositeinterest.RateEffectyPeriodic).
		MustBuild()

	future, _ := ci.Future()
	fmt.Println("Present: $1000, Rate: 5%, 12 months")
	fmt.Println("Future value:", future.StringFixed(2))
}

func annuityExample() {
	fmt.Println("\n=== Annuity Payment ===")

	payment := annuities.NewAnnuity().
		Present(300000, money.USD).
		AnnualRate(0.06).
		Periods(360).
		Monthly().
		MustPayment()

	fmt.Println("Loan: $300,000, Rate: 6%, 360 months")
	fmt.Println("Monthly payment:", payment.StringFixed(2))

	schedule, err := annuities.BuildSchedule(
		money.MustMoneyFromFloat64(300000, money.USD),
		money.MustFromFloat64(0.005),
		payment,
		money.MustFromFloat64(360),
	)
	if err != nil {
		fmt.Println("schedule error:", err)
		return
	}
	fmt.Println("Schedule rows:", len(schedule))
}

func simpleExample() {
	fmt.Println("\n=== Simple Interest ===")

	future, _ := simpleinterest.NewSimple().
		Present(5000, money.USD).
		AnnualRate(0.12).
		Periods(18).
		Months().
		FutureValue()

	fmt.Println("Present: $5000, Rate: 12%, 18 months")
	fmt.Println("Future value:", future.StringFixed(2))
}
