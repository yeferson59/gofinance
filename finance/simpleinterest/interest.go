package simpleinterest

import (
	"github.com/yeferson59/gofinance/money"
)

// Interest calculates the interest using present value, rate, and periods.
// Formula: Interest = Present × Periods × Rate
// Returns an error if periods is invalid.
func (s SimpleInterest) Interest() (money.Money, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	interest := s.present.Mul(numberOfPeriods.Decimal).Mul(s.rateInterest.Decimal)

	return money.Money{Decimal: interest}, nil
}

func (s SimpleInterest) InterestWithPresentAndFuture() (money.Money, error) {
	interest := s.future.Sub(s.present.Decimal)

	return money.Money{Decimal: interest}, nil
}
