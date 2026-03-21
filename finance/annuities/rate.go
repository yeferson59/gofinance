package annuities

import "github.com/yeferson59/gofinance/money"

func (a Annuity) Rate() (money.Decimal, money.Decimal, error) {
	return a.compositeInterest.GetEqualsRateInterestPeriods()
}
