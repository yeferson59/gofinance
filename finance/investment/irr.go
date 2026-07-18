package investment

import (
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/money"
)

const (
	maxNewtonIter = 100
	maxBisectIter = 200
)

var (
	irrRateTol    = decimal.MustFromString("0.0000000001") // 1e-10
	irrBracketTol = decimal.MustFromString("0.0000000001") // 1e-10
	irrHalf       = decimal.MustFromFloat64(0.5)
	irrMinusOne   = decimal.One.Neg()
)

// IRR returns the internal rate of return of cashFlows: the periodic rate at
// which the net present value is zero. The returned value is a decimal.Decimal
// fraction per period (e.g. 0.08 for 8% per period).
//
// cashFlows must be non-empty, all in the same currency, and contain at least
// one sign change (an investment followed by returns, or vice versa). A first
// Newton–Raphson pass from a 10% guess handles the common case; if it leaves
// the valid domain or fails to converge, a bracketed bisection search takes
// over.
//
// It returns ErrNoCashFlows for an empty slice, money.ErrCurrencyMismatch on
// mixed currencies, ErrNoSignChange when no sign change is present, and
// ErrNoConvergence if no root can be located.
func IRR(cashFlows []money.Money) (decimal.Decimal, error) {
	amounts, _, err := decimalFlows(cashFlows)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if !hasSignChange(amounts) {
		return decimal.Decimal{}, ErrNoSignChange
	}

	// Newton–Raphson from a 10% guess.
	rate := decimal.MustFromFloat64(0.1)

	for i := 0; i < maxNewtonIter; i++ {
		f, fPrime, err := npvAndDerivative(rate, amounts)
		if err != nil || fPrime.IsZero() {
			break
		}

		step, err := f.Div(fPrime)
		if err != nil {
			break
		}

		next := rate.Sub(step)

		// Keep the iterate inside the (−1, ∞) domain; otherwise hand off to
		// the more robust bisection search.
		if next.LessThanOrEqual(irrMinusOne) {
			break
		}

		if next.Sub(rate).Abs().LessThan(irrRateTol) {
			return next, nil
		}

		rate = next
	}

	return irrBisection(amounts)
}

// MustIRR is like IRR but panics on error.
func MustIRR(cashFlows []money.Money) decimal.Decimal {
	d, err := IRR(cashFlows)
	if err != nil {
		panic(err)
	}

	return d
}

// hasSignChange reports whether the non-zero amounts change sign at least
// once, a necessary condition for an internal rate of return to exist.
func hasSignChange(amounts []decimal.Decimal) bool {
	lastSign := 0

	for _, a := range amounts {
		s := a.Sign()
		if s == 0 {
			continue
		}

		if lastSign != 0 && s != lastSign {
			return true
		}

		lastSign = s
	}

	return false
}

// npvAndDerivative returns the NPV of amounts at rate and its first derivative
// with respect to rate, computed together in a single pass for the Newton
// step.
func npvAndDerivative(rate decimal.Decimal, amounts []decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	onePlus := decimal.One.Add(rate)
	if !onePlus.IsPos() {
		return decimal.Decimal{}, decimal.Decimal{}, ErrInvalidRate
	}

	f := amounts[0]
	fPrime := decimal.Zero
	factor := onePlus // (1+rate)^t for the current t, starting at t = 1

	for t := 1; t < len(amounts); t++ {
		discounted, err := amounts[t].Div(factor)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		f = f.Add(discounted)

		// d/dr [ CFₜ (1+r)^-t ] = −t · CFₜ (1+r)^-(t+1) = −t · discounted / (1+r)
		tDec, err := decimal.NewFromInt64(int64(t), 0)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		dTerm, err := discounted.Mul(tDec).Div(onePlus)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		fPrime = fPrime.Sub(dTerm)
		factor = factor.Mul(onePlus)
	}

	return f, fPrime, nil
}

// irrBisection scans a range of candidate rates for a change in the sign of
// NPV and, once one is bracketed, bisects to locate the root.
func irrBisection(amounts []decimal.Decimal) (decimal.Decimal, error) {
	candidates := irrCandidates()

	prevRate := candidates[0]

	prevNPV, err := npvDecimal(prevRate, amounts)
	if err != nil {
		return decimal.Decimal{}, err
	}

	for i := 1; i < len(candidates); i++ {
		curRate := candidates[i]

		curNPV, err := npvDecimal(curRate, amounts)
		if err != nil {
			continue
		}

		if prevNPV.IsZero() {
			return prevRate, nil
		}

		if curNPV.IsZero() {
			return curRate, nil
		}

		if prevNPV.Sign() != curNPV.Sign() {
			return bisect(prevRate, curRate, prevNPV, amounts)
		}

		prevRate, prevNPV = curRate, curNPV
	}

	return decimal.Decimal{}, ErrNoConvergence
}

// bisect narrows the bracket [lo, hi], on which NPV changes sign, until it is
// tighter than the tolerance, and returns the enclosed root. loNPV is the NPV
// at lo, used to decide which half to keep.
func bisect(lo, hi, loNPV decimal.Decimal, amounts []decimal.Decimal) (decimal.Decimal, error) {
	mid := lo.Add(hi).Mul(irrHalf)

	for i := 0; i < maxBisectIter; i++ {
		mid = lo.Add(hi).Mul(irrHalf)

		midNPV, err := npvDecimal(mid, amounts)
		if err != nil {
			return decimal.Decimal{}, err
		}

		if midNPV.IsZero() || hi.Sub(lo).Abs().LessThan(irrBracketTol) {
			return mid, nil
		}

		if midNPV.Sign() == loNPV.Sign() {
			lo, loNPV = mid, midNPV
		} else {
			hi = mid
		}
	}

	return mid, nil
}

// irrCandidates builds the ordered list of trial rates the bisection search
// scans for a sign change: a point just above −1, a fine sweep across the
// common −99%…100% range, then a coarse sweep up to 10000%. Rates are built
// from exact integer ratios to keep them at a representable precision.
func irrCandidates() []decimal.Decimal {
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
