package simpleinterest

import (
	"github.com/quagmt/udecimal"
	"github.com/yeferson59/gofinance/money"
)

// Future calculates the future value using present value and interest.
// Formula: Future = Present + Interest
// Returns an error if calculation fails, but currently always succeeds.
func (s SimpleInterest) Future() (money.Money, error) {
	future := s.present.Add(s.interest.Decimal)

	return money.Money{Decimal: future}, nil
}

// FutureWithRateInterest calculates the future value using present value, rate, and periods.
// Formula: Future = Present × (1 + Periods × Rate)
// Returns an error if periods is invalid.
func (s SimpleInterest) FutureWithRateInterest() (money.Money, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return money.Money{}, err
	}

	// Step 1: Calculate (Periods × Rate)
	periodRate := numberOfPeriods.Mul(s.rateInterest.Decimal)

	// Step 2: Calculate (1 + Periods × Rate)
	onePlusRate := periodRate.Add(udecimal.One)

	// Step 3: Calculate Present × (1 + Periods × Rate)
	future := s.present.Mul(onePlusRate)

	return money.Money{Decimal: future}, nil
}
