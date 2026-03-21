package annuities

import (
	"encoding/csv"
	"os"

	"github.com/yeferson59/gofinance/money"
)

type Schedule struct {
	Period      money.Decimal
	Balance     money.Money
	Payment     money.Money
	Interest    money.Money
	SumInterest money.Money
	Principal   money.Money
}

func BuildSchedule(pv, rate, payment money.Money, nper int) []Schedule {
	balance, rows := pv, make([]Schedule, 0, nper)
	sumInterest := money.MoneyZero

	rows = append(rows, Schedule{
		Period:      money.Zero,
		Balance:     pv,
		Payment:     money.MoneyZero,
		Interest:    money.MoneyZero,
		SumInterest: money.MoneyZero,
		Principal:   money.MoneyZero,
	})

	for p := 1; p <= nper; p++ {
		interest := balance.Mul(rate.ToDecimal().ToMoney(balance.Currency()))
		principal := payment.Sub(interest)
		balance = balance.Sub(principal)
		sumInterest = sumInterest.Add(interest)

		rows = append(rows, Schedule{
			Period:      money.MustFromInt64(int64(p), 0),
			Balance:     balance,
			Payment:     payment,
			Interest:    interest,
			SumInterest: sumInterest,
			Principal:   principal,
		})
	}

	return rows
}

func WriteCSV(filenamePath string, headers []string, rows []Schedule) error {
	f, err := os.Create(filenamePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(headers); err != nil {
		return err
	}

	for _, r := range rows {
		rec := []string{
			r.Period.String(),
			r.Balance.RoundBank(2).StringFixed(2),
			r.Payment.RoundBank(2).StringFixed(2),
			r.Interest.RoundBank(2).StringFixed(2),
			r.SumInterest.RoundBank(2).StringFixed(2),
			r.Principal.RoundBank(2).StringFixed(2),
		}

		if err := w.Write(rec); err != nil {
			return err
		}
	}

	return nil
}
