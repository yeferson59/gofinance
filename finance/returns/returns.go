package returns

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// CAGR returns the compound annual growth rate that takes begin to end over
// the given number of periods (typically years):
//
//	CAGR = (end / begin)^(1/periods) − 1
//
// begin and end must share the same currency and begin must be strictly
// positive; periods must be strictly positive. The result is a decimal.Decimal
// fraction (e.g. 0.15 for 15% per period). A negative result is a decline.
//
// It returns money.ErrCurrencyMismatch on mixed currencies,
// ErrNonPositiveValue if begin is not positive, and ErrNonPositivePeriods if
// periods is not positive.
func CAGR(begin, end money.Money, periods decimal.Decimal) (decimal.Decimal, error) {
	if begin.Currency() != end.Currency() {
		return decimal.Decimal{}, money.ErrCurrencyMismatch
	}

	if !begin.IsPositive() {
		return decimal.Decimal{}, ErrNonPositiveValue
	}

	if !periods.IsPos() {
		return decimal.Decimal{}, ErrNonPositivePeriods
	}

	ratio, err := end.ToDecimal().Div(begin.ToDecimal())
	if err != nil {
		return decimal.Decimal{}, err
	}

	exponent, err := decimal.One.Div(periods)
	if err != nil {
		return decimal.Decimal{}, err
	}

	growth, err := ratio.Pow(exponent)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return growth.Sub(decimal.One), nil
}

// MustCAGR is like CAGR but panics on error.
func MustCAGR(begin, end money.Money, periods decimal.Decimal) decimal.Decimal {
	d, err := CAGR(begin, end, periods)
	if err != nil {
		panic(err)
	}

	return d
}

// ROI returns the simple return on investment, the profit relative to the
// amount invested:
//
//	ROI = (final − initial) / initial
//
// initial and final must share the same currency and initial must be
// strictly positive. The result is a decimal.Decimal fraction (e.g. 0.25 for a
// 25% gain, −0.10 for a 10% loss).
//
// It returns money.ErrCurrencyMismatch on mixed currencies and
// ErrNonPositiveValue if initial is not positive.
func ROI(initial, final money.Money) (decimal.Decimal, error) {
	if initial.Currency() != final.Currency() {
		return decimal.Decimal{}, money.ErrCurrencyMismatch
	}

	if !initial.IsPositive() {
		return decimal.Decimal{}, ErrNonPositiveValue
	}

	// TrySub rather than Sub: with amounts at opposite ends of the decimal
	// range the difference can overflow, and a function that returns an error
	// must report that rather than panic.
	profit, err := final.TrySub(initial)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return profit.ToDecimal().Div(initial.ToDecimal())
}

// MustROI is like ROI but panics on error.
func MustROI(initial, final money.Money) decimal.Decimal {
	d, err := ROI(initial, final)
	if err != nil {
		panic(err)
	}

	return d
}

// HoldingPeriodReturn returns the total return over a holding period,
// counting both the change in value and any income (dividends, coupons,
// rent) received:
//
//	HPR = (final − initial + income) / initial
//
// All three amounts must share the same currency and initial must be
// strictly positive. The result is a decimal.Decimal fraction.
//
// It returns money.ErrCurrencyMismatch on mixed currencies and
// ErrNonPositiveValue if initial is not positive.
func HoldingPeriodReturn(initial, final, income money.Money) (decimal.Decimal, error) {
	// ResolveCurrency rather than comparing the three directly: income is
	// optional — a share that paid no dividend — and an unset money.Money
	// carries no currency, so a direct comparison rejected the ordinary case
	// of leaving it out.
	if _, err := money.ResolveCurrency(initial, final, income); err != nil {
		return decimal.Decimal{}, err
	}

	if !initial.IsPositive() {
		return decimal.Decimal{}, ErrNonPositiveValue
	}

	// The Try variants rather than Sub/Add: with amounts at opposite ends of
	// the decimal range the gain can overflow, and a function that returns an
	// error must report that rather than panic.
	gain, err := final.ToDecimal().TrySub(initial.ToDecimal())
	if err != nil {
		return decimal.Decimal{}, err
	}

	total, err := gain.TryAdd(income.ToDecimal())
	if err != nil {
		return decimal.Decimal{}, err
	}

	return total.Div(initial.ToDecimal())
}

// MustHoldingPeriodReturn is like HoldingPeriodReturn but panics on error.
func MustHoldingPeriodReturn(initial, final, income money.Money) decimal.Decimal {
	d, err := HoldingPeriodReturn(initial, final, income)
	if err != nil {
		panic(err)
	}

	return d
}

// Annualized converts a total return accumulated over the given number of
// periods into an equivalent per-period (typically annual) rate:
//
//	annualized = (1 + totalReturn)^(1/periods) − 1
//
// totalReturn is a fraction (e.g. 0.20 for a cumulative 20% gain) and must be
// greater than −1 (you cannot lose more than everything). periods must be
// strictly positive. The result is a decimal.Decimal fraction.
//
// It returns ErrNonPositiveValue if 1 + totalReturn is not positive and
// ErrNonPositivePeriods if periods is not positive.
func Annualized(totalReturn, periods decimal.Decimal) (decimal.Decimal, error) {
	if !periods.IsPos() {
		return decimal.Decimal{}, ErrNonPositivePeriods
	}

	growth := decimal.One.Add(totalReturn)
	if !growth.IsPos() {
		return decimal.Decimal{}, ErrNonPositiveValue
	}

	exponent, err := decimal.One.Div(periods)
	if err != nil {
		return decimal.Decimal{}, err
	}

	annualGrowth, err := growth.Pow(exponent)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return annualGrowth.Sub(decimal.One), nil
}

// MustAnnualized is like Annualized but panics on error.
func MustAnnualized(totalReturn, periods decimal.Decimal) decimal.Decimal {
	d, err := Annualized(totalReturn, periods)
	if err != nil {
		panic(err)
	}

	return d
}
