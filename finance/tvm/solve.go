package tvm

import "github.com/yeferson59/gofinance/v2/decimal"

// factors returns (1+rate)ᴺ and the coefficient that multiplies PMT in the
// TVM equation, (1 + rate·type)·((1+rate)ᴺ − 1)/rate, using the annuity-timing
// setting. At a zero rate the coefficient collapses to its limit, N. It
// returns ErrInvalidRate when 1+rate is not positive.
func (t Config) factors(rate decimal.Decimal) (pow, pmtCoef decimal.Decimal, err error) {
	onePlus := decimal.One.Add(rate)
	if !onePlus.IsPos() {
		return decimal.Decimal{}, decimal.Decimal{}, ErrInvalidRate
	}

	pow, err = onePlus.Pow(t.n)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, err
	}

	typeFactor := decimal.One
	if t.due {
		typeFactor = onePlus
	}

	if rate.IsZero() {
		// limit of ((1+i)ᴺ − 1)/i as i → 0 is N.
		return pow, t.n.Mul(typeFactor), nil
	}

	annuity, err := pow.Sub(decimal.One).Div(rate)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, err
	}

	return pow, typeFactor.Mul(annuity), nil
}

// SolveFV returns the future value implied by the other four variables:
//
//	FV = −(PV·(1+i)ᴺ + PMT·coef)
func (t Config) SolveFV() (decimal.Decimal, error) {
	pow, pmtCoef, err := t.factors(t.rate)
	if err != nil {
		return decimal.Decimal{}, err
	}

	balance, err := t.balance(pow, pmtCoef)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return balance.Neg(), nil
}

// balance returns PV·(1+i)ᴺ + PMT·coef, the part of the TVM equation shared by
// SolveFV and the rate solver's residual.
//
// Every step uses the Try variants: with a large principal or payment the
// products overflow, and a function that returns an error must report that
// rather than panic.
func (t Config) balance(pow, pmtCoef decimal.Decimal) (decimal.Decimal, error) {
	grown, err := t.pv.TryMul(pow)
	if err != nil {
		return decimal.Decimal{}, err
	}

	payments, err := t.pmt.TryMul(pmtCoef)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return grown.TryAdd(payments)
}

// MustSolveFV is like SolveFV but panics on error.
func (t Config) MustSolveFV() decimal.Decimal {
	d, err := t.SolveFV()
	if err != nil {
		panic(err)
	}

	return d
}

// SolvePV returns the present value implied by the other four variables:
//
//	PV = −(FV + PMT·coef) / (1+i)ᴺ
func (t Config) SolvePV() (decimal.Decimal, error) {
	pow, pmtCoef, err := t.factors(t.rate)
	if err != nil {
		return decimal.Decimal{}, err
	}

	payments, err := t.pmt.TryMul(pmtCoef)
	if err != nil {
		return decimal.Decimal{}, err
	}

	numerator, err := t.fv.TryAdd(payments)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return numerator.Neg().Div(pow)
}

// MustSolvePV is like SolvePV but panics on error.
func (t Config) MustSolvePV() decimal.Decimal {
	d, err := t.SolvePV()
	if err != nil {
		panic(err)
	}

	return d
}

// SolvePMT returns the per-period payment implied by the other four variables:
//
//	PMT = −(PV·(1+i)ᴺ + FV) / coef
//
// It returns ErrNoSolution when the payment coefficient is zero (for example
// when N is zero).
func (t Config) SolvePMT() (decimal.Decimal, error) {
	pow, pmtCoef, err := t.factors(t.rate)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if pmtCoef.IsZero() {
		return decimal.Decimal{}, ErrNoSolution
	}

	grown, err := t.pv.TryMul(pow)
	if err != nil {
		return decimal.Decimal{}, err
	}

	numerator, err := grown.TryAdd(t.fv)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return numerator.Neg().Div(pmtCoef)
}

// MustSolvePMT is like SolvePMT but panics on error.
func (t Config) MustSolvePMT() decimal.Decimal {
	d, err := t.SolvePMT()
	if err != nil {
		panic(err)
	}

	return d
}

// SolveN returns the number of periods implied by the other four variables.
// At a zero rate it uses N = −(PV + FV)/PMT; otherwise it solves the equation
// for (1+i)ᴺ and takes logarithms.
//
// It returns ErrInvalidRate if 1+rate is not positive and ErrNoSolution when
// the inputs admit no finite, positive-argument logarithm (for example a zero
// payment at a zero rate, or values that force a non-positive growth factor).
func (t Config) SolveN() (decimal.Decimal, error) {
	if t.rate.IsZero() {
		if t.pmt.IsZero() {
			return decimal.Decimal{}, ErrNoSolution
		}

		return t.pv.Add(t.fv).Neg().Div(t.pmt)
	}

	onePlus := decimal.One.Add(t.rate)
	if !onePlus.IsPos() {
		return decimal.Decimal{}, ErrInvalidRate
	}

	typeFactor := decimal.One
	if t.due {
		typeFactor = onePlus
	}

	// k = PMT·typeFactor / i, so the equation PV·powᴺ + PMT·coef + FV = 0
	// rearranges to (1+i)ᴺ = (k − FV)/(PV + k).
	k, err := t.pmt.Mul(typeFactor).Div(t.rate)
	if err != nil {
		return decimal.Decimal{}, err
	}

	denominator := t.pv.Add(k)
	if denominator.IsZero() {
		return decimal.Decimal{}, ErrNoSolution
	}

	pow, err := k.Sub(t.fv).Div(denominator)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if !pow.IsPos() {
		return decimal.Decimal{}, ErrNoSolution
	}

	lnPow, err := pow.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	lnBase, err := onePlus.Ln()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return lnPow.Div(lnBase)
}

// MustSolveN is like SolveN but panics on error.
func (t Config) MustSolveN() decimal.Decimal {
	d, err := t.SolveN()
	if err != nil {
		panic(err)
	}

	return d
}
