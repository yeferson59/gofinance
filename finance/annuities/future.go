package annuities

import "math"

func (a Annuity) Future() (float64, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	future := a.value * ((math.Pow(1+rateInterest, periods) - 1) / rateInterest)

	return future, nil
}
