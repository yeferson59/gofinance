package compositeinterest

import (
	"errors"

	"github.com/yeferson59/gofinance/money"
)

// ErrInvalidPeriods is returned by BuildGrowthSchedule when nper doesn't
// represent a positive whole number of periods.
var ErrInvalidPeriods = errors.New("compositeinterest: number of periods must be positive")

// GrowthSchedule is one row of a period-by-period compound interest growth
// table: the balance at that period, how it changed from the previous
// period (in absolute and percentage terms), and the interest accumulated
// so far.
type GrowthSchedule struct {
	Period        money.Decimal
	Balance       money.Money
	Change        money.Money
	ChangePercent money.Decimal
	SumInterest   money.Money
}

// BuildGrowthSchedule generates a period-by-period table showing how a lump
// sum present grows under compound interest at the periodic rate, with no
// further contributions: Balance[p] = Balance[p-1] × (1+rate).
//
// Change and ChangePercent describe how each period's balance compares to
// the previous one, which for plain compound interest is the period's
// interest in absolute terms (growing every period, since it's earned on an
// ever-larger balance) and the constant configured rate in percentage
// terms. This is meant for feeding charts/statistics; use Future for a
// single final value.
//
// It returns ErrInvalidPeriods if nper isn't a positive whole number.
func BuildGrowthSchedule(present money.Money, rate money.Decimal, nper money.Decimal) ([]GrowthSchedule, error) {
	until, err := nper.Int64()
	if err != nil {
		return nil, err
	}
	if until <= 0 {
		return nil, ErrInvalidPeriods
	}

	currency := present.Currency()
	zero := money.MustMoneyFromFloat64(0, currency)

	balance := present
	sumInterest := zero
	rows := make([]GrowthSchedule, 0, until+1)

	rows = append(rows, GrowthSchedule{
		Period:        money.Zero,
		Balance:       present,
		Change:        zero,
		ChangePercent: money.Zero,
		SumInterest:   zero,
	})

	for p := 1; p <= int(until); p++ {
		previous := balance
		interest := balance.Mul(rate.ToMoney(currency))
		balance = balance.Add(interest)
		sumInterest = sumInterest.Add(interest)

		changePercent, err := interest.ToDecimal().Div(previous.ToDecimal())
		if err != nil {
			return nil, err
		}

		rows = append(rows, GrowthSchedule{
			Period:        money.MustFromInt64(int64(p), 0),
			Balance:       balance,
			Change:        interest,
			ChangePercent: changePercent,
			SumInterest:   sumInterest,
		})
	}

	return rows, nil
}
