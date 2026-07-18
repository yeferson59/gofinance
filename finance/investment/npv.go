package investment

import "github.com/yeferson59/gofinance/money"

// NPV returns the net present value of cashFlows discounted at the given
// periodic rate:
//
//	NPV = Σ CFₜ / (1 + rate)ᵗ   for t = 0, 1, … , n−1
//
// The flow at index 0 is not discounted. cashFlows must be non-empty and all
// in the same currency; rate must be greater than −1. The result carries the
// cash flows' currency.
//
// It returns ErrNoCashFlows for an empty slice, money.ErrCurrencyMismatch on
// mixed currencies, and ErrInvalidRate if rate ≤ −1.
func NPV(rate money.Decimal, cashFlows []money.Money) (money.Money, error) {
	amounts, currency, err := decimalFlows(cashFlows)
	if err != nil {
		return money.Money{}, err
	}

	value, err := npvDecimal(rate, amounts)
	if err != nil {
		return money.Money{}, err
	}

	return value.ToMoney(currency), nil
}

// MustNPV is like NPV but panics on error.
func MustNPV(rate money.Decimal, cashFlows []money.Money) money.Money {
	m, err := NPV(rate, cashFlows)
	if err != nil {
		panic(err)
	}

	return m
}

// npvDecimal computes the net present value of the already-decimal amounts at
// the given periodic rate. The discount factor (1+rate)ᵗ is built
// incrementally so each period costs one multiplication and one division.
func npvDecimal(rate money.Decimal, amounts []money.Decimal) (money.Decimal, error) {
	onePlus := money.One.Add(rate)
	if !onePlus.IsPos() {
		return money.Decimal{}, ErrInvalidRate
	}

	// t = 0 is undiscounted.
	sum := amounts[0]
	factor := onePlus // (1+rate)^t for the current t, starting at t = 1

	for t := 1; t < len(amounts); t++ {
		discounted, err := amounts[t].Div(factor)
		if err != nil {
			return money.Decimal{}, err
		}

		sum = sum.Add(discounted)
		factor = factor.Mul(onePlus)
	}

	return sum, nil
}
