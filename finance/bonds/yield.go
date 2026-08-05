package bonds

import "github.com/yeferson59/gofinance/v2/decimal"

const maxYieldBisectIter = 200

var (
	yieldBracketTol = decimal.MustFromString("0.0000000001") // 1e-10
	yieldHalf       = decimal.MustFromFloat64(0.5)
)

// YTM returns the annual yield to maturity implied by the configured clean
// price: the yield at which the bond's discounted cash flows equal that price.
// Since price falls monotonically with yield, the root is found by scanning for
// a sign change in (price(yield) − target) and bisecting.
//
// It returns ErrNonPositivePrice if the price is not positive,
// ErrInvalidFrequency or ErrInvalidPeriods on invalid terms, and
// ErrNoConvergence if no yield reproduces the price.
func (b Config) YTM() (decimal.Decimal, error) {
	if !b.price.IsPos() {
		return decimal.Decimal{}, ErrNonPositivePrice
	}

	var (
		prevYield decimal.Decimal
		prevDiff  decimal.Decimal
		bracketed bool
	)

	for _, curYield := range yieldCandidates() {
		curDiff, err := b.priceDiff(curYield)
		if err != nil {
			// At the top of the candidate range the discount factor of a
			// long bond overflows. That says nothing about where the yield
			// is, so skip the candidate and start a fresh pair after the gap
			// rather than abandoning the search.
			bracketed = false

			continue
		}

		if curDiff.IsZero() {
			return curYield, nil
		}

		if bracketed && prevDiff.Sign() != curDiff.Sign() {
			return b.bisectYield(prevYield, curYield, prevDiff)
		}

		prevYield, prevDiff, bracketed = curYield, curDiff, true
	}

	return decimal.Decimal{}, ErrNoConvergence
}

// MustYTM is like YTM but panics on error.
func (b Config) MustYTM() decimal.Decimal {
	d, err := b.YTM()
	if err != nil {
		panic(err)
	}

	return d
}

// priceDiff returns the bond's price at the given annual yield minus the target
// price. Its root is the yield to maturity.
func (b Config) priceDiff(annualYield decimal.Decimal) (decimal.Decimal, error) {
	y, err := annualYield.Div(decimal.MustFromInt64(int64(b.freq), 0))
	if err != nil {
		return decimal.Decimal{}, err
	}

	price, _, _, err := b.cashflowSums(y)
	if err != nil {
		return decimal.Decimal{}, err
	}

	// TrySub rather than Sub: at the extremes of the candidate sweep the
	// modelled price and the target can sit far enough apart to overflow, and
	// the scan needs that reported so it can skip the candidate.
	return price.TrySub(b.price)
}

// bisectYield narrows the bracket [lo, hi], across which the price difference
// changes sign, until it is tighter than the tolerance. loDiff is the price
// difference at lo.
func (b Config) bisectYield(lo, hi, loDiff decimal.Decimal) (decimal.Decimal, error) {
	mid := lo.Add(hi).Mul(yieldHalf)

	for i := 0; i < maxYieldBisectIter; i++ {
		mid = lo.Add(hi).Mul(yieldHalf)

		midDiff, err := b.priceDiff(mid)
		if err != nil {
			return decimal.Decimal{}, err
		}

		if midDiff.IsZero() || hi.Sub(lo).Abs().LessThan(yieldBracketTol) {
			return mid, nil
		}

		if midDiff.Sign() == loDiff.Sign() {
			lo, loDiff = mid, midDiff
		} else {
			hi = mid
		}
	}

	return mid, nil
}

// yieldCandidates builds the ordered list of annual yields the search scans for
// a sign change: a fine sweep across −99%…100% then a coarse sweep up to
// 10000%. Values come from exact integer ratios to stay at a representable
// precision.
func yieldCandidates() []decimal.Decimal {
	hundred := decimal.MustFromInt64(100, 0)

	candidates := make([]decimal.Decimal, 0, 300)

	for k := int64(-99); k <= 100; k++ {
		candidates = append(candidates, decimal.MustFromInt64(k, 0).MustDiv(hundred))
	}

	for x := int64(2); x <= 100; x++ {
		candidates = append(candidates, decimal.MustFromInt64(x, 0))
	}

	return candidates
}
