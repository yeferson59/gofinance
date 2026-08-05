package investment

import "github.com/yeferson59/gofinance/v2/decimal"

// XIRR returns the internal rate of return of cash flows that occur on specific
// dates: the annual rate at which their date-based net present value (XNPV) is
// zero. The returned value is a decimal.Decimal annual fraction (e.g. 0.12 for
// 12%/yr).
//
// flows must be non-empty, all in the same currency, dated on or after the
// first flow, and contain at least one sign change. A Newton–Raphson pass from
// a 10% guess handles the common case, with a bracketed bisection fallback.
//
// It returns ErrNoCashFlows for an empty slice, money.ErrCurrencyMismatch on
// mixed currencies, ErrDatesBeforeBase if a date precedes the first,
// ErrNoSignChange when no sign change is present, and ErrNoConvergence if no
// root can be located.
func XIRR(flows []DatedCashFlow) (decimal.Decimal, error) {
	amounts, times, _, err := datedFlows(flows)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if !hasSignChange(amounts) {
		return decimal.Decimal{}, ErrNoSignChange
	}

	rate := decimal.MustFromFloat64(0.1)

	for i := 0; i < maxNewtonIter; i++ {
		f, fPrime, err := xnpvAndDerivative(rate, amounts, times)
		if err != nil || fPrime.IsZero() {
			break
		}

		step, err := f.Div(fPrime)
		if err != nil {
			break
		}

		next := rate.Sub(step)
		if next.LessThanOrEqual(irrMinusOne) {
			break
		}

		if next.Sub(rate).Abs().LessThan(irrRateTol) {
			return next, nil
		}

		rate = next
	}

	return xirrBisection(amounts, times)
}

// MustXIRR is like XIRR but panics on error.
func MustXIRR(flows []DatedCashFlow) decimal.Decimal {
	d, err := XIRR(flows)
	if err != nil {
		panic(err)
	}

	return d
}

// xnpvAndDerivative returns the XNPV of the amounts at rate and its first
// derivative with respect to rate, computed together for the Newton step.
func xnpvAndDerivative(rate decimal.Decimal, amounts, times []decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	onePlus := decimal.One.Add(rate)
	if !onePlus.IsPos() {
		return decimal.Decimal{}, decimal.Decimal{}, ErrInvalidRate
	}

	f := decimal.Zero
	fPrime := decimal.Zero

	for i, amount := range amounts {
		discounted, err := discountToBase(amount, onePlus, times[i])
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		f, err = f.TryAdd(discounted)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		// A flow at the base date is not discounted, so it does not move with
		// the rate and contributes nothing to the derivative.
		if times[i].IsZero() {
			continue
		}

		// d/dr [ a·(1+r)^-t ] = −t·a·(1+r)^-(t+1) = −t·discounted / (1+r)
		scaled, err := times[i].TryMul(discounted)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		dTerm, err := scaled.Div(onePlus)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		fPrime, err = fPrime.TrySub(dTerm)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}
	}

	return f, fPrime, nil
}

// xirrBisection scans the candidate annual rates for a sign change in XNPV and
// bisects to the root once one is bracketed.
//
// Like irrBisection, a candidate whose discount factors overflow or underflow
// is skipped rather than aborting the scan: it says nothing about where the
// root is.
func xirrBisection(amounts, times []decimal.Decimal) (decimal.Decimal, error) {
	var (
		prevRate  decimal.Decimal
		prevNPV   decimal.Decimal
		bracketed bool
	)

	for _, curRate := range irrCandidates() {
		curNPV, err := xnpvDecimal(curRate, amounts, times)
		if err != nil {
			bracketed = false

			continue
		}

		if curNPV.IsZero() {
			return curRate, nil
		}

		if bracketed && prevNPV.Sign() != curNPV.Sign() {
			return xbisect(prevRate, curRate, prevNPV, amounts, times)
		}

		prevRate, prevNPV, bracketed = curRate, curNPV, true
	}

	return decimal.Decimal{}, ErrNoConvergence
}

// xbisect narrows the bracket [lo, hi], across which XNPV changes sign, to the
// root. loNPV is the XNPV at lo.
func xbisect(lo, hi, loNPV decimal.Decimal, amounts, times []decimal.Decimal) (decimal.Decimal, error) {
	mid := lo.Add(hi).Mul(irrHalf)

	for i := 0; i < maxBisectIter; i++ {
		mid = lo.Add(hi).Mul(irrHalf)

		midNPV, err := xnpvDecimal(mid, amounts, times)
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
