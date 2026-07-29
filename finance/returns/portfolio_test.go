package returns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/investment"
	"github.com/yeferson59/gofinance/v2/money"
)

func eur(amount float64) money.Money {
	return money.MustMoneyFromFloat64(amount, money.EUR)
}

// rate builds a return from its exact decimal string, so no float64 rounding
// creeps into the expected values.
func rate(value string) decimal.Decimal {
	return decimal.MustFromString(value)
}

// portfolioHistory is the running example: 100k grows to 110k, then 50k is
// deposited and the 160k at work ends the year at 180k.
func portfolioHistory() []Subperiod {
	return []Subperiod{
		{Begin: usd(100000), End: usd(110000)},
		{Begin: usd(110000), Flow: usd(50000), End: usd(180000)},
	}
}

func TestTimeWeightedReturn(t *testing.T) {
	// 1.10 × (180/160) = 1.2375.
	twr, err := TimeWeightedReturn(portfolioHistory())
	require.NoError(t, err)
	assert.InDelta(t, 0.2375, twr.InexactFloat64(), 1e-12)
}

func TestTimeWeightedReturnIgnoresFlowTiming(t *testing.T) {
	// The same two subperiod returns with a much larger deposit in between
	// give the same time-weighted result: the manager is not credited for the
	// investor's timing.
	bigger := []Subperiod{
		{Begin: usd(100000), End: usd(110000)},
		{Begin: usd(110000), Flow: usd(890000), End: usd(1125000)},
	}

	twr, err := TimeWeightedReturn(bigger)
	require.NoError(t, err)
	assert.InDelta(t, 0.2375, twr.InexactFloat64(), 1e-12)
}

func TestTimeWeightedReturnWithdrawal(t *testing.T) {
	// A withdrawal shrinks the capital at work but not the measured return:
	// 90/100 − 1 = −10%.
	subperiods := []Subperiod{{Begin: usd(150000), Flow: usd(-50000), End: usd(90000)}}

	twr, err := TimeWeightedReturn(subperiods)
	require.NoError(t, err)
	assert.InDelta(t, -0.10, twr.InexactFloat64(), 1e-12)
}

func TestTimeWeightedReturnErrors(t *testing.T) {
	_, err := TimeWeightedReturn(nil)
	assert.ErrorIs(t, err, ErrNoSubperiods)

	_, err = TimeWeightedReturn([]Subperiod{{Begin: usd(1000), End: eur(1100)}})
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	_, err = TimeWeightedReturn([]Subperiod{{Begin: usd(1000), Flow: eur(100), End: usd(1100)}})
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	// A withdrawal that empties the portfolio leaves nothing at work.
	_, err = TimeWeightedReturn([]Subperiod{{Begin: usd(1000), Flow: usd(-1000), End: usd(0)}})
	assert.ErrorIs(t, err, ErrNonPositiveValue)
}

func TestChainReturns(t *testing.T) {
	// The subperiod returns behind portfolioHistory: +10% then +12.5%.
	total, err := ChainReturns([]decimal.Decimal{rate("0.10"), rate("0.125")})
	require.NoError(t, err)
	assert.InDelta(t, 0.2375, total.InexactFloat64(), 1e-12)
}

func TestChainReturnsRecoversFromLoss(t *testing.T) {
	// −50% then +100% is flat, not +50%.
	total, err := ChainReturns([]decimal.Decimal{rate("-0.5"), rate("1")})
	require.NoError(t, err)
	assert.InDelta(t, 0.0, total.InexactFloat64(), 1e-12)
}

func TestChainReturnsErrors(t *testing.T) {
	_, err := ChainReturns(nil)
	assert.ErrorIs(t, err, ErrNoReturns)

	_, err = ChainReturns([]decimal.Decimal{rate("-1.5")})
	assert.ErrorIs(t, err, ErrNonPositiveValue)
}

func TestMoneyWeightedReturn(t *testing.T) {
	// −100000, −50000, +180000 solves 10x² + 5x − 18 = 0 with x = 1+r.
	mwr, err := MoneyWeightedReturn(usd(100000), []money.Money{usd(50000)}, usd(180000))
	require.NoError(t, err)
	assert.InDelta(t, 0.1147344, mwr.InexactFloat64(), 1e-7)
}

func TestMoneyWeightedReturnDiffersFromTimeWeighted(t *testing.T) {
	twr := MustTimeWeightedReturn(portfolioHistory())
	mwr := MustMoneyWeightedReturn(usd(100000), []money.Money{usd(50000)}, usd(180000))

	// The deposit landed before the weaker subperiod, so the investor's own
	// experience trails the manager's record.
	assert.Less(t, mwr.InexactFloat64(), twr.InexactFloat64())
}

func TestMoneyWeightedReturnWithoutFlows(t *testing.T) {
	// With nothing added in between it is just the per-period growth rate:
	// 121/100 over two periods = 10%.
	mwr, err := MoneyWeightedReturn(usd(100), []money.Money{{}}, usd(121))
	require.NoError(t, err)
	assert.InDelta(t, 0.10, mwr.InexactFloat64(), 1e-9)
}

func TestMoneyWeightedReturnErrors(t *testing.T) {
	_, err := MoneyWeightedReturn(usd(100000), []money.Money{eur(50000)}, usd(180000))
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)

	// Money only ever going out never comes back to zero.
	_, err = MoneyWeightedReturn(usd(100), []money.Money{usd(50)}, usd(0))
	assert.ErrorIs(t, err, investment.ErrNoSignChange)
}

func TestPortfolioMustHelpers(t *testing.T) {
	assert.InDelta(t, 0.2375, MustChainReturns([]decimal.Decimal{rate("0.10"), rate("0.125")}).InexactFloat64(), 1e-12)

	assert.Panics(t, func() { MustTimeWeightedReturn(nil) })
	assert.Panics(t, func() { MustChainReturns(nil) })
	assert.Panics(t, func() { MustMoneyWeightedReturn(usd(100), []money.Money{usd(50)}, usd(0)) })
}
