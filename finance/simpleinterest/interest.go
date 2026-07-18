package simpleinterest

import (
	"github.com/yeferson59/gofinance/v2/money"
)

// Interest calculates the interest using present value, rate, and periods.
// Formula: Interest = Present × Periods × Rate
// Returns an error if periods is invalid.
func (s SimpleInterest) Interest() (money.Money, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	interest := s.present.MulDecimal(numberOfPeriods).MulDecimal(s.rateInterest)

	return interest, nil
}

func (s SimpleInterest) InterestWithPresentAndFuture() (money.Money, error) {
	interest := s.future.Sub(s.present)

	return interest, nil
}
