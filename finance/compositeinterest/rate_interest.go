package compositeinterest

import (
	"math"
)

func (c *CompositeInterest) Interest() (float64, error) {
	periods, _, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	rateInterest := math.Pow((c.future/c.present), (1/periods)) - 1

	return rateInterest, nil
}
