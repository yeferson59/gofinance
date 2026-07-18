package annuities

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// ErrInvalidPeriods is returned by BuildSchedule when nper doesn't
// represent a positive whole number of periods.
var ErrInvalidPeriods = errors.New("annuities: number of periods must be positive")

type Schedule struct {
	Period      decimal.Decimal
	Balance     money.Money
	Payment     money.Money
	Interest    money.Money
	SumInterest money.Money
	Principal   money.Money
}

// BuildSchedule generates a period-by-period amortization table for a loan
// with present value pv, periodic rate rate, fixed periodic payment
// payment, and nper total periods.
//
// It returns ErrCurrencyMismatch if pv and payment aren't in the same
// currency, ErrInvalidPeriods if nper isn't a positive whole number, and
// wraps any error from parsing nper as an integer.
func BuildSchedule(pv money.Money, rate decimal.Decimal, payment money.Money, nper decimal.Decimal) ([]Schedule, error) {
	if pv.Currency() != payment.Currency() {
		return nil, money.ErrCurrencyMismatch
	}

	until, err := nper.Int64()
	if err != nil {
		return nil, err
	}
	if until <= 0 {
		return nil, ErrInvalidPeriods
	}

	currency := pv.Currency()
	zero := money.MustMoneyFromFloat64(0, currency)

	balance, rows := pv, make([]Schedule, 0, until+1)
	sumInterest := zero

	rows = append(rows, Schedule{
		Period:      decimal.Zero,
		Balance:     pv,
		Payment:     zero,
		Interest:    zero,
		SumInterest: zero,
		Principal:   zero,
	})

	for p := 1; p <= int(until); p++ {
		interest := balance.MulDecimal(rate)
		principal := payment.Sub(interest)
		balance = balance.Sub(principal)
		sumInterest = sumInterest.Add(interest)

		rows = append(rows, Schedule{
			Period:      decimal.MustFromInt64(int64(p), 0),
			Balance:     balance,
			Payment:     payment,
			Interest:    interest,
			SumInterest: sumInterest,
			Principal:   principal,
		})
	}

	return rows, nil
}

// WriteCSV writes rows to filenamePath as CSV. It is a convenience wrapper
// around WriteCSVTo for the common file-on-disk case.
func WriteCSV(filenamePath string, headers []string, rows []Schedule) (err error) {
	f, err := os.Create(filenamePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()

	return WriteCSVTo(f, headers, rows)
}

// WriteCSVTo writes rows to out as CSV, rounding each monetary column to
// its own currency's standard precision (e.g. 0 decimals for JPY, 3 for
// BHD) rather than assuming two decimal places. Callers choose the
// destination — a file, an HTTP response, a buffer — so the schedule
// domain logic stays free of filesystem concerns.
func WriteCSVTo(out io.Writer, headers []string, rows []Schedule) (err error) {
	bw := bufio.NewWriterSize(out, 65536)
	defer func() {
		if ferr := bw.Flush(); err == nil {
			err = ferr
		}
	}()

	w := csv.NewWriter(bw)
	defer func() {
		w.Flush()
		if werr := w.Error(); err == nil {
			err = werr
		}
	}()

	if err = w.Write(headers); err != nil {
		return err
	}

	rec := make([]string, 6)
	for _, r := range rows {
		prec, perr := r.Balance.Currency().GetCurrencyPrecisionCode()
		if perr != nil {
			return perr
		}

		rec[0] = r.Period.String()
		rec[1] = r.Balance.RoundBankString(prec)
		rec[2] = r.Payment.RoundBankString(prec)
		rec[3] = r.Interest.RoundBankString(prec)
		rec[4] = r.SumInterest.RoundBankString(prec)
		rec[5] = r.Principal.RoundBankString(prec)

		if err = w.Write(rec); err != nil {
			return err
		}
	}

	return err
}
