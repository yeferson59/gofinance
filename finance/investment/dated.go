package investment

import (
	"time"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/daycount"
	"github.com/yeferson59/gofinance/v2/money"
)

// DatedCashFlow is a cash flow that occurs on a specific calendar date, used by
// the date-based metrics XNPV and XIRR. As with the period-indexed helpers, an
// outflow is negative and an inflow positive.
type DatedCashFlow struct {
	Date   time.Time
	Amount money.Money
}

// datedFlows validates a dated cash-flow series and reduces it to the amounts
// and their time offsets (in years) from the first flow's date, measured with
// the Actual/365 Fixed convention — the same basis spreadsheet XNPV/XIRR use.
// All dates must be on or after the base date and every amount must share one
// currency.
func datedFlows(flows []DatedCashFlow) (amounts, times []decimal.Decimal, currency money.Currency, err error) {
	if len(flows) == 0 {
		return nil, nil, 0, ErrNoCashFlows
	}

	base := flows[0].Date
	currency = flows[0].Amount.GetCurrency()

	amounts = make([]decimal.Decimal, len(flows))
	times = make([]decimal.Decimal, len(flows))

	for i, flow := range flows {
		if flow.Amount.GetCurrency() != currency {
			return nil, nil, 0, money.ErrCurrencyMismatch
		}

		if flow.Date.Before(base) {
			return nil, nil, 0, ErrDatesBeforeBase
		}

		yf, err := daycount.YearFraction(base, flow.Date, daycount.Actual365Fixed)
		if err != nil {
			return nil, nil, 0, err
		}

		amounts[i] = flow.Amount.GetDecimal()
		times[i] = yf
	}

	return amounts, times, currency, nil
}
