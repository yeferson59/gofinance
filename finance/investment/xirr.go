package investment

import "github.com/yeferson59/gofinance/money"

// XIRR returns the internal rate of return of cash flows that occur on specific
// dates: the annual rate at which their date-based net present value (XNPV) is
// zero. The returned value is a money.Decimal annual fraction (e.g. 0.12 for
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
func XIRR(flows []DatedCashFlow) (money.Decimal, error) {
	amounts, times, _, err := datedFlows(flows)
	if err != nil {
		return money.Decimal{}, err
	}

	if !hasSignChange(amounts) {
		return money.Decimal{}, ErrNoSignChange
	}

	rate := money.MustFromFloat64(0.1)

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
func MustXIRR(flows []DatedCashFlow) money.Decimal {
	d, err := XIRR(flows)
	if err != nil {
		panic(err)
	}

	return d
}

// xnpvAndDerivative returns the XNPV of the amounts at rate and its first
// derivative with respect to rate, computed together for the Newton step.
func xnpvAndDerivative(rate money.Decimal, amounts, times []money.Decimal) (money.Decimal, money.Decimal, error) {
	onePlus := money.One.Add(rate)
	if !onePlus.IsPos() {
		return money.Decimal{}, money.Decimal{}, ErrInvalidRate
	}

	f := money.Zero
	fPrime := money.Zero

	for i, amount := range amounts {
		if times[i].IsZero() {
			f = f.Add(amount)
			continue
		}

		factor, err := onePlus.Pow(times[i])
		if err != nil {
			return money.Decimal{}, money.Decimal{}, err
		}

		discounted, err := amount.Div(factor)
		if err != nil {
			return money.Decimal{}, money.Decimal{}, err
		}

		f = f.Add(discounted)

		// d/dr [ a·(1+r)^-t ] = −t·a·(1+r)^-(t+1) = −t·discounted / (1+r)
		dTerm, err := times[i].Mul(discounted).Div(onePlus)
		if err != nil {
			return money.Decimal{}, money.Decimal{}, err
		}

		fPrime = fPrime.Sub(dTerm)
	}

	return f, fPrime, nil
}

// xirrBisection scans the candidate annual rates for a sign change in XNPV and
// bisects to the root once one is bracketed.
func xirrBisection(amounts, times []money.Decimal) (money.Decimal, error) {
	candidates := irrCandidates()

	prevRate := candidates[0]

	prevNPV, err := xnpvDecimal(prevRate, amounts, times)
	if err != nil {
		return money.Decimal{}, err
	}

	for i := 1; i < len(candidates); i++ {
		curRate := candidates[i]

		curNPV, err := xnpvDecimal(curRate, amounts, times)
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
			return xbisect(prevRate, curRate, prevNPV, amounts, times)
		}

		prevRate, prevNPV = curRate, curNPV
	}

	return money.Decimal{}, ErrNoConvergence
}

// xbisect narrows the bracket [lo, hi], across which XNPV changes sign, to the
// root. loNPV is the XNPV at lo.
func xbisect(lo, hi, loNPV money.Decimal, amounts, times []money.Decimal) (money.Decimal, error) {
	mid := lo.Add(hi).Mul(irrHalf)

	for i := 0; i < maxBisectIter; i++ {
		mid = lo.Add(hi).Mul(irrHalf)

		midNPV, err := xnpvDecimal(mid, amounts, times)
		if err != nil {
			return money.Decimal{}, err
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
