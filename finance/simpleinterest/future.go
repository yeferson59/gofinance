package simpleinterest

import "github.com/quagmt/udecimal"

// Future calculates the future value using present value and interest.
// Formula: Future = Present + Interest
// Returns an error if calculation fails, but currently always succeeds.
func (s SimpleInterest) Future() (Decimal, error) {
	future := s.present.Add(s.interest.Decimal)

	return Decimal{future}, nil
}

// FutureWithRateInterest calculates the future value using present value, rate, and periods.
// Formula: Future = Present × (1 + Periods × Rate)
// Returns an error if periods is invalid.
func (s SimpleInterest) FutureWithRateInterest() (Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return Decimal{}, err
	}

	// Step 1: Calculate (Periods × Rate)
	periodRate := numberOfPeriods.Mul(s.rateInterest.Decimal)

	// Step 2: Calculate (1 + Periods × Rate)
	onePlusRate := periodRate.Add(udecimal.One)

	// Step 3: Calculate Present × (1 + Periods × Rate)
	future := s.present.Mul(onePlusRate)

	return Decimal{future}, nil
}
