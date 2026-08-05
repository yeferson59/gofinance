package gradients_test

import (
	"fmt"

	"github.com/yeferson59/gofinance/v2/finance/gradients"
	"github.com/yeferson59/gofinance/v2/money"
)

func ExampleArithmeticConfig() {
	// Maintenance costs starting at 1,000 and rising 100 a year for five
	// years, discounted at 10%.
	present, err := gradients.NewArithmeticSeries().
		FirstPayment(1000, money.USD).
		Gradient(100, money.USD).
		Rate(0.10).
		Periods(5).
		Annually().
		Present()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", present.InexactFloat64())
	// Output: 4476.97
}

func ExampleGeometricConfig() {
	// A salary starting at 50,000 and growing 4% a year for ten years,
	// discounted at 8%.
	present, err := gradients.NewGeometricSeries().
		FirstPayment(50000, money.USD).
		GrowthRate(0.04).
		Rate(0.08).
		Periods(10).
		Annually().
		Present()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", present.InexactFloat64())
	// Output: 392950.61
}

func ExampleArithmeticConfig_zeroRate() {
	// With no interest the series is just the sum of its payments:
	// 1000 + 1100 + 1200 + 1300 + 1400.
	present, err := gradients.NewArithmeticSeries().
		FirstPayment(1000, money.USD).
		Gradient(100, money.USD).
		Rate(0).
		Periods(5).
		Annually().
		Present()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", present.InexactFloat64())
	// Output: 6000.00
}
