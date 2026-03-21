package annuities

import (
	"encoding/csv"
	"os"

	"github.com/quagmt/udecimal"
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

func (a Annuity) BuildSchedule(pv, rate, payment money.Money, nper int) []Schedule {
	balance, rows, sumInterest := pv, make([]Schedule, 0, nper), udecimal.Zero

	rows = append(rows, Schedule{
		Period:      money.Decimal{Decimal: udecimal.Zero},
		Balance:     balance,
		Payment:     money.Money{Decimal: udecimal.Zero},
		Interest:    money.Money{Decimal: udecimal.Zero},
		SumInterest: money.Money{Decimal: udecimal.Zero},
		Principal:   money.Money{Decimal: udecimal.Zero},
	})

	for p := 1; p <= nper; p++ {
		interest := balance.Mul(rate.Decimal)
		principal := payment.Sub(interest)
		balance, sumInterest = money.Money{Decimal: balance.Sub(principal)}, sumInterest.Add(interest)

		rows = append(rows, Schedule{
			Period:      money.MustFromInt64(int64(p), 0),
			Balance:     balance,
			Payment:     payment,
			Interest:    money.Money{Decimal: interest},
			SumInterest: money.Money{Decimal: sumInterest},
			Principal:   money.Money{Decimal: principal},
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
