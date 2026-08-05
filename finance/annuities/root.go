// Package annuities provides functionality for annuity calculations.
package annuities

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// Annuity is a series of equal periodic payments (value) evaluated against a
// compound interest configuration: a present value, a future value, a rate,
// and a number of periods.
type Annuity struct {
	value            money.Money
	currency         money.Currency
	compoundInterest compoundinterest.CompoundInterest
}

// New creates an Annuity from the periodic payment (value), the present and
// future values, the number of periods, and the interest rate. Leave present
// or future at zero when they are the unknown being solved for.
//
// The amounts that are set must share one currency; New returns
// money.ErrCurrencyMismatch otherwise, since an annuity denominated in two
// currencies has no meaning.
func New(value, present, future money.Money, period compoundinterest.Period, rateInterest compoundinterest.RateInterest) (Annuity, error) {
	currency, err := resolveCurrency(value, present, future)
	if err != nil {
		return Annuity{}, err
	}

	ci, err := compoundinterest.New(present, future, rateInterest, period)
	if err != nil {
		return Annuity{}, err
	}

	return Annuity{
		value:            value,
		currency:         currency,
		compoundInterest: ci,
	}, nil
}

// resolveCurrency returns the single currency the given amounts are expressed
// in, ignoring those that were never set.
//
// An annuity is often described by only some of its three amounts — one given
// a present value alone has no periodic payment — and an unset money.Money
// carries money.XXX, the ISO code for "no currency". Deriving every result
// from one particular field therefore produced XXX whenever that field was the
// unset one, and adding such a result to an amount in a real currency panicked
// with a currency mismatch.
//
// It returns money.ErrCurrencyMismatch when two amounts that are set disagree.
func resolveCurrency(amounts ...money.Money) (money.Currency, error) {
	resolved := money.XXX

	for _, amount := range amounts {
		currency := amount.Currency()
		if currency == money.XXX {
			continue
		}

		if resolved != money.XXX && currency != resolved {
			return money.XXX, money.ErrCurrencyMismatch
		}

		resolved = currency
	}

	return resolved, nil
}

// paymentFactor returns i(1+i)^n / [(1+i)^n - 1], the factor turning a present
// value into the fixed payment of an ordinary annuity.
//
// At a zero rate the factor degenerates to 1/n: with no interest to service,
// the present value is simply split evenly across the periods. Computing it
// through the general formula would divide by zero, so the limit is returned
// directly.
func paymentFactor(rate, periods decimal.Decimal) (decimal.Decimal, error) {
	if rate.IsZero() {
		return decimal.One.Div(periods)
	}

	growthPower, err := rate.Add(decimal.One).Pow(periods)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return rate.Mul(growthPower).Div(growthPower.Sub(decimal.One))
}

// sinkingFundFactor returns i / [(1+i)^n - 1], the factor turning a future
// value into the fixed payment of an ordinary annuity. Like paymentFactor it
// degenerates to 1/n at a zero rate.
func sinkingFundFactor(rate, periods decimal.Decimal) (decimal.Decimal, error) {
	if rate.IsZero() {
		return decimal.One.Div(periods)
	}

	growthPower, err := rate.Add(decimal.One).Pow(periods)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return rate.Div(growthPower.Sub(decimal.One))
}

// dueFactor converts an ordinary-annuity factor into its annuity-due
// counterpart by dividing by (1+i): paying one period earlier lets a smaller
// payment reach the same value.
func dueFactor(factor, rate decimal.Decimal) (decimal.Decimal, error) {
	return factor.Div(rate.Add(decimal.One))
}

// PaymentFromPresentValue returns the fixed periodic payment that amortizes
// the configured present value over the configured number of periods, with
// each payment made at the end of its period (ordinary annuity).
//
//	PMT = PV × i(1+i)^n / [(1+i)^n - 1]
//
// At a zero rate this reduces to PV/n.
func (a Annuity) PaymentFromPresentValue() (money.Money, error) {
	periods, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	present, err := a.compoundInterest.Present()
	if err != nil {
		return money.Money{}, err
	}

	factor, err := paymentFactor(rateInterest, periods)
	if err != nil {
		return money.Money{}, err
	}

	return present.MulDecimal(factor), nil
}

// PaymentFromFutureValue returns the fixed periodic payment that accumulates
// to the configured future value over the configured number of periods, with
// each payment made at the end of its period (ordinary annuity).
//
//	PMT = FV × i / [(1+i)^n - 1]
//
// At a zero rate this reduces to FV/n.
func (a Annuity) PaymentFromFutureValue() (money.Money, error) {
	periods, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	future, err := a.compoundInterest.Future()
	if err != nil {
		return money.Money{}, err
	}

	factor, err := sinkingFundFactor(rateInterest, periods)
	if err != nil {
		return money.Money{}, err
	}

	return future.MulDecimal(factor), nil
}

// AnticipatePaymentFromPresentValue is like PaymentFromPresentValue, but
// assumes each payment is made at the beginning of its period (annuity due)
// instead of at the end (ordinary annuity).
//
// Formula: PMT = PV × [i(1+i)^n] / {[(1+i)^n - 1] × (1+i)}
// This is PaymentFromPresentValue divided by (1+i): paying one period
// earlier lets a smaller payment reach the same present value.
//
// At a zero rate this reduces to PV/n, the same as the ordinary annuity:
// with no interest, when the payment falls inside the period makes no
// difference.
func (a Annuity) AnticipatePaymentFromPresentValue() (money.Money, error) {
	periods, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	present, err := a.compoundInterest.Present()
	if err != nil {
		return money.Money{}, err
	}

	ordinary, err := paymentFactor(rateInterest, periods)
	if err != nil {
		return money.Money{}, err
	}

	factor, err := dueFactor(ordinary, rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	return present.MulDecimal(factor), nil
}

// AnticipatePaymentFromFutureValue is like PaymentFromFutureValue, but
// assumes each payment is made at the beginning of its period (annuity due)
// instead of at the end (ordinary annuity).
//
// Formula: PMT = FV × i / {[(1+i)^n - 1] × (1+i)}
// This is PaymentFromFutureValue divided by (1+i).
//
// At a zero rate this reduces to FV/n.
func (a Annuity) AnticipatePaymentFromFutureValue() (money.Money, error) {
	periods, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	future, err := a.compoundInterest.Future()
	if err != nil {
		return money.Money{}, err
	}

	ordinary, err := sinkingFundFactor(rateInterest, periods)
	if err != nil {
		return money.Money{}, err
	}

	factor, err := dueFactor(ordinary, rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	return future.MulDecimal(factor), nil
}
