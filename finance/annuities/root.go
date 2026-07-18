// Package annuities provides functionality for annuity calculations.
package annuities

import (
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

type Annuity struct {
	value             money.Money
	compositeInterest compositeinterest.CompositeInterest
}

func New(value, present, future money.Money, period compositeinterest.Period, rateInterest compositeinterest.RateInterest) (Annuity, error) {
	ci, err := compositeinterest.New(present, future, rateInterest, period)
	if err != nil {
		return Annuity{}, err
	}

	return Annuity{
		value:             value,
		compositeInterest: ci,
	}, nil
}

func (a Annuity) PaymentFromPresentValue() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return money.Money{}, err
	}

	growthFactor := rateInterest.Add(decimal.One)

	growthPower := growthFactor.MustPow(periods)

	numerator := rateInterest.Mul(growthPower)

	denominator := growthPower.Sub(decimal.One)

	annuity := present.MulDecimal(numerator.MustDiv(denominator))

	return annuity, nil
}

func (a Annuity) PaymentFromFutureValue() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return money.Money{}, err
	}

	growthFactor := rateInterest.Add(decimal.One)

	growthPower := growthFactor.MustPow(periods)

	denominator := growthPower.Sub(decimal.One)

	annuity := future.MulDecimal(rateInterest.MustDiv(denominator))

	return annuity, nil
}

// AnticipatePaymentFromPresentValue is like PaymentFromPresentValue, but
// assumes each payment is made at the beginning of its period (annuity due)
// instead of at the end (ordinary annuity).
//
// Formula: PMT = PV × [i(1+i)^n] / {[(1+i)^n - 1] × (1+i)}
// This is PaymentFromPresentValue divided by (1+i): paying one period
// earlier lets a smaller payment reach the same present value.
func (a Annuity) AnticipatePaymentFromPresentValue() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return money.Money{}, err
	}

	growthFactor := rateInterest.Add(decimal.One)

	growthPower := growthFactor.MustPow(periods)

	numerator := rateInterest.Mul(growthPower)

	denominator := growthPower.Sub(decimal.One).Mul(growthFactor)

	annuity := present.MulDecimal(numerator.MustDiv(denominator))

	return annuity, nil
}

// AnticipatePaymentFromFutureValue is like PaymentFromFutureValue, but
// assumes each payment is made at the beginning of its period (annuity due)
// instead of at the end (ordinary annuity).
//
// Formula: PMT = FV × i / {[(1+i)^n - 1] × (1+i)}
// This is PaymentFromFutureValue divided by (1+i).
func (a Annuity) AnticipatePaymentFromFutureValue() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return money.Money{}, err
	}

	growthFactor := rateInterest.Add(decimal.One)

	growthPower := growthFactor.MustPow(periods)

	denominator := growthPower.Sub(decimal.One).Mul(growthFactor)

	annuity := future.MulDecimal(rateInterest.MustDiv(denominator))

	return annuity, nil
}
