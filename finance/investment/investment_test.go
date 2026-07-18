package investment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func usd(amount float64) money.Money {
	return money.MustMoneyFromFloat64(amount, money.USD)
}

func flows(amounts ...float64) []money.Money {
	out := make([]money.Money, len(amounts))
	for i, a := range amounts {
		out[i] = usd(a)
	}

	return out
}

func TestNPV(t *testing.T) {
	// -1000 + 400/1.1 + 400/1.21 + 400/1.331 ≈ -5.2592.
	npv, err := NPV(decimal.MustFromFloat64(0.10), flows(-1000, 400, 400, 400))
	require.NoError(t, err)
	assert.InDelta(t, -5.2592, npv.InexactFloat64(), 1e-3)
	assert.Equal(t, money.USD, npv.Currency())
}

func TestNPVZeroRate(t *testing.T) {
	// At a 0% discount rate NPV is just the sum of the flows.
	npv, err := NPV(decimal.Zero, flows(-1000, 400, 400, 400))
	require.NoError(t, err)
	assert.InDelta(t, 200.0, npv.InexactFloat64(), 1e-9)
}

func TestNPVErrors(t *testing.T) {
	_, err := NPV(decimal.MustFromFloat64(0.1), nil)
	assert.ErrorIs(t, err, ErrNoCashFlows)

	_, err = NPV(decimal.MustFromFloat64(-1), flows(-1000, 400))
	assert.ErrorIs(t, err, ErrInvalidRate)

	mixed := []money.Money{usd(-1000), money.MustMoneyFromFloat64(400, money.EUR)}
	_, err = NPV(decimal.MustFromFloat64(0.1), mixed)
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)
}

func TestIRRExact(t *testing.T) {
	// -100 today, 110 next period → exactly 10%.
	irr, err := IRR(flows(-100, 110))
	require.NoError(t, err)
	assert.InDelta(t, 0.10, irr.InexactFloat64(), 1e-9)
}

func TestIRRNegative(t *testing.T) {
	// -100 today, 90 next period → −10%.
	irr, err := IRR(flows(-100, 90))
	require.NoError(t, err)
	assert.InDelta(t, -0.10, irr.InexactFloat64(), 1e-9)
}

func TestIRRMultiPeriod(t *testing.T) {
	fs := flows(-1000, 400, 400, 400)

	irr, err := IRR(fs)
	require.NoError(t, err)
	assert.InDelta(t, 0.09701, irr.InexactFloat64(), 1e-4)

	// The NPV discounted at the IRR must be essentially zero.
	npv, err := NPV(irr, fs)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, npv.InexactFloat64(), 1e-4)
}

func TestIRRErrors(t *testing.T) {
	_, err := IRR(nil)
	assert.ErrorIs(t, err, ErrNoCashFlows)

	// All inflows, no sign change.
	_, err = IRR(flows(100, 200, 300))
	assert.ErrorIs(t, err, ErrNoSignChange)

	mixed := []money.Money{usd(-1000), money.MustMoneyFromFloat64(400, money.EUR)}
	_, err = IRR(mixed)
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)
}

func TestMustHelpers(t *testing.T) {
	assert.InDelta(t, 0.10, MustIRR(flows(-100, 110)).InexactFloat64(), 1e-9)
	assert.InDelta(t, 200.0, MustNPV(decimal.Zero, flows(-1000, 400, 400, 400)).InexactFloat64(), 1e-9)
	assert.Panics(t, func() { MustIRR(nil) })
}
