package bonds

import "github.com/yeferson59/gofinance/v2/decimal"

// MacaulayDuration returns the Macaulay duration in years: the present-value
// weighted average time to the bond's cash flows, discounted at the configured
// yield.
//
//	D = ( Σ t·PV(cashflowₜ) / price ) / frequency
//
// It returns ErrInvalidFrequency, ErrInvalidPeriods, or ErrInvalidYield on
// invalid terms.
func (b Config) MacaulayDuration() (decimal.Decimal, error) {
	y, err := b.periodicYield()
	if err != nil {
		return decimal.Decimal{}, err
	}

	price, sumT, _, err := b.cashflowSums(y)
	if err != nil {
		return decimal.Decimal{}, err
	}

	weightedPeriods, err := sumT.Div(price)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return weightedPeriods.Div(decimal.MustFromInt64(int64(b.freq), 0))
}

// MustMacaulayDuration is like MacaulayDuration but panics on error.
func (b Config) MustMacaulayDuration() decimal.Decimal {
	d, err := b.MacaulayDuration()
	if err != nil {
		panic(err)
	}

	return d
}

// ModifiedDuration returns the modified duration in years: the approximate
// percentage change in price for a one-unit (100%) change in yield. It is the
// Macaulay duration divided by (1 + yield/frequency).
//
// It returns ErrInvalidFrequency, ErrInvalidPeriods, or ErrInvalidYield on
// invalid terms.
func (b Config) ModifiedDuration() (decimal.Decimal, error) {
	macaulay, err := b.MacaulayDuration()
	if err != nil {
		return decimal.Decimal{}, err
	}

	y, err := b.periodicYield()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return macaulay.Div(decimal.One.Add(y))
}

// MustModifiedDuration is like ModifiedDuration but panics on error.
func (b Config) MustModifiedDuration() decimal.Decimal {
	d, err := b.ModifiedDuration()
	if err != nil {
		panic(err)
	}

	return d
}

// Convexity returns the bond's convexity in years², the second-order
// sensitivity of price to yield:
//
//	C = ( Σ t(t+1)·PV(cashflowₜ) / price ) / (1 + y)² / frequency²
//
// where y is the per-period yield.
//
// It returns ErrInvalidFrequency, ErrInvalidPeriods, or ErrInvalidYield on
// invalid terms.
func (b Config) Convexity() (decimal.Decimal, error) {
	y, err := b.periodicYield()
	if err != nil {
		return decimal.Decimal{}, err
	}

	price, _, sumTT, err := b.cashflowSums(y)
	if err != nil {
		return decimal.Decimal{}, err
	}

	weighted, err := sumTT.Div(price)
	if err != nil {
		return decimal.Decimal{}, err
	}

	onePlusSquared := decimal.One.Add(y).Mul(decimal.One.Add(y))

	perPeriod, err := weighted.Div(onePlusSquared)
	if err != nil {
		return decimal.Decimal{}, err
	}

	freqSquared := decimal.MustFromInt64(int64(b.freq*b.freq), 0)

	return perPeriod.Div(freqSquared)
}

// MustConvexity is like Convexity but panics on error.
func (b Config) MustConvexity() decimal.Decimal {
	d, err := b.Convexity()
	if err != nil {
		panic(err)
	}

	return d
}
