package daycount_test

import (
	"fmt"
	"time"

	"github.com/yeferson59/gofinance/v2/finance/daycount"
)

func ExampleDays() {
	// The same six months measured under each convention. 30/360 counts
	// idealised months; the Actual conventions count real calendar days.
	start := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, time.July, 1, 0, 0, 0, 0, time.UTC)

	for _, convention := range []daycount.Convention{
		daycount.Thirty360, daycount.Actual360, daycount.Actual365Fixed,
	} {
		days, err := daycount.Days(start, end, convention)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%-20s %d days\n", convention, days)
	}

	// Output:
	// 30/360               180 days
	// Actual/360           182 days
	// Actual/365 Fixed     182 days
}

func ExampleDays_endOfFebruary() {
	// The US 30/360 convention treats the last day of February as the 30th,
	// so a period from a February month end to another month end measures a
	// whole number of 30-day months.
	days, err := daycount.Days(
		time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2024, time.August, 31, 0, 0, 0, 0, time.UTC),
		daycount.Thirty360,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(days)
	// Output: 180
}

func ExampleYearFraction() {
	// Actual/Actual ISDA splits the period at the calendar-year boundary and
	// divides each part by its own year length, so a leap year counts 366.
	fraction, err := daycount.YearFraction(
		time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		daycount.ActualActualISDA,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f\n", fraction.InexactFloat64())
	// Output: 1.0000
}
