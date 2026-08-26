package returns

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/investment"
	"github.com/yeferson59/gofinance/v2/money"
)

// Subperiod is one valuation interval of a portfolio, bounded by the external
// cash flows that would otherwise distort the measured performance. Split the
// history at every contribution or withdrawal and describe each piece here.
type Subperiod struct {
	// Begin is the portfolio's value at the start of the subperiod, before the
	// flow.
	Begin money.Money

	// Flow is the external cash flow at the start of the subperiod: positive
	// when money is put in, negative when money is taken out. A flow at the
	// end of a subperiod is the same thing as a flow at the start of the next
	// one, which is why subperiods are cut at the flow dates.
	Flow money.Money

	// End is the portfolio's value at the end of the subperiod.
	End money.Money
}

// TimeWeightedReturn returns the time-weighted rate of return over the given
// subperiods: the return the manager earned, with the timing and size of the
// investor's own deposits and withdrawals stripped out. Each subperiod's growth
// is measured against the capital actually at work in it and the results are
// chained:
//
//	TWR = Π [ Endₖ / (Beginₖ + Flowₖ) ] − 1
//
// Every amount must share the same currency, and the invested base of each
// subperiod (Begin + Flow) must be positive. The result is a decimal.Decimal
// fraction over the whole horizon — annualize it with Annualized.
//
// It returns ErrNoSubperiods for an empty slice, money.ErrCurrencyMismatch on
// mixed currencies, and ErrNonPositiveValue when a subperiod starts with no
// capital at work.
func TimeWeightedReturn(subperiods []Subperiod) (decimal.Decimal, error) {
	if len(subperiods) == 0 {
		return decimal.Decimal{}, ErrNoSubperiods
	}

	currency := subperiods[0].Begin.GetCurrency()
	growth := decimal.One

	for _, sub := range subperiods {
		if sub.Begin.GetCurrency() != currency || sub.End.GetCurrency() != currency {
			return decimal.Decimal{}, money.ErrCurrencyMismatch
		}

		invested := sub.Begin

		// A flow of zero carries no currency to check: leaving Flow unset is
		// the natural way to describe a subperiod with no deposit or
		// withdrawal.
		if !sub.Flow.IsZero() {
			if sub.Flow.GetCurrency() != currency {
				return decimal.Decimal{}, money.ErrCurrencyMismatch
			}

			invested = sub.Begin.Add(sub.Flow)
		}

		if !invested.IsPositive() {
			return decimal.Decimal{}, ErrNonPositiveValue
		}

		factor, err := sub.End.GetDecimal().Div(invested.GetDecimal())
		if err != nil {
			return decimal.Decimal{}, err
		}

		growth, err = growth.TryMul(factor)
		if err != nil {
			return decimal.Decimal{}, err
		}
	}

	return growth.Sub(decimal.One), nil
}

// MustTimeWeightedReturn is like TimeWeightedReturn but panics on error.
func MustTimeWeightedReturn(subperiods []Subperiod) decimal.Decimal {
	d, err := TimeWeightedReturn(subperiods)
	if err != nil {
		panic(err)
	}

	return d
}

// ChainReturns links a series of per-period returns into the cumulative return
// over the whole series:
//
//	total = Π (1 + rₖ) − 1
//
// Use it when the per-period returns are already known — a fund's monthly
// performance, say — rather than the valuations behind them. Each rate is a
// fraction and must be greater than −1.
//
// It returns ErrNoReturns for an empty slice and ErrNonPositiveValue if a
// return wipes out more than the whole capital.
func ChainReturns(rates []decimal.Decimal) (decimal.Decimal, error) {
	if len(rates) == 0 {
		return decimal.Decimal{}, ErrNoReturns
	}

	growth := decimal.One

	for _, rate := range rates {
		factor := decimal.One.Add(rate)
		if !factor.IsPos() {
			return decimal.Decimal{}, ErrNonPositiveValue
		}

		var err error

		growth, err = growth.TryMul(factor)
		if err != nil {
			return decimal.Decimal{}, err
		}
	}

	return growth.Sub(decimal.One), nil
}

// MustChainReturns is like ChainReturns but panics on error.
func MustChainReturns(rates []decimal.Decimal) decimal.Decimal {
	d, err := ChainReturns(rates)
	if err != nil {
		panic(err)
	}

	return d
}

// MoneyWeightedReturn returns the money-weighted (dollar-weighted) rate of
// return: the periodic rate at which every cash flow the investor made, plus
// the portfolio's final value, discounts back to zero. It is the internal rate
// of return of the investor's own experience, so unlike TimeWeightedReturn it
// rewards or punishes the timing of the flows.
//
// initial is the amount put in at the start; interimFlows holds the net
// external flow at the end of each period in between — positive when money is
// added, negative when it is withdrawn — and final is the portfolio's value at
// the end. The horizon is therefore len(interimFlows)+1 periods, and the result
// is a rate per period.
//
// It returns money.ErrCurrencyMismatch on mixed currencies and the solver
// errors of investment.IRR (ErrNoSignChange, ErrNoConvergence) when no rate
// balances the flows.
func MoneyWeightedReturn(initial money.Money, interimFlows []money.Money, final money.Money) (decimal.Decimal, error) {
	flows := make([]money.Money, 0, len(interimFlows)+2)
	quiet := money.NewFromDecimal(decimal.Zero, initial.GetCurrency())

	// Money paid into the portfolio is an outflow for the investor; the value
	// recovered at the end is the inflow that has to justify it.
	flows = append(flows, initial.Neg())

	for _, flow := range interimFlows {
		// A period with no flow can be left as the zero Money, which carries
		// no currency of its own.
		if flow.IsZero() {
			flows = append(flows, quiet)
			continue
		}

		flows = append(flows, flow.Neg())
	}

	flows = append(flows, final)

	return investment.IRR(flows)
}

// MustMoneyWeightedReturn is like MoneyWeightedReturn but panics on error.
func MustMoneyWeightedReturn(initial money.Money, interimFlows []money.Money, final money.Money) decimal.Decimal {
	d, err := MoneyWeightedReturn(initial, interimFlows, final)
	if err != nil {
		panic(err)
	}

	return d
}
