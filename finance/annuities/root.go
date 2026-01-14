package annuities

import (
	"math"

	"github.com/yeferson59/gofinance/finance/compositeinterest"
)

type Annuity struct {
	value             float64
	compositeInterest compositeinterest.CompositeInterest
}

func New(value, present, future float64, period compositeinterest.Period, rateInterest compositeinterest.RateInterest) (Annuity, error) {
	compositeinterest, err := compositeinterest.New(present, future, rateInterest, period)
	if err != nil {
		return Annuity{}, err
	}

	return Annuity{
		value:             value,
		compositeInterest: compositeinterest,
	}, nil
}

func (a Annuity) GetWithPresent() (float64, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return 0, err
	}

	annuity := present * (rateInterest * math.Pow(1+rateInterest, periods) / (math.Pow(1+rateInterest, periods) - 1))

	return annuity, nil
}

func (a Annuity) GetWithFuture() (float64, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return 0, err
	}

	annuity := future * (rateInterest / (math.Pow(1+rateInterest, periods) - 1))

	return annuity, nil
}
