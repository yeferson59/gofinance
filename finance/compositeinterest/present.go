package compositeinterest

import "math"

func (c *CompositeInterest) Present() (float64, error) {
	periods, rateInterest, err := c.GetEqualsRateInterestPeriods()
	if err != nil {
		return 0, err
	}

	present := (c.future / math.Pow(1+rateInterest, periods))

	return present, nil
}
