package returns_test

import (
	"fmt"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/returns"
	"github.com/yeferson59/gofinance/v2/money"
)

func usd(amount float64) money.Money {
	return money.MustMoneyFromFloat64(amount, money.USD)
}

func ExampleCAGR() {
	// 10,000 grown to 16,105 over five years.
	rate, err := returns.CAGR(usd(10000), usd(16105), decimal.MustFromInt64(5, 0))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f a year\n", rate.InexactFloat64())
	// Output: 0.1000 a year
}

func ExampleHoldingPeriodReturn() {
	// A share bought at 100, sold at 110, having paid 5 in dividends.
	rate, err := returns.HoldingPeriodReturn(usd(100), usd(110), usd(5))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f\n", rate.InexactFloat64())
	// Output: 0.1500
}

func ExampleMustHoldingPeriodReturn() {
	// The Must variants suit configuration known to be valid; they panic
	// rather than returning an error.
	rate := returns.MustHoldingPeriodReturn(usd(100), usd(110), usd(5))

	fmt.Printf("%.4f\n", rate.InexactFloat64())
	// Output: 0.1500
}

func ExampleAnnualized() {
	// A 33.1% return earned over three years.
	rate, err := returns.Annualized(decimal.MustFromFloat64(0.331), decimal.MustFromInt64(3, 0))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f a year\n", rate.InexactFloat64())
	// Output: 0.1000 a year
}

func ExampleMustAnnualized() {
	rate := returns.MustAnnualized(decimal.MustFromFloat64(0.331), decimal.MustFromInt64(3, 0))

	fmt.Printf("%.4f a year\n", rate.InexactFloat64())
	// Output: 0.1000 a year
}

func ExampleRealValue() {
	// What 10,000 twenty years from now is worth in today's money at 2%
	// inflation.
	real, err := returns.RealValue(usd(10000), decimal.MustFromFloat64(0.02), decimal.MustFromInt64(20, 0))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", real.InexactFloat64())
	// Output: 6729.71
}

func ExampleMustNominalValue() {
	// The inverse: what today's 6,729.71 must become to keep its purchasing
	// power over twenty years of 2% inflation.
	nominal := returns.MustNominalValue(usd(6729.71), decimal.MustFromFloat64(0.02), decimal.MustFromInt64(20, 0))

	fmt.Printf("%.2f\n", nominal.InexactFloat64())
	// Output: 10000.00
}

func ExampleRealRate() {
	// The Fisher relation: an 8% nominal return under 3% inflation.
	rate, err := returns.RealRate(decimal.MustFromFloat64(0.08), decimal.MustFromFloat64(0.03))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f\n", rate.InexactFloat64())
	// Output: 0.0485
}

func ExampleSharpeRatio() {
	monthly := []decimal.Decimal{
		decimal.MustFromFloat64(0.02),
		decimal.MustFromFloat64(-0.01),
		decimal.MustFromFloat64(0.03),
		decimal.MustFromFloat64(0.015),
	}

	ratio, err := returns.SharpeRatio(monthly, decimal.MustFromFloat64(0.002))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f\n", ratio.InexactFloat64())
	// Output: 0.6905
}

func ExampleMoneyWeightedReturn() {
	// Contributions are positive: money the investor put in.
	rate, err := returns.MoneyWeightedReturn(
		usd(10000),
		[]money.Money{usd(2000), usd(-1000), usd(0)},
		usd(13000),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f a period\n", rate.InexactFloat64())
	// Output: 0.0427 a period
}
