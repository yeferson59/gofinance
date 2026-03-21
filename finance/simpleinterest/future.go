package simpleinterest

import (
	"github.com/yeferson59/gofinance/money"
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

	onePlusRate := periodRate.Add(money.One)

	future := s.present.Mul(onePlusRate.ToMoney(s.present.Currency()))

	return future, nil
}
