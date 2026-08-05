package tvm

import "github.com/yeferson59/gofinance/v2/decimal"

const maxRateBisectIter = 200

var (
	rateBracketTol = decimal.MustFromString("0.0000000001") // 1e-10
	rateHalf       = decimal.MustFromFloat64(0.5)
)

// residual evaluates the left-hand side of the TVM equation,
// PV·(1+i)ᴺ + PMT·coef + FV, at the given rate. Its root is the rate that
// balances the cash flows.
func (t Config) residual(rate decimal.Decimal) (decimal.Decimal, error) {
	pow, pmtCoef, err := t.factors(rate)
	if err != nil {
		return decimal.Decimal{}, err
	}

	balance, err := t.balance(pow, pmtCoef)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return balance.TryAdd(t.fv)
}

// SolveRate returns the per-period interest rate implied by the other four
// variables. Because the equation has no closed form in the rate, it scans a
// range of candidate rates for a sign change in the residual and then bisects
// to the root, so it works for savings, loan, and mixed cash-flow setups
// without needing a starting guess.
//
// It returns ErrNoConvergence when no rate in the searched range balances the
// equation (which also happens when the cash flows never cross zero).
func (t Config) SolveRate() (decimal.Decimal, error) {
	var (
		prevRate  decimal.Decimal
		prevRes   decimal.Decimal
		bracketed bool
	)

	for _, curRate := range rateCandidates() {
		curRes, err := t.residual(curRate)
		if err != nil {
			// A candidate whose growth factor overflows or underflows says
			// nothing about where the root is, so skip it and start a fresh
			// pair after the gap instead of abandoning the search.
			bracketed = false

			continue
		}

		if curRes.IsZero() {
			return curRate, nil
		}

		if bracketed && prevRes.Sign() != curRes.Sign() {
			return t.bisectRate(prevRate, curRate, prevRes)
		}

		prevRate, prevRes, bracketed = curRate, curRes, true
	}

	return decimal.Decimal{}, ErrNoConvergence
}

// MustSolveRate is like SolveRate but panics on error.
func (t Config) MustSolveRate() decimal.Decimal {
	d, err := t.SolveRate()
	if err != nil {
		panic(err)
	}

	return d
}

// bisectRate narrows the bracket [lo, hi], across which the residual changes
// sign, until it is tighter than the tolerance and returns the enclosed root.
// loRes is the residual at lo.
func (t Config) bisectRate(lo, hi, loRes decimal.Decimal) (decimal.Decimal, error) {
	mid := lo.Add(hi).Mul(rateHalf)

	for i := 0; i < maxRateBisectIter; i++ {
		mid = lo.Add(hi).Mul(rateHalf)

		midRes, err := t.residual(mid)
		if err != nil {
			return decimal.Decimal{}, err
		}

		if midRes.IsZero() || hi.Sub(lo).Abs().LessThan(rateBracketTol) {
			return mid, nil
		}

		if midRes.Sign() == loRes.Sign() {
			lo, loRes = mid, midRes
		} else {
			hi = mid
		}
	}

	return mid, nil
}

// rateCandidates builds the ordered list of trial rates the search scans: a
// point just above −1, a fine sweep across −99%…100% per period, then a coarse
// sweep up to 10000% per period. Rates are built from exact integer ratios to
// keep them at a representable precision.
func rateCandidates() []decimal.Decimal {
	hundred := decimal.MustFromInt64(100, 0)

	candidates := []decimal.Decimal{decimal.MustFromString("-0.9999")}

	for k := int64(-99); k <= 100; k++ {
		candidates = append(candidates, decimal.MustFromInt64(k, 0).MustDiv(hundred))
	}

	for x := int64(2); x <= 100; x++ {
		candidates = append(candidates, decimal.MustFromInt64(x, 0))
	}

	return candidates
}
