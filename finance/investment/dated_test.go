package investment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestXNPVExactYear(t *testing.T) {
	// One year apart (2009 is a common year → 365/365 = 1.0):
	// -1000 + 1100/1.1 = 0.
	fs := []DatedCashFlow{
		{Date: day(2009, 1, 1), Amount: usd(-1000)},
		{Date: day(2010, 1, 1), Amount: usd(1100)},
	}

	npv, err := XNPV(decimal.MustFromFloat64(0.10), fs)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, npv.InexactFloat64(), 1e-6)
	assert.Equal(t, money.USD, npv.Currency())
}

func TestXNPVErrors(t *testing.T) {
	_, err := XNPV(decimal.MustFromFloat64(0.1), nil)
	assert.ErrorIs(t, err, ErrNoCashFlows)

	_, err = XNPV(decimal.MustFromFloat64(-1), []DatedCashFlow{{Date: day(2020, 1, 1), Amount: usd(10)}})
	assert.ErrorIs(t, err, ErrInvalidRate)

	before := []DatedCashFlow{
		{Date: day(2010, 1, 1), Amount: usd(-1000)},
		{Date: day(2009, 1, 1), Amount: usd(1100)},
	}
	_, err = XNPV(decimal.MustFromFloat64(0.1), before)
	assert.ErrorIs(t, err, ErrDatesBeforeBase)

	mixed := []DatedCashFlow{
		{Date: day(2010, 1, 1), Amount: usd(-1000)},
		{Date: day(2011, 1, 1), Amount: money.MustMoneyFromFloat64(1100, money.EUR)},
	}
	_, err = XNPV(decimal.MustFromFloat64(0.1), mixed)
	assert.ErrorIs(t, err, money.ErrCurrencyMismatch)
}

func TestXIRRExcelExample(t *testing.T) {
	// The canonical spreadsheet XIRR example → ≈ 0.373363.
	fs := []DatedCashFlow{
		{Date: day(2008, 1, 1), Amount: usd(-10000)},
		{Date: day(2008, 3, 1), Amount: usd(2750)},
		{Date: day(2008, 10, 30), Amount: usd(4250)},
		{Date: day(2009, 2, 15), Amount: usd(3250)},
		{Date: day(2009, 4, 1), Amount: usd(2750)},
	}

	irr, err := XIRR(fs)
	require.NoError(t, err)
	assert.InDelta(t, 0.373363, irr.InexactFloat64(), 1e-4)

	// XNPV discounted at the XIRR must be essentially zero.
	npv, err := XNPV(irr, fs)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, npv.InexactFloat64(), 1e-3)
}

func TestXIRRErrors(t *testing.T) {
	_, err := XIRR(nil)
	assert.ErrorIs(t, err, ErrNoCashFlows)

	noChange := []DatedCashFlow{
		{Date: day(2020, 1, 1), Amount: usd(100)},
		{Date: day(2021, 1, 1), Amount: usd(200)},
	}
	_, err = XIRR(noChange)
	assert.ErrorIs(t, err, ErrNoSignChange)
}

func TestMustXHelpers(t *testing.T) {
	fs := []DatedCashFlow{
		{Date: day(2009, 1, 1), Amount: usd(-1000)},
		{Date: day(2010, 1, 1), Amount: usd(1100)},
	}
	assert.InDelta(t, 0.10, MustXIRR(fs).InexactFloat64(), 1e-6)
	assert.Panics(t, func() { MustXNPV(decimal.MustFromFloat64(0.1), nil) })
}
