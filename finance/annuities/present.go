package annuities

import "math"

func (a Annuity) Present() (float64, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	pow := math.Pow(1+rateInterest, periods)
	present := a.value * ((pow - 1) / (rateInterest * pow))

	return present, nil
}
