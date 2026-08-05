package simpleinterest_test

import (
	"fmt"

	"github.com/yeferson59/gofinance/v2/finance/simpleinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

func ExampleSimpleConfig_FutureValue() {
	// Simple interest does not compound: 1,000 at 6% a year for two years
	// reaches 1,120, where compound interest would reach 1,123.60.
	future, err := simpleinterest.NewSimple().
		Present(1000, money.USD).
		AnnualRate(0.06).
		Periods(2).
		Years().
		FutureValue()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", future.InexactFloat64())
	// Output: 1120.00
}

func ExampleSimpleConfig_PresentValue() {
	// What must be invested today at 6% simple interest to hold 1,120 in two
	// years: the inverse of the example above.
	present, err := simpleinterest.NewSimple().
		Future(1120, money.USD).
		AnnualRate(0.06).
		Periods(2).
		Years().
		PresentValue()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", present.InexactFloat64())
	// Output: 1000.00
}

func ExampleSimpleInterest_Interest() {
	// The interest alone, through the plain constructor rather than the
	// builder.
	simple := simpleinterest.NewSimple().
		Present(1000, money.USD).
		AnnualRate(0.06).
		Periods(2).
		Years().
		Build()

	interest, err := simple.Interest()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", interest.InexactFloat64())
	// Output: 120.00
}

func ExampleSimpleInterest_RateInterest() {
	// Recover the rate from a known principal, interest and term.
	simple := simpleinterest.NewSimple().
		Present(1000, money.USD).
		Interest(120, money.USD).
		Periods(2).
		Years().
		Build()

	rate, err := simple.RateInterest()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f a year\n", rate.InexactFloat64())
	// Output: 0.0600 a year
}
