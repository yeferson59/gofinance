package compositeinterest

import "math"

func (c *CompositeInterest) Periods() (float64, error) {
	_, rateInterest, err := c.getEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	periods := (math.Log((c.future / c.present)) / math.Log(1+rateInterest))

	return periods, nil
}
