package annuities

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// discountByPeriods divides value by (1+i)^count, the shared step behind the
// deferred-annuity calculations: pushing a value back count extra periods
// without payments.
func (a Annuity) discountByPeriods(value money.Money, count int) (money.Money, error) {
	_, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthPower, err := decimal.One.Add(rateInterest).Pow(decimal.MustFromInt64(int64(count), 0))
	if err != nil {
		return money.Money{}, err
	}

	return value.DivDecimal(growthPower)
}

// PresentDeferred returns the present value of a deferred ordinary annuity:
// one whose payments don't start right away. deferPeriods periods pass with
// no payment (the grace period) before the first of the n regular payments
// (Value, at the end of each period) is made. It discounts the ordinary
// present value back the extra deferPeriods periods:
//
//	PV_deferred = PV_ordinary / (1+i)^deferPeriods
//
// deferPeriods = 0 reduces to Present().
func (a Annuity) PresentDeferred(deferPeriods int) (money.Money, error) {
	present, err := a.Present()
	if err != nil {
		return money.Money{}, err
	}

	return a.discountByPeriods(present, deferPeriods)
}

// AnticipatePresentDeferred is like PresentDeferred, but assumes each
// payment, once the grace period ends, is made at the beginning of its
// period (annuity due) instead of the end.
func (a Annuity) AnticipatePresentDeferred(deferPeriods int) (money.Money, error) {
	present, err := a.AnticipatePresent()
	if err != nil {
		return money.Money{}, err
	}

	return a.discountByPeriods(present, deferPeriods)
}

// PaymentFromPresentValueDeferred returns the fixed periodic payment for a
// deferred ordinary annuity: the payment whose ordinary-annuity present
// value, once discounted back the extra deferPeriods grace periods, equals
// the configured Present value.
//
//	PMT = PV × (1+i)^deferPeriods × i(1+i)^n / [(1+i)^n - 1]
//
// deferPeriods = 0 reduces to PaymentFromPresentValue().
func (a Annuity) PaymentFromPresentValueDeferred(deferPeriods int) (money.Money, error) {
	periods, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	present, err := a.compoundInterest.Present()
	if err != nil {
		return money.Money{}, err
	}

	deferGrowth, err := decimal.One.Add(rateInterest).Pow(decimal.MustFromInt64(int64(deferPeriods), 0))
	if err != nil {
		return money.Money{}, err
	}

	adjustedPresent := present.MulDecimal(deferGrowth)

	factor, err := paymentFactor(rateInterest, periods)
	if err != nil {
		return money.Money{}, err
	}

	return adjustedPresent.MulDecimal(factor), nil
}
