package money_test

import (
	"encoding/json"
	"fmt"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func ExampleMoney_Allocate() {
	// Split a bill three ways. The parts always add back to the total, so no
	// cent is lost or invented by rounding.
	bill := money.MustMoneyFromFloat64(100.01, money.USD)

	parts, err := bill.Allocate(1, 1, 1)
	if err != nil {
		panic(err)
	}

	for _, part := range parts {
		fmt.Println(part.RoundBankString(2))
	}

	// Output:
	// 33.34
	// 33.34
	// 33.33
}

func ExampleMoney_Allocate_byShares() {
	// Split profit between partners holding 50%, 30% and 20%.
	profit := money.MustMoneyFromFloat64(10000, money.USD)

	shares, err := profit.Allocate(5, 3, 2)
	if err != nil {
		panic(err)
	}

	for _, share := range shares {
		fmt.Println(share.RoundBankString(2))
	}

	// Output:
	// 5000.00
	// 3000.00
	// 2000.00
}

func ExampleMoney_MulDecimal() {
	// Applying a rate to an amount: the rate is a plain decimal, so it never
	// needs a placeholder currency.
	principal := money.MustMoneyFromFloat64(1500, money.EUR)
	rate := decimal.MustFromFloat64(0.075)

	interest := principal.MulDecimal(rate)

	code, err := interest.GetCurrency().GetCurrencyISOCode()
	if err != nil {
		panic(err)
	}

	fmt.Println(interest.RoundBankString(2), code)
	// Output: 112.50 EUR
}

func ExampleMoney_MustDivDecimal() {
	// Splitting an amount by a plain count. The Must variant suits a divisor
	// known to be non-zero; DivDecimal returns an error instead.
	total := money.MustMoneyFromFloat64(1200, money.USD)

	perMonth := total.MustDivDecimal(decimal.MustFromInt64(12, 0))

	fmt.Println(perMonth.RoundBankString(2))
	// Output: 100.00
}

func ExampleNewFromDecimal() {
	// FromDecimal is the bridge from a computed decimal back to an amount.
	computed := decimal.MustFromString("1234.5678")

	amount := money.NewFromDecimal(computed, money.USD)

	fmt.Println(amount.RoundBankString(2))
	// Output: 1234.57
}

func ExampleMoney_MarshalJSON() {
	amount := money.MustMoneyFromFloat64(99.95, money.EUR)

	encoded, err := json.Marshal(amount)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(encoded))
	// Output: {"value":"99.95","currency":"EUR"}
}

func ExampleMoney_TryAdd() {
	// TryAdd reports a currency mismatch instead of panicking, which is what
	// Add would do.
	dollars := money.MustMoneyFromFloat64(10, money.USD)
	euros := money.MustMoneyFromFloat64(10, money.EUR)

	if _, err := dollars.TryAdd(euros); err != nil {
		fmt.Println("refused:", err)
	}

	total, err := dollars.TryAdd(money.MustMoneyFromFloat64(5, money.USD))
	if err != nil {
		panic(err)
	}

	fmt.Println(total.RoundBankString(2))
	// Output:
	// refused: money: currency mismatch
	// 15.00
}
