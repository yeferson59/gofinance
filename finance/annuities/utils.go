package annuities

import (
	"bufio"
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

func BuildSchedule(pv money.Money, rate money.Decimal, payment money.Money, nper money.Decimal) []Schedule {
	until, err := nper.Int64()
	if err != nil {
		return nil
	}

	balance, rows := pv, make([]Schedule, 0, until)
	sumInterest := money.MoneyZero

	rows = append(rows, Schedule{
		Period:      money.Zero,
		Balance:     pv,
		Payment:     money.MoneyZero,
		Interest:    money.MoneyZero,
		SumInterest: money.MoneyZero,
		Principal:   money.MoneyZero,
	})

	for p := 1; p <= int(until); p++ {
		interest := balance.Mul(rate.ToMoney(balance.Currency()))
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

	bw := bufio.NewWriterSize(f, 65536)
	defer bw.Flush()

	w := csv.NewWriter(bw)
	defer w.Flush()

	if err := w.Write(headers); err != nil {
		return err
	}

	rec := make([]string, 6)
	for _, r := range rows {
		rec[0] = r.Period.String()
		rec[1] = r.Balance.RoundBankString(2)
		rec[2] = r.Payment.RoundBankString(2)
		rec[3] = r.Interest.RoundBankString(2)
		rec[4] = r.SumInterest.RoundBankString(2)
		rec[5] = r.Principal.RoundBankString(2)

		if err := w.Write(rec); err != nil {
			return err
		}
	}

	return nil
}
