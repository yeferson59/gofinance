package returns

import "github.com/yeferson59/gofinance/v2/decimal"

// Mean returns the arithmetic mean of a series of per-period returns.
//
// It returns ErrNoReturns for an empty slice.
func Mean(rates []decimal.Decimal) (decimal.Decimal, error) {
	if len(rates) == 0 {
		return decimal.Decimal{}, ErrNoReturns
	}

	sum := decimal.Zero

	for _, rate := range rates {
		var err error

		sum, err = sum.TryAdd(rate)
		if err != nil {
			return decimal.Decimal{}, err
		}
	}

	count, err := decimal.NewFromInt64(int64(len(rates)), 0)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return sum.Div(count)
}

// MustMean is like Mean but panics on error.
func MustMean(rates []decimal.Decimal) decimal.Decimal {
	d, err := Mean(rates)
	if err != nil {
		panic(err)
	}

	return d
}

// Variance returns the sample variance of a series of returns, dividing the
// squared deviations by n−1. This is the estimator to use when the returns are
// a sample of a longer history — the usual case for realized performance.
//
// It returns ErrInsufficientReturns when given fewer than two returns.
func Variance(rates []decimal.Decimal) (decimal.Decimal, error) {
	if len(rates) < 2 {
		return decimal.Decimal{}, ErrInsufficientReturns
	}

	return variance(rates, len(rates)-1)
}

// MustVariance is like Variance but panics on error.
func MustVariance(rates []decimal.Decimal) decimal.Decimal {
	d, err := Variance(rates)
	if err != nil {
		panic(err)
	}

	return d
}

// PopulationVariance returns the population variance of a series of returns,
// dividing the squared deviations by n. Use it when the series is the entire
// population of interest rather than a sample of it.
//
// It returns ErrNoReturns for an empty slice.
func PopulationVariance(rates []decimal.Decimal) (decimal.Decimal, error) {
	if len(rates) == 0 {
		return decimal.Decimal{}, ErrNoReturns
	}

	return variance(rates, len(rates))
}

// MustPopulationVariance is like PopulationVariance but panics on error.
func MustPopulationVariance(rates []decimal.Decimal) decimal.Decimal {
	d, err := PopulationVariance(rates)
	if err != nil {
		panic(err)
	}

	return d
}

// Volatility returns the sample standard deviation of a series of returns: the
// square root of Variance, and the usual measure of how much a return series
// swings around its own average. It is expressed per period, in the same units
// as the returns — annualize it with AnnualizedVolatility.
//
// It returns ErrInsufficientReturns when given fewer than two returns.
func Volatility(rates []decimal.Decimal) (decimal.Decimal, error) {
	v, err := Variance(rates)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return v.Sqrt()
}

// MustVolatility is like Volatility but panics on error.
func MustVolatility(rates []decimal.Decimal) decimal.Decimal {
	d, err := Volatility(rates)
	if err != nil {
		panic(err)
	}

	return d
}

// PopulationVolatility returns the population standard deviation of a series
// of returns: the square root of PopulationVariance.
//
// It returns ErrNoReturns for an empty slice.
func PopulationVolatility(rates []decimal.Decimal) (decimal.Decimal, error) {
	v, err := PopulationVariance(rates)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return v.Sqrt()
}

// MustPopulationVolatility is like PopulationVolatility but panics on error.
func MustPopulationVolatility(rates []decimal.Decimal) decimal.Decimal {
	d, err := PopulationVolatility(rates)
	if err != nil {
		panic(err)
	}

	return d
}

// AnnualizedVolatility scales a per-period volatility to a yearly figure by the
// square-root-of-time rule:
//
//	σannual = σperiod × √k
//
// where k is the number of periods in a year (12 for monthly returns, 252 for
// daily trading returns). The rule assumes returns are serially independent,
// which is the standard convention for quoting volatility.
//
// It returns ErrNonPositivePeriods if periodsPerYear is not positive.
func AnnualizedVolatility(volatility, periodsPerYear decimal.Decimal) (decimal.Decimal, error) {
	scale, err := timeScale(periodsPerYear)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return volatility.TryMul(scale)
}

// MustAnnualizedVolatility is like AnnualizedVolatility but panics on error.
func MustAnnualizedVolatility(volatility, periodsPerYear decimal.Decimal) decimal.Decimal {
	d, err := AnnualizedVolatility(volatility, periodsPerYear)
	if err != nil {
		panic(err)
	}

	return d
}

// SharpeRatio returns the risk-adjusted return of a series: how much return
// each unit of volatility bought, above what a risk-free asset paid.
//
//	Sharpe = (mean(r) − riskFree) / σ(r)
//
// riskFree is the risk-free rate over the same period as the returns (a
// monthly rate for monthly returns) and the volatility is the sample standard
// deviation. The result is per period — use AnnualizedSharpeRatio for the
// figure funds usually quote.
//
// It returns ErrInsufficientReturns when given fewer than two returns and
// ErrZeroVolatility when the returns never vary.
func SharpeRatio(rates []decimal.Decimal, riskFree decimal.Decimal) (decimal.Decimal, error) {
	mean, err := Mean(rates)
	if err != nil {
		return decimal.Decimal{}, err
	}

	volatility, err := Volatility(rates)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if volatility.IsZero() {
		return decimal.Decimal{}, ErrZeroVolatility
	}

	excess, err := mean.TrySub(riskFree)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return excess.Div(volatility)
}

// MustSharpeRatio is like SharpeRatio but panics on error.
func MustSharpeRatio(rates []decimal.Decimal, riskFree decimal.Decimal) decimal.Decimal {
	d, err := SharpeRatio(rates, riskFree)
	if err != nil {
		panic(err)
	}

	return d
}

// AnnualizedSharpeRatio scales a per-period Sharpe ratio to a yearly figure by
// the same square-root-of-time rule as AnnualizedVolatility:
//
//	Sharpeannual = Sharpeperiod × √k
//
// rates and riskFree are both per period; periodsPerYear is how many of those
// periods fit in a year.
//
// It returns the errors of SharpeRatio plus ErrNonPositivePeriods if
// periodsPerYear is not positive.
func AnnualizedSharpeRatio(rates []decimal.Decimal, riskFree, periodsPerYear decimal.Decimal) (decimal.Decimal, error) {
	sharpe, err := SharpeRatio(rates, riskFree)
	if err != nil {
		return decimal.Decimal{}, err
	}

	scale, err := timeScale(periodsPerYear)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return sharpe.TryMul(scale)
}

// MustAnnualizedSharpeRatio is like AnnualizedSharpeRatio but panics on error.
func MustAnnualizedSharpeRatio(rates []decimal.Decimal, riskFree, periodsPerYear decimal.Decimal) decimal.Decimal {
	d, err := AnnualizedSharpeRatio(rates, riskFree, periodsPerYear)
	if err != nil {
		panic(err)
	}

	return d
}

// variance returns the sum of squared deviations from the mean divided by
// denominator, which is n−1 for the sample estimator and n for the population
// one.
func variance(rates []decimal.Decimal, denominator int) (decimal.Decimal, error) {
	mean, err := Mean(rates)
	if err != nil {
		return decimal.Decimal{}, err
	}

	sumSquares := decimal.Zero

	for _, rate := range rates {
		deviation, err := rate.TrySub(mean)
		if err != nil {
			return decimal.Decimal{}, err
		}

		square, err := deviation.TryMul(deviation)
		if err != nil {
			return decimal.Decimal{}, err
		}

		sumSquares, err = sumSquares.TryAdd(square)
		if err != nil {
			return decimal.Decimal{}, err
		}
	}

	divisor, err := decimal.NewFromInt64(int64(denominator), 0)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return sumSquares.Div(divisor)
}

// timeScale returns √periodsPerYear, the factor that converts a per-period
// dispersion measure into an annual one.
func timeScale(periodsPerYear decimal.Decimal) (decimal.Decimal, error) {
	if !periodsPerYear.IsPos() {
		return decimal.Decimal{}, ErrNonPositivePeriods
	}

	return periodsPerYear.Sqrt()
}
