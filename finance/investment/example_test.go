package investment_test

import (
	"fmt"
	"time"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/investment"
	"github.com/yeferson59/gofinance/v2/money"
)

func usd(amount float64) money.Money {
	return money.MustMoneyFromFloat64(amount, money.USD)
}

func ExampleNPV() {
	// Invest 1,000 now and collect 400 for three years. At a 10% discount
	// rate the project is slightly value-destroying.
	value, err := investment.NPV(
		decimal.MustFromFloat64(0.10),
		[]money.Money{usd(-1000), usd(400), usd(400), usd(400)},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", value.InexactFloat64())
	// Output: -5.26
}

func ExampleIRR() {
	// The rate at which those same flows break even.
	rate, err := investment.IRR([]money.Money{usd(-1000), usd(400), usd(400), usd(400)})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f\n", rate.InexactFloat64())
	// Output: 0.0970
}

func ExampleIRR_noSignChange() {
	// A series that only ever pays out has no break-even rate, and says so.
	_, err := investment.IRR([]money.Money{usd(-1000), usd(-400)})

	fmt.Println(err)
	// Output: investment: cash flows must contain at least one sign change for an IRR to exist
}

func ExampleXIRR() {
	// Flows on real dates rather than even periods, measured Actual/365.
	base := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

	rate, err := investment.XIRR([]investment.DatedCashFlow{
		{Date: base, Amount: usd(-10000)},
		{Date: base.AddDate(0, 3, 0), Amount: usd(3000)},
		{Date: base.AddDate(0, 9, 0), Amount: usd(4000)},
		{Date: base.AddDate(1, 0, 0), Amount: usd(4000)},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f a year\n", rate.InexactFloat64())
	// Output: 0.1460 a year
}

func ExamplePerpetuity() {
	// A payment of 100 forever, discounted at 5%.
	value, err := investment.Perpetuity(usd(100), decimal.MustFromFloat64(0.05))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", value.InexactFloat64())
	// Output: 2000.00
}

func ExampleGrowingPerpetuity() {
	// The Gordon model: a payment of 100 growing 2% a year, discounted at 5%.
	value, err := investment.GrowingPerpetuity(
		usd(100), decimal.MustFromFloat64(0.05), decimal.MustFromFloat64(0.02))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", value.InexactFloat64())
	// Output: 3333.33
}
