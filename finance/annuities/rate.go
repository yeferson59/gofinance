package annuities

import (
	"errors"

	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/money"
)

// ErrRateNotFound is returned when RateWithPresent, RateWithFuture,
// AnticipateRateWithPresent, or AnticipateRateWithFuture cannot bracket a
// periodic rate that reproduces the given payment, present/future value, and
// number of periods.
var ErrRateNotFound = errors.New("annuities: could not find a rate that satisfies the given values")

const rateSolverMaxIterations = 200

var (
	// rateSolverMaxGrowth caps how large (1+rate)^periods (or its
	// reciprocal) is allowed to get while searching for a bracket, well
	// inside Decimal's representable range, leaving ample headroom.
	rateSolverMaxGrowth = money.MustFromFloat64(1e15)
	rateSolverTolerance = money.MustFromFloat64(1e-12)
	rateSolverTwo       = money.MustFromFloat64(2)
	// overflowSentinel stands in for a genuinely overflowing Pow result.
	// Every formula here keeps its base (1+rate) positive, so overflow can
	// only occur in the +infinity direction, making a large positive
	// sentinel directionally correct for bisection.
	overflowSentinel = money.MustFromFloat64(1e18)
)

// rateBounds returns a bracket [lower, upper] symmetric enough to contain any
// realistic periodic rate for the given number of periods, while keeping
// (1+rate)^periods and (1+rate)^-periods within rateSolverMaxGrowth for both
// bounds (so neither the present- nor future-value formulas overflow at
// either end of the bracket, regardless of how large periods is).
func rateBounds(periods money.Decimal) (money.Decimal, money.Decimal, error) {
	reciprocal, err := money.One.Div(periods)
	if err != nil {
		return money.Decimal{}, money.Decimal{}, err
	}

	upperGrowth, err := rateSolverMaxGrowth.Pow(reciprocal)
	if err != nil {
		return money.Decimal{}, money.Decimal{}, err
	}

	lowerGrowth, err := rateSolverMaxGrowth.Pow(reciprocal.Neg())
	if err != nil {
		return money.Decimal{}, money.Decimal{}, err
	}

	return lowerGrowth.Sub(money.One), upperGrowth.Sub(money.One), nil
}

// safePow is Decimal.Pow, except a genuine overflow yields overflowSentinel
// instead of an error, so bisection can keep narrowing the bracket instead
// of aborting when it briefly evaluates an extreme candidate rate.
func safePow(base, exponent money.Decimal) (money.Decimal, error) {
	result, err := base.Pow(exponent)
	if errors.Is(err, decimal.ErrOverflow) {
		return overflowSentinel, nil
	}

	return result, err
}

// sign returns -1, 0, or 1 according to d's sign.
func sign(d money.Decimal) int {
	switch {
	case d.IsPos():
		return 1
	case d.IsNeg():
		return -1
	default:
		return 0
	}
}

// solveRate finds the rate for which pv(rate) equals target, using
// bisection entirely in Decimal arithmetic. pv must be monotonic in rate,
// which holds for the ordinary and due present/future annuity formulas.
func solveRate(target, periods money.Decimal, pv func(rate money.Decimal) (money.Decimal, error)) (money.Decimal, error) {
	lower, upper, err := rateBounds(periods)
	if err != nil {
		return money.Decimal{}, err
	}

	lowValue, err := pv(lower)
	if err != nil {
		return money.Decimal{}, err
	}
	flow := lowValue.Sub(target)

	highValue, err := pv(upper)
	if err != nil {
		return money.Decimal{}, err
	}
	fhigh := highValue.Sub(target)

	if sign(flow)*sign(fhigh) > 0 {
		return money.Decimal{}, ErrRateNotFound
	}

	for range rateSolverMaxIterations {
		mid := lower.Add(upper).MustDiv(rateSolverTwo)

		midValue, err := pv(mid)
		if err != nil {
			return money.Decimal{}, err
		}
		fmid := midValue.Sub(target)

		if fmid.Abs().LessThan(rateSolverTolerance) || upper.Sub(lower).Abs().LessThan(rateSolverTolerance) {
			return mid, nil
		}

		if sign(fmid) == sign(flow) {
			lower, flow = mid, fmid
		} else {
			upper = mid
		}
	}

	return lower.Add(upper).MustDiv(rateSolverTwo), nil
}

// presentValueOrdinary is PV = PMT × [1-(1+i)^-n]/i, the present value of an
// ordinary annuity (payments at the end of each period).
func presentValueOrdinary(payment, rate, periods money.Decimal) (money.Decimal, error) {
	if rate.IsZero() {
		return payment.Mul(periods), nil
	}

	growthFactor := rate.Add(money.One)

	discountFactor, err := safePow(growthFactor, periods.Neg())
	if err != nil {
		return money.Decimal{}, err
	}

	numerator := payment.Mul(money.One.Sub(discountFactor))

	return numerator.Div(rate)
}

// futureValueOrdinary is FV = PMT × [(1+i)^n-1]/i, the future value of an
// ordinary annuity (payments at the end of each period).
func futureValueOrdinary(payment, rate, periods money.Decimal) (money.Decimal, error) {
	if rate.IsZero() {
		return payment.Mul(periods), nil
	}

	growthFactor := rate.Add(money.One)

	growthPower, err := safePow(growthFactor, periods)
	if err != nil {
		return money.Decimal{}, err
	}

	numerator := payment.Mul(growthPower.Sub(money.One))

	return numerator.Div(rate)
}

// presentValueDue is presentValueOrdinary × (1+i): the present value of an
// annuity due (payments at the beginning of each period).
func presentValueDue(payment, rate, periods money.Decimal) (money.Decimal, error) {
	ordinary, err := presentValueOrdinary(payment, rate, periods)
	if err != nil {
		return money.Decimal{}, err
	}

	return ordinary.Mul(rate.Add(money.One)), nil
}

// futureValueDue is futureValueOrdinary × (1+i): the future value of an
// annuity due (payments at the beginning of each period).
func futureValueDue(payment, rate, periods money.Decimal) (money.Decimal, error) {
	ordinary, err := futureValueOrdinary(payment, rate, periods)
	if err != nil {
		return money.Decimal{}, err
	}

	return ordinary.Mul(rate.Add(money.One)), nil
}

// RateWithPresent solves for the periodic interest rate implied by a known
// payment, present value, and number of periods, assuming payments are made
// at the end of each period (ordinary annuity). There is no closed-form
// solution for i in PV = PMT × [1-(1+i)^-n]/i, so it's found numerically.
//
// Returns:
//   - The periodic interest rate as a Decimal
//   - An error if periods/present can't be obtained, or if no rate brackets
//     the target (see ErrRateNotFound)
//
// Example:
//
//	ann, _ := New(500, 10000, 0, period, rate)
//	rate, err := ann.RateWithPresent()
//	// rate is the periodic rate at which $500 payments amortize $10,000
func (a Annuity) RateWithPresent() (money.Decimal, error) {
	periods, _, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return money.Decimal{}, err
	}

	payment := a.value.ToDecimal()
	target := present.ToDecimal()

	return solveRate(target, periods, func(r money.Decimal) (money.Decimal, error) {
		return presentValueOrdinary(payment, r, periods)
	})
}

// RateWithFuture solves for the periodic interest rate implied by a known
// payment, future value, and number of periods, assuming payments are made
// at the end of each period (ordinary annuity). There is no closed-form
// solution for i in FV = PMT × [(1+i)^n-1]/i, so it's found numerically.
//
// Returns:
//   - The periodic interest rate as a Decimal
//   - An error if periods/future can't be obtained, or if no rate brackets
//     the target (see ErrRateNotFound)
//
// Example:
//
//	ann, _ := New(500, 0, 10000, period, rate)
//	rate, err := ann.RateWithFuture()
//	// rate is the periodic rate at which $500 payments accumulate to $10,000
func (a Annuity) RateWithFuture() (money.Decimal, error) {
	periods, _, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return money.Decimal{}, err
	}

	payment := a.value.ToDecimal()
	target := future.ToDecimal()

	return solveRate(target, periods, func(r money.Decimal) (money.Decimal, error) {
		return futureValueOrdinary(payment, r, periods)
	})
}

// AnticipateRateWithPresent is like RateWithPresent, but assumes each payment
// is made at the beginning of its period (annuity due) instead of the end:
// PV = PMT × [1-(1+i)^-n]/i × (1+i).
func (a Annuity) AnticipateRateWithPresent() (money.Decimal, error) {
	periods, _, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	present, err := a.compositeInterest.Present()
	if err != nil {
		return money.Decimal{}, err
	}

	payment := a.value.ToDecimal()
	target := present.ToDecimal()

	return solveRate(target, periods, func(r money.Decimal) (money.Decimal, error) {
		return presentValueDue(payment, r, periods)
	})
}

// AnticipateRateWithFuture is like RateWithFuture, but assumes each payment
// is made at the beginning of its period (annuity due) instead of the end:
// FV = PMT × [(1+i)^n-1]/i × (1+i).
func (a Annuity) AnticipateRateWithFuture() (money.Decimal, error) {
	periods, _, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Decimal{}, err
	}

	future, err := a.compositeInterest.Future()
	if err != nil {
		return money.Decimal{}, err
	}

	payment := a.value.ToDecimal()
	target := future.ToDecimal()

	return solveRate(target, periods, func(r money.Decimal) (money.Decimal, error) {
		return futureValueDue(payment, r, periods)
	})
}
