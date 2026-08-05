package depreciation_test

import (
	"fmt"

	"github.com/yeferson59/gofinance/v2/finance/depreciation"
	"github.com/yeferson59/gofinance/v2/money"
)

func ExampleStraightLine() {
	// A 10,000 machine with a 1,000 salvage value over five years: the same
	// charge every year.
	schedule, err := depreciation.StraightLine(
		money.MustMoneyFromFloat64(10000, money.USD),
		money.MustMoneyFromFloat64(1000, money.USD),
		5,
	)
	if err != nil {
		panic(err)
	}

	for _, row := range schedule {
		fmt.Printf("year %d: %.2f, book value %.2f\n",
			row.Year, row.Depreciation.InexactFloat64(), row.BookValue.InexactFloat64())
	}

	// Output:
	// year 1: 1800.00, book value 8200.00
	// year 2: 1800.00, book value 6400.00
	// year 3: 1800.00, book value 4600.00
	// year 4: 1800.00, book value 2800.00
	// year 5: 1800.00, book value 1000.00
}

func ExampleDoubleDecliningBalance() {
	// Twice the straight-line rate, switching to a level charge once that
	// deducts more, so the asset still lands exactly on its salvage value.
	schedule, err := depreciation.DoubleDecliningBalance(
		money.MustMoneyFromFloat64(10000, money.USD),
		money.MustMoneyFromFloat64(0, money.USD),
		5,
	)
	if err != nil {
		panic(err)
	}

	for _, row := range schedule {
		fmt.Printf("year %d: %.2f\n", row.Year, row.Depreciation.InexactFloat64())
	}

	// Output:
	// year 1: 4000.00
	// year 2: 2400.00
	// year 3: 1440.00
	// year 4: 1080.00
	// year 5: 1080.00
}

func ExampleMACRS() {
	// The US tax tables ignore salvage and spread a 5-year class over six
	// years, thanks to the half-year convention.
	schedule, err := depreciation.MACRS(money.MustMoneyFromFloat64(100000, money.USD), 5)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d years, first %.2f, last %.2f\n",
		len(schedule),
		schedule[0].Depreciation.InexactFloat64(),
		schedule[len(schedule)-1].Depreciation.InexactFloat64())
	// Output: 6 years, first 20000.00, last 5760.00
}

func ExampleMustSumOfYearsDigits() {
	// Weights of 5/15, 4/15, … over a five-year life.
	schedule := depreciation.MustSumOfYearsDigits(
		money.MustMoneyFromFloat64(10000, money.USD),
		money.MustMoneyFromFloat64(1000, money.USD),
		5,
	)

	fmt.Printf("first %.2f, last %.2f\n",
		schedule[0].Depreciation.InexactFloat64(),
		schedule[len(schedule)-1].Depreciation.InexactFloat64())
	// Output: first 3000.00, last 600.00
}
