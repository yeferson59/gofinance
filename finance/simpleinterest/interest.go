package simpleinterest

// Interest calculates the interest using present value, rate, and periods.
// Formula: Interest = Present × Periods × Rate
// Returns an error if periods is invalid.
func (s SimpleInterest) Interest() (Decimal, error) {
	numberOfPeriods, err := s.periods.getPeriod()
	if err != nil {
		return Decimal{}, err
	}

	interest := s.present.Mul(numberOfPeriods.Decimal).Mul(s.rateInterest.Decimal)

	return Decimal{interest}, nil
}

func (s SimpleInterest) InterestWithPresentAndFuture() (Decimal, error) {
	interest := s.future.Sub(s.present.Decimal)

	return Decimal{interest}, nil
}
