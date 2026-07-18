package annuities

import (
	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/money"
)

// InvestmentSchedule is one row of a period-by-period investment growth
// table: the balance at that period, how it changed from the previous
// period (in absolute and percentage terms), and the contributions/interest
// accumulated so far.
type InvestmentSchedule struct {
	Period           decimal.Decimal
	Balance          money.Money
	Contribution     money.Money
	Change           money.Money
	ChangePercent    decimal.Decimal
	SumContributions money.Money
	SumInterest      money.Money
}

// BuildInvestmentSchedule generates a period-by-period table showing how an
// initial principal grows under compound interest while also receiving a
// fixed contribution at the end of every period (ordinary annuity):
// Balance[p] = Balance[p-1] × (1+rate) + contribution.
//
// principal may be zero (starting from nothing, growing only from
// contributions). ChangePercent is left as zero for any period whose
// previous balance was zero, since a percentage change from a zero base is
// undefined.
//
// This is meant for feeding charts/statistics — e.g. it shows how a fixed
// contribution's relative weight (ChangePercent) shrinks over time as the
// compounding balance grows past it. Use Annuity.FutureWithContributions for
// a single final value.
//
// It returns ErrInvalidPeriods if nper isn't a positive whole number.
func BuildInvestmentSchedule(principal, contribution money.Money, rate decimal.Decimal, nper decimal.Decimal) ([]InvestmentSchedule, error) {
	return buildInvestmentSchedule(principal, contribution, rate, nper, false)
}

// BuildAnticipateInvestmentSchedule is like BuildInvestmentSchedule, but
// assumes each contribution is made at the beginning of its period (annuity
// due) instead of the end, so it also earns interest during its own first
// period: Balance[p] = (Balance[p-1] + contribution) × (1+rate).
func BuildAnticipateInvestmentSchedule(principal, contribution money.Money, rate decimal.Decimal, nper decimal.Decimal) ([]InvestmentSchedule, error) {
	return buildInvestmentSchedule(principal, contribution, rate, nper, true)
}

func buildInvestmentSchedule(principal, contribution money.Money, rate decimal.Decimal, nper decimal.Decimal, anticipated bool) ([]InvestmentSchedule, error) {
	if principal.Currency() != contribution.Currency() {
		return nil, money.ErrCurrencyMismatch
	}

	until, err := nper.Int64()
	if err != nil {
		return nil, err
	}

	if until <= 0 {
		return nil, ErrInvalidPeriods
	}

	currency := principal.Currency()
	zero := money.MustMoneyFromFloat64(0, currency)

	balance := principal
	sumContributions, sumInterest := zero, zero
	rows := make([]InvestmentSchedule, 0, until+1)

	rows = append(rows, InvestmentSchedule{
		Period:           decimal.Zero,
		Balance:          principal,
		Contribution:     zero,
		Change:           zero,
		ChangePercent:    decimal.Zero,
		SumContributions: zero,
		SumInterest:      zero,
	})

	for p := 1; p <= int(until); p++ {
		previous := balance
		var interest money.Money

		if anticipated {
			balance = balance.Add(contribution)
			interest = balance.MulDecimal(rate)
			balance = balance.Add(interest)
		} else {
			interest = balance.MulDecimal(rate)
			balance = balance.Add(interest).Add(contribution)
		}

		sumContributions = sumContributions.Add(contribution)
		sumInterest = sumInterest.Add(interest)

		changePercent := decimal.Zero
		if !previous.IsZero() {
			changePercent, err = balance.Sub(previous).ToDecimal().Div(previous.ToDecimal())
			if err != nil {
				return nil, err
			}
		}

		rows = append(rows, InvestmentSchedule{
			Period:           decimal.MustFromInt64(int64(p), 0),
			Balance:          balance,
			Contribution:     contribution,
			Change:           balance.Sub(previous),
			ChangePercent:    changePercent,
			SumContributions: sumContributions,
			SumInterest:      sumInterest,
		})
	}

	return rows, nil
}
