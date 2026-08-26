package loans

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

const maxAPRBisectIter = 200

var (
	aprBracketTol = decimal.MustFromString("0.0000000001") // 1e-10
	aprHalf       = decimal.MustFromFloat64(0.5)
)

// Payment returns the level amount due at the end of each period that
// amortizes the principal over the full term at the note rate:
//
//	PMT = PV · i / (1 − (1+i)^−n)
//
// Fees and extra payments do not affect it — they change what the loan costs,
// not what the contract schedules. The result carries the principal's
// currency.
//
// It returns ErrInvalidFrequency, ErrInvalidPeriods, ErrNonPositivePrincipal,
// or ErrInvalidRate on invalid terms.
func (l Config) Payment() (money.Money, error) {
	rate, n, err := l.terms()
	if err != nil {
		return money.Money{}, err
	}

	return levelPayment(l.principal, rate, n)
}

// MustPayment is like Payment but panics on error.
func (l Config) MustPayment() money.Money {
	m, err := l.Payment()
	if err != nil {
		panic(err)
	}

	return m
}

// NetProceeds returns the cash the borrower actually receives: the principal
// minus the up-front fees. It is the amount the APR discounts the payment
// stream back to.
//
// It returns ErrNegativeFees if the fees are negative and
// ErrFeesExceedPrincipal if they leave nothing to finance.
func (l Config) NetProceeds() (money.Money, error) {
	if l.fees.IsNeg() {
		return money.Money{}, ErrNegativeFees
	}

	net := l.principal.GetDecimal().Sub(l.fees)
	if !net.IsPos() {
		return money.Money{}, ErrFeesExceedPrincipal
	}

	return money.NewFromDecimal(net, l.principal.GetCurrency()), nil
}

// MustNetProceeds is like NetProceeds but panics on error.
func (l Config) MustNetProceeds() money.Money {
	m, err := l.NetProceeds()
	if err != nil {
		panic(err)
	}

	return m
}

// EffectiveAnnualRate returns the effective annual rate of the note rate
// alone, i.e. what the periodic rate compounds to over a year:
//
//	EAR = (1 + i)^k − 1
//
// where k is the number of payments per year. It ignores fees; for the
// fee-inclusive figure use EffectiveAPR.
//
// It returns ErrInvalidFrequency, ErrInvalidPeriods, ErrNonPositivePrincipal,
// or ErrInvalidRate on invalid terms.
func (l Config) EffectiveAnnualRate() (decimal.Decimal, error) {
	rate, _, err := l.terms()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return effectiveAnnual(rate, l)
}

// MustEffectiveAnnualRate is like EffectiveAnnualRate but panics on error.
func (l Config) MustEffectiveAnnualRate() decimal.Decimal {
	d, err := l.EffectiveAnnualRate()
	if err != nil {
		panic(err)
	}

	return d
}

// PeriodicAPR returns the per-period rate that equates the scheduled payments
// to the cash actually received (principal − fees):
//
//	netProceeds = PMT · (1 − (1+i)^−n) / i
//
// With no fees it is exactly the note rate; every fee pushes it higher, since
// the same payments now buy less cash. Because the equation has no closed form
// in the rate, the root is bracketed by scanning for a sign change and then
// bisected.
//
// It returns ErrNegativeFees, ErrFeesExceedPrincipal, the term errors of
// PeriodicRate, and ErrNoConvergence if no rate reproduces the net proceeds.
func (l Config) PeriodicAPR() (decimal.Decimal, error) {
	rate, n, err := l.terms()
	if err != nil {
		return decimal.Decimal{}, err
	}

	net, err := l.NetProceeds()
	if err != nil {
		return decimal.Decimal{}, err
	}

	// Without finance charges the borrower receives the full principal, so the
	// APR is the note rate itself — no need to solve for it.
	if l.fees.IsZero() {
		return rate, nil
	}

	payment, err := levelPayment(l.principal, rate, n)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return solveAPR(payment.GetDecimal(), net.GetDecimal(), n)
}

// MustPeriodicAPR is like PeriodicAPR but panics on error.
func (l Config) MustPeriodicAPR() decimal.Decimal {
	d, err := l.PeriodicAPR()
	if err != nil {
		panic(err)
	}

	return d
}

// APR returns the annual percentage rate: the periodic APR annualized by
// simple multiplication (i × k), which is the nominal convention lenders quote
// under US Regulation Z and the EU's APRC-equivalent nominal figure.
//
// It returns the same errors as PeriodicAPR.
func (l Config) APR() (decimal.Decimal, error) {
	periodic, err := l.PeriodicAPR()
	if err != nil {
		return decimal.Decimal{}, err
	}

	periodsPerYear, err := l.periodsPerYear()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return periodic.Mul(periodsPerYear), nil
}

// MustAPR is like APR but panics on error.
func (l Config) MustAPR() decimal.Decimal {
	d, err := l.APR()
	if err != nil {
		panic(err)
	}

	return d
}

// EffectiveAPR returns the fee-inclusive cost of the loan compounded over a
// year:
//
//	(1 + periodicAPR)^k − 1
//
// It is always at least the nominal APR and is the figure to use when
// comparing loans that pay on different schedules.
//
// It returns the same errors as PeriodicAPR.
func (l Config) EffectiveAPR() (decimal.Decimal, error) {
	periodic, err := l.PeriodicAPR()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return effectiveAnnual(periodic, l)
}

// MustEffectiveAPR is like EffectiveAPR but panics on error.
func (l Config) MustEffectiveAPR() decimal.Decimal {
	d, err := l.EffectiveAPR()
	if err != nil {
		panic(err)
	}

	return d
}

// effectiveAnnual compounds a periodic rate over one year of the loan's
// payment frequency.
func effectiveAnnual(rate decimal.Decimal, l Config) (decimal.Decimal, error) {
	periodsPerYear, err := l.periodsPerYear()
	if err != nil {
		return decimal.Decimal{}, err
	}

	growth, err := decimal.One.Add(rate).Pow(periodsPerYear)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return growth.Sub(decimal.One), nil
}

// solveAPR finds the periodic rate at which n payments of pmt are worth
// exactly net today. The present value falls monotonically as the rate rises,
// so the search scans candidate rates for a sign change in
// (PV(rate) − net) and bisects the bracket it finds. Candidates deep enough
// below zero for (1+i)ⁿ to underflow are skipped rather than aborting the
// search: no fee-bearing loan has its APR down there.
func solveAPR(pmt, net decimal.Decimal, n int) (decimal.Decimal, error) {
	var (
		prevRate, prevDiff decimal.Decimal
		bracketed          bool
	)

	for _, rate := range aprCandidates() {
		diff, err := aprResidual(rate, pmt, net, n)
		if err != nil {
			continue
		}

		if diff.IsZero() {
			return rate, nil
		}

		if bracketed && prevDiff.Sign() != diff.Sign() {
			return bisectAPR(prevRate, rate, prevDiff, pmt, net, n)
		}

		prevRate, prevDiff, bracketed = rate, diff, true
	}

	return decimal.Decimal{}, ErrNoConvergence
}

// aprResidual returns the present value of n payments of pmt at the given
// periodic rate, minus the net proceeds. Its root is the periodic APR.
func aprResidual(rate, pmt, net decimal.Decimal, n int) (decimal.Decimal, error) {
	factor, err := annuityFactor(rate, n)
	if err != nil {
		return decimal.Decimal{}, err
	}

	// Try* variants: candidate rates far below the true root can blow the
	// annuity factor past the representable range, and an overflow there is a
	// candidate to skip, not a panic.
	present, err := pmt.TryMul(factor)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return present.TrySub(net)
}

// bisectAPR narrows the bracket [lo, hi], across which the residual changes
// sign, until it is tighter than the tolerance and returns the enclosed rate.
// loDiff is the residual at lo.
func bisectAPR(lo, hi, loDiff, pmt, net decimal.Decimal, n int) (decimal.Decimal, error) {
	mid := lo.Add(hi).Mul(aprHalf)

	for i := 0; i < maxAPRBisectIter; i++ {
		mid = lo.Add(hi).Mul(aprHalf)

		midDiff, err := aprResidual(mid, pmt, net, n)
		if err != nil {
			return decimal.Decimal{}, err
		}

		if midDiff.IsZero() || hi.Sub(lo).Abs().LessThan(aprBracketTol) {
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

// aprCandidates builds the ordered list of periodic rates the search scans for
// a sign change: a point just above −1, a fine sweep across −99%…100% per
// period, then a coarse sweep up to 10000%. Rates come from exact integer
// ratios to stay at a representable precision.
func aprCandidates() []decimal.Decimal {
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
