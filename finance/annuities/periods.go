package annuities

import "math"

func (a Annuity) PeriodsWithPresent() (float64, error) {
	_, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return 0, err
	}

	logBase := math.Log(1 + rateInterest)
	periods := (math.Log(a.value/(a.value-(present*rateInterest))) / logBase)

	return periods, nil
}

func (a Annuity) PeriodsWithFuture() (float64, error) {
	_, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return 0, err
	}

	logBase := math.Log(1 + rateInterest)
	periods := ((math.Log((rateInterest*future)+a.value) - math.Log(a.value)) / logBase)

	return periods, nil
}
