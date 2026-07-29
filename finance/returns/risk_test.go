package returns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
)

// monthlyReturns is the running series of these tests: five monthly returns
// averaging 1.2% with a sample standard deviation of ≈1.5248%.
func monthlyReturns() []decimal.Decimal {
	return []decimal.Decimal{
		rate("0.02"), rate("-0.01"), rate("0.03"), rate("0.005"), rate("0.015"),
	}
}

func TestMean(t *testing.T) {
	mean, err := Mean(monthlyReturns())
	require.NoError(t, err)
	assert.InDelta(t, 0.012, mean.InexactFloat64(), 1e-15)
}

func TestMeanEmpty(t *testing.T) {
	_, err := Mean(nil)
	assert.ErrorIs(t, err, ErrNoReturns)
}

func TestVariance(t *testing.T) {
	// Σ(r − r̄)² = 0.00093, divided by n−1 = 4.
	sample, err := Variance(monthlyReturns())
	require.NoError(t, err)
	assert.InDelta(t, 0.0002325, sample.InexactFloat64(), 1e-15)

	population, err := PopulationVariance(monthlyReturns())
	require.NoError(t, err)
	assert.InDelta(t, 0.000186, population.InexactFloat64(), 1e-15)

	// Dividing by n instead of n−1 always understates the dispersion.
	assert.Less(t, population.InexactFloat64(), sample.InexactFloat64())
}

func TestVarianceOfAConstantSeriesIsZero(t *testing.T) {
	flat := []decimal.Decimal{rate("0.01"), rate("0.01"), rate("0.01")}

	v, err := Variance(flat)
	require.NoError(t, err)
	assert.True(t, v.IsZero())
}

func TestVarianceErrors(t *testing.T) {
	_, err := Variance([]decimal.Decimal{rate("0.01")})
	assert.ErrorIs(t, err, ErrInsufficientReturns)

	_, err = PopulationVariance(nil)
	assert.ErrorIs(t, err, ErrNoReturns)
}

func TestVolatility(t *testing.T) {
	volatility, err := Volatility(monthlyReturns())
	require.NoError(t, err)
	assert.InDelta(t, 0.0152479506819769, volatility.InexactFloat64(), 1e-15)

	population, err := PopulationVolatility(monthlyReturns())
	require.NoError(t, err)
	assert.InDelta(t, 0.0136381816969859, population.InexactFloat64(), 1e-15)
}

func TestVolatilityIsTheSquareRootOfVariance(t *testing.T) {
	volatility := MustVolatility(monthlyReturns())
	squared := volatility.Mul(volatility)

	assert.InDelta(t, MustVariance(monthlyReturns()).InexactFloat64(), squared.InexactFloat64(), 1e-15)
}

func TestAnnualizedVolatility(t *testing.T) {
	// Monthly volatility scaled by √12.
	annual, err := AnnualizedVolatility(MustVolatility(monthlyReturns()), rate("12"))
	require.NoError(t, err)
	assert.InDelta(t, 0.0528204505849770, annual.InexactFloat64(), 1e-15)
}

func TestAnnualizedVolatilityErrors(t *testing.T) {
	_, err := AnnualizedVolatility(rate("0.05"), decimal.Zero)
	assert.ErrorIs(t, err, ErrNonPositivePeriods)
}

func TestSharpeRatio(t *testing.T) {
	// (0.012 − 0.002) / 0.01524795.
	sharpe, err := SharpeRatio(monthlyReturns(), rate("0.002"))
	require.NoError(t, err)
	assert.InDelta(t, 0.655825835783953, sharpe.InexactFloat64(), 1e-14)
}

func TestSharpeRatioIsNegativeBelowTheRiskFreeRate(t *testing.T) {
	sharpe, err := SharpeRatio(monthlyReturns(), rate("0.05"))
	require.NoError(t, err)
	assert.Negative(t, sharpe.InexactFloat64())
}

func TestAnnualizedSharpeRatio(t *testing.T) {
	annual, err := AnnualizedSharpeRatio(monthlyReturns(), rate("0.002"), rate("12"))
	require.NoError(t, err)
	assert.InDelta(t, 2.27184733698826, annual.InexactFloat64(), 1e-13)
}

func TestSharpeRatioErrors(t *testing.T) {
	_, err := SharpeRatio([]decimal.Decimal{rate("0.01")}, decimal.Zero)
	assert.ErrorIs(t, err, ErrInsufficientReturns)

	// A series that never moves offers no risk to be paid for.
	flat := []decimal.Decimal{rate("0.01"), rate("0.01")}

	_, err = SharpeRatio(flat, decimal.Zero)
	assert.ErrorIs(t, err, ErrZeroVolatility)

	_, err = AnnualizedSharpeRatio(monthlyReturns(), decimal.Zero, decimal.Zero)
	assert.ErrorIs(t, err, ErrNonPositivePeriods)
}

func TestRiskMustHelpers(t *testing.T) {
	rates := monthlyReturns()

	assert.InDelta(t, 0.012, MustMean(rates).InexactFloat64(), 1e-15)
	assert.InDelta(t, 0.0002325, MustVariance(rates).InexactFloat64(), 1e-15)
	assert.InDelta(t, 0.000186, MustPopulationVariance(rates).InexactFloat64(), 1e-15)
	assert.InDelta(t, 0.0136381816969859, MustPopulationVolatility(rates).InexactFloat64(), 1e-15)
	assert.InDelta(t, 0.0528204505849770, MustAnnualizedVolatility(MustVolatility(rates), rate("12")).InexactFloat64(), 1e-15)
	assert.InDelta(t, 0.655825835783953, MustSharpeRatio(rates, rate("0.002")).InexactFloat64(), 1e-14)
	assert.InDelta(t, 2.27184733698826, MustAnnualizedSharpeRatio(rates, rate("0.002"), rate("12")).InexactFloat64(), 1e-13)

	assert.Panics(t, func() { MustMean(nil) })
	assert.Panics(t, func() { MustVariance(nil) })
	assert.Panics(t, func() { MustPopulationVariance(nil) })
	assert.Panics(t, func() { MustVolatility(nil) })
	assert.Panics(t, func() { MustPopulationVolatility(nil) })
	assert.Panics(t, func() { MustAnnualizedVolatility(rate("0.05"), decimal.Zero) })
	assert.Panics(t, func() { MustSharpeRatio(nil, decimal.Zero) })
	assert.Panics(t, func() { MustAnnualizedSharpeRatio(nil, decimal.Zero, rate("12")) })
}
