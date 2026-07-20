package annuities

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// NewGeneral builds an Annuity for a "general annuity": one whose payment
// frequency differs from the interest rate's compounding frequency (e.g.
// monthly payments on a quarterly-compounded rate). A "simple" annuity, by
// contrast, pays and compounds on the same frequency — that's what New
// builds.
//
// paymentPeriods is the total number of payments, made every
// paymentFrequency. rateInterest carries its own compounding frequency,
// which may be different; GetEqualsRateInterestPeriods (used internally by
// every Annuity calculation) reconciles the two by rebasing the period
// count onto the rate's compounding frequency, so the payment schedule and
// the rate stay consistent without the caller having to convert rates by
// hand.
//
// Example:
//
//	// Monthly payments against a quarterly-compounded 8% nominal rate.
//	rate, _ := compoundinterest.NewRateInterest(
//	    decimal.MustFromFloat64(0.08), compoundinterest.Quarterly, compoundinterest.RateEffectyNominal)
//	annuity, _ := NewGeneral(
//	    money.MustMoneyFromFloat64(100, money.USD),
//	    money.MoneyZero, money.MoneyZero,
//	    24, compoundinterest.Monthly, rate)
func NewGeneral(value, present, future money.Money, paymentPeriods int, paymentFrequency compoundinterest.CompoundingFrequency, rateInterest compoundinterest.RateInterest) (Annuity, error) {
	period, err := compoundinterest.NewPeriod(decimal.MustFromFloat64(float64(paymentPeriods)), paymentFrequency)
	if err != nil {
		return Annuity{}, err
	}

	return New(value, present, future, period, rateInterest)
}
