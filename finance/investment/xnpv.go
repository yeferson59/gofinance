package investment

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// XNPV returns the net present value of a series of cash flows that occur on
// specific dates, discounted at the given annual rate:
//
//	XNPV = Σ amountᵢ / (1 + rate)^(daysᵢ/365)
//
// where daysᵢ is the number of days from the first flow's date to flow i,
// measured with the Actual/365 Fixed convention. The first flow is not
// discounted. flows must be non-empty and all in the same currency; rate must
// be greater than −1. The result carries the flows' currency.
//
// It returns ErrNoCashFlows for an empty slice, money.ErrCurrencyMismatch on
// mixed currencies, ErrDatesBeforeBase if a date precedes the first, and
// ErrInvalidRate if rate ≤ −1.
func XNPV(rate decimal.Decimal, flows []DatedCashFlow) (money.Money, error) {
	amounts, times, currency, err := datedFlows(flows)
	if err != nil {
		return money.Money{}, err
	}

	value, err := xnpvDecimal(rate, amounts, times)
	if err != nil {
		return money.Money{}, err
	}

	return money.NewFromDecimal(value, currency), nil
}

// MustXNPV is like XNPV but panics on error.
func MustXNPV(rate decimal.Decimal, flows []DatedCashFlow) money.Money {
	m, err := XNPV(rate, flows)
	if err != nil {
		panic(err)
	}

	return m
}

// xnpvDecimal computes the date-based net present value of the amounts at the
// given annual rate, using the precomputed year offsets in times. Each flow is
// discounted by (1+rate) raised to its (possibly fractional) year offset.
func xnpvDecimal(rate decimal.Decimal, amounts, times []decimal.Decimal) (decimal.Decimal, error) {
	onePlus := decimal.One.Add(rate)
	if !onePlus.IsPos() {
		return decimal.Decimal{}, ErrInvalidRate
	}

	sum := decimal.Zero

	for i, amount := range amounts {
		term, err := discountToBase(amount, onePlus, times[i])
		if err != nil {
			return decimal.Decimal{}, err
		}

		// The running sum can overflow at an extreme rate, so report that
		// instead of panicking inside a function that returns an error.
		sum, err = sum.TryAdd(term)
		if err != nil {
			return decimal.Decimal{}, err
		}
	}

	return sum, nil
}

// discountToBase discounts amount back to the base date by (1+rate) raised to
// its year offset. A flow at the base date (offset zero) is returned as is.
func discountToBase(amount, onePlus, years decimal.Decimal) (decimal.Decimal, error) {
	if years.IsZero() {
		return amount, nil
	}

	factor, err := onePlus.Pow(years)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return amount.Div(factor)
}
