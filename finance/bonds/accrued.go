package bonds

import (
	"time"

	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/finance/daycount"
	"github.com/yeferson59/gofinance/money"
)

// AccruedInterest returns the coupon interest accrued from the last coupon date
// up to the settlement date, the portion of the current period's coupon the
// buyer owes the seller:
//
//	accrued = couponPerPeriod × days(last, settlement) / days(last, next)
//
// where both day counts use the given convention. Adding accrued interest to a
// bond's clean price gives its dirty (invoice) price.
//
// It returns daycount errors for out-of-order dates and ErrInvalidPeriods if
// the coupon period has zero length.
func AccruedInterest(couponPerPeriod money.Money, lastCoupon, settlement, nextCoupon time.Time, conv daycount.Convention) (money.Money, error) {
	accruedDays, err := daycount.Days(lastCoupon, settlement, conv)
	if err != nil {
		return money.Money{}, err
	}

	periodDays, err := daycount.Days(lastCoupon, nextCoupon, conv)
	if err != nil {
		return money.Money{}, err
	}

	if periodDays == 0 {
		return money.Money{}, ErrInvalidPeriods
	}

	fraction, err := decimal.MustFromInt64(int64(accruedDays), 0).Div(decimal.MustFromInt64(int64(periodDays), 0))
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(couponPerPeriod.ToDecimal().Mul(fraction), couponPerPeriod.Currency()), nil
}
