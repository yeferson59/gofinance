package compoundinterest_test

import (
	"fmt"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

func ExampleRateInterest_RateEffectyAnnually() {
	// A 12% nominal rate compounded monthly is really 12.68% a year.
	rate, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.12),
		compoundinterest.Monthly,
		compoundinterest.RateEffectyNominal,
	)
	if err != nil {
		panic(err)
	}

	effective, err := rate.RateEffectyAnnually()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.6f\n", effective.InexactFloat64())
	// Output: 0.126825
}

func ExampleRateInterest_RateAnticipatePeriodic() {
	// The anticipated (discount) rate equivalent to an ordinary 1% a month:
	// d = i/(1+i). Charging it up front costs the borrower the same.
	rate, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.01),
		compoundinterest.Monthly,
		compoundinterest.RateEffectyPeriodic,
	)
	if err != nil {
		panic(err)
	}

	discount, err := rate.RateAnticipatePeriodic()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.6f\n", discount.InexactFloat64())
	// Output: 0.009901
}

func ExampleCompoundInterest_Future() {
	// 10,000 growing at 1% a month for three years.
	period, err := compoundinterest.NewPeriod(decimal.MustFromInt64(36, 0), compoundinterest.Monthly)
	if err != nil {
		panic(err)
	}

	rate, err := compoundinterest.NewRateInterest(
		decimal.MustFromFloat64(0.01),
		compoundinterest.Monthly,
		compoundinterest.RateEffectyPeriodic,
	)
	if err != nil {
		panic(err)
	}

	compound, err := compoundinterest.New(
		money.MustMoneyFromFloat64(10000, money.USD),
		money.MustMoneyFromFloat64(0, money.USD),
		rate, period,
	)
	if err != nil {
		panic(err)
	}

	future, err := compound.Future()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", future.InexactFloat64())
	// Output: 14307.69
}

func ExampleRateInterest_RatePeriodic_anyForm() {
	// Every rate type converts to every other, including across the ordinary
	// and anticipated families.
	forms := []struct {
		name  string
		value float64
		kind  compoundinterest.TypeRate
	}{
		{"periodic", 0.01, compoundinterest.RateEffectyPeriodic},
		{"nominal", 0.12, compoundinterest.RateEffectyNominal},
		{"anticipated periodic", 0.00990099009900990099, compoundinterest.RateAnticipateEffectyPeriodic},
	}

	for _, form := range forms {
		rate, err := compoundinterest.NewRateInterest(
			decimal.MustFromFloat64(form.value), compoundinterest.Monthly, form.kind)
		if err != nil {
			panic(err)
		}

		periodic, err := rate.RatePeriodic()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%-20s -> %.6f\n", form.name, periodic.InexactFloat64())
	}

	// Output:
	// periodic             -> 0.010000
	// nominal              -> 0.010000
	// anticipated periodic -> 0.010000
}
