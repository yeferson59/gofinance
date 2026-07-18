package simpleinterest

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func (s SimpleInterest) Future() (money.Money, error) {
	future := s.present.Add(s.interest)

	return future, nil
}

func (s SimpleInterest) FutureWithRateInterest() (money.Money, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	periodRate := numberOfPeriods.Mul(s.rateInterest)

	onePlusRate := periodRate.Add(decimal.One)

	future := s.present.MulDecimal(onePlusRate)

	return future, nil
}
