package compositeinterest

import (
	"math"
)

func (c *CompositeInterest) Future() (float64, error) {
	periods, rateInterest, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	future := c.present * math.Pow(1+rateInterest, periods)

	return future, nil
}
