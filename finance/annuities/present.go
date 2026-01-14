package annuities

import "math"

func (a Annuity) Present() (float64, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	present := a.value * ((math.Pow(1+rateInterest, periods) - 1) / (rateInterest * math.Pow(1+rateInterest, periods)))

	return present, nil
}
